package web

import (
	. "belajar-golang/app/Helper"
	. "belajar-golang/app/Http/Controllers"
	. "belajar-golang/app/Http/Middleware"
	"github.com/gorilla/mux"
)

func SetRoutes() *mux.Router {

	router := mux.NewRouter()

	//input routes user in here
	router.HandleFunc("/user/get", MainController{}.GetUser).Methods("GET")
	router.HandleFunc("/test", MainController{}.Test).Methods("GET")
	router.HandleFunc("/user/register", MainController{}.Register).Methods("POST")
	router.HandleFunc("/user/update", MainController{}.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/delete", MainController{}.DeleteUser).Methods("POST")
	router.HandleFunc("/upload/avatar", MainController{}.UploadAvatar).Methods("POST")
	router.HandleFunc("/user/login", MainController{}.LoginUser).Methods("POST")

	//murid routes
	router.HandleFunc("/murid/get", MuridController{}.GetAllMurid).Methods("GET")

	//routes verification email
	router.HandleFunc("/user/verification/{id}", AppHelper{}.ActivateUser)

	//routes Attendance
	router.HandleFunc("/user/attendance", AttendanceController{}.Attendance).Methods("POST")
	router.HandleFunc("/user/check-attendance", AttendanceController{}.CheckAttendance).Methods("POST")

	//if use middleware
	router.HandleFunc("/auth/user/get", ApiMiddleware{}.Auth(MainController{}.GetUser)).Methods("GET")

	return router

}
