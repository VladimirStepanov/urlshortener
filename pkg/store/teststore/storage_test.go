package teststore

import (
	"reflect"
	"testing"
	"time"

	"github.com/VladimirStepanov/urlshortener/pkg/store"
)

var (
	defaultItem = &store.Item{ID: 1, BaseItem: store.BaseItem{URL: "https://vk.com", Visits: 100, Expire: "10.1.2380 1:0:0", Once: true}}
)

func GetTestStore() *TestStorage {
	return &TestStorage{
		items: map[uint64]*store.Item{
			1: defaultItem,
		},
	}
}

func TestIsExistsTestStorage(t *testing.T) {
	rs := GetTestStore()
	tests := map[string]struct {
		id     uint64
		result bool
	}{
		"Key is found":     {1, true},
		"Key is not found": {666, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, _ := rs.isExists(test.id)

			if res != test.result {
				t.Fatalf("Expected %v, but got %v", test.result, res)
			}
		})
	}
}

func TestSaveTestStorage(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		url     string
		expire  time.Time
		isError bool
		err     error
		once    bool
	}{
		"Expire error": {
			expire:  now.AddDate(-1, 0, 0),
			isError: true,
			err:     store.ErrExpired,
		},
		"Succes add[once=true]": {
			expire: now.AddDate(1, 0, 0),
			once:   true,
		},
		"Succes add[once=false]": {
			expire: now.AddDate(1, 0, 0),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rs := GetTestStore()
			_, err := rs.Save(test.url, test.expire, test.once)

			if test.isError && err == nil {
				t.Fatalf("Expected error: %v, but got nil", test.err)
			}

			if test.isError && err != nil {
				if test.err.Error() != err.Error() {
					t.Fatalf("Expected errror: %v, but got: %v", test.err, err)
				}
			}

			if !test.isError && err != nil {
				t.Fatalf("isError false, but got %v", err)
			}
		})
	}

}

//LRTestCase - common test case type for testing Load and Remove methods
type LRTestCase struct {
	id      uint64
	item    *store.Item
	isError bool
	err     error
}

//Wrapper for resting rs.Remove and rs.Load methods of RedisStorage
func LRWrapper(t *testing.T, tests map[string]LRTestCase, rs *TestStorage, testFunc func(uint64) (*store.Item, error)) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			item, err := testFunc(test.id)

			if test.isError && err == nil {
				t.Fatalf("Expected error: %v, but got nil", test.err)
			}

			if test.isError && err != nil {
				if test.err.Error() != err.Error() {
					t.Fatalf("Expected errror: %v, but got: %v", test.err, err)
				}
			} else if !test.isError && err == nil {
				if !reflect.DeepEqual(*item, *test.item) {
					t.Fatalf("Expected: %v, but got: %v", test.item, item)
				}
			}
		})
	}
}

func TestLoadTestStorage(t *testing.T) {
	rs := GetTestStore()
	tests := map[string]LRTestCase{
		"Success load": {
			id:   defaultItem.ID,
			item: defaultItem,
		},
		"Error: item not found": {
			id:      defaultItem.ID + 1,
			isError: true,
			err:     store.ErrItemNotFound,
		},
	}

	LRWrapper(t, tests, rs, rs.Load)
}

func TestRemoveTestStorage(t *testing.T) {
	rs := GetTestStore()
	tests := map[string]LRTestCase{
		"Success delete": {
			id:   defaultItem.ID,
			item: defaultItem,
		},
		"Delete error: item not found": {
			id:      defaultItem.ID + 6000,
			isError: true,
			err:     store.ErrItemNotFound,
		},
	}

	LRWrapper(t, tests, rs, rs.Remove)
}

func TestIncVisitsTestStorage(t *testing.T) {
	rs := GetTestStore()
	defer rs.Close()

	tests := map[string]struct {
		id  uint64
		err error
	}{
		"Success increment": {
			id: defaultItem.ID,
		},
		"Error: item not found": {
			id:  defaultItem.ID + 1,
			err: store.ErrItemNotFound,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := rs.IncVisits(tc.id)

			if err != tc.err {
				t.Fatalf("Expected errror: %v, but got: %v", tc.err, err)
			}
		})
	}
}
