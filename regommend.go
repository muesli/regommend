/*
 * Simple recommendation engine
 *     Copyright (c) 2014, Christian Muehlhaeuser <muesli@gmail.com>
 *                   
 *   For license see LICENSE.txt
 */

package regommend

import (
	"sync"
)

var (
	tables = make(map[string]*RegommendTable)
	mutex sync.RWMutex
)

// Table returns the existing engine table with given name or creates a new one
// if the table does not exist yet.
func Table(table string) *RegommendTable {
	mutex.RLock()
	t, ok := tables[table]
	mutex.RUnlock()

	if !ok {
		t = &RegommendTable{
			name:  table,
			items: make(map[interface{}]*RegommendItem),
		}

		mutex.Lock()
		tables[table] = t
		mutex.Unlock()
	}

	return t
}
