package main

import (
	"errors"
	"strings"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
)

type GetUsersFilter struct {
	Id       int
	Uuid     string
	Username string
}

func (g *GetUsersFilter) GenerateQuery() (string, []interface{}) {
	query := "SELECT username FROM users"
	var conditions []string
	var whereValues []interface{}
	if g.Id != 0 {
		conditions = append(conditions, "id = ?")
		whereValues = append(whereValues, g.Id)
	}
	if g.Uuid != "" {
		conditions = append(conditions, "uuid = ?")
		whereValues = append(whereValues, g.Uuid)
	}
	if g.Username != "" {
		conditions = append(conditions, "username = ?")
		whereValues = append(whereValues, g.Username)
	}
	if len(whereValues) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query, whereValues
}

type GetCarsFilter struct {
	Id       int
	Uuid     string
	Plate    string
	UserUuid string
}

func (g *GetCarsFilter) GenerateQuery() (string, []interface{}) {
	query := "SELECT plate, user_uuid FROM cars"
	var conditions []string
	var whereValues []interface{}
	if g.Id != 0 {
		conditions = append(conditions, "id = ?")
		whereValues = append(whereValues, g.Id)
	}
	if g.Uuid != "" {
		conditions = append(conditions, "uuid = ?")
		whereValues = append(whereValues, g.Uuid)
	}
	if g.Plate != "" {
		conditions = append(conditions, "plate = ?")
		whereValues = append(whereValues, g.Plate)
	}
	if g.UserUuid != "" {
		conditions = append(conditions, "user_uuid = ?")
		whereValues = append(whereValues, g.UserUuid)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query, whereValues
}

type GetAppointmentsFilter struct {
	Id        int
	Uuid      string
	UserUuid  string
	CarUuid   string
	StartTime string
	EndTime   string
}

// TODO: add Time filter.
func (g *GetAppointmentsFilter) GenerateQuery() (string, []interface{}) {
	query := "SELECT * FROM appointments"
	var conditions []string
	var whereValues []interface{}
	if g.Id != 0 {
		conditions = append(conditions, "id = ?")
		whereValues = append(whereValues, g.Id)
	}
	if g.Uuid != "" {
		conditions = append(conditions, "uuid = ?")
		whereValues = append(whereValues, g.Uuid)
	}
	if g.UserUuid != "" {
		conditions = append(conditions, "user_uuid = ?")
		whereValues = append(whereValues, g.UserUuid)
	}
	if g.CarUuid != "" {
		conditions = append(conditions, "car_uuid = ?")
		whereValues = append(whereValues, g.CarUuid)
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	return query, whereValues
}

type CarBookingRepository interface {
	CreateUser(*User) error
	DeleteUser(*User) error
	GetUser(uuid string) (User, error)
	GetUsers(*GetUsersFilter) ([]User, error)
	CreateCar(*Car) error
	DeleteCar(*Car) error
	GetCar(uuid string) (Car, error)
	GetCars(*GetCarsFilter) ([]Car, error)
	CreateAppointment(*Appointment) error
	DeleteAppointment(*Appointment) error
	GetAppointment(uuid string) (Appointment, error)
	GetAppointments(*GetAppointmentsFilter) ([]Appointment, error)
}
