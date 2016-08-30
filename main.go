package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/negroni"
)

var (
	slaveConnection  string
	slavePassword    string
	masterConnection string
	masterPassword   string
	appVersion       string
	showVersion      bool
	healthCheckTime  int
)

func init() {
	// Metrics have to be registered to be exposed:
	registerMetrics()
}

func main() {
	flag.StringVar(&masterConnection, "master", "redis-master:6379", "The connection string to the Redis master as <hostname/ip>:<port>")
	flag.StringVar(&slaveConnection, "slave", "redis-slave:6379", "The connection string to the Redis slave as <hostname/ip>:<port>")
	flag.StringVar(&masterPassword, "master-password", "", "The password used to connect to the master")
	flag.StringVar(&slavePassword, "slave-password", "", "The password used to connect to the slave")
	flag.IntVar(&healthCheckTime, "health-check", 15, "Period to check all connections")
	flag.BoolVar(&showVersion, "version", false, "Shows the version")
	flag.Parse()

	if showVersion {
		log.Printf("Version: %s\n", appVersion)
		return
	}
	// Iniitialize metrics
	healthCheckTimer := time.NewTicker(time.Duration(healthCheckTime) * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-healthCheckTimer.C:
				log.Println("Called Health check")
				getHealthStatus()
			case <-quit:
				healthCheckTimer.Stop()
				return
			}
		}
	}()

	r := mux.NewRouter()
	r.Path("/read/{key}").Methods("GET").HandlerFunc(readTodoHandler)
	r.Path("/insert/{key}/{value}").Methods("GET").HandlerFunc(insertTodoHandler)
	r.Path("/delete/{key}/{value}").Methods("GET").HandlerFunc(deleteTodoHandler)
	r.Path("/health").Methods("GET").HandlerFunc(healthCheckHandler)
	r.Path("/metrics").Methods("GET").Handler(prometheus.Handler())
	r.Path("/whoami").Methods("GET").HandlerFunc(whoAmIHandler)
	r.Path("/version").Methods("GET").HandlerFunc(versionHandler)

	n := negroni.Classic()
	n.UseHandler(r)
	http.ListenAndServe(":3000", n)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)
	close(quit)
}
