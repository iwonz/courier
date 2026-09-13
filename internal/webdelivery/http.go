package webdelivery

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/iwonz/courier/internal/archive"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/policy"
)

const (
	defaultHeaderTimeout = 10 * time.Second
	maximumLoginBody     = 4096
	streamBufferSize     = 128 * 1024
	sessionCookieName    = "courier_session"
)

var uploadRandom = rand.Read

type peerAddress string

func (address peerAddress) Network() string { return "tcp" }
func (address peerAddress) String() string  { return string(address) }

func (host *Host) authorize(response http.ResponseWriter, request *http.Request, hosted *hostedDelivery, incoming, transfer bool, declared int64) (*policy.Authorization, bool) {
	authentication := policy.Authentication{}
	configured := hosted.policy.Policy()
	if configured.Auth == delivery.AuthBasic {
		authentication.Username, authentication.BasicSecret, _ = basicCredentials(request)
	} else if configured.Auth == delivery.AuthPassword {
		if cookie, err := request.Cookie(sessionCookieName); err == nil {
			authentication.Session = cookie.Value
		}
		authentication.CSRF = request.Header.Get("X-Courier-CSRF")
	}
	authorization, err := hosted.policy.Authorize(request.Context(), policy.Request{
		Peer: peerAddress(request.RemoteAddr), Authentication: authentication,
		StateChanging: request.Method != http.MethodGet && request.Method != http.MethodHead,
		Incoming:      incoming, Transfer: transfer, DeclaredSize: declared,
	})
	clear(authentication.BasicSecret)
	if err == nil {
		return authorization, true
	}
	if errors.Is(err, policy.ErrStopRequested) {
		host.requestStop(hosted.record.ID)
	}
	if configured.Auth == delivery.AuthBasic {
		response.Header().Set("WWW-Authenticate", `Basic realm="Courier", charset="UTF-8"`)
	}
	status := http.StatusUnauthorized
	if errors.Is(err, policy.ErrPeerDenied) || errors.Is(err, policy.ErrBanned) {
		status = http.StatusForbidden
	} else if errors.Is(err, policy.ErrCSRF) {
		status = http.StatusForbidden
	} else if errors.Is(err, policy.ErrStopRequested) {
		status = http.StatusGone
	} else if errors.Is(err, policy.ErrTransferLimit) || errors.Is(err, policy.ErrFileTooLarge) {
		status = http.StatusTooManyRequests
	}
	writeError(response, status, "request denied")
	return nil, false
}

func basicCredentials(request *http.Request) (string, []byte, bool) {
	username, password, ok := request.BasicAuth()
	return username, []byte(password), ok
}

func (host *Host) root(response http.ResponseWriter, request *http.Request) {
	hosted, exists := host.lookup(request)
	if !exists {
		writeError(response, http.StatusNotFound, "delivery not found")
		return
	}
	if hosted.policy.Policy().Auth == delivery.AuthPassword {
		if _, err := request.Cookie(sessionCookieName); err != nil && !hosted.definition.NoUI {
			response.Header().Set("Content-Type", "text/html; charset=utf-8")
			response.Header().Set("Cache-Control", "no-store")
			_, _ = io.WriteString(response, dataPage)
			return
		}
	}
	authorization, ok := host.authorize(response, request, hosted, false, false, 0)
	if !ok {
		return
	}
	defer authorization.Release()
	if hosted.definition.NoUI {
		writeJSON(response, http.StatusOK, map[string]string{"api": "v1", "delivery": string(hosted.record.Route)})
		return
	}
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.Header().Set("Cache-Control", "no-store")
	_, _ = io.WriteString(response, dataPage)
}

