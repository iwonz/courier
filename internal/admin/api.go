// Package admin provides Courier's loopback-only administration process,
// guarded HTTP API, and embedded administration application.
package admin

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iwonz/courier/internal/control"
	"github.com/iwonz/courier/internal/delivery"
)

const maxAPIRequestBytes = 1 << 20

type Controller interface {
	List(context.Context) ([]control.ServerView, error)
	Stop(context.Context, control.StopRequest) (control.StopResult, error)
	UpdatePolicy(context.Context, delivery.ID, uint64, delivery.Policy) error
}

type APIConfig struct {
	Host           string
	Control        Controller
	EventInterval  time.Duration
	MaxSubscribers int
}

type API struct {
	host          string
	origin        string
	control       Controller
	eventInterval time.Duration
	subscribers   chan struct{}
	handler       http.Handler
}

type Snapshot struct {
	Servers []ServerSnapshot `json:"servers"`
}

type ServerSnapshot struct {
	ID         delivery.ID        `json:"id"`
	Bind       string             `json:"bind"`
	ProcessID  int                `json:"processId"`
	State      delivery.State     `json:"state"`
	Status     string             `json:"status"`
	StartedAt  time.Time          `json:"startedAt"`
	UpdatedAt  time.Time          `json:"updatedAt"`
	Deliveries []DeliverySnapshot `json:"deliveries"`
}

type DeliverySnapshot struct {
	ID          delivery.ID              `json:"id"`
	Route       delivery.Route           `json:"route"`
	Source      string                   `json:"source"`
	Destination string                   `json:"destination"`
	State       delivery.State           `json:"state"`
	Policy      delivery.Policy          `json:"policy"`
	Counters    delivery.CounterSnapshot `json:"counters"`
	CreatedAt   time.Time                `json:"createdAt"`
	UpdatedAt   time.Time                `json:"updatedAt"`
}

type policyUpdate struct {
	ExpectedVersion uint64          `json:"expectedVersion"`
	Policy          delivery.Policy `json:"policy"`
}

func NewAPI(config APIConfig) (*API, error) {
	if config.Control == nil || strings.TrimSpace(config.Host) == "" {
		return nil, fmt.Errorf("%w: admin API dependencies are incomplete", delivery.ErrInvalid)
	}
	interval := config.EventInterval
	if interval <= 0 {
		interval = time.Second
	}
	maximum := config.MaxSubscribers
	if maximum <= 0 {
		maximum = 16
	}
	api := &API{
		host: config.Host, origin: "http://" + config.Host, control: config.Control,
		eventInterval: interval, subscribers: make(chan struct{}, maximum),
	}
	router := chi.NewRouter()
	router.Use(api.guard)
	router.Get("/api/v1/servers", api.list)
	router.Get("/api/v1/events", api.events)
	router.Put("/api/v1/deliveries/{id}/policy", api.updatePolicy)
	router.Post("/api/v1/deliveries/{id}/stop", api.stopDelivery)
	router.Post("/api/v1/servers/{id}/stop", api.stopServer)
	router.Handle("/assets/*", adminAssetHandler())
	router.Get("/", func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.WriteString(response, adminPage)
	})
	api.handler = router
	return api, nil
}

func (api *API) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	api.handler.ServeHTTP(response, request)
}

func (api *API) guard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		setAdminHeaders(response.Header())
		if request.Host != api.host || !loopbackRemote(request.RemoteAddr) {
			http.Error(response, "request rejected", http.StatusForbidden)
			return
		}
		if request.Method != http.MethodGet && request.Method != http.MethodHead && request.Header.Get("Origin") != api.origin {
			http.Error(response, "request rejected", http.StatusForbidden)
			return
		}
		next.ServeHTTP(response, request)
	})
}

