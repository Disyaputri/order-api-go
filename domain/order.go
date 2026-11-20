// Package domain berisi inti aplikasi: struktur data, error, dan kontrak (interface).
// Lapisan ini tidak bergantung pada database maupun HTTP.
package domain

import "errors"

// Error umum yang dipakai antar lapisan.
var (
	ErrNotFound     = errors.New("order not found")
	ErrInvalidInput = errors.New("invalid input")
)

// Order adalah data lengkap satu order (sesuai kolom tabel orders).
// CreatedAt berbentuk tanggal "2006-01-02" sesuai kontrak API.
type Order struct {
	ID        int64  `json:"id"`
	BuyerName string `json:"buyer_name"`
	StoreName string `json:"store_name"`
	ItemName  string `json:"item_name"`
	ItemQty   int    `json:"item_qty"`
	CreatedAt string `json:"created_at"`
}

// OrderInput adalah isi request body untuk POST dan PUT.
type OrderInput struct {
	BuyerName string `json:"buyer_name"`
	StoreName string `json:"store_name"`
	ItemName  string `json:"item_name"`
	ItemQty   int    `json:"item_qty"`
}

// OrderRepository adalah kontrak akses data (diimplementasi di folder repository).
type OrderRepository interface {
	Create(o *Order) error
	GetByID(id int64) (*Order, error)
	GetAll() ([]Order, error)
	Update(o *Order) error
	Delete(id int64) error
}

// OrderUsecase adalah kontrak logika bisnis (diimplementasi di folder usecase).
type OrderUsecase interface {
	Create(in OrderInput) (*Order, error)
	GetByID(id int64) (*Order, error)
	GetAll() ([]Order, error)
	Update(id int64, in OrderInput) (*Order, error)
	Delete(id int64) error
}
