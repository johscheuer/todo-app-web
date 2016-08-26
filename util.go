package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"


	"gopkg.in/redis.v4"

	"github.com/prometheus/client_golang/prometheus"
)

func createRedisClient(addr, password string) *(redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})
}

func checkConnection(connection string, password string, okString string) (string) {
	if _, err := createRedisClient(connection, password).Ping().Result(); err != nil {
		return err.Error()
	} else {
		return okString
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

var redisMastersTotal = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "todoapp_redis_masters_total",
		Help: "Total count of available redis masters",
	},
	[]string{"instance", "version"},
)

var redisMastersHealthyTotal = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "todoapp_redis_masters_healthy_total",
		Help: "Total count of healthy redis masters",
	},
	[]string{"instance", "version"},
)

var redisSlavesTotal = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "todoapp_redis_slaves_total",
		Help: "Total count of available redis slaves",
	},
	[]string{"instance", "version"},
)

var redisSlavesHealthyTotal = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "todoapp_redis_slaves_healthy_total",
		Help: "Total count of healthy redis slaves",
	},
	[]string{"instance", "version"},
)

func registerMetrics() {
	prometheus.MustRegister(redisMastersTotal)
	prometheus.MustRegister(redisMastersHealthyTotal)
	prometheus.MustRegister(redisSlavesTotal)
	prometheus.MustRegister(redisSlavesHealthyTotal)
}

func getHealthStatus() (map[string]string) {
	okString := "ok"
	result := map[string]string{"self": okString}
	hostname, err := os.Hostname()

	if err != nil { //TODO we'll just ignore any errors :)
		hostname = "UNKNOWN"
	}

	result["redis-master"] = checkConnection(masterConnection, masterPassword, okString)
	redisMastersTotal.WithLabelValues(hostname, APP_VERSION).Set(1.0)

	if result["redis-master"] == okString {
		redisMastersHealthyTotal.WithLabelValues(hostname, APP_VERSION).Set(1.0)
	} else {
		redisMastersHealthyTotal.WithLabelValues(hostname, APP_VERSION).Set(0.0)
	}


	result["redis-slave"] = checkConnection(slaveConnection, slavePassword, okString)
	redisSlavesTotal.WithLabelValues(hostname, APP_VERSION).Set(1.0)

	if result["redis-slave"] == okString {
		redisSlavesHealthyTotal.WithLabelValues(hostname, APP_VERSION).Set(1.0)
	} else {
		redisSlavesHealthyTotal.WithLabelValues(hostname, APP_VERSION).Set(0.0)
	}

	return result
}
