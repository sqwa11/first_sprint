package main

import (
	"fmt"
	"github.com/sqwa11/first_sprint/internal/app/get"
	"github.com/sqwa11/first_sprint/internal/app/post"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", post.HandleShorten)
	http.HandleFunc("/{id}", get.HandleRedirect)

	fmt.Println("Server listening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
