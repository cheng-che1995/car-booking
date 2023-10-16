package main

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

type BoltRepository struct {
	dbPath string
}

func (b BoltRepository) openDB() (*bolt.DB, error) {
	return bolt.Open(b.dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
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
	//TODO: use better way to check empty string
	if *a.Username == "" {
		a.Username = nil
	}
	if *a.DateStart == "" {
		a.DateStart = nil
	}
	if *a.DateEnd == "" {
		a.DateEnd = nil
	}
	//
	// checkEmptyString(a.Username)
	// checkEmptyString(a.DateStart)
	// checkEmptyString(a.DateEnd)
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
			if (a.Username == nil || *a.Username == string(v)) &&
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

func (b BoltRepository) Delete(a *Appointment) error {
	db, err := b.openDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Appointments"))
		dateString := a.Date.Format("2006-01-02")
		if v := b.Get([]byte(dateString)); v == nil {
			return ErrNotFound
		} else if (v != nil) && (string(v) != a.Username) {
			return ErrUnauthorized
		}
		b.Delete([]byte(dateString))
		return nil

	})
	return err
}
