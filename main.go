package main

import (
	"go-rest-api/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1")
	{
		v1.POST("/", createTodo)
		v1.GET("/", fetchAllTodo)
		v1.GET("/:id", fetchSingleTodo)
		v1.PUT("/:id", updateTodo)
		v1.DELETE("/:id", deleteTodo)
	}

	models.ConnectDatabase()

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createTodo(c *gin.Context) {
	var json models.Todo

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error 1": err.Error()})
		return
	}

	json.Done = false

	success, err := models.InsertTodo(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Todo Created Successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error 2": err.Error()})
	}
}

func fetchAllTodo(c *gin.Context) {
	todos, err := models.GetTodos(10)
	checkErr(err)
	if todos == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": todos})
	}
}

func fetchSingleTodo(c *gin.Context) {
	id := c.Param("id")
	todo, err := models.GetTodoById(id)
	checkErr(err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		if todo.Id == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": todo})
	}
}

func updateTodo(c *gin.Context) {
	var json models.Todo

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error 1": err.Error()})
		return
	}

	todoId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	json.Id = todoId

	success, err := models.UpdateTodo(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Todo Updated Successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error 2": err.Error()})
	}
}

func deleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
	}

	success, err := models.DeleteTodo(id)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Todo Deleted Successfully"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error 2": err.Error()})
	}
}
