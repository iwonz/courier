package policy

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/iwonz/courier/internal/delivery"
)

const opaqueTokenBytes = 32

type SessionConfig struct {
	IdleTimeout     time.Duration
	AbsoluteTimeout time.Duration
}

func DefaultSessionConfig() SessionConfig {
	return SessionConfig{IdleTimeout: 30 * time.Minute, AbsoluteTimeout: 12 * time.Hour}
}

type SessionTokens struct {
	Session string
	CSRF    string
}

type sessionRecord struct {
	deliveryID delivery.ID
	csrfHash   [32]byte
	createdAt  time.Time
	lastSeen   time.Time
}

type Sessions struct {
	mutex  sync.Mutex
	config SessionConfig
	random io.Reader
	now    func() time.Time
	items  map[[32]byte]sessionRecord
}

func NewSessions(config SessionConfig, random io.Reader, now func() time.Time) (*Sessions, error) {
	if config.IdleTimeout <= 0 || config.AbsoluteTimeout <= 0 || config.IdleTimeout > config.AbsoluteTimeout {
		return nil, errors.New("invalid session timeouts")
	}
	if random == nil || now == nil {
		return nil, errors.New("session dependencies are required")
	}
	return &Sessions{config: config, random: random, now: now, items: make(map[[32]byte]sessionRecord)}, nil
}

func (sessions *Sessions) Issue(id delivery.ID) (SessionTokens, error) {
	session, err := sessions.token()
	if err != nil {
		return SessionTokens{}, err
	}
	csrf, err := sessions.token()
	if err != nil {
		return SessionTokens{}, err
	}
	now := sessions.now()
	record := sessionRecord{deliveryID: id, csrfHash: sha256.Sum256([]byte(csrf)), createdAt: now, lastSeen: now}
	sessions.mutex.Lock()
	sessions.items[sha256.Sum256([]byte(session))] = record
	sessions.mutex.Unlock()
	return SessionTokens{Session: session, CSRF: csrf}, nil
}

func (sessions *Sessions) token() (string, error) {
	buffer := make([]byte, opaqueTokenBytes)
	if _, err := io.ReadFull(sessions.random, buffer); err != nil {
		return "", fmt.Errorf("generate session token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}

func (sessions *Sessions) Validate(id delivery.ID, token string) bool {
	key := sha256.Sum256([]byte(token))
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	record, exists := sessions.items[key]
	if !exists || record.deliveryID != id || sessions.expired(record) {
		if exists && sessions.expired(record) {
			delete(sessions.items, key)
		}
		return false
	}
	record.lastSeen = sessions.now()
	sessions.items[key] = record
	return true
}

func (sessions *Sessions) ValidateCSRF(id delivery.ID, token, csrf string) bool {
	key := sha256.Sum256([]byte(token))
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	record, exists := sessions.items[key]
	if !exists || record.deliveryID != id || sessions.expired(record) {
		if exists && sessions.expired(record) {
			delete(sessions.items, key)
		}
		return false
	}
	if record.csrfHash != sha256.Sum256([]byte(csrf)) {
		return false
	}
	record.lastSeen = sessions.now()
	sessions.items[key] = record
	return true
}

func (sessions *Sessions) expired(record sessionRecord) bool {
	now := sessions.now()
	return !now.Before(record.createdAt.Add(sessions.config.AbsoluteTimeout)) || !now.Before(record.lastSeen.Add(sessions.config.IdleTimeout))
}

func (sessions *Sessions) Revoke(token string) {
	sessions.mutex.Lock()
	delete(sessions.items, sha256.Sum256([]byte(token)))
	sessions.mutex.Unlock()
}

func (sessions *Sessions) Purge() int {
	sessions.mutex.Lock()
	defer sessions.mutex.Unlock()
	removed := 0
	for key, record := range sessions.items {
		if sessions.expired(record) {
			delete(sessions.items, key)
			removed++
		}
	}
	return removed
}

func SessionCookie(name, token string, secure bool, lifetime time.Duration) *http.Cookie {
	return &http.Cookie{
		Name: name, Value: token, Path: "/", HttpOnly: true, Secure: secure,
		SameSite: http.SameSiteStrictMode, MaxAge: int(lifetime.Seconds()),
	}
}
