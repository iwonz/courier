// Package webhook implements Courier's intentionally small outgoing webhook
// profile: one bounded multipart POST containing one file field.
package webhook

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/progress"
	"golang.org/x/time/rate"
)

const (
	streamBufferSize = 128 * 1024
	maximumRateBurst = 128 * 1024
)

// HTTPDoer is the one-attempt boundary used by Sender.
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// Request describes one fully preflighted outgoing webhook.
type Request struct {
	URL        string
	SourceFS   fsx.Backend
	SourcePath string
	Name       string
	Username   string
	Password   []byte
	Rate       int64
	Unlimited  bool
	Progress   progress.Sink
}

// Result reports payload bytes consumed by the HTTP transport and the known
// response status. Bytes exclude multipart framing.
type Result struct {
	StatusCode int
	Bytes      int64
	Elapsed    time.Duration
}

// RejectedError is a received, non-2xx HTTP response.
type RejectedError struct{ StatusCode int }

func (err *RejectedError) Error() string {
	return fmt.Sprintf("webhook receiver returned HTTP %d", err.StatusCode)
}

// UnknownOutcomeError means the transport consumed payload bytes but no
// response arrived, so receiver-side acceptance cannot be determined safely.
type UnknownOutcomeError struct {
	Sent  int64
	Cause error
}

func (err *UnknownOutcomeError) Error() string {
	return fmt.Sprintf("webhook outcome is unknown after sending %d bytes: %v", err.Sent, err.Cause)
}

func (err *UnknownOutcomeError) Unwrap() error { return err.Cause }

// Sender performs exactly one HTTP request. The configured client must not
// implement retries; DefaultClient also refuses redirects.
type Sender struct{ Client HTTPDoer }

// DefaultClient returns Courier's native, TLS-verifying, no-redirect client.
func DefaultClient() *http.Client {
	dialer := &net.Dialer{Timeout: 10 * time.Second, KeepAlive: 30 * time.Second}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment, DialContext: dialer.DialContext,
		ForceAttemptHTTP2: true, DisableKeepAlives: true, TLSHandshakeTimeout: 10 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second, ExpectContinueTimeout: time.Second,
		MaxResponseHeaderBytes: 1 << 20,
	}
	return &http.Client{
		Transport:     transport,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
}

// Send validates and streams one request without redirect or retry logic.
func (sender Sender) Send(ctx context.Context, input Request) (result Result, resultErr error) {
	if err := validateRequest(input); err != nil {
		return Result{}, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, input.URL, http.NoBody)
	if err != nil {
		return Result{}, err
	}
	if request.URL.Host == "" || request.URL.User != nil || request.URL.Fragment != "" || request.URL.Scheme != "http" && request.URL.Scheme != "https" {
		return Result{}, errors.New("outgoing webhook requires an HTTP(S) URL without credentials or fragment")
	}
	info, err := input.SourceFS.Lstat(input.SourcePath)
	if err != nil {
		return Result{}, err
	}
	if !info.Mode().IsRegular() {
		return Result{}, errors.New("outgoing webhook source must be a regular file")
	}
	reader, err := input.SourceFS.Open(input.SourcePath)
	if err != nil {
		return Result{}, err
	}
	defer func() { resultErr = errors.Join(resultErr, reader.Close()) }()

	started := time.Now()
	tracker := progress.New(info.Size(), nil, input.Progress)
	tracker.Stage(progress.StageTransfer)
	pipeReader, pipeWriter := io.Pipe()
	multipartWriter := multipart.NewWriter(pipeWriter)
	contentType := multipartWriter.FormDataContentType()
	var consumed atomic.Int64
	writeResult := make(chan error, 1)
	go func() {
		writeErr := writeMultipart(ctx, multipartWriter, reader, input.Name, info.Size(), input.Rate, input.Unlimited, tracker, &consumed)
		writeResult <- writeErr
		_ = pipeWriter.CloseWithError(writeErr)
	}()

	request.Body = pipeReader
	request.GetBody = nil
	request.ContentLength = -1
	request.Header.Set("Content-Type", contentType)
	request.Header.Set("Expect", "100-continue")
	if input.Username != "" {
		request.SetBasicAuth(input.Username, string(input.Password))
		defer request.Header.Del("Authorization")
	}
	client := sender.Client
	if client == nil {
		client = DefaultClient()
	}
	response, sendErr := client.Do(request)
	if sendErr != nil {
		_ = pipeReader.CloseWithError(sendErr)
	}
	writeErr := <-writeResult
	result = Result{Bytes: consumed.Load(), Elapsed: time.Since(started)}
	if response == nil {
		cause := errors.Join(sendErr, writeErr)
		if result.Bytes != 0 || writeErr == nil {
			return result, &UnknownOutcomeError{Sent: result.Bytes, Cause: cause}
		}
		return result, cause
	}
	result.StatusCode = response.StatusCode
	closeErr := response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return result, errors.Join(&RejectedError{StatusCode: response.StatusCode}, closeErr)
	}
	if sendErr != nil || writeErr != nil || closeErr != nil {
		return result, errors.Join(sendErr, writeErr, closeErr)
	}
	_ = tracker.AddConfirmed(result.Bytes)
	tracker.Stage(progress.StageComplete)
	return result, nil
}

func validateRequest(input Request) error {
	if input.SourceFS == nil || strings.TrimSpace(input.SourcePath) == "" || input.Name == "" || strings.ContainsAny(input.Name, "/\\\x00\r\n") {
		return errors.New("outgoing webhook source and safe basename are required")
	}
	if (input.Username == "") != (len(input.Password) == 0) {
		return errors.New("outgoing Basic credentials are incomplete")
	}
	if !input.Unlimited && input.Rate <= 0 {
		return errors.New("outgoing webhook rate must be positive or unlimited")
	}
	return nil
}

func writeMultipart(ctx context.Context, writer *multipart.Writer, source io.Reader, name string, expected, bytesPerSecond int64, unlimited bool, tracker *progress.Tracker, consumed *atomic.Int64) error {
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		return err
	}
	stream := &payloadReader{ctx: ctx, source: source, tracker: tracker, consumed: consumed}
	if !unlimited {
		burst := bytesPerSecond
		if burst > maximumRateBurst {
			burst = maximumRateBurst
		}
		stream.limiter = rate.NewLimiter(rate.Limit(bytesPerSecond), int(burst))
		stream.burst = int(burst)
	}
	written, copyErr := io.CopyBuffer(part, stream, make([]byte, streamBufferSize))
	if copyErr == nil && written != expected {
		copyErr = fmt.Errorf("webhook source changed while streaming: read %d bytes, expected %d", written, expected)
	}
	return errors.Join(copyErr, writer.Close())
}

type payloadReader struct {
	ctx      context.Context
	source   io.Reader
	limiter  *rate.Limiter
	burst    int
	tracker  *progress.Tracker
	consumed *atomic.Int64
}

func (reader *payloadReader) Read(buffer []byte) (int, error) {
	if reader.limiter != nil && len(buffer) > reader.burst {
		buffer = buffer[:reader.burst]
	}
	count, err := reader.source.Read(buffer)
	if count != 0 && reader.limiter != nil {
		if waitErr := reader.limiter.WaitN(reader.ctx, count); waitErr != nil {
			return 0, waitErr
		}
	}
	if count != 0 {
		reader.consumed.Add(int64(count))
		if progressErr := reader.tracker.AddCounters(int64(count), int64(count), 0); progressErr != nil {
			return 0, progressErr
		}
	}
	return count, err
}
