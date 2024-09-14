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

func ShowLoginPage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	if user != nil {
		c.Redirect(http.StatusFound, "/dashboard")
		return
	}

	if c.GetHeader("HX-Request") != "" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Title": "Login",
			"Error": c.Query("error"),
		})
	} else {
		c.HTML(http.StatusOK, "base.html", gin.H{
			"Title": "Login",
			"Error": c.Query("error"),
		})
	}
}

func PerformLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	err := auth.Authenticate(username, password)
	if err != nil {
		if c.GetHeader("HX-Request") != "" {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"Title": "Login",
				"Error": "Authentication failed",
			})
		} else {
			c.Redirect(http.StatusFound, "/login?error=Authentication failed")
		}
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

	if c.GetHeader("HX-Request") != "" {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"Title": "Dashboard",
			"User":  user,
		})
	} else {
		c.HTML(http.StatusOK, "base.html", gin.H{
			"Title": "Dashboard",
			"User":  user,
		})
	}
}

func PerformLogout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusFound, "/")
}
