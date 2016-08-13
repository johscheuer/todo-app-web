package main

import (
	"gopkg.in/redis.v4"
)

func createRedisClient(addr string) *(redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func checkConnection(result map[string]string, key, connection, okString string) {
  if _, err := createRedisClient(connection).Ping().Result(); err != nil {
    result[key] = err.Error()
  } else {
    result[key] = okString
  }
}
