package db

import (
	"fmt"
	"github.com/fatih/color"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Sql struct {
	Db     *sqlx.DB
	Host   string
	Port   int
	Uid    string
	Pwd    string
	DbName string
}

func (s *Sql) Connect() {
	dataSource := fmt.Sprintf(`%s:%s@tcp(%s:%d)/%s`, s.Uid, s.Pwd, s.Host, s.Port, s.DbName)
	db, err := sqlx.Connect("mysql", dataSource)
	if err != nil {
	} else {
		color.Green("Connect success db")
		s.Db = db
	}
}

func (s *Sql) Close() {
	s.Db.Close()
}
