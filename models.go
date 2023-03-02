package main

import (
	"crypto/sha256"
	"encoding"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Appointment struct {
	StartTime time.Time
	EndTime   time.Time
	UserUuid  string
	CarUuid   string
}

type User struct {
	Uuid     string
	Username string
	Password string
}

func (u *User) GenerateUuid() {
	u.Uuid = uuid.NewV4().String()
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

func (u User) Validate() error {
	if u.Username == "" {
		return errors.New("使用者名稱不得為空值！")
	}
	if u.Password == "" {
		return errors.New("密碼不得為空值！")
	}
	if u.Uuid == "" {
		return errors.New("uuid不得為空值！")
	}
	return nil
}

type Car struct {
	Plate    string
	Uuid     string
	UserUuid string
}

func (c *Car) GenerateUuid() {
	c.Uuid = uuid.NewV4().String()
}

func (c Car) Validate() error {
	if c.Plate == "" {
		return errors.New("車牌不得為空值！")
	}
	if c.Uuid == "" {
		return errors.New("uuid不得為空值！")
	}
	if c.UserUuid == "" {
		return errors.New("UserUuid不得為空值！")
	}
	return nil
}
