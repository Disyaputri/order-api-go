// Package repository berisi akses data ke MySQL.
package repository

import (
	"database/sql"
	"errors"

	"lc1-order-api/domain"
)

// kolom yang dibaca; created_at diubah ke format tanggal "2006-01-02" langsung oleh MySQL.
const selectColumns = "id, buyer_name, store_name, item_name, item_qty, DATE_FORMAT(created_at, '%Y-%m-%d')"

type orderMySQL struct {
	db *sql.DB
}

// NewOrderMySQL membuat repository order berbasis MySQL.
func NewOrderMySQL(db *sql.DB) domain.OrderRepository {
	return &orderMySQL{db: db}
}

// scan membaca satu baris hasil query ke struct Order.
func scan(row interface{ Scan(...any) error }, o *domain.Order) error {
	return row.Scan(&o.ID, &o.BuyerName, &o.StoreName, &o.ItemName, &o.ItemQty, &o.CreatedAt)
}

func (r *orderMySQL) Create(o *domain.Order) error {
	res, err := r.db.Exec(
		"INSERT INTO orders (buyer_name, store_name, item_name, item_qty) VALUES (?, ?, ?, ?)",
		o.BuyerName, o.StoreName, o.ItemName, o.ItemQty,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	// Ambil ulang agar created_at terisi dari database.
	saved, err := r.GetByID(id)
	if err != nil {
		return err
	}
	*o = *saved
	return nil
}

func (r *orderMySQL) GetByID(id int64) (*domain.Order, error) {
	var o domain.Order
	err := scan(r.db.QueryRow("SELECT "+selectColumns+" FROM orders WHERE id = ?", id), &o)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *orderMySQL) GetAll() ([]domain.Order, error) {
	rows, err := r.db.Query("SELECT " + selectColumns + " FROM orders ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := []domain.Order{} // slice kosong (bukan nil) supaya JSON tampil [] bukan null
	for rows.Next() {
		var o domain.Order
		if err := scan(rows, &o); err != nil {
			return nil, err
		}
		list = append(list, o)
	}
	return list, rows.Err()
}

func (r *orderMySQL) Update(o *domain.Order) error {
	_, err := r.db.Exec(
		"UPDATE orders SET buyer_name = ?, store_name = ?, item_name = ?, item_qty = ? WHERE id = ?",
		o.BuyerName, o.StoreName, o.ItemName, o.ItemQty, o.ID,
	)
	return err
}

func (r *orderMySQL) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM orders WHERE id = ?", id)
	return err
}
