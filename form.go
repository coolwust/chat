package main

import (
	"regexp"
)

type form interface {
	Validate() bool
}

type control interface {
	Value()    string
	Hint()     string
	Validate() bool
	Error()    string
}

var _ control = &emailControl{}

type emailControl struct {
	value     string
	ok        *bool
	exists    bool
	activated bool
}

func (f *emailControl) Value() string {
	return f.value
}

func (f *emailControl) Hint() string {
	return "example: mike@example.com"
}

func (f *emailControl) Validate() (ok bool) {
	if !f.activated {
		return true
	}
	if f.ok != nil {
		return *f.ok
	}
	f.exists = false
	if ok, _ = regexp.MatchString(`^[[:alnum:]]+@[[:alnum:]]+\.[[:alpha:]]{2,}$`, f.value); ok {
		if _, err := getUser("email", f.value); err != ErrUserNotFound {
			ok = false
			f.exists = true
		}
	}
	f.ok = &ok
	return
}

func (f *emailControl) Error() string {
	if f.Validate() {
		return ""
	}
	if f.exists {
		return "Email is already in use"
	}
	return "Email address is invalid"
}

var _ control = &passwordControl{}

type passwordControl struct {
	value     string
	ok        *bool
	activated bool
}

func (f *passwordControl) Value() string {
	return f.value
}

func (f *passwordControl) Hint() string {
	return "6-100 characters"
}

func (f *passwordControl) Validate() bool {
	if !f.activated {
		return true
	}
	if f.ok != nil {
		return *f.ok
	}
	ma, _ := regexp.MatchString(`^.{6,100}$`, f.value)
	f.ok = &ma
	return ma
}

func (f *passwordControl) Error() string {
	if f.Validate() {
		return ""
	}
	return "Password must be 6-100 characters"
}

var _ control = &nameControl{}

type nameControl struct {
	value     string
	ok        *bool
	activated bool
}

func (f *nameControl) Value() string {
	return f.value
}

func (f *nameControl) Hint() string {
	return "3-30 English characters"
}

func (f *nameControl) Validate() bool {
	if !f.activated {
		return true
	}
	if f.ok != nil {
		return *f.ok
	}
	ma, _ := regexp.MatchString(`^[[:alpha:]]{3,30}$`, f.value)
	f.ok = &ma
	return ma
}

func (f *nameControl) Error() string {
	if f.Validate() {
		return ""
	}
	return "Name must be 3-30 English characters"
}
