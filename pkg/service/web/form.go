package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/cmj0121/zoe/pkg/service/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var (
	FORM_SVC_NAME = "form"
)

// The honeypot service for Web Login Form
type Form struct {
	Bind string // the address to listen
}

// Run the Form honeypot service
func (f *Form) Run(ch chan<- *types.Message) error {
	log.Info().Msg("running the Form honeypot service ...")

	// set the gin to release mode
	gin.SetMode(gin.ReleaseMode)
	routes := gin.New()
	routes.Use(gin.Recovery())
	routes.Use(f.redirectToIndex())

	// setup HTML templates
	switch tmpl, err := template.ParseFS(templates, "web/template/*.htm"); err {
	case nil:
		routes.SetHTMLTemplate(tmpl)
	default:
		log.Warn().Err(err).Msg("failed to parse the HTML templates")
		return err
	}

	routes.GET("/", f.formPage)
	routes.GET("/static/:filepath", f.static)
	routes.POST("/v/login", f.formPage)
	routes.Run(f.Bind)

	return nil
}

func (f *Form) redirectToIndex() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() == http.StatusNotFound {
			c.Redirect(http.StatusFound, "/")
		}
	}
}

// Show the honeypot service for web login form
func (f *Form) formPage(c *gin.Context) {
	var errorMsg string

	if c.Request.Method == http.MethodPost {
		errorMsg = "Invalid username or password"

		username := c.PostForm("username")
		password := c.PostForm("password")

		log.Info().Str("username", username).Str("password", password).Msg("login failed")
	}

	c.HTML(http.StatusOK, "index.htm", gin.H{
		"year":       time.Now().Year(),
		"traceToken": uuid.New().String(),
		"error":      errorMsg,
	})
}

// Get the static file from the embed.FS
func (f *Form) static(c *gin.Context) {
	filepath := fmt.Sprintf("web/static/%v", c.Param("filepath"))
	c.FileFromFS(filepath, http.FS(static))
}
