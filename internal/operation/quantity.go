package operation

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
)

var quantityPattern = regexp.MustCompile(`^([0-9]+(?:[.][0-9]+)?)(B|kB|MB|GB|TB|KiB|MiB|GiB|TiB)(/s)?$`)

var quantityUnits = map[string]float64{
	"B": 1, "kB": 1_000, "MB": 1_000_000, "GB": 1_000_000_000, "TB": 1_000_000_000_000,
	"KiB": 1 << 10, "MiB": 1 << 20, "GiB": 1 << 30, "TiB": 1 << 40,
}

// Quantity is a byte size or byte-per-second rate. Unlimited is explicit so
// zero remains a meaningful parsed quantity rather than an implicit default.
type Quantity struct {
	Unlimited bool
	Rate      bool
	Value     int64
}

// ParseQuantity parses documented SI/IEC sizes and rates.
func ParseQuantity(value string, rate bool) (Quantity, error) {
	if value == "unlimited" {
		return Quantity{Unlimited: true, Rate: rate}, nil
	}
	parts := quantityPattern.FindStringSubmatch(value)
	if parts == nil || (parts[3] == "/s") != rate {
		return Quantity{}, fmt.Errorf("expected %s or unlimited", quantityKind(rate))
	}
	number, err := strconv.ParseFloat(parts[1], 64)
	bytes := number * quantityUnits[parts[2]]
	if err != nil || math.IsInf(bytes, 0) || bytes > math.MaxInt64 || bytes != math.Trunc(bytes) {
		return Quantity{}, fmt.Errorf("quantity is outside the supported byte range")
	}
	return Quantity{Rate: rate, Value: int64(bytes)}, nil
}

// MustParseQuantity parses a built-in default and panics if it is invalid.
func MustParseQuantity(value string, rate bool) Quantity {
	result, err := ParseQuantity(value, rate)
	if err != nil {
		panic(err)
	}
	return result
}

func quantityKind(rate bool) string {
	if rate {
		return "rate such as 50MiB/s"
	}
	return "size such as 10GiB"
}
