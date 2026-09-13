package webdelivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/iwonz/courier/internal/archive"
	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/fsx"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
	"github.com/iwonz/courier/internal/selection"
)

type Resource struct {
	Endpoint endpoint.Endpoint
	Backend  fsx.Backend
	Path     string
	Name     string
	Close    func() error
	Temps    []delivery.OwnedTemp
}

type OpenFunc func(context.Context, endpoint.Endpoint, EndpointRuntime) (*Resource, error)
type SelectFunc func([]operation.SelectionRule) (selection.Selector, error)
type ArchiveFunc func(context.Context, *Resource, string, selection.Selector) (*Resource, error)

type HostOptions struct {
	Open               OpenFunc
	Select             SelectFunc
	Archive            ArchiveFunc
	PolicyDependencies policy.Dependencies
	Server             *http.Server
	StopRequested      func(delivery.ID)
}

type hostedDelivery struct {
	record     delivery.Delivery
	definition Definition
	resource   *Resource
	selector   selection.Selector
	policy     *policy.Engine
	temps      []delivery.OwnedTemp
	closeOnce  sync.Once
	closeErr   error
}

func (hosted *hostedDelivery) close() error {
	hosted.closeOnce.Do(func() {
		if hosted.resource != nil && hosted.resource.Close != nil {
			hosted.closeErr = hosted.resource.Close()
		}
	})
	return hosted.closeErr
}

type Host struct {
	mutex         sync.RWMutex
	open          OpenFunc
	selectRules   SelectFunc
	archive       ArchiveFunc
	dependencies  policy.Dependencies
	server        *http.Server
	assets        http.Handler
	stopRequested func(delivery.ID)
	byToken       map[string]*hostedDelivery
	byID          map[delivery.ID]*hostedDelivery
}

func NewHost(options HostOptions) (*Host, error) {
	if options.Open == nil || options.Select == nil {
		return nil, errors.New("web delivery host dependencies are incomplete")
	}
	server := options.Server
	if server == nil {
		server = &http.Server{ReadHeaderTimeout: defaultHeaderTimeout}
	}
	host := &Host{
		open: options.Open, selectRules: options.Select, dependencies: options.PolicyDependencies,
		archive: options.Archive, server: server, assets: dataAssetHandler(), stopRequested: options.StopRequested,
		byToken: make(map[string]*hostedDelivery), byID: make(map[delivery.ID]*hostedDelivery),
	}
	if host.archive == nil {
		host.archive = prepareArchiveResource
	}
	host.server.Handler = host.routes()
	return host, nil
}

func (host *Host) Register(ctx context.Context, record delivery.Delivery, runtime json.RawMessage) error {
	if err := record.Validate(); err != nil {
		return err
	}
	var definition Definition
	decoder := json.NewDecoder(bytes.NewReader(runtime))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&definition); err != nil {
		return fmt.Errorf("decode web delivery definition: %w", err)
	}
	defer definition.ClearSecrets()
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("trailing web delivery definition")
	}
	if err := definition.Validate(record.Route); err != nil {
		return err
	}
	selector, err := host.selectRules(definition.Selection)
	if err != nil {
		return err
	}
	openedEndpoint := definition.Source
	if record.Route == delivery.RouteWebToPath || record.Route == delivery.RouteWebhookToPath {
		openedEndpoint = definition.Destination
	}
	parsed, _ := endpoint.Parse(openedEndpoint)
	endpointRuntime := definition.Endpoint
	definition.Endpoint = EndpointRuntime{}
	defer endpointRuntime.Clear()
	resource, err := host.open(ctx, parsed, endpointRuntime)
	if err != nil {
		return err
	}
	success := false
	defer func() {
		if !success && resource.Close != nil {
			_ = resource.Close()
		}
	}()
	info, err := resource.Backend.Lstat(resource.Path)
	if err != nil {
		return err
	}
	if (record.Route == delivery.RouteWebToPath || record.Route == delivery.RouteWebhookToPath) && !info.IsDir() {
		return errors.New("incoming delivery destination must be an existing directory")
	}
	if record.Route == delivery.RoutePathToWeb && !info.Mode().IsRegular() && !info.IsDir() {
		return errors.New("browser download source must be a file or directory")
	}
	if record.Route == delivery.RoutePathToWeb && definition.Archive {
		prepared, prepareErr := host.archive(ctx, resource, info.Name(), selector)
		if prepareErr != nil {
			return prepareErr
		}
		if prepared == nil || prepared.Backend == nil || prepared.Path == "" || prepared.Close == nil {
			if prepared != nil && prepared.Close != nil {
				_ = prepared.Close()
			}
			return errors.New("browser archive preparation returned an incomplete resource")
		}
		resource = prepared
		info, err = resource.Backend.Lstat(resource.Path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return errors.New("prepared browser archive is not a regular file")
		}
	}
	ownedTemps := append([]delivery.OwnedTemp(nil), resource.Temps...)
	for index := range ownedTemps {
		ownedTemps[index].OwnerID = record.ID
		ownedTemps[index].CreatedAt = record.CreatedAt
		if err := ownedTemps[index].Validate(); err != nil {
			return err
		}
	}
	engine, err := policy.NewEngine(record.ID, record.Policy, definition.Credentials, host.dependencies)
	definition.Credentials = policy.Credentials{}
	if err != nil {
		return err
	}
	hosted := &hostedDelivery{record: record, definition: definition, resource: resource, selector: selector, policy: engine, temps: ownedTemps}
	host.mutex.Lock()
	defer host.mutex.Unlock()
	if _, exists := host.byToken[definition.Token]; exists {
		return errors.New("duplicate web delivery resource token")
	}
	if _, exists := host.byID[record.ID]; exists {
		return delivery.ErrDuplicate
	}
	host.byToken[definition.Token] = hosted
	host.byID[record.ID] = hosted
	success = true
	return nil
}

