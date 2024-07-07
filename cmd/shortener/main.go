package main

import (
	"fmt"
	"github.com/sqwa11/first_sprint/internal/app/GET"
	"github.com/sqwa11/first_sprint/internal/app/POST"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", POST.HandleShorten)
	http.HandleFunc("/{id}", GET.HandleRedirect)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
