package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func init() {

	conn := "root:root@tcp(127.0.0.1:3306)/test"

	db, err := sql.Open("mysql", conn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
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