func setAdminHeaders(header http.Header) {
	header.Set("Cache-Control", "no-store")
	header.Set("Content-Security-Policy", "default-src 'self'; connect-src 'self'; img-src 'self'; style-src 'self'; script-src 'self'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
	header.Set("Referrer-Policy", "no-referrer")
	header.Set("X-Content-Type-Options", "nosniff")
	header.Set("X-Frame-Options", "DENY")
}

func loopbackRemote(remote string) bool {
	host, _, err := net.SplitHostPort(remote)
	if err != nil {
		return false
	}
	address, err := netip.ParseAddr(host)
	return err == nil && address.IsLoopback()
}

func (api *API) list(response http.ResponseWriter, request *http.Request) {
	snapshot, err := api.snapshot(request.Context())
	if err != nil {
		writeAPIError(response, err)
		return
	}
	writeJSON(response, http.StatusOK, snapshot)
}

func (api *API) snapshot(ctx context.Context) (Snapshot, error) {
	views, err := api.control.List(ctx)
	if err != nil {
		return Snapshot{}, err
	}
	result := Snapshot{Servers: make([]ServerSnapshot, 0, len(views))}
	for _, view := range views {
		status := "unreachable"
		if view.Live {
			status = "live"
		}
		server := ServerSnapshot{
			ID: view.Server.ID, Bind: view.Server.Bind, ProcessID: view.Server.ProcessID, State: view.Server.State,
			Status: status, StartedAt: view.Server.StartedAt, UpdatedAt: view.Server.UpdatedAt,
			Deliveries: make([]DeliverySnapshot, 0, len(view.Deliveries)),
		}
		for _, item := range view.Deliveries {
			server.Deliveries = append(server.Deliveries, DeliverySnapshot{
				ID: item.ID, Route: item.Route, Source: item.Source, Destination: item.Destination, State: item.State,
				Policy: item.Policy, Counters: item.Counters, CreatedAt: item.CreatedAt, UpdatedAt: item.UpdatedAt,
			})
		}
		result.Servers = append(result.Servers, server)
	}
	return result, nil
}

func (api *API) updatePolicy(response http.ResponseWriter, request *http.Request) {
	id, ok := requestID(response, request)
	if !ok {
		return
	}
	var update policyUpdate
	if err := decodeJSON(response, request, &update); err != nil {
		writeAPIError(response, err)
		return
	}
	if err := api.control.UpdatePolicy(request.Context(), id, update.ExpectedVersion, update.Policy); err != nil {
		writeAPIError(response, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (api *API) stopDelivery(response http.ResponseWriter, request *http.Request) {
	api.stop(response, request, delivery.TargetDelivery)
}

func (api *API) stopServer(response http.ResponseWriter, request *http.Request) {
	api.stop(response, request, delivery.TargetServer)
}

func (api *API) stop(response http.ResponseWriter, request *http.Request, expected delivery.TargetKind) {
	id, ok := requestID(response, request)
	if !ok {
		return
	}
	result, err := api.control.Stop(request.Context(), control.StopRequest{ID: id, Kind: expected})
	if err != nil {
		writeAPIError(response, err)
		return
	}
	if result.Kind != expected {
		writeAPIError(response, delivery.ErrNotFound)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

func (api *API) events(response http.ResponseWriter, request *http.Request) {
	flusher, ok := response.(http.Flusher)
	if !ok {
		writeAPIError(response, errors.New("streaming unavailable"))
		return
	}
	select {
	case api.subscribers <- struct{}{}:
		defer func() { <-api.subscribers }()
	default:
		http.Error(response, "too many event subscribers", http.StatusTooManyRequests)
		return
	}
	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Connection", "keep-alive")
	ticker := time.NewTicker(api.eventInterval)
	defer ticker.Stop()
	for {
		snapshot, err := api.snapshot(request.Context())
		if err != nil {
			return
		}
		data, _ := json.Marshal(snapshot) // snapshot contains only JSON-safe value types
		if _, err := fmt.Fprintf(response, "event: snapshot\ndata: %s\n\n", data); err != nil {
			return
		}
		flusher.Flush()
		select {
		case <-request.Context().Done():
			return
		case <-ticker.C:
		}
	}
}

func requestID(response http.ResponseWriter, request *http.Request) (delivery.ID, bool) {
	id, err := delivery.ParseID(chi.URLParam(request, "id"))
	if err != nil {
		writeAPIError(response, err)
		return "", false
	}
	return id, true
}

func decodeJSON(response http.ResponseWriter, request *http.Request, destination any) error {
	if request.Header.Get("Content-Type") != "application/json" {
		return fmt.Errorf("%w: application/json is required", delivery.ErrInvalid)
	}
	request.Body = http.MaxBytesReader(response, request.Body, maxAPIRequestBytes)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("%w: invalid JSON", delivery.ErrInvalid)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("%w: trailing JSON", delivery.ErrInvalid)
	}
	return nil
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}

func writeAPIError(response http.ResponseWriter, err error) {
	status := http.StatusServiceUnavailable
	switch {
	case errors.Is(err, delivery.ErrInvalid):
		status = http.StatusBadRequest
	case errors.Is(err, delivery.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, delivery.ErrRevisionConflict):
		status = http.StatusConflict
	}
	writeJSON(response, status, map[string]string{"error": http.StatusText(status)})
}
