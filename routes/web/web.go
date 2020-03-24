package web

import (
	"github.com/gorilla/mux"
	. "belajar-golang/app/Http/Middleware"
	. "belajar-golang/app/Http/Controllers"
	. "belajar-golang/app/Helper"
)

func SetRoutes() *mux.Router {
	
	router := mux.NewRouter()

	//input routes in here
	router.HandleFunc("/user/get", MainController{}.GetUser).Methods("GET")
	router.HandleFunc("/test", MainController{}.Test).Methods("GET")
	router.HandleFunc("/user/register", MainController{}.Register).Methods("POST")
	router.HandleFunc("/user/update", MainController{}.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/delete", MainController{}.DeleteUser).Methods("POST")
	router.HandleFunc("/upload/avatar", MainController{}.UploadAvatar).Methods("POST")
	router.HandleFunc("/user/login", MainController{}.LoginUser).Methods("POST")

	router.HandleFunc("/user/verification/{id}", AppHelper{}.ActivateUser)

	//if use middleware
	router.HandleFunc("/auth/user/get", ApiMiddleware{}.Auth(MainController{}.GetUser)).Methods("GET")

	return router

}