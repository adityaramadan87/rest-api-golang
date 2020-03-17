package Helper

import (
	"belajar-golang/database"
	"belajar-golang/app/Model"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"time"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
)

type AppHelper struct {}

func (AppHelper) SendEmail(email string, id int, randInt int) bool {
	const CONFIG_SMTP_HOST = "smtp.gmail.com"
	const CONFIG_SMTP_PORT = 587
	const CONFIG_EMAIL = "youremail"
	const CONFIG_PASSWORD = "yourpass"

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_EMAIL)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Email Verification")
	mailer.SetBody("text/html", "<a href=\"http://localhost:9000/user/verification/"+strconv.Itoa(AppHelper{}.GenerateRandomInt())+strconv.Itoa(id)+strconv.Itoa(randInt)+"\"><button type=\"submit\">ACITVATE</button></a>")

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_EMAIL,
		CONFIG_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err)
		return false
	}

	log.Print("Success send mail")
	return true
}

func (AppHelper) QueryUser(email string) Model.User {
	var users Model.User
	db := database.Connect()
	defer db.Close()
	err := db.QueryRow(
		`SELECT id,
		fullname,
		email,
		phone,
		avatar,
		password,
		is_activate
		FROM users WHERE email = $1`,
		email).
		Scan(
			&users.Id,
			&users.Fullname,
			&users.Email,
			&users.Phone,
			&users.Avatar,
			&users.Password,
			&users.IsActivate,
		)
	_ = err
	return users
}

func (AppHelper) ActivateUser(w http.ResponseWriter, r *http.Request) {
	fromMuxUrl := mux.Vars(r)
	log.Print(fromMuxUrl["id"])

	//add query in here and update the is_activated to true

	fmt.Fprintln(w, "Halo !! \n  Apa Kabar!!")
}

func (AppHelper) GenerateRandomInt() int {
    rand.Seed(time.Now().UnixNano())
    slice := rand.Intn(99999)
    return slice
}