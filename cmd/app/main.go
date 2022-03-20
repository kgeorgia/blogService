package main

import (
    "blog_service2/internal/app"
    "blog_service2/internal/controller/handler"
    "blog_service2/internal/repository"
    _ "github.com/lib/pq"
    "log"
    "net/http"
    "os"
)

func main() {
    port := ":8080"

    dbUser, dbPassword, dbName :=
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME")

    repo, err := repository.New(dbUser, dbPassword, dbName)
    if err != nil {
        log.Fatal(err)
    }

    handler := handler.New(repo)

    router := app.NewRouter(handler)

    err = http.ListenAndServe(port, router)
    log.Fatal(err)
}
