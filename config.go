package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type TodoAppConfig struct {
	HealthCheckTime int
	DBDriver        string
	DBConfig        map[string]string
}

func readConfig(configFile string) *TodoAppConfig {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", string(file))

	config := &TodoAppConfig{}
	json.Unmarshal(file, config)
	return config
}

/*
	flag.StringVar(&masterConnection, "master", "redis-master:6379", "The connection string to the Redis master as <hostname/ip>:<port>")
	flag.StringVar(&slaveConnection, "slave", "redis-slave:6379", "The connection string to the Redis slave as <hostname/ip>:<port>")
	flag.StringVar(&masterPassword, "master-password", "", "The password used to connect to the master")
	flag.StringVar(&slavePassword, "slave-password", "", "The password used to connect to the slave")
*/
