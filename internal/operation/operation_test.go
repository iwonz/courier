package operation

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/iwonz/courier/internal/endpoint"
)

func TestRouteAndDirectionNames(t *testing.T) {
	for route, want := range map[Route]string{
		RoutePathToPath:    "path-to-path",
		RouteWebToPath:     "web-to-path",
		RoutePathToWeb:     "path-to-web",
		RouteWebhookToPath: "webhook-to-path",
		RoutePathToHTTP:    "path-to-http",
		Route(255):         "unknown",
	} {
		if got := route.String(); got != want {
			t.Errorf("route %d=%q want %q", route, got, want)
		}
	}
	for direction, want := range map[Direction]string{
		DirectionNone:         "none",
		DirectionLocalToLocal: "local-to-local",
		DirectionLocalToSSH:   "local-to-ssh",
		DirectionSSHToLocal:   "ssh-to-local",
		DirectionSSHToSSH:     "ssh-to-ssh",
		Direction(255):        "unknown",
	} {
		if got := direction.String(); got != want {
			t.Errorf("direction %d=%q want %q", direction, got, want)
		}
	}
}

func TestBuildRoutes(t *testing.T) {
	for _, test := range []struct {
		name        string
		source      string
		destination string
		route       Route
		direction   Direction
	}{
		{"local to local", "./source", "./destination", RoutePathToPath, DirectionLocalToLocal},
		{"local to SSH", "./source", "host:/destination", RoutePathToPath, DirectionLocalToSSH},
		{"SSH to local", "host:/source", "./destination", RoutePathToPath, DirectionSSHToLocal},
		{"SSH to SSH", "one:/source", "two:/destination", RoutePathToPath, DirectionSSHToSSH},
		{"web to local", "web://", "./destination", RouteWebToPath, DirectionNone},
		{"web to SSH", "web://", "host:/destination", RouteWebToPath, DirectionNone},
		{"local to web", "./source", "web://", RoutePathToWeb, DirectionNone},
		{"SSH to web", "host:/source", "web://", RoutePathToWeb, DirectionNone},
		{"webhook to local", "webhook://", "./destination", RouteWebhookToPath, DirectionNone},
		{"webhook to SSH", "webhook://", "host:/destination", RouteWebhookToPath, DirectionNone},
		{"local to HTTP", "./source", "http://example.test/hook", RoutePathToHTTP, DirectionNone},
		{"SSH to HTTPS", "host:/source", "https://example.test/hook", RoutePathToHTTP, DirectionNone},
	} {
		t.Run(test.name, func(t *testing.T) {
			plan, err := Build(Request{Source: test.source, Destination: test.destination})
			if err != nil || plan.Route != test.route || plan.Direction != test.direction {
				t.Fatalf("plan=%+v err=%v", plan, err)
			}
			if plan.Source.Raw != test.source || plan.Destination.Raw != test.destination {
				t.Fatalf("endpoints=%+v %+v", plan.Source, plan.Destination)
			}
		})
	}
}

func TestBuildRouteErrors(t *testing.T) {
	for _, test := range []struct {
		name        string
		source      string
		destination string
		class       error
		field       string
	}{
		{"invalid source", "ftp://example.test", "./out", ErrInvalidEndpoint, "source"},
		{"invalid destination", "./in", "ftp://example.test", ErrInvalidEndpoint, "destination"},
		{"HTTP source", "https://example.test/file", "./out", ErrUnsupportedRoute, ""},
		{"web to web", "web://", "web://", ErrUnsupportedRoute, ""},
		{"webhook to HTTP", "webhook://", "https://example.test", ErrUnsupportedRoute, ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Build(Request{Source: test.source, Destination: test.destination})
			var operationError *Error
			if !errors.Is(err, test.class) || !errors.As(err, &operationError) || operationError.Field != test.field {
				t.Fatalf("error=%v", err)
			}
			if !strings.Contains(err.Error(), test.class.Error()) {
				t.Fatalf("error string=%q", err)
			}
		})
	}
}

