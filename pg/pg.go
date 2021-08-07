package pg

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	Host     = "127.0.0.1"
	Port     = 5432
	DBName   = "ui_test"
	Username = "ui_test"
	Password = "ui_test"
)

type PG struct {
	db *gorm.DB
}

func (pg *PG) Connect() {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
			Host,
			Port,
			DBName,
			Username,
			Password,
		),
	)
	if err != nil {
		panic(err)
	}
	pg.db = db
}

func (pg *PG) DB() *gorm.DB {
	return pg.db
}

func (pg *PG) Disconnect() {
	pg.db.Close()
}

type Table string

func (t Table) String() string {
	return string(t)
}

type Field string

func (f Field) String() string {
	return string(f)
}

const (
	TableUsers Table = "users"

	FieldUserAcct      Field = "acct"
	FieldUserPwd       Field = "pwd"
	FieldUserFullname  Field = "fullname"
	FieldUserCreatedAt Field = "created_at"
	FieldUserUpdatedAt Field = "updated_at"
)

type User struct {
	Acct       string    `json:"account"`
	Pwd        string    `json:"password"`
	Fullname   string    `json:"fullname"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}
