package main

import (
	"crypto/sha256"
	"encoding"
	"errors"
	"fmt"
	"time"
	"unicode"

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

func (u User) CheckPassword() (bool, error) {
	const (
		minLenth = 10
		maxLenth = 25
	)
	var (
		ErrTooShort = errors.New(fmt.Sprintf("密碼長度不足，請大於%d字元!", minLenth))
		ErrTooLong  = errors.New(fmt.Sprintf("密碼長度過長，請小於%d字元！", maxLenth))
		hasSpecial  = false
		hasUpper    = false
		hasLower    = false
		hasDigit    = false
	)
	if len(u.Password) < minLenth {
		return false, ErrTooShort
	}
	if len(u.Password) > maxLenth {
		return false, ErrTooLong
	}
	for _, v := range u.Password {
		if unicode.IsSymbol(v) {
			hasSpecial = true
		}
		if unicode.IsUpper(v) {
			hasUpper = true
		}
		if unicode.IsLower(v) {
			hasLower = true
		}
		if unicode.IsNumber(v) {
			hasDigit = true
		}
		if hasSpecial && hasUpper && hasLower && hasDigit {
			break
		}
	}
	return true, nil
}

var ErrHashFailed = errors.New("密碼雜湊轉換失敗！")

func (u User) HashPassword() ([]byte, error) {
	hash := sha256.New()
	hash.Write([]byte(u.Password))
	marshaler, ok := hash.(encoding.BinaryMarshaler)
	if !ok {
		return nil, ErrHashFailed
	}
	pwd, err := marshaler.MarshalBinary()
	if err != nil {
		return nil, err
	}
	return pwd, nil
}

var (
	ErrUsernameEmpty = errors.New("使用者名稱不得為空值！")
	ErrPasswordEmpty = errors.New("密碼不得為空值！")
)

func (u *User) Validate() error {
	if u == nil {
		return nil
	}
	if u.Username == "" {
		return ErrUsernameEmpty
	}
	if u.Password == "" {
		return ErrPasswordEmpty
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

// 註：車牌命名規範為 AAA-0001 ~ ZZZ-9999（實際上根據交通部定義，除了須以車牌號碼範圍區別車種，還有額外的規範，如：避免出現4444或是不雅字母組合等，在此忽略不計。）
func (c *Car) Validate() error {
	if c == nil {
		return nil
	}
	if c.UserUuid == "" {
		return errors.New("UserUuid不得為空值！")
	}
	if c.Plate == "" {
		return errors.New("車牌不得為空值！")
	}
	if len(c.Plate) != 8 {
		return errors.New("車牌長度不符合命名規範！")
	}
	for _, r := range c.Plate {
		switch r {
		case 0, 1, 2:
			if !unicode.IsLetter(r) || !unicode.IsUpper(r) {
				return errors.New("車牌格式不符合命名規範！")
			}
		case 3:
			if r != '-' {
				return errors.New("車牌格式不符合命名規範！")
			}
		default:
			if !unicode.IsDigit(r) {
				return errors.New("車牌格式不符合命名規範！")
			}
		}
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
