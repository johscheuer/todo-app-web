package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/gin-gonic/gin"
)

type TodoAppConfig struct {
	HealthCheckTime int
	DBDriver        string
	DBConfig        map[string]string
	ReleaseMode     string
}

func readConfig(configFile string) (*TodoAppConfig, error) {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		return &TodoAppConfig{
			DBDriver:    "redis",
			DBConfig:    map[string]string{},
			ReleaseMode: gin.DebugMode,
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

	if config.ReleaseMode == "" {
		config.ReleaseMode = gin.DebugMode
	}

	return config, err
}
