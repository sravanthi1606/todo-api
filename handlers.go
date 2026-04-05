package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getTodos(c *gin.Context) {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       "success",
		"totalRecords": len(todos),
		"data":         todos,
	})
}

func createTodos(c *gin.Context) {
	var todo Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if todo.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  "Todo body cannot be empty",
		})
		return
	}

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = nil
	todo.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   todo,
	})
}

func updateTodos(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var body struct {
		Body string `json:"body"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if body.Body == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Todo body cannot be empty"})
		return
	}

	filter := bson.M{"_id": objectId}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"body":      body.Body,
			"updatedAt": now,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func completeTodo(c *gin.Context) {
	id := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	filter := bson.M{"_id": objectId}
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"completed": true,
			"updatedAt": now,
		},
	}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func deleteTodos(c *gin.Context) {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	filter := bson.M{"_id": objectID}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"deletedCount": result.DeletedCount,
	})
}
