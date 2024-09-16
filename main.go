package main

import (
	"html/template"
	"io"
	"log"
	"os"

	"github.com/catalog/virtbrowser/internal/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	tmpl   *template.Template
	logger *log.Logger
)

func main() {

	// # TODO: update to config
	// Set up logging to file and console with rotation
	logFile := &lumberjack.Logger{
		Filename:   "virtbrowser.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     2,    // days
		Compress:   true, // disabled by default
	}

	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	logger = log.New(gin.DefaultWriter, "", log.LstdFlags)
	handlers.SetLogger(logger)

	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	var err error
	tmpl, err = template.ParseGlob("web/templates/*.html")
	if err != nil {
		logger.Fatalf("Error parsing templates: %v", err)
	}

	router.SetHTMLTemplate(tmpl)
	router.Static("/static", "./web/static")

	// router.GET("/", handlers.Root)
	router.GET("/login", handlers.ShowLoginPage)
	router.GET("/logout", handlers.PerformLogout)
	router.POST("/login", handlers.PerformLogin)
	router.GET("/dashboard", handlers.ShowDashboard)

	router.GET("/terminal", handlers.ShowTerminal)             // Serve the HTML page
	router.GET("/ws/terminal", handlers.ShowTerminalWebSocket) // Handle WebSocket connections
	router.GET("/logs", handlers.TailLogFile)                  // Serve the HTML page
	router.GET("/ws/logs", handlers.TailLogWebSocket)          // Handle WebSocket connections

	// group routes for vms
	routeGrpVms := router.Group("/vms")
	routeGrpVms.GET("", handlers.ShowVMList)
	// routeGrpVms.POST("", CreateVM)
	// routeGrpVms.PUT("/:id", UpdateVM)
	// routeGrpVms.DELETE("/:id", DeleteVM)

	router.Use(handlers.RedirectToDashboardOrLogin)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		logger.Fatalf("Failed to run server: %v", err)
	}
}
