package handlers

import (
	"net/http"

	"github.com/catalog/virtbrowser/internal/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RedirectToDashboardOrLogin(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.Redirect(http.StatusFound, "/dashboard")
	}
}

func ShowLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Title": "Login",
		"Error": c.Query("error"),
	})
}

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
	c.Redirect(http.StatusFound, "/login")
}
