package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	addAuthSubrouter(r)

	fmt.Println("Initialized. Serving on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}