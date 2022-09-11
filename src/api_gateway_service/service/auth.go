package service

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
	"net/http"
	"time"
)

var jwtKey = []byte("secret")

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

	creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to complete auth for %s: %s",
			provider, err), http.StatusInternalServerError)
		return
	}

	user, err := provider.GetUser(creds)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to get user from %s: %s",
			provider, err), http.StatusInternalServerError)
		return
	}

	authCookieValue := objx.New(map[string]interface{}{
		"name":  user.Name(),
		"email": user.Email(),
	}).MustBase64()

	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: authCookieValue,
		Path:  "/",
	})

	claim := jwt.StandardClaims{
		Id:        user.Email(),
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error when trying to make jwt token"), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: tokenString,
		Path:  "/",
	})

	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func logoutHandler(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, fmt.Sprintln(err), http.StatusBadRequest)
			return
		}
		tokenString := token.Value

		_, err = jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			w.Header().Set("Location", "/logout")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)

	})
}
