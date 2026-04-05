package main

import "github.com/gin-gonic/gin"

func setupRoutes(app *gin.Engine) {
	app.GET("/api/todos", getTodos)
	app.POST("/api/todos", createTodos)
	app.PATCH("/api/todos/:id", updateTodos)
	app.PATCH("/api/todos/:id/complete", completeTodo)
	app.DELETE("/api/todos/:id", deleteTodos)
}
