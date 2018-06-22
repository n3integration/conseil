package sql

import (
	"database/sql"

	_ "{{ .Import }}"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("{{ .Driver }}", "{{ .Conn }}")
	if err != nil {
		panic(err)
	}

  	db.SetMaxOpenConns(50)
}
