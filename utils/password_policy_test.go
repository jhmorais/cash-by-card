package utils

import "testing"

func TestValidatePasswordPolicy(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valida completa", "NovaSenha@1", false},
		{"valida com simbolos variados", "Abc123!x", false},
		{"curta", "Aa1!aa", true},
		{"sem maiuscula", "novasenha@1", true},
		{"sem minuscula", "NOVASENHA@1", true},
		{"sem numero", "NovaSenha@!", true},
		{"sem especial", "NovaSenha12", true},
		{"vazia", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := ValidatePasswordPolicy(c.password)
			if (err != nil) != c.wantErr {
				t.Fatalf("ValidatePasswordPolicy(%q) err='%v', wantErr=%v", c.password, err, c.wantErr)
			}
		})
	}
}
