package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ShowVMList ...
func ShowVMList(c *gin.Context) {
	renderPartial(c, "vm.html", "Virtual Machine")
}

// ShowTerminal ...
func ShowTerminal(c *gin.Context) {
	renderPartial(c, "term.html", "Terminal")
	// c.HTML(http.StatusOK, "term_.html", nil)
}

// VMList ...
func VMList(c *gin.Context) {
	data := []Data{
		{ID: 1, Name: "John Doe", Email: "john@example.com"},
		{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
	}
	c.JSON(http.StatusOK, data)
}
