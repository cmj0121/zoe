package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIIndex(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{})
}
