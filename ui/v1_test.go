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

	UsersHdl         ui.QueryUserHandlerFunc
	FullnameQueryHdl ui.QueryUserHandlerFunc
}

func (s *_Suite) SetupSuite() {
	s.UI = ui.New()
}

func (s *_Suite) TearDownSuite() {
}

func (s *_Suite) SetupTest() {
	s.UsersHdl = ui.UsersHdl
	ui.UsersHdl = nil
	s.FullnameQueryHdl = ui.FullnameQueryHdl
	ui.FullnameQueryHdl = nil
}

func (s *_Suite) TearDownTest() {
	ui.UsersHdl = s.UsersHdl
	s.UsersHdl = nil
	ui.FullnameQueryHdl = s.FullnameQueryHdl
	s.FullnameQueryHdl = nil
}

func (s *_Suite) TestUsers() {
	// normal case
	ui.UsersHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return []pg.User{
			{Acct: "User1"},
			{Acct: "User2"},
		}, nil
	}

	rcd := httptest.NewRecorder()
	http.HandlerFunc(s.UI.Users).ServeHTTP(rcd, nil)
	s.Equal(http.StatusOK, rcd.Code)
	s.Equal("User1\nUser2\n", rcd.Body.String())

	// error case
	ui.UsersHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return nil, errors.New("mock error")
	}
	rcd = httptest.NewRecorder()
	http.HandlerFunc(s.UI.Users).ServeHTTP(rcd, nil)
	s.Equal(http.StatusInternalServerError, rcd.Code)
	s.Equal("", rcd.Body.String())
}

func (s *_Suite) TestFullnameQuery() {
	// normal case
	ui.FullnameQueryHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return []pg.User{
			{Acct: "User1"},
			{Acct: "User2"},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	rcd := httptest.NewRecorder()
	http.HandlerFunc(s.UI.FullnameQuery).ServeHTTP(rcd, req)
	s.Equal(http.StatusOK, rcd.Code)
	s.Equal("User1\nUser2\n", rcd.Body.String())

	// error case
	ui.FullnameQueryHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return nil, errors.New("mock error")
	}
	rcd = httptest.NewRecorder()
	http.HandlerFunc(s.UI.FullnameQuery).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
	s.Equal("", rcd.Body.String())
}

func TestRun(t *testing.T) {
	suite.Run(t, new(_Suite))
}
