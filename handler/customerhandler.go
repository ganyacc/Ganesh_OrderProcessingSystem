package handler

import (
	"net/http"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/entities"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/logger"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type customersHandler struct {
	CustomerRepo repository.CustomerHandler
}

// NewCustomerHandler returns the new instace of type customersHandler
func NewCustomerHandler(customerRepository repository.CustomerHandler) CustomerHandler {
	return &customersHandler{
		CustomerRepo: customerRepository,
	}
}

// GetAllCustomers returns the all available customers
func (cm customersHandler) GetAllCustomers(c echo.Context) error {

	logger.Log.Info("GET /api/customers - Retrieving all customers")

	customers, err := cm.CustomerRepo.GetAllCustomers()
	if err != nil {
		logger.Log.Error("Error retrieving customers: ", err)
		return c.JSON(err.HttpStatusCode(), err.Error())
	}

	logger.Log.Infof("Successfully retrieved %d customers", len(customers))
	return c.JSON(http.StatusOK, customers)

}

// GetCustomerByID returns the all details of customer by customerId
func (cm customersHandler) GetCustomerByID(c echo.Context) error {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id.")
	}

	customer, errs := cm.CustomerRepo.GetCustomerByID(id)
	if errs != nil {
		return c.JSON(errs.HttpStatusCode(), errs.Error())
	}

	return c.JSON(http.StatusOK, customer)

}

// CreateOrder handler
func (cm customersHandler) CreateOrder(c echo.Context) error {

	logger.Log.Info("POST /api/orders - Creating a new order")

	var orderRequest entities.OrderRequest

	//parse request body
	err := c.Bind(&orderRequest)
	if err != nil {
		logger.Log.Warn("Invalid request payload for creating order")
		return c.JSON(http.StatusInternalServerError, err)
	}

	//validate customerId
	_, err = uuid.Parse(orderRequest.CustomerID)
	if err != nil {
		logger.Log.Warn("Invalid CustomerId: ", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	//validate productIds
	if len(orderRequest.ProductIDs) > 0 {
		for _, id := range orderRequest.ProductIDs {
			_, err = uuid.Parse(id)
			if err != nil {
				logger.Log.Warn("Invalid productId: ", err)
				return c.JSON(http.StatusInternalServerError, err)
			}
		}
	}

	logger.Log.Infof("Processing order for customer_id: %v", orderRequest.CustomerID)

	order, errs := cm.CustomerRepo.CreateOrder(orderRequest.CustomerID, orderRequest.ProductIDs)
	if errs != nil {
		logger.Log.Warn("Error creating an order: ", errs.Error())
		return c.JSON(errs.HttpStatusCode(), errs.Error())
	}

	logger.Log.Infof("Order created successfully with ID: %d", order.ID)
	return c.JSON(http.StatusCreated, order)
}

// GetOrderByID handler retrieves order by id
func (cm customersHandler) GetOrderByID(c echo.Context) error {

	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id.")
	}

	order, errs := cm.CustomerRepo.GetOrderByID(id)
	if err != nil {
		return c.JSON(errs.HttpStatusCode(), errs.Error())
	}

	return c.JSON(http.StatusOK, order)

}
