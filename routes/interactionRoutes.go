package routes

import (
    "projectk/controllers"
    "github.com/gin-gonic/gin"
)

func SetupInteractionRoutes(router *gin.Engine) {
    interactionRoutes := router.Group("/interactions")
    {
        interactionRoutes.POST("/", controllers.CreateInteraction)
        interactionRoutes.GET("/:id", controllers.GetInteractionByID)
		interactionRoutes.GET("/", controllers.GetAllInteractions) 
        interactionRoutes.GET("/customer/:customer_id", controllers.GetInteractionsByCustomerID)
        interactionRoutes.PUT("/:id", controllers.UpdateInteraction)
    }
}
