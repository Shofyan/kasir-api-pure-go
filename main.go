// Package main provides a simple REST API for Kasir (Cashier) system
// @title Kasir API
// @version 1.0
// @description A simple REST API for managing products and categories
// @host localhost:8080
// @BasePath /
package main

import (
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBPass string `mapstructure:"DB_PASSWORD"`
	DBName string `mapstructure:"DB_NAME"`
	DBCON  string `mapstructure:"DB_CON"`
}

func main() {

	viper.SetConfigType("env")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	viper.AutomaticEnv()

	config := Config{
		Port:   viper.GetString("PORT"),
		DBHost: viper.GetString("DB_HOST"),
		DBPort: viper.GetString("DB_PORT"),
		DBUser: viper.GetString("DB_USER"),
		DBPass: viper.GetString("DB_PASSWORD"),
		DBName: viper.GetString("DB_NAME"),
	}

	// Build Postgres connection string (keyword format)
	dbConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require connect_timeout=30", config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)

	// Log connection info (without password)
	log.Printf("Connecting to database: postgres://%s:***@%s:%s/%s", config.DBUser, config.DBHost, config.DBPort, config.DBName)

	// Setup database
	db, err := database.InitDB(dbConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// new routes

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// category routes
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)
	http.HandleFunc("/api/categories", categoryHandler.HandleCategories)
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID)

	log.Println("REST API running at :" + config.Port)
	log.Println("Dashboard: http://localhost:" + config.Port)
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
