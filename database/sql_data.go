package sql_data

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

// uretim_takip veritabanındaki parça kayıtlarını temsil eder
type uretim_takip struct {
	ID               int       `json:"ID"`
	ParcaAdi         string    `json:"ParcaAdi"`
	ParcaMalzemesi   string    `json:"ParcaMalzemesi"`
	UretimSekli      string    `json:"UretimSekli"`
	UretimAdedi      int       `json:"UretimAdedi"`
	KayitTarihi      time.Time `json:"KayitTarihi"`
	GuncellemeTarihi time.Time `json:"update_time"`
}

// HomePage basit bir karşılama mesajı döndürür
func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hoşgeldiniz! Parça eklemek için /addpiece adresine POST isteği gönderin.")
}

// _piece fonksiyonu POST isteği ile gelen JSON verisini alıp MySQL'e ekler
func Piece(w http.ResponseWriter, r *http.Request) {
	// Gelen HTTP metodunu terminale yazdıralım
	//log.Println("Gelen HTTP Metodu:", r.Method)
	if r.Method != http.MethodGet {

		http.Error(w, "Lütfen GET isteği gönderin", http.StatusMethodNotAllowed)
		return
	}

	// Verileri MySQL'den çek
	rows, err := db.Query("SELECT ID, ParcaAdi, ParcaMalzemesi, UretimSekli, UretimAdedi, KayitTarihi FROM uretim_takip")
	if err != nil {
		log.Println("❌ SQL Hatası:", err) // Terminale hatayı yazdır
		http.Error(w, "Veri çekilirken hata oluştu", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pieces []uretim_takip

	for rows.Next() {
		var p uretim_takip
		if err := rows.Scan(&p.ID, &p.ParcaAdi, &p.ParcaMalzemesi, &p.UretimSekli, &p.UretimAdedi, &p.KayitTarihi); err != nil {
			log.Println("❌ SQL Hatası:", err) // Terminale hatayı yazdır
			http.Error(w, "Veri işlenirken hata oluştu", http.StatusInternalServerError)
			return
		}
		pieces = append(pieces, p)
	}
	// JSON olarak yanıt döndür
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pieces)
}

// Kullanıcıdan gelen JSON verisini alıp veritabanına ekleyen fonksiyon
func AddPiece(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		http.Error(w, "Lütfen POST isteği gönderin", http.StatusMethodNotAllowed)
		return
	}

	var p uretim_takip

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "geçersiz json formatı", http.StatusBadRequest)
		log.Println("json parse hatası", err)
		return
	}

	// Kayıt tarihini ekleyelim
	//const time.DateTime untyped string = "2006-01-02 15:04:05"
	p.KayitTarihi = time.Now()

	// MySQL'e ekleme sorgusu
	query := `INSERT INTO uretim_takip (ParcaAdi, ParcaMalzemesi, UretimSekli, UretimAdedi, kayitTarihi) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, p.ParcaAdi, p.ParcaMalzemesi, p.UretimSekli, p.UretimAdedi, p.KayitTarihi)
	if err != nil {
		http.Error(w, fmt.Sprintf("Veri eklenirken hata: %v", err), http.StatusInternalServerError)
		log.Println("❌ SQL Hatası:", err) // Terminale hatayı yazdır
		return
	}

	// Eklenen kaydın ID'sini al
	insertedID, err := result.LastInsertId()
	if err != nil {
		log.Println("❌ SQL Hatası:", err) // Terminale hatayı yazdırs
		http.Error(w, "Kayıt ID alınırken hata", http.StatusInternalServerError)
		return
	}
	p.ID = int(insertedID)

	// JSON olarak başarılı yanıt döndür
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Yeni kayıt başarıyla eklendi!",
		"data":    p,
	})
}

//kayıt silme

func Delete_Piece(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Lütfen Delete isteği gönderin", http.StatusMethodNotAllowed)
		return
	}
	//json formatındaki veriyi alma
	var p uretim_takip
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Geçersiz JSON formatı", http.StatusBadRequest)
		log.Println("json perse hatası", err)
		return
	}
	// parçayı veri tabanından silme
	result, err := db.Exec("DELETE FROM uretim_takip WHERE ParcaAdi= ?", p.ParcaAdi)
	if err != nil {
		http.Error(w, "Parça silinirken hata oluştur", http.StatusInternalServerError)
		log.Println("❌ SQL Hatası:", err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "✅ Kullanıcı '%s' başarıyla silindi.\n", p.ParcaAdi)

}
func UpdatePiece(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Lütfen PUT İSteği gönderin", http.StatusMethodNotAllowed)
		return
	}
	var p uretim_takip
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}
	//update time
	p.GuncellemeTarihi = time.Now()

	//update query
	query := `
	UPDATE uretim_takip
	SET ParcaAdi=?,ParcaMlazemesi=?, UretimSekli=?,UretimAdedi=?,GuncellemeTarihi=?
	WHERE ID=?
	`
	result, err := db.Exec(query, p.ParcaAdi, p.ParcaMalzemesi, p.UretimSekli, p.UretimAdedi, p.GuncellemeTarihi, p.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Güncelleme sırasında hata oluştu:%v", err), http.StatusInternalServerError)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Güncellenecek kayıt bulunamadı", http.StatusNotFound)
		return
	}
	//send reply
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"massage": "Parça başarıyla güncellendi",
		"data":    p,
	})
}
