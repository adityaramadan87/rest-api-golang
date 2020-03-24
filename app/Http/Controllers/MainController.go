package Controllers

import (
	"log"
	"net/http"
	"encoding/json"
	"belajar-golang/database"
	"belajar-golang/app/Model"
	Constant "belajar-golang/app/Constant"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	. "belajar-golang/app/Helper"
	_ "gopkg.in/gomail.v2"
	_ "encoding/base64"
	_ "strings"
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
		if err := rows.Scan(&users.Id, &users.Fullname, &users.Email, &users.Phone, &users.Avatar, &users.Password, &users.IsActivate); err != nil {
			log.Fatal(err)
		} else {
			dataUsers = append(dataUsers, users)
		}
	}

	responseUser.Status = Constant.SuccessRequest
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

func (MainController) Register(w http.ResponseWriter, r *http.Request){
	var users Model.User
	var dataUsers []Model.User
	var responseUser Model.Response
	var idUser int


	randInt := AppHelper{}.GenerateRandomInt()

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	users.Id, err = strconv.Atoi(r.Form.Get("id"))
	users.Fullname = r.Form.Get("fullname")
	users.Email = r.Form.Get("email")
	users.Phone = r.Form.Get("phone")
	users.Avatar = r.Form.Get("avatar")
	users.Password = r.Form.Get("password")
	users.IsActivate = false
	log.Print(users.Password)

	var summaryUser Model.User = AppHelper{}.QueryUser(users.Email)
	if (Model.User{}) != summaryUser {
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "the user who uses the email already exists"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return 
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)

	_, err = db.Exec("INSERT INTO users (fullname, email, phone, avatar, password, is_activate) values ($1,$2,$3,$4,$5,$6)", users.Fullname, users.Email, users.Phone, users.Avatar, hashPassword, users.IsActivate)
	if err != nil {
		log.Print(err)
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "Register Failed"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return 
	}
	
	data, err := db.Query("SELECT * FROM users WHERE email = $1", users.Email)
	for data.Next() {
		if err:= data.Scan(&users.Id, &users.Fullname, &users.Email, &users.Phone, &users.Avatar, &users.Password, &users.IsActivate); err != nil {
			log.Print(err)
		}else {
			idUser = users.Id
			dataUsers = append(dataUsers, users)
		}
	}
	var userAdded Model.User = AppHelper{}.QueryUser(users.Email)
	if (Model.User{}) != userAdded {
		idUser = userAdded.Id
		log.Print(idUser)
	}

	isSend := AppHelper{}.SendEmail(users.Email, idUser, randInt)

	if !isSend {
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "failed sending email"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return 
	}

	responseUser.Status = Constant.SuccessRequest
	responseUser.Message = "Success register \n We send email verification to your email \n please check your email for activation"
	responseUser.Data = dataUsers
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
	avatar := r.Form.Get("avatar")
	user_id, _ := strconv.Atoi(id)

	_, err = db.Query("UPDATE users set fullname = $1, email = $2, phone = $3, avatar = $4 where id = $5", fullname, email, phone, avatar, user_id)
	
	if err != nil {
		log.Print(err)
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "Failed Update Data"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return
	}

	responseUser.Status = Constant.SuccessRequest
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

	_, err = db.Exec("DELETE FROM users WHERE id = $1",user_id)
	
	if err != nil {
		log.Print(err)
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "failed delete data"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return
	}

	responseUser.Status = Constant.SuccessRequest
	responseUser.Message = "Success delete data"
	log.Print("delete data to database in table users")

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseUser)

}

func (MainController) UploadAvatar(w http.ResponseWriter, r *http.Request){
	var responseUser Model.Response
	var users Model.User
	var dataUsers []Model.User

	db := database.Connect()
	defer db.Close()
	
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	avatar := r.Form.Get("avatar")
	id := r.Form.Get("id")

	_, err = db.Exec("UPDATE users set avatar = $1 where id = $2", avatar, id)
	if err != nil {
		log.Print(err)
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "failed add avatar data"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return
	}

	
	data, err := db.Query("SELECT * FROM users WHERE id = $1", id,)
	for data.Next() {
		if err:= data.Scan(&users.Id, &users.Fullname, &users.Email, &users.Phone, &users.Avatar, &users.IsActivate); err != nil {
			log.Print(err)
		}else {
			dataUsers = append(dataUsers, users)
		}
	}
	
	responseUser.Status = Constant.SuccessRequest
	responseUser.Message = "Success add avatar"
	responseUser.Data = dataUsers
	log.Print("Insert to database table users")

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseUser)

}

func (MainController) LoginUser(w http.ResponseWriter, r *http.Request){
	var responseUser Model.Response

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	var summaryUser Model.User = AppHelper{}.QueryUser(email)
	if (Model.User{}) == summaryUser {
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "Wrong email"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return 
	} else if (Model.User{}.IsActivate) == summaryUser.IsActivate {
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "Account isn't activated, check your email for activation"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return 
	}

	var errPassword = bcrypt.CompareHashAndPassword([]byte(summaryUser.Password), []byte(password))

	if errPassword != nil {
		log.Print(errPassword)
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "Wrong password"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return 
	}

	responseUser.Status = Constant.SuccessRequest
	responseUser.Message = "Success Login"

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseUser)

}

