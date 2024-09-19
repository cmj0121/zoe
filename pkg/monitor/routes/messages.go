package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cmj0121/zoe/pkg/types"
)

// Get the daily popular messages based on the passed-in field.
func APIMessagePopular(ctx *gin.Context) {
	field := ctx.Param("field")
	switch field {
	case "client_ip":
	case "username":
	case "password":
	default:
		// show the default 404 page
		ctx.String(http.StatusNotFound, "404 page not found")
	}

	report := types.DailyPopularMessages(ctx, field, 10)
	ctx.JSON(http.StatusOK, report)
}
