package teststore

import (
	"math/rand"
	"time"

	"github.com/VladimirStepanov/urlshortener/pkg/store"
)

//TestStorage ...
type TestStorage struct {
	items map[uint64]*store.Item
}

//New ...
func New(items map[uint64]*store.Item) *TestStorage {

	return &TestStorage{items: items}
}

func (rs *TestStorage) isExists(id uint64) (bool, error) {

	_, ok := rs.items[id]

	return ok, nil
}

// Save ...
func (rs *TestStorage) Save(url string, expire time.Time, once bool) (uint64, error) {
	now := time.Now()
	var id uint64

	if expire.Before(now.UTC()) {
		return 0, store.ErrExpired
	}
	for {
		id = rand.Uint64()
		exists, err := rs.isExists(id)
		if err != nil {
			return 0, err
		}

		if !exists {
			break
		}
	}

	rs.items[id] = &store.Item{ID: id, BaseItem: store.BaseItem{URL: url, Visits: 0, Expire: expire.Format("2.1.2006 15:4:5"), Once: once}}

	return id, nil
}

func (rs *TestStorage) getItem(id uint64) (*store.Item, error) {
	item, ok := rs.items[id]

	if !ok {
		return nil, store.ErrItemNotFound
	}

	t, err := time.Parse("2.1.2006 15:4:5", item.Expire)

	if err != nil {
		return nil, err
	}

	if t.Before(time.Now()) {
		// delete(rs.items, id)
		return nil, store.ErrItemNotFound
	}

	return item, nil
}

//Load ...
func (rs *TestStorage) Load(id uint64) (*store.Item, error) {
	return rs.getItem(id)

}

//Remove ...
func (rs *TestStorage) Remove(id uint64) (*store.Item, error) {

	res, err := rs.getItem(id)

	if err != nil {
		return nil, err
	}

	delete(rs.items, id)

	return res, nil
}

//IncVisits ...
func (rs *TestStorage) IncVisits(id uint64) error {

	item, ok := rs.items[id]
	if !ok {
		return store.ErrItemNotFound
	}

	item.Visits++

	return nil
}

//Close - close pool
func (rs *TestStorage) Close() error {
	return nil
}
