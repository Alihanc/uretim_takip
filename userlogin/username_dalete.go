package username

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func Username_delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Lütfen Delete İsteği Gönderin.", http.StatusMethodNotAllowed)
		return
	}

	//json formatındaki veriyi alma
	var u usernames
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Geçersiz Json Formatı", http.StatusBadRequest)
		log.Println("JSON perse hatası", err)
		return
	}
	//kullanıcıyı veri tabanından silme
	result, err := db.Exec("DELETE FROM usernames WHERE Username = ?", u.Username)
	if err != nil {
		http.Error(w, "kullanıcı silinirken hata oluştu", http.StatusInternalServerError)
		log.Println("❌ SQL Hatası:", err)
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Println(w, "✅ Kullanıcı '%s' başarıyla silindi.\n", u.Username)

}
