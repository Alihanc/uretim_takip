package main

import (
	"log"
	"net/http"

	sql_data "github.com/Alihanc/uretim_takip/database"
	username "github.com/Alihanc/uretim_takip/userlogin"
	"github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", sql_data.HomePage)
	myRouter.HandleFunc("/allpiece", sql_data.Piece).Methods("GET") // /allpiece adresine POST isteği gönderilmeli
	myRouter.HandleFunc("/addpiece", sql_data.AddPiece).Methods("POST")
	myRouter.HandleFunc("/login", username.Login).Methods("POST")
	myRouter.HandleFunc("/signup", username.Username_signup).Methods("POST")
	myRouter.HandleFunc("/signdelete", username.Username_delete).Methods("DELETE")
	myRouter.HandleFunc("/updatepiece", sql_data.UpdatePiece).Methods("PUT")

	log.Println("✅ Server çalışıyor: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
