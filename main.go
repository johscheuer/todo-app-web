package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/johscheuer/todo-app-web/tododb"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	appVersion  string
	showVersion bool
	database    tododb.TodoDB
)

func main() {
	configFile := flag.String("config-file", "./default.config", "Path to the configuration file")
	flag.BoolVar(&showVersion, "version", false, "Shows the version")
	flag.Parse()

	if showVersion {
		log.Printf("Version: %s\n", appVersion)
		return
	}

	config, err := readConfig(*configFile)
	if err != nil {
		log.Println(err)
	}

	gin.SetMode(config.ReleaseMode)
	if strings.ToLower(config.DBDriver) == "mysql" {
		database = tododb.NewMySQLDB(config.DBConfig, appVersion)
	} else if strings.ToLower(config.DBDriver) == "redis" {
		database = tododb.NewRedisDB(config.DBConfig, appVersion)
	}

	database.RegisterMetrics()

	// Iniitialize metrics
	quit := make(chan struct{})
	defer close(quit)
	if config.HealthCheckTime > 0 {
		healthCheckTimer := time.NewTicker(time.Duration(config.HealthCheckTime) * time.Second)
		go func() {
			for {
				select {
				case <-healthCheckTimer.C:
					log.Println("Called Health check")
					database.GetHealthStatus()
				case <-quit:
					healthCheckTimer.Stop()
					return
				}
			}
		}()
	}

	router := gin.Default()
	router.GET("/read/todo", readTodoHandler)
	router.GET("/insert/todo/:value", insertTodoHandler)
	router.GET("/delete/todo/:value", deleteTodoHandler)
	router.GET("/health", healthCheckHandler)
	router.GET("/metrics", gin.WrapH(prometheus.Handler()))
	router.GET("/whoami", whoAmIHandler)
	router.GET("/version", versionHandler)

	router.Use(static.Serve("/", static.LocalFile("./public", true)))
	router.Run(":3000")
}
