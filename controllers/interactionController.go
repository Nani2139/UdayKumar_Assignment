package controllers

import (
	"context"
	"net/http"
	"projectk/config"
	"projectk/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllInteractions handles GET requests to retrieve all interactions
func GetAllInteractions(c *gin.Context) {
    var interactions []models.Interaction

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := config.DB.Collection("interactions").Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interactions"})
        return
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &interactions); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse interactions"})
        return
    }

    c.JSON(http.StatusOK, interactions)
}

// CreateInteraction handles POST requests to create a new interaction (e.g., scheduling a meeting)
func CreateInteraction(c *gin.Context) {
    var newInteraction models.Interaction

    if err := c.ShouldBindJSON(&newInteraction); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newInteraction.ID = primitive.NewObjectID()
    newInteraction.CreatedAt = time.Now()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := config.DB.Collection("interactions").InsertOne(ctx, newInteraction)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create interaction"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Interaction created successfully"})
}

// GetInteractionByID handles GET requests to retrieve a specific interaction by ID
func GetInteractionByID(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    var interaction models.Interaction
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = config.DB.Collection("interactions").FindOne(ctx, bson.M{"_id": objID}).Decode(&interaction)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Interaction not found"})
        return
    }

    c.JSON(http.StatusOK, interaction)
}

// GetInteractionsByCustomerID handles GET requests to retrieve all interactions for a specific customer
func GetInteractionsByCustomerID(c *gin.Context) {
    customerID := c.Param("customer_id")
    objID, err := primitive.ObjectIDFromHex(customerID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Customer ID format"})
        return
    }

    var interactions []models.Interaction
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := config.DB.Collection("interactions").Find(ctx, bson.M{"customer_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch interactions"})
        return
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &interactions); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse interactions"})
        return
    }

    c.JSON(http.StatusOK, interactions)
}

// UpdateInteraction handles PUT requests to update an interaction (e.g., marking a ticket as resolved)
func UpdateInteraction(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    var updatedInteraction models.Interaction
    if err := c.ShouldBindJSON(&updatedInteraction); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := config.DB.Collection("interactions").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
        "$set": updatedInteraction,
    })
    if err != nil || result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "Interaction not found"})
        return
    }

    c.JSON(http.StatusOK, updatedInteraction)
}
