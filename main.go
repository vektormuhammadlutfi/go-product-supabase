package main

import (
	"encoding/json"
	"fmt"
	"gocats/internal/config"
	"gocats/internal/database"
	"gocats/internal/handlers"
	"gocats/internal/repository"
	"gocats/internal/services"
	"gocats/migrations"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the database
	db, err := database.New(database.Config{
		DSN: cfg.Database.DSN,
	})

	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Run database migrations if enabled
	if cfg.Database.AutoMigrate {
		log.Println("Running database migrations...")
		if err := migrations.RunMigrations(db); err != nil {
			log.Fatalf("Error running migrations: %v", err)
		}
	} else {
		log.Println("Auto migration disabled, skipping...")
	}

	// initialize repositories
	categoryRepo := repository.NewCategoryRepository(db.DB)
	productRepo := repository.NewProductRepository(db.DB)

	// initialize services
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo, categoryRepo)

	// initialize HTTP Handlers
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)

	// setup routes
	// health check endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {

		if err := db.HealthCheck(); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "error",
				"message": "Database connection failed",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"message": "Gocats services is running",
		})

	})

	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetAllCategories(w, r)
		case http.MethodPost:
			categoryHandler.CreateCategory(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetCategoryByID(w, r)
		case http.MethodPut:
			categoryHandler.UpdateCategory(w, r)
		case http.MethodDelete:
			categoryHandler.DeleteCategory(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Product routes
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Check if filtering by category
			if r.URL.Query().Get("category_id") != "" {
				productHandler.GetProductsByCategoryID(w, r)
			} else {
				productHandler.GetAllProducts(w, r)
			}
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.GetProductByID(w, r)
		case http.MethodPut:
			productHandler.UpdateProduct(w, r)
		case http.MethodDelete:
			productHandler.DeleteProduct(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("🚀 Server starting on %s...", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("❌ HTTP server failed: %v", err)
	}

}
