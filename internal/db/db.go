package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"golang/tutorial/todo/internal/config"
)

type Client struct {
	db *sql.DB
}

func (c *Client) Connect(cfg config.DBConfig) error {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}
	c.db = db
	return nil
}

func (c *Client) DB() *sql.DB {
	return c.db
}

func (c *Client) Close() error {
	return c.db.Close()
}
