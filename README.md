**CRM System - README**
1. Introduction
This CRM (Customer Relationship Management) system is designed to manage users, customers, and their interactions. It allows users to create, update, delete, and retrieve data related to customers and interactions. The system is built using Go, MongoDB, and the Gin framework, ensuring scalability and efficiency.

2. Technologies Used
Go (Golang): The primary programming language used to build the backend API.
Gin: A web framework used for building the API.
MongoDB: A NoSQL database used for storing user, customer, and interaction data.
JWT (JSON Web Tokens): Used for securing API endpoints.
Postman: For testing API endpoints.

3. Database Schema Design
Database Structure
The system uses MongoDB as the database. The primary collections in the database are:
Users: Stores data about the system users (e.g., name, email, password, etc.).
Customers: Stores information about the customers.
Interactions: Tracks interactions between users and customers (e.g., meetings, tickets).

Schema Design

Users Collection:

_id: ObjectID (Primary Key)
name: String
email: String (unique)
password: String (hashed)
company: String
status: String (active/inactive)
Customers Collection:

_id: ObjectID (Primary Key)
name: String
email: String (unique)
company: String
status: String (active/inactive)
Interactions Collection:

_id: ObjectID (Primary Key)
user_id: ObjectID (references Users)
customer_id: ObjectID (references Customers)
type: String (e.g., "Meeting", "Ticket")
description: String
status: String (e.g., "Open", "Resolved")
scheduled_at: DateTime

Database Schema Diagram
+----------------+       +-----------------+       +-----------------+
|    Users       |       |   Customers     |       |  Interactions   |
+----------------+       +-----------------+       +-----------------+
| _id: ObjectID  |       | _id: ObjectID   |       | _id: ObjectID   |
| name: String   |       | name: String    |       | user_id: ObjectID|
| email: String  |       | email: String   |       | customer_id: ObjectID|
| password: String|      | company: String |       | type: String     |
| company: String|       | status: String  |       | description: String|
| status: String |       +-----------------+       | status: String   |
+----------------+                                | scheduled_at: DateTime|
                                                  +-----------------+
4. System Architecture
Overview
The system is built as a RESTful API with three main entities: Users, Customers, and Interactions. Each entity has its own set of CRUD operations, and the system is secured using JWT for authentication.System Architecture Diagram
+-----------------+       +-------------------+       +-------------------+
|   Client        |       |     API Server    |       |   MongoDB         |
|  (Postman/Browser)|----->|   (Go + Gin)     |<----->|    Database       |
+-----------------+       +-------------------+       +-------------------+
                                |  |  |
                                |  |  +--------------------+
                                |  +--> JWT Authentication |
                                +-------------------------+
Key Components:
Client: The front-end client or Postman used to make HTTP requests to the API.
API Server: The backend server running the Go application using the Gin framework. It handles the business logic and communicates with the database.
MongoDB: The database that stores the data for users, customers, and interactions.

5. API Documentation
For detailed information on each API endpoint, refer to the API Documentation.

6. Setup and Installation
###Prerequisites:
-Go installed on your machine.
-MongoDB instance running locally or on MongoDB Atlas.
###Installation Steps:
-git clone https://github.com/Nani2139/UdayKumar_Assignment.git
-cd UdayKumar_Assignment
-Install dependencies
-go run main.go
-Access the API at http://localhost:8080.

8. Usage
Example API Requests:
Create a User:
POST http://localhost:8080/users
Login:
POST http://localhost:8080/users/login
Create a Customer:
POST http://localhost:8080/customers
