package formatError

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RecordNotFound(c *gin.Context, err error, errMessage ...string) {
	errorMessage := "The record not found"
	if len(errMessage) > 0 {
		errorMessage = errMessage[0]
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": errorMessage,
		})
		return
	}

	// Else show internal server error
	InternalServerError(c, err)
}

func InternalServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
	return
}
