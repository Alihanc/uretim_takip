package main

/*import (
	"log"
	"net/http"

	//_ "github.com/Alihanc/uretim_takip"

	 sql_data ".\üretim-takip-uygulaması\database"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", sql_data.HomePage)
	myRouter.HandleFunc("/allpiece", sql_data._piece)

	log.Fatal(http.ListenAndServe(":8080", myRouter))

}*/

import (
	"log"
	"net/http"

	sql_data "github.com/Alihanc/uretim_takip/database"
	"github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", sql_data.HomePage)
	myRouter.HandleFunc("/allpiece", sql_data.Piece).Methods("GET") // /allpiece adresine POST isteği gönderilmeli
	myRouter.HandleFunc("/addpiece", sql_data.AddPiece).Methods("POST")
	myRouter.HandleFunc("/signup", sql_data.username_signup).Methods("POST")

	log.Println("✅ Server çalışıyor: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
