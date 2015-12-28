package main

import (
	"testing"
)

var emailControlTests = []struct {
	email string
	ok    bool
}{
	{"foo@example.com", true},
	{"foo,@example.com", false},
	{"foo@,example.com", false},
}

func TestEmailControl(t *testing.T) {
	for i, test := range emailControlTests {
		c := &emailControl{value: test.email, activated: true}
		if ok := c.Validate(); ok != test.ok {
			t.Errorf("%d: ok = %v, want %v", i, ok, test.ok)
		}
	}
}

var passwordControlTests = []struct {
	password string
	ok       bool
}{
	{"hello world", true},
	{"hello", false},
	{"world", false},
}

func TestPasswordControl(t *testing.T) {
	for i, test := range passwordControlTests {
		c := &passwordControl{value: test.password, activated: true}
		if ok := c.Validate(); ok != test.ok {
			t.Errorf("%d: ok = %v, want %v", i, ok, test.ok)
		}
	}
}

var nameControlTests = []struct {
	name string
	ok       bool
}{
	{"foo", true},
	{"x", false},
}

func TestNameControl(t *testing.T) {
	for i, test := range nameControlTests {
		c := &nameControl{value: test.name, activated: true}
		if ok := c.Validate(); ok != test.ok {
			t.Errorf("%d: ok = %v, want %v", i, ok, test.ok)
		}
	}
}
