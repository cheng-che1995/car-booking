package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type AppointmentsResponse struct {
	Appointments []Appointment `json:"appointments"`
	Status       string        `json:"status"`
	Message      string        `json:"message"`
}

// type NewAppointmentsResponse struct {
// 	NewAppointments []NewAppointment `json:"appointments"`
// 	Status          string           `json:"status"`
// 	Message         string           `json:"message"`
// }

const (
	SuccessResponse      string = "success"
	ConflictResponse     string = "conflict"
	NotFoundResponse     string = "notFound"
	UnauthorizedResponse string = "unauthorized"
)

// type CustomFunc func(echo.Context) error
var users map[string]string = map[string]string{
	"tony":   "4321",
	"wilson": "1234",
}

var mysqlRepo Repository

func showUsers(c echo.Context) error {
	var account []string
	for u := range users {
		account = append(account, u)
	}
	return c.JSON(http.StatusOK, account)
}

// TODO: use mysql database.

// func createAppointments(c echo.Context) error {
// 	token := c.Get("token").(*jwt.Token)
// 	claims := token.Claims.(*jwtCustomClaims)
// 	username := claims.Name
// 	selectedDate := c.FormValue("selectedDate")
// 	t, err := time.Parse("2006-01-02", selectedDate)
// 	if err != nil {
// 		return err
// 	}
// 	errMessage := fmt.Sprintf("%s，此日期已被預訂，請您重新選擇其他日期！", username)
// 	successMessage := fmt.Sprintf("預約成功！%s，您的預約日期為： %s", username, t.Format("2006-01-02"))
// 	//TODO: Put this elsewhere
// 	br := BoltRepository{dbPath: "car-booking.db"}
// 	selectAppointments := Appointment{
// 		Username: username,
// 		Date:     t,
// 	}
// 	if err := br.Create(&selectAppointments); err != nil {
// 		return c.JSON(http.StatusConflict, AppointmentsResponse{Status: ConflictResponse, Message: errMessage})
// 	}
// 	return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
// }

// func searchAppointments(c echo.Context) error {
// 	filterByUsername := c.FormValue("filterByUsername")
// 	filterByDateStart := c.FormValue("filterByDateStart")
// 	filterByDateEnd := c.FormValue("filterByDateEnd")

// 	br := BoltRepository{dbPath: "car-booking.db"}
// 	//TODO: 補上filterByItem
// 	selectedFilter := SearchFilter{
// 		Username:  &filterByUsername,
// 		DateStart: &filterByDateStart,
// 		DateEnd:   &filterByDateEnd,
// 	}
// 	FilteredAppointments, err := br.Search(&selectedFilter)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Appointments: FilteredAppointments})
// }

// func cancelAppointments(c echo.Context) error {
// 	token := c.Get("token").(*jwt.Token)
// 	claims := token.Claims.(*jwtCustomClaims)
// 	username := claims.Name
// 	selectedDate := c.FormValue("selectedDate")
// 	t, err := time.Parse("2006-01-02", selectedDate)
// 	if err != nil {
// 		return err
// 	}
// 	successMessage := fmt.Sprintf("取消成功！%s，您已將 %s預約取消！", username, t.Format("2006-01-02"))
// 	errMessage := fmt.Sprintf("此%s日期不屬於%s您的預約！", t.Format("2006-01-02"), username)
// 	notFoundMessage := fmt.Sprintf("查無此預約！%s請您重新選擇日期！", username)

// 	br := BoltRepository{dbPath: "car-booking.db"}
// 	selectAppointments := Appointment{
// 		Username: username,
// 		Date:     t,
// 	}

// 	if err = br.Delete(&selectAppointments); err == nil {
// 		return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
// 	} else if err == ErrNotFound {
// 		return c.JSON(http.StatusNotFound, AppointmentsResponse{Status: NotFoundResponse, Message: notFoundMessage})
// 	} else if err == ErrUnauthorized {
// 		return c.JSON(http.StatusConflict, AppointmentsResponse{Status: UnauthorizedResponse, Message: errMessage})
// 	}
// 	return err

// }

// TODO: 新增驗證
func deleteUser(c echo.Context) error {
	newUser := User{
		Uuid: c.Param("uuid"),
	}
	if err := mysqlRepo.DeleteUser(&newUser); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("使用者刪除成功！"))
}
func getUser(c echo.Context) error {
	uuid := c.Param("uuid")
	user, err := mysqlRepo.GetUser(uuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user)
}

func getUsers(c echo.Context) error {
	g := new(GetUsersFilter)
	if err := c.Bind(g); err != nil {
		return err
	}
	users, err := mysqlRepo.GetUsers(g)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, users)
}

