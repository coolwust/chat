package main

import (
	"errors"
	"strings"
	"golang.org/x/crypto/bcrypt"
	r "github.com/dancannon/gorethink"
)

const bcryptConst int = 12

var (
	ErrUserNotFound = errors.New("The user is not found in database")
	ErrEmailInUse = errors.New("The email is already in use")
)

type User struct {
	ID       string   `gorethink:"id"`
	Email    string   `gorethink:"email"`
	Password string   `gorethink:"password"`
	Profile  *Profile `gorethink:"profile"`
}

type Profile struct {
	Name    string `gorethink:"name"`
	URL     string `gorethink:"url"`
	Company string `gorethink:"company"`
}

func getUser(index, value string) (*User, error) {
	cur, err := r.Table(rUserTable).GetAllByIndex(index, value).Run(rSession)
	if err != nil {
		return nil, err
	}
	defer cur.Close()
	if cur.IsNil() {
		return nil, ErrUserNotFound
	}
	user := &User{}
	if err := cur.One(user); err != nil {
		return nil, err
	}
	return user, nil
}

func insertUser(email, passwd, name string) error {
	passwd, err := encryptPassword(passwd)
	if err != nil {
		return err
	}
	id, err := uuid()
	if err != nil {
		return err
	}
	user := &User{
		ID:       id,
		Email:    email,
		Password: passwd,
		Profile:  &Profile{Name: name},
	}
	if err := r.Branch(
		r.Table(rUserTable).GetAllByIndex("email", email).IsEmpty(),
		r.Table(rUserTable).Insert(user),
		r.Error("%%% email in use %%%"),
	).Exec(rSession); err != nil {
		if strings.Contains(err.Error(), "%%% email in use %%%") {
			return ErrEmailInUse
		}
		return err
	}
	return nil
}

func deleteUser(id string) error {
	return r.Table(rUserTable).Get(id).Delete().Exec(rSession)
}

func updateEmail(id, email string) error {
	update := map[string]interface{}{"email": email}
	return r.Table(rUserTable).Get(id).Update(update).Exec(rSession)
}

func updatePassword(id, passwd string) error {
	passwd, err := encryptPassword(passwd)
	if err != nil {
		return err
	}
	update := map[string]interface{}{"password": passwd}
	return r.Table(rUserTable).Get(id).Update(update).Exec(rSession)
}

func updateProfile(id string, p *Profile) error {
	update := map[string]interface{}{"profile": p}
	return r.Table(rUserTable).Get(id).Update(update).Exec(rSession)
}

func authenticateUser(email, passwd string) (*User, string, bool) {
	user, err := getUser("email", email)
	if err != nil {
		if err == ErrUserNotFound {
			goto INCORRECT
		}
		return nil, "Internel server error. Please try again later.", false
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd)) == nil {
		return user, "", true
	}
INCORRECT:
	return nil, "The email or password you entered is incorrect.", false
}

func encryptPassword(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcryptConst)
	return string(hash), err
}
