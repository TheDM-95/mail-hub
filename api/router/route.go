package router

import (
	"time"

	"github.com/etherlabsio/healthcheck/v2"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"

	"github.com/TheDM-95/mail-hub/api/handler"
	"github.com/TheDM-95/mail-hub/api/middleware"
)

func ResolveRoute(r *mux.Router) {
	r.Handle("/health-check", healthcheck.Handler(
		healthcheck.WithTimeout(5*time.Second),
	))

	subMailRouter := mux.NewRouter().PathPrefix("/api/mail").Subrouter().StrictSlash(true)

	sendHandler := handler.NewSendMailHandler()
	subMailRouter.Methods("POST").PathPrefix("/send").HandlerFunc(sendHandler.Handle)

	r.PathPrefix("/api/mail").Handler(negroni.New(
		middleware.Authenticated(),
		negroni.Wrap(subMailRouter),
	))
}
