package Model

type Attendance struct{
	Id int `form:"id" json:"id"`
	UserID int `form:"user_id" json:"user_id"`
	AbsentIN string `form:"absent_in" json:"absent_in"`
	AbsentOUT string `form:"absent_out" json:"absent_out"`
	Type int `form:"type"`
}

type ResponseAttendance struct{
	Status int `json:"status"`
	Message string `json:"message"`
	Data Attendance
}