# Order API (Go, Gin, MySQL)

REST API Go + MySQL untuk mengelola order (CRUD), dengan clean architecture.

## Endpoint
| Method | Path | Fungsi |
|---|---|---|
| GET | /orders | daftar order |
| POST | /orders | buat order |
| GET | /orders/:id | detail order |
| PUT | /orders/:id | ubah order |
| DELETE | /orders/:id | hapus order |

## Struktur folder
```
domain/      struktur data, error, interface (inti, tanpa database/HTTP)
usecase/     logika bisnis + validasi input
repository/  query MySQL
delivery/    router + handler HTTP
ddl.sql      DDL tabel orders + data contoh
main.go      menyambungkan semua lapisan
```

## Menjalankan lokal
1. Jalankan `ddl.sql` di MySQL.
2. `cp .env.example .env`, isi `DB_*` sesuai database.
3. `go mod tidy` lalu `go run .`
4. Tes: `go test ./...`
