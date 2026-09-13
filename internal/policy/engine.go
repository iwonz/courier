package policy

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/iwonz/courier/internal/delivery"
)

var (
	ErrPeerDenied     = errors.New("connection peer is not allowed")
	ErrAuthentication = errors.New("authentication failed")
	ErrBanned         = errors.New("authentication fingerprint is banned")
	ErrStopRequested  = errors.New("authentication policy requests delivery stop")
	ErrCSRF           = errors.New("CSRF validation failed")
)

type Stage string

const (
	StageBounds      Stage = "request-bounds"
	StageAdmission   Stage = "peer-admission"
	StageBan         Stage = "ban-check"
	StageAuth        Stage = "authentication"
	StageCSRF        Stage = "csrf"
	StageReservation Stage = "reservation"
	StageRate        Stage = "rate-limit"
	StageHandler     Stage = "handler"
)

func EnforcementOrder() []Stage {
	return []Stage{StageBounds, StageAdmission, StageBan, StageAuth, StageCSRF, StageReservation, StageRate, StageHandler}
}

type Dependencies struct {
	Random          io.Reader
	Now             func() time.Time
	Argon2          Argon2Parameters
	AuthConcurrency uint32
	Sessions        SessionConfig
}

func DefaultDependencies() Dependencies {
	return Dependencies{Random: rand.Reader, Now: time.Now, Argon2: DefaultArgon2Parameters(), AuthConcurrency: 2, Sessions: DefaultSessionConfig()}
}

type Engine struct {
	mutex         sync.RWMutex
	id            delivery.ID
	policy        delivery.Policy
	authenticator *Authenticator
	sessions      *Sessions
	admission     *Admission
	attempts      *AttemptTracker
	reservations  *Reservations
	rates         *Rates
}

func NewEngine(id delivery.ID, configured delivery.Policy, credentials Credentials, dependencies Dependencies) (*Engine, error) {
	if !id.Valid() {
		return nil, errors.New("delivery identity is invalid")
	}
	if err := configured.Validate(); err != nil {
		return nil, err
	}
	authenticator, err := NewAuthenticator(credentials, dependencies.Random, dependencies.Argon2, dependencies.AuthConcurrency)
	if err != nil {
		return nil, err
	}
	if configured.Auth == delivery.AuthBasic && authenticator.basic == nil {
		return nil, errors.New("Basic credentials are required")
	}
	if configured.Auth == delivery.AuthPassword && authenticator.password == nil {
		return nil, errors.New("password credential is required")
	}
	sessions, err := NewSessions(dependencies.Sessions, dependencies.Random, dependencies.Now)
	if err != nil {
		return nil, err
	}
	admission, _ := NewAdmission(configured.AllowIP)
	attempts, _ := NewAttemptTracker(configured.AuthAttempts, configured.AuthFailAction)
	reservations, _ := NewReservations(configured.DeliveryLimit, configured.MaxFileSize)
	rates, _ := NewRates(configured.UploadRate, configured.DownloadRate)
	return &Engine{id: id, policy: clonePolicy(configured), authenticator: authenticator, sessions: sessions, admission: admission, attempts: attempts, reservations: reservations, rates: rates}, nil
}

func (engine *Engine) Policy() delivery.Policy {
	engine.mutex.RLock()
	defer engine.mutex.RUnlock()
	return clonePolicy(engine.policy)
}

func clonePolicy(policy delivery.Policy) delivery.Policy {
	policy.AllowIP = append([]string(nil), policy.AllowIP...)
	return policy
}

func (engine *Engine) Update(expected uint64, next delivery.Policy) error {
	if err := next.Validate(); err != nil {
		return err
	}
	admission, _ := NewAdmission(next.AllowIP)
	engine.mutex.Lock()
	defer engine.mutex.Unlock()
	if engine.policy.Version != expected || next.Version != expected+1 {
		return delivery.ErrRevisionConflict
	}
	if next.Auth != engine.policy.Auth {
		return errors.New("authentication mode changes require new credential material")
	}
	_ = engine.attempts.Update(next.AuthAttempts, next.AuthFailAction)
	_ = engine.reservations.Update(next.DeliveryLimit, next.MaxFileSize)
	_ = engine.rates.Update(next.UploadRate, next.DownloadRate)
	engine.admission = admission
	engine.policy = clonePolicy(next)
	return nil
}

type Authentication struct {
	Username    string
	BasicSecret []byte
	Password    []byte
	Session     string
	CSRF        string
}

type Request struct {
	Peer           net.Addr
	Authentication Authentication
	StateChanging  bool
	Incoming       bool
	DeclaredSize   int64
}

type Authorization struct {
	Session     *SessionTokens
	reservation *Reservation
	rates       *Rates
}

func (authorization *Authorization) Consume(bytes int64) error {
	if authorization.reservation == nil {
		return errors.New("authorization has no incoming reservation")
	}
	return authorization.reservation.Consume(bytes)
}

func (authorization *Authorization) Wait(ctx context.Context, direction Direction, bytes int) error {
	return authorization.rates.Wait(ctx, direction, bytes)
}

func (authorization *Authorization) Release() {
	if authorization.reservation != nil {
		authorization.reservation.Release()
	}
}

func (engine *Engine) Authorize(ctx context.Context, request Request) (*Authorization, error) {
	if ctx == nil {
		return nil, errors.New("authorization context is required")
	}
	if request.Incoming && request.DeclaredSize < -1 {
		return nil, errors.New("declared file size must be non-negative or unknown")
	}
	address, err := PeerIP(request.Peer)
	if err != nil {
		return nil, err
	}
	engine.mutex.RLock()
	defer engine.mutex.RUnlock()
	if !engine.admission.Allows(address) {
		return nil, ErrPeerDenied
	}
	blocked := engine.attempts.Blocked(engine.id, address)
	if blocked.Stop {
		return nil, ErrStopRequested
	}
	if blocked.Banned {
		return nil, ErrBanned
	}
	authorized := &Authorization{rates: engine.rates}
	var issued *SessionTokens
	var authenticated bool
	switch engine.policy.Auth {
	case delivery.AuthNone:
		authenticated = true
	case delivery.AuthBasic:
		authenticated, err = engine.authenticator.VerifyBasic(ctx, request.Authentication.Username, request.Authentication.BasicSecret)
	case delivery.AuthPassword:
		if request.Authentication.Session != "" {
			if request.StateChanging {
				authenticated = engine.sessions.ValidateCSRF(engine.id, request.Authentication.Session, request.Authentication.CSRF)
				if !authenticated {
					return nil, ErrCSRF
				}
			} else {
				authenticated = engine.sessions.Validate(engine.id, request.Authentication.Session)
			}
		} else {
			authenticated, err = engine.authenticator.VerifyPassword(ctx, request.Authentication.Password)
			if err == nil && authenticated {
				tokens, issueErr := engine.sessions.Issue(engine.id)
				if issueErr != nil {
					return nil, issueErr
				}
				issued = &tokens
			}
		}
	}
	if err != nil {
		return nil, err
	}
	if !authenticated {
		decision := engine.attempts.Failure(engine.id, address)
		if decision.Stop {
			return nil, ErrStopRequested
		}
		if decision.Banned {
			return nil, ErrBanned
		}
		return nil, ErrAuthentication
	}
	engine.attempts.Success(engine.id, address)
	authorized.Session = issued
	if request.Incoming {
		authorized.reservation, err = engine.reservations.Reserve(request.DeclaredSize)
		if err != nil {
			return nil, err
		}
	}
	return authorized, nil
}
