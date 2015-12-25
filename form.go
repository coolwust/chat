package main

import (
	"regexp"
)

type field struct {
	hint     string
	validate func(input string) bool
	message  string
}

func (f *field) GetHint() string {
	return f.hint
}

func (f *field) Validate(input string) (string, bool) {
	if f.validate != nil && !f.validate(input) {
		return f.message, false
	}
	return "", true
}

var emailField = &field{
	hint: "Email",
	validate: func(input string) bool {
		ma, _ := regexp.MatchString(`^[[:alnum:]]+@[[:alnum:]]+\.[[:alpha:]]{2,}$`, input)
		return ma
	},
	message: "Email address is invalid",
}

var passwordField = &field{
	hint: "Password (6-100 any characters)",
	validate: func(input string) bool {
		ma, _ := regexp.MatchString(`^.{6,100}$`, input)
		return ma
	},
	message: "Password must be 6-100 characters",
}

var nicknameField = &field{
	hint: "Nickname (3-30 English characters)",
	validate: func(input string) bool {
		ma, _ := regexp.MatchString(`^[[:alpha:]]{3,30}$`, input)
		return ma
	},
	message: "Nickname must be 3-30 English characters",
}
