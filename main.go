package main

import (
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

func main() {
	// var e *echo.Echo // = echo.New()
	e := echo.New()
	e.POST("/login", login)

	b := e.Group("/booking")
	b.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
		ErrorHandlerWithContext: func(err error, c echo.Context) error {
			return c.String(http.StatusUnauthorized, "JWT驗證無效")
		},
	}))

	e.Use(middleware.Logger())
	e.Logger.Fatal(e.Start(":1323"))

}
