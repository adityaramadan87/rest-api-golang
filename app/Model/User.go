package Model

type User struct {
	Id          int    `form:"id" json:"id"`
	MuridID     int    `form:"murid_id" json:"murid_id"`
	Password    string `form:"password" json:"password"`
	Avatar      string `form:"avatar" json:"avatar"`
	IsActivate  bool   `json:"is_activate"`
	ReferalCode string `form:"referal_code" json:"referal_code"`
}

type ResponseUser struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User
}
