package main

import (
	"19_11_2026_go/internal/app"
	"19_11_2026_go/internal/checkers"
	"19_11_2026_go/internal/storages"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const filePath = "storages.json"

func main() {
	storage, err := storages.NewStorage(filePath)
	if err != nil {
		log.Fatalf("Error creating new storage %v", err)
	}

	checker := checkers.NewChecker(5 * time.Second)

	linkService := app.NewService(storage, checker)

	handler := NewHandler(linkService)

	server := http.Server{
		Addr:    ":8080",
		Handler: routes(handler),
	}

	go func() {
		log.Printf("Starting web server at localhost:%s\n", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting web server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Saving data to file")
	if err := storage.Save(); err != nil {
		log.Printf("Error saving data: %v", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
