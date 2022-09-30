package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type DogString string

// type CustomFunc func(echo.Context) error
var users map[string]string = map[string]string{
	"tony":   "tonytsai",
	"wilson": "1234",
}

func customHandler(s string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, s)
	}
}

func handle(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func handle2(c echo.Context) error {
	return c.String(http.StatusOK, "Hello Sec!")
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	pw, doesKeyExist := users[username]

	if pw != password || !doesKeyExist {
		return echo.ErrUnauthorized
	}

	return c.String(http.StatusOK, fmt.Sprintf("username:%s\npassword:%s\n", username, password))
}

func main() {
	fmt.Println("jghfghfdgfs")
	fmt.Println("jghfghgggggg")
	fmt.Printf("hello %s %d\n", "cat", 5)

	// var e *echo.Echo // = echo.New()
	e := echo.New()
	e.GET("/dog", handle)
	e.POST("/login", login)
	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))

}
