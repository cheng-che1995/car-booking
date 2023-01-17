package main

import (
	"errors"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
)

type SearchFilter struct {
	Username  *string
	Item      *string
	DateStart *string
	DateEnd   *string
}

type AppointmentRepository interface {
	Create(*Appointment) error
	Search(*SearchFilter) ([]Appointment, error)
	Delete(*Appointment) error
}

// TODO: 找其他地方放
func checkEmptyString(a *string) {
	if len(*a) > 0 {
		a = nil
	}
}
