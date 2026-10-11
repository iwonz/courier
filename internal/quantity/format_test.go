package quantity

import (
	"math"
	"testing"
)

func TestDataQuantityFormatting(t *testing.T) {
	for _, test := range []struct {
		name string
		got  string
		want string
	}{
		{"zero bytes", Bytes(0), "0 B"},
		{"grouped bytes", Bytes(1000), "1 000 B"},
		{"scaled bytes", Bytes(411789), "402.1 KiB"},
		{"mebibyte", Bytes(1 << 20), "1 MiB"},
		{"exbibyte", Bytes(1 << 60), "1 EiB"},
		{"bits", Bits(2048), "2 Kibit"},
		{"detailed", BytesDetailed(411789), "402.1 KiB (411 789 B)"},
		{"detailed base", BytesDetailed(42), "42 B"},
		{"exact", ExactBytes(411789), "411 789 B"},
		{"byte rate", BytesPerSecond(1536), "1.5 KiB/s"},
		{"base byte rate", BytesPerSecond(42), "42 B/s"},
		{"bit rate", BitsPerSecond(1024), "1 Kibit/s"},
		{"negative", Bytes(-1), "—"},
		{"negative rate", BytesPerSecond(-1), "—"},
		{"nan rate", BytesPerSecond(math.NaN()), "—"},
		{"infinite rate", BitsPerSecond(math.Inf(1)), "—"},
		{"grouped bits", Bits(1000), "1 000 bit"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("got %q want %q", test.got, test.want)
			}
		})
	}
	if grouped(-1) != "—" {
		t.Fatal("negative grouped quantity accepted")
	}
}
