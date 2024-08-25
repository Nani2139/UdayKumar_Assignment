package routes

import (
    "projectk/controllers"
    "github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
    userRoutes := router.Group("/users")
    {
        userRoutes.POST("/", controllers.CreateUser)
        userRoutes.POST("/login", controllers.LoginUser)

        // Protected routes
        userRoutes.Use(controllers.AuthMiddleware())
        userRoutes.GET("/", controllers.GetUsers)
        userRoutes.GET("/:id", controllers.GetUserByID)
        userRoutes.PUT("/:id", controllers.UpdateUser)
        userRoutes.DELETE("/:id", controllers.DeleteUser)
    }
}
func SetupCustomerRoutes(router *gin.Engine) {
    customerRoutes := router.Group("/customers")
    {
        customerRoutes.POST("/", controllers.CreateCustomer)
        customerRoutes.GET("/", controllers.GetCustomers)
        customerRoutes.GET("/:id", controllers.GetCustomerByID)
        customerRoutes.PUT("/:id", controllers.UpdateCustomer)
        customerRoutes.DELETE("/:id", controllers.DeleteCustomer)
    }
}