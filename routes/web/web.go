package web

import (
	"github.com/gorilla/mux"
	. "belajar-golang/app/Http/Middleware"
	. "belajar-golang/app/Http/Controllers"
)

func SetRoutes() *mux.Router {
	
	router := mux.NewRouter()

	//input routes in here
	router.HandleFunc("/user/get", MainController{}.GetUser).Methods("GET")
	router.HandleFunc("/test", MainController{}.Test).Methods("GET")
	router.HandleFunc("/user/add", MainController{}.InsertUser).Methods("POST")
	router.HandleFunc("/user/update", MainController{}.UpdateUser).Methods("PUT")
	router.HandleFunc("/user/delete", MainController{}.DeleteUser).Methods("POST")
	//if use middleware
	router.HandleFunc("/auth/user/get", ApiMiddleware{}.Auth(MainController{}.GetUser)).Methods("GET")

	return router

}