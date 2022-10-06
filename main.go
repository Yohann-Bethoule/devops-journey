package main

import (
	"go-rest-api/models"
	"log"
	"net/http"
	"strconv"

	_ "go-rest-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Devops Journey Go API
// @version         1.0
// @description     Petit test pour le développement d'une API en Go
// @termsOfService  http://swagger.io/terms/

// @contact.name   Yohann Bethoule
// @contact.url    http://www.swagger.io/support
// @contact.email  ybethoule@figarocms.fr

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      ybethoule-rest-api.endpoints.project-staging1.cloud.goog
// @BasePath  /
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

// createTodo godoc
// @Summary      Create a todo
// @Description  Create a todo with a label
// @ID create-todo
// @Accept       json
// @Produce      json
// @Param        label   body      string  true  "Label of the task to do"
// @Success      200  {string}  string "Ok"
// @Failure      400  {string}  string "Invalid request payload"
// @Failure      404  {string}  string "Not found"
// @Failure      500  {string}  string "Internal server error"
// @Router       / [post]
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

// fetchAllTodo godoc
// @Summary      Fetch all todos
// @Description  Fetch all todos in the database
// @ID fetch-all-todos
// @Produce      json
// @Success      200  {array}  models.Todo
// @Failure      400  {string}  Invalid request payload
// @Failure      404  {string}  Not found
// @Failure      500  {string}  Internal server error
// @Router       / [get]
func fetchAllTodo(c *gin.Context) {
	todos, err := models.GetTodos()
	checkErr(err)
	if todos == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": todos})
	}
}

// fetchSingleTodo godoc
// @Summary      Fetch a todo
// @Description  Fetch a todo from its id
// @ID fetch-single-todo
// @Param        id   path      integer  true  "todo ID"
// @Produce      json
// @Success      200  {object}  models.Todo
// @Failure      400  {string}  Invalid request payload
// @Failure      404  {string}  Not found
// @Failure      500  {string}  Internal server error
// @Router       /{id} [get]
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

// updateTodo godoc
// @Summary      Update a todo
// @Description  Update label and done status of a todo
// @ID update-todo
// @Param        id   	path	integer  true   "Todo ID"
// @Param        label  body    string   false  "Todo label"
// @Param        isDone body    boolean  false  "True if the todo is done"
// @Accept       json
// @Produce      json
// @Success      200  {string}  Todo Updated Successfully
// @Failure      400  {string}  Invalid request payload
// @Failure      404  {string}  Not found
// @Failure      500  {string}  Internal server error
// @Router       /{id} [put]
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

// deleteTodo godoc
// @Summary      Delete a todo
// @Description  Delete a todo from its id
// @ID delete-todo
// @Param        id   	path	integer  true   "Todo ID"
// @Produce      json
// @Success      200  {string}  Todo Deleted Successfully
// @Failure      400  {string}  Invalid request payload
// @Failure      404  {string}  Not found
// @Failure      500  {string}  Internal server error
// @Router       /{id} [delete]
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
