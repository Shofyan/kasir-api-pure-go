// Package main provides a simple REST API for Kasir (Cashier) system
// @title Kasir API
// @version 1.0
// @description A simple REST API for managing products and categories
// @host localhost:8080
// @BasePath /
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Produk struct {
	ID         int    `json:"id"`
	Nama       string `json:"nama"`
	Harga      int    `json:"harga"`
	Stok       int    `json:"stok"`
	CategoryID int    `json:"category_id"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var (
	produks    = make([]Produk, 0)
	categories = make([]Category, 0)
	lastID     = 0
	lastCatID  = 0
	mu         sync.RWMutex
)

func main() {
	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/dashboard", indexHandler)
	http.HandleFunc("/products.html", htmlHandler("products.html"))
	http.HandleFunc("/categories.html", htmlHandler("categories.html"))

	// API routes
	http.HandleFunc("/produk", produkHandler)
	http.HandleFunc("/produk/", produkByIDHandler)
	http.HandleFunc("/categories", categoriesHandler)
	http.HandleFunc("/categories/", categoriesByIDHandler)

	// Helper routes
	http.HandleFunc("/api/stats", statsHandler)
	http.HandleFunc("/api/category-options", categoryOptionsHandler)

	log.Println("REST API running at :8080")
	log.Println("Dashboard: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// HTML handlers
func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func htmlHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/"+filename)
	}
}

// @Summary Get API statistics
// @Description Get statistics about products and categories count
// @Tags dashboard
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/stats [get]

// Stats handler for dashboard
func statsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mu.RLock()
	defer mu.RUnlock()

	stats := map[string]interface{}{
		"total_products":   len(produks),
		"total_categories": len(categories),
		"last_product_id":  lastID,
		"last_category_id": lastCatID,
	}

	err := json.NewEncoder(w).Encode(stats)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Helper function to check if request is from HTMX
func isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

// Helper function to render products as HTML table
func renderProductsHTML(products []Produk) string {
	if len(products) == 0 {
		return "<p class='text-muted'>No products found.</p>"
	}

	html := `<div class="table-responsive">
		<table class="table table-hover">
		<thead>
			<tr>
				<th>ID</th>
				<th>Name</th>
				<th>Price</th>
				<th>Stock</th>
				<th>Category</th>
				<th>Actions</th>
			</tr>
		</thead>
		<tbody>`

	for _, p := range products {
		// Find category name
		categoryName := "No Category"
		for _, c := range categories {
			if c.ID == p.CategoryID {
				categoryName = c.Name
				break
			}
		}

		html += fmt.Sprintf(`
			<tr>
				<td>%d</td>
				<td>%s</td>
				<td>Rp %d</td>
				<td>%d</td>
				<td>%s</td>
				<td>
					<button class="btn btn-sm btn-warning me-1" onclick="editProduct(%d, '%s', %d, %d, %d)">Edit</button>
					<button class="btn btn-sm btn-danger" onclick="deleteProduct(%d)">Delete</button>
				</td>
			</tr>`, p.ID, p.Nama, p.Harga, p.Stok, categoryName, p.ID, p.Nama, p.Harga, p.Stok, p.CategoryID, p.ID)
	}

	html += "</tbody></table></div>"
	return html
}

// Helper function to render categories as HTML table
func renderCategoriesHTML(categories []Category) string {
	if len(categories) == 0 {
		return "<p class='text-muted'>No categories found.</p>"
	}

	html := `<div class="table-responsive">
		<table class="table table-hover">
		<thead>
			<tr>
				<th>ID</th>
				<th>Name</th>
				<th>Description</th>
				<th>Actions</th>
			</tr>
		</thead>
		<tbody>`

	for _, c := range categories {
		html += fmt.Sprintf(`
			<tr>
				<td>%d</td>
				<td>%s</td>
				<td>%s</td>
				<td>
					<button class="btn btn-sm btn-warning me-1" onclick="editCategory(%d, '%s', '%s')">Edit</button>
					<button class="btn btn-sm btn-danger" onclick="deleteCategory(%d)">Delete</button>
				</td>
			</tr>`, c.ID, c.Name, c.Description, c.ID, c.Name, c.Description, c.ID)
	}

	html += "</tbody></table></div>"
	return html
}

// Helper function to render categories as HTML options for select dropdown
func renderCategoryOptionsHTML(categories []Category) string {
	html := ""
	for _, c := range categories {
		html += fmt.Sprintf(`<option value="%d">%s</option>`, c.ID, c.Name)
	}
	return html
}

// @Summary Get category options for dropdown
// @Description Get categories formatted as HTML options for select dropdown
// @Tags categories
// @Produce text/html
// @Success 200 {string} string "HTML options"
// @Router /api/category-options [get]

// Handler for category options dropdown
func categoryOptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	mu.RLock()
	defer mu.RUnlock()

	w.Write([]byte(renderCategoryOptionsHTML(categories)))
}

// @Summary Create a new product
// @Description Create a new product with name, price, stock and category
// @Tags products
// @Accept json
// @Produce json
// @Param product body Produk true "Product data"
// @Success 201 {object} Produk
// @Failure 400 {string} string "Bad request"
// @Router /produk [post]

// @Summary Get all products
// @Description Get all products with optional filtering by name and price range
// @Tags products
// @Produce json
// @Param nama query string false "Product name filter"
// @Param minHarga query int false "Minimum price filter"
// @Param maxHarga query int false "Maximum price filter"
// @Success 200 {array} Produk
// @Router /produk [get]

/*
POST   /produk
GET    /produk?nama=kopi&minHarga=10000&maxHarga=50000
*/
func produkHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		var p Produk

		// Check if this is form data from HTMX
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			p.Nama = r.FormValue("nama")
			p.Harga, _ = strconv.Atoi(r.FormValue("harga"))
			p.Stok, _ = strconv.Atoi(r.FormValue("stok"))
			p.CategoryID, _ = strconv.Atoi(r.FormValue("category_id"))
		} else {
			// Handle JSON data
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		mu.Lock()
		lastID++
		p.ID = lastID
		produks = append(produks, p)
		mu.Unlock()

		err := json.NewEncoder(w).Encode(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

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

		// Check if this is an HTMX request
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			_, err := w.Write([]byte(renderProductsHTML(result)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(result)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Get a product by ID
// @Description Get a single product by its ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} Produk
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Product not found"
// @Router /produk/{id} [get]

// @Summary Update a product
// @Description Update an existing product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Produk true "Updated product data"
// @Success 200 {object} Produk
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Product not found"
// @Router /produk/{id} [put]

// @Summary Delete a product
// @Description Delete a product by ID
// @Tags products
// @Param id path int true "Product ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Product not found"
// @Router /produk/{id} [delete]

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
				err := json.NewEncoder(w).Encode(p)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
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
				produks[i].CategoryID = input.CategoryID
				err := json.NewEncoder(w).Encode(produks[i])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
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

// @Summary Create a new category
// @Description Create a new category with name and description
// @Tags categories
// @Accept json
// @Produce json
// @Param category body Category true "Category data"
// @Success 201 {object} Category
// @Failure 400 {string} string "Bad request"
// @Router /categories [post]

// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} Category
// @Router /categories [get]

/*
POST   /categories
GET    /categories
*/
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		w.Header().Set("Content-Type", "application/json")
		var c Category

		// Check if this is form data from HTMX
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			// Parse form data
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			c.Name = r.FormValue("name")
			c.Description = r.FormValue("description")
		} else {
			// Handle JSON data
			if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		mu.Lock()
		lastCatID++
		c.ID = lastCatID
		categories = append(categories, c)
		mu.Unlock()

		err := json.NewEncoder(w).Encode(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodGet:
		mu.RLock()
		defer mu.RUnlock()

		// Check if this is an HTMX request
		if isHTMXRequest(r) {
			w.Header().Set("Content-Type", "text/html")
			_, err := w.Write([]byte(renderCategoriesHTML(categories)))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(categories)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// @Summary Get a category by ID
// @Description Get a single category by its ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} Category
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [get]

// @Summary Update a category
// @Description Update an existing category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body Category true "Updated category data"
// @Success 200 {object} Category
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [put]

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Param id path int true "Category ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Category not found"
// @Router /categories/{id} [delete]

/*
GET    /categories/{id}
PUT    /categories/{id}
DELETE /categories/{id}
*/
func categoriesByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		mu.RLock()
		defer mu.RUnlock()

		for _, c := range categories {
			if c.ID == id {
				err := json.NewEncoder(w).Encode(c)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}
		http.Error(w, "category not found", http.StatusNotFound)

	case http.MethodPut:
		var input Category
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()
		defer mu.Unlock()

		for i, c := range categories {
			if c.ID == id {
				categories[i].Name = input.Name
				categories[i].Description = input.Description
				err := json.NewEncoder(w).Encode(categories[i])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}
		http.Error(w, "category not found", http.StatusNotFound)

	case http.MethodDelete:
		mu.Lock()
		defer mu.Unlock()

		for i, c := range categories {
			if c.ID == id {
				categories = append(categories[:i], categories[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "category not found", http.StatusNotFound)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
