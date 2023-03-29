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

const (
	unauthorizedMessage = "使用者名稱或密碼錯誤！"
	successMessage      = "登入成功！"
)

func login(c echo.Context) error {
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

type CreateCarRequest struct {
	Car
}

type CreateCarResponse struct {
	Car Car `json:"car"`
}

func createCar(c echo.Context) error {
	request := CreateCarRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	newCar := Car{
		Plate:    request.Plate,
		Uuid:     request.Uuid,
		UserUuid: request.UserUuid}

	if err := mysqlRepo.CreateCar(&newCar); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, CreateCarResponse{Car: newCar})
}

type CreateAppointmentRequest struct {
	Appointment
}

type CreateAppontmentRespone struct {
	Appointment Appointment `json:"appointment"`
}

func createAppointment(c echo.Context) error {
	request := CreateAppointmentRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	newAppointment := Appointment{
		StartTime: request.StartTime,
		EndTime:   request.EndTime,
		Uuid:      request.Uuid,
		UserUuid:  request.UserUuid,
		CarUuid:   request.CarUuid,
	}
	if err := mysqlRepo.CreateAppointment(&newAppointment); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, CreateAppontmentRespone{Appointment: newAppointment})
}

type GetUserRequest struct {
	Uuid string `json:"uuid"`
}

type GetUserResponse struct {
	User User `json:"user"`
}

func getUser(c echo.Context) error {
	request := GetUserRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	user, err := mysqlRepo.GetUser(request.Uuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetUserResponse{User: *user})
}

type GetUsersRequest struct {
	User
}

type GetUsersResponse struct {
	Users []User `json:"users"`
}

func getUsers(c echo.Context) error {
	request := GetUsersRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	g := GetUsersFilter{
		Uuid:     request.Uuid,
		Username: request.Username,
	}
	users, err := mysqlRepo.GetUsers(&g)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetUsersResponse{Users: users})
}

type GetCarRequest struct {
	Uuid string `json:"uuid"`
}

type GetCarResponse struct {
	Car Car `json:"car"`
}

func getCar(c echo.Context) error {
	request := GetCarRequest{}
	if err := c.Bind(&request); err != nil {
		return nil
	}
	car, err := mysqlRepo.GetCar(request.Uuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetCarResponse{Car: *car})
}

type GetCarsRequest struct {
	Car
}

type GetCarsResponse struct {
	Cars []Car `json:"cars"`
}

func getCars(c echo.Context) error {
	request := GetCarsRequest{}
	if err := c.Bind(&request); err != nil {
		return nil
	}
	g := GetCarsFilter{
		Uuid:     request.Uuid,
		Plate:    request.Plate,
		UserUuid: request.UserUuid,
	}
	cars, err := mysqlRepo.GetCars(&g)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetCarsResponse{Cars: cars})
}

type GetAppointmentRequest struct {
	Uuid string `json:"uuid"`
}

type GetAppointmentResponse struct {
	Appointment Appointment `json:"appointment"`
}

func getAppointment(c echo.Context) error {
	request := GetAppointmentRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	appointent, err := mysqlRepo.GetAppointment(request.Uuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, GetAppointmentResponse{Appointment: *appointent})
}

type GetAppointmentsRequest struct {
	Fields []string `json:"fields"`
	Appointment
}

type GetAppointmentsResponse struct {
	Appointments []Appointment `json:"appointments`
}

func getAppointments(c echo.Context) error {
	fieldsList := map[int]string{
		1: "appointment_uuid",
		2: "user_uuid",
		3: "car_uuid",
		4: "start_time",
		5: "end_time",
	}
	fields := []string{"appointment_uuid", "user_uuid", "car_uuid", "start_time", "end_time"}
	request := GetAppointmentsRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	g := GetAppointmentsFilter{
		Uuid:      request.Uuid,
		UserUuid:  request.UserUuid,
		CarUuid:   request.CarUuid,
		StartTime: request.StartTime.Format("2006-01-02"),
		EndTime:   request.EndTime.Format("2006-01-02"),
	}
	appointments, err := mysqlRepo.GetAppointments(request.Fields, &g)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, appointments)
}
