package handlers

import (
	"bufio"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// TailLogFile serves the log viewer HTML page
func TailLogFile(c *gin.Context) {
	filePath := c.Query("file")
	if filePath == "" {
		c.String(http.StatusBadRequest, "File path is required")
		return
	}

	// c.HTML(http.StatusOK, "log_viewer.html", gin.H{
	// 	"file": filePath,
	// })

	renderPartial(c, "logger.html", "Log")
}

// TailLogWebSocket handles the WebSocket connection for tailing the log file
func TailLogWebSocket(c *gin.Context) {
	filePath := c.Query("file")
	if filePath == "" {
		c.String(http.StatusBadRequest, "File path is required")
		return
	}

	wsConn, err := upgradeWebSocket(c)
	if err != nil {
		logger.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	defer wsConn.Close()

	file, err := os.Open(filePath)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Failed to open log file: "+err.Error()))
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		err = wsConn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			logger.Println("Failed to write message to WebSocket:", err)
			break
		}
	}
}
