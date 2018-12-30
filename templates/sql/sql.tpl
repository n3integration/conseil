package sql

import (
    "database/sql"

    _ "{{ .Import }}"
)

var db *sql.DB

func Open() error {
    var err error
    db, err = sql.Open("{{ .Driver }}", "{{ .Conn }}")
    if err != nil {
        return nil, err
    }

    if err := db.Ping(); err != nil {
        return nil, err
    }

    db.SetMaxOpenConns(50)
    return db, nil
}

func Close() error {
    if db != nil {
        return db.Close()
    }
    return nil
}