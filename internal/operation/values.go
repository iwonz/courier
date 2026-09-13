package operation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// ErrRepeatedValue identifies a duplicate assignment to a single-use value.
var ErrRepeatedValue = errors.New("value may only be set once")

// SingleValue is a reusable pflag-compatible value that records explicit use.
type SingleValue struct {
	value string
	set   bool
}

// Set records the value and rejects a second occurrence.
func (v *SingleValue) Set(value string) error {
	if v.set {
		return ErrRepeatedValue
	}
	v.value = value
	v.set = true
	return nil
}

// String returns the current value.
func (v *SingleValue) String() string { return v.value }

// Type returns the value type shown by pflag.
func (v *SingleValue) Type() string { return "string" }

// Value returns the value and whether it was explicitly set.
func (v *SingleValue) Value() (string, bool) { return v.value, v.set }

// BoolValue is a pflag-compatible single-use boolean with bare-flag support.
type BoolValue struct {
	value bool
	set   bool
}

// Set parses and records one boolean occurrence.
func (v *BoolValue) Set(value string) error {
	if v.set {
		return ErrRepeatedValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	v.value = parsed
	v.set = true
	return nil
}

// String returns the current boolean value.
func (v *BoolValue) String() string { return strconv.FormatBool(v.value) }

// Type returns the value type shown by pflag.
func (v *BoolValue) Type() string { return "bool" }

// IsBoolFlag lets pflag accept the option without an explicit value.
func (v *BoolValue) IsBoolFlag() bool { return true }

// Value returns the value and whether it was explicitly set.
func (v *BoolValue) Value() (bool, bool) { return v.value, v.set }

// RepeatedValue is a reusable pflag-compatible ordered value list.
type RepeatedValue struct{ values []string }

// Set appends one occurrence.
func (v *RepeatedValue) Set(value string) error {
	v.values = append(v.values, value)
	return nil
}

// String renders the accumulated values for pflag diagnostics.
func (v *RepeatedValue) String() string { return strings.Join(v.values, ",") }

// Type returns the value type shown by pflag.
func (v *RepeatedValue) Type() string { return "strings" }

// Values returns a defensive copy in occurrence order.
func (v *RepeatedValue) Values() []string { return append([]string(nil), v.values...) }

// OrderedValues records interleaved occurrences from multiple option names.
type OrderedValues struct{ options []Option }

// For returns a pflag-compatible value bound to one option name.
func (v *OrderedValues) For(name OptionName) *OrderedValue {
	return &OrderedValue{name: name, owner: v}
}

// Options returns a defensive copy in global occurrence order.
func (v *OrderedValues) Options() []Option { return append([]Option(nil), v.options...) }

// OrderedValue binds one option name into an OrderedValues collection.
type OrderedValue struct {
	name  OptionName
	owner *OrderedValues
}

// Set appends one named occurrence to the shared collection.
func (v *OrderedValue) Set(value string) error {
	if v.owner == nil {
		return fmt.Errorf("ordered value has no owner")
	}
	v.owner.options = append(v.owner.options, Option{Name: v.name, Value: value})
	return nil
}

// String renders values belonging to this binding.
func (v *OrderedValue) String() string {
	if v.owner == nil {
		return ""
	}
	values := make([]string, 0)
	for _, option := range v.owner.options {
		if option.Name == v.name {
			values = append(values, option.Value)
		}
	}
	return strings.Join(values, ",")
}

// Type returns the value type shown by pflag.
func (v *OrderedValue) Type() string { return "strings" }
