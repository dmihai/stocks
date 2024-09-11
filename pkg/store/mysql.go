package store

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type Conn struct {
	db *sql.DB
}

func NewConn(addr, user, pass, database string) (*Conn, error) {
	cfg := mysql.Config{
		User:                 user,
		Passwd:               pass,
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               database,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	conn, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	return &Conn{
		db: conn,
	}, nil
}

func (c *Conn) Ping() error {
	return c.db.Ping()
}
