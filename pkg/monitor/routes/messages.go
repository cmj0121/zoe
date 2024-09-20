package routes

import (
	"embed"
	"net/http"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zoe/pkg/types"
)

//go:embed *.md
var index embed.FS

// Get the daily popular messages, return the markdown page.
func MessagePopular(ctx *gin.Context) {
	// the escape function for the cell in table
	fn := template.FuncMap{
		"escapeTable": func(value string) string {
			return strings.ReplaceAll(value, "|", "\\|")
		},
	}

	// get the template from embeded file
	tmpl, err := template.New("index.md").Funcs(fn).ParseFS(index, "*.md")
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse the template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var report = struct {
		ClientIP []*types.Report  `json:"client_ip"`
		Username []*types.Report  `json:"username"`
		Password []*types.Report  `json:"password"`
		Command  []*types.Message `json:"command"`
	}{}

	report.ClientIP = types.DailyPopularMessages(ctx, "client_ip", 10)
	report.Username = types.DailyPopularMessages(ctx, "username", 10)
	report.Password = types.DailyPopularMessages(ctx, "password", 10)
	report.Command = types.DailyMessage(ctx, "command")

	// render the template
	if err := tmpl.Execute(ctx.Writer, report); err != nil {
		log.Warn().Err(err).Msg("failed to render the template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// make the response as markdown
	ctx.Header("Content-Type", "text/markdown")
	ctx.Status(http.StatusOK)
}

// Get the daily popular messages based on the passed-in field.
func APIMessagePopular(ctx *gin.Context) {
	field := ctx.Param("field")
	switch field {
	case "client_ip":
	case "username":
	case "password":
	case "command":
	default:
		// show the default 404 page
		ctx.String(http.StatusNotFound, "404 page not found")
	}

	report := types.DailyPopularMessages(ctx, field, 10)
	ctx.JSON(http.StatusOK, report)
}
