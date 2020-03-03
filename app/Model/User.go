package Model

type User struct {
	Id int `form:"id" json:"id"`
	Fullname string `form:"fullname" json:"fullname"`
	Email string `form:"email" json:"email"`
	Phone string `form:"phone" json:"phone"`
}

type Response struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data []User
}