package main

import (
	"testing"
)

func testCreateUser(t *testing.T) {
	request := CreateUserRequest{
		Username: "Wilson",
		Password: "88888888",
	}
	response := CreateUserResponse{
		User: User{
			Username: "Wilson",
		},
	}
}
