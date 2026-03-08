package domain_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/akemoon/crowdfunding-app-user/modules/user/domain"
)

func TestValidateDisplayName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// empty is allowed
		{name: "empty", input: "", wantErr: false},

		// valid
		{name: "latin word", input: "John", wantErr: false},
		{name: "cyrillic word", input: "Иван", wantErr: false},
		{name: "latin with space", input: "John Doe", wantErr: false},
		{name: "cyrillic with space", input: "Иван Иванов", wantErr: false},
		{name: "digits", input: "User 42", wantErr: false},
		{name: "punctuation", input: "O'Brien", wantErr: false},
		{name: "mixed scripts", input: "Ivan Иванов", wantErr: false},
		{name: "max length", input: strings.Repeat("a", domain.MaxDisplayNameLength), wantErr: false},
		{name: "max length cyrillic", input: strings.Repeat("а", domain.MaxDisplayNameLength), wantErr: false},

		// length
		{name: "over max length", input: strings.Repeat("a", domain.MaxDisplayNameLength+1), wantErr: true},
		{name: "over max length cyrillic runes", input: strings.Repeat("а", domain.MaxDisplayNameLength+1), wantErr: true},

		// whitespace
		{name: "leading space", input: " John", wantErr: true},
		{name: "trailing space", input: "John ", wantErr: true},
		{name: "tab character", input: "John\tDoe", wantErr: true},
		{name: "newline", input: "John\nDoe", wantErr: true},

		// invalid characters
		{name: "emoji", input: "John😀", wantErr: true},
		{name: "chinese", input: "约翰", wantErr: true},
		{name: "arabic", input: "يوحنا", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := domain.ValidateDisplayName(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateDisplayName(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			}
			if err != nil && test.wantErr {
				if err != domain.ErrInvalidDisplayName {
					t.Errorf("ValidateDisplayName(%q) = %v, want ErrInvalidDisplayName", test.input, err)
				}
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	longDesc := strings.Repeat("a", domain.MaxDescriptionLength)
	// make sure it's exactly at max in runes
	_ = utf8.RuneCountInString(longDesc)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// empty is allowed
		{name: "empty", input: "", wantErr: false},

		// valid
		{name: "simple latin", input: "Hello world", wantErr: false},
		{name: "cyrillic", input: "Привет мир", wantErr: false},
		{name: "with newline", input: "Line one\nLine two", wantErr: false},
		{name: "digits and punct", input: "Version 2.0: released!", wantErr: false},
		{name: "mixed scripts", input: "Hello Мир", wantErr: false},
		{name: "max length", input: strings.Repeat("a", domain.MaxDescriptionLength), wantErr: false},
		{name: "max length cyrillic", input: strings.Repeat("а", domain.MaxDescriptionLength), wantErr: false},

		// length
		{name: "over max length", input: strings.Repeat("a", domain.MaxDescriptionLength+1), wantErr: true},
		{name: "over max length cyrillic", input: strings.Repeat("а", domain.MaxDescriptionLength+1), wantErr: true},

		// whitespace
		{name: "leading space", input: " Hello", wantErr: true},
		{name: "trailing space", input: "Hello ", wantErr: true},
		{name: "tab character", input: "Hello\tWorld", wantErr: true},

		// invalid characters
		{name: "emoji", input: "Great description 🚀", wantErr: true},
		{name: "chinese", input: "描述", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := domain.ValidateDescription(test.input)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateDescription(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			}
			if err != nil && test.wantErr {
				if err != domain.ErrInvalidDescription {
					t.Errorf("ValidateDescription(%q) = %v, want ErrInvalidDescription", test.input, err)
				}
			}
		})
	}
}
