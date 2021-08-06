package pg

import "time"

const (
	Host     = "127.0.0.1"
	Port     = 5432
	DBName   = "ui_test"
	Username = "ui_test"
	Password = "ui_test"
)

const (
	UsersTable = "users"
)

type User struct {
	Acct       string
	Pwd        string
	Fullname   string
	Created_at time.Time
	Updated_at time.Time
}
