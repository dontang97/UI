package ui

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/dontang97/ui/pg"
	"github.com/gorilla/mux"
)

type QueryUserHandlerFunc func(*UI, ...interface{}) ([]pg.User, error)

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

var UserInfoQueryHdl QueryUserHandlerFunc = func(ui *UI, args ...interface{}) ([]pg.User, error) {
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
	users, err := UserInfoQueryHdl(ui, acct)

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
