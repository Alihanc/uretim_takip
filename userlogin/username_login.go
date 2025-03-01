package username_login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

// db global bağlantı değişkeni
var db *sql.DB

// init fonksiyonu veri tabanına bağlanır
func init() {
	var err error
	dns := "root:kemal1938@tcp(127.0.0.1:3306)/uretim_takip1?parseTime==true"
	db, err = sql.Open("mysql", dns)
	if err != nil {
		log.Fatalf("Veri tabanına bağlanırken hata oluştu: %v", err)

	}
	//Bağlantı kontrol
	err = db.Ping()
	if err != nil {
		log.Fatalf("Veritabanı bağlantısı kontrol eilemedir: %v", err)

	}
	fmt.Println("MYSQL veritabanına başarıyla bağlandı.")

}

//users tablosundaki kullanıcı kayıtlarını temsin eder

/*type usernames struct {
	Id        int       `json: "ID"`
	Username  string    `json:"UserName`
	Password  string    `json:"PassWord"`
	Mail      string    `json:"Mail"`
	CreatedAT time.Time `json:"CreatedAT"`
}*/

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Lütfen post isteği gönderin.", http.StatusMethodNotAllowed)
		return
	}

	// JSON formatındaki veriyi a
	var u usernames
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		log.Println("JSON parse hatası", err)
		return
	}

	// kullanıcıları veri tabanında çekme

	var storedHash string
	err = db.QueryRow("SELECT Password FROM Usernames WHERE username=?").Scan(&storedHash)
	if err != nil {
		http.Error(w, "kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}
	//şifreleri karşılaştırma

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(u.Password))
	if err != nil {
		http.Error(w, "hatalı şifre", http.StatusUnauthorized)
		return
	}

	//başarılı giriş
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Giriş Başarılı")

}
