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
