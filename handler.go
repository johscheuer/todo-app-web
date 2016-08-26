package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func readTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := createRedisClient(slaveConnection, slavePassword).LRange(mux.Vars(req)["key"], -100, 100)
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}

	fmt.Println(cmd.Val())
	generateJSONResponse(rw, cmd.Val())
}

func insertTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := createRedisClient(masterConnection, masterPassword).RPush(mux.Vars(req)["key"], mux.Vars(req)["value"])
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}
	readTodoHandler(rw, req)
}

func deleteTodoHandler(rw http.ResponseWriter, req *http.Request) {
	cmd := createRedisClient(masterConnection, masterPassword).LRem(mux.Vars(req)["key"], 1, mux.Vars(req)["value"])
	if cmd.Err() != nil {
		fmt.Println(cmd.Err())
		http.Error(rw, cmd.Err().Error(), 500)
	}
	readTodoHandler(rw, req)
}

func healthCheckHandler(rw http.ResponseWriter, req *http.Request) {
	generateJSONResponse(rw, getHealthStatus())
}

func whoAmIHandler(rw http.ResponseWriter, r *http.Request) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}

	addresses, err := getAllAddresses(ifaces)
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}

	generateJSONResponse(rw, addresses)
}
