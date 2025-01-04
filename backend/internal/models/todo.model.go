package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID        primitive.ObjectID  `bson:"_id,omitempty"`       // MongoDB ObjectID
	Title     string              `bson:"title"`               // TODO item title
	Completed bool                `bson:"completed"`           // Status of TODO item
	CreatedBy *primitive.ObjectID `bson:"createdBy,omitempty"` // Reference to User
	CreatedAt int64               `bson:"created_at"`          // Timestamp
}

//
