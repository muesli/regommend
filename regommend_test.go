/*
 * Simple recommendation engine
 *     Copyright (c) 2014, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package regommend

import (
	"testing"
)

var (
)

func TestEngine(t *testing.T) {
	books := Table("books")

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

	// check if both items are still there
	p, err := books.Value("Chris")
	if err != nil || p == nil {
		t.Error("Error retrieving item from engine", err)
	}
	p, err = books.Value("Jay")
	if err != nil || p == nil {
		t.Error("Error retrieving item from engine", err)
	}
}

func TestNeighbors(t *testing.T) {
	books := Table("books")

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

	booksMaryRead := make(map[interface{}]float64)
	booksMaryRead["1984"] = 4.0
	booksMaryRead["Robinson Crusoe"] = 3.0
	booksMaryRead["Gulliver's Travels"] = 4.5
	books.Add("Mary", booksMaryRead)

	booksJackRead := make(map[interface{}]float64)
	booksJackRead["1984"] = 3.0
	booksJackRead["Robinson Crusoe"] = 1.0
	books.Add("Jack", booksJackRead)

	nbs, _ := books.Neighbors("Chris")
	if len(nbs) != 3 {
		t.Error("Expected 3 neighbours, got", len(nbs))
	}
	if nbs[0].Key != "Jay" || nbs[1].Key != "Mary" || nbs[2].Key != "Jack" {
		t.Error("Unexpected similarity order")
	}
}

func TestRecommendations(t *testing.T) {
	books := Table("books")

	booksChrisRead := make(map[interface{}]float64)
	booksChrisRead["1984"] = 5.0
	booksChrisRead["Robinson Crusoe"] = 4.0
	booksChrisRead["Moby-Dick"] = 3.0
	books.Add("Chris", booksChrisRead)

	booksJayRead := make(map[interface{}]float64)
	booksJayRead["1984"] = 4.0
	booksJayRead["Robinson Crusoe"] = 3.0
	booksJayRead["Gulliver's Travels"] = 4.5
	booksJayRead["A Tale of Two Cities"] = 3.5
	books.Add("Jay", booksJayRead)

	recs, _ := books.Recommend("Chris")
	if len(recs) != 2 {
		t.Error("Expected 2 recommendations, got", len(recs))
	}
	if recs[0].Key != "Gulliver's Travels" || recs[1].Key != "A Tale of Two Cities" {
		t.Error("Unexpected recommendation order")
	}
}
