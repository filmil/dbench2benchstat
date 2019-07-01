package main

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestGetResult(t *testing.T) {
	tests := []struct {
		in  string
		out Result
	}{
		{
			in: `basic/increment -> avg 0.00927782162588819ms out of 2534 samples. (std dev 0.00161196365508829, min 0.008, max 0.082)`,
			out: Result{
				Name:       "basic/increment",
				Average:    time.Duration(9277 * time.Nanosecond),
				NumSamples: 2534,
				StdDev:     time.Duration(1611 * time.Nanosecond),
				Min:        time.Duration(8 * time.Microsecond),
				Max:        time.Duration(82 * time.Microsecond),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.in, func(t *testing.T) {
			res, err := GetResult(test.in)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(test.out, res) {
				t.Errorf("want: %+v, got: %+v", test.out, res)
			}
		})
	}
}
func TestProcess(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		out  string
	}{
		{
			name: "basic",
			in: []string{
				"basic/increment -> avg 0.00927782162588819ms out of 2534 samples. (std dev 0.00161196365508829, min 0.008, max 0.082)",
			},
			out: "basic/increment\t2534\t9277 ns/op\n",
		},
		{
			name: "complete",
			in: []string{
				"basic/increment -> avg 0.00927782162588819ms out of 2534 samples. (std dev 0.00161196365508829, min 0.008, max 0.082)",
				"basic/setInt -> avg 0.00841712707182296ms out of 2534 samples. (std dev 0.0006938950962008, min 0.008, max 0.018)",
				"basic/delNode -> avg 0.013989340702724ms out of 2533 samples. (std dev 0.00110354419501747, min 0.013, max 0.047)",
				"basic/addNode -> avg 0.0237091554853995ms out of 2534 samples. (std dev 0.00162248501181094, min 0.023, max 0.062)",
				"basic/setByte -> avg 0.0249545992893796ms out of 2533 samples. (std dev 0.00133135809515253, min 0.024, max 0.061)",
				"basic/resetByte -> avg 0.0291954204500601ms out of 2533 samples. (std dev 0.0019506613490976, min 0.028, max 0.068)",
				"basic/setByteLong -> avg 0.275580734307151ms out of 2533 samples. (std dev 0.0136697085424698, min 0.271, max 0.769)",
				"basic/resetByteLong -> avg 0.37006474536123ms out of 2533 samples. (std dev 0.0143529116400364, min 0.365, max 0.764)",
				"basic/setString -> avg 0.0195100671140934ms out of 2533 samples. (std dev 0.00258201386776676, min 0.018, max 0.058)",
				"basic/resetString -> avg 0.0214441373864983ms out of 2533 samples. (std dev 0.00123627010570329, min 0.021, max 0.058)",
				"basic/setStringLong -> avg 0.29572404263719ms out of 2533 samples. (std dev 0.0133516930974411, min 0.291, max 0.634)",
				"basic/resetStringLong -> avg 0.390995262534538ms out of 2533 samples. (std dev 0.0200222925814532, min 0.385, max 1.027)",
				"basic/incDouble -> avg 0.00923875295974771ms out of 2534 samples. (std dev 0.000715857828541273, min 0.008, max 0.023)",
				"basic/setDouble -> avg 0.00862352012628243ms out of 2534 samples. (std dev 0.00383421425772501, min 0.008, max 0.198)",
				"basic/alloc -> avg 0.0794449269640761ms out of 2533 samples. (std dev 0.00567592930146688, min 0.077, max 0.239)",
				"basic/url -> avg 0.0912913541255455ms out of 2533 samples. (std dev 0.00704704681367013, min 0.089, max 0.382)",
			},
			out: "basic/increment\t2534\t9277 ns/op\n" +
				"basic/setInt\t2534\t8417 ns/op\n" +
				"basic/delNode\t2533\t13989 ns/op\n" +
				"basic/addNode\t2534\t23709 ns/op\n" +
				"basic/setByte\t2533\t24954 ns/op\n" +
				"basic/resetByte\t2533\t29195 ns/op\n" +
				"basic/setByteLong\t2533\t275580 ns/op\n" +
				"basic/resetByteLong\t2533\t370064 ns/op\n" +
				"basic/setString\t2533\t19510 ns/op\n" +
				"basic/resetString\t2533\t21444 ns/op\n" +
				"basic/setStringLong\t2533\t295724 ns/op\n" +
				"basic/resetStringLong\t2533\t390995 ns/op\n" +
				"basic/incDouble\t2534\t9238 ns/op\n" +
				"basic/setDouble\t2534\t8623 ns/op\n" +
				"basic/alloc\t2533\t79444 ns/op\n" +
				"basic/url\t2533\t91291 ns/op\n" +
				"",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var b strings.Builder
			if err := Process(strings.NewReader(strings.Join(test.in, "\n")), &b); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if test.out != b.String() {
				t.Errorf("mismatch:\n\twant:\n%v\n\tgot:\n%v", test.out, b.String())
			}
		})
	}
}
