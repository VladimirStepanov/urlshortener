package store

import (
	"fmt"
	"time"
)

var (
	//ErrExpired ...
	ErrExpired = fmt.Errorf("Date is expired")
	//ErrItemNotFound ...
	ErrItemNotFound = fmt.Errorf("Item not found")
)

//BaseItem ...
type BaseItem struct {
	URL    string `redis:"url" json:"url"`
	Visits uint64 `redis:"visits" json:"visits"`
	Expire string `redis:"expire" json:"expire"`
	Once   bool   `redis:"once" json:"once"`
}

//Item ...
type Item struct {
	ID uint64 `redis:"id"`
	BaseItem
}

//Storage ...
type Storage interface {
	Save(url string, expire time.Time, once bool) (uint64, error)
	Load(id uint64) (*Item, error)
	Remove(id uint64) (*Item, error)
	Close() error
	IncVisits(id uint64) error
}