func TestBuildEffectiveOptions(t *testing.T) {
	request := Request{
		Source:      "web://",
		Destination: "./uploads",
		Options: []Option{
			{Name: OptionExtract},
			{Name: OptionListen, Value: "[::1]:9000"},
			{Name: OptionBackground},
			{Name: OptionAuth, Value: "basic"},
			{Name: OptionAuthAttempts, Value: "7"},
			{Name: OptionAuthFailAction, Value: "stop"},
			{Name: OptionLimit, Value: "3"},
			{Name: OptionAllowIP, Value: "192.0.2.42"},
			{Name: OptionAllowIP, Value: "2001:db8::/64"},
			{Name: OptionExclude, Value: "*.tmp"},
			{Name: OptionExcludeFrom, Value: "rules.txt"},
			{Name: OptionExcludeRegex, Value: `^private/`},
			{Name: OptionMaxFileSize, Value: "2GB"},
			{Name: OptionMaxExtractedSize, Value: "3GiB"},
			{Name: OptionUploadRate, Value: "50MiB/s"},
		},
	}
	plan, err := Build(request)
	if err != nil {
		t.Fatal(err)
	}
	options := plan.Options
	if !options.Extract || options.Listen != "[::1]:9000" || !options.Background || options.Auth != AuthBasic || options.AuthAttempts != 7 || options.AuthFailAction != AuthFailStop {
		t.Fatalf("scalar options=%+v", options)
	}
	if options.Limit.Unlimited || options.Limit.Value != 3 || options.NoUI || options.MaxFileSize.Value != 2_000_000_000 || options.MaxExtractedSize.Value != 3<<30 || options.UploadRate.Value != 50<<20 || !options.DownloadRate.Unlimited {
		t.Fatalf("limit/quantity options=%+v", options)
	}
	wantNetworks := []string{"192.0.2.42/32", "2001:db8::/64"}
	for index, prefix := range options.AllowIP {
		if prefix.String() != wantNetworks[index] {
			t.Fatalf("prefixes=%v", options.AllowIP)
		}
	}
	wantSelection := []SelectionRule{
		{Kind: SelectionGitignore, Value: "*.tmp", Position: 9},
		{Kind: SelectionFile, Value: "rules.txt", Position: 10},
		{Kind: SelectionRegex, Value: `^private/`, Position: 11},
	}
	if !reflect.DeepEqual(options.Selection, wantSelection) || len(options.Occurrences) != len(request.Options) {
		t.Fatalf("selection=%+v occurrences=%+v", options.Selection, options.Occurrences)
	}
	if !options.Explicit(OptionLimit) || options.Explicit(OptionNoUI) {
		t.Fatalf("explicit metadata=%+v", options)
	}
}

