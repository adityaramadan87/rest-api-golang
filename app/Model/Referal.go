package Model

type Referal struct {
	Id          int    `form:"id" json:"id"`
	ReferalCode string `form:"referal_code" json:"referal_code"`
	MuridID     int    `form:"murid_id" json:"murid_id"`
	Used        bool   `form:"used" json:"used"`
}

type ResponseReferal struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Data    []Referal `json:"data"`
}

type ResponseReferall struct {
	Status  int     `json:"status"`
	Message string  `json:"message"`
	Data    Referal `json:"data"`
}
