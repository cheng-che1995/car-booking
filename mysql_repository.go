package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL,
		username VARCHAR(100) NOT NULL,
		password VARCHAR(255) NOT NULL
		)`,
	`CREATE TABLE IF NOT EXISTS cars(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL,
		plate VARCHAR(12) NOT NULL UNIQUE,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS appointments(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL,
		FOREIGN KEY (car_id) REFERENCES cars (id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
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

func (m *Repository) CreateUser(u *User) error {
	if u == nil {
		return nil
	}

	if u.Uuid == "" {
		u.GenerateUuid()
	}

	if err := u.Validate(); err != nil {
		return err
	}

	pwd, err := u.HashPassword()
	if err != nil {
		return err
	}
	q := `INSERT INTO users SET username = ?, password = ?, uuid = ?`

	if _, err := m.db.Exec(q, u.Username, pwd, u.Uuid); err != nil {
		return err
	}
	return nil
}

func (m *Repository) DeleteUser(u *User) error {
	if u == nil {
		return nil
	}
	q := "DELETE FROM users WHERE uuid = ?"
	if _, err := m.db.Exec(q, u.Uuid); err != nil {
		return err
	}
	return nil
}

func (m *Repository) CreateCar(c *Car) error {
	if c == nil {
		return nil
	}

	if c.Uuid == "" {
		c.GenerateUuid()
	}

	if err := c.Validate(); err != nil {
		return err
	}

	q := `INSERT INTO cars SET plate = ?, uuid = ?, user_id = (SELECT id FROM users WHERE uuid = ?)`
	if _, err := m.db.Exec(q, c.Plate, c.Uuid, c.UserUuid); err != nil {
		return err
	}
	return nil
}

func (m *Repository) DeleteCar(c *Car) error {
	if c == nil {
		return nil
	}
	q := `DELETE FROM cars WHERE uuid = ?`
	if _, err := m.db.Exec(q, c.Uuid); err != nil {
		return err
	}
	return nil
}

func (m *Repository) GetCars(g *GetCarsFilter) ([]Car, error) {
	cars := []Car{}
	if g == nil {
		return nil, nil
	}
	q := `SELECT uuid, plate FROM cars`
	rows, err := m.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		car := Car{}
		if err := rows.Scan(&car.Uuid, &car.Plate); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func (m *Repository) CreateAppointment(a *Appointment) error {
	if a == nil {
		return nil
	}
	if err := a.Vaildate(); err != nil {
		return err
	}
	if a.Uuid == "" {
		a.generateUuid()
	}

	q := `INSERT INTO appointments SET start_time = ?, end_time = ?, uuid = ?, user_id = (SELECT id FROM users WHERE uuid = ?), car_id = (SELECT id FROM cars WHERE uuid = ?)`
	if _, err := m.db.Exec(q, a.StartTime, a.EndTime, a.Uuid, a.UserUuid, a.CarUuid); err != nil {
		return err
	}
	return nil
}

func (m *Repository) DeleteAppointment(a *Appointment) error {
	if a == nil {
		return nil
	}
	q := `DELETE FROM appointments WHERE uuid = ?`
	if _, err := m.db.Exec(q, a.Uuid); err != nil {
		return nil
	}
	return nil
}

func (m *Repository) GetAppointments(g *GetAppointmentsFilter) ([]Appointment, error) {
	appointment := Appointment{}

}

func (m *Repository) Create(username string, item string, date string) error {
	tx, err := m.db.BeginTx(context.TODO(), nil)
	if err != nil {
		fmt.Printf("BeginTx failed: %v\n", err)
		return err
	}
	defer tx.Rollback()

	// TODO: Make sure the order won't conflict
	row, err := tx.Query("SELECT * FROM appointments WHERE item = ? AND order_at = ?", item, date)
	if err != nil {
		fmt.Printf("check exsits part failed: %v\n", err)
		return err
	}
	if row.Next() {
		//exsits
		return ErrConflict
	}

	stmt1, err := tx.Prepare("INSERT appointments SET uuid = ?, item = ?, order_at = ?, order_by= ?")
	if err != nil {
		fmt.Printf("Prepare insert table appointments failed: %v\n", err)
		return err
	}
	res1, err := stmt1.Exec(uuid.NewV4(), item, date, username)
	if err != nil {
		fmt.Printf("Insert table appointments failed: %v\n", err)
		return err
	}

	id, _ := res1.LastInsertId()

	stmt2, err := tx.Prepare("INSERT users SET appointment_id = ?, username = ?")
	if err != nil {
		fmt.Printf("Prepare insert table users failed: %v\n", err)
		return err
	}
	res2, err := stmt2.Exec(id, username)
	if err != nil {
		fmt.Printf("Insert table users failed: %v\n", err)
	}
	res2.RowsAffected()

	if err = tx.Commit(); err != nil {
		fmt.Printf("tx.Commit failed: %v\n", err)
		return err
	}
	return nil
}

func (m *Repository) Search(a *SearchFilter) ([]NewAppointment, error) {
	var (
		FilteredAppointments []NewAppointment
		dateStart            time.Time
		dateEnd              time.Time
		id                   int
		uuid                 string
		item                 string
		orderAt              time.Time
		createBy             string
		createTime           time.Time
	)
	checkEmptyString(a.Username)
	checkEmptyString(a.Item)
	checkEmptyString(a.DateStart)
	checkEmptyString(a.DateEnd)

	if a.DateStart != nil {
		dateStart, _ = time.Parse("2006-01-02", *a.DateStart)
	}
	if a.DateEnd != nil {
		dateEnd, _ = time.Parse("2006-01-02", *a.DateEnd)
	}

	rows, err := m.db.Query("SELECT * FROM appointments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&id, &uuid, &item, &orderAt, &createBy, &createTime)
		if err != nil {
			fmt.Printf("Scan failed: %v\n", err)
			return nil, err
		}
		if (*a.Username == "" || *a.Username == createBy) &&
			(*a.Item == "" || *a.Item == item) &&
			(dateStart.Before(orderAt) || dateStart.Equal(orderAt) || *a.DateStart == "") &&
			(dateEnd.After(orderAt) || dateEnd.Equal(orderAt) || *a.DateEnd == "") {
			FilteredAppointments = append(FilteredAppointments, NewAppointment{Username: createBy, Item: item, Date: orderAt})
		}
	}
	return FilteredAppointments, nil
}

func (m *Repository) Delete(username string, item string, date string) error {
	tx, err := m.db.BeginTx(context.TODO(), nil)
	if err != nil {
		fmt.Printf("BeginTx failed: %v\n", err)
		return err
	}
	defer tx.Rollback()

	rows, err := m.db.Query("SELECT * FROM appointments WHERE item = ? AND order_at = ?", item, date)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return ErrNotFound
	}
	rows2, err := m.db.Query("SELECT * FROM appointments WHERE order_by = ? AND item = ? AND order_at = ?", username, item, date)
	if err != nil {
		return err
	}
	if !rows2.Next() {
		return ErrUnauthorized
	}

	res, err := tx.Exec("DELETE FROM appointments WHERE order_by = ? AND item = ? AND order_at = ?", username, item, date)
	if err != nil {
		fmt.Printf("Delete from appointments failed: %v\n", err)
		return err
	}

	res.RowsAffected()

	if err = tx.Commit(); err != nil {
		fmt.Printf("tx.Commit failed: %v\n", err)
		return err
	}
	return nil
}
