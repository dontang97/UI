package ui

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/dontang97/ui/pg"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

var validAcctPwd = regexp.MustCompile(`[A-Za-z0-9_]{8,20}`)

type QueryUserHandlerFunc func(*UI, ...interface{}) ([]pg.User, error)
type AddUserHandlerFunc func(*UI, *pg.User) error
type DeleteUserHandlerFunc func(*UI, *pg.User) error

func scanUsers(ui *UI, rows *sql.Rows) ([]pg.User, error) {
	var users []pg.User
	for rows.Next() {
		var user pg.User
		if err := ui.DB().ScanRows(rows, &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

//////////////////////////////////
/////    GET /ui/v1/users    /////
//////////////////////////////////

var UsersHdl QueryUserHandlerFunc = func(ui *UI, _ ...interface{}) ([]pg.User, error) {
	rows, err := ui.DB().
		Table(pg.TableUsers.String()).
		Select(pg.FieldUserAcct.String()).
		Rows()
	if err != nil {
		return nil, err
	}

	return scanUsers(ui, rows)
}

func (ui *UI) Users(w http.ResponseWriter, _ *http.Request) {
	users, err := UsersHdl(ui)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accts := []string{}
	for _, user := range users {
		accts = append(accts, user.Acct)
	}

	data := map[string][]string{"users": accts}
	WriteJsonResponse(StatusOK, data, w)
}

///////////////////////////////////////////////////////
//////    GET /ui/v1/user?fullname={fullname}    //////
///////////////////////////////////////////////////////

var FullnameQueryHdl QueryUserHandlerFunc = func(ui *UI, args ...interface{}) ([]pg.User, error) {
	rows, err := ui.DB().
		Table(pg.TableUsers.String()).
		Select(pg.FieldUserAcct.String()).
		Where(pg.FieldUserFullname.String()+" = ?", args[0]).Rows()
	if err != nil {
		return nil, err
	}

	return scanUsers(ui, rows)
}

func (ui *UI) FullnameQuery(w http.ResponseWriter, r *http.Request) {
	fullname := r.URL.Query().Get(pg.FieldUserFullname.String())
	users, err := FullnameQueryHdl(ui, fullname)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accts := []string{}
	for _, user := range users {
		accts = append(accts, user.Acct)
	}

	data := map[string]interface{}{
		pg.FieldUserFullname.String(): fullname,
		"users":                       accts,
	}

	WriteJsonResponse(StatusOK, data, w)
}

//////////////////////////////////////////////////////////////
//////    GET /ui/v1/user/{acct:[A-Za-z0-9_]{8,20}}}    //////
//////////////////////////////////////////////////////////////

var UserInfoHdl QueryUserHandlerFunc = func(ui *UI, args ...interface{}) ([]pg.User, error) {
	rows, err := ui.DB().
		Table(pg.TableUsers.String()).
		Select("*").
		Where(pg.FieldUserAcct.String()+" = ?", args[0]).
		Limit(1).
		Rows()
	if err != nil {
		return nil, err
	}

	return scanUsers(ui, rows)
}

func (ui *UI) UserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acct := vars[pg.FieldUserAcct.String()]
	users, err := UserInfoHdl(ui, acct)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) > 1 {
		log.Print(fmt.Errorf("Error: %v records were found for account %v", len(users), acct))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		return
	}

	WriteJsonResponse(StatusOK, users[0], w)
}

//////////////////////////////////////
//////    POST /ui/v1/signup    //////
//////////////////////////////////////

var SignUpHdl AddUserHandlerFunc = func(ui *UI, user *pg.User) error {

	if res := ui.DB().Table(pg.TableUsers.String()).Create(user); res.Error != nil {
		err := res.Error
		return err
	}

	return nil
}

func (ui *UI) SignUp(w http.ResponseWriter, r *http.Request) {
	// TODO: check content-type
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsmap := map[string]interface{}{}
	err = json.Unmarshal(body, &jsmap)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var ok bool
	user := pg.User{}

	if user.Acct, ok = jsmap["account"].(string); !ok {
		WriteJsonResponse(StatusInvalidContent, map[string]string{"missing_field": "account"}, w)
		return
	}

	if user.Pwd, ok = jsmap["password"].(string); !ok {
		WriteJsonResponse(StatusInvalidContent, map[string]string{"missing_field": "password"}, w)
		return
	}

	if user.Fullname, ok = jsmap["fullname"].(string); !ok {
		log.Print("Not provide fullname")
	}

	// check valid acct and pwd
	if !validAcctPwd.MatchString(user.Acct) {
		WriteJsonResponse(StatusInvalidContent,
			map[string]map[string]string{"invalid": {"field": "account", "value": user.Acct}}, w)
		return
	}
	if !validAcctPwd.MatchString(user.Pwd) {
		WriteJsonResponse(StatusInvalidContent,
			map[string]map[string]string{"invalid": {"field": "password", "value": user.Pwd}}, w)
		return
	}

	err = SignUpHdl(ui, &user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == pq.ErrorCode("23505") {
			WriteJsonResponse(StatusUserExisted, map[string]string{"user": user.Acct}, w)
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJsonResponse(StatusOK, map[string]string{"user": user.Acct}, w)
}

//////////////////////////////////////////////////////////////
//////   DELETE /ui/v1/user/{acct:[A-Za-z0-9_]{8,20}}}   /////
//////////////////////////////////////////////////////////////

var DeleteHdl DeleteUserHandlerFunc = func(ui *UI, user *pg.User) error {
	if res := ui.DB().
		Table(pg.TableUsers.String()).
		Delete(&pg.User{}, pg.FieldUserAcct.String()+" = ?", user.Acct); res.Error != nil {
		err := res.Error
		return err
	}
	return nil
}

func (ui *UI) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	acct := vars[pg.FieldUserAcct.String()]
	err := DeleteHdl(ui, &pg.User{Acct: acct})

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJsonResponse(StatusOK, map[string]string{"user": acct}, w)
}
