package trip

import (
	"fmt"
	"slices"
)

type Trips struct {
	storage    map[string]Trip
	arrStorage []Trip
	// ban-list
}

func NewTrips() *Trips {
	return &Trips{
		storage: make(map[string]Trip),
	}
}

func (ts *Trips) AddTrip(t Trip) { // add ban service?
	v, ok := ts.storage[t.CacheID]
	if ok {
		if v.GetPrice() > t.GetPrice() {
			ts.Set(t.CacheID, t)
		}
	} else {
		ts.Set(t.CacheID, t)
	}
}

func (ts *Trips) RemoveTrip(t *Trip) {
	delete(ts.storage, t.CacheID)
	ts.updateArr()
}

func (ts *Trips) Merge(tsm *Trips) {
	for k := range tsm.storage {
		ts.AddTrip(tsm.storage[k])
	}
}

func (ts *Trips) Set(key string, t Trip) {
	ts.storage[t.CacheID] = t
	ts.updateArr()
}

func (ts *Trips) Get(key string) (Trip, error) {
	v, ok := ts.storage[key]
	if !ok {
		return Trip{}, fmt.Errorf("invalid key")
	}
	return v, nil
}

func (ts *Trips) GetFirst() Trip {
	if len(ts.arrStorage) == 0 {
		return Trip{}
	}
	return ts.arrStorage[0]
}

func (ts *Trips) ToArray() []Trip {
	return ts.arrStorage
}

func (ts *Trips) Count() int {
	return len(ts.arrStorage)
}

func (ts *Trips) IsEmpty() bool {
	return len(ts.arrStorage) == 0
}

func (ts *Trips) Contains(key string) bool {
	_, ok := ts.storage[key]
	return ok
}

func (ts *Trips) SortByPrice() {
	slices.SortFunc(ts.arrStorage, func(a, b Trip) int {
		if a.GetPrice() > b.GetPrice() {
			return 1
		} else if a.GetPrice() < b.GetPrice() {
			return -1
		} else {
			return 0
		}
	})
}

func (ts *Trips) SortByDirection() {
	slices.SortFunc(ts.arrStorage, func(a, b Trip) int {
		return a.Metadata.RouteDuration - b.Metadata.RouteDuration
	})
}

func (ts *Trips) updateArr() {
	ts.arrStorage = ts.arrStorage[:0]
	for _, v := range ts.storage {
		ts.arrStorage = append(ts.arrStorage, v)
	}
}

// func (ts *Trips) getCacheID(t *Trip) string {
// 	return t.SRO.GetToken()
// }
