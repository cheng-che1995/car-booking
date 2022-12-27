package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
)

// database info

const (
	USERNAME = "root"
	PASSWORD = "root"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "testdb"
)

var schema = []string{
	`CREATE TABLE IF NOT EXISTS users(
		id INT NOT NULL AUTO_INCREMENT,
		uuid VARCHAR(36) NOT NULL,
		username VARCHAR(100) NOT NULL DEFAULT '',
		UNIQUE (uuid),
		PRIMARY KEY (id)
		)`,
	`CREATE TABLE IF NOT EXISTS appointments(
			id INT NOT NULL AUTO_INCREMENT,
			parant_uuid VARCHAR(36) NOT NULL,
			item VARCHAR(100) NOT NULL DEFAULT '',
			order_time DATETIME NOT NULL,
			create_by VARCHAR(100) NOT NULL DEFAULT '',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			FOREIGN KEY(parant_uuid) REFERENCES users (uuid)
			)`,
}

type Repository struct {
	db     *sql.DB
	config DbConfig
}
type DbConfig struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	DbName     string
}

func NewRepository(config DbConfig) *Repository {
	return &Repository{config: config}
}

func (m *Repository) initialize() error {
	tx, err := m.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, e := range schema {
		_, err := tx.Exec(e)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (m *Repository) OpenConn() error {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return err
	}

	m.db = db

	err = m.initialize()
	if err != nil {
		return err
	}
	return nil
}

func (m *Repository) CloseConn() error {
	m.db.Close()
	return nil
}

func (m *Repository) Create(username string, item string, orderTime string) error {
	tx, err := m.db.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if row := m.db.QueryRow("SELECT * FROM users NATRUAL JOIN appointments WHERE create_by = ? AND item = ? AND order_time = ?", username, item, orderTime); row == nil {
		return errors.New(ConflictResponse)
	}

	_, err = m.db.Exec("INSERT INTO users(uuid, username) VALUE (?,?)", uuid.NewV4(), username)
	if err != nil {
		return err
	}

	_, err = m.db.Exec("INSERT INTO appointments (item, order_time, create_by) VALUE (?,?,?)", item, orderTime, username)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (m *Repository) Search() error {

}
