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
	"github.com/dontang97/ui/secret"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

var validAcctPwd = regexp.MustCompile(`[A-Za-z0-9_]{8,20}`)

type QueryUserHandlerFunc func(*UI, ...interface{}) ([]pg.User, error)
type AddUserHandlerFunc func(*UI, *pg.User) error
type DeleteUserHandlerFunc func(*UI, *pg.User) error
type UpdateUserHandlerFunc func(*UI, *pg.User) error

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
		WriteJsonResponse(StatusInvalidContent,
			map[string]string{"error": err.Error()},
			w,
		)
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
		WriteJsonResponse(StatusInvalidContent, map[string]string{"missing_field": "fullname"}, w)
		return
	}

	// check valid acct, pwd and fullname
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
	if user.Fullname == "" || len(user.Fullname) > pg.FieldUserFullnameMaxLen {
		WriteJsonResponse(StatusInvalidContent,
			map[string]map[string]string{"invalid": {"field": "fullname", "value": user.Fullname}}, w)
		return
	}

	err = SignUpHdl(ui, &user)
	if err != nil {
		// user has existed
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

//////////////////////////////////////////////////////////////
//////   UPDATE /ui/v1/user/{acct:[A-Za-z0-9_]{8,20}}}   /////
//////////////////////////////////////////////////////////////

var UpdateHdl UpdateUserHandlerFunc = func(ui *UI, user *pg.User) error {
	values := map[string]interface{}{}
	if user.Pwd != "" {
		values[pg.FieldUserPwd.String()] = user.Pwd
	}
	if user.Fullname != "" {
		values[pg.FieldUserFullname.String()] = user.Fullname
	}
	if len(values) == 0 {
		return nil
	}

	if res := ui.DB().
		Table(pg.TableUsers.String()).
		Where(pg.FieldUserAcct.String()+" = ?", user.Acct).
		Updates(values); res.Error != nil {
		err := res.Error
		return err
	}

	return nil
}

func (ui *UI) Update(w http.ResponseWriter, r *http.Request) {
	// TODO: check content-type
	vars := mux.Vars(r)
	user := &pg.User{
		Acct: vars[pg.FieldUserAcct.String()],
	}

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
		WriteJsonResponse(StatusInvalidContent,
			map[string]string{"error": err.Error()},
			w,
		)
		return
	}

	var ok bool
	if user.Pwd, ok = jsmap["password"].(string); ok {
		if !validAcctPwd.MatchString(user.Pwd) {
			WriteJsonResponse(StatusInvalidContent,
				map[string]map[string]string{"invalid": {"field": "password", "value": user.Pwd}}, w)
			return
		}
	}

	if user.Fullname, ok = jsmap["fullname"].(string); ok {
		if user.Fullname == "" || len(user.Fullname) > pg.FieldUserFullnameMaxLen {
			WriteJsonResponse(StatusInvalidContent,
				map[string]map[string]string{"invalid": {"field": "fullname", "value": user.Fullname}}, w)
			return
		}
	}

	err = UpdateHdl(ui, user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJsonResponse(StatusOK, nil, w)
}

//////////////////////////////////////
//////    POST /ui/v1/login     //////
//////////////////////////////////////

var LoginHdl QueryUserHandlerFunc = func(ui *UI, args ...interface{}) ([]pg.User, error) {
	rows, err := ui.DB().
		Table(pg.TableUsers.String()).
		Select(pg.FieldUserPwd.String()).
		Where(pg.FieldUserAcct.String()+" = ?", args[0]).Rows()
	if err != nil {
		return nil, err
	}

	return scanUsers(ui, rows)
}

func (ui *UI) Login(w http.ResponseWriter, r *http.Request) {
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
		WriteJsonResponse(StatusInvalidContent,
			map[string]string{"error": err.Error()},
			w,
		)
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

	// check valid acct, pwd and fullname
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

	users, err := LoginHdl(ui, user.Acct)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) > 1 {
		log.Print(fmt.Errorf("Error: %v records were found for account %v", len(users), user.Acct))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		WriteJsonResponse(StatusUserNotFound,
			map[string]string{"account": user.Acct}, w)
		return
	}

	if users[0].Pwd != user.Pwd {
		WriteJsonResponse(StatusWrongPassword,
			map[string]string{"account": user.Acct}, w)
		return
	}

	// JWT token return
	token, err := secret.CreateUserJWT(user.Acct)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	WriteJsonResponse(StatusOK, map[string]string{"user": user.Acct, "JWT": token}, w)
}
