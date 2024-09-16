package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// TailLogFile serves the log file content over HTTP
func TailLogFile(c *gin.Context) {
	logFile, err := os.Open("virtbrowser.log")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to open log file: %v", err)
		return
	}
	defer logFile.Close()

	stat, err := logFile.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get log file info: %v", err)
		return
	}

	http.ServeContent(c.Writer, c.Request, "virtbrowser.log", stat.ModTime(), logFile)
}
