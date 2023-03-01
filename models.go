package main

import (
	"crypto/sha256"
	"encoding"
	"errors"
	"time"
)

type Appointment struct {
	StartTime time.Time
	EndTime   time.Time
	UserUuid  string
	CarUuid   string
}

type Car struct {
	Plate    string
	UserUuid string
}

type User struct {
	Uuid     string
	Username string
	Password string
}

func (u User) HashPassword() ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(u.Password))
	marshaler, ok := hash.(encoding.BinaryMarshaler)
	if !ok {
		return nil, errors.New("hash failed.")
	}
	pwd, err := marshaler.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pwd, nil
}
