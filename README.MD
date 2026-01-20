# 🛒 Golang REST API CRUD Produk (Standard Library Only)

Project ini adalah contoh **REST API CRUD Produk** yang dibangun menggunakan **library bawaan Golang (standard library)** tanpa framework tambahan seperti Gin, Echo, atau Fiber.

Tujuan utama project ini adalah untuk memahami **fondasi backend Golang yang benar**, ringan, cepat, dan mudah dikembangkan ke arsitektur yang lebih besar.

---

## ✨ Fitur

* Create produk
* Read semua produk
* Read produk berdasarkan ID
* Update produk
* Delete produk
* Search produk berdasarkan:

  * Nama (partial match)
  * Harga minimum
  * Harga maksimum
* Thread-safe (menggunakan `sync.RWMutex`)
* Tanpa dependency eksternal

---

## 🧱 Struktur Data

```go
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}
```

---

## 🚀 Cara Menjalankan

### 1. Pastikan Go sudah terinstall

```bash
go version
```

Minimal disarankan:

```
go 1.20+
```

---

### 2. Jalankan aplikasi

```bash
go run main.go
```

Server akan berjalan di:

```
http://localhost:8080
```

---

## 📡 Endpoint API

### ➕ Create Produk

**POST** `/produk`

```json
{
  "nama": "Laptop",
  "harga": 15000000,
  "stok": 5
}
```

Response:

```json
{
  "id": 1,
  "nama": "Laptop",
  "harga": 15000000,
  "stok": 5
}
```

---

### 📄 Get Semua Produk

**GET** `/produk`

```bash
curl http://localhost:8080/produk
```

---

### 🔍 Search Produk

**GET** `/produk`

Query parameter:

| Parameter | Keterangan                            |
| --------- | ------------------------------------- |
| nama      | Cari berdasarkan nama (partial match) |
| minHarga  | Harga minimum                         |
| maxHarga  | Harga maksimum                        |

Contoh:

```bash
curl "http://localhost:8080/produk?nama=lap&minHarga=1000000"
```

---

### 🔎 Get Produk by ID

**GET** `/produk/{id}`

```bash
curl http://localhost:8080/produk/1
```

---

### ✏️ Update Produk

**PUT** `/produk/{id}`

```json
{
  "nama": "Laptop Pro",
  "harga": 17000000,
  "stok": 3
}
```

---

### ❌ Delete Produk

**DELETE** `/produk/{id}`

```bash
curl -X DELETE http://localhost:8080/produk/1
```

Response:

```
204 No Content
```

---

## 🧠 Catatan Teknis

* Data disimpan **in-memory**
* Tidak menggunakan database
* Cocok untuk:

  * belajar REST API Golang
  * interview backend
  * prototype
  * bahan mentoring / onboarding

---

## 🔥 Kenapa Menggunakan Standard Library?

* Performa tinggi
* Dependency nol
* Binary kecil
* Mudah di-debug
* Cocok untuk microservice dan system core

Framework bisa digunakan kemudian, namun **fondasi stdlib wajib dipahami terlebih dahulu**.

---

## 🚧 Pengembangan Selanjutnya

Beberapa improvement yang bisa dikembangkan:

* Integrasi database (PostgreSQL / MySQL)
* Clean Architecture (handler → service → repository)
* Pagination
* Sorting
* Validation
* Middleware (logging, request-id)
* Graceful shutdown
* Docker support

---

## 👨‍💻 Author

**Sofyan Sky**
Backend Engineer – Golang

---

## 📜 License

MIT License

---

> "Framework helps you move fast.
> Understanding the standard library helps you move far."

🔥
