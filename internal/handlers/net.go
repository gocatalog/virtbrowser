package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NetList ...
func NetList(c *gin.Context) {
	data := []Data{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	c.JSON(http.StatusOK, data)
}
