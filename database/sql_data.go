package sql_data

/*import (
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


func All(){
	type pieces struct{
	parca_adi  varchar
	parca_malzemesi varchar
	uretim_sekli  varchar
	uretim_adedi  int
	}
	table, err := vt.Query('SELECT,id, parca_adi,parca_malzemesi,uretim_sekli,uretim_adedi')
	control(err)
	defer table.Close()
	var uretimtakip [] pieces
	for table.Next(){
		var P pieces
		err:= table.Scan(&P.parca_adi,&P.parca_malzemesi,&P.uretim_sekli,&P.uretim_adedi)
		control(err)
		pieces = append(urtimtakip, P)


	}
	err:= rows.Err()
	control(err)

}
func delete_piece(){
	silinecek_parca:= 1
	_,err := vt.Exec('DELETE FROM uretim_takip WHERE id =?',silinecek_parca)
	control(err)
}*/

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// db global bağlantı değişkeni
var db *sql.DB

// init fonksiyonu veritabanına bağlanır ve tabloyu oluşturur
func init() {
	var err error
	dns := "root:kemal1938@tcp(127.0.0.1:3306)/uretim-takip"
	db, err = sql.Open("mysql", dns)
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata oluştu: %v", err)
	}

	// Bağlantıyı kontrol et
	err = db.Ping()
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı kontrol edilemedi: %v", err)
	}
	fmt.Println("MySQL veritabanına başarıyla bağlanıldı.")

	// Tabloyu oluştur (eğer yoksa)
	createTableQuery := `CREATE TABLE IF NOT EXISTS pieces (
		id INT AUTO_INCREMENT PRIMARY KEY,
		parca_adi VARCHAR(255),
		parca_malzemesi VARCHAR(255),
		uretim_sekli VARCHAR(255),
		uretim_adedi INT,
		kayit_tarihi DATETIME
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Tablo oluşturulurken hata: %v", err)
	}
}

// Piece veritabanındaki parça kayıtlarını temsil eder
type Piece struct {
	ID             int       `json:"id"`
	ParcaAdi       string    `json:"parca_adi"`
	ParcaMalzemesi string    `json:"parca_malzemesi"`
	UretimSekli    string    `json:"uretim_sekli"`
	UretimAdedi    int       `json:"uretim_adedi"`
	KayitTarihi    time.Time `json:"kayit_tarihi"`
}

// HomePage basit bir karşılama mesajı döndürür
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hoşgeldiniz! Parça eklemek için /allpiece adresine POST isteği gönderin.")
}

// _piece fonksiyonu POST isteği ile gelen JSON verisini alıp MySQL'e ekler
func _piece(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Lütfen POST isteği gönderin."})
		return
	}

	// JSON verisini al
	var p Piece
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	// Kayıt tarihini güncel zamana ayarla
	p.KayitTarihi = time.Now()

	// Veriyi veritabanına ekle
	query := `INSERT INTO pieces (parca_adi, parca_malzemesi, uretim_sekli, uretim_adedi, kayit_tarihi) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, p.ParcaAdi, p.ParcaMalzemesi, p.UretimSekli, p.UretimAdedi, p.KayitTarihi)
	if err != nil {
		http.Error(w, fmt.Sprintf("Veri eklenirken hata: %v", err), http.StatusInternalServerError)
		return
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, fmt.Sprintf("Kayıt ID alınırken hata: %v", err), http.StatusInternalServerError)
		return
	}
	p.ID = int(insertedID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}
