/*
 * Simple recommendation engine
 *     Copyright (c) 2014, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package regommend

import (
	"errors"
	"log"
	"fmt"
	_ "sort"
	"math"
	"sort"
	"sync"
	_ "time"
)

// Structure of a table with items in the engine.
type RegommendTable struct {
	sync.RWMutex

	// The table's name.
	name string
	// All items in the table.
	items map[interface{}]*RegommendItem

	// The logger used for this table.
	logger *log.Logger

	// Callback method triggered when trying to load a non-existing key.
	loadData func(key interface{}) *RegommendItem
	// Callback method triggered when adding a new item to the engine.
	addedItem func(item *RegommendItem)
	// Callback method triggered before deleting an item from the engine.
	aboutToDeleteItem func(item *RegommendItem)
}
