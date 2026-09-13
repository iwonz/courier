package policy

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

// Argon2Parameters records every work factor needed to verify a credential.
type Argon2Parameters struct {
	Time      uint32
	MemoryKiB uint32
	Threads   uint8
	KeyBytes  uint32
	SaltBytes uint32
}

// DefaultArgon2Parameters returns Courier's production credential work factors.
func DefaultArgon2Parameters() Argon2Parameters {
	return Argon2Parameters{Time: 1, MemoryKiB: 64 * 1024, Threads: 4, KeyBytes: 32, SaltBytes: 16}
}

func (parameters Argon2Parameters) validate() error {
	if parameters.Time == 0 || parameters.MemoryKiB < 8 || parameters.Threads == 0 || parameters.KeyBytes < 16 || parameters.SaltBytes < 16 {
		return errors.New("invalid Argon2id parameters")
	}
	return nil
}

// Verifier contains only salted one-way credential material.
type Verifier struct {
	parameters Argon2Parameters
	salt       []byte
	digest     []byte
}

// NewVerifier converts a credential to an in-memory Argon2id verifier.
func NewVerifier(secret []byte, random io.Reader, parameters Argon2Parameters) (*Verifier, error) {
	if len(secret) == 0 {
		return nil, errors.New("credential must not be empty")
	}
	if random == nil {
		return nil, errors.New("credential random source is required")
	}
	if err := parameters.validate(); err != nil {
		return nil, err
	}
	salt := make([]byte, parameters.SaltBytes)
	if _, err := io.ReadFull(random, salt); err != nil {
		return nil, fmt.Errorf("generate credential salt: %w", err)
	}
	digest := argon2.IDKey(secret, salt, parameters.Time, parameters.MemoryKiB, parameters.Threads, parameters.KeyBytes)
	return &Verifier{parameters: parameters, salt: salt, digest: digest}, nil
}

func (verifier *Verifier) verify(secret []byte) bool {
	if verifier == nil {
		return false
	}
	candidate := argon2.IDKey(secret, verifier.salt, verifier.parameters.Time, verifier.parameters.MemoryKiB, verifier.parameters.Threads, verifier.parameters.KeyBytes)
	matched := subtle.ConstantTimeCompare(candidate, verifier.digest) == 1
	clear(candidate)
	return matched
}

// Credentials are consumed during engine construction and are never retained.
type Credentials struct {
	BasicUsername string
	BasicPassword []byte
	Password      []byte
}

type Authenticator struct {
	basicUsername []byte
	basic         *Verifier
	password      *Verifier
	slots         chan struct{}
}

func NewAuthenticator(credentials Credentials, random io.Reader, parameters Argon2Parameters, concurrency uint32) (*Authenticator, error) {
	if concurrency == 0 {
		return nil, errors.New("authentication concurrency must be positive")
	}
	authenticator := &Authenticator{slots: make(chan struct{}, concurrency)}
	if credentials.BasicUsername != "" || len(credentials.BasicPassword) != 0 {
		if credentials.BasicUsername == "" || len(credentials.BasicPassword) == 0 {
			return nil, errors.New("Basic authentication requires username and password")
		}
		verifier, err := NewVerifier(credentials.BasicPassword, random, parameters)
		if err != nil {
			return nil, err
		}
		authenticator.basicUsername = []byte(credentials.BasicUsername)
		authenticator.basic = verifier
	}
	if len(credentials.Password) != 0 {
		verifier, err := NewVerifier(credentials.Password, random, parameters)
		if err != nil {
			return nil, err
		}
		authenticator.password = verifier
	}
	return authenticator, nil
}

func (authenticator *Authenticator) VerifyBasic(ctx context.Context, username string, password []byte) (bool, error) {
	return authenticator.verify(ctx, password, func() bool {
		userMatch := subtle.ConstantTimeCompare([]byte(username), authenticator.basicUsername) == 1
		return authenticator.basic.verify(password) && userMatch
	})
}

func (authenticator *Authenticator) VerifyPassword(ctx context.Context, password []byte) (bool, error) {
	return authenticator.verify(ctx, password, func() bool { return authenticator.password.verify(password) })
}

func (authenticator *Authenticator) verify(ctx context.Context, secret []byte, verify func() bool) (bool, error) {
	if ctx == nil {
		return false, errors.New("authentication context is required")
	}
	select {
	case authenticator.slots <- struct{}{}:
		defer func() { <-authenticator.slots }()
	case <-ctx.Done():
		return false, ctx.Err()
	}
	return verify(), nil
}
