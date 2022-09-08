package service

import (
	"flag"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func Run() {
	var addr = flag.String("addr", ":5000", "The addr of the application.")
	flag.Parse()

	gomniauth.SetSecurityKey("SECURITY KEY")
	gomniauth.WithProviders(google.New(
		"key",
		"Secret",
		"http://localhost:8080/auth/callback/google",
	))

	lmURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	smURL, err := url.Parse("http://localhost:9090")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/links", httputil.NewSingleHostReverseProxy(lmURL))
	http.Handle("/followers", httputil.NewSingleHostReverseProxy(smURL))
	http.Handle("/folliwing", httputil.NewSingleHostReverseProxy(smURL))
	http.Handle("/follow", httputil.NewSingleHostReverseProxy(smURL))
	http.Handle("/unfollow", httputil.NewSingleHostReverseProxy(smURL))

	log.Println("Starting api gateway service", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
