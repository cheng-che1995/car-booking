package main

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
)

type SearchFilter struct {
	Username  *string
	DateStart *string
	DateEnd   *string
}

type AppointmentRepository interface {
	Create(*Appointment) error
	Search(*SearchFilter) ([]Appointment, error)
	Delete(*Appointment) error
}

type BoltRepository struct {
	dbPath string
}

func (b BoltRepository) openDB() (*bolt.DB, error) {
	return bolt.Open(b.dbPath, 0600, nil)
}

func (b BoltRepository) Create(a *Appointment) error {
	db, err := b.openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Appointments"))
		dateString := a.Date.Format("2006-01-02")
		if b.Get([]byte(dateString)) == nil {
			b.Put([]byte(dateString), []byte(a.Username))
			return nil
		} else {
			return errors.New(ConflictResponse)
		}
	})
	return err
}

func (b BoltRepository) Search(a *SearchFilter) ([]Appointment, error) {
	db, err := b.openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var startDate, endDate time.Time
	FilteredAppointments := []Appointment{}
	if a.DateStart != nil {
		startDate, _ = time.Parse("2006-01-02", *a.DateStart)
	}
	if a.DateEnd != nil {
		endDate, _ = time.Parse("2006-01-02", *a.DateEnd)
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Appointments"))
		b.ForEach(func(k, v []byte) error {
			kt, _ := time.Parse("2006-01-02", string(k))
			if (a.Username == nil || string(v) != *a.Username) &&
				(a.DateStart == nil || (startDate.Before(kt) || startDate.Equal(kt))) &&
				(a.DateEnd == nil || (endDate.After(kt)) || endDate.Equal(kt)) {
				FilteredAppointments = append(FilteredAppointments, Appointment{Username: string(v), Date: kt})
			}
			return nil
		})
		return nil
	})
	return FilteredAppointments, err
}