func (host *Host) login(response http.ResponseWriter, request *http.Request) {
	hosted, exists := host.lookup(request)
	if !exists {
		writeError(response, http.StatusNotFound, "delivery not found")
		return
	}
	if hosted.policy.Policy().Auth != delivery.AuthPassword {
		writeError(response, http.StatusNotFound, "delivery not found")
		return
	}
	if !sameOrigin(request) {
		writeError(response, http.StatusForbidden, "request origin denied")
		return
	}
	request.Body = http.MaxBytesReader(response, request.Body, maximumLoginBody)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	var input struct {
		Password string `json:"password"`
	}
	if err := decoder.Decode(&input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid login request")
		return
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(response, http.StatusBadRequest, "invalid login request")
		return
	}
	authorization, err := hosted.policy.Authorize(request.Context(), policy.Request{
		Peer: peerAddress(request.RemoteAddr), Authentication: policy.Authentication{Password: []byte(input.Password)},
	})
	input.Password = ""
	if err != nil || authorization.Session == nil {
		if errors.Is(err, policy.ErrStopRequested) {
			host.requestStop(hosted.record.ID)
		}
		writeError(response, http.StatusUnauthorized, "authentication failed")
		return
	}
	defer authorization.Release()
	cookie := policy.SessionCookie(sessionCookieName, authorization.Session.Session, request.TLS != nil, hosted.sessionsLifetime())
	cookie.Path = "/d/" + hosted.definition.Token + "/"
	http.SetCookie(response, cookie)
	writeJSON(response, http.StatusOK, map[string]string{"csrf": authorization.Session.CSRF})
}

func (hosted *hostedDelivery) sessionsLifetime() time.Duration {
	return 12 * time.Hour
}

type metadataResponse struct {
	Name    string          `json:"name"`
	Path    string          `json:"path"`
	Size    int64           `json:"size"`
	Type    string          `json:"type"`
	Entries []metadataEntry `json:"entries,omitempty"`
}

type metadataEntry struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

func (host *Host) metadata(response http.ResponseWriter, request *http.Request) {
	hosted, exists := host.lookup(request)
	if !exists {
		writeError(response, http.StatusNotFound, "delivery not found")
		return
	}
	authorization, ok := host.authorize(response, request, hosted, false, false, 0)
	if !ok {
		return
	}
	defer authorization.Release()
	if hosted.record.Route == delivery.RouteWebToPath {
		writeJSON(response, http.StatusOK, metadataResponse{Name: "Upload", Type: "upload"})
		return
	}
	relative, target, err := resolveRequestedPath(hosted, request.URL.Query().Get("path"))
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid path")
		return
	}
	target, info, err := inspectRequestedPath(hosted, relative)
	if err != nil || !info.Mode().IsRegular() && !info.IsDir() {
		writeError(response, http.StatusNotFound, "object not found")
		return
	}
	result := metadataResponse{Name: publishedName(hosted.resource, info), Path: relative, Size: info.Size(), Type: objectType(info)}
	if info.IsDir() {
		entries, readErr := hosted.resource.Backend.ReadDir(target)
		if readErr != nil {
			writeError(response, http.StatusInternalServerError, "metadata unavailable")
			return
		}
		for _, entry := range entries {
			childRelative := path.Join(relative, entry.Name())
			childInfo, infoErr := entry.Info()
			if infoErr != nil || !hosted.selector.Include(childRelative, entry.IsDir()) || childInfo.Mode()&fs.ModeSymlink != 0 || !childInfo.Mode().IsRegular() && !childInfo.IsDir() {
				continue
			}
			result.Entries = append(result.Entries, metadataEntry{Name: entry.Name(), Size: childInfo.Size(), Type: objectType(childInfo)})
		}
	}
	writeJSON(response, http.StatusOK, result)
}

func objectType(info fs.FileInfo) string {
	if info.IsDir() {
		return "directory"
	}
	return "file"
}

func publishedName(resource *Resource, info fs.FileInfo) string {
	if resource.Name != "" {
		return resource.Name
	}
	return info.Name()
}

func resolveRequestedPath(hosted *hostedDelivery, value string) (string, string, error) {
	if strings.Contains(value, "\\") || strings.ContainsRune(value, 0) || strings.HasPrefix(value, "/") {
		return "", "", errors.New("unsafe relative path")
	}
	clean := path.Clean(value)
	if clean == "." {
		clean = ""
	}
	if clean == ".." || strings.HasPrefix(clean, "../") || clean != value && value != "" {
		return "", "", errors.New("unsafe relative path")
	}
	target := hosted.resource.Path
	if clean != "" {
		target = hosted.resource.Backend.Join(target, clean)
	}
	return clean, target, nil
}

