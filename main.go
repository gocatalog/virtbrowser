package main

import (
	"log"

	"github.com/catalog/virtbrowser/internal/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("web/templates/*")

	router.GET("/", handlers.RedirectToDashboardOrLogin)
	router.GET("/login", handlers.ShowLoginPage)
	router.GET("/logout", handlers.PerformLogout)
	router.POST("/login", handlers.PerformLogin)
	router.GET("/dashboard", handlers.ShowDashboard)

	// Start the server
	// TODO: make port dynamic
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
