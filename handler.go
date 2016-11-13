package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gorilla/mux"
)

func readTodoHandler(rw http.ResponseWriter, req *http.Request) {
	todos, err := database.GetAllTodos()
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}
	fmt.Println(todos)
	generateJSONResponse(rw, todos)
}

func insertTodoHandler(rw http.ResponseWriter, req *http.Request) {
	if err := database.SaveTodo(mux.Vars(req)["value"]); err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}

	readTodoHandler(rw, req)
}

func deleteTodoHandler(rw http.ResponseWriter, req *http.Request) {
	if err := database.DeleteTodo(mux.Vars(req)["value"]); err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), 500)
	}

	readTodoHandler(rw, req)
}

func healthCheckHandler(rw http.ResponseWriter, req *http.Request) {
	generateJSONResponse(rw, database.GetHealthStatus())
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

func versionHandler(rw http.ResponseWriter, req *http.Request) {
	generateJSONResponse(rw, map[string]string{"version": appVersion})
}
