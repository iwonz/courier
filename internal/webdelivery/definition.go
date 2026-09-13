package webdelivery

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/iwonz/courier/internal/delivery"
	"github.com/iwonz/courier/internal/endpoint"
	"github.com/iwonz/courier/internal/operation"
	"github.com/iwonz/courier/internal/policy"
)

const (
	DefinitionVersion  = 1
	ResourceTokenBytes = 32
	maxEndpointPrompts = 32
	maxPromptBytes     = 4096
	maxSecretBytes     = 4096
)

type EndpointCredential struct {
	Prompt string `json:"prompt"`
	Secret []byte `json:"secret"`
}

type EndpointRuntime struct {
	Credentials []EndpointCredential `json:"credentials,omitempty"`
	AllowHelper bool                 `json:"allowHelper,omitempty"`
}

func (runtime *EndpointRuntime) Clear() {
	for index := range runtime.Credentials {
		clear(runtime.Credentials[index].Secret)
		runtime.Credentials[index] = EndpointCredential{}
	}
	runtime.Credentials = nil
	runtime.AllowHelper = false
}

func (runtime EndpointRuntime) Prompt(label string) ([]byte, error) {
	for _, credential := range runtime.Credentials {
		if credential.Prompt == label {
			return append([]byte(nil), credential.Secret...), nil
		}
	}
	return nil, errors.New("interactive SSH credential is unavailable in the delivery worker")
}

type Definition struct {
	Version     int                       `json:"version"`
	Token       string                    `json:"token"`
	Source      string                    `json:"source"`
	Destination string                    `json:"destination"`
	Archive     bool                      `json:"archive"`
	Extract     bool                      `json:"extract"`
	NoUI        bool                      `json:"noUi"`
	Selection   []operation.SelectionRule `json:"selection,omitempty"`
	Credentials policy.Credentials        `json:"credentials"`
	Endpoint    EndpointRuntime           `json:"endpointRuntime,omitempty"`
}

func (definition *Definition) ClearSecrets() {
	clear(definition.Credentials.BasicPassword)
	clear(definition.Credentials.Password)
	definition.Credentials = policy.Credentials{}
	definition.Endpoint.Clear()
}

func NewToken(random io.Reader) (string, error) {
	if random == nil {
		return "", errors.New("resource-token random source is required")
	}
	value := make([]byte, ResourceTokenBytes)
	if _, err := io.ReadFull(random, value); err != nil {
		return "", fmt.Errorf("generate resource token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func NewDefinition(plan operation.Plan, credentials policy.Credentials, endpointRuntime EndpointRuntime, random io.Reader) (Definition, error) {
	token, err := NewToken(random)
	if err != nil {
		return Definition{}, err
	}
	definition := Definition{
		Version: DefinitionVersion, Token: token, Source: plan.Source.Raw, Destination: plan.Destination.Raw,
		Archive: plan.Options.Archive, Extract: plan.Options.Extract, NoUI: plan.Options.NoUI,
		Selection: append([]operation.SelectionRule(nil), plan.Options.Selection...), Credentials: cloneCredentials(credentials),
		Endpoint: cloneEndpointRuntime(endpointRuntime),
	}
	if err := definition.Validate(routeForPlan(plan)); err != nil {
		definition.ClearSecrets()
		return Definition{}, err
	}
	return definition, nil
}

func routeForPlan(plan operation.Plan) delivery.Route {
	if plan.Route == operation.RouteWebToPath {
		return delivery.RouteWebToPath
	}
	if plan.Route == operation.RoutePathToWeb {
		return delivery.RoutePathToWeb
	}
	if plan.Route == operation.RouteWebhookToPath {
		return delivery.RouteWebhookToPath
	}
	return ""
}

func (definition Definition) Validate(route delivery.Route) error {
	if definition.Version != DefinitionVersion || strings.TrimSpace(definition.Source) != definition.Source || strings.TrimSpace(definition.Destination) != definition.Destination {
		return errors.New("invalid web delivery definition")
	}
	decoded, err := base64.RawURLEncoding.DecodeString(definition.Token)
	if err != nil || len(decoded) != ResourceTokenBytes {
		return errors.New("invalid web delivery resource token")
	}
	source, sourceErr := endpoint.Parse(definition.Source)
	destination, destinationErr := endpoint.Parse(definition.Destination)
	if sourceErr != nil || destinationErr != nil {
		return errors.Join(sourceErr, destinationErr)
	}
	switch route {
	case delivery.RouteWebToPath:
		if source.Kind != endpoint.KindWeb || !destination.IsPath() || definition.Archive {
			return errors.New("invalid browser-upload definition")
		}
	case delivery.RouteWebhookToPath:
		if source.Kind != endpoint.KindWebhook || !destination.IsPath() || definition.Archive || definition.NoUI {
			return errors.New("invalid incoming-webhook definition")
		}
	case delivery.RoutePathToWeb:
		if !source.IsPath() || destination.Kind != endpoint.KindWeb || definition.Extract {
			return errors.New("invalid browser-download definition")
		}
	default:
		return errors.New("unsupported web delivery route")
	}
	remote := source.Remote || destination.Remote
	if !remote && (len(definition.Endpoint.Credentials) != 0 || definition.Endpoint.AllowHelper) {
		return errors.New("local browser delivery contains SSH credentials")
	}
	if err := validateEndpointRuntime(definition.Endpoint); err != nil {
		return err
	}
	for _, rule := range definition.Selection {
		if rule.Value == "" {
			return errors.New("empty web delivery selection rule")
		}
	}
	return nil
}

func RandomDefinition(plan operation.Plan, credentials policy.Credentials) (Definition, error) {
	return NewDefinition(plan, credentials, EndpointRuntime{}, rand.Reader)
}

func cloneEndpointRuntime(runtime EndpointRuntime) EndpointRuntime {
	result := EndpointRuntime{Credentials: make([]EndpointCredential, len(runtime.Credentials)), AllowHelper: runtime.AllowHelper}
	for index, credential := range runtime.Credentials {
		result.Credentials[index] = EndpointCredential{Prompt: credential.Prompt, Secret: append([]byte(nil), credential.Secret...)}
	}
	return result
}

func cloneCredentials(credentials policy.Credentials) policy.Credentials {
	return policy.Credentials{
		BasicUsername: credentials.BasicUsername,
		BasicPassword: append([]byte(nil), credentials.BasicPassword...),
		Password:      append([]byte(nil), credentials.Password...),
	}
}

func validateEndpointRuntime(runtime EndpointRuntime) error {
	if len(runtime.Credentials) > maxEndpointPrompts {
		return errors.New("too many SSH credential prompts")
	}
	for _, credential := range runtime.Credentials {
		if strings.TrimSpace(credential.Prompt) == "" || len(credential.Prompt) > maxPromptBytes || strings.ContainsAny(credential.Prompt, "\x00\r\n") || len(credential.Secret) > maxSecretBytes {
			return errors.New("invalid SSH credential prompt")
		}
	}
	return nil
}
