package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/johscheuer/todo-app-web/tododb"
	"github.com/mcuadros/go-gin-prometheus"
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
		os.Exit(0)
	}

	config, err := readConfig(*configFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	gin.SetMode(config.ReleaseMode)
	if strings.ToLower(config.DBDriver) == "redis" {
		database = tododb.NewRedisDB(config.DBConfig, appVersion)
	} else {
		log.Printf("Datebase: %s is not supported", config.DBDriver)
		os.Exit(1)
	}

	p := ginprometheus.NewPrometheus("gin")
	database.RegisterMetrics()

	router := gin.Default()

	p.Use(router)
	router.GET("/read/todo", readTodoHandler)
	router.GET("/insert/todo/:value", insertTodoHandler)
	router.GET("/delete/todo/:value", deleteTodoHandler)
	router.GET("/health", healthCheckHandler)
	router.GET("/whoami", whoAmIHandler)
	router.GET("/version", versionHandler)

	router.Use(static.Serve("/", static.LocalFile("./public", true)))
	router.Run(":3000")
}
