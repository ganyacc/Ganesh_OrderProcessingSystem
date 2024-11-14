package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/config"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/database"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/handler"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/logger"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type EchoServer struct {
	app  *echo.Echo
	db   database.Database
	conf *config.Config
}

func NewEchoServer(conf *config.Config, db database.Database) Server {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG)

	return &EchoServer{
		app:  echoApp,
		db:   db,
		conf: conf,
	}
}

func (s *EchoServer) Start() error {
	s.app.Use(middleware.Recover())
	s.app.Use(middleware.Logger())
	s.app.Use(handler.LatencyLogger)

	//initialize logger
	logger.Init()

	// Health check adding
	s.app.GET("/v1/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Ok")
	})

	//initialize routes
	s.Routes()

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	return s.app.Start(serverUrl)
}

// Shutdown gracefully stops the server with a given context
func (s *EchoServer) Shutdown(ctx context.Context) error {
	logrus.Println("Attempting to gracefully shutdown the server...")
	return s.app.Shutdown(ctx)
}

// Routes define the new routes
func (s *EchoServer) Routes() {
	customerRepo := repository.NewCustomerRepository(s.db)
	customerHandler := handler.NewCustomerHandler(customerRepo)

	route := s.app.Group("/api")

	route.GET("/customers", customerHandler.GetAllCustomers)
	route.GET("/customers/:id", customerHandler.GetCustomerByID)
	route.POST("/orders", customerHandler.CreateOrder)
	route.GET("/orders/:id", customerHandler.GetOrderByID)
}
