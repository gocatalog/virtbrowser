package main

import (
	"html/template"
	"log"

	"github.com/catalog/virtbrowser/internal/handlers"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var tmpl *template.Template

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	// Parse templates
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

	router.Use(handlers.RedirectToDashboardOrLogin)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
