package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var (
	produks = make([]Produk, 0)
	lastID  = 0
	mu      sync.RWMutex
)

func main() {
	http.HandleFunc("/produk", produkHandler)
	http.HandleFunc("/produk/", produkByIDHandler)

	log.Println("REST API running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

/*
POST   /produk
GET    /produk?nama=kopi&minHarga=10000&maxHarga=50000
*/
func produkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodPost:
		var p Produk
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		lastID++
		p.ID = lastID
		produks = append(produks, p)
		mu.Unlock()

		json.NewEncoder(w).Encode(p)

	case http.MethodGet:
		nama := strings.ToLower(r.URL.Query().Get("nama"))
		minHarga, _ := strconv.Atoi(r.URL.Query().Get("minHarga"))
		maxHarga, _ := strconv.Atoi(r.URL.Query().Get("maxHarga"))

		mu.RLock()
		defer mu.RUnlock()

		result := make([]Produk, 0)

		for _, p := range produks {
			if nama != "" && !strings.Contains(strings.ToLower(p.Nama), nama) {
				continue
			}
			if minHarga > 0 && p.Harga < minHarga {
				continue
			}
			if maxHarga > 0 && p.Harga > maxHarga {
				continue
			}
			result = append(result, p)
		}

		json.NewEncoder(w).Encode(result)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

/*
GET    /produk/{id}
PUT    /produk/{id}
DELETE /produk/{id}
*/
func produkByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		mu.RLock()
		defer mu.RUnlock()

		for _, p := range produks {
			if p.ID == id {
				json.NewEncoder(w).Encode(p)
				return
			}
		}
		http.Error(w, "produk not found", http.StatusNotFound)

	case http.MethodPut:
		var input Produk
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		for i, p := range produks {
			if p.ID == id {
				produks[i].Nama = input.Nama
				produks[i].Harga = input.Harga
				produks[i].Stok = input.Stok
				json.NewEncoder(w).Encode(produks[i])
				return
			}
		}
		http.Error(w, "produk not found", http.StatusNotFound)

	case http.MethodDelete:
		mu.Lock()
		defer mu.Unlock()

		for i, p := range produks {
			if p.ID == id {
				produks = append(produks[:i], produks[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "produk not found", http.StatusNotFound)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
