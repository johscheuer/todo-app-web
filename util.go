package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

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

func getAllAddresses(ifaces []net.Interface) ([]string, error) {
	var addresses []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			return addresses, err
		}

		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}
	}

	return addresses, nil
}

func generateJSONResponse(rw http.ResponseWriter, toMarshal interface{}) {
	responseJSON, err := json.MarshalIndent(toMarshal, "", "  ")
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}
	rw.Write(responseJSON)
}
