package Controllers

import (
	"belajar-golang/app/Constant"
	"belajar-golang/app/Model"
	"belajar-golang/database"
	"encoding/json"
	"log"
	"net/http"
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
