# Order Processing System

This is a Golang-based Order Processing System that allows users to manage customers and orders through RESTful API endpoints. This project uses the Echo web framework, Gorm ORM for database operations, and includes unit tests for endpoint testing.

## Features

- Retrieve customers (you need to add customers manually in database.)
- Create and retrieve orders
- Integrated with Docker for containerized deployment
- Includes unit tests to validate functionality

## Prerequisites

- [Go 1.18+](https://golang.org/dl/)
- [Docker](https://www.docker.com/get-started)
- [PostgreSQL](https://www.postgresql.org/download/) (optional if using Docker for database)

## Project Structure

```
.
├── config      # Database Configuration
├── database    # Connection to database and schema migration
├── entities    # Entity definitions (Customer, Order, Product)
├── handler     # HTTP handlers for Customer, Order and Product
├── logger      # log initializer
├── pkg         # public package (contains only error pkg now)
├── repository  # Repository layer for DB interactions
├── server      # echo server to run applicatiom
├── testCases   # Unit tests for endpoints
├── config.yaml # configuration file for database connection
├── Dockerfile  # Multistage dockerfile
├── main.go     # Main application file
└── README.md   # Project documentation
```

Getting Started

1. Clone the Repository

```
   git clone https://github.com/yourusername/Ganesh_OrderProcessingSystem.git
   cd Ganesh_OrderProcessingSystem
```

2. Run with Docker
   docker build -t your-img-name .
   docker your-img-name .

3. Run Locally (Without Docker)
   go mod download
   go run main.go

Endpoints

```
Customers
GET /customers - Retrieve all customers
GET /customers/:id - Retrieve a specific customer by ID
```

POST /api/orders -
Create a new order

```
Request Body
{
"customer_id": "your-customer-id",
"product_ids": ["product-id-1", "product-id-2"]
}
```

Get Order By Id

```
GET /orders/:id - Retrieve a specific order by ID
```

Testing
Run Unit Tests
The application includes unit tests for each endpoint. You can run them with:
go test ./testCases

MIT License
This README provides instructions on:

1. Cloning and configuring the project.
2. Running the app locally or in Docker.
3. Using endpoints.
4. Running tests with Go and Docker Compose.
