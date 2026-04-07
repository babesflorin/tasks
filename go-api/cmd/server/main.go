package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/workspace/go-api/internal/config"
	"github.com/workspace/go-api/internal/database"
	"github.com/workspace/go-api/internal/handler"
	"github.com/workspace/go-api/internal/repository"
	"github.com/workspace/go-api/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db, "/migrations"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	repo := repository.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	r := chi.NewRouter()
	h.RegisterRoutes(r)

	addr := ":" + cfg.ServerPort
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
