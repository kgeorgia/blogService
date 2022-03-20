package main

import (
	"blog_service2/internal/app"
	"blog_service2/internal/controller/handler"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	port := ":8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	handler, err := handler.New()
	if err != nil {
		log.Fatal(err)
	}

	router := app.NewRouter(handler)

	err = http.ListenAndServe(port, router)
	log.Fatal(err)
}
