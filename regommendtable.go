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
	_ "fmt"
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

// Count returns how many items are currently stored in the engine.
func (table *RegommendTable) Count() int {
	table.RLock()
	defer table.RUnlock()
	return len(table.items)
}

// Configures a data-loader callback, which will be called when trying
// to use access a non-existing key.
func (table *RegommendTable) SetDataLoader(f func(interface{}) *RegommendItem) {
	table.Lock()
	defer table.Unlock()
	table.loadData = f
}

// Configures a callback, which will be called every time a new item
// is added to the engine.
func (table *RegommendTable) SetAddedItemCallback(f func(*RegommendItem)) {
	table.Lock()
	defer table.Unlock()
	table.addedItem = f
}

// Configures a callback, which will be called every time an item
// is about to be removed from the engine.
func (table *RegommendTable) SetAboutToDeleteItemCallback(f func(*RegommendItem)) {
	table.Lock()
	defer table.Unlock()
	table.aboutToDeleteItem = f
}

// Sets the logger to be used by this engine table.
func (table *RegommendTable) SetLogger(logger *log.Logger) {
	table.Lock()
	defer table.Unlock()
	table.logger = logger
}

// Add adds a key/value pair to the engine.
// Parameter key is the item's engine-key.
// Parameter data is the item's value.
func (table *RegommendTable) Add(key interface{}, data map[interface{}]float64) *RegommendItem {
	item := CreateRegommendItem(key, data)

	// Add item to engine.
	table.Lock()
	table.items[key] = &item

	// engine values so we don't keep blocking the mutex.
	addedItem := table.addedItem
	table.Unlock()

	// Trigger callback after adding an item to engine.
	if addedItem != nil {
		addedItem(&item)
	}

	return &item
}

// Delete an item from the engine.
func (table *RegommendTable) Delete(key interface{}) (*RegommendItem, error) {
	table.RLock()
	r, ok := table.items[key]
	if !ok {
		table.RUnlock()
		return nil, errors.New("Key not found in engine")
	}

	// engine value so we don't keep blocking the mutex.
	aboutToDeleteItem := table.aboutToDeleteItem
	table.RUnlock()

	// Trigger callbacks before deleting an item from engine.
	if aboutToDeleteItem != nil {
		aboutToDeleteItem(r)
	}

	r.RLock()
	defer r.RUnlock()

	table.Lock()
	defer table.Unlock()
	delete(table.items, key)

	return r, nil
}

// Test whether an item exists in the engine. Unlike the Value method
// Exists neither tries to fetch data via the loadData callback nor
// does it keep the item alive in the engine.
func (table *RegommendTable) Exists(key interface{}) bool {
	table.RLock()
	defer table.RUnlock()
	_, ok := table.items[key]

	return ok
}

// Value gets an item from the engine and mark it to be kept alive.
func (table *RegommendTable) Value(key interface{}) (*RegommendItem, error) {
	table.RLock()
	r, ok := table.items[key]
	loadData := table.loadData
	table.RUnlock()

	if ok {
		return r, nil
	}

	// Item doesn't exist in engine. Try and fetch it with a data-loader.
	if loadData != nil {
		item := loadData(key)
		if item != nil {
			table.Add(key, item.data)
			return item, nil
		}

		return nil, errors.New("Key not found and could not be loaded into engine")
	}

	return nil, errors.New("Key not found in engine")
}

// Delete all items from engine.
func (table *RegommendTable) Flush() {
	table.Lock()
	defer table.Unlock()

	table.log("Flushing table", table.name)

	table.items = make(map[interface{}]*RegommendItem)
}

type DistancePair struct {
	Key interface{}
	Distance float64
}
type DistancePairList []DistancePair

func (p DistancePairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p DistancePairList) Len() int { return len(p) }
func (p DistancePairList) Less(i, j int) bool { return p[i].Distance > p[j].Distance }

func (table *RegommendTable) Recommend(key interface{}) (DistancePairList, error) {
	dists, err := table.Neighbors(key)
	if err != nil {
		return dists, err
	}
	sitem, err := table.Value(key)
	if err != nil {
		return dists, err
	}
	smap := sitem.Data()

	totalDistance := 0.0
	for _, v := range dists {
		//fmt.Println("Comparing to", v.Key, "-", v.Distance)
		totalDistance += v.Distance
	}

	recs := make(map[interface{}]float64)
	for _, v := range dists {
		weight := v.Distance / totalDistance
		if weight <= 0 {
			continue
		}
		if weight > 1 {
			weight = 1
		}

		ditem, _ := table.Value(v.Key)
		recMap := ditem.Data()
		for key, x := range recMap {
			_, ok := smap[key]
			if ok {
				// key already knows this item, don't recommend it
				continue
			}

			//fmt.Println("Adding to recs:", key)
			score, ok := recs[key]
			if ok {
				recs[key] = score + x * weight
			} else {
				recs[key] = x * weight
			}
		}
	}

	recsList := make(DistancePairList, len(recs))
	i := 0
	for key, score := range recs {
		recsList[i] = DistancePair{
			Key: key,
			Distance: score,
		}
		i++
	}
	sort.Sort(recsList)

	return recsList, nil
}

func (table *RegommendTable) Neighbors(key interface{}) (DistancePairList, error) {
	dists := DistancePairList{}

	sitem, err := table.Value(key)
	if err != nil {
		return dists, err
	}
	smap := sitem.Data()

	table.RLock()
	defer table.RUnlock()
	for k, ditem := range table.items {
		if err != nil {
			continue
		}
		if k == key {
			continue
		}

		//fmt.Println("Analyzing:", k)
		distance := DistancePair{
			Key: k,
			Distance: cosineSim(smap, ditem.Data()),
		}
		//fmt.Println("Distance:", distance.Distance)
		dists = append(dists, distance)
	}
	sort.Sort(dists)

	return dists, nil
}

// Internal logging method for convenience.
func (table *RegommendTable) log(v ...interface{}) {
	if table.logger == nil {
		return
	}

	table.logger.Println(v...)
}
