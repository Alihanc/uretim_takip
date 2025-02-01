package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", uretim_takip.HomePage)
	myRouter.HandleFunc("/allpiece", database.All)

	log.Fatal(http.ListenAndServe(":8080", myRouter))

}
