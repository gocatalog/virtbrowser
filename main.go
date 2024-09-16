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
)

var tmpl *template.Template

func main() {

	// Set up logging to file and console
	// # TODO: update to config
	logFile, err := os.OpenFile("virtbrowser.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(gin.DefaultWriter)

	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

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

	router.Use(handlers.RedirectToDashboardOrLogin)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
