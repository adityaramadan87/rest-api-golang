package Helper

import (
	"belajar-golang/app/Model"
	"belajar-golang/database"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/gomail.v2"
	"log"
	"math/rand"
	"net/http"
	_ "reflect"
	"strconv"
	"time"
)

type AppHelper struct{}

var expiredActivationTimer *time.Timer
var timeExpired bool = false
var hasBeenVerified = false

func (AppHelper) SendEmail(email string, id int, randInt int) bool {
	const CONFIG_SMTP_HOST = "smtp.gmail.com"
	const CONFIG_SMTP_PORT = 587
	const CONFIG_EMAIL = "yourEmail"
	const CONFIG_PASSWORD = "yourPassword"

	hashId := strconv.Itoa(AppHelper{}.GenerateRandomInt()) + strconv.Itoa(id) + strconv.Itoa(randInt)

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_EMAIL)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Email Verification")
	mailer.SetBody("text/html", "<a href=\"http://localhost:9000/user/verification/"+hashId+"\"><button type=\"submit\">ACITVATE</button></a>")

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

	expiredActivationTimer = time.NewTimer(25 * time.Second)

	go func() {
		<-expiredActivationTimer.C
		log.Print("expired")
		timeExpired = true
	}()

	log.Print("Success send mail")
	return true
}

func (AppHelper) QueryUser(murid_id int) Model.User {
	var users Model.User

	db := database.Connect()
	defer db.Close()
	err := db.QueryRow(
		`SELECT id,
		murid_id,
		avatar,
		password,
		is_activate,
		referal_code
		FROM users WHERE murid_id = $1`,
		murid_id).
		Scan(
			&users.Id,
			&users.MuridID,
			&users.Avatar,
			&users.Password,
			&users.IsActivate,
			&users.ReferalCode,
		)
	_ = err
	return users
}

//Belom Berguna
func (AppHelper) ActivateUser(w http.ResponseWriter, r *http.Request) {
	log.Print(timeExpired)
	if timeExpired {
		if hasBeenVerified {
			fmt.Fprintln(w, "Your email has been verified,")
			return
		}

		fmt.Fprintln(w, "Email verification expired, try to register again")
		return
	}

	fromMuxUrl := mux.Vars(r)
	log.Print(fromMuxUrl["id"])

	id := fromMuxUrl["id"]
	isActive := true

	slice := []rune(id)

	resultSlice := string(slice[5:11])

	log.Print(resultSlice + " " + strconv.FormatBool(isActive))
	//add query in here and update the is_activated to true
	db := database.Connect()
	defer db.Close()

	_, err := db.Exec("UPDATE users set is_activate = $1 where id = $2", isActive, resultSlice)
	if err != nil {
		log.Print(err)
		fmt.Fprintln(w, "failed verification email")
	}

	timeExpired = true
	hasBeenVerified = true
	fmt.Fprintln(w, "Email successfully verified")
}

func (AppHelper) GenerateRandomInt() int {
	rand.Seed(time.Now().UnixNano())
	slice := rand.Intn(99999)
	return slice
}
