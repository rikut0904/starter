package main

import (
	apphttp "github.com/rikut0904/starter/create/next-go/backend/internal/interface/http"
	"github.com/rikut0904/starter/create/next-go/backend/internal/usecase"
	"log"
	"net/http"
)

func main() {
	handler := apphttp.NewHandler(usecase.NewHealth())
	log.Println("listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}
