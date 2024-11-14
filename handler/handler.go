package handler

import (
	"net/http"
	"time"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/entities"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/logger"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomersHandler struct {
	CustomerRepo repository.CustomerHandler
}

// NewCustomerHandler returns the new instace of type customersHandler
func NewCustomerHandler(customerRepository repository.CustomerHandler) CustomerHandler {
	return &CustomersHandler{
		CustomerRepo: customerRepository,
	}
}

// GetAllCustomers returns the all available customers
func (cm CustomersHandler) GetAllCustomers(c echo.Context) error {

	logger.Log.Info("GET /api/customers - Retrieving all customers")

	customers, err := cm.CustomerRepo.GetAllCustomers()
	if err != nil {
		logger.Log.Error("Error retrieving customers: ", err)
		return c.JSON(err.HttpStatusCode(), err.Error())
	}

	return c.JSON(http.StatusOK, customers)

}

// GetCustomerByID returns the all details of customer by customerId
func (cm CustomersHandler) GetCustomerByID(c echo.Context) error {
	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		logger.Log.Warn("Invalid request parameter for fetching customer.")
		return c.JSON(http.StatusBadRequest, "invalid id.")
	}

	customer, errs := cm.CustomerRepo.GetCustomerByID(id)
	if errs != nil {
		logger.Log.Warn("Error fetching customer with id: ", id)
		return c.JSON(errs.HttpStatusCode(), errs.Error())
	}

	return c.JSON(http.StatusOK, customer)

}

// CreateOrder handler
func (cm CustomersHandler) CreateOrder(c echo.Context) error {

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

	return c.JSON(http.StatusCreated, order)
}

// GetOrderByID handler retrieves order by id
func (cm CustomersHandler) GetOrderByID(c echo.Context) error {

	id := c.Param("id")

	_, err := uuid.Parse(id)
	if err != nil {
		logger.Log.Warn("Invalid request parameter to fetch order by id: ", err)
		return c.JSON(http.StatusBadRequest, "invalid id.")
	}

	order, errs := cm.CustomerRepo.GetOrderByID(id)
	if errs != nil {
		logger.Log.Warn("Error fetching order by id: ", errs.Error())
		return c.JSON(errs.HttpStatusCode(), errs.Error())
	}

	return c.JSON(http.StatusOK, order)

}

// Middleware to log API latency
func LatencyLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		err := next(c)
		stop := time.Now()

		latency := stop.Sub(start)
		c.Logger().Infof("Request to %s %s took %v", c.Request().Method, c.Request().URL.Path, latency)
		return err
	}
}
