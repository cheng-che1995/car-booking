package main

import (
	"context"
	"database/sql"
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
		username VARCHAR(100) NOT NULL DEFAULT '',
		PRIMARY KEY (id)
		)`,
	`CREATE TABLE IF NOT EXISTS appointments(
			id INT NOT NULL,
			uuid VARCHAR(36) NOT NULL,
			item VARCHAR(100) NOT NULL DEFAULT '',
			order_at DATETIME NOT NULL,
			create_by VARCHAR(100) NOT NULL DEFAULT '',
			create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(id) REFERENCES users (id)
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
			fmt.Printf("Range schema failed: %v\n", err)
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
		fmt.Printf("OpenConn failed: %v\n", err)
		return err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		fmt.Printf("Ping failed: %v\n", err)
		return err
	}

	m.db = db

	err = m.initialize()
	if err != nil {
		fmt.Printf("init failed: %v\n", err)
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
		fmt.Printf("BeginTx failed: %v\n", err)
		return err
	}
	defer tx.Rollback()

	// TODO: Make sure the order won't conflict
	row, err := tx.Query("SELECT * FROM appointments WHERE item = ? AND order_at = ?", item, orderTime)
	if err != nil {
		fmt.Printf("check exsits part failed: %v\n", err)
		return err
	}
	if row.Next() {
		//exsits
		return ErrConflict
	}

	stmt1, err := tx.Prepare("INSERT users SET username=?")
	if err != nil {
		fmt.Printf("Prepare insert table users failed: %v\n", err)
		return err
	}

	stmt2, err := tx.Prepare("INSERT appointments SET id=?, uuid=?, item=?, order_at=?, create_by=?")
	if err != nil {
		fmt.Printf("Prepare insert table appointments failed: %v\n", err)
		return err
	}

	res, err := stmt1.Exec(username)
	if err != nil {
		fmt.Printf("Insert table users failed: %v\n", err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Get last insert id failed: %v\n", err)
		return err
	}

	res2, err := stmt2.Exec(id, uuid.NewV4(), item, orderTime, username)
	if err != nil {
		fmt.Printf("Insert table appointments failed: %v\n", err)
		return err
	}
	res2.RowsAffected()

	if err = tx.Commit(); err != nil {
		fmt.Printf("tx.Commit failed: %v\n", err)
		return err
	}
	return nil
}

// func (m *Repository) Search() error {
// 	rows, err := m.db.Exec("SELECT u.id, u.uuid, a.item, a.order_time, a.create_by, a.create_time FROM users AS u JOIN appointments AS a ON u.id = a.id;")
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

//select u.id, u.uuid, a.item, a.order_at, a.create_by, a.create_time from users as u join appointments as a on u.id = a.id;
