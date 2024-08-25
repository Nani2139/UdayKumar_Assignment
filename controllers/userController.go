package controllers

import (
	"context"
	"net/http"
	"projectk/config"
	"projectk/models"
	"time"

	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)
var jwtKey = []byte("mySuperSecretKey12345!@#$%67890") // Replace with your secret key

// Claims structure
type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

// CreateUser handles POST requests to add a new user
func CreateUser(c *gin.Context) {
    var newUser models.User

    if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Generate a new ObjectID for the user
    newUser.ID = primitive.NewObjectID()

    // Hash the password before storing
	log.Println(newUser.Password)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
    if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
        return
    }
    newUser.Password = string(hashedPassword)
	log.Println(newUser.Password)

    // Set up the context for the database operation
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Insert the new user into the database
    _, err = config.DB.Collection("users").InsertOne(ctx, newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Respond with success
    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// GetUsers handles GET requests to fetch all users
func GetUsers(c *gin.Context) {
    var users []models.User

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := config.DB.Collection("users").Find(ctx, primitive.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &users); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse users"})
        return
    }

    c.JSON(http.StatusOK, users)
}

// GetUserByID handles GET requests to fetch a user by ID
func GetUserByID(c *gin.Context) {
    id := c.Param("id")
    objID, _ := primitive.ObjectIDFromHex(id)
    var user models.User

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    err := config.DB.Collection("users").FindOne(ctx, primitive.M{"_id": objID}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT requests to update a user by ID
func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    objID, _ := primitive.ObjectIDFromHex(id)
    var updatedUser models.User

    if err := c.ShouldBindJSON(&updatedUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := config.DB.Collection("users").UpdateOne(ctx, primitive.M{"_id": objID}, primitive.M{
        "$set": updatedUser,
    })
    if err != nil || result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
        return
    }

    c.JSON(http.StatusOK, updatedUser)
}

// DeleteUser handles DELETE requests to remove a user by ID
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    objID, _ := primitive.ObjectIDFromHex(id)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := config.DB.Collection("users").DeleteOne(ctx, primitive.M{"_id": objID})
    if err != nil || result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "user not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")

        if tokenString == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Request does not contain an access token"})
            c.Abort()
            return
        }
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
            c.Abort()
            return
        }

        // Store the user's email in the context
        c.Set("email", claims.Email)

        c.Next()
    }
}

// LoginUser handles POST requests to log in a user
func LoginUser(c *gin.Context) {
    var loginData struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Bind the incoming JSON to the loginData struct
    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON provided"})
        return
    }

    var user models.User
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Find the user by email
    err := config.DB.Collection("users").FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Log the retrieved user's details for debugging (remove in production)
    log.Println("User found:", user.Email)
    log.Println("Stored hashed password:", user.Password)
    log.Println("Provided plain text password:", loginData.Password)

    // Compare the stored hashed password with the provided password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password))
    if err != nil {
        log.Println("Password comparison error:", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
        return
    }

    // Log for debugging (remove in production)
    log.Println("Password matched for user:", user.Email)

    // Create a JWT token
    expirationTime := time.Now().Add(24 * time.Hour) // 24-hour expiration
    claims := &Claims{
        Email: loginData.Email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    // Respond with the JWT token
    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// CreateCustomer handles POST requests to create a new customer
func CreateCustomer(c *gin.Context) {
    var newCustomer models.Customer

    if err := c.ShouldBindJSON(&newCustomer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    newCustomer.ID = primitive.NewObjectID()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := config.DB.Collection("customers").InsertOne(ctx, newCustomer)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "Customer created successfully"})
}

// GetCustomers handles GET requests to retrieve all customers
func GetCustomers(c *gin.Context) {
    var customers []models.Customer

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    cursor, err := config.DB.Collection("customers").Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
        return
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &customers); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse customers"})
        return
    }

    c.JSON(http.StatusOK, customers)
}

func GetCustomerByID(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
        return
    }

    var customer models.Customer
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = config.DB.Collection("customers").FindOne(ctx, bson.M{"_id": objID}).Decode(&customer)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
        return
    }

    c.JSON(http.StatusOK, customer)
}


// UpdateCustomer handles PUT requests to update a customer by ID
func UpdateCustomer(c *gin.Context) {
    id := c.Param("id")
    objID, _ := primitive.ObjectIDFromHex(id)
    var updatedCustomer models.Customer

    if err := c.ShouldBindJSON(&updatedCustomer); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := config.DB.Collection("customers").UpdateOne(ctx, bson.M{"_id": objID}, bson.M{
        "$set": updatedCustomer,
    })
    if err != nil || result.MatchedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
        return
    }

    c.JSON(http.StatusOK, updatedCustomer)
}

// DeleteCustomer handles DELETE requests to remove a customer by ID
func DeleteCustomer(c *gin.Context) {
    id := c.Param("id")
    objID, _ := primitive.ObjectIDFromHex(id)

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := config.DB.Collection("customers").DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil || result.DeletedCount == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "Customer not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}
