package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB ObjectID
	Title     string             `bson:"title"`         // TODO item title
	Completed bool               `bson:"completed"`     // Status of TODO item
	CreatedAt int64              `bson:"created_at"`    // Timestamp
}

// UserID    primitive.ObjectID `bson:"user_id"`       // Reference to User
