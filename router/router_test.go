package router_test

import (
	"context"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/dontang97/ui/pg"
	"github.com/dontang97/ui/router"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
)

type _Suite struct {
	suite.Suite
	srv *http.Server

	// for mocking middle func
	JWTMiddleFunc mux.MiddlewareFunc

	// test fields
	flagLogin         bool
	flagLogout        bool
	flagUsers         bool
	flagFullnameQuery bool

	flagUserInfo    bool
	acctVarUserInfo string

	flagSignup bool
	flagDelete bool
	flagUpdate bool
}

func (s *_Suite) Login(http.ResponseWriter, *http.Request) {
	s.flagLogin = true
}

func (s *_Suite) Logout(http.ResponseWriter, *http.Request) {
	s.flagLogout = true
}

func (s *_Suite) Users(http.ResponseWriter, *http.Request) {
	s.flagUsers = true
}

func (s *_Suite) FullnameQuery(http.ResponseWriter, *http.Request) {
	s.flagFullnameQuery = true
}

func (s *_Suite) UserInfo(_ http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s.acctVarUserInfo = vars[pg.FieldUserAcct.String()]
	s.flagUserInfo = true
}

func (s *_Suite) SignUp(http.ResponseWriter, *http.Request) {
	s.flagSignup = true
}

func (s *_Suite) Delete(http.ResponseWriter, *http.Request) {
	s.flagDelete = true
}

func (s *_Suite) Update(http.ResponseWriter, *http.Request) {
	s.flagUpdate = true
}

func (s *_Suite) SetupSuite() {
	s.JWTMiddleFunc, router.JWTMiddleFunc = router.JWTMiddleFunc, func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	s.srv = router.Route(s)
	go func() {
		s.Equal(http.ErrServerClosed, s.srv.ListenAndServe())
	}()

	// wait the mock server
	time.Sleep(1 * time.Second)
}

func (s *_Suite) TearDownSuite() {
	s.srv.Shutdown(context.Background())
	router.JWTMiddleFunc, s.JWTMiddleFunc = s.JWTMiddleFunc, nil
}

func (s *_Suite) SetupTest() {
	s.flagLogin = false
	s.flagLogout = false
	s.flagUsers = false
	s.flagFullnameQuery = false

	s.flagUserInfo = false
	s.acctVarUserInfo = ""

	s.flagSignup = false
	s.flagDelete = false
	s.flagUpdate = false
}

func (s *_Suite) TearDownTest() {
}

func (s *_Suite) TestRoute() {
	// Get /ui
	resp, err := http.Get("http://" + router.Addr + "/ui")
	s.Equal(nil, err)
	s.Equal(http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	s.Equal(nil, err)
	s.Equal("This is UI project.", string(body))

	// Post /ui
	resp, err = http.Post("http://"+router.Addr+"/ui", "", nil)
	s.Equal(nil, err)
	s.Equal(http.StatusOK, resp.StatusCode)

	body, err = io.ReadAll(resp.Body)
	s.Equal(nil, err)
	s.Equal("This is UI project.", string(body))

	// Get /ui/v1/users
	_, err = http.Get("http://" + router.Addr + "/ui/v1/users")
	s.Equal(nil, err)
	s.Equal(true, s.flagUsers)

	// Get /ui/v1/user?fullname={fullname}
	_, err = http.Get("http://" + router.Addr + "/ui/v1/user?fullname=test")
	s.Equal(nil, err)
	s.Equal(true, s.flagFullnameQuery)

	// Get /ui/v1/user/{acct:[A-Za-z0-9_]{8,20}}
	_, err = http.Get("http://" + router.Addr + "/ui/v1/user/user_acct")
	s.Equal(nil, err)
	s.Equal("user_acct", s.acctVarUserInfo)
	s.Equal(true, s.flagUserInfo)

	// Post /ui/v1/signup
	_, err = http.Post("http://"+router.Addr+"/ui/v1/signup", "", nil)
	s.Equal(nil, err)
	s.Equal(true, s.flagSignup)

	// Delete /ui/v1/user/{acct:[A-Za-z0-9_]{8,20}}
	req, err := http.NewRequest(http.MethodDelete, "http://"+router.Addr+"/ui/v1/user/user_acct", nil)
	s.Equal(nil, err)
	c := http.Client{}
	_, err = c.Do(req)
	s.Equal(nil, err)
	s.Equal(true, s.flagDelete)

	// Put /ui/v1/user/{acct:[A-Za-z0-9_]{8,20}}
	req, err = http.NewRequest(http.MethodPut, "http://"+router.Addr+"/ui/v1/user/user_acct", nil)
	s.Equal(nil, err)
	c = http.Client{}
	_, err = c.Do(req)
	s.Equal(nil, err)
	s.Equal(true, s.flagUpdate)

	// Post /ui/v1/login
	_, err = http.Post("http://"+router.Addr+"/ui/v1/login", "", nil)
	s.Equal(nil, err)
	s.Equal(true, s.flagLogin)
}

func TestRun(t *testing.T) {
	suite.Run(t, new(_Suite))
}
