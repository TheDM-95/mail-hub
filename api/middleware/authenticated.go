package middleware

import (
	"net/http"
)

const (
	SimpleBasicAuthUser     = "mail-hub"
	SimpleBasicAuthPassword = "mail-hub@123"
)

type AuthenticatedMiddleware struct {
}

func Authenticated() *AuthenticatedMiddleware {
	return &AuthenticatedMiddleware{}
}

func (m *AuthenticatedMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username, password, ok := r.BasicAuth()
	if ok {
		usernameMatch := username == SimpleBasicAuthUser
		passwordMatch := password == SimpleBasicAuthPassword

		if usernameMatch && passwordMatch {
			next(rw, r)
			return
		}
	}

	rw.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	http.Error(rw, "Unauthorized", http.StatusUnauthorized)
	// do some stuff after
}
