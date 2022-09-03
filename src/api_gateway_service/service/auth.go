package service

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"net/http"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.validateServeHTTP(w, r); err != nil {
		return
	}
	h.next.ServeHTTP(w, r)
}

func (h *authHandler) validateServeHTTP(w http.ResponseWriter, r *http.Request) (err error) {
	_, err = r.Cookie("auth")
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := gomniauth.Provider("google")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",
			provider, err), http.StatusBadRequest)
		return
	}
	loginUrl, err := provider.GetBeginAuthURL(nil, nil)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s",
			provider, err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Location", loginUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := gomniauth.Provider("google")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get provider %s: %s",
			provider, err), http.StatusBadRequest)
		return
	}

	_, err = provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s: %s",
			provider, err), http.StatusInternalServerError)
		return
	}
}
