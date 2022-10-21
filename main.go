package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type LoginResponse struct {
	Token   string `json:"token"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type AppointmentsResponse struct {
	Appointments []Appointment `json:"appointments"`
	Status       string        `json:"status"`
	Message      string        `json:"message"`
}
type AppointmentsResponse1 struct {
	Appointments []Appointment1 `json:"appointments"`
	Status       string         `json:"status"`
	Message      string         `json:"message"`
}

const (
	SuccessResponse      string = "success"
	ConflictResponse     string = "conflict"
	NotFoundResponse     string = "notFound"
	UnauthorizedResponse string = "unauthorized"
)

type Appointment struct {
	Username string
	Date     time.Time
}
type Appointment1 struct {
	Item string
	Date time.Time
}

var appoint2 []Appointment

// type CustomFunc func(echo.Context) error
var users map[string]string = map[string]string{
	"tony":   "4321",
	"wilson": "1234",
}

func showUsers(c echo.Context) error {
	var account []string
	for u := range users {
		account = append(account, u)
	}
	return c.JSON(http.StatusOK, account)
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	unauthorizedMessage := "密碼錯誤！"
	successMessage := "驗證成功！"
	pw, ok := users[username]

	if pw != password || !ok {
		return c.JSON(http.StatusUnauthorized, LoginResponse{Status: UnauthorizedResponse, Message: unauthorizedMessage})
	}

	claims := &jwtCustomClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: t, Status: SuccessResponse, Message: successMessage})
}

func createAppointments(c echo.Context) error {
	db, err := bolt.Open("car-booking.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var tx = &bolt.Tx{}
	b := tx.Bucket([]byte("Appointments"))
	token := c.Get("token").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)
	username := claims.Name
	selectedItem := c.FormValue("selectedItem")
	selectedDate := c.FormValue("selectedDate")
	t, err := time.Parse("2006-01-02", selectedDate)
	if err != nil {
		return err
	}
	errMessage := fmt.Sprintf("%s，此日期已被預訂，請您重新選擇其他日期！", username)
	successMessage := fmt.Sprintf("預約成功！%s，您的預約日期為： %s", username, t.Format("2006-01-02"))
	appoint1 := Appointment1{}

	b.ForEach(func(k, v []byte) error {
		json.Unmarshal(v, &appoint1)
		if appoint1.Date == t && appoint1.Item == selectedItem {
			return c.JSON(http.StatusConflict, AppointmentsResponse{Status: ConflictResponse, Message: errMessage})
		}
		return nil
	})
	data, _ := json.Marshal(appoint2)
	b.Put([]byte(username), data)
	return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
}

func searchAppointments(c echo.Context) error {
	db, err := bolt.Open("car-booking.db", 0600, nil)
	if err != nil {
		return nil
	}
	defer db.Close()
	var tx = &bolt.Tx{}
	b := tx.Bucket([]byte("Appointments"))

	filterByUsername := c.FormValue("filterByUsername")
	filterByItem := c.FormValue("filterByItem")
	filterByDateStart := c.FormValue("filterByDateStart")
	filterByDateEnd := c.FormValue("filterByDateEnd")
	startDate, _ := time.Parse("2006-01-02", filterByDateStart)
	endDate, _ := time.Parse("2006-01-02", filterByDateEnd)
	appoint1 := Appointment1{}
	var filteredAppointments []Appointment1
	b.ForEach(func(k, v []byte) error {
		err = json.Unmarshal(v, &appoint1)
		if err != nil {
			return err
		}
		if (filterByUsername != "" && string(k) != filterByUsername) ||
			(filterByItem != "" && appoint1.Item != filterByItem) ||
			(filterByDateStart != "" && appoint1.Date.Before(startDate)) ||
			(filterByDateEnd != "" && appoint1.Date.After(endDate)) {
			filteredAppointments = append(filteredAppointments, appoint1)
		}
		return nil
	})
	return c.JSON(http.StatusOK, AppointmentsResponse1{Status: SuccessResponse, Appointments: filteredAppointments})
}

func cancelAppointments(c echo.Context) error {
	token := c.Get("token").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)
	username := claims.Name

	selectedDate := c.FormValue("selectedDate")
	t, err := time.Parse("2006-01-02", selectedDate)
	if err != nil {
		return err
	}
	successMessage := fmt.Sprintf("取消成功！%s，您已將 %s預約取消！", username, t.Format("2006-01-02"))
	errMessage := fmt.Sprintf("此%s日期不屬於%s您的預約！", t.Format("2006-01-02"), username)
	notFoundMessage := fmt.Sprintf("查無此預約！%s請您重新選擇日期！", username)

	var found bool

	for i := range appoint2 {
		if t == appoint2[i].Date {
			if appoint2[i].Username != username {
				return c.JSON(http.StatusConflict, AppointmentsResponse{Status: ConflictResponse, Message: errMessage})
			}
			found = true
			for j := range appoint2[i : len(appoint2)-1] {
				appoint2[i+j].Date = appoint2[i+j+1].Date
			}
			appoint2 = appoint2[:len(appoint2)-1]
			break
		}

	}
	if found {
		return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
	}
	return c.JSON(http.StatusNotFound, AppointmentsResponse{Status: NotFoundResponse, Message: notFoundMessage})
}

func main() {
	// Create a db named "car-booking.db" in current directory.
	// It will be created if doesn't exsit.
	// And keep it connected.
	db, err := bolt.Open("car-booking.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Create a bucket(table) named "appointments".
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Appointments"))
		if err != nil {
			return fmt.Errorf("create bucket err: %s", err)
		}
		return nil
	})
	// Close the connection.
	db.Close()

	// var e *echo.Echo // = echo.New()
	e := echo.New()
	e.POST("/login", login)
	e.GET("/users", showUsers)
	b := e.Group("/booking")
	b.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		// Claims:     &jwtCustomClaims{},
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
		ContextKey: "token",
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			return c.String(http.StatusUnauthorized, err.Error())
		},
	}))
	b.POST("/appointments", createAppointments)
	b.GET("/appointments", searchAppointments)
	b.DELETE("/appointments", cancelAppointments)
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(":1323"))
}
