package main

import (
	"testing"
)

var emailFieldTests = []struct {
	email string
	ok    bool
}{
	{"foo@example.com", true},
	{"foo,@example.com", false},
	{"foo@,example.com", false},
}

func TestEmailField(t *testing.T) {
	for i, test := range emailFieldTests {
		_, ok := emailField.Validate(test.email); if ok != test.ok {
			t.Errorf("%d: expect %s to be %v, got %v", i, test.email, test.ok, ok)
		}
	}
}

var passwordFieldTests = []struct {
	password string
	ok       bool
}{
	{"hello world", true},
	{"hello", false},
	{"world", false},
}

func TestPasswordField(t *testing.T) {
	for i, test := range passwordFieldTests {
		_, ok := passwordField.Validate(test.password); if ok != test.ok {
			t.Errorf("%d: expect %s to be %v, got %v", i, test.password, test.ok, ok)
		}
	}
}

var nicknameFieldTests = []struct {
	nickname string
	ok       bool
}{
	{"foo", true},
	{"x", false},
}

func TestNicknameField(t *testing.T) {
	for i, test := range nicknameFieldTests {
		_, ok := nicknameField.Validate(test.nickname); if ok != test.ok {
			t.Errorf("%d: expect %s to be %v, got %v", i, test.nickname, test.ok, ok)
		}
	}
}
