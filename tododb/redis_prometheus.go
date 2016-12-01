package tododb

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

func (redisDB RedisDB) RegisterMetrics() {
	log.Println("Registered Redis Metrics")
	prometheus.MustRegister(redisMastersTotal)
	prometheus.MustRegister(redisMastersHealthyTotal)
	prometheus.MustRegister(redisSlavesTotal)
	prometheus.MustRegister(redisSlavesHealthyTotal)
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

func (redisDB RedisDB) GetHealthStatus() map[string]string {
	result := map[string]string{"self": okString}
	hostname, err := os.Hostname()

	if err != nil { //TODO we'll just ignore any errors :)
		hostname = "UNKNOWN"
	}

	redisMasterHost := getHostnameFromConnection(redisDB.master, "redis-master")
	redisSlaveHost := getHostnameFromConnection(redisDB.slave, "redis-slave")
	var wg sync.WaitGroup
	results := make(chan *checkConnectionResult, 2)
	wg.Add(2)
	go func() {

		results <- checkConnections(redisMasterHost, hostname, redisDB.master, redisDB.masterPassword)
		wg.Done()
	}()

	go func() {

		results <- checkConnections(redisSlaveHost, hostname, redisDB.slave, redisDB.slavePassword)
		wg.Done()
	}()
	wg.Wait()

	close(results)

	// Merge Results
	for res := range results {
		if res.name == redisMasterHost {
			redisMastersTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.total))
			redisMastersHealthyTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.healthy))
		}
		if res.name == redisSlaveHost {
			redisSlavesTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.total))
			redisSlavesHealthyTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.healthy))
		}

		for k, v := range res.results {
			result[k] = v
		}
	}

	return result
}

type checkConnectionResult struct {
	results map[string]string
	total   int
	healthy int
	name    string
}

func getHostnameFromConnection(connection, defaultHost string) string {
	host, _, err := net.SplitHostPort(connection)
	if err != nil {
		host = defaultHost
		fmt.Println(err)
	}

	return host
}

func newCheckConnectionResult(name string) *checkConnectionResult {
	return &checkConnectionResult{
		results: map[string]string{},
		total:   0,
		healthy: 0,
		name:    name,
	}
}

func checkConnection(connection string, password string) string {
	client := createRedisClient(connection, password)
	defer client.Close()
	if _, err := client.Ping().Result(); err != nil {
		return err.Error()
	}

	return okString
}

func checkConnections(name, hostname, connection, password string) *checkConnectionResult {
	res := newCheckConnectionResult(name)
	connections, err := getAllConnections(connection)
	if err != nil {
		log.Println(err)
		// Simple fallback
		connections = []string{connection}
	}

	for index, connection := range connections {
		conName := fmt.Sprintf("%s-%d", name, index)
		res.results[conName] = checkConnection(connection, password)
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
