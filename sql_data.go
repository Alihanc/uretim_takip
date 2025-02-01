package sql_data

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func sql_open() {

	dns := "root:kemal1938@tcp(127.0.0.1:3306)/uretim-takip"

	db, err := sql.Open("mysql", dns)

	if err != nil {
		log.Fatal("veri tabanına bağlanırken sorun oluştu: %v", err)
	}
	defer db.Close()
}

func add_piece() {
	var id int
	var parca_adi string
	var parca_malzemesi string
	var uretim_sekli string
	var uretim_adedi int
	kayit_tarihi := time.Now()
	sonuc, err := db.Exec(`INSERT INTO usernames ( username , password, createdAt) VALUES (?, ?, ?)`, username, password, createdAt)
	kontrol(err)
	eklenenid, err := sonuc.LastInsertId()
	fmt.Println(eklenenid)
}
