package redis

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/VladimirStepanov/urlshortener/pkg/config"
	"github.com/VladimirStepanov/urlshortener/pkg/store"
	"github.com/gomodule/redigo/redis"
)

var (
	defaultConf = &config.Config{RedisHost: "127.0.0.1", RedisPort: "6379"}
	defaultItem = &store.Item{ID: 1, BaseItem: store.BaseItem{URL: "https://vk.com", Visits: 100, Expire: "10.1.2380 1:0:0", Once: true}}
)

func addKey(rs *RedisStorage, item *store.Item) error {
	pool := rs.pool.Get()
	defer pool.Close()

	_, err := pool.Do(
		"HMSET", fmt.Sprintf("url:%d", item.ID),
		"url", item.URL,
		"visits", item.Visits,
		"once", item.Once,
		"expire", item.Expire,
	)

	if err != nil {
		return err
	}

	return nil
}

func removeKey(rs *RedisStorage, ID uint64) error {
	pool := rs.pool.Get()
	defer pool.Close()

	_, err := pool.Do("DEL", fmt.Sprintf("url:%d", ID))

	if err != nil {
		return err
	}

	return nil
}

func NewTestRedisStore(c *config.Config) *RedisStorage {
	s := &RedisStorage{
		&redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort))
			},
		},
	}

	addKey(s, defaultItem)

	return s
}

func CloseTestRedisStore(rs *RedisStorage) {
	removeKey(rs, defaultItem.ID)
}

func TestIsExistsRedisStorage(t *testing.T) {
	rs := NewTestRedisStore(defaultConf)

	defer CloseTestRedisStore(rs)

	tests := map[string]struct {
		id     uint64
		result bool
	}{
		"Key is found":     {1, true},
		"Key is not found": {666, false},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			res, _ := rs.isExists(test.id, rs.pool.Get())

			if res != test.result {
				t.Fatalf("Expected %v, but got %v", test.result, res)
			}
		})
	}
}

func TestVariousSaveIDRedisStorage(t *testing.T) {
	rs := NewTestRedisStore(defaultConf)
	defer CloseTestRedisStore(rs)

	now := time.Now()

	firstID, err := rs.Save("https://vk.com", now.AddDate(1, 0, 0), true)

	if err != nil {
		t.Fatalf("For first Save got error: %v", err)
	}

	defer removeKey(rs, firstID)

	secondID, err := rs.Save("https://vk.com", now.AddDate(2, 0, 0), true)

	if err != nil {
		t.Fatalf("For first Save got error: %v", err)
	}

	defer removeKey(rs, secondID)

	if firstID == secondID {
		t.Fatalf("Error! %d == %d", firstID, secondID)
	}
}

func TestSaveRedisStorage(t *testing.T) {
	now := time.Now()

	tests := map[string]struct {
		url     string
		expire  time.Time
		conf    *config.Config
		isError bool
		err     error
		once    bool
	}{
		"Expire error": {
			expire:  now.AddDate(-1, 0, 0),
			conf:    defaultConf,
			isError: true,
			err:     store.ErrExpired,
		},
		"Connection refused": {
			expire:  now.AddDate(1, 0, 0),
			conf:    &config.Config{RedisHost: "127.0.1.1", RedisPort: "6379"},
			isError: true,
			err:     fmt.Errorf("dial tcp 127.0.1.1:6379: connect: connection refused"),
		},
		"Succes add[once=true]": {
			expire: now.AddDate(1, 0, 0),
			conf:   defaultConf,
			once:   true,
		},
		"Succes add[once=false]": {
			expire: now.AddDate(1, 0, 0),
			conf:   defaultConf,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			rs := NewTestRedisStore(test.conf)
			defer CloseTestRedisStore(rs)
			id, err := rs.Save(test.url, test.expire, test.once)

			if err == nil {
				defer removeKey(rs, id)
			}

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
func LRWrapper(t *testing.T, tests map[string]LRTestCase, rs *RedisStorage, testFunc func(uint64) (*store.Item, error)) {
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

func TestRemoveRedisStorage(t *testing.T) {
	rs := NewTestRedisStore(defaultConf)
	defer CloseTestRedisStore(rs)
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

func TestLoadRedisStorage(t *testing.T) {
	rs := NewTestRedisStore(defaultConf)
	defer CloseTestRedisStore(rs)
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

func TestIncVisitsRedisStorage(t *testing.T) {
	rs := NewTestRedisStore(defaultConf)
	defer CloseTestRedisStore(rs)

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
