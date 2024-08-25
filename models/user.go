package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
    ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name           string             `json:"name" bson:"name"`
    Email          string             `json:"email" bson:"email"`
    Password       string             `json:"password" bson:"password"`
    ContactInfo    string             `json:"contact_info" bson:"contact_info"`
    Company        string             `json:"company" bson:"company"`
    Status         string             `json:"status" bson:"status"`
    Notes          string             `json:"notes" bson:"notes"`
    Role           string             `json:"role" bson:"role"` // e.g., Admin, User
}

type Customer struct {
    ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Name           string             `json:"name" bson:"name"`
    ContactInfo    string             `json:"contact_info" bson:"contact_info"`
    Company        string             `json:"company" bson:"company"`
    Status         string             `json:"status" bson:"status"`
    Notes          string             `json:"notes" bson:"notes"`
}

