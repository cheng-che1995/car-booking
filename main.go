package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type LoginResponse struct {
	Token string `json:"keyyy"`
}

type AppointmentsResponse struct {
	Appointments []Appointment `json:"appointments"`
}

// var appointments []time.Time

var appoint2 []Appointment

type Appointment struct {
	Username string
	Date     time.Time
}

// type CustomFunc func(echo.Context) error
var users map[string]string = map[string]string{
	"tony":   "tonytsai",
	"wilson": "1234",
}

func handle(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	pw, ok := users[username]

	if pw != password || !ok {
		return echo.ErrUnauthorized
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

	return c.JSON(http.StatusOK, LoginResponse{t})
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

	for _, a := range appoint2 {
		if a.Date == t {
			return c.String(http.StatusOK, fmt.Sprintf("%s，此日期已被預訂，請您重新選擇其他日期！", username))
		}
	}

	appoint2 = append(appoint2, Appointment{username, t})
	return c.String(http.StatusOK, fmt.Sprintf("預約成功！%s，您的預約日期為： %s", username, t.Format("2006-01-02")))
}

func searchAppointments(c echo.Context) error {
	var fitlerByUsername []Appointment
	filterSelectUsername := c.FormValue("filterSelectUsername")
	if filterSelectUsername == "" {
		fitlerByUsername = appoint2
	} else {
		for _, a := range appoint2 {
			if a.Username == filterSelectUsername {
				fitlerByUsername = append(fitlerByUsername, a)
			}
		}
	}
	return c.JSON(http.StatusOK, AppointmentsResponse{Appointments: fitlerByUsername})
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

	var found bool

	for i := range appoint2 {
		if t == appoint2[i].Date {
			if appoint2[i].Username != username {
				return c.String(http.StatusOK, fmt.Sprintf("此%s日期不屬於%s您的預約！", t.Format("2006-01-02"), username))
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
		return c.String(http.StatusOK, fmt.Sprintf("取消成功！%s，您已將 %s預約取消！", username, t.Format("2006-01-02")))
	}
	return c.String(http.StatusOK, "查無此預約！")
}

func main() {
	// var e *echo.Echo // = echo.New()
	e := echo.New()
	e.POST("/login", login)
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
	b.GET("/dog", handle)
	b.POST("/appointments", createAppointments)
	b.GET("/appointments", searchAppointments)
	b.DELETE("/appointments", cancelAppointments)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))

}
