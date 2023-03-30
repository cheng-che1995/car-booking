package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		username VARCHAR(100) NOT NULL UNIQUE,
		password BLOB NOT NULL
		)`,
	`CREATE TABLE IF NOT EXISTS cars(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		plate VARCHAR(12) NOT NULL UNIQUE,
		user_uuid VARCHAR(36) NOT NULL,
		FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE
	)`,
	`CREATE TABLE IF NOT EXISTS appointments(
		id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		uuid VARCHAR(36) NOT NULL UNIQUE,
		user_uuid VARCHAR(36) NOT NULL,
		car_uuid VARCHAR(36) NOT NULL,
		FOREIGN KEY (user_uuid) REFERENCES users (uuid) ON DELETE CASCADE,
		FOREIGN KEY (car_uuid) REFERENCES cars (uuid) ON DELETE CASCADE,
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
func (m *Repository) AuthUser(u *User) (bool, error) {
	if err := u.Validate(); err != nil {
		return false, err
	}
	pwd, err := u.HashPassword()
	if err != nil {
		return false, err
	}
	q := `SELECT uuid FROM users WHERE username = ? AND password = ?`
	row := m.db.QueryRow(q, u.Username, pwd)
	if err := row.Scan(&u.Uuid); err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (m *Repository) CreateUser(u *User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if ok, err := u.CheckPassword(); !ok {
		return err
	}
	if u.Uuid == "" {
		u.GenerateUuid()
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

func (m *Repository) GetUser(uuid string) (*User, error) {
	if &uuid == nil {
		return nil, nil
	}
	query := "SELECT username FROM users WHERE uuid = ?"
	rows, err := m.db.Query(query, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	user := User{}
	for rows.Next() {
		if err := rows.Scan(&user.Username); err != nil {
			return nil, err
		}
	}
	return &user, nil
}
func (m *Repository) GetUsers(g *GetUsersFilter) ([]User, error) {
	if g == nil {
		return nil, nil
	}
	users := []User{}
	query, whereValues := g.GenerateQuery()
	rows, err := m.db.Query(query, whereValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := User{}
		if err := rows.Scan(&user.Username); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, err
}

func (m *Repository) CreateCar(c *Car) error {
	if err := c.Validate(); err != nil {
		return err
	}

	if c.Uuid == "" {
		c.GenerateUuid()
	}

	q := `INSERT INTO cars SET plate = ?, uuid = ?, user_uuid = ?`
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
func (m *Repository) GetCar(uuid string) (*Car, error) {
	if &uuid == nil {
		return nil, nil
	}
	query := `SELECT plate, user_uuid FROM cars WHERE uuid = ?`
	rows, err := m.db.Query(query, uuid)
	if err != nil {
		return nil, err
	}
	car := Car{}
	if err := rows.Scan(&car.Plate, &car.UserUuid); err != nil {
		return nil, err
	}
	return &car, nil
}
func (m *Repository) GetCars(g *GetCarsFilter) ([]Car, error) {
	if g == nil {
		return nil, nil
	}
	cars := []Car{}
	query, whereValues := g.GenerateQuery()
	rows, err := m.db.Query(query, whereValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		car := Car{}
		if err := rows.Scan(&car.Plate, &car.UserUuid); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func (m *Repository) CreateAppointment(a *Appointment) error {
	if err := a.Vaildate(); err != nil {
		return err
	}
	s := `SELECT COUNT(*) FROM appointments 
			WHERE car_uuid = ? 
			AND start_time < ?
			OR end_time > ?`
	count := 0
	err := m.db.QueryRow(s, a.CarUuid, a.EndTime, a.StartTime).Scan(&count)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("預約時間重疊，請重新選擇！")
	}
	if a.Uuid == "" {
		a.generateUuid()
	}

	q := `INSERT INTO appointments SET start_time = ?, end_time = ?, uuid = ?, user_uuid = ?, car_uuid = ?`
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

func (m *Repository) GetAppointment(uuid string) (*Appointment, error) {
	if uuid == "" || &uuid == nil {
		return nil, nil
	}
	appointment := Appointment{}
	q := `SELECT user_uuid, car_uuid, start_time, end_time FROM appointments WHERE uuid = ?`
	rows, err := m.db.Query(q, uuid)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	if err := rows.Scan(&appointment.UserUuid, &appointment.CarUuid, &appointment.StartTime, &appointment.EndTime); err != nil {
		return nil, err
	}
	return &appointment, nil
}

// TODO: 動態變更SELECT欄位，用map
func (m *Repository) GetAppointments(fields []string, g *GetAppointmentsFilter) ([]Appointment, error) {
	if g == nil {
		return nil, nil
	}
	appointment := Appointment{}
	appointments := []Appointment{}
	fieldsMap := map[string]interface{}{
		"appointment_uuid": &appointment.Uuid,
		"user_uuid":        &appointment.UserUuid,
		"car_uuid":         &appointment.CarUuid,
		"start_time":       &appointment.StartTime,
		"end_time":         &appointment.EndTime,
	}
	var scanArgs []interface{}
	for _, v := range fields {
		if val, ok := fieldsMap[v]; ok {
			scanArgs = append(scanArgs, val)
		}
	}
	query, whereValues := g.GenerateQuery(fields)
	rows, err := m.db.Query(query, whereValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}
