package accountnum

import (
	"crypto/rand"
	"math/big"
)

// GenerateAccountNumber creates a valid 8-digit account number
// where the last digit is a Luhn check digit.
func GenerateAccountNumber() (string, error) {
	baseDigits := make([]int, 7)

	for i := 0; i < 7; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		baseDigits[i] = int(n.Int64())
	}

	checkDigit := calculateLuhnCheckDigit(baseDigits)
	fullNumber := append(baseDigits, checkDigit)

	return digitsToString(fullNumber), nil
}

// calculateLuhnCheckDigit computes the Luhn check digit for given digits.
func calculateLuhnCheckDigit(digits []int) int {
	sum := 0
	double := true

	// Traverse from right to left
	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]

		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}

		sum += d
		double = !double
	}

	return (10 - (sum % 10)) % 10
}

// ValidateAccountNumber verifies if an 8-digit number is Luhn-valid.
func ValidateAccountNumber(number string) bool {
	if len(number) != 8 {
		return false
	}

	digits := make([]int, 8)
	for i, ch := range number {
		if ch < '0' || ch > '9' {
			return false
		}
		digits[i] = int(ch - '0')
	}

	checkDigit := digits[7]
	expected := calculateLuhnCheckDigit(digits[:7])

	return checkDigit == expected
}

// digitsToString converts a slice of digits to a string.
func digitsToString(digits []int) string {
	result := make([]byte, len(digits))
	for i, d := range digits {
		result[i] = byte('0' + d)
	}
	return string(result)
}
