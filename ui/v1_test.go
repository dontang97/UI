package ui_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dontang97/ui/pg"
	"github.com/dontang97/ui/ui"
	"github.com/stretchr/testify/suite"
)

type _v1Suite struct {
	suite.Suite
	UI *ui.UI

	UsersHdl         ui.QueryUserHandlerFunc
	FullnameQueryHdl ui.QueryUserHandlerFunc
	UserInfoQueryHdl ui.QueryUserHandlerFunc
	SignUpAddOneHdl  ui.AddOneUserHandlerFunc
}

func (s *_v1Suite) SetupSuite() {
	s.UI = ui.New()
}

func (s *_v1Suite) TearDownSuite() {
}

func (s *_v1Suite) SetupTest() {
	s.UsersHdl = ui.UsersHdl
	ui.UsersHdl = nil
	s.FullnameQueryHdl = ui.FullnameQueryHdl
	ui.FullnameQueryHdl = nil
	s.UserInfoQueryHdl = ui.UserInfoQueryHdl
	ui.UserInfoQueryHdl = nil
	s.SignUpAddOneHdl = ui.SignUpAddOneHdl
	ui.SignUpAddOneHdl = nil
}

func (s *_v1Suite) TearDownTest() {
	ui.UsersHdl = s.UsersHdl
	s.UsersHdl = nil
	ui.FullnameQueryHdl = s.FullnameQueryHdl
	s.FullnameQueryHdl = nil
	ui.UserInfoQueryHdl = s.UserInfoQueryHdl
	s.UserInfoQueryHdl = nil
	ui.SignUpAddOneHdl = s.SignUpAddOneHdl
	s.SignUpAddOneHdl = nil
}

func (s *_v1Suite) TestUsers() {
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

	body := map[string]interface{}{}
	err := json.Unmarshal(rcd.Body.Bytes(), &body)
	s.Equal(nil, err)
	v := body["data"].(map[string]interface{})
	s.Equal([]interface{}([]interface{}{"User1", "User2"}), v["users"])

	// error case
	ui.UsersHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return nil, errors.New("mock error")
	}
	rcd = httptest.NewRecorder()
	http.HandlerFunc(s.UI.Users).ServeHTTP(rcd, nil)
	s.Equal(http.StatusInternalServerError, rcd.Code)
	s.Equal("", rcd.Body.String())
}

func (s *_v1Suite) TestFullnameQuery() {
	// normal case
	ui.FullnameQueryHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return []pg.User{
			{Acct: "User1"},
			{Acct: "User2"},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "http://test.com?fullname=ABC", nil)
	rcd := httptest.NewRecorder()
	http.HandlerFunc(s.UI.FullnameQuery).ServeHTTP(rcd, req)
	s.Equal(http.StatusOK, rcd.Code)

	body := map[string]interface{}{}
	err := json.Unmarshal(rcd.Body.Bytes(), &body)
	s.Equal(nil, err)
	v := body["data"].(map[string]interface{})
	s.Equal(interface{}("ABC"), v["fullname"])
	s.Equal([]interface{}([]interface{}{"User1", "User2"}), v["users"])

	// error case
	ui.FullnameQueryHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return nil, errors.New("mock error")
	}
	rcd = httptest.NewRecorder()
	http.HandlerFunc(s.UI.FullnameQuery).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
	s.Equal("", rcd.Body.String())
}

func (s *_v1Suite) TestUserInfo() {
	// normal case
	ui.UserInfoQueryHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return []pg.User{
			{
				Acct:       "User1",
				Pwd:        "Pwd1",
				Fullname:   "Fullname1",
				Created_at: time.Time{},
				Updated_at: time.Time{},
			},
		}, nil
	}

	req := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	rcd := httptest.NewRecorder()
	http.HandlerFunc(s.UI.UserInfo).ServeHTTP(rcd, req)
	s.Equal(http.StatusOK, rcd.Code)

	body := map[string]interface{}{}
	err := json.Unmarshal(rcd.Body.Bytes(), &body)
	s.Equal(nil, err)
	v := body["data"].(map[string]interface{})
	s.Equal(interface{}("User1"), v["account"])
	s.Equal(interface{}("Pwd1"), v["password"])
	s.Equal(interface{}("Fullname1"), v["fullname"])
	//s.Equal(interface{}(time.Time{}.String()), v["created_at"])
	//s.Equal(interface{}(time.Time{}.String()), v["updated_at"])

	// error case
	ui.FullnameQueryHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return nil, errors.New("mock error")
	}
	rcd = httptest.NewRecorder()
	http.HandlerFunc(s.UI.FullnameQuery).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
	s.Equal("", rcd.Body.String())
}

func (s *_v1Suite) TestSignUp() {
	// normal case
	ui.SignUpAddOneHdl = func(ui *ui.UI, user interface{}) error {
		return nil
	}

	user := struct {
		Acct     string `json:"account"`
		Pwd      string `json:"password"`
		Fullname string `json:"fullname"`
	}{
		Acct:     "123456789",
		Pwd:      "123456789",
		Fullname: "",
	}

	js, err := json.Marshal(user)
	s.Equal(nil, err)

	req := httptest.NewRequest(http.MethodPost, "http://test.com", bytes.NewBuffer(js))
	rcd := httptest.NewRecorder()

	http.HandlerFunc(s.UI.SignUp).ServeHTTP(rcd, req)
	s.Equal(http.StatusOK, rcd.Code)

	// error case
	ui.SignUpAddOneHdl = func(ui *ui.UI, user interface{}) error {
		return errors.New("mock error")
	}

	req = httptest.NewRequest(http.MethodPost, "http://test.com", bytes.NewBuffer(js))
	rcd = httptest.NewRecorder()

	http.HandlerFunc(s.UI.SignUp).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
}

func TestRunV1(t *testing.T) {
	suite.Run(t, new(_v1Suite))
}
