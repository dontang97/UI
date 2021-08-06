package pg

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Client struct {
	db *gorm.DB
}

func (cli *Client) Query() ([]User, error) {
	rows, err := cli.db.Table(UsersTable).Select("acct").Rows()
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		if err := cli.db.ScanRows(rows, &user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (cli *Client) Connect() {
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
	cli.db = db
}

func (cli *Client) Disconnect() {
	cli.db.Close()
}
