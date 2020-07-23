package Controllers

import (
	"log"
	"net/http"
	"encoding/json"
	"belajar-golang/database"
	"belajar-golang/app/Model"
	Constant "belajar-golang/app/Constant"
	"strconv"
	"time"
)

type AttendanceController struct{}

func (AttendanceController) Attendance(w http.ResponseWriter, r *http.Request) {
	var attendance Model.Attendance
	var responseAttendance Model.ResponseAttendance

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		return
	}

	dt := time.Now()
	formatted := dt.Format("02/01/2006 15:04:05")

	ID, _ := strconv.Atoi(r.Form.Get("id"))
	UserID, _ := strconv.Atoi(r.Form.Get("user_id"))
	AbsentIN := formatted
	AbsentOUT := formatted
	Type, _ := strconv.Atoi(r.Form.Get("type"))

	//type have 3 values, if 1 is attendance in, if 2 is attendance out, if 3 is get all attendance user in that day

	if Type == 1 {
		_, err = db.Exec("INSERT INTO attendance (user_id, absent_in, absent_out) values ($1,$2,$3)", UserID, AbsentIN, "")
		if err != nil {
			log.Print(err)
			responseAttendance.Status = Constant.BadRequest
			responseAttendance.Message = "Attendance IN Failed"
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(responseAttendance)
		} else {

			data, err := db.Query("SELECT * FROM attendance WHERE user_id = $1 AND absent_in = $2", UserID, AbsentIN)
			if err != nil {
				log.Print(err)
				return
			}
			for data.Next() {
				if err:= data.Scan(&attendance.Id, &attendance.UserID, &attendance.AbsentIN, &attendance.AbsentOUT); err != nil {
					log.Print(err)
				}else {
					attendance.Type = 1
					responseAttendance.Status = Constant.SuccessRequest
					responseAttendance.Message = "Attendance IN Success"
					attendance.Type = Type
					responseAttendance.Data = attendance
					w.Header().Set("Content-type", "application/json")
					json.NewEncoder(w).Encode(responseAttendance)
				}
			}
		}
	} else if Type == 2 {
		_, err = db.Query("UPDATE attendance set absent_out = $1 WHERE id = $2", AbsentOUT, ID)
		if err != nil {
			responseAttendance.Status = Constant.BadRequest
			responseAttendance.Message = "Attendance OUT Failed"
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(responseAttendance)
			return
		}

		data, err := db.Query("SELECT * FROM attendance WHERE user_id = $1 AND absent_out = $2", UserID, AbsentOUT)
		if err != nil {
			log.Print(err)
			return
		}
		for data.Next() {
			if err:= data.Scan(&attendance.Id, &attendance.UserID, &attendance.AbsentIN, &attendance.AbsentOUT); err != nil {
				log.Print(err)
			}else {
				
				responseAttendance.Status = Constant.SuccessRequest
				responseAttendance.Message = "Attendance OUT Success"
				attendance.Type = Type
				responseAttendance.Data = attendance
				w.Header().Set("Content-type", "application/json")
				json.NewEncoder(w).Encode(responseAttendance)
			}
		}

	} else {

		data, err := db.Query("SELECT * FROM attendance WHERE id = $1", ID)
		if err != nil {
			log.Print(err)
			return
		}
		for data.Next() {
			if err:= data.Scan(&attendance.Id, &attendance.UserID, &attendance.AbsentIN, &attendance.AbsentOUT); err != nil {
				log.Print(err)
			}else {
				
				responseAttendance.Status = Constant.SuccessRequest
				responseAttendance.Message = "Attendance Clear"
				attendance.Type = Type
				responseAttendance.Data = attendance
				w.Header().Set("Content-type", "application/json")
				json.NewEncoder(w).Encode(responseAttendance)
			}
		}

	}

}

func (AttendanceController) CheckAttendance(w http.ResponseWriter, r *http.Request) {
	var attendance Model.Attendance
	var responseAttendance Model.ResponseAttendance

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Print(err)
		return
	}

	dt := time.Now()
	formatted := dt.Format("02/01/2006")

	// ID, _ := strconv.Atoi(r.Form.Get("id"))
	UserID, _ := strconv.Atoi(r.Form.Get("user_id"))
	AbsentIN := formatted + "%"
	// AbsentOUT := formatted

	data, err := db.Query("SELECT * FROM attendance WHERE user_id = $1 AND absent_in LIKE $2", UserID, AbsentIN)
		if err != nil {
			log.Print(err)
			return
		}
		if data != nil {
			for data.Next() {
				if err:= data.Scan(&attendance.Id, &attendance.UserID, &attendance.AbsentIN, &attendance.AbsentOUT); err != nil {
					log.Print(err)
				}else {
					
					responseAttendance.Status = Constant.SuccessRequest
					responseAttendance.Message = "Check Attendance"
					responseAttendance.Data = attendance
					w.Header().Set("Content-type", "application/json")
					json.NewEncoder(w).Encode(responseAttendance)
				}
			}
		}else {
			var att Model.Attendance
			responseAttendance.Status = Constant.SuccessRequest
			responseAttendance.Message = "Check Attendance"
			responseAttendance.Data = att
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(responseAttendance)
		}
		

}