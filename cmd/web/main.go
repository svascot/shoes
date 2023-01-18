package main

import (
	"fmt"
	"net/http"
	"shoes/cmd/web/pkg/shoes"
)

const portNumber = ":9090"

func main() {

	http.HandleFunc("/", shoes.GetShoes)

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))
	_ = http.ListenAndServe(portNumber, nil)
}
