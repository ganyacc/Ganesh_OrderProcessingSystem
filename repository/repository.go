package repository

import (
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/entities"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/pkg/errorPkg"
)

type CustomerHandler interface {
	GetAllCustomers() ([]entities.Customer, errorPkg.CustomErrors)
	GetCustomerByID(id string) (*entities.Customer, errorPkg.CustomErrors)
	CreateOrder(customerID string, productIds []string) (*entities.Order, errorPkg.CustomErrors)
	GetOrderByID(orderId string) (*entities.Order, errorPkg.CustomErrors)
}
