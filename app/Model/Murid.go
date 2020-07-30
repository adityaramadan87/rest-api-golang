package Model

type Murid struct {
	Id       int    `form:"id" json:"id"`
	Fullname string `form:"fullname" json:"fullname"`
	Email    string `form:"email" json:"email"`
	Phone    string `form:"phone" json:"phone"`
	Jurusan  string `form:"jurusan" json:"jurusan"`
	Class    int    `form:"class" json:"class"`
	SubClass int    `form:"sub_class" json:"sub_class"`
}

type ResponseMurid struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Murid
}
