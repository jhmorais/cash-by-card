package utils

import "fmt"

// ValidatePasswordPolicy exige min 8 caracteres com maiuscula, minuscula, numero e especial.
func ValidatePasswordPolicy(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("a senha deve ter pelo menos 8 caracteres, com letra maiúscula, minúscula, número e caractere especial")
	}
	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, r := range password {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return fmt.Errorf("a senha deve ter pelo menos 8 caracteres, com letra maiúscula, minúscula, número e caractere especial")
	}
	return nil
}
