package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	/* Bolt database
	dbBolt, err := bolt.Open("car-booking.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	dbBolt.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Appointments"))
		if err != nil {
			return fmt.Errorf("create bucket err: %s", err)
		}
		return nil
	})
	dbBolt.Close()
	*/

	//Mysql database
	mysqlRepo.OpenConn()
	defer mysqlRepo.CloseConn()
	e := echo.New()
	e.POST("/login", login)
	// e.GET("/users", showUsers)
	e.POST("/users", createUser)
	b := e.Group("/booking")
	b.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &jwtCustomClaims{},
		SigningKey:  []byte("secret"),
		ContextKey:  "token",
		TokenLookup: "cookie:jwt_access",
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			return c.String(http.StatusUnauthorized, err.Error())
		},
	}))
	b.DELETE("/users/:uuid", deleteUser)
	b.GET("/users/:uuid", getUser)
	b.GET("/users", getUsers)
	b.POST("/cars", createCar)
	b.DELETE("/cars/:uuid", deleteCar)
	b.GET("/cars/:uuid", getCar)
	b.GET("/cars", getCars)
	b.POST("/appointments", createAppointment)
	b.DELETE("/appointments/:uuid", deleteAppointment)
	b.GET("/appointments/:uuid", getAppointment)
	b.GET("/appointments", getAppointments)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "time=${time_custom}, status=${status}, method=${method}, uri=${uri}\nerror:{${error}}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))
	//TODO: 暫且先用 "*" 方便測試。
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))
	e.Logger.Fatal(e.Start(":1323"))
}
