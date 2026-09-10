package main

import (
	"log"
	"net/http"

	apphttp "github.com/rikut0904/starter/create/go/internal/interface/http"
	"github.com/rikut0904/starter/create/go/internal/usecase"
)

func main() {
	handler := apphttp.NewHandler(usecase.NewHealth())
	log.Println("listening on 0.0.0.0:8080")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
}
