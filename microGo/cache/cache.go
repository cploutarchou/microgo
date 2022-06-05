package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Exist(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Delete(string) error
	DeleteIfMatch(string) error
	Empty() error
}

type RedisCache struct {
	Connect *redis.Pool
	Prefix  string
}

type Entry map[string]interface{}

//Exist :  Check if the key exists
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

//encode : Encode a string/s  of type entry
func encode(item Entry) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(item)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

//decode : Decode  a string/s  of type entry
func decode(str string) (Entry, error) {
	item := Entry{}
	b := bytes.Buffer{}
	b.Write([]byte(str))
	d := gob.NewDecoder(&b)
	err := d.Decode(&item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

//Get : Return Key values from Redis if it exists.
func (c *RedisCache) Get(key string) (interface{}, error) {
	_key := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Connect.Get()
	defer conn.Close()
	cacheEntry, err := redis.Bytes(conn.Do("GET", _key))
	if err != nil {
		return nil, err
	}
	decoded, err := decode(string(cacheEntry))
	if err != nil {
		return nil, err
	}
	item := decoded[key]
	return item, nil
}

//Set : Set a key value in Redis
func (c *RedisCache) Set(key string, value interface{}, expires ...int) error {
	_key := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Connect.Get()
	defer conn.Close()

	entry := Entry{}
	entry[_key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		_, err = conn.Do("SEMTEX", _key, expires[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		_, err = conn.Do("SET", _key, string(encoded))
		if err != nil {
			return err
		}
	}
	return nil
}

//Delete : Delete a key value in Redis.
func (c *RedisCache) Delete(key string) error {
	_key := fmt.Sprintf("%s:%s", c.Prefix, key)
	conn := c.Connect.Get()
	defer conn.Close()
	_, err := conn.Do("DELETE", _key)

	if err != nil {
		return err
	}
	return nil
}

//DeleteIfMatch : Delete a key values where match the with the key value
func (c *RedisCache) DeleteIfMatch(key string) error {
	_key := fmt.Sprintf("%s:%s", c.Prefix, key)
	keys, err := c.getKeys(_key)
	if err != nil {
		return err
	}
	for _, i := range keys {
		err := c.Delete(i)
		if err != nil {
			return err
		}
	}
	return nil
}

//Empty : Delete all entries from redis.
func (c *RedisCache) Empty() error {
	_keys := fmt.Sprintf("%s:", c.Prefix)
	conn := c.Connect.Get()
	defer conn.Close()
	keys, err := c.getKeys(_keys)
	if err != nil {
		return err
	}
	for _, i := range keys {
		err := c.Delete(i)
		if err != nil {
			return err
		}
	}
	return nil
}

//getKeys : Return all keys that match to the pattern.
func (c *RedisCache) getKeys(pattern string) ([]string, error) {
	conn := c.Connect.Get()
	defer conn.Close()
	iter := 0
	var keys []string
	for {
		arr, err := redis.Values(conn.Do("SCAN", iter, "MATCH", fmt.Sprintf("%s*", pattern)))
		if err != nil {
			return keys, err
		}
		iter, _ = redis.Int(arr[0], nil)
		k, _ := redis.Strings(arr[1], nil)
		keys = append(keys, k...)

		if iter == 0 {
			break
		}

	}

	return keys, nil
}
