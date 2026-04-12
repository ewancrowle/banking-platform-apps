package luhn

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Generate creates a Luhn-valid numeric string of the specified length.
// The last digit is the calculated Luhn check digit.
func Generate(length int) (string, error) {
	if length < 2 {
		// 2 is the minimum length, since the last digit is the check digit
		return "", errors.New("length must be at least 2")
	}

	l := length - 1 // Exclude the check digit from the desired length
	digits := make([]int, l)

	for i := range l { // Iterate over the base length to generate random digits
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err // If rand.Read returns one
		}
		digits[i] = int(n.Int64())
	}

	check := CheckDigit(digits)
	complete := append(digits, check)

	return digitsToString(complete), nil
}

// CheckDigit computes the Luhn check digit for a given slice of digits.
func CheckDigit(digits []int) int {
	s := 0
	everyOther := true

	// Traverse from right to left (starting from the digit immediately before the check digit)
	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]
		if everyOther {
			d *= 2 // Double every other digit
			if d > 9 {
				d -= 9 // If it exceeds 9, subtract 9
			}
		}

		s += d
		everyOther = !everyOther
	}

	return (10 - (s % 10)) % 10 // // Smallest number to make a multiple of 10
}

// Validate verifies if a numeric string satisfies the Luhn algorithm.
func Validate(s string) bool {
	if len(s) < 2 {
		return false // 2 is the minimum length, since the last digit is the check digit
	}

	digits := make([]int, len(s))
	for i, r := range s { // Iterate through each rune in the numeric string
		if r < '0' || r > '9' {
			return false // If the rune is not a digit, return false
		}
		digits[i] = int(r - '0') // Convert the rune to an integer digit
	}

	// Extract the provided check digit and calculate what it should be
	v := digits[:len(digits)-1]
	check := digits[len(digits)-1]

	return check == CheckDigit(v) // Compare the check digit with the calculated value
}

// digitsToString converts a slice of integers to a string of digits.
func digitsToString(digits []int) string {
	res := make([]byte, len(digits)) // Allocate a byte slice of the same length as the digit slice
	for i, d := range digits {
		res[i] = byte('0' + d) // Convert each digit to a byte and store in the result slice
	}
	return string(res) // Convert the byte slice to a string and return
}
