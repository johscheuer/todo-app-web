package tododb

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	redis "gopkg.in/redis.v5"

	"github.com/prometheus/client_golang/prometheus"
)

type RedisDB struct {
	master         string
	masterPassword string
	slave          string
	slavePassword  string
	appVersion     string
}

const redisKey = "todo"

var _ TodoDB = RedisDB{}

func NewRedisDB(config map[string]string, appVersion string) RedisDB {
	if _, exists := config["master"]; !exists {
		config["master"] = "redis-master:6379"
	}

	if _, exists := config["masterPassword"]; !exists {
		config["masterPassword"] = ""
	}

	if _, exists := config["slave"]; !exists {
		config["slave"] = "redis-slave:6379"
	}

	if _, exists := config["slavePassword"]; !exists {
		config["slavePassword"] = ""
	}

	return RedisDB{
		master:         config["master"],
		masterPassword: config["masterPassword"],
		slave:          config["slave"],
		slavePassword:  config["slavePassword"],
		appVersion:     appVersion,
	}
}

func (redisDB RedisDB) GetAllTodos() ([]string, error) {
	cmd := createRedisClient(redisDB.slave, redisDB.slavePassword).LRange(redisKey, -100, 100)

	// Fallback to read from master
	if cmd.Err() != nil {
		log.Println("Fallback using Redis Master")
		cmd = createRedisClient(redisDB.master, redisDB.masterPassword).LRange(redisKey, -100, 100)
	}
	return cmd.Val(), cmd.Err()
}

func (redisDB RedisDB) SaveTodo(todo string) error {
	cmd := createRedisClient(redisDB.master, redisDB.masterPassword).RPush(redisKey, todo)
	return cmd.Err()
}

func (redisDB RedisDB) DeleteTodo(todo string) error {
	cmd := createRedisClient(redisDB.master, redisDB.masterPassword).LRem(redisKey, 1, todo)
	return cmd.Err()
}

func createRedisClient(addr, password string) *(redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})
}

func (redisDB RedisDB) RegisterMetrics() {
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

	var wg sync.WaitGroup
	results := make(chan *checkConnectionResult, 2)
	wg.Add(2)
	go func() {
		host, _, err := net.SplitHostPort(redisDB.master)
		if err != nil {
			host = "redis-master"
			fmt.Println(err)
		}
		results <- checkConnections(host, hostname, redisDB.master, redisDB.masterPassword)
		wg.Done()
	}()

	go func() {
		host, _, err := net.SplitHostPort(redisDB.slave)
		if err != nil {
			host = "redis-slave"
			fmt.Println(err)
		}
		results <- checkConnections(host, hostname, redisDB.slave, redisDB.slavePassword)
		wg.Done()
	}()
	wg.Wait()

	close(results)

	// Merge Results
	for res := range results {
		if res.name == "redis-master" {
			redisMastersTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.total))
			redisMastersHealthyTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.healthy))
		}
		if res.name == "redis-slave" {
			redisSlavesTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.total))
			redisSlavesHealthyTotal.WithLabelValues(hostname, redisDB.appVersion).Set(float64(res.healthy))
		}

		for k, v := range res.results {
			result[k] = v
		}
	}

	return result
}

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

func checkConnection(connection string, password string) string {
	if _, err := createRedisClient(connection, password).Ping().Result(); err != nil {
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
