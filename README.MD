# 🛒 Golang REST API CRUD Produk & Kategori (Standard Library Only)

Project ini adalah contoh **REST API CRUD Produk & Kategori** yang dibangun menggunakan **library bawaan Golang (standard library)** tanpa framework tambahan seperti Gin, Echo, atau Fiber.

Tujuan utama project ini adalah untuk memahami **fondasi backend Golang yang benar**, ringan, cepat, dan mudah dikembangkan ke arsitektur yang lebih besar.

---

## ✨ Fitur

* Create, Read, Update, Delete produk
* Create, Read, Update, Delete kategori
* Clean Architecture (handler → service → repository)
* Integrasi database (PostgreSQL)
* Tanpa dependency eksternal

---

## 🧱 Struktur Data

```go
type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

### ➕ Create Category

**POST** `/api/categories`

- Accepts JSON, `application/x-www-form-urlencoded`, or `multipart/form-data`.

Contoh JSON:
```json
{
  "name": "Sabun",
  "description": "Sabun mandi"
}
```

Contoh curl form-data:
```bash
curl --location 'http://localhost:8080/api/categories' \
  --form 'name="sabun"' \
  --form 'description="sabun"'
```

Response:
```json
{
  "id": 1,
  "name": "Sabun",
  "description": "Sabun mandi"
}
```

---

### 📄 Get Semua Category

**GET** `/api/categories`

```bash
curl http://localhost:8080/api/categories
```

---

### 🔎 Get Category by ID

**GET** `/api/categories/{id}`

```bash
curl http://localhost:8080/api/categories/1
```

---

### ✏️ Update Category

**PUT** `/api/categories/{id}`

- Accepts JSON or form-data.

```json
{
  "name": "Sabun Cuci",
  "description": "Sabun untuk mencuci"
}
```

---

### ❌ Delete Category

**DELETE** `/api/categories/{id}`

```bash
curl -X DELETE http://localhost:8080/api/categories/1
```

Response:
```json
{"message": "Category deleted successfully"}
```

---

### 🔍 Search Produk

**GET** `/produk`

Query parameter:

| Parameter | Keterangan                            |
| --------- | ------------------------------------- |
| nama      | Cari berdasarkan nama (partial match) |

Contoh:

```bash
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

## � Testing

### Unit Tests

Jalankan unit tests:

```bash
go test -v
```

### API Testing

Gunakan script testing untuk menguji endpoint categories:

```bash
./test_categories.sh
```

Atau test manual dengan curl:

```bash
# Get all categories
curl -X GET "http://127.0.0.1:8080/api/categories"

# Create category
curl -X POST "http://127.0.0.1:8080/api/categories" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "description": "Test category"}'
```

---

## �🧠 Catatan Teknis

* Data produk/kategori disimpan di database (PostgreSQL)
* Clean Architecture: handler → service → repository
* Mendukung request dari HTMX (form) dan API (JSON)
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
