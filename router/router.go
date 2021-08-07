package router

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	Addr = "127.0.0.1:9900"
)

type API interface {
	// v1 api
	Users(http.ResponseWriter, *http.Request)
	FullnameQuery(http.ResponseWriter, *http.Request)
	UserInfo(http.ResponseWriter, *http.Request)
	SignUp(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

func Route(api API) *http.Server {
	root := mux.NewRouter()
	ui := root.PathPrefix("/ui").Subrouter()

	ui.HandleFunc("", func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write([]byte("This is UI project.")); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}).Methods(http.MethodGet, http.MethodPost)

	v1 := ui.PathPrefix("/v1").Subrouter()

	v1.HandleFunc("/users", api.Users).Methods(http.MethodGet)
	v1.HandleFunc("/user", api.FullnameQuery).Queries("fullname", "{fullname}")
	v1.HandleFunc("/user/{acct:[A-Za-z0-9_]{8,20}}", api.UserInfo).Methods(http.MethodGet)
	v1.HandleFunc("/signup", api.SignUp).Methods(http.MethodPost)
	v1.HandleFunc("/user/{acct:[A-Za-z0-9_]{8,20}}", api.Delete).Methods(http.MethodDelete)
	//r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Addr:         Addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      root,
	}

	return srv
}
