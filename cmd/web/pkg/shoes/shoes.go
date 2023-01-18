package shoes

import (
	"encoding/json"
	"log"
	"net/http"
)

type Shoe struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

func GetShoes(w http.ResponseWriter, r *http.Request) {
	log.Println("GetShoes method")

	var moreSneakers []Shoe

	var shoeOne Shoe
	shoeOne.Color = "White"
	shoeOne.Name = "Jordan 1 low SE triple White"
	shoeOne.Type = "Jordan 1 low"

	moreSneakers = append(moreSneakers, shoeOne)

	var shoeTwo Shoe
	shoeTwo.Color = "Red"
	shoeTwo.Name = "Jordan 4 Metallic Red"
	shoeTwo.Type = "Jordan 4"

	moreSneakers = append(moreSneakers, shoeTwo)

	newJson, err := json.Marshal(moreSneakers)

	if err != nil {
		log.Println("Error marshalling json", err)
	}

	log.Println(newJson)
	w.Write(newJson)

}
