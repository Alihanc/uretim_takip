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

	sql_data ".\üretim-takip-uygulaması\\database"
	"github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", sql_data.HomePage)
	myRouter.HandleFunc("/allpiece", sql_data._piece) // /allpiece adresine POST isteği gönderilmeli

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}
