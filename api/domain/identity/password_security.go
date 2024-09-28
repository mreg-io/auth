package identity

import (
	"unicode"
)

func IsSecure(password string) bool {
	// Example criteria:
	// 1. At least 8 characters long block by interceptor
	// 2. At least one uppercase letter
	// 3. At least one lowercase letter
	// 4. At least one number
	// 5. At least one special character
	number, lower, upper, special := false, false, false, false
	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsLower(c):
			lower = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}
	return number && lower && upper && special
}
