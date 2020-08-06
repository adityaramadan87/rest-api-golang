package Controllers

import (
	"belajar-golang/app/Constant"
	"belajar-golang/app/Helper"
	"belajar-golang/app/Model"
	"belajar-golang/database"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type ReferalController struct{}

func (ReferalController) GetAllReferalCode(w http.ResponseWriter, r *http.Request) {
	var referal Model.Referal
	var referals []Model.Referal
	var responseReferal Model.ResponseReferal

	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT id,murid_id,referal_code,used FROM new_user_access")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&referal.Id, &referal.MuridID, &referal.ReferalCode, &referal.Used); err != nil {
			log.Fatal(err)
		} else {
			referals = append(referals, referal)
		}
	}

	responseReferal.Status = Constant.SuccessRequest
	responseReferal.Message = "Success"
	responseReferal.Data = referals

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseReferal)
}

func (ReferalController) InsertReferalCode(w http.ResponseWriter, r *http.Request) {

	var referal Model.Referal
	var referals []Model.Referal
	var responseReferal Model.ResponseReferal

	var refCode string

	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	referal.Id, _ = strconv.Atoi(r.Form.Get("id"))
	referal.MuridID, _ = strconv.Atoi(r.Form.Get("murid_id"))
	//referal.ReferalCode = Helper.AppHelper{}.StringWithCharset(16)
	referal.Used = false

	for {
		refCode = Helper.AppHelper{}.StringWithCharset(16)

		var summaryReferal Model.Referal = Helper.AppHelper{}.QueryReferalCode(refCode, 0)

		var murid Model.Murid
		err := db.QueryRow(`SELECT id,fullname,email,phone,jurusan,class,sub_class FROM murid WHERE id = $1`, referal.MuridID).
			Scan(
				&murid.Id,
				&murid.Fullname,
				&murid.Email,
				&murid.Phone,
				&murid.Jurusan,
				&murid.Class,
				&murid.SubClass,
			)
		_ = err

		var sumarryRCmurid Model.Referal = Helper.AppHelper{}.QueryReferalCode("", murid.Id)

		fmt.Print("Referal Code : " + refCode)

		if (Model.Referal{}.MuridID) != sumarryRCmurid.Id {
			mId := strconv.Itoa(murid.Id)
			log.Print("murid id not null " + mId)

			responseReferal.Status = Constant.BadRequest
			responseReferal.Message = "murid id already exists"
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(responseReferal)
			return
		} else if (Model.Referal{}.ReferalCode) != summaryReferal.ReferalCode {
			log.Print("ref code not null " + summaryReferal.ReferalCode)
			continue
		} else if (Model.Murid{}) == murid {
			responseReferal.Status = Constant.BadRequest
			responseReferal.Message = "murid not exists in database"
			w.Header().Set("Content-type", "application/json")
			json.NewEncoder(w).Encode(responseReferal)
			return
		} else {
			break
		}
	}

	_, err = db.Exec("INSERT INTO new_user_access (referal_code, murid_id, used) values ($1,$2,$3)", refCode, referal.MuridID, referal.Used)
	if err != nil {
		responseReferal.Status = Constant.BadRequest
		responseReferal.Message = "Insert Failed  " + err.Error()
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseReferal)
		return
	}

	data, err := db.Query("SELECT id,referal_code,murid_id,used FROM new_user_access WHERE referal_code = $1 AND murid_id = $2", referal.ReferalCode, referal.MuridID)
	for data.Next() {
		if err := data.Scan(&referal.Id, &referal.ReferalCode, &referal.MuridID, &referal.Used); err != nil {
			log.Print(err)
		} else {
			referals = append(referals, referal)
		}
	}

	responseReferal.Status = Constant.SuccessRequest
	responseReferal.Message = "Success Insert Referal Code"
	responseReferal.Data = referals
	log.Print("Insert to database table new_user_access")

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(responseReferal)

}

func (ReferalController) GetByReferalCode(w http.ResponseWriter, r *http.Request) {
	var referal Model.Referal
	var referals []Model.Referal
	var responseReferal Model.ResponseReferall

	fromMuxUrl := mux.Vars(r)
	log.Print(fromMuxUrl["refcode"])

	referalCode := fromMuxUrl["refcode"]

	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT id,murid_id,referal_code,used FROM new_user_access WHERE referal_code = $1", referalCode)
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&referal.Id, &referal.MuridID, &referal.ReferalCode, &referal.Used); err != nil {
			log.Fatal(err)
		} else {
			if (Model.Referal{}) != referal {
				_, err = db.Query("UPDATE new_user_access set used = $1 where id = $2", true, referal.Id)
			}

			referals = append(referals, referal)
		}
	}

	if (Model.Referal{}) != referal {
		if referal.Used {
			responseReferal.Status = Constant.BadRequest
			responseReferal.Message = "Referal code already used"
			responseReferal.Data = referal
		} else {
			responseReferal.Status = Constant.SuccessRequest
			responseReferal.Message = "Success"
			responseReferal.Data = referal
		}

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseReferal)

	} else {
		responseReferal.Status = Constant.BadRequest
		responseReferal.Message = "No data"
		responseReferal.Data = referal

		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(responseReferal)
	}
}
