package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"gopkg.in/redis.v4"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func createRedisClient(addr string) *(redis.Client) {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func readTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := createRedisClient("redis-slave:6379").LRange(mux.Vars(req)["key"], -100, 100)
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}

	fmt.Println(cmd.Val())

	membersJSON, err := json.MarshalIndent(cmd.Val(), "", "  ")
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}
	rw.Write(membersJSON)
}

func insertTodoHandler(rw http.ResponseWriter, req *http.Request) {
	//Insert todo test[negroni]
	fmt.Printf("Insert %s %s\n", mux.Vars(req)["key"], mux.Vars(req)["value"])
	cmd := createRedisClient("redis-master:6379").RPush(mux.Vars(req)["key"], mux.Vars(req)["value"])
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}
	readTodoHandler(rw, req)
}

func deleteTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := createRedisClient("redis-master:6379").LRem(mux.Vars(req)["key"], 1, mux.Vars(req)["value"])
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}
	readTodoHandler(rw, req)
}

func healthCheck(rw http.ResponseWriter, req *http.Request) {
	okString := "ok"
	result := map[string]string{"self": okString}

	if _, err := createRedisClient("redis-master:6379").Ping().Result(); err != nil {
		result["redis-master"] = err.Error()
	} else {
		result["redis-master"] = okString
	}

	if _, err := createRedisClient("redis-slave:6379").Ping().Result(); err != nil {
		result["redis-slave"] = err.Error()
	} else {
		result["redis-slave"] = okString
	}

	aliveJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}
	rw.Write(aliveJSON)
}

func responseWithIPs(rw http.ResponseWriter, r *http.Request) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}

	var addresses []string
	for _, i := range ifaces {
		addrs, erro := i.Addrs()
		if erro != nil {
			fmt.Println(err)
			http.Error(rw, err.Error(), 500)
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
	r.Path("/read/{key}").Methods("GET").HandlerFunc(readTodoHandler)
	r.Path("/insert/{key}/{value}").Methods("GET").HandlerFunc(insertTodoHandler)
	r.Path("/delete/{key}/{value}").Methods("GET").HandlerFunc(deleteTodoHandler)
	r.Path("/health").Methods("GET").HandlerFunc(healthCheck)
	r.Path("/whoami").Methods("GET").HandlerFunc(responseWithIPs)

	n := negroni.Classic()
	n.UseHandler(r)
	http.ListenAndServe(":3000", n)
}
