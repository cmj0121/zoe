package monitor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Show the health check of the monitor service
func (m *Monitor) livez(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

// Check the readiness of the monitor service
func (m *Monitor) readyz(c *gin.Context) {
	switch err := m.DB.Ping(); err {
	case nil:
		c.Writer.WriteHeader(http.StatusOK)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
