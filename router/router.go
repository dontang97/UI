package router

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dontang97/ui/pg"
	"github.com/dontang97/ui/secret"
	"github.com/dontang97/ui/ui"
	"github.com/gorilla/mux"
)

const (
	Addr = ":9900"
)

type API interface {
	// v1 api
	Users(http.ResponseWriter, *http.Request)
	FullnameQuery(http.ResponseWriter, *http.Request)
	UserInfo(http.ResponseWriter, *http.Request)
	SignUp(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
}

var JWTMiddleFunc mux.MiddlewareFunc = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		acct, _ := vars[pg.FieldUserAcct.String()]

		auth, ok := r.Header["Authorization"]
		if !ok || len(auth) <= 0 || !strings.HasPrefix(auth[0], "Bearer ") {
			ui.WriteJsonResponse(ui.StatusNoAuth,
				map[string]string{"account": acct}, w)
			return
		}

		token := auth[0][len("Bearer "):]
		if err := secret.VerifyUserJWT(token, acct); err != nil {
			if je, ok := err.(*secret.JWTError); ok {
				switch je.Code() {
				case secret.JWTUnknownError,
					secret.JWTNotActiveError,
					secret.JWTExpiredError,
					secret.JWTAcctNotMatchError,
					secret.JWTNotAuthError:
					ui.WriteJsonResponse(ui.StatusNoAuth, map[string]string{"error": je.Error()}, w)
					return
				}
			}

			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
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

	v1.HandleFunc("/signup", api.SignUp).Methods(http.MethodPost)
	v1.HandleFunc("/login", api.Login).Methods(http.MethodPost)

	users := v1.PathPrefix("/users").Subrouter()
	users.Use(JWTMiddleFunc)
	users.HandleFunc("", api.Users).Methods(http.MethodGet)

	user := v1.PathPrefix("/user").Subrouter()
	user.Use(JWTMiddleFunc)
	user.HandleFunc("", api.FullnameQuery).Queries("fullname", "{fullname}")

	acct := user.PathPrefix("/{acct:[A-Za-z0-9_]{8,20}}").Subrouter()
	acct.HandleFunc("", api.UserInfo).Methods(http.MethodGet)
	acct.HandleFunc("", api.Delete).Methods(http.MethodDelete)
	acct.HandleFunc("", api.Update).Methods(http.MethodPut)

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
