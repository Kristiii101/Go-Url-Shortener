package main

import (
	"fmt"
	"net/http"

	myhttp "github.com/Kristiii101/GO_URL_Shortener_ATAD/internal/http"
)

func main() {
	fmt.Println("Starting URL shortener on :8080")

	r := myhttp.NewRouter()

	http.ListenAndServe(":8080", r)
}
