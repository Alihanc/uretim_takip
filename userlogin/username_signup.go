package username_login

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
//var db *sql.DB

// init fonksiyonu veritabanına bağlanır ve tabloyu oluşturur
func init() {
	var err error
	dns := "root:kemal1938@tcp(127.0.0.1:3306)/uretim_takip1?parseTime=true"
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

}

type usernames struct {
	ID        int       `json:"ID"`
	Username  string    `json:"UserName"`
	Password  string    `json:"PassWord"`
	Mail      string    `json:"Mail"`
	CreatedAT time.Time `json:"CreatedAT"`
}

func username_signup(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Lütfen POST isteği gönderin", http.StatusMethodNotAllowed)
		return
	}

	var user usernames

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "geçersiz json formatı", http.StatusBadRequest)
		log.Println("json parse hatası", err)
	}

	//Kayıt tarihi ekleme
	//const time.DateTime untyped string = "2006-01-02 15:04:05"
	user.CreatedAT = time.Now()

	//MYSQL'e ekleme sorgusu
	query := `INSERT INTO username(Username,Password,Mail,CreatedAT) VALUES(?,?,?,?)`
	result, err := db.Exec(query, user.Username, user.Password, user.Mail, user.CreatedAT)

	if err != nil {
		http.Error(w, fmt.Sprintf("Veri eklenirken hata: %v", err), http.StatusInternalServerError)
		log.Println("❌ SQL Hatası:", err) // Terminale hatayı yazdır
		return
	}

	//eklenen kaydın ID sini oluşturma
	insertedID, err := result.LastInsertId()
	if err != nil {
		log.Println("❌ SQL Hatası:", err) // Terminale hatayı yazdırs
		http.Error(w, "Kayıt ID alınırken hata", http.StatusInternalServerError)
		return

	}
	user.ID = int(insertedID)

	// JSON olarak yanıt döndür
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"massege": "Yeni kayıt başarıyla eklendi.",
		"data":    user,
	})

}

//
