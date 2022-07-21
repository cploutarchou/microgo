package cache

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	Exists(string) (bool, error)
	Get(string) (interface{}, error)
	Set(string, interface{}, ...int) error
	Delete(string) error
	DeleteIfMatch(string) error
	Clean() error
}

type RedisCache struct {
	Connection *redis.Pool
	Prefix     string
}

type Entry map[string]interface{}

//Exists :  Check if the key exists
func (c *RedisCache) Exists(str string) (bool, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Connection.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
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
func (c *RedisCache) Get(str string) (interface{}, error) {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Connection.Get()
	defer conn.Close()

	cacheEntry, err := redis.Bytes(conn.Do("GET", key))
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
func (c *RedisCache) Set(str string, value interface{}, expires ...int) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Connection.Get()
	defer conn.Close()

	entry := Entry{}
	entry[key] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}

	if len(expires) > 0 {
		_, err := conn.Do("SETEX", key, expires[0], string(encoded))
		if err != nil {
			return err
		}
	} else {
		_, err := conn.Do("SET", key, string(encoded))
		if err != nil {
			return err
		}
	}

	return nil
}

//Delete : Delete a key value in Redis.
func (c *RedisCache) Delete(str string) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Connection.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}

//DeleteIfMatch : Delete a key values where match the with the key value
func (c *RedisCache) DeleteIfMatch(str string) error {
	key := fmt.Sprintf("%s:%s", c.Prefix, str)
	conn := c.Connection.Get()
	defer conn.Close()

	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err := conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

//Clean : Delete all entries from redis.
func (c *RedisCache) Clean() error {
	key := fmt.Sprintf("%s:", c.Prefix)
	conn := c.Connection.Get()
	defer conn.Close()

	keys, err := c.getKeys(key)
	if err != nil {
		return err
	}

	for _, x := range keys {
		_, err := conn.Do("DEL", x)
		if err != nil {
			return err
		}
	}

	return nil
}

//getKeys : Return all keys that match to the pattern.
func (c *RedisCache) getKeys(pattern string) ([]string, error) {
	conn := c.Connection.Get()
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
