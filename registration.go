package main

import (
	"log"
	"reflect"
	"net/http"
	"html/template"
)

type registrationForm struct {
	Email    *emailControl
	Password *passwordControl
	Name     *nameControl
}

func parseRegistrationForm(r *http.Request) (*registrationForm, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	return &registrationForm{
		Email:    &emailControl{value: r.PostForm.Get("email"), activated: true},
		Password: &passwordControl{value: r.PostForm.Get("password"), activated: true},
		Name:     &nameControl{value: r.PostForm.Get("name"), activated: true},
	}, nil
}

func (f *registrationForm) Validate() bool {
	ok := true
	rv := reflect.ValueOf(f).Elem()
	for i := 0; i < rv.NumField(); i++ {
		c := rv.Field(i).Interface().(control)
		if !c.Validate() {
			ok = false
			break
		}
	}
	return ok
}

type registrationData struct {
	Form *registrationForm
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		data := &registrationData{Form: &registrationForm{
			Email:    &emailControl{},
			Password: &passwordControl{},
			Name:     &nameControl{},
		}}
		template.Must(template.ParseFiles("./view/registration.html")).Execute(w, data)
		return
	}
	f, err := parseRegistrationForm(r)
	if err != nil {
		goto InternalServerError
	}
	if !f.Validate() {
		data := &registrationData{Form: f}
		template.Must(template.ParseFiles("./view/registration.html")).Execute(w, data)
		return
	}
	if err := insertUser(f.Email.Value(), f.Password.Value(), f.Name.Value()); err != nil {
	log.Println(err)
		goto InternalServerError
	}
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	return
InternalServerError:
	code := http.StatusInternalServerError
	http.Error(w, http.StatusText(code), code)
}
