package operation

import (
	"errors"
	"reflect"
	"testing"
)

func TestSingleValue(t *testing.T) {
	var value SingleValue
	if value.String() != "" || value.Type() != "string" {
		t.Fatalf("initial value=%q type=%q", value.String(), value.Type())
	}
	if actual, set := value.Value(); actual != "" || set {
		t.Fatalf("initial Value=%q,%t", actual, set)
	}
	if err := value.Set("first"); err != nil {
		t.Fatal(err)
	}
	if actual, set := value.Value(); actual != "first" || !set || value.String() != "first" {
		t.Fatalf("Value=%q,%t String=%q", actual, set, value.String())
	}
	if err := value.Set("second"); !errors.Is(err, ErrRepeatedValue) {
		t.Fatalf("duplicate error=%v", err)
	}
}

func TestBoolValue(t *testing.T) {
	var value BoolValue
	if value.String() != "false" || value.Type() != "bool" || !value.IsBoolFlag() {
		t.Fatalf("initial bool=%q type=%q", value.String(), value.Type())
	}
	if actual, set := value.Value(); actual || set {
		t.Fatalf("initial Value=%t,%t", actual, set)
	}
	if err := value.Set("invalid"); err == nil {
		t.Fatal("expected invalid bool")
	}
	if err := value.Set("true"); err != nil {
		t.Fatal(err)
	}
	if actual, set := value.Value(); !actual || !set || value.String() != "true" {
		t.Fatalf("Value=%t,%t String=%q", actual, set, value.String())
	}
	if err := value.Set("false"); !errors.Is(err, ErrRepeatedValue) {
		t.Fatalf("duplicate error=%v", err)
	}
}

func TestRepeatedValue(t *testing.T) {
	var value RepeatedValue
	if value.Type() != "strings" || value.String() != "" || value.Values() != nil {
		t.Fatalf("initial repeated value=%+v", value)
	}
	if err := value.Set("one"); err != nil {
		t.Fatal(err)
	}
	if err := value.Set("two"); err != nil {
		t.Fatal(err)
	}
	values := value.Values()
	values[0] = "changed"
	if value.String() != "one,two" || !reflect.DeepEqual(value.Values(), []string{"one", "two"}) {
		t.Fatalf("repeated value=%q values=%v", value.String(), value.Values())
	}
}

func TestOrderedValues(t *testing.T) {
	var values OrderedValues
	exclude := values.For(OptionExclude)
	from := values.For(OptionExcludeFrom)
	if exclude.Type() != "strings" || exclude.String() != "" {
		t.Fatalf("empty binding type=%q value=%q", exclude.Type(), exclude.String())
	}
	for _, call := range []struct {
		value *OrderedValue
		text  string
	}{{exclude, "one"}, {from, "rules"}, {exclude, "two"}} {
		if err := call.value.Set(call.text); err != nil {
			t.Fatal(err)
		}
	}
	want := []Option{{Name: OptionExclude, Value: "one"}, {Name: OptionExcludeFrom, Value: "rules"}, {Name: OptionExclude, Value: "two"}}
	got := values.Options()
	got[0].Value = "changed"
	if exclude.String() != "one,two" || from.String() != "rules" || !reflect.DeepEqual(values.Options(), want) {
		t.Fatalf("exclude=%q from=%q options=%v", exclude.String(), from.String(), values.Options())
	}
	orphan := &OrderedValue{name: OptionExclude}
	if orphan.String() != "" || orphan.Set("x") == nil {
		t.Fatal("orphan binding should be inert and reject Set")
	}
}
