package monitor

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	"github.com/cmj0121/zoe/pkg/service/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

//go:embed web/static/*.css
var static embed.FS

// List the messages from the database
func (m *Monitor) messages() ([]*types.Message, error) {
	stmt := `
		SELECT
			message.service,
			message.client_ip,
			message.username,
			message.password,
			message.created_at
		FROM message
		WHERE created_at < ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	now_ns := time.Now().UnixNano()
	rows, err := m.Query(stmt, now_ns, 20)
	if err != nil {
		log.Warn().Err(err).Msg("failed to query the messages")
		return nil, err
	}

	var messages []*types.Message
	for rows.Next() {
		switch message, err := types.MessageFromRows(rows); err {
		case nil:
			messages = append(messages, message)
		default:
			log.Warn().Err(err).Msg("failed to parse the message")
			continue
		}
	}

	return messages, nil
}

// Show the index page of the monitor service
func (m *Monitor) index(c *gin.Context) {
	switch messages, err := m.messages(); err {
	case nil:
		c.HTML(http.StatusOK, "index.htm", gin.H{
			"year":     time.Now().Year(),
			"messages": messages,
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Get the static file from the embed.FS
func (m *Monitor) static(c *gin.Context) {
	filepath := fmt.Sprintf("web/static/%v", c.Param("filepath"))
	c.FileFromFS(filepath, http.FS(static))
}
