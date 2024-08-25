package main

import (
    "projectk/config"
    "projectk/routes"
    "github.com/gin-gonic/gin"
)

func main() {
    // Connect to MongoDB
    config.ConnectDB()

    r := gin.Default()

    // Set up user routes
    routes.SetupUserRoutes(r)

    // Set up customer routes
    routes.SetupCustomerRoutes(r)


    // Set up interaction routes
    routes.SetupInteractionRoutes(r)

    // Start the server
    r.Run(":8080")
}