func TestRouteSpecificRateOptions(t *testing.T) {
	for _, test := range []struct {
		name        string
		source      string
		destination string
		options     []Option
	}{
		{"local to SSH upload", "./source", "host:/out", []Option{{Name: OptionUploadRate, Value: "1MB/s"}}},
		{"SSH to local download", "host:/source", "./out", []Option{{Name: OptionDownloadRate, Value: "1MB/s"}}},
		{"SSH to SSH both", "one:/source", "two:/out", []Option{{Name: OptionUploadRate, Value: "1MB/s"}, {Name: OptionDownloadRate, Value: "2MB/s"}}},
		{"path to web download", "./source", "web://", []Option{{Name: OptionArchive}, {Name: OptionListen, Value: "127.0.0.1:8088"}, {Name: OptionBackground, Value: "false"}, {Name: OptionAuth, Value: "password"}, {Name: OptionAuthAttempts, Value: "5"}, {Name: OptionAuthFailAction, Value: "ban"}, {Name: OptionLimit, Value: "unlimited"}, {Name: OptionNoUI}, {Name: OptionAllowIP, Value: "127.0.0.1/8"}, {Name: OptionDownloadRate, Value: "unlimited"}}},
		{"webhook upload", "webhook://", "./out", []Option{{Name: OptionExtract}, {Name: OptionAuth, Value: "basic"}, {Name: OptionUploadRate, Value: "1kB/s"}}},
		{"path to HTTP upload", "./source", "https://example.test/hook", []Option{{Name: OptionArchive}, {Name: OptionAuth, Value: "basic"}, {Name: OptionUploadRate, Value: "1kB/s"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			if _, err := Build(Request{Source: test.source, Destination: test.destination, Options: test.options}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestOptionValidationErrors(t *testing.T) {
	validWeb := Request{Source: "web://", Destination: "./out"}
	for _, test := range []struct {
		name    string
		request Request
		option  OptionName
	}{
		{"unknown", withOptions(validWeb, Option{Name: "missing"}), "missing"},
		{"duplicate", withOptions(validWeb, Option{Name: OptionBackground}, Option{Name: OptionBackground}), OptionBackground},
		{"not applicable", Request{Source: "./in", Destination: "./out", Options: []Option{{Name: OptionListen, Value: "127.0.0.1:80"}}}, OptionListen},
		{"all is control only", withOptions(validWeb, Option{Name: OptionAll}), OptionAll},
		{"invalid archive bool", Request{Source: "./in", Destination: "./out", Options: []Option{{Name: OptionArchive, Value: "sometimes"}}}, OptionArchive},
		{"invalid extract bool", withOptions(validWeb, Option{Name: OptionExtract, Value: "sometimes"}), OptionExtract},
		{"invalid listen", withOptions(validWeb, Option{Name: OptionListen, Value: "localhost"}), OptionListen},
		{"zero listen port", withOptions(validWeb, Option{Name: OptionListen, Value: "localhost:0"}), OptionListen},
		{"invalid auth", withOptions(validWeb, Option{Name: OptionAuth, Value: "token"}), OptionAuth},
		{"invalid attempts", withOptions(validWeb, Option{Name: OptionAuthAttempts, Value: "0"}), OptionAuthAttempts},
		{"invalid fail action", withOptions(validWeb, Option{Name: OptionAuthFailAction, Value: "ignore"}), OptionAuthFailAction},
		{"invalid limit", withOptions(validWeb, Option{Name: OptionLimit, Value: "many"}), OptionLimit},
		{"invalid no-ui bool", Request{Source: "./in", Destination: "web://", Options: []Option{{Name: OptionNoUI, Value: "sometimes"}}}, OptionNoUI},
		{"invalid IP", withOptions(validWeb, Option{Name: OptionAllowIP, Value: "not-an-ip"}), OptionAllowIP},
		{"empty exclude", withOptions(validWeb, Option{Name: OptionExclude}), OptionExclude},
		{"invalid regex", withOptions(validWeb, Option{Name: OptionExcludeRegex, Value: "["}), OptionExcludeRegex},
		{"empty exclude file", withOptions(validWeb, Option{Name: OptionExcludeFrom}), OptionExcludeFrom},
		{"invalid max file size", withOptions(validWeb, Option{Name: OptionMaxFileSize, Value: "large"}), OptionMaxFileSize},
		{"invalid max extracted size", withOptions(validWeb, Option{Name: OptionExtract}, Option{Name: OptionMaxExtractedSize, Value: "large"}), OptionMaxExtractedSize},
		{"invalid upload rate", withOptions(validWeb, Option{Name: OptionUploadRate, Value: "1MiB"}), OptionUploadRate},
		{"invalid download rate", Request{Source: "./in", Destination: "web://", Options: []Option{{Name: OptionDownloadRate, Value: "1MiB"}}}, OptionDownloadRate},
		{"archive extract conflict", Request{Source: "./in", Destination: "./out", Options: []Option{{Name: OptionArchive}, {Name: OptionExtract}}}, OptionArchive},
		{"extracted limit dependency", Request{Source: "./in", Destination: "./out", Options: []Option{{Name: OptionMaxExtractedSize, Value: "1GiB"}}}, OptionMaxExtractedSize},
		{"password without browser", Request{Source: "./in", Destination: "https://example.test", Options: []Option{{Name: OptionAuth, Value: "password"}}}, OptionAuth},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := Build(test.request)
			var operationError *Error
			if !errors.Is(err, ErrInvalidOption) || !errors.As(err, &operationError) || operationError.Field != "--"+string(test.option) {
				t.Fatalf("error=%v", err)
			}
		})
	}

	options := defaultOptions()
	err := applyOption(&options, Occurrence{Option: Option{Name: OptionAll}})
	if !errors.Is(err, ErrInvalidOption) {
		t.Fatalf("unhandled known option error=%v", err)
	}
}

func withOptions(request Request, options ...Option) Request {
	request.Options = options
	return request
}

func TestQuantityParsing(t *testing.T) {
	for _, test := range []struct {
		value     string
		rate      bool
		want      int64
		unlimited bool
	}{
		{"unlimited", false, 0, true},
		{"unlimited", true, 0, true},
		{"0B", false, 0, false},
		{"1.5KiB", false, 1536, false},
		{"2TB", false, 2_000_000_000_000, false},
		{"3TiB", false, 3 << 40, false},
		{"50MiB/s", true, 50 << 20, false},
	} {
		got, err := ParseQuantity(test.value, test.rate)
		if err != nil || got.Value != test.want || got.Unlimited != test.unlimited || got.Rate != test.rate {
			t.Errorf("ParseQuantity(%q,%t)=%+v,%v", test.value, test.rate, got, err)
		}
	}
	for _, test := range []struct {
		value string
		rate  bool
	}{
		{"1MiB/s", false},
		{"1MiB", true},
		{"0.1B", false},
		{strings.Repeat("9", 400) + "TB", false},
		{"-1MiB", false},
	} {
		if _, err := ParseQuantity(test.value, test.rate); err == nil {
			t.Errorf("expected quantity error for %q", test.value)
		}
	}
	if got := MustParseQuantity("1kB", false); got.Value != 1000 {
		t.Fatalf("must parse=%+v", got)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("expected MustParseQuantity panic")
		}
	}()
	MustParseQuantity("invalid", false)
}

func TestOperationErrorFormatting(t *testing.T) {
	cause := fmt.Errorf("detail")
	withoutField := (&Error{Class: ErrUnsupportedRoute, Cause: cause}).Error()
	withField := (&Error{Class: ErrInvalidOption, Field: "--archive", Cause: cause}).Error()
	if withoutField != "unsupported operation route: detail" || withField != `invalid operation option "--archive": detail` {
		t.Fatalf("without=%q with=%q", withoutField, withField)
	}
	if !errors.Is(&Error{Class: ErrInvalidOption, Cause: cause}, cause) {
		t.Fatal("underlying cause is not exposed")
	}
}

func TestLegacyRemoteDirection(t *testing.T) {
	source := endpoint.Endpoint{Remote: true}
	destination := endpoint.Endpoint{Remote: true}
	if direction := pathDirection(source.Kind, destination.Kind); direction != DirectionLocalToLocal {
		t.Fatalf("typed direction unexpectedly consumed legacy fields: %s", direction)
	}
}

func TestContractMatrix(t *testing.T) {
	matrix := ContractMatrix()
	if len(matrix.EndpointKinds) != 5 || len(matrix.Routes) != 5 || len(matrix.Options) != len(optionDefinitions) || len(matrix.Allowed) != 5 {
		t.Fatalf("matrix dimensions=%+v", matrix)
	}
	if matrix.EndpointKinds[1].Kind != "ssh" || matrix.Routes[0].Route != RoutePathToPath || matrix.Routes[0].Sources[1] != endpoint.KindSSH {
		t.Fatalf("matrix identity=%+v", matrix)
	}
	if !reflect.DeepEqual(matrix.Options[OptionUploadRate], []string{"local-to-ssh", "path-to-http", "ssh-to-ssh", "web-to-path", "webhook-to-path"}) {
		t.Fatalf("upload scopes=%v", matrix.Options[OptionUploadRate])
	}
	if !reflect.DeepEqual(matrix.Allowed[RoutePathToPath], []OptionName{OptionArchive, OptionDownloadRate, OptionExclude, OptionExcludeFrom, OptionExcludeRegex, OptionExtract, OptionMaxExtractedSize, OptionUploadRate}) {
		t.Fatalf("path options=%v", matrix.Allowed[RoutePathToPath])
	}
	matrix.Routes[0].Sources[0] = endpoint.KindHTTP
	matrix.Routes[0].Destinations[0] = endpoint.KindHTTP
	matrix.Options[OptionArchive][0] = "changed"
	matrix.Allowed[RoutePathToPath][0] = OptionAll
	fresh := ContractMatrix()
	if fresh.Routes[0].Sources[0] != endpoint.KindLocal || fresh.Routes[0].Destinations[0] != endpoint.KindLocal || fresh.Options[OptionArchive][0] != "path-to-http" || fresh.Allowed[RoutePathToPath][0] != OptionArchive {
		t.Fatalf("matrix did not return defensive data: %+v", fresh)
	}
}
