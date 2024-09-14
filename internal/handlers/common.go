package handlers

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func renderPartial(c *gin.Context, partial string, title string) {

	session := sessions.Default(c)
	user := session.Get("user")
	if user == nil {
		c.Redirect(http.StatusFound, "/login")
		return
	}

	tmpl, err := template.ParseFiles("web/templates/" + partial)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error rendering partial: %v", err)
		return
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, gin.H{"User": user}); err != nil {
		c.String(http.StatusInternalServerError, "Error executing template: %v", err)
		return
	}
	c.HTML(http.StatusOK, "base.html", gin.H{"Title": title, "User": user, "Content": template.HTML(buf.String())})
}
