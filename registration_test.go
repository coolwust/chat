package main

import (
	"testing"
)

var formValidateTests = []struct {
	email    string
	passwd   string
	name string
	ok       bool
}{
	{"foo@example.com", "hello world", "foo", true},
	{"bar@,example.com", "world hello", "bar", false},
	{"baz@example.com", "hello", "baz", false},
}

func TestFormValidate(t *testing.T) {
	rSetUp()
	for i, tt := range formValidateTests {
		f := &registrationForm{
			Email:    &emailControl{value: tt.email},
			Password: &passwordControl{value: tt.passwd},
			Name:     &nameControl{value: tt.name},
		}
		if ok := f.Validate(); ok != tt.ok {
			t.Error("%d: ok = %v, want %v", i, ok, tt.ok)
		}
	}
}

func TestRegistrationHandler(t *testing.T) {
}
