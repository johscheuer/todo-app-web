package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"gopkg.in/redis.v3"
)

func CreateRedisClient(addr string) *(redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func ReadTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := CreateRedisClient("redis-slave:6379").LRange(mux.Vars(req)["key"], -100, 100)
	if cmd.Err() != nil {
		http.Error(rw, cmd.Err().Error(), 500)
		fmt.Println(cmd.Err())
	}
	membersJSON, err := json.MarshalIndent(cmd.Val(), "", "  ")
	if err != nil {
		http.Error(rw, err.Error(), 500)
		fmt.Println(err)
	}
	rw.Write(membersJSON)
}

func InsertTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := CreateRedisClient("redis-master:6379").RPush(mux.Vars(req)["key"], mux.Vars(req)["value"])
	if cmd.Err() != nil {
		http.Error(rw, cmd.Err().Error(), 500)
		fmt.Println(cmd.Err())
	}
	ReadTodoHandler(rw, req)
}

func DeleteTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := CreateRedisClient("redis-master:6379").LRem(mux.Vars(req)["key"], 1, mux.Vars(req)["value"])
	if cmd.Err() != nil {
		http.Error(rw, cmd.Err().Error(), 500)
		fmt.Println(cmd.Err())
	}
	ReadTodoHandler(rw, req)
}

func HealthCheck(rw http.ResponseWriter, req *http.Request) {
	aliveJSON, err := json.MarshalIndent("Alive", "", "  ")
	if err != nil {
		http.Error(rw, err.Error(), 500)
		fmt.Println(err)
	}
	rw.Write(aliveJSON)
}

func ResponseWithIPs(rw http.ResponseWriter, r *http.Request) {
	ifaces, err := net.Interfaces()
	if err != nil {
		http.Error(rw, err.Error(), 500)
		fmt.Println(err)
	}

	var addresses []string
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			http.Error(rw, err.Error(), 500)
			fmt.Println(err)
		}

		for _, addr := range addrs {
			addresses = append(addresses, addr.String())
		}

	}

	addressJSON, err := json.MarshalIndent(addresses, "", "  ")
	rw.Write(addressJSON)
}

func main() {
	r := mux.NewRouter()
	r.Path("/read/{key}").Methods("GET").HandlerFunc(ReadTodoHandler)
	r.Path("/insert/{key}/{value}").Methods("GET").HandlerFunc(InsertTodoHandler)
	r.Path("/delete/{key}/{value}").Methods("GET").HandlerFunc(DeleteTodoHandler)
	r.Path("/health").Methods("GET").HandlerFunc(HealthCheck)
	r.Path("/whoami").Methods("GET").HandlerFunc(ResponseWithIPs)

	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":3000")
}
