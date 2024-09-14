package handlers

import (
	"net/http"

	"github.com/catalog/virtbrowser/internal/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RedirectToDashboardOrLogin ...
func RedirectToDashboardOrLogin(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil && c.Request.URL.Path != "/login" {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	} else if user != nil && c.Request.URL.Path != "/dashboard" {
		c.Redirect(http.StatusFound, "/dashboard")
		c.Abort()
	} else {
		c.Next()
	}
}

// ShowLoginPage ...
func ShowLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title": "Login",
		"Error": c.Query("error"),
	})
}

// PerformLogin ...
func PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	err := auth.Authenticate(username, password)
	if err != nil {
		c.Redirect(http.StatusFound, "/login?error=Authentication failed")
		return
	}

	session := sessions.Default(c)
	session.Set("user", username)
	session.Save()

	c.Redirect(http.StatusFound, "/dashboard")
}

func ShowDashboard(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Title": "Dashboard",
		"User":  user,
	})
}

func PerformLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
