package Controllers

import (
	Constant "belajar-golang/app/Constant"
	. "belajar-golang/app/Helper"
	"belajar-golang/app/Model"
	"belajar-golang/database"
	_ "encoding/base64"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	_ "gopkg.in/gomail.v2"
	"log"
	"net/http"
	"strconv"
	_ "strings"
)

type MainController struct{}

func (MainController) GetUser(res http.ResponseWriter, req *http.Request) {

	var users Model.User
	var dataUsers []Model.User
	var responseUser Model.ResponseUser

	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&users.Id, &users.MuridID, &users.Password, &users.Avatar, &users.ReferalCode, &users.IsActivate); err != nil {
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

func (MainController) Test(res http.ResponseWriter, req *http.Request) {
	var message Hello
	message.Message = "Testing"
	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(message)
}

func (MainController) Register(w http.ResponseWriter, r *http.Request) {
	var users Model.User
	var dataUsers []Model.User
	var responseUser Model.ResponseUser
	var idUser int

	//randInt := AppHelper{}.GenerateRandomInt()

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	users.Id, err = strconv.Atoi(r.Form.Get("id"))
	users.MuridID, err = strconv.Atoi(r.Form.Get("murid_id"))
	users.ReferalCode = r.Form.Get("referal_code")
	users.Avatar = r.Form.Get("avatar")
	users.Password = r.Form.Get("password")
	users.IsActivate = false
	log.Print(users.Password)

	var summaryUser Model.User = AppHelper{}.QueryUser(users.MuridID)
	if (Model.User{}) != summaryUser {
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "the user who uses the email already exists"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(users.Password), bcrypt.DefaultCost)

	_, err = db.Exec("INSERT INTO users (murid_id, avatar, password, is_activate, referal_code) values ($1,$2,$3,$4,$5)", users.MuridID, users.Avatar, hashPassword, users.IsActivate, users.ReferalCode)
	if err != nil {
		log.Print(err)
		responseUser.Status = Constant.BadRequest
		responseUser.Message = "Register Failed"
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseUser)
		return
	}

	data, err := db.Query("SELECT * FROM users WHERE murid_id = $1", users.MuridID)
	for data.Next() {
		if err := data.Scan(&users.Id, &users.MuridID, &users.ReferalCode, &users.Avatar, &users.Password, &users.IsActivate); err != nil {
			log.Print(err)
		} else {
			idUser = users.Id
			dataUsers = append(dataUsers, users)
		}
	}
	var userAdded Model.User = AppHelper{}.QueryUser(users.MuridID)
	if (Model.User{}) != userAdded {
		idUser = userAdded.Id
		log.Print(idUser)
	}

	//isSend := AppHelper{}.SendEmail(users.MuridID, idUser, randInt)

	//if !isSend {
	//	responseUser.Status = Constant.BadRequest
	//	responseUser.Message = "failed sending email"
	//	w.Header().Set("Content-type", "application/json")
	//	json.NewEncoder(w).Encode(responseUser)
	//	return
	//}

	responseUser.Status = Constant.SuccessRequest
	responseUser.Message = "Success register \n Login with your nim "
	responseUser.Data = dataUsers
	log.Print("Insert to database table users")

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseUser)

}

func (MainController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var responseUser Model.ResponseUser

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	id := r.Form.Get("id")
	muridID, _ := strconv.Atoi(r.Form.Get("murid_id"))
	referalCode := r.Form.Get("referalCode")
	avatar := r.Form.Get("avatar")
	user_id, _ := strconv.Atoi(id)

	_, err = db.Query("UPDATE users set murid_id = $1, referal_code = $2,  avatar = $4 where id = $5", muridID, referalCode, avatar, user_id)

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

func (MainController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var responseUser Model.ResponseUser
	//var id int

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	id := r.Form.Get("id")

	user_id, _ := strconv.Atoi(id)

	_, err = db.Exec("DELETE FROM users WHERE id = $1", user_id)

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

func (MainController) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	var responseUser Model.ResponseUser
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

	data, err := db.Query("SELECT * FROM users WHERE id = $1", id)
	for data.Next() {
		if err := data.Scan(&users.Id, &users.MuridID, &users.ReferalCode, &users.Avatar, &users.IsActivate); err != nil {
			log.Print(err)
		} else {
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

func (MainController) LoginUser(w http.ResponseWriter, r *http.Request) {
	var responseUser Model.ResponseUser

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	muridID, _ := strconv.Atoi(r.Form.Get("murid_id"))
	password := r.Form.Get("password")

	var summaryUser Model.User = AppHelper{}.QueryUser(muridID)
	//if (Model.User{}) == summaryUser {
	//	responseUser.Status = Constant.BadRequest
	//	responseUser.Message = "Wrong email"
	//	w.Header().Set("Content-type", "application/json")
	//	json.NewEncoder(w).Encode(responseUser)
	//	return
	//} else if (Model.User{}.IsActivate) == summaryUser.IsActivate {
	//	responseUser.Status = Constant.BadRequest
	//	responseUser.Message = "Account isn't activated, check your email for activation"
	//	w.Header().Set("Content-type", "application/json")
	//	json.NewEncoder(w).Encode(responseUser)
	//	return
	//}

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
