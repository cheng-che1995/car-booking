package main

import (
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
	token := c.Get("token").(*jwt.Token)
	claims := token.Claims.(*jwtCustomClaims)
	username := claims.Name
	selectedDate := c.FormValue("selectedDate")
	t, err := time.Parse("2006-01-02", selectedDate)
	if err != nil {
		return err
	}
	errMessage := fmt.Sprintf("%s，此日期已被預訂，請您重新選擇其他日期！", username)
	successMessage := fmt.Sprintf("預約成功！%s，您的預約日期為： %s", username, t.Format("2006-01-02"))
	//TODO: Put this elsewhere
	br := BoltRepository{dbPath: "car-booking.db"}

	selectAppointments := Appointment{
		Username: username,
		Date:     t,
	}
	if err := br.Create(&selectAppointments); err != nil {
		return c.JSON(http.StatusConflict, AppointmentsResponse{Status: ConflictResponse, Message: errMessage})
	}
	return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
}

func searchAppointments(c echo.Context) error {
	filterByUsername := c.FormValue("filterByUsername")
	filterByDateStart := c.FormValue("filterByDateStart")
	filterByDateEnd := c.FormValue("filterByDateEnd")
	br := BoltRepository{dbPath: "car-booking.db"}
	//TODO:
	selectedFilter := SearchFilter{
		Username:  &filterByUsername,
		DateStart: &filterByDateStart,
		DateEnd:   &filterByDateEnd,
	}
	FilteredAppointments, err := br.Search(&selectedFilter)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Appointments: FilteredAppointments})
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

	br := BoltRepository{dbPath: "car-booking.db"}
	selectAppointments := Appointment{
		Username: username,
		Date:     t,
	}

	if err = br.Delete(&selectAppointments); err == nil {
		return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
	} else if err == ErrNotFound {
		return c.JSON(http.StatusNotFound, AppointmentsResponse{Status: NotFoundResponse, Message: notFoundMessage})
	} else if err == ErrUnauthorized {
		return c.JSON(http.StatusConflict, AppointmentsResponse{Status: UnauthorizedResponse, Message: errMessage})
	}
	return err

}

func main() {
	// Create a db named "car-booking.db" in current directory.
	// It will be created if doesn't exsit.
	// And keep it connected.
	db, err := bolt.Open("car-booking.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
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
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "time=${time_custom}, status=${status}, method=${method}, uri=${uri}\nerror:{${error}}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(":1323"))
}
