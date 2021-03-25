package base62

import (
	"testing"
)

var (
	tests = map[string]struct {
		inpNum  uint64
		outStr  string
		isError bool
		err     string
	}{
		"111111":                {111111, "h4C", false, ""},
		"284772472784":          {284772472784, "Ubrm0af", false, ""},
		"Error! unknown symbol": {0, "Ubrm0af.", true, "invalid character: ."},
	}
)

func TestDecoder(t *testing.T) {
	b := New()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := b.Decode(tc.outStr)

			if tc.isError && err == nil {
				t.Fatalf("Expected error: %v, but got nil", tc.err)
			}

			if tc.isError && err != nil {
				if tc.err != err.Error() {
					t.Fatalf("Expected errror: %v, but got: %v", tc.err, err)
				}
			} else {
				if res != tc.inpNum {
					t.Fatalf("Expected: %v, got: %v", tc.inpNum, tc.outStr)
				}
			}
		})
	}
}

func TestEncoder(t *testing.T) {

	b := New()

	for name, tc := range tests {
		if !tc.isError {
			t.Run(name, func(t *testing.T) {
				res := b.Encode(tc.inpNum)

				if res != tc.outStr {
					t.Fatalf("Expected: %v, got: %v", tc.inpNum, tc.outStr)

				}
			})
		}
	}
}
