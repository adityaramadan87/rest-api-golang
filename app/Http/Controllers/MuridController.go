package Controllers

import (
	"belajar-golang/app/Constant"
	"belajar-golang/app/Helper"
	"belajar-golang/app/Model"
	"belajar-golang/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type MuridController struct{}

func (MuridController) GetAllMurid(res http.ResponseWriter, req *http.Request) {

	var murid Model.Murid
	var murids []Model.Murid
	var responseMurid Model.ResponseMurid

	db := database.Connect()
	defer db.Close()
	rows, err := db.Query("SELECT * FROM murid")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&murid.Id, &murid.Fullname, &murid.Email, &murid.Phone, &murid.Jurusan, &murid.Class, &murid.SubClass); err != nil {
			log.Fatal(err)
		} else {
			murids = append(murids, murid)
		}
	}

	responseMurid.Status = Constant.SuccessRequest
	responseMurid.Message = "Success"
	responseMurid.Data = murids

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(responseMurid)

}

func (MuridController) InsertMurid(res http.ResponseWriter, req *http.Request) {
	var murid Model.Murid
	var murids []Model.Murid
	var responseMurid Model.ResponseMurid

	db := database.Connect()
	defer db.Close()

	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	murid.Id, _ = strconv.Atoi(req.Form.Get("id"))
	murid.Fullname = req.Form.Get("fullname")
	murid.Email = req.Form.Get("email")
	murid.Phone = req.Form.Get("phone")
	murid.Jurusan = req.Form.Get("jurusan")
	murid.Class, _ = strconv.Atoi(req.Form.Get("class"))
	murid.SubClass, _ = strconv.Atoi(req.Form.Get("sub_class"))

	var summaryMurid Model.Murid = Helper.AppHelper{}.QueryMurid(murid.Email)

	if (Model.Murid{}) != summaryMurid {
		responseMurid.Status = Constant.BadRequest
		responseMurid.Message = "the Murid already exists"
		res.Header().Set("Content-type", "application/json")
		json.NewEncoder(res).Encode(responseMurid)
		return
	}

	_, err = db.Exec("INSERT INTO murid (fullname, email, phone, jurusan, class, sub_class) values ($1,$2,$3,$4,$5,$6)", murid.Fullname, murid.Email, murid.Phone, murid.Jurusan, murid.Class, murid.SubClass)

	if err != nil {
		responseMurid.Status = Constant.BadRequest
		responseMurid.Message = "Insert Failed  " + err.Error()
		res.Header().Set("Content-type", "application/json")
		json.NewEncoder(res).Encode(responseMurid)
		return
	}

	data, err := db.Query("SELECT * FROM murid WHERE email = $1", murid.Email)
	for data.Next() {
		if err := data.Scan(&murid.Id, &murid.Fullname, &murid.Email, &murid.Phone, &murid.Jurusan, &murid.Class, &murid.SubClass); err != nil {
			log.Print(err)
		} else {
			murids = append(murids, murid)
		}
	}

	responseMurid.Status = Constant.SuccessRequest
	responseMurid.Message = "Success Insert murid"
	responseMurid.Data = murids
	log.Print("Insert to database table murid")

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(responseMurid)

}
