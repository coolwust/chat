package main

import (
	"os"
	"net"
	"net/http"
	"flag"
	"crypto/tls"
	r "github.com/dancannon/gorethink"
	"github.com/coldume/mux"
)

var (
	// Environments
	rAddress      = os.Getenv("CHAT_R_ADDRESS")
	rDatabase     = os.Getenv("CHAT_R_DATABASE")
	rSessionTable = os.Getenv("CHAT_R_SESSION_TABLE")
	rUserTable    = os.Getenv("CHAT_R_USER_TABLE")
	sessionKey    = os.Getenv("CHAT_SESSION_KEY")

	// Flags
	httpPort  string
	httpsPort string

	// Variables
	rSession *r.Session
)

func init() {
	flag.Parse()
	flag.StringVar(&httpPort, "http-port", "8080", "HTTP port")
	flag.StringVar(&httpsPort, "https-port", "8443", "HTTPS port")

	var err error
	if rSession, err = r.Connect(r.ConnectOpts{
		Address:  rAddress,
		Database: rDatabase,
	}); err != nil {
		panic(err)
	}
}

func main() {
	go func() {
		listener, err := net.Listen("tcp", ":" + httpPort)
		if err != nil {
			panic(err)
		}
		serv := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://" + r.Host + r.RequestURI, http.StatusMovedPermanently)
		})}
		if err := serv.Serve(listener); err != nil {
			panic(err)
		}
	}()
	cert, err := tls.LoadX509KeyPair("tls/cert.pem", "tls/key.pem")
	if err != nil {
		panic(err)
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: false,
	}
	listener, err := tls.Listen("tcp", ":" + httpsPort, &config)
	if err != nil {
		panic(err)
	}
	m := mux.NewServeMux()
	//m.HandleFunc(`^/$`, indexHandler)
	m.HandleFunc(`^/registration$`, registrationHandler)
	m.Handle("^/public/", http.FileServer(http.Dir(".")))
	serv := &http.Server{Handler: m}
	if err := serv.Serve(listener); err != nil {
		panic(err)
	}
}
