package handlers

import (
	"github.com/gin-gonic/gin"
)

// ShowVMList ...
func ShowVMList(c *gin.Context) {
	renderPartial(c, "vm.html", "Virtual Machine")
}

// ShowTerminal ...
func ShowTerminal(c *gin.Context) {
	renderPartial(c, "term.html", "Terminal")
	// c.HTML(http.StatusOK, "term_.html", nil)
}