func (host *Host) Stop(_ context.Context, id delivery.ID) error {
	host.mutex.Lock()
	hosted, exists := host.byID[id]
	if exists {
		delete(host.byID, id)
		delete(host.byToken, hosted.definition.Token)
	}
	host.mutex.Unlock()
	if !exists {
		return nil
	}
	return hosted.close()
}

func (host *Host) PreparePolicy(_ context.Context, id delivery.ID, expected uint64, next delivery.Policy) (func(), func(), error) {
	host.mutex.RLock()
	hosted, exists := host.byID[id]
	host.mutex.RUnlock()
	if !exists {
		return nil, nil, delivery.ErrNotFound
	}
	prepared, err := hosted.policy.PrepareUpdate(expected, next)
	if err != nil {
		return nil, nil, err
	}
	return prepared.Commit, prepared.Cancel, nil
}

func (host *Host) OwnedTemps(id delivery.ID) []delivery.OwnedTemp {
	host.mutex.RLock()
	defer host.mutex.RUnlock()
	hosted := host.byID[id]
	if hosted == nil {
		return nil
	}
	return append([]delivery.OwnedTemp(nil), hosted.temps...)
}

func (host *Host) Close(ctx context.Context) error {
	host.mutex.Lock()
	items := make([]*hostedDelivery, 0, len(host.byID))
	for _, hosted := range host.byID {
		items = append(items, hosted)
	}
	host.byID = make(map[delivery.ID]*hostedDelivery)
	host.byToken = make(map[string]*hostedDelivery)
	host.mutex.Unlock()
	var result error
	for _, hosted := range items {
		result = errors.Join(result, hosted.close())
	}
	return errors.Join(result, host.server.Shutdown(ctx))
}

func (host *Host) Serve(listener net.Listener) error {
	if listener == nil {
		return errors.New("web delivery listener is required")
	}
	err := host.server.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (host *Host) Handler() http.Handler { return host.server.Handler }

func (host *Host) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(dataSecurityHeaders)
	router.Handle("/assets/*", host.assets)
	router.Route("/d/{token}", func(router chi.Router) {
		router.Get("/", host.root)
		router.Post("/upload", host.webhookUpload)
		router.Post("/api/v1/session", host.login)
		router.Get("/api/v1/meta", host.metadata)
		router.Get("/api/v1/download", host.download)
		router.Post("/api/v1/upload", host.upload)
	})
	return router
}

func dataSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self'; connect-src 'self'; form-action 'self'; frame-ancestors 'none'; base-uri 'none'")
		response.Header().Set("Referrer-Policy", "no-referrer")
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("X-Frame-Options", "DENY")
		next.ServeHTTP(response, request)
	})
}

func (host *Host) lookup(request *http.Request) (*hostedDelivery, bool) {
	host.mutex.RLock()
	defer host.mutex.RUnlock()
	hosted, exists := host.byToken[chi.URLParam(request, "token")]
	return hosted, exists
}

func (host *Host) requestStop(id delivery.ID) {
	if host.stopRequested != nil {
		host.stopRequested(id)
		return
	}
	_ = host.Stop(context.Background(), id)
}

func prepareArchiveResource(ctx context.Context, source *Resource, sourceName string, selector selection.Selector) (*Resource, error) {
	artifact, err := archive.CreateSelected(ctx, source.Backend, source.Path, sourceName, os.TempDir(), selector, nil)
	if err != nil {
		return nil, err
	}
	return &Resource{
		Endpoint: source.Endpoint, Backend: fsx.Local{}, Path: artifact.Path, Name: artifact.Name,
		Temps: []delivery.OwnedTemp{{ID: delivery.NewID(), Location: delivery.TempLocal, Path: artifact.Path}},
		Close: func() error {
			var sourceErr error
			if source.Close != nil {
				sourceErr = source.Close()
			}
			return errors.Join(artifact.Cleanup(), sourceErr)
		},
	}, nil
}
