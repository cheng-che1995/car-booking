package main

import (
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

func init() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

}

// CREATE TABLE `users`(
// 	`uuid` VARCHAR(36) PRIMARY KEY NOT NULL,
// 	`username` VARCHAR(64) DEFAULT NULL,
// )

// CREATE TABLE `appointments`(
// 	`uuid` VARCHAR(36) PRIMARY KEY NOT NULL,
// 	`item` VVARCHAR(64) DEFAULT NULL,
// 	`order_by` VARCHAR(64) DEFAULT NULL,
// 	`created_time` DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
// 	`selected_time` DATETIME NULL DEFAULT NULL,
// )
