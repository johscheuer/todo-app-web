package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type TodoAppConfig struct {
	HealthCheckTime int
	DBDriver        string
	DBConfig        map[string]string
}

func readConfig(configFile string) (*TodoAppConfig, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &TodoAppConfig{
			DBDriver: "redis",
			DBConfig: map[string]string{},
		}, err
	}
	config := &TodoAppConfig{}
	json.Unmarshal(file, config)

	if config.DBDriver == "" {
		log.Println("Use redis as default")
		config.DBDriver = "redis"
	}

	if config.DBConfig == nil {
		config.DBConfig = map[string]string{}
	}

	return config, nil
}
