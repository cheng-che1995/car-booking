package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
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

func init() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	for _, e := range schema {
		_, err := db.Exec(e)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
}
