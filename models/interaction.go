package models

import (
    "time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Interaction struct {
    ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
    CustomerID  primitive.ObjectID `json:"customer_id" bson:"customer_id"`
    Type        string             `json:"type" bson:"type"` // e.g., Ticket, Meeting
    Description string             `json:"description" bson:"description"`
    Status      string             `json:"status" bson:"status"` // e.g., Open, Resolved
    ScheduledAt time.Time          `json:"scheduled_at" bson:"scheduled_at"`
    CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
