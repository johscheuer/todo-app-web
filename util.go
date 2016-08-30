package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	"gopkg.in/redis.v4"

	"github.com/prometheus/client_golang/prometheus"
)

const okString string = "ok"

type checkConnectionResult struct {
	results map[string]string
	total   int
	healthy int
	name    string
}

func newCheckConnectionResult(name string) *checkConnectionResult {
	return &checkConnectionResult{
		results: map[string]string{},
		total:   0,
		healthy: 0,
		name:    name,
	}
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

	var wg sync.WaitGroup
	results := make(chan *checkConnectionResult, 2)
	wg.Add(2)
	go func() {
		results <- checkConnections("redis-master", hostname, masterConnection)
		wg.Done()
	}()

	go func() {
		results <- checkConnections("redis-slave", hostname, slaveConnection)
		wg.Done()
	}()
	wg.Wait()

	close(results)

	// Merge Results
	for res := range results {
		if res.name == "redis-master" {
			redisMastersTotal.WithLabelValues(hostname, appVersion).Set(float64(res.total))
			redisMastersHealthyTotal.WithLabelValues(hostname, appVersion).Set(float64(res.healthy))
		}
		if res.name == "redis-slave" {
			redisSlavesTotal.WithLabelValues(hostname, appVersion).Set(float64(res.total))
			redisSlavesHealthyTotal.WithLabelValues(hostname, appVersion).Set(float64(res.healthy))
		}

		for k, v := range res.results {
			result[k] = v
		}
	}

	return result
}

func checkConnections(name, hostname, connection string) *checkConnectionResult {
	res := newCheckConnectionResult(name)
	masterConnections, err := getAllConnections(connection)
	if err != nil {
		log.Println(err)
		// Simple fallback
		masterConnections = []string{connection}
	}

	for index, connection := range masterConnections {
		conName := fmt.Sprintf("%s-%d", name, index)
		res.results[conName] = checkConnection(connection, masterPassword)
		res.total++

		if res.results[conName] == okString {
			res.healthy++
		}
	}

	return res
}

//TODO add function with SRV lookup
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
