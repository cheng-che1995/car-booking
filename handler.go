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
	request := CreateUserRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}

	//
	newUser := User{Username: request.Username, Password: request.Password}
	if err := mysqlRepo.CreateUser(&newUser); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, CreateUserResponse{User: newUser})
}