func deleteCar(c echo.Context) error {
	car := new(Car)
	if err := c.Bind(car); err != nil {
		return err
	}
	if err := mysqlRepo.DeleteCar(car); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("車輛%s刪除成功！", car.Plate))
}

func getCar(c echo.Context) error {
	uuid := c.Param("uuid")
	car, err := mysqlRepo.GetCar(uuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, car)
}

func getCars(c echo.Context) error {
	g := new(GetCarsFilter)
	if err := c.Bind(g); err != nil {
		return err
	}
	cars, err := mysqlRepo.GetCars(g)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, cars)
}

func deleteAppointment(c echo.Context) error {
	appointment := Appointment{
		Uuid: c.FormValue("appointment_uuid"),
	}
	if err := mysqlRepo.DeleteAppointment(&appointment); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("預約刪除成功！"))
}

func getAppointment(c echo.Context) error {
	uuid := c.FormValue("appointment_uuid")
	appointent, err := mysqlRepo.GetAppointment(uuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, appointent)
}

func getAppointments(c echo.Context) error {
	fields := []string{"appointment_uuid", "user_uuid", "car_uuid", "start_time", "end_time"}
	g := GetAppointmentsFilter{
		Uuid:      c.FormValue("appointment_uuid"),
		UserUuid:  c.FormValue("user_uuid"),
		CarUuid:   c.FormValue("car_uuid"),
		StartTime: c.FormValue("start_time"),
		EndTime:   c.FormValue("end_time"),
	}
	appointments, err := mysqlRepo.GetAppointments(fields, &g)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, appointments)
}

// func createAppointmentsByMysql(c echo.Context) error {
// 	token := c.Get("token").(*jwt.Token)
// 	claims := token.Claims.(*jwtCustomClaims)
// 	username := claims.Name
// 	selectedDate := c.FormValue("selectedDate")
// 	selectedItem := c.FormValue("selectedItem")

// 	errMessage := fmt.Sprintf("%s，此日期已被預訂，請您重新選擇其他日期！", username)
// 	successMessage := fmt.Sprintf("預約成功！%s，您的預約日期為： %s", username, selectedDate)

// 	if err := mysqlRepo.Create(username, selectedItem, selectedDate); err != nil {
// 		return c.JSON(http.StatusConflict, AppointmentsResponse{Status: ConflictResponse, Message: errMessage})
// 	}
// 	return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
// }

// func searchAppointmentsByMysql(c echo.Context) error {
// 	filterByUsername := c.FormValue("filterByUsername")
// 	filterByItem := c.FormValue("filterByItem")
// 	filterByDateStart := c.FormValue("filterByDateStart")
// 	filterByDateEnd := c.FormValue("filterByDateEnd")

// 	selectedFilter := SearchFilter{
// 		Username:  &filterByUsername,
// 		Item:      &filterByItem,
// 		DateStart: &filterByDateStart,
// 		DateEnd:   &filterByDateEnd,
// 	}

// 	FilteredAppointments, err := mysqlRepo.Search(&selectedFilter)
// 	if err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusOK, NewAppointmentsResponse{Status: SuccessResponse, NewAppointments: FilteredAppointments})
// }

// func cancelAppointmentsByMysql(c echo.Context) error {
// 	token := c.Get("token").(*jwt.Token)
// 	claims := token.Claims.(*jwtCustomClaims)
// 	username := claims.Name
// 	selectedDate := c.FormValue("selectedDate")
// 	selectedItem := c.FormValue("selectedItem")
// 	// t, err := time.Parse("2006-01-02", selectedDate)
// 	// if err != nil {
// 	// 	return err
// 	// }
// 	successMessage := fmt.Sprintf("取消成功！%s，您已將 %s預約取消！", username, selectedDate)
// 	unauthorizedMessage := fmt.Sprintf("此%s日期不屬於%s您的預約！", selectedDate, username)
// 	notFoundMessage := fmt.Sprintf("查無此預約!%s請您重新選擇日期！", username)

// 	if err := mysqlRepo.Delete(username, selectedItem, selectedDate); err == nil {
// 		return c.JSON(http.StatusOK, AppointmentsResponse{Status: SuccessResponse, Message: successMessage})
// 	} else if err == ErrNotFound {
// 		return c.JSON(http.StatusNotFound, AppointmentsResponse{Status: NotFoundResponse, Message: notFoundMessage})
// 	} else if err == ErrUnauthorized {
// 		return c.JSON(http.StatusUnauthorized, AppointmentsResponse{Status: UnauthorizedResponse, Message: unauthorizedMessage})
// 	}
// 	return nil
// }

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
