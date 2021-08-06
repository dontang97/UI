package ui_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dontang97/ui/pg"
	"github.com/dontang97/ui/ui"
	"github.com/stretchr/testify/suite"
)

type _Suite struct {
	suite.Suite
	UI *ui.UI

	// test fields
	QueryResp struct {
		Users []pg.User
		err   error
	}
}

func (s *_Suite) Query() ([]pg.User, error) {
	return s.QueryResp.Users, s.QueryResp.err
}

func (s *_Suite) SetupSuite() {
	s.UI = ui.New(s)
}

func (s *_Suite) TearDownSuite() {
}

func (s *_Suite) SetupTest() {
	s.QueryResp.Users = nil
	s.QueryResp.err = nil
}

func (s *_Suite) TearDownTest() {
}

func (s *_Suite) TestQuery() {
	// normal case
	s.QueryResp.Users = []pg.User{
		{Acct: "User1"},
		{Acct: "User2"},
	}
	s.QueryResp.err = nil

	rcd := httptest.NewRecorder()
	http.HandlerFunc(s.UI.Users).ServeHTTP(rcd, nil)
	s.Equal(http.StatusOK, rcd.Code)
	s.Equal("User1\nUser2\n", rcd.Body.String())

	// error case
	s.QueryResp.err = errors.New("mock err")
	rcd = httptest.NewRecorder()
	http.HandlerFunc(s.UI.Users).ServeHTTP(rcd, nil)
	s.Equal(http.StatusInternalServerError, rcd.Code)
	s.Equal("", rcd.Body.String())
}

func TestRun(t *testing.T) {
	suite.Run(t, new(_Suite))
}
