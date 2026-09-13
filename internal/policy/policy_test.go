package policy

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"net/netip"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/iwonz/courier/internal/delivery"
)

const testDeliveryID delivery.ID = "00000000-0000-4000-8000-000000000022"

type testAddress string

func (address testAddress) Network() string { return "test" }
func (address testAddress) String() string  { return string(address) }

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func fastParameters() Argon2Parameters {
	return Argon2Parameters{Time: 1, MemoryKiB: 8, Threads: 1, KeyBytes: 16, SaltBytes: 16}
}

func randomBytes() *bytes.Reader {
	return bytes.NewReader(bytes.Repeat([]byte{0x5a}, 4096))
}

func testPolicy(mode delivery.AuthMode) delivery.Policy {
	configured := delivery.DefaultPolicy()
	configured.Auth = mode
	return configured
}

func testDependencies(now func() time.Time) Dependencies {
	return Dependencies{
		Random: randomBytes(), Now: now, Argon2: fastParameters(), AuthConcurrency: 1,
		Sessions: SessionConfig{IdleTimeout: time.Minute, AbsoluteTimeout: time.Hour},
	}
}

func TestVerifierAndAuthenticator(t *testing.T) {
	if DefaultArgon2Parameters().MemoryKiB != 64*1024 {
		t.Fatal("unexpected production Argon2id parameters")
	}
	invalid := fastParameters()
	invalid.Time = 0
	if _, err := NewVerifier([]byte("secret"), randomBytes(), invalid); err == nil {
		t.Fatal("expected invalid parameters")
	}
	for name, parameters := range map[string]Argon2Parameters{
		"memory":  {Time: 1, MemoryKiB: 7, Threads: 1, KeyBytes: 16, SaltBytes: 16},
		"threads": {Time: 1, MemoryKiB: 8, Threads: 0, KeyBytes: 16, SaltBytes: 16},
		"key":     {Time: 1, MemoryKiB: 8, Threads: 1, KeyBytes: 15, SaltBytes: 16},
		"salt":    {Time: 1, MemoryKiB: 8, Threads: 1, KeyBytes: 16, SaltBytes: 15},
	} {
		t.Run(name, func(t *testing.T) {
			if err := parameters.validate(); err == nil {
				t.Fatal("expected parameter error")
			}
		})
	}
	if _, err := NewVerifier(nil, randomBytes(), fastParameters()); err == nil {
		t.Fatal("expected empty credential error")
	}
	if _, err := NewVerifier([]byte("secret"), nil, fastParameters()); err == nil {
		t.Fatal("expected random source error")
	}
	if _, err := NewVerifier([]byte("secret"), failingReader{}, fastParameters()); err == nil {
		t.Fatal("expected salt error")
	}
	verifier, err := NewVerifier([]byte("secret"), randomBytes(), fastParameters())
	if err != nil || !verifier.verify([]byte("secret")) || verifier.verify([]byte("wrong")) || (*Verifier)(nil).verify([]byte("secret")) {
		t.Fatalf("unexpected verifier result: %v", err)
	}

	if _, err := NewAuthenticator(Credentials{}, randomBytes(), fastParameters(), 0); err == nil {
		t.Fatal("expected concurrency error")
	}
	if _, err := NewAuthenticator(Credentials{BasicUsername: "user"}, randomBytes(), fastParameters(), 1); err == nil {
		t.Fatal("expected incomplete Basic credentials")
	}
	if _, err := NewAuthenticator(Credentials{BasicUsername: "user", BasicPassword: []byte("pass")}, failingReader{}, fastParameters(), 1); err == nil {
		t.Fatal("expected Basic verifier error")
	}
	if _, err := NewAuthenticator(Credentials{Password: []byte("pass")}, failingReader{}, fastParameters(), 1); err == nil {
		t.Fatal("expected password verifier error")
	}
	authenticator, err := NewAuthenticator(Credentials{BasicUsername: "user", BasicPassword: []byte("basic"), Password: []byte("password")}, randomBytes(), fastParameters(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := authenticator.VerifyBasic(context.Background(), "user", []byte("basic")); err != nil || !ok {
		t.Fatalf("expected Basic match: %v", err)
	}
	if ok, err := authenticator.VerifyBasic(context.Background(), "other", []byte("basic")); err != nil || ok {
		t.Fatalf("expected username mismatch: %v", err)
	}
	if ok, err := authenticator.VerifyPassword(context.Background(), []byte("wrong")); err != nil || ok {
		t.Fatalf("expected password mismatch: %v", err)
	}
	if _, err := authenticator.VerifyPassword(nil, []byte("password")); err == nil {
		t.Fatal("expected nil context error")
	}
	authenticator.slots <- struct{}{}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := authenticator.VerifyPassword(canceled, []byte("password")); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled authentication, got %v", err)
	}
	<-authenticator.slots
	empty, err := NewAuthenticator(Credentials{}, randomBytes(), fastParameters(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := empty.VerifyBasic(context.Background(), "", nil); err != nil || ok {
		t.Fatalf("missing verifier must reject: %v", err)
	}
}

func TestAdmissionCanonicalization(t *testing.T) {
	if _, err := NewAdmission([]string{""}); err == nil {
		t.Fatal("expected empty admission error")
	}
	if _, err := NewAdmission([]string{"not-an-ip"}); err == nil {
		t.Fatal("expected invalid admission error")
	}
	if _, err := NewAdmission([]string{"::ffff:192.0.2.0/95"}); err == nil {
		t.Fatal("expected ambiguous mapped prefix error")
	}
	admission, err := NewAdmission([]string{"192.0.2.1", "192.0.2.1/32", "2001:db8::/32", "::ffff:198.51.100.0/120"})
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"192.0.2.1", "2001:db8::5", "198.51.100.3"} {
		if !admission.Allows(netip.MustParseAddr(value)) {
			t.Fatalf("expected %s to be admitted", value)
		}
	}
	if admission.Allows(netip.MustParseAddr("203.0.113.1")) {
		t.Fatal("unexpected admission")
	}
	open, err := NewAdmission(nil)
	if err != nil || !open.Allows(netip.MustParseAddr("203.0.113.1")) {
		t.Fatalf("empty allowlist must allow: %v", err)
	}
	if _, err := PeerIP(nil); err == nil {
		t.Fatal("expected nil peer error")
	}
	if _, err := PeerIP(&net.TCPAddr{IP: net.IP{1, 2}}); err == nil {
		t.Fatal("expected invalid TCP IP")
	}
	address, err := PeerIP(&net.TCPAddr{IP: net.ParseIP("::ffff:192.0.2.8"), Port: 8})
	if err != nil || address.String() != "192.0.2.8" {
		t.Fatalf("unexpected mapped peer: %s %v", address, err)
	}
	address, err = PeerIP(testAddress("[2001:db8::9]:443"))
	if err != nil || address.String() != "2001:db8::9" {
		t.Fatalf("unexpected IPv6 peer: %s %v", address, err)
	}
	if _, err := ParsePeerAddress("missing-port"); err == nil {
		t.Fatal("expected split error")
	}
	if _, err := ParsePeerAddress("hostname:80"); err == nil {
		t.Fatal("expected numeric IP error")
	}
}

func TestAttemptTrackerConcurrentDecisions(t *testing.T) {
	if _, err := NewAttemptTracker(0, delivery.AuthFailBan); err == nil {
		t.Fatal("expected threshold error")
	}
	if _, err := NewAttemptTracker(1, "invalid"); err == nil {
		t.Fatal("expected action error")
	}
	address := netip.MustParseAddr("192.0.2.2")
	tracker, err := NewAttemptTracker(4, delivery.AuthFailBan)
	if err != nil {
		t.Fatal(err)
	}
	var wait sync.WaitGroup
	decisions := make(chan FailureDecision, 4)
	for range 4 {
		wait.Add(1)
		go func() { defer wait.Done(); decisions <- tracker.Failure(testDeliveryID, address) }()
	}
	wait.Wait()
	close(decisions)
	transitions := 0
	for decision := range decisions {
		if decision.Transition {
			transitions++
		}
	}
	if transitions != 1 || !tracker.Blocked(testDeliveryID, address).Banned {
		t.Fatalf("expected one atomic ban, got %d", transitions)
	}
	if decision := tracker.Failure(testDeliveryID, address); decision.Transition || !decision.Banned {
		t.Fatal("an existing ban must not transition again")
	}
	tracker.Success(testDeliveryID, address)
	if tracker.Blocked(testDeliveryID, address).Attempts != 0 {
		t.Fatal("success must clear the fingerprint")
	}
	if err := tracker.Update(0, delivery.AuthFailBan); err == nil {
		t.Fatal("expected update threshold error")
	}
	if err := tracker.Update(1, "invalid"); err == nil {
		t.Fatal("expected update action error")
	}
	if err := tracker.Update(1, delivery.AuthFailStop); err != nil {
		t.Fatal(err)
	}
	decision := tracker.Failure(testDeliveryID, address)
	if !decision.Stop || !decision.Transition || !tracker.Blocked(testDeliveryID, address).Stop {
		t.Fatal("expected stop transition")
	}
	if next := tracker.Failure(testDeliveryID, address); next.Transition || !next.Stop {
		t.Fatal("stop must remain stable")
	}
}

func TestSessionsIsolationExpirationAndCookies(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	clock := func() time.Time { return now }
	if _, err := NewSessions(SessionConfig{}, randomBytes(), clock); err == nil {
		t.Fatal("expected invalid timeout error")
	}
	if _, err := NewSessions(SessionConfig{IdleTimeout: 2, AbsoluteTimeout: 1}, randomBytes(), clock); err == nil {
		t.Fatal("expected inverted timeout error")
	}
	if _, err := NewSessions(DefaultSessionConfig(), nil, clock); err == nil {
		t.Fatal("expected dependency error")
	}
	if _, err := NewSessions(DefaultSessionConfig(), randomBytes(), nil); err == nil {
		t.Fatal("expected clock error")
	}
	failing, err := NewSessions(DefaultSessionConfig(), failingReader{}, clock)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := failing.Issue(testDeliveryID); err == nil {
		t.Fatal("expected first token error")
	}
	shortRandom := bytes.NewReader(bytes.Repeat([]byte{1}, opaqueTokenBytes))
	secondFailure, err := NewSessions(DefaultSessionConfig(), shortRandom, clock)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := secondFailure.Issue(testDeliveryID); err == nil {
		t.Fatal("expected second token error")
	}
	sessions, err := NewSessions(SessionConfig{IdleTimeout: time.Minute, AbsoluteTimeout: 3 * time.Minute}, randomBytes(), clock)
	if err != nil {
		t.Fatal(err)
	}
	tokens, err := sessions.Issue(testDeliveryID)
	if err != nil {
		t.Fatal(err)
	}
	otherID := delivery.ID("00000000-0000-4000-8000-000000000023")
	if sessions.Validate(otherID, tokens.Session) || sessions.Validate(testDeliveryID, "missing") {
		t.Fatal("session must be delivery-scoped and opaque")
	}
	if !sessions.Validate(testDeliveryID, tokens.Session) {
		t.Fatal("expected valid session")
	}
	if sessions.ValidateCSRF(testDeliveryID, tokens.Session, "wrong") || sessions.ValidateCSRF(otherID, tokens.Session, tokens.CSRF) {
		t.Fatal("CSRF must be session and delivery scoped")
	}
	if !sessions.ValidateCSRF(testDeliveryID, tokens.Session, tokens.CSRF) {
		t.Fatal("expected valid CSRF")
	}
	now = now.Add(2 * time.Minute)
	if sessions.Validate(testDeliveryID, tokens.Session) {
		t.Fatal("expected idle expiration")
	}
	tokens, err = sessions.Issue(testDeliveryID)
	if err != nil {
		t.Fatal(err)
	}
	now = now.Add(4 * time.Minute)
	if sessions.ValidateCSRF(testDeliveryID, tokens.Session, tokens.CSRF) {
		t.Fatal("expected absolute expiration")
	}
	tokens, err = sessions.Issue(testDeliveryID)
	if err != nil {
		t.Fatal(err)
	}
	sessions.Revoke(tokens.Session)
	if sessions.Validate(testDeliveryID, tokens.Session) {
		t.Fatal("revoked session remained valid")
	}
	tokens, err = sessions.Issue(testDeliveryID)
	if err != nil {
		t.Fatal(err)
	}
	now = now.Add(4 * time.Minute)
	if removed := sessions.Purge(); removed != 1 || sessions.Purge() != 0 {
		t.Fatalf("unexpected purge count: %d", removed)
	}
	cookie := SessionCookie("courier", tokens.Session, true, time.Hour)
	if cookie.Name != "courier" || cookie.Path != "/" || !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteStrictMode || cookie.MaxAge != 3600 {
		t.Fatalf("unsafe cookie metadata: %#v", cookie)
	}
}

func TestReservationsAreAtomicAndBounded(t *testing.T) {
	unlimited := delivery.Limit{Unlimited: true}
	one := delivery.Limit{Value: 1}
	ten := delivery.Limit{Value: 10}
	if _, err := NewReservations(delivery.Limit{}, ten); err == nil {
		t.Fatal("expected invalid transfer limit")
	}
	reservations, err := NewReservations(one, ten)
	if err != nil {
		t.Fatal(err)
	}
	if err := reservations.Update(one, delivery.Limit{}); err == nil {
		t.Fatal("expected invalid file limit")
	}
	if _, err := reservations.Reserve(-2); err == nil {
		t.Fatal("expected negative declared size error")
	}
	if _, err := reservations.Reserve(11); !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected declared size rejection: %v", err)
	}
	reservation, err := reservations.Reserve(-1)
	if err != nil || reservations.Active() != 1 {
		t.Fatalf("expected reservation: %v", err)
	}
	if _, err := reservations.Reserve(0); !errors.Is(err, ErrTransferLimit) {
		t.Fatalf("expected aggregate limit: %v", err)
	}
	if err := reservation.Consume(-1); err == nil {
		t.Fatal("expected negative consume error")
	}
	if err := reservation.Consume(10); err != nil || reservation.Consumed() != 10 {
		t.Fatalf("unexpected consume: %v", err)
	}
	if err := reservation.Consume(1); !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected incremental size rejection: %v", err)
	}
	reservation.consumed = math.MaxInt64
	if err := reservation.Consume(1); !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected overflow rejection: %v", err)
	}
	reservation.Release()
	reservation.Release()
	if reservations.Active() != 0 || reservation.Consume(0) == nil {
		t.Fatal("release must be idempotent and terminal")
	}
	if err := reservations.Update(unlimited, unlimited); err != nil {
		t.Fatal(err)
	}
	declared, err := reservations.Reserve(2)
	if err != nil {
		t.Fatal(err)
	}
	if err := declared.Consume(3); !errors.Is(err, ErrFileTooLarge) {
		t.Fatalf("expected declared-size enforcement: %v", err)
	}
	declared.Release()
}

