package repository

import (
	"fmt"
	"math"
	"net/http"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/database"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/entities"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/logger"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/pkg/errorPkg"
	"gorm.io/gorm"

	"github.com/google/uuid"
)

type customerRepository struct {
	db database.Database
}

func NewCustomerRepository(db database.Database) CustomerHandler {
	return &customerRepository{db: db}
}

// GetAllCustomers returns the list of all customers
func (c customerRepository) GetAllCustomers() ([]entities.Customer, errorPkg.CustomErrors) {
	if c.db == nil {
		logger.Log.Warnf("Database connection not available.")
		return nil, errorPkg.CustomErrorHandle(http.StatusInternalServerError, "Database connection not available")
	}

	db := c.db.GetDb().Begin()

	var customers []entities.Customer
	if err := db.Debug().Model(&entities.Customer{}).Find(&customers).Error; err != nil {
		logger.Log.Error("Error fetching customers: ", err)
		return nil, errorPkg.HandleError(db, err)
	}

	if len(customers) == 0 {
		return nil, errorPkg.CustomErrorHandle(http.StatusNotFound, "no customer available")
	}

	logger.Log.Infof("Customers fetched successfully!. Total number of customers are: %d", len(customers))
	return customers, nil
}

// GetCustomerByID retrives the customer from database by provided Id
func (c customerRepository) GetCustomerByID(id string) (*entities.Customer, errorPkg.CustomErrors) {
	if c.db == nil {
		logger.Log.Warnf("Database connection not available.")
		return nil, errorPkg.CustomErrorHandle(http.StatusInternalServerError, "Database connection not available")
	}

	db := c.db.GetDb().Begin()

	var customer *entities.Customer
	tx := db.Debug().Model(&entities.Customer{}).Find(&customer).Where("id=?", id)
	if tx.Error != nil {
		logger.Log.Error("Error fetching customer: ", tx.Error)
		return nil, errorPkg.HandleError(tx, tx.Error)
	}

	logger.Log.Infof("Customer fetched successfully with ID: %v", id)
	return customer, nil

}

// CreateOrder
func (c customerRepository) CreateOrder(customerID string, productIds []string) (*entities.Order, errorPkg.CustomErrors) {

	if c.db == nil {
		logger.Log.Warnf("Database connection not available.")
		return nil, errorPkg.CustomErrorHandle(http.StatusInternalServerError, "Database connection not available")
	}

	db := c.db.GetDb().Begin()

	var customer entities.Customer
	if err := db.Debug().Where("id = ?", customerID).Find(&customer).Error; err != nil {
		logger.Log.Warnf("Customer with id %v not found.", customerID)
		return nil, errorPkg.HandleError(db, err)
	}

	var lastOrder *entities.Order
	err := db.Debug().Model(&entities.Order{}).Where("customer_id = ?", customerID).Last(&lastOrder).Error
	if err == nil {
		if lastOrder.Status == entities.Unfulfilled {
			logger.Log.Warn("Customer has an unfulfilled order")
			return nil, errorPkg.CustomErrorHandle(http.StatusBadRequest, fmt.Sprintf("Customer with id '%v' has an unfulfilled order", customerID))

		}
	}

	var products []entities.Product
	if err := db.Debug().Model(&entities.Product{}).Find(&products, productIds).Error; err != nil {
		logger.Log.Error("Error retrieving products: ", err)
		return nil, errorPkg.HandleError(db, err)
	}

	totalPrice := 0.0

	for _, product := range products {
		totalPrice += product.Price

	}

	logger.Log.Infof("Calculated total price for order: %.2f", totalPrice)

	totalPrice = math.Round(totalPrice*100) / 100

	custId, err := uuid.Parse(customerID)
	if err != nil {
		logger.Log.Error("Error parsing customerId: ", err)
		return nil, errorPkg.CustomErrorHandle(http.StatusInternalServerError, err.Error())
	}

	order := &entities.Order{
		Customer:   customer,
		CustomerID: custId,
		Products:   products,
		TotalPrice: totalPrice,
		Status:     entities.Unfulfilled,
	}

	if err := db.Debug().Model(&entities.Order{}).Create(&order).Error; err != nil {
		logger.Log.Error("Could not create order: ", err)
		return nil, errorPkg.CustomErrorHandle(http.StatusInternalServerError, "Could not create order.")
	}

	err = db.Commit().Error
	if err != nil {
		logger.Log.Error("Error commiting order transaction: ", err)
		db.Rollback()
		return nil, errorPkg.HandleError(db, err)
	}

	logger.Log.Infof("Order created successfully with ID: %v", order.ID)
	return order, nil
}

// GetOrderByID
func (c customerRepository) GetOrderByID(orderId string) (*entities.Order, errorPkg.CustomErrors) {
	if c.db == nil {
		logger.Log.Warnf("Database connection not available.")
		return nil, errorPkg.CustomErrorHandle(http.StatusInternalServerError, "Database connection not available")
	}

	db := c.db.GetDb().Begin()

	var order *entities.Order
	if err := db.Debug().Preload("Products").Where("id=?", orderId).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Log.Error("Order not found: ", err)
			return nil, errorPkg.CustomErrorHandle(http.StatusNotFound, "Order not found.")
		}

		logger.Log.Error("Error fetching order: ", err)
		return nil, errorPkg.CustomErrorHandle(http.StatusNotFound, err.Error())
	}

	logger.Log.Infof("Order fetched successfully with ID: %v", order.ID)
	return order, nil
}