func inspectRequestedPath(hosted *hostedDelivery, relative string) (string, fs.FileInfo, error) {
	target := hosted.resource.Path
	info, err := hosted.resource.Backend.Lstat(target)
	if err != nil || info.Mode()&fs.ModeSymlink != 0 {
		return "", nil, errors.Join(err, errors.New("unsafe requested path"))
	}
	if relative == "" {
		return target, info, nil
	}
	current := ""
	for _, component := range strings.Split(relative, "/") {
		if !info.IsDir() {
			return "", nil, errors.New("requested path crosses a non-directory")
		}
		current = path.Join(current, component)
		target = hosted.resource.Backend.Join(target, component)
		info, err = hosted.resource.Backend.Lstat(target)
		if err != nil || info.Mode()&fs.ModeSymlink != 0 || !hosted.selector.Include(current, info.IsDir()) {
			return "", nil, errors.Join(err, errors.New("unsafe or excluded requested path"))
		}
	}
	return target, info, nil
}

func (host *Host) download(response http.ResponseWriter, request *http.Request) {
	hosted, exists := host.lookup(request)
	if !exists || hosted.record.Route != delivery.RoutePathToWeb {
		writeError(response, http.StatusNotFound, "delivery not found")
		return
	}
	authorization, ok := host.authorize(response, request, hosted, false, true, -1)
	if !ok {
		return
	}
	defer authorization.Release()
	relative, _, err := resolveRequestedPath(hosted, request.URL.Query().Get("path"))
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid path")
		return
	}
	target, info, err := inspectRequestedPath(hosted, relative)
	if err != nil {
		writeError(response, http.StatusNotFound, "object not found")
		return
	}
	if info.IsDir() {
		if request.URL.Query().Get("archive") != "tar.gz" {
			writeError(response, http.StatusBadRequest, "directory download requires archive=tar.gz")
			return
		}
		manifest, walkErr := buildManifest(request.Context(), hosted, target, info.Name(), "")
		if walkErr != nil {
			writeError(response, http.StatusUnprocessableEntity, "directory is not downloadable")
			return
		}
		response.Header().Set("Content-Type", "application/gzip")
		setDownloadName(response, publishedName(hosted.resource, info)+".tar.gz")
		if err := streamArchive(request.Context(), response, hosted.resource.Backend, manifest, authorization); err != nil {
			return
		}
		return
	}
	if !info.Mode().IsRegular() {
		writeError(response, http.StatusUnprocessableEntity, "object is not downloadable")
		return
	}
	reader, err := hosted.resource.Backend.Open(target)
	if err != nil {
		writeError(response, http.StatusInternalServerError, "download unavailable")
		return
	}
	defer reader.Close()
	response.Header().Set("Content-Type", "application/octet-stream")
	response.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))
	setDownloadName(response, publishedName(hosted.resource, info))
	_ = copyRateLimited(request.Context(), response, reader, authorization, policy.Download)
}

func setDownloadName(response http.ResponseWriter, name string) {
	response.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
}

type manifestEntry struct {
	path string
	name string
	info fs.FileInfo
}

func buildManifest(ctx context.Context, hosted *hostedDelivery, target, rootName, relative string) ([]manifestEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	info, err := hosted.resource.Backend.Lstat(target)
	if err != nil {
		return nil, err
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.Mode().IsRegular() && !info.IsDir() {
		return nil, errors.New("unsafe archive object")
	}
	if relative != "" && !hosted.selector.Include(relative, info.IsDir()) {
		return nil, nil
	}
	name := rootName
	if relative != "" {
		name = path.Join(rootName, relative)
	}
	result := []manifestEntry{{path: target, name: name, info: info}}
	if !info.IsDir() {
		return result, nil
	}
	entries, err := hosted.resource.Backend.ReadDir(target)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		childRelative := path.Join(relative, entry.Name())
		children, childErr := buildManifest(ctx, hosted, hosted.resource.Backend.Join(target, entry.Name()), rootName, childRelative)
		if childErr != nil {
			return nil, childErr
		}
		result = append(result, children...)
	}
	return result, nil
}

