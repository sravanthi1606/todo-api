package main

import "github.com/gofiber/fiber/v2"

func setupRoutes(app *fiber.App) {
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Patch("/api/todos/:id/complete", completeTodo)
	app.Delete("/api/todos/:id", deleteTodos)
}
