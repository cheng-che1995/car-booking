package main

import (
	"crypto/sha256"
	"encoding"
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	Uuid     string `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"-"`
}

func (u *User) GenerateUuid() {
	u.Uuid = uuid.NewV4().String()
}

func (u User) HashPassword() ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(u.Password))
	marshaler, ok := hash.(encoding.BinaryMarshaler)
	if !ok {
		return nil, errors.New("hash password failed.")
	}
	pwd, err := marshaler.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pwd, nil
}

func (u *User) Validate() error {
	if u == nil {
		return nil
	}
	if u.Username == "" {
		return errors.New("使用者名稱不得為空值！")
	}
	if u.Password == "" {
		return errors.New("密碼不得為空值！")
	}
	return nil
}

type Car struct {
	Plate    string `json:"plate"`
	Uuid     string `json:"uuid"`
	UserUuid string `json:"user_uuid"`
}

func (c *Car) GenerateUuid() {
	c.Uuid = uuid.NewV4().String()
}

func (c *Car) Validate() error {
	if c == nil {
		return nil
	}
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

type Appointment struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Uuid      string    `json:"uuid"`
	UserUuid  string    `json:"user_uuid"`
	CarUuid   string    `json:"car_uuid"`
}

func (a *Appointment) generateUuid() {
	a.Uuid = uuid.NewV4().String()
}

func (a *Appointment) Vaildate() error {
	if a == nil {
		return nil
	}
	if a.StartTime.Format("2006-01-02") == "" || a.StartTime.Format("2006-01-02") == "" {
		return errors.New("預約起始、結束時間不得為空值！")
	}
	if a.StartTime.After(a.EndTime) {
		return errors.New("預約起始、結束時間順序有誤，請重新選擇！")
	}
	if a.UserUuid == "" {
		return errors.New("UserUuid不得為空值！")
	}
	if a.CarUuid == "" {
		return errors.New("CarsUuid不得為空值！")
	}
	return nil
}
