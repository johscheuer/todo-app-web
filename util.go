package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"gopkg.in/redis.v4"

	"github.com/prometheus/client_golang/prometheus"
)

const okString string = "ok"

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

func createRedisClient(addr, password string) *(redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})
}

func checkConnection(connection string, password string) string {
	if _, err := createRedisClient(connection, password).Ping().Result(); err != nil {
		return err.Error()
	}

	return okString
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

func registerMetrics() {
	prometheus.MustRegister(redisMastersTotal)
	prometheus.MustRegister(redisMastersHealthyTotal)
	prometheus.MustRegister(redisSlavesTotal)
	prometheus.MustRegister(redisSlavesHealthyTotal)
}

func getHealthStatus() map[string]string {
	result := map[string]string{"self": okString}
	hostname, err := os.Hostname()

	if err != nil { //TODO we'll just ignore any errors :)
		hostname = "UNKNOWN"
	}

	//TODO can we generalize this to use only one function?
	checkMasterConnections(result, hostname)
	checkSlaveConnections(result, hostname)

	return result
}

func checkMasterConnections(result map[string]string, hostname string) {
	masterConnections, err := getAllConnections(masterConnection)
	if err != nil {
		log.Println(err)
		// Simple fallback
		masterConnections = []string{masterConnection}
	}

	redisMastersTotal.WithLabelValues(hostname, appVersion).Set(1.0)
	redisMastersHealthyTotal.WithLabelValues(hostname, appVersion).Set(1.0)

	for index, connection := range masterConnections {
		name := fmt.Sprintf("redis-master-%d", index)
		result[name] = checkConnection(connection, masterPassword)
		redisMastersTotal.WithLabelValues(hostname, appVersion).Add(1.0)

		if result[name] == okString {
			redisMastersHealthyTotal.WithLabelValues(hostname, appVersion).Add(1.0)
		}
	}

}

func checkSlaveConnections(result map[string]string, hostname string) {
	slaveConnections, err := getAllConnections(slaveConnection)
	if err != nil {
		log.Println(err)
		// Simple fallback
		slaveConnections = []string{slaveConnection}
	}

	redisSlavesTotal.WithLabelValues(hostname, appVersion).Set(0.0)
	redisSlavesHealthyTotal.WithLabelValues(hostname, appVersion).Set(0.0)

	for index, connection := range slaveConnections {
		name := fmt.Sprintf("redis-slave-%d", index)
		result[name] = checkConnection(connection, slavePassword)
		redisSlavesTotal.WithLabelValues(hostname, appVersion).Add(1.0)

		if result[name] == okString {
			redisSlavesHealthyTotal.WithLabelValues(hostname, appVersion).Add(1.0)
		}
	}
}

func getAllConnections(connection string) ([]string, error) {
	connections := []string{}

	hostname, port, err := net.SplitHostPort(connection)
	if err != nil {
		return connections, err
	}

	hosts, err := net.LookupHost(hostname)
	if err != nil {
		return connections, err
	}

	for _, host := range hosts {
		connections = append(connections, fmt.Sprintf("%s:%s", host, port))
	}

	return connections, nil
}
