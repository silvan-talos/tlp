package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/silvan-talos/tlp/example"
)

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, example.ErrInternal):
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error occurred"})
		return
	case errors.Is(err, example.ErrNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	default:
		// fallback to a 400 Bad Request error response, even tough it's not the case for all errors, I know
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}
