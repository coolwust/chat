package main

import (
	"testing"
	"os"
	"sync"
	r "github.com/dancannon/gorethink"
)

var rTestDatabase = os.Getenv("CHAT_R_TEST_DATABASE")
var rOnce sync.Once

func rSetUp() {
	rOnce.Do(func() {
		var err error
		rSession, err = r.Connect(r.ConnectOpts{
			Address:  rAddress,
			Database: rTestDatabase,
		})
		if err != nil {
			panic(err)
		}
		var name string
		cur, _ := r.TableList().Run(rSession)
		var exist bool
		for cur.Next(&name) {
			if name == rUserTable {
				exist = true
			}
		}
		if !exist {
			r.TableCreate(rUserTable).RunWrite(rSession)
		}
		_, err = r.Table(rUserTable).IndexCreate("email").RunWrite(rSession)
	})
	r.Table(rUserTable).Delete().Exec(rSession)
}

var insertUserTests = []struct{
	email  string
	passwd string
	name   string
	err    error
}{
	{"foo@example.com", "hello", "foo", nil},
	{"bar@example.com", "world", "bar", nil},
	{"bar@example.com", "world", "bar", ErrEmailInUse},
}

func TestInsertUser(t *testing.T) {
	rSetUp()
	for i, tt := range insertUserTests {
		if err := insertUser(tt.email, tt.passwd, tt.name); err != tt.err {
			t.Errorf("%d: error = %s, want %s", i, err, tt.err)
			continue
		}
		if tt.err != nil {
			continue
		}
		u, err := getUser("email", tt.email)
		if err != nil {
			t.Error("%d: insert error = %s", i, err)
			continue
		}
		if tt.email != u.Email {
			t.Error("%d: email = %s, want %s", i, u.Email, tt.email)
		}
		if tt.name != u.Profile.Name {
			t.Error("%d: name = %s, want %s", i, u.Profile.Name, tt.name)
		}
	}
}

var deleteUserTest = struct {
	email  string
	passwd string
	name   string
}{
	email:  "foo@example.com",
	passwd: "hello",
	name:   "foo",
}

func TestDeleteUser(t *testing.T) {
	rSetUp()
	tt := deleteUserTest
	if err := insertUser(tt.email, tt.passwd, tt.name); err != nil {
		t.Fatal(err)
	}
	u, err := getUser("email", tt.email)
	if err != nil {
		t.Fatal(err)
	}
	if err := deleteUser(u.ID); err != nil {
		t.Fatal(err)
	}
}

var updateEmailTest = struct {
	email1 string
	email2 string
	passwd string
	name   string
}{
	email1: "foo@example.com",
	email2: "bar@example.com",
	passwd: "hello",
	name:   "foo",
}

func TestUpdateEmail(t *testing.T) {
	rSetUp()
	tt := updateEmailTest
	if err := insertUser(tt.email1, tt.passwd, tt.name); err != nil {
		t.Fatal(err)
	}
	u, err := getUser("email", tt.email1)
	if err != nil {
		t.Fatal(err)
	}
	if err := updateEmail(u.ID, tt.email2); err != nil {
		t.Fatal(err)
	}
	if _, err := getUser("email", tt.email2); err != nil {
		t.Fatal(err)
	}
}

var authenticateUserTest = struct {
	email  string
	passwd string
	name   string
}{
	email:  "foo@example.com",
	passwd: "hello",
	name:   "foo",
}

func TestAuthenticateUser(t *testing.T) {
	rSetUp()
	tt := authenticateUserTest
	if err := insertUser(tt.email, tt.passwd, tt.name); err != nil {
		t.Fatal(err)
	}
	if _, _, ok := authenticateUser(tt.email, tt.passwd); !ok {
		t.Fatal("login should be successful with password %s", tt.passwd)
	}
	wrong := tt.passwd + " "
	if _, _, ok := authenticateUser(tt.email, wrong); ok {
		t.Fatal("login should be failed with password %s", wrong)
	}
}

var updatePasswordTest = struct {
	email   string
	passwd1 string
	passwd2 string
	name    string
}{
	email:   "foo@example.com",
	passwd1: "hello",
	passwd2: "world",
	name:    "foo",
}

func TestUpdatePassword(t *testing.T) {
	rSetUp()
	tt := updatePasswordTest
	if err := insertUser(tt.email, tt.passwd1, tt.name); err != nil {
		t.Fatal(err)
	}
	u, err := getUser("email", tt.email)
	if err != nil {
		t.Fatal(err)
	}
	if err := updatePassword(u.ID, tt.passwd2); err != nil {
		t.Fatal(err)
	}
	if _, _, ok := authenticateUser(tt.email, tt.passwd2); !ok {
		t.Fatal("login should be successful with password %s", tt.passwd2)
	}
}

var updateProfileTest = struct {
	email  string
	passwd string
	name   string
	profile *Profile
}{
	email:   "foo@example.com",
	passwd:  "hello",
	name:    "foo",
	profile: &Profile{Name: "bar"},
}

func TestUpdateProfile(t *testing.T) {
	rSetUp()
	tt := updateProfileTest
	if err := insertUser(tt.email, tt.passwd, tt.name); err != nil {
		t.Fatal(err)
	}
	u, err := getUser("email", tt.email)
	if err != nil {
		t.Fatal(err)
	}
	if err := updateProfile(u.ID, tt.profile); err != nil {
		t.Fatal(err)
	}
	u, err = getUser("email", tt.email)
	if err != nil {
		t.Fatal(err)
	}
	if u.Profile.Name != tt.profile.Name {
		t.Fatalf("name = %s, want %s", u.Profile.Name, tt.profile.Name)
	}
}
