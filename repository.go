package main

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
)

type GetAppointmentsFilter struct {
	Username  *string
	Item      *string
	DateStart *string
	DateEnd   *string
}

type GetCarsFilter struct {
	Id      int
	Uuid    string
	Plate   string
	User_id string
}

type CarBookingRepository interface {
	CreateUser(*User) error
	DeleteUser(*User) error
	CreateCar(*Car) error
	DeleteCar(*Car) error
	GetCars(*GetCarsFilter) ([]Car, error)
	CreateAppointment(*Appointment) error
	DeleteAppointment(*Appointment) error
	GetAppointments(*GetAppointmentsFilter) ([]Appointment, error)
}
