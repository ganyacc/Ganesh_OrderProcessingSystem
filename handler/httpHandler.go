package handler

import "github.com/labstack/echo/v4"

type CustomerHandler interface {
	GetAllCustomers(c echo.Context) error
	GetCustomerByID(c echo.Context) error
	CreateOrder(c echo.Context) error
	GetOrderByID(c echo.Context) error
}
