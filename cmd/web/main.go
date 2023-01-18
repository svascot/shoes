package main

import (
	"fmt"
	"net/http"
	"shoes/cmd/web/pkg/shoes"
)

const portNumber = ":9090"

func main() {

	http.HandleFunc("/", shoes.GetShoes)
	http.HandleFunc("/mockShoes", shoes.CreateMockData)

	shoes.CreateTable()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
