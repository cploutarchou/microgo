package cache

import (
	"github.com/dgraph-io/badger/v3"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gomodule/redigo/redis"
)

var testRedisCache RedisCache
var testBadgerCache BadgerCache

func TestMain(m *testing.M) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer s.Close()

	pool := redis.Pool{
		MaxIdle:     50,
		MaxActive:   1000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", s.Addr())
		},
	}

	testRedisCache.Connection = &pool
	testRedisCache.Prefix = "test-microGO"

	defer func(Connect *redis.Pool) {
		_ = Connect.Close()
	}(testRedisCache.Connection)

	_ = os.RemoveAll("./testdata/tmp/badger")
	// create a badger database

	if _, err := os.Stat("./testdata/tmp"); os.IsNotExist(err) {
		err := os.Mkdir("./testdata/tmp", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = os.MkdirAll("./testdata/tmp/badger", 0755)
	if err != nil {
		log.Fatal(err)
	}
	db, _ := badger.Open(badger.DefaultOptions("./testdata/tmp/badger"))
	testBadgerCache.Connection = db

	os.Exit(m.Run())
}
