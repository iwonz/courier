package contract

import (
	"fmt"
	"sort"

	"github.com/iwonz/courier/internal/operation"
)

// CheckPlanner verifies that runtime endpoint, route, and option definitions
// match the canonical operation portion of the contract.
func (c Contract) CheckPlanner(matrix operation.Matrix) error {
	contractEndpoints := make([]string, len(c.EndpointKinds))
	for index, value := range c.EndpointKinds {
		contractEndpoints[index] = value.Name
	}
	runtimeEndpoints := make([]string, len(matrix.EndpointKinds))
	for index, value := range matrix.EndpointKinds {
		runtimeEndpoints[index] = value.Kind
	}
	if err := equalValues("planner endpoint", contractEndpoints, runtimeEndpoints); err != nil {
		return err
	}

	runtimeRoutes := make(map[string]operation.RouteDefinition, len(matrix.Routes))
	for _, definition := range matrix.Routes {
		runtimeRoutes[definition.Route.String()] = definition
	}
	if len(runtimeRoutes) != len(matrix.Routes) {
		return fmt.Errorf("planner route definitions contain duplicates")
	}
	for _, route := range c.Routes {
		definition, ok := runtimeRoutes[route.Name]
		if !ok {
			return fmt.Errorf("planner route %q is missing", route.Name)
		}
		sources := make([]string, len(definition.Sources))
		for index, kind := range definition.Sources {
			sources[index] = kind.String()
		}
		destinations := make([]string, len(definition.Destinations))
		for index, kind := range definition.Destinations {
			destinations[index] = kind.String()
		}
		if err := equalValues("planner route "+route.Name+" source", route.Source, sources); err != nil {
			return err
		}
		if err := equalValues("planner route "+route.Name+" destination", route.Destination, destinations); err != nil {
			return err
		}
		allowed := make([]string, len(matrix.Allowed[definition.Route]))
		for index, name := range matrix.Allowed[definition.Route] {
			allowed[index] = string(name)
		}
		if err := equalValues("planner route "+route.Name+" option", route.AllowedFlags, allowed); err != nil {
			return err
		}
		delete(runtimeRoutes, route.Name)
	}
	if len(runtimeRoutes) != 0 {
		return fmt.Errorf("planner has undocumented routes")
	}

	operationScopes := map[string]bool{
		"path-to-path": true, "web-to-path": true, "path-to-web": true, "webhook-to-path": true, "path-to-http": true,
		"local-to-local": true, "local-to-ssh": true, "ssh-to-local": true, "ssh-to-ssh": true,
	}
	seenOptions := make(map[operation.OptionName]bool, len(c.Flags))
	for _, flag := range c.Flags {
		name := operation.OptionName(flag.Name)
		runtimeScopes, ok := matrix.Options[name]
		if !ok {
			return fmt.Errorf("planner option %q is missing", flag.Name)
		}
		contractScopes := make([]string, 0, len(flag.AppliesTo))
		for _, scope := range flag.AppliesTo {
			if operationScopes[scope] {
				contractScopes = append(contractScopes, scope)
			}
		}
		if err := equalValues("planner option "+flag.Name+" applicability", contractScopes, runtimeScopes); err != nil {
			return err
		}
		seenOptions[name] = true
	}
	for name := range matrix.Options {
		if !seenOptions[name] {
			return fmt.Errorf("planner option %q is undocumented", name)
		}
	}
	return nil
}

func equalValues(kind string, expected, actual []string) error {
	expected = append([]string(nil), expected...)
	actual = append([]string(nil), actual...)
	sort.Strings(expected)
	sort.Strings(actual)
	if len(expected) != len(actual) {
		return fmt.Errorf("%s mismatch: expected=%v actual=%v", kind, expected, actual)
	}
	for index := range expected {
		if expected[index] != actual[index] {
			return fmt.Errorf("%s mismatch: expected=%v actual=%v", kind, expected, actual)
		}
	}
	return nil
}
