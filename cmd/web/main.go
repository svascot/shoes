package main

import (
	"fmt"
	"net/http"
	"shoes/cmd/web/pkg/shoes"

	"github.com/gorilla/mux"
)

const portNumber = ":9090"

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/shoes", shoes.GetAllShoes).Methods("GET")
	router.HandleFunc("/shoes/{id}", shoes.GetShoe).Methods("GET")
	router.HandleFunc("/shoes", shoes.AddShoe).Methods("POST")
	router.HandleFunc("/shoes/{id}/wear-today", shoes.WearToday).Methods("PUT")
	router.HandleFunc("/shoes/{id}", shoes.DeleteShoe).Methods("DELETE")

	router.HandleFunc("/createMockData", shoes.CreateMockData).Methods("GET")
	//shoes.CreateTable()

	fmt.Println(fmt.Sprintf("Starting application on port %s", portNumber))

	_ = http.ListenAndServe(portNumber, router)
}
