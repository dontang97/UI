package ui

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dontang97/ui/pg"
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
		Table(pg.TableUsers.ToString()).
		Select(pg.FieldUserAcct.ToString()).
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

	for _, user := range users {
		if _, err := w.Write([]byte(user.Acct + "\n")); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

///////////////////////////////////////////////////////
//////    GET /ui/v1/user?fullname={fullname}    //////
///////////////////////////////////////////////////////
var FullnameQueryHdl QueryUserHandlerFunc = func(ui *UI, args ...interface{}) ([]pg.User, error) {
	rows, err := ui.DB().
		Table(pg.TableUsers.ToString()).
		Select("acct").
		Where("fullname = ?", args[0]).Rows()
	if err != nil {
		return nil, err
	}

	return scanUsers(ui, rows)
}

func (ui *UI) FullnameQuery(w http.ResponseWriter, r *http.Request) {
	users, err := FullnameQueryHdl(ui, r.URL.Query().Get("fullname"))

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, user := range users {
		if _, err := w.Write([]byte(user.Acct + "\n")); err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
