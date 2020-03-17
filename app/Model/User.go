package Model

type User struct {
	Id int `form:"id" json:"id"`
	Fullname string `form:"fullname" json:"fullname"`
	Email string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Phone string `form:"phone" json:"phone"`
	Avatar string `form:"avatar" json:"avatar"`
	IsActivate bool `json:"is_activate"`
}

type Response struct {
	Status int `json:"status"`
	Message string `json:"message"`
	Data []User
}