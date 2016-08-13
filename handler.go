package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func readTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := createRedisClient(slaveConnection).LRange(mux.Vars(req)["key"], -100, 100)
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
	cmd := createRedisClient(masterConnection).RPush(mux.Vars(req)["key"], mux.Vars(req)["value"])
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}
	readTodoHandler(rw, req)
}

func deleteTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := createRedisClient(masterConnection).LRem(mux.Vars(req)["key"], 1, mux.Vars(req)["value"])
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}
	readTodoHandler(rw, req)
}

func healthCheckHandler(rw http.ResponseWriter, req *http.Request) {
	okString := "ok"
	result := map[string]string{"self": okString}

	checkConnection(result, "redis-master", masterConnection, okString)
	checkConnection(result, "redis-slave", slaveConnection, okString)

	aliveJSON, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}
	rw.Write(aliveJSON)
}

func whoAmIHandler(rw http.ResponseWriter, r *http.Request) {
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
