package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"packCalculator/server/db"
	"packCalculator/server/handlers"
	"packCalculator/server/repository"
	"packCalculator/server/service"
)

func main() {
	// Initialize with default pack sizes
	database, err := db.NewDatabase("packs.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	packRepo := repository.NewPackRepository(database.DB)
	packService, err := service.NewPackService(packRepo)
	if err != nil {
		log.Fatalf("Failed to create pack service: %v", err)
	}

	// Set up handlers
	packHandler := handlers.NewPackHandler(packService)

	// Create router
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// API routes
	router.Route("/api", func(r chi.Router) {
		r.Get("/calculate", packHandler.CalculatePacks)
		r.Get("/packs", packHandler.GetPackSizes)
		r.Post("/packs", packHandler.UpdatePackSizes)
		r.Post("/packs/add", packHandler.AddPackSize)
		r.Delete("/packs/{size}", packHandler.RemovePackSize)
	})

	// Serve frontend files
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/*", fs)

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
