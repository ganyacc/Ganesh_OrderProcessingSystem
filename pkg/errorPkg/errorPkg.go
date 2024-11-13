package errorPkg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/logger"
	"gorm.io/gorm"
)

type CustomError struct {
	ErrorMsg   string
	StatusCode int
}

func HandleError(tx *gorm.DB, err error) *CustomError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Error("Error checking last order status: ", err)
		tx.Rollback()
		serverErr := &CustomError{
			StatusCode: http.StatusBadRequest,
			ErrorMsg:   fmt.Sprintf("Error: %v", err),
		}
		return serverErr
	}

	serverErr := &CustomError{
		StatusCode: http.StatusInternalServerError,
		ErrorMsg:   fmt.Sprintf("Error: %v", err),
	}
	return serverErr
}

func (e *CustomError) Error() string {
	return e.ErrorMsg
}

func (e *CustomError) HttpStatusCode() int {
	return e.StatusCode
}

func CustomErrorHandle(statusCode int, message string) *CustomError {
	serverErr := &CustomError{
		StatusCode: statusCode,
		ErrorMsg:   fmt.Sprintf("Error: %v", message),
	}
	return serverErr
}
