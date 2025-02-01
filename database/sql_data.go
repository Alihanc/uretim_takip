package sql_data

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

func sql_open() {

	dns := "root:kemal1938@tcp(127.0.0.1:3306)/uretim_takip1?Time=true"

	db, err := sql.Open("mysql", dns)

	if err != nil {
		log.Fatal("veri tabanına bağlanırken sorun oluştu: %v", err)
	}
	defer db.Close()
}

func control(error error) {
	if error != nil {
		log.Fatal(error)
	}
}
func new_piece() {
	parca_adi := fmt.Scanf("parça numarasını giriniz", &parca_adi)
	parca_malzemesi := fmt.Scanf("parça malzemesini giriniz", &parca_malzemesi)
	uretim_sekli := fmt.Scanf("üretim yöntemini giriniz", &uretim_sekli)
	uretim_adedi := fmt.Scanf("uretim adedini giriniz", &uretim_adedi)


	uretim-takip,err := db.Exec('INSERT INTO uretim_takip (parca_adi, parca_malzemesi, uretim_sekli, uretim_adedi) VALUES (?,?,?,?)',
	parca_adi,parca_malzemesi, uretim_sekli, uretim_adedi)


control(err)
eklenenparca, err :=uretim-takip.LastInsertId()
fmt.Println(eklenenparca)

}

func delete_piece(){
	silinecek_parca:= 1
	_,err := vt.Exec('DELETE FROM uretim_takip WHERE id =?',silinecek_parca)
	control(err)
}