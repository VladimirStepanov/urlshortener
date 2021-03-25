package redis

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/VladimirStepanov/urlshortener/pkg/config"
	"github.com/VladimirStepanov/urlshortener/pkg/store"
	"github.com/gomodule/redigo/redis"
)

//RedisStorage ...
type RedisStorage struct {
	pool *redis.Pool
}

//New - constructor for RedisStorage
func New(c *config.Config) store.Storage {
	s := &RedisStorage{
		&redis.Pool{
			MaxIdle:     10,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort))
			},
		},
	}
	return s
}

// IsExists - check key in store
func (rs *RedisStorage) isExists(id uint64, conn redis.Conn) (bool, error) {

	exists, err := redis.Bool(conn.Do("EXISTS", fmt.Sprintf("url:%d", id)))

	if err != nil {
		return false, err
	}

	if !exists {
		return false, nil
	}

	return true, nil
}

// Save data to redis store
func (rs *RedisStorage) Save(url string, expire time.Time, once bool) (uint64, error) {
	now := time.Now()
	var id uint64

	if expire.Before(now.UTC()) {
		return 0, store.ErrExpired
	}
	conn := rs.pool.Get()
	defer conn.Close()

	for {
		id = rand.Uint64()
		exists, err := rs.isExists(id, conn)
		if err != nil {
			return 0, err
		}

		if !exists {
			break
		}
	}

	key := fmt.Sprintf("url:%d", id)

	_, err := conn.Do(
		"HMSET", key,
		"url", url,
		"visits", 0,
		"once", once,
		"expire", expire.Format("2.1.2006 15:4:5"),
	)

	if err != nil {
		return 0, err
	}

	_, err = conn.Do(
		"EXPIREAT", key, expire.Unix(),
	)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (rs *RedisStorage) getItem(id uint64, conn redis.Conn) (*store.Item, error) {
	values, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("url:%d", id)))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, store.ErrItemNotFound
	}

	res := &store.Item{ID: id}

	err = redis.ScanStruct(values, res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//Load - get Item from Redis store
func (rs *RedisStorage) Load(id uint64) (*store.Item, error) {
	conn := rs.pool.Get()
	defer conn.Close()

	return rs.getItem(id, conn)

}

//Remove - remove item from redis
func (rs *RedisStorage) Remove(id uint64) (*store.Item, error) {
	conn := rs.pool.Get()
	defer conn.Close()

	res, err := rs.getItem(id, conn)

	if err != nil {
		return nil, err
	}

	_, err = conn.Do(
		"DEL", fmt.Sprintf("url:%d", id),
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//IncVisits ...
func (rs *RedisStorage) IncVisits(id uint64) error {
	conn := rs.pool.Get()
	defer conn.Close()

	exists, err := rs.isExists(id, conn)

	if err != nil {
		return err
	}

	if !exists {
		return store.ErrItemNotFound
	}

	err = conn.Send("HINCRBY", fmt.Sprintf("url:%d", id), "visits", 1)

	return err
}

//Close - close pool
func (rs *RedisStorage) Close() error {
	return rs.pool.Close()
}
