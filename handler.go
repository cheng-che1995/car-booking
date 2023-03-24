package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Token   string `json:"token"`
	Status  string `json:"status"`
	Message string `json:"message"`
	User    User   `json:"user"`
}

func login(c echo.Context) error {
	unauthorizedMessage := "使用者名稱或密碼錯誤！"
	successMessage := "登入成功！"
	r := LoginRequest{}
	if err := c.Bind(&r); err != nil {
		return err
	}
	user := User{Username: r.Username, Password: r.Password}
	if ok, err := mysqlRepo.AuthUser(&user); err != nil {
		return err
	} else if !ok {
		return c.JSON(http.StatusUnauthorized, LoginResponse{Status: UnauthorizedResponse, Message: unauthorizedMessage})
	}
	var expireTime time.Time
	expireTime = time.Now().Add(time.Hour * 72)
	claims := &jwtCustomClaims{
		r.Username,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	resp := LoginResponse{
		Token:   t,
		Status:  SuccessResponse,
		Message: successMessage,
		User:    user}
	cookie := new(http.Cookie)
	cookie.Name = "jwt_access"
	cookie.Value = t
	cookie.Expires = expireTime
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Secure = true
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, resp)
}