func streamArchive(ctx context.Context, output io.Writer, backend fsx.Backend, manifest []manifestEntry, authorization *policy.Authorization) (resultErr error) {
	rateWriter := &limitedWriter{ctx: ctx, destination: output, authorization: authorization, direction: policy.Download}
	gzipWriter := gzip.NewWriter(rateWriter)
	tarWriter := tar.NewWriter(gzipWriter)
	defer func() { resultErr = errors.Join(resultErr, tarWriter.Close(), gzipWriter.Close()) }()
	buffer := make([]byte, streamBufferSize)
	for _, entry := range manifest {
		header, err := tar.FileInfoHeader(entry.info, "")
		if err != nil {
			return err
		}
		header.Name = entry.name
		if entry.info.IsDir() {
			header.Name += "/"
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}
		if !entry.info.Mode().IsRegular() {
			continue
		}
		reader, err := backend.Open(entry.path)
		if err != nil {
			return err
		}
		written, copyErr := io.CopyBuffer(tarWriter, reader, buffer)
		closeErr := reader.Close()
		if copyErr != nil || closeErr != nil {
			return errors.Join(copyErr, closeErr)
		}
		if written != entry.info.Size() {
			return errors.New("archive source changed while streaming")
		}
	}
	return nil
}

type limitedWriter struct {
	ctx           context.Context
	destination   io.Writer
	authorization *policy.Authorization
	direction     policy.Direction
}

func (writer *limitedWriter) Write(data []byte) (int, error) {
	if err := writer.authorization.Wait(writer.ctx, writer.direction, len(data)); err != nil {
		return 0, err
	}
	return writer.destination.Write(data)
}

func copyRateLimited(ctx context.Context, destination io.Writer, source io.Reader, authorization *policy.Authorization, direction policy.Direction) error {
	_, err := io.CopyBuffer(&limitedWriter{ctx: ctx, destination: destination, authorization: authorization, direction: direction}, source, make([]byte, streamBufferSize))
	return err
}

