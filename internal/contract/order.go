package contract

import "sort"

func requiredArgumentsFirst(values []Argument) []Argument {
	ordered := append([]Argument(nil), values...)
	sort.SliceStable(ordered, func(left, right int) bool {
		return requiredValue(ordered[left].Required) && !requiredValue(ordered[right].Required)
	})
	return ordered
}

func requiredFlagsFirst(values []Flag) []Flag {
	ordered := append([]Flag(nil), values...)
	sort.SliceStable(ordered, func(left, right int) bool {
		return requiredValue(ordered[left].Required) && !requiredValue(ordered[right].Required)
	})
	return ordered
}
