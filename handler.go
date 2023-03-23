package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	User User `json:"user"`
}

func createUser(c echo.Context) error {
	newUser := User{}
	if err := c.Bind(&newUser); err != nil {
		return err
	}
	if err := mysqlRepo.CreateUser(&newUser); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, CreateUserResponse{newUser})
}
