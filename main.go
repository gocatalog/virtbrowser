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

var tmpl *template.Template

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
	log.SetOutput(gin.DefaultWriter)

	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	var err error
	tmpl, err = template.ParseGlob("web/templates/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	router.SetHTMLTemplate(tmpl)
	router.Static("/static", "./web/static")

	// router.GET("/", handlers.Root)
	router.GET("/login", handlers.ShowLoginPage)
	router.GET("/logout", handlers.PerformLogout)
	router.POST("/login", handlers.PerformLogin)
	router.GET("/dashboard", handlers.ShowDashboard)

	router.GET("/terminal", handlers.ShowTerminal)
	router.GET("/ssh", handlers.HandleTerm)

	// group routes for vms
	routeGrpVms := router.Group("/vms")
	routeGrpVms.GET("", handlers.ShowVMList)
	// routeGrpVms.POST("", CreateVM)
	// routeGrpVms.PUT("/:id", UpdateVM)
	// routeGrpVms.DELETE("/:id", DeleteVM)

	router.GET("/logs", handlers.TailLogFile) // Add this line

	router.Use(handlers.RedirectToDashboardOrLogin)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
