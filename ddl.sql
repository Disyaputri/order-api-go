-- DDL Live Code 1 (P2): REST API Order
-- Jalankan file ini di MySQL sebelum aplikasi dijalankan.

CREATE DATABASE IF NOT EXISTS lc1_order;
USE lc1_order;

CREATE TABLE IF NOT EXISTS orders (
    id         INT UNSIGNED NOT NULL AUTO_INCREMENT,          -- identitas unik order (PRIMARY KEY otomatis ter-index, dipakai oleh query berdasarkan ID)
    buyer_name VARCHAR(100) NOT NULL,                         -- nama pembeli
    store_name VARCHAR(100) NOT NULL,                         -- nama toko penjual
    item_name  VARCHAR(100) NOT NULL,                         -- nama item yang dibeli
    item_qty   INT UNSIGNED NOT NULL,                         -- jumlah item yang dibeli
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- waktu order dibuat
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Data contoh (dummy) untuk percobaan endpoint.
INSERT INTO orders (buyer_name, store_name, item_name, item_qty) VALUES
    ('Budi Santoso',  'Toko Maju Jaya',   'Kaos Polos',     2),
    ('Siti Aminah',   'Toko Sinar Baru',  'Tas Selempang',  1),
    ('Andi Wijaya',   'Toko Maju Jaya',   'Sepatu Sneaker', 1),
    ('Dewi Lestari',  'Toko Cantik',      'Parfum 50ml',    3),
    ('Rizky Pratama', 'Toko Sinar Baru',  'Topi Baseball',  4);
