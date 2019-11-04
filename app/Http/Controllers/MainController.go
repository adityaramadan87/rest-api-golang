package Controllers

import (
	"log"
	"net/http"
	"encoding/json"
	"belajar-golang/database"
	"belajar-golang/app/Model"
	"strconv"
)

type MainController struct{}

func (MainController) GetUser(res http.ResponseWriter, req *http.Request)  {

	var users Model.User
	var dataUsers []Model.User
	var responseUser Model.Response

	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Print(err)
	}

	for rows.Next(){
		if err := rows.Scan(&users.Id, &users.Fullname, &users.Email, &users.Phone); err != nil {
			log.Fatal(err)
		} else {
			dataUsers = append(dataUsers, users)
		}
	}

	responseUser.Status = 200
	responseUser.Message = "OKE!"
	responseUser.Data = dataUsers

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(responseUser)
	
}

type Hello struct {
	Message string
}

func (MainController) Test(res http.ResponseWriter, req *http.Request){
	var message Hello
	message.Message = "Testing"
	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(message)
}

func (MainController) InsertUser(w http.ResponseWriter, r *http.Request){
//	var users Model.User
//	var dataUsers []Model.User
	var responseUser Model.Response

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	id := r.Form.Get("id")
	fullname := r.Form.Get("fullname")
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")

	_, err = db.Exec("INSERT INTO users (id, fullname, email, phone) values ($1,$2,$3,$4)",id, fullname, email, phone,)
	if err != nil {
		log.Print(err)
	}
	responseUser.Status = 200
	responseUser.Message = "Add 1 Data"
	log.Print("Insert to database table users")

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseUser)

}
func (MainController) UpdateUser(w http.ResponseWriter, r *http.Request){

	var responseUser Model.Response


	db := database.Connect()
	defer db.Close()
	
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	id := r.Form.Get("id")
	fullname := r.Form.Get("fullname")
	email := r.Form.Get("email")
	phone := r.Form.Get("phone")

	_, err = db.Exec("UPDATE users set fullname = $1, email = $2, phone = $3 where id = $4",fullname, email, phone, id,)
	if err != nil {
		log.Print(err)
	}

	responseUser.Status = 200
	responseUser.Message = "Success Update Data"
	log.Print("Update data to database table users")

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseUser)

}

func (MainController) DeleteUser(w http.ResponseWriter, r *http.Request){
	var responseUser Model.Response
	//var id int

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id := r.Form.Get("id")

	user_id, _ := strconv.Atoi(id)

	_, err = db.Exec("DELETE from users where id=$1",user_id)
	
	if err != nil {
		log.Print(err)
	}

	responseUser.Status = 200
	responseUser.Message = "Success delete data"
	log.Print("deleta data to database in table users")

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseUser)

}



