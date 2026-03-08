package domain_test

import (
	"strings"
	"testing"

	"github.com/akemoon/crowdfunding-app-user/lib/domain"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// valid
		{name: "min length", input: "abc", wantErr: false},
		{name: "max length", input: strings.Repeat("a", 32), wantErr: false},
		{name: "with digits", input: "user123", wantErr: false},
		{name: "with dash in middle", input: "user-name", wantErr: false},
		{name: "digits only", input: "123", wantErr: false},
		{name: "mixed dash and digits", input: "a1-b2", wantErr: false},

		// length
		{name: "empty", input: "", wantErr: true},
		{name: "one char", input: "a", wantErr: true},
		{name: "two chars", input: "ab", wantErr: true},
		{name: "too long", input: strings.Repeat("a", 33), wantErr: true},

		// invalid characters
		{name: "uppercase", input: "Username", wantErr: true},
		{name: "space", input: "user name", wantErr: true},
		{name: "underscore", input: "user_name", wantErr: true},
		{name: "dot", input: "user.name", wantErr: true},
		{name: "at sign", input: "user@name", wantErr: true},
		{name: "cyrillic", input: "пользователь", wantErr: true},
		{name: "emoji", input: "user😀", wantErr: true},

		// dash position
		{name: "leading dash", input: "-username", wantErr: true},
		{name: "trailing dash", input: "username-", wantErr: true},
		{name: "only dashes", input: "---", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := domain.ValidateUsername(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateUsername(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			}
			if err != nil && test.wantErr {
				if err != domain.ErrInvalidUsername {
					t.Errorf("ValidateUsername(%q) = %v, want ErrInvalidUsername", test.input, err)
				}
			}
		})
	}
}
