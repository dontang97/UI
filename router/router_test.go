package router_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/dontang97/ui/router"
	"github.com/stretchr/testify/suite"
)

type _Suite struct {
	suite.Suite
	srv *http.Server

	// test fields
	flagLogin  bool
	flagLogout bool
	flagUsers  bool
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

func (s *_Suite) SetupSuite() {

	s.srv = router.Route(s)
	go func() {
		s.Equal(http.ErrServerClosed, s.srv.ListenAndServe())
	}()
}

func (s *_Suite) TearDownSuite() {
	s.srv.Shutdown(context.Background())
}

func (s *_Suite) SetupTest() {
	s.flagLogin = false
	s.flagLogout = false
	s.flagUsers = false
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
}

func TestRun(t *testing.T) {
	suite.Run(t, new(_Suite))
}