func (host *Host) upload(response http.ResponseWriter, request *http.Request) {
	hosted, exists := host.lookup(request)
	if !exists || hosted.record.Route != delivery.RouteWebToPath {
		writeError(response, http.StatusNotFound, "delivery not found")
		return
	}
	if !sameOrigin(request) {
		writeError(response, http.StatusForbidden, "request origin denied")
		return
	}
	declared := int64(-1)
	if value := request.Header.Get("X-Courier-File-Size"); value != "" {
		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil || parsed < 0 {
			writeError(response, http.StatusBadRequest, "invalid declared file size")
			return
		}
		declared = parsed
	}
	authorization, ok := host.authorize(response, request, hosted, true, true, declared)
	if !ok {
		return
	}
	defer authorization.Release()
	mediaType, parameters, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "multipart/form-data" || parameters["boundary"] == "" {
		writeError(response, http.StatusBadRequest, "multipart upload required")
		return
	}
	reader := multipart.NewReader(request.Body, parameters["boundary"])
	part, name, err := firstFilePart(reader)
	if err != nil {
		writeError(response, http.StatusBadRequest, "one file part is required")
		return
	}
	defer part.Close()
	if !safeUploadName(name) || !hosted.definition.Extract && !hosted.selector.Include(name, false) {
		writeError(response, http.StatusBadRequest, "unsafe or excluded file name")
		return
	}
	finalPath := hosted.resource.Backend.Join(hosted.resource.Path, name)
	if _, err := hosted.resource.Backend.Lstat(finalPath); err == nil {
		writeError(response, http.StatusConflict, "destination already exists")
		return
	} else if !errors.Is(err, fs.ErrNotExist) {
		writeError(response, http.StatusInternalServerError, "destination unavailable")
		return
	}
	tokenBytes := make([]byte, 8)
	if _, err := uploadRandom(tokenBytes); err != nil {
		writeError(response, http.StatusInternalServerError, "staging unavailable")
		return
	}
	stagePath := finalPath + ".courier-upload-" + hex.EncodeToString(tokenBytes)
	defer hosted.resource.Backend.RemoveAll(stagePath)
	writer, err := hosted.resource.Backend.Create(stagePath, 0o600)
	if err != nil {
		writeError(response, http.StatusInternalServerError, "staging unavailable")
		return
	}
	limited := &uploadWriter{ctx: request.Context(), destination: writer, authorization: authorization}
	_, copyErr := io.CopyBuffer(limited, part, make([]byte, streamBufferSize))
	closeErr := errors.Join(writer.Sync(), writer.Close())
	if copyErr != nil || closeErr != nil {
		writeError(response, http.StatusUnprocessableEntity, "upload failed")
		return
	}
	if extra, nextErr := reader.NextPart(); nextErr != io.EOF {
		if extra != nil {
			_ = extra.Close()
		}
		writeError(response, http.StatusBadRequest, "exactly one file part is supported")
		return
	}
	if hosted.definition.Extract {
		configured := hosted.policy.Policy()
		result, extractErr := archive.DefaultRegistry().Extract(request.Context(), archive.ExtractionRequest{
			SourceFS: hosted.resource.Backend, SourcePath: stagePath, SourceName: name,
			DestinationFS: hosted.resource.Backend, DestinationRoot: hosted.resource.Path, Selector: hosted.selector,
			Limits: archive.Limits{MaxBytes: configured.MaxExtractedSize.Value, Unlimited: configured.MaxExtractedSize.Unlimited, Configured: true},
		})
		if extractErr != nil {
			status := http.StatusUnprocessableEntity
			if errors.Is(extractErr, fsx.ErrDestinationExists) {
				status = http.StatusConflict
			}
			writeError(response, status, "archive extraction failed")
			return
		}
		writeJSON(response, http.StatusCreated, map[string]any{"name": name, "bytes": authorizationBytes(authorization), "entries": result.Entries})
		return
	}
	if err := hosted.resource.Backend.CommitAbsent(stagePath, finalPath); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, fsx.ErrDestinationExists) {
			status = http.StatusConflict
		}
		writeError(response, status, "upload commit failed")
		return
	}
	writeJSON(response, http.StatusCreated, map[string]any{"name": name, "bytes": authorizationBytes(authorization)})
}

type uploadWriter struct {
	ctx           context.Context
	destination   io.Writer
	authorization *policy.Authorization
}

func (writer *uploadWriter) Write(data []byte) (int, error) {
	if err := writer.authorization.Wait(writer.ctx, policy.Upload, len(data)); err != nil {
		return 0, err
	}
	if err := writer.authorization.Consume(int64(len(data))); err != nil {
		return 0, err
	}
	return writer.destination.Write(data)
}

func authorizationBytes(authorization *policy.Authorization) int64 {
	return authorization.Consumed()
}

func firstFilePart(reader *multipart.Reader) (*multipart.Part, string, error) {
	part, err := reader.NextPart()
	if err != nil {
		return nil, "", err
	}
	if part.FileName() == "" {
		return nil, "", errors.New("first multipart part is not a file")
	}
	return part, part.FileName(), nil
}

func safeUploadName(name string) bool {
	return name != "" && name != "." && name != ".." && path.Base(name) == name && !strings.ContainsAny(name, "\\/\x00\r\n")
}

func sameOrigin(request *http.Request) bool {
	origin := request.Header.Get("Origin")
	if origin == "" {
		return !strings.EqualFold(request.Header.Get("Sec-Fetch-Site"), "cross-site")
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.User != nil || parsed.Host == "" || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	scheme := "http"
	if request.TLS != nil {
		scheme = "https"
	}
	return parsed.Scheme == scheme && strings.EqualFold(parsed.Host, request.Host)
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.Header().Set("Cache-Control", "no-store")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}

func writeError(response http.ResponseWriter, status int, message string) {
	writeJSON(response, status, map[string]string{"error": message})
}

func URL(bind, token string) (string, error) {
	host, port, err := net.SplitHostPort(bind)
	if err != nil {
		return "", err
	}
	if host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	value := url.URL{Scheme: "http", Host: net.JoinHostPort(host, port), Path: "/d/" + token + "/"}
	return value.String(), nil
}
