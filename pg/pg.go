package pg

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	DBName   = "ui_test"
	Username = "ui_test"
	Password = "ui_test"

	RetryCount    = 3
	RetryTimeStep = time.Second * 10
)

type PG struct {
	db *gorm.DB
}

func (pg *PG) Connect(host string, port int) {
	var err error
	var db *gorm.DB
	for i := 0; i < RetryCount; i++ {
		db, err = gorm.Open(
			"postgres",
			fmt.Sprintf(
				"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
				host,
				port,
				DBName,
				Username,
				Password,
			),
		)
		if err != nil {
			log.Print(err)
			log.Print(fmt.Sprintf("Retry connecting to DB...%v", i+1))
			time.Sleep(RetryTimeStep)
		}
	}
	if err != nil {
		panic(err)
	}
	pg.db = db
	pg.initDBSQL()
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
	initSQLFile string = "./pg/users.sql"

	TableUsers Table = "users"

	FieldUserAcct      Field = "acct"
	FieldUserPwd       Field = "pwd"
	FieldUserFullname  Field = "fullname"
	FieldUserCreatedAt Field = "created_at"
	FieldUserUpdatedAt Field = "updated_at"

	FieldUserFullnameMaxLen = 50
)

type User struct {
	Acct       string    `json:"account"`
	Pwd        string    `json:"password"`
	Fullname   string    `json:"fullname"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

func (pg *PG) initDBSQL() {
	if pg.DB().HasTable(TableUsers.String()) {
		return
	}

	var sql []byte
	var err error
	if sql, err = ioutil.ReadFile(initSQLFile); err != nil {
		log.Fatal(err)
	}

	if res := pg.DB().Exec(string(sql)); res.Error != nil {
		err = res.Error
		log.Fatal(err)
	}
}
