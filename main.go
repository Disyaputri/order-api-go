// Program utama: membaca konfigurasi, membuka koneksi MySQL, lalu menjalankan server HTTP.
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"lc1-order-api/delivery"
	"lc1-order-api/repository"
	"lc1-order-api/usecase"
)

// env membaca environment variable, memakai nilai bawaan bila kosong.
func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	// File .env hanya ada di komputer lokal.
	_ = godotenv.Load()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		env("DB_USER", "root"), os.Getenv("DB_PASS"),
		env("DB_HOST", "localhost"), env("DB_PORT", "3306"), env("DB_NAME", "lc1_order"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("gagal terhubung ke database: ", err)
	}

	// Susunan lapisan: repository -> usecase -> delivery (HTTP).
	repo := repository.NewOrderMySQL(db)
	uc := usecase.NewOrderUsecase(repo)
	router := delivery.NewRouter(uc)

	log.Fatal(router.Run(":" + env("PORT", "8080")))
}
