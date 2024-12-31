package controllers

import (
	"backend/config"
	"backend/constants"
	"backend/internal/models"
	"backend/pkg/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTodos(c *fiber.Ctx) error {
	collection := config.GetCollection(constants.TODO_COLLECTION)
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Unable to fetch all todos")
	}
	var todos []models.Todo
	if err = cursor.All(context.Background(), &todos); err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Unable to parse all todos")
	}
	return utils.JSONResponse(c, 200, todos, "All todos are passes")
}

func AddTodo(c *fiber.Ctx) error {
	collection := config.GetCollection(constants.TODO_COLLECTION)
	var todo models.Todo
	if err := c.BodyParser(&todo); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Invalid Input")
	}
	todo.ID = primitive.NewObjectID()
	todo.CreatedAt = time.Now().Unix()
	_, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusInternalServerError, nil, "Could not create todo")
	}
	return utils.JSONResponse(c, fiber.StatusCreated, todo, "Todo created successfully")
}

func UpdateTodo(c *fiber.Ctx) error {
	//get collection
	collection := config.GetCollection(constants.TODO_COLLECTION)
	// get id from params
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Invalid ID")
	}
	//todo model
	var todo models.Todo
	// parser title
	if err := c.BodyParser(&todo); err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Invalid Input")
	}
	// create updated obj
	update := bson.M{
		"$set": bson.M{
			"title":    todo.Title,
			"complete": todo.Completed,
		},
	}
	// update title in db
	result, err := collection.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	if err != nil || result.MatchedCount == 0 {
		return utils.JSONResponse(c, fiber.StatusNotFound, nil, "Todo not found")
	}
	// return
	return utils.JSONResponse(c, fiber.StatusOK, update, "Todo updated Successfully")
}
func DeleteTodo(c *fiber.Ctx) error {
	// get collection
	collection := config.GetCollection(constants.TODO_COLLECTION)
	//get id from params
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return utils.JSONResponse(c, fiber.StatusBadRequest, nil, "Invalid ID")
	}
	// delete todo with that id
	result, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil || result.DeletedCount == 0 {
		return utils.JSONResponse(c, fiber.StatusNotFound, nil, "Todo not found")
	}
	//return
	return utils.JSONResponse(c, fiber.StatusOK, nil, "Todo deleted successfully")
}
