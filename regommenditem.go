/*
 * Simple recommendation engine
 *     Copyright (c) 2014, Christian Muehlhaeuser <muesli@gmail.com>
 *
 *   For license see LICENSE.txt
 */

package regommend

import (
	"sync"
	_ "time"
)

// Structure of an item in the recommendation engine.
// Parameter data contains the user-set value in the engine.
type RegommendItem struct {
	sync.RWMutex

	// The item's key.
	key interface{}
	// All items for this key.
	data map[interface{}]float64
}

// CreateRegommendItem returns a newly created RegommendItem.
// Parameter key is the item's key.
// Parameter data is the item's value.
func CreateRegommendItem(key interface{}, data map[interface{}]float64) RegommendItem {
	return RegommendItem{
		key:           key,
		data:          data,
	}
}

// Key returns the key of this item.
func (item *RegommendItem) Key() interface{} {
	// immutable
	return item.key
}

// Data returns the value of this item.
func (item *RegommendItem) Data() map[interface{}]float64 {
	// immutable
	return item.data
}
