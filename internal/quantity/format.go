// Package quantity formats human-facing binary data quantities without
// confusing bits and bytes. Machine-readable output should retain raw values.
package quantity

import (
	"math"
	"strconv"
	"strings"
)

var byteUnits = [...]string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
var bitUnits = [...]string{"bit", "Kibit", "Mibit", "Gibit", "Tibit", "Pibit", "Eibit"}

// Bytes formats a non-negative byte count using binary byte units.
func Bytes(value int64) string { return formatInteger(value, byteUnits[:]) }

// Bits formats a non-negative bit count using binary bit units.
func Bits(value int64) string { return formatInteger(value, bitUnits[:]) }

// ExactBytes formats an exact byte count with space-grouped decimal digits.
func ExactBytes(value int64) string { return grouped(value) + " B" }

// BytesDetailed adds an exact base-unit value when Bytes scales the number.
func BytesDetailed(value int64) string {
	compact := Bytes(value)
	exact := ExactBytes(value)
	if compact == exact {
		return exact
	}
	return compact + " (" + exact + ")"
}

// BytesPerSecond formats a non-negative byte rate.
func BytesPerSecond(value float64) string { return formatRate(value, byteUnits[:]) }

// BitsPerSecond formats a non-negative bit rate.
func BitsPerSecond(value float64) string { return formatRate(value, bitUnits[:]) }

func formatRate(value float64, units []string) string {
	formatted := formatFloat(value, units)
	if formatted == "—" {
		return formatted
	}
	return formatted + "/s"
}

func formatInteger(value int64, units []string) string {
	if value < 0 {
		return "—"
	}
	if value < 1024 {
		return grouped(value) + " " + units[0]
	}
	return formatFloat(float64(value), units)
}

func formatFloat(value float64, units []string) string {
	if value < 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return "—"
	}
	unit := 0
	for value >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}
	if unit == 0 {
		return grouped(int64(math.Round(value))) + " " + units[unit]
	}
	formatted := strconv.FormatFloat(value, 'f', 1, 64)
	formatted = strings.TrimSuffix(formatted, ".0")
	return formatted + " " + units[unit]
}

func grouped(value int64) string {
	if value < 0 {
		return "—"
	}
	digits := strconv.FormatInt(value, 10)
	first := len(digits) % 3
	if first == 0 {
		first = 3
	}
	var result strings.Builder
	result.Grow(len(digits) + len(digits)/3)
	result.WriteString(digits[:first])
	for index := first; index < len(digits); index += 3 {
		result.WriteByte(' ')
		result.WriteString(digits[index : index+3])
	}
	return result.String()
}