func TestAggregateRates(t *testing.T) {
	unlimited := delivery.Limit{Unlimited: true}
	if _, err := NewRates(delivery.Limit{}, unlimited); err == nil {
		t.Fatal("expected invalid rate")
	}
	rates, err := NewRates(unlimited, unlimited)
	if err != nil {
		t.Fatal(err)
	}
	if err := rates.Update(unlimited, delivery.Limit{}); err == nil {
		t.Fatal("expected download rate error")
	}
	if err := rates.Wait(nil, Upload, 1); err == nil {
		t.Fatal("expected context error")
	}
	if err := rates.Wait(context.Background(), Upload, -1); err == nil {
		t.Fatal("expected byte count error")
	}
	if err := rates.Wait(context.Background(), Direction("sideways"), 1); err == nil {
		t.Fatal("expected direction error")
	}
	if err := rates.Wait(context.Background(), Upload, 100); err != nil {
		t.Fatalf("unlimited rate failed: %v", err)
	}
	if err := rates.Update(delivery.Limit{Value: 100000}, delivery.Limit{Value: math.MaxInt64}); err != nil {
		t.Fatal(err)
	}
	first := rates.upload
	if err := rates.Wait(context.Background(), Upload, 0); err != nil {
		t.Fatal(err)
	}
	if err := rates.Wait(context.Background(), Upload, maximumRateBurst+1); err != nil {
		t.Fatal(err)
	}
	if err := rates.Update(delivery.Limit{Value: 200000}, unlimited); err != nil || rates.upload != first {
		t.Fatalf("expected limiter reconfiguration: %v", err)
	}
	if err := rates.Wait(context.Background(), Download, 1); err != nil {
		t.Fatal(err)
	}
	limited, err := NewRates(delivery.Limit{Value: 1}, unlimited)
	if err != nil {
		t.Fatal(err)
	}
	if err := limited.Wait(context.Background(), Upload, 1); err != nil {
		t.Fatal(err)
	}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if err := limited.Wait(canceled, Upload, 1); err == nil {
		t.Fatal("expected canceled rate wait")
	}
}

