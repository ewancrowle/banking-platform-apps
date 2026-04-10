package luhn_test

import (
	"account-svc/internal/luhn"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		length  int
		want    string
		wantErr bool
	}{
		// Test case for a standard length (e.g., 8 digits)
		{
			name:    "StandardLength8",
			length:  8,
			want:    "", // We only check length/error here, content is random
			wantErr: false,
		},
		// Test case for minimum valid length (e.g., 2 digits)
		{
			name:    "MinimumLength2",
			length:  2,
			want:    "",
			wantErr: false,
		},
		// Test case for zero length (should error)
		{
			name:    "ZeroLength",
			length:  0,
			want:    "",
			wantErr: true,
		},
		// Test case for negative length (should error)
		{
			name:    "NegativeLength",
			length:  -5,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := luhn.Generate(tt.length)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Generate() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Generate() succeeded unexpectedly")
			}
			// Check if the generated string has the expected length
			if got != "" {
				if len(got) != tt.length {
					t.Errorf("Generate() = %s, want length %d, got length %d", got, tt.length, len(got))
				}
			}
		})
	}
}

func TestCheckDigit(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		digits []int
		want   int
	}{
		{
			name:   "CheckDigit8",
			digits: []int{4, 9, 9, 2, 7, 3, 9, 8, 7, 1, 6},
			want:   8,
		},
		{
			name:   "CheckDigit6",
			digits: []int{4, 9, 9, 2, 7, 3, 9, 8, 7, 1, 7},
			want:   6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := luhn.CheckDigit(tt.digits)
			if got != tt.want {
				t.Errorf("CheckDigit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want bool
	}{
		{
			name: "ValidNumber",
			s:    "49927398716",
			want: true,
		},
		{
			name: "InvalidNumber",
			s:    "49927398717", // Changed last digit from 6 to 7
			want: false,
		},
		{
			name: "EmptyString",
			s:    "",
			want: false,
		},
		{
			name: "NonDigitChars",
			s:    "499A7398716",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := luhn.Validate(tt.s)
			if got != tt.want {
				t.Errorf("Validate() = %t, want %t", got, tt.want)
			}
		})
	}
}
