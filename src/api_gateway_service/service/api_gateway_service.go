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
	gomniauth.WithProviders(
		google.New(
			"624312092698-n5djffqsu9cocq008fr729osg7657l8m.apps.googleusercontent.com",
			"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",
			"http://localhost:8080/auth/callback/google",
		))

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/login/", loginHandler)
	http.HandleFunc("/auth/callback/", callbackHandler)

	log.Println("Starting api gateway service", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