func TestEngineAuthorizationAndUpdates(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	clock := func() time.Time { return now }
	peer := testAddress("192.0.2.10:1234")
	if dependencies := DefaultDependencies(); dependencies.Random == nil || dependencies.Now == nil || dependencies.AuthConcurrency != 2 {
		t.Fatal("default dependencies are incomplete")
	}
	if _, err := NewEngine("invalid", testPolicy(delivery.AuthNone), Credentials{}, testDependencies(clock)); err == nil {
		t.Fatal("expected identity error")
	}
	invalidPolicy := testPolicy(delivery.AuthNone)
	invalidPolicy.AuthAttempts = 0
	if _, err := NewEngine(testDeliveryID, invalidPolicy, Credentials{}, testDependencies(clock)); err == nil {
		t.Fatal("expected policy error")
	}
	dependencies := testDependencies(clock)
	dependencies.AuthConcurrency = 0
	if _, err := NewEngine(testDeliveryID, testPolicy(delivery.AuthNone), Credentials{}, dependencies); err == nil {
		t.Fatal("expected authentication dependency error")
	}
	if _, err := NewEngine(testDeliveryID, testPolicy(delivery.AuthBasic), Credentials{}, testDependencies(clock)); err == nil {
		t.Fatal("expected missing Basic credentials")
	}
	if _, err := NewEngine(testDeliveryID, testPolicy(delivery.AuthPassword), Credentials{}, testDependencies(clock)); err == nil {
		t.Fatal("expected missing password credential")
	}
	dependencies = testDependencies(clock)
	dependencies.Sessions = SessionConfig{}
	if _, err := NewEngine(testDeliveryID, testPolicy(delivery.AuthNone), Credentials{}, dependencies); err == nil {
		t.Fatal("expected session dependency error")
	}
	mappedPolicy := testPolicy(delivery.AuthNone)
	mappedPolicy.AllowIP = []string{"::ffff:192.0.2.0/95"}
	if _, err := NewEngine(testDeliveryID, mappedPolicy, Credentials{}, testDependencies(clock)); err == nil {
		t.Fatal("expected canonical admission error")
	}

	configured := testPolicy(delivery.AuthBasic)
	configured.AuthAttempts = 2
	configured.AllowIP = []string{"192.0.2.0/24"}
	configured.DeliveryLimit = delivery.Limit{Value: 1}
	configured.MaxFileSize = delivery.Limit{Value: 3}
	engine, err := NewEngine(testDeliveryID, configured, Credentials{BasicUsername: "u", BasicPassword: []byte("p")}, testDependencies(clock))
	if err != nil {
		t.Fatal(err)
	}
	if got := EnforcementOrder(); !reflect.DeepEqual(got, []Stage{StageBounds, StageAdmission, StageBan, StageAuth, StageCSRF, StageReservation, StageRate, StageHandler}) {
		t.Fatalf("unexpected order: %v", got)
	}
	copyPolicy := engine.Policy()
	copyPolicy.AllowIP[0] = "modified"
	if engine.Policy().AllowIP[0] == "modified" {
		t.Fatal("policy leaked mutable allowlist")
	}
	if _, err := engine.Authorize(nil, Request{}); err == nil {
		t.Fatal("expected context error")
	}
	if _, err := engine.Authorize(context.Background(), Request{Incoming: true, DeclaredSize: -2}); err == nil {
		t.Fatal("expected bounds error")
	}
	if _, err := engine.Authorize(context.Background(), Request{Peer: testAddress("invalid")}); err == nil {
		t.Fatal("expected peer parsing error")
	}
	if _, err := engine.Authorize(context.Background(), Request{Peer: testAddress("203.0.113.1:1")}); !errors.Is(err, ErrPeerDenied) {
		t.Fatalf("expected peer denial: %v", err)
	}
	bad := Request{Peer: peer, Authentication: Authentication{Username: "u", BasicSecret: []byte("bad")}}
	if _, err := engine.Authorize(context.Background(), bad); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("expected authentication failure: %v", err)
	}
	if _, err := engine.Authorize(context.Background(), bad); !errors.Is(err, ErrBanned) {
		t.Fatalf("expected threshold ban: %v", err)
	}
	if _, err := engine.Authorize(context.Background(), bad); !errors.Is(err, ErrBanned) {
		t.Fatalf("expected existing ban: %v", err)
	}
	goodPeer := testAddress("192.0.2.11:2")
	good := Request{Peer: goodPeer, Incoming: true, DeclaredSize: 3, Authentication: Authentication{Username: "u", BasicSecret: []byte("p")}}
	authorization, err := engine.Authorize(context.Background(), good)
	if err != nil {
		t.Fatal(err)
	}
	if authorization.Session != nil || authorization.Consume(3) != nil || authorization.Consume(1) == nil {
		t.Fatal("unexpected Basic authorization reservation")
	}
	if err := authorization.Wait(context.Background(), Upload, 0); err != nil {
		t.Fatal(err)
	}
	if _, err := engine.Authorize(context.Background(), good); !errors.Is(err, ErrTransferLimit) {
		t.Fatalf("expected active reservation limit: %v", err)
	}
	authorization.Release()
	authorization.Release()
	nonIncoming, err := engine.Authorize(context.Background(), Request{Peer: goodPeer, Authentication: good.Authentication})
	if err != nil {
		t.Fatal(err)
	}
	if nonIncoming.Consume(0) == nil {
		t.Fatal("non-incoming authorization must not consume")
	}
	nonIncoming.Release()

	next := engine.Policy()
	next.Version++
	next.AllowIP = []string{"198.51.100.0/24"}
	if err := engine.Update(99, next); !errors.Is(err, delivery.ErrRevisionConflict) {
		t.Fatalf("expected revision conflict: %v", err)
	}
	changedMode := next
	changedMode.Auth = delivery.AuthPassword
	if err := engine.Update(configured.Version, changedMode); err == nil {
		t.Fatal("expected auth mode update rejection")
	}
	invalidNext := next
	invalidNext.AuthAttempts = 0
	if err := engine.Update(configured.Version, invalidNext); err == nil {
		t.Fatal("expected invalid update")
	}
	badAdmission := next
	badAdmission.AllowIP = []string{"::ffff:192.0.2.0/95"}
	if err := engine.Update(configured.Version, badAdmission); err == nil {
		t.Fatal("expected admission update rejection")
	}
	if err := engine.Update(configured.Version, next); err != nil {
		t.Fatal(err)
	}
	if _, err := engine.Authorize(context.Background(), good); !errors.Is(err, ErrPeerDenied) {
		t.Fatalf("updated admission not applied: %v", err)
	}

	passwordPolicy := testPolicy(delivery.AuthPassword)
	passwordPolicy.AuthAttempts = 1
	passwordPolicy.AuthFailAction = delivery.AuthFailStop
	passwordEngine, err := NewEngine(testDeliveryID, passwordPolicy, Credentials{Password: []byte("secret")}, testDependencies(clock))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := passwordEngine.Authorize(context.Background(), Request{Peer: peer, Authentication: Authentication{Password: []byte("wrong")}}); !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected stop decision: %v", err)
	}
	if _, err := passwordEngine.Authorize(context.Background(), Request{Peer: peer}); !errors.Is(err, ErrStopRequested) {
		t.Fatalf("expected retained stop decision: %v", err)
	}
	passwordEngine, err = NewEngine(testDeliveryID, testPolicy(delivery.AuthPassword), Credentials{Password: []byte("secret")}, testDependencies(clock))
	if err != nil {
		t.Fatal(err)
	}
	login, err := passwordEngine.Authorize(context.Background(), Request{Peer: peer, Authentication: Authentication{Password: []byte("secret")}})
	if err != nil || login.Session == nil {
		t.Fatalf("expected password session: %v", err)
	}
	if _, err := passwordEngine.Authorize(context.Background(), Request{Peer: peer, StateChanging: true, Authentication: Authentication{Session: login.Session.Session, CSRF: "wrong"}}); !errors.Is(err, ErrCSRF) {
		t.Fatalf("expected CSRF error: %v", err)
	}
	if _, err := passwordEngine.Authorize(context.Background(), Request{Peer: peer, StateChanging: true, Authentication: Authentication{Session: login.Session.Session, CSRF: login.Session.CSRF}}); err != nil {
		t.Fatalf("expected CSRF authorization: %v", err)
	}
	if _, err := passwordEngine.Authorize(context.Background(), Request{Peer: peer, Authentication: Authentication{Session: login.Session.Session}}); err != nil {
		t.Fatalf("expected session authorization: %v", err)
	}
	if _, err := passwordEngine.Authorize(context.Background(), Request{Peer: peer, Authentication: Authentication{Session: "invalid"}}); !errors.Is(err, ErrAuthentication) {
		t.Fatalf("expected invalid session authentication: %v", err)
	}
	passwordEngine.authenticator.slots <- struct{}{}
	canceled, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := passwordEngine.Authorize(canceled, Request{Peer: testAddress("192.0.2.12:3"), Authentication: Authentication{Password: []byte("secret")}}); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected canceled password verification: %v", err)
	}
	<-passwordEngine.authenticator.slots

	shortDependencies := testDependencies(clock)
	shortDependencies.Random = bytes.NewReader(bytes.Repeat([]byte{7}, 16))
	shortEngine, err := NewEngine(testDeliveryID, testPolicy(delivery.AuthPassword), Credentials{Password: []byte("secret")}, shortDependencies)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := shortEngine.Authorize(context.Background(), Request{Peer: peer, Authentication: Authentication{Password: []byte("secret")}}); err == nil {
		t.Fatal("expected session issuance failure")
	}

	openEngine, err := NewEngine(testDeliveryID, testPolicy(delivery.AuthNone), Credentials{}, testDependencies(clock))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := openEngine.Authorize(context.Background(), Request{Peer: peer}); err != nil {
		t.Fatalf("open policy rejected request: %v", err)
	}
}
