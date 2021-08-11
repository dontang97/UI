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
	"github.com/dontang97/ui/secret"
	"github.com/dontang97/ui/ui"
	"github.com/stretchr/testify/suite"
)

type _v1Suite struct {
	suite.Suite
	UI *ui.UI

	UsersHdl         ui.QueryUserHandlerFunc
	FullnameQueryHdl ui.QueryUserHandlerFunc
	UserInfoHdl      ui.QueryUserHandlerFunc
	SignUpHdl        ui.AddUserHandlerFunc
	DeleteHdl        ui.DeleteUserHandlerFunc
	UpdateHdl        ui.UpdateUserHandlerFunc
	LoginHdl         ui.QueryUserHandlerFunc
}

func (s *_v1Suite) SetupSuite() {
	secret.InitSecretKey("../secret")
	s.UI = ui.New()
}

func (s *_v1Suite) TearDownSuite() {
}

func (s *_v1Suite) SetupTest() {
	s.UsersHdl, ui.UsersHdl = ui.UsersHdl, nil
	s.FullnameQueryHdl, ui.FullnameQueryHdl = ui.FullnameQueryHdl, nil
	s.UserInfoHdl, ui.UserInfoHdl = ui.UserInfoHdl, nil
	s.SignUpHdl, ui.SignUpHdl = ui.SignUpHdl, nil
	s.DeleteHdl, ui.DeleteHdl = ui.DeleteHdl, nil
	s.UpdateHdl, ui.UpdateHdl = ui.UpdateHdl, nil
	s.LoginHdl, ui.LoginHdl = ui.LoginHdl, nil
}

func (s *_v1Suite) TearDownTest() {
	ui.UsersHdl, s.UsersHdl = s.UsersHdl, nil
	ui.FullnameQueryHdl, s.FullnameQueryHdl = s.FullnameQueryHdl, nil
	ui.UserInfoHdl, s.UserInfoHdl = s.UserInfoHdl, nil
	ui.SignUpHdl, s.SignUpHdl = s.SignUpHdl, nil
	ui.DeleteHdl, s.DeleteHdl = s.DeleteHdl, nil
	ui.UpdateHdl, s.UpdateHdl = s.UpdateHdl, nil
	ui.LoginHdl, s.LoginHdl = s.LoginHdl, nil
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
	ui.UserInfoHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
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
	ui.SignUpHdl = func(ui *ui.UI, user *pg.User) error {
		return nil
	}

	user := struct {
		Acct     string `json:"account"`
		Pwd      string `json:"password"`
		Fullname string `json:"fullname"`
	}{
		Acct:     "123456789",
		Pwd:      "123456789",
		Fullname: "123456789",
	}

	js, err := json.Marshal(user)
	s.Equal(nil, err)

	req := httptest.NewRequest(http.MethodPost, "http://test.com", bytes.NewBuffer(js))
	rcd := httptest.NewRecorder()

	http.HandlerFunc(s.UI.SignUp).ServeHTTP(rcd, req)
	s.Equal(http.StatusOK, rcd.Code)

	// error case
	ui.SignUpHdl = func(ui *ui.UI, user *pg.User) error {
		return errors.New("mock error")
	}

	req = httptest.NewRequest(http.MethodPost, "http://test.com", bytes.NewBuffer(js))
	rcd = httptest.NewRecorder()

	http.HandlerFunc(s.UI.SignUp).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
}

func (s *_v1Suite) TestDelete() {
	// normal case
	ui.DeleteHdl = func(ui *ui.UI, user *pg.User) error {
		return nil
	}
	req := httptest.NewRequest(http.MethodDelete, "http://test.com", nil)
	rcd := httptest.NewRecorder()

	http.HandlerFunc(s.UI.Delete).ServeHTTP(rcd, req)
	s.Equal(http.StatusOK, rcd.Code)

	// error case
	ui.DeleteHdl = func(ui *ui.UI, user *pg.User) error {
		return errors.New("mock error")
	}
	req = httptest.NewRequest(http.MethodDelete, "http://test.com", nil)
	rcd = httptest.NewRecorder()

	http.HandlerFunc(s.UI.Delete).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
}

func (s *_v1Suite) TestUpdate() {
	// normal case
	ui.UpdateHdl = func(ui *ui.UI, user *pg.User) error {
		return nil
	}
	user := struct {
		Pwd      string `json:"password"`
		Fullname string `json:"fullname"`
	}{
		Pwd:      "123456789",
		Fullname: "123456789",
	}

	js, err := json.Marshal(user)
	s.Equal(nil, err)

	req := httptest.NewRequest(http.MethodPut, "http://test.com/", bytes.NewBuffer(js))
	rcd := httptest.NewRecorder()

	http.HandlerFunc(s.UI.Update).ServeHTTP(rcd, req)
	s.Equal(http.StatusOK, rcd.Code)

	// error case
	ui.UpdateHdl = func(ui *ui.UI, user *pg.User) error {
		return errors.New("mock error")
	}
	req = httptest.NewRequest(http.MethodPut, "http://test.com/", bytes.NewBuffer(js))
	rcd = httptest.NewRecorder()

	http.HandlerFunc(s.UI.Update).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
}

func (s *_v1Suite) TestLogin() {
	// normal case
	ui.LoginHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return []pg.User{{Pwd: "123456789"}}, nil
	}
	user := struct {
		Acct string `json:"account"`
		Pwd  string `json:"password"`
	}{
		Acct: "123456789",
		Pwd:  "123456789",
	}

	js, err := json.Marshal(user)
	s.Equal(nil, err)

	req := httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewBuffer(js))
	rcd := httptest.NewRecorder()

	http.HandlerFunc(s.UI.Login).ServeHTTP(rcd, req)
	//str := rcd.Body.String()
	//fmt.Println(str)
	s.Equal(http.StatusOK, rcd.Code)

	// error case
	ui.LoginHdl = func(ui *ui.UI, args ...interface{}) ([]pg.User, error) {
		return nil, errors.New("mock error")
	}

	req = httptest.NewRequest(http.MethodPost, "http://test.com/", bytes.NewBuffer(js))
	rcd = httptest.NewRecorder()

	http.HandlerFunc(s.UI.Login).ServeHTTP(rcd, req)
	s.Equal(http.StatusInternalServerError, rcd.Code)
}

func TestRunV1(t *testing.T) {
	suite.Run(t, new(_v1Suite))
}
