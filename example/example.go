package main

import (
	"github.com/muesli/regommend"
	"fmt"
)

func main() {
	// Accessing a new regommend table for the first time will create it.
	books := regommend.Table("books")

	booksChrisRead := make(map[interface{}]float64)
	booksChrisRead["1984"] = 5.0
	booksChrisRead["Robinson Crusoe"] = 4.0
	booksChrisRead["Moby-Dick"] = 3.0
	books.Add("Chris", booksChrisRead)

	booksJayRead := make(map[interface{}]float64)
	booksJayRead["1984"] = 5.0
	booksJayRead["Robinson Crusoe"] = 4.0
	booksJayRead["Gulliver's Travels"] = 4.5
	books.Add("Jay", booksJayRead)

	recs, _ := books.Recommend("Chris")
	for _, rec := range recs {
		fmt.Println("Recommending", rec.Key, "with score:", rec.Distance)
	}
}
