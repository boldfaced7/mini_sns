package service

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
)

func Run() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()

	gomniauth.SetSecurityKey("SECURITY KEY")
	gomniauth.WithProviders(google.New(
		"key",
		"Secret",
		"http://localhost:8080/auth/callback/google",
	))

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/auth/", loginHandler)

	log.Println("Starting api gateway service", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
