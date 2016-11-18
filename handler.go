package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
)

func readTodoHandler(c *gin.Context) {
	todos, err := database.GetAllTodos()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}
	fmt.Println(todos)
	c.JSON(http.StatusOK, todos)
}

func insertTodoHandler(c *gin.Context) {
	if err := database.SaveTodo(c.Param("value")); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	readTodoHandler(c)
}

func deleteTodoHandler(c *gin.Context) {
	if err := database.DeleteTodo(c.Param("value")); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	readTodoHandler(c)
}

func healthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, database.GetHealthStatus())
}

func whoAmIHandler(c *gin.Context) {
	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	addresses, err := getAllAddresses(ifaces)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, addresses)
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": appVersion,
	})
}
