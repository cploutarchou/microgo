package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Exist(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Delete(string) error
	EmptyIfMatch(string)
	Empty() error
}

type RedisCache struct {
	Connect *redis.Pool
	Prefix  string
}

type Entry map[string]interface{}

func (c *RedisCache) Exist(key string) (bool, error) {
	_key := fmt.Sprintf("%s:%s", c.Prefix, key)

	conn := c.Connect.Get()
	defer conn.Close()
	ok, err := redis.Bool(conn.Do("EXIST", _key))
	if err != nil {
		return false, err
	}
	return ok, nil
}
