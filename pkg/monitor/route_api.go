package monitor

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cmj0121/zoe/pkg/service/types"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

var (
	// The default page size
	DEFAULT_SIZE = 40
)

// List all the records with pagination
func (m *Monitor) listMessages(c *gin.Context) {
	after := m.PageAfter(c)
	size := m.PageSize(c)

	stmt, err := m.Prepare(`
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
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Info().Str("after", after.String()).Int("size", size).Msg("list the messages")
	rows, err := stmt.Query(after, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var messages []*types.Message
	for rows.Next() {
		switch message, err := types.MessageFromRows(rows); err {
		case nil:
			messages = append(messages, message)
		default:
			log.Error().Err(err).Msg("failed to parse the message")
			continue
		}
	}

	c.JSON(http.StatusOK, messages)
}

// Get the pagination size from the query string and return the
// default size if the given size is invalid or not exists
func (m *Monitor) PageSize(c *gin.Context) int {
	size := c.DefaultQuery("size", fmt.Sprintf("%d", DEFAULT_SIZE))
	switch size, err := strconv.Atoi(size); err {
	case nil:
		return size
	default:
		return DEFAULT_SIZE
	}
}

// Get the pagination after time from the query string and return
// the now if the given time is invalid or not exists
func (m *Monitor) PageAfter(c *gin.Context) time.Time {
	after := c.DefaultQuery("after", "")
	switch after, err := time.Parse("2006-01-02T15:04:05", after); err {
	case nil:
		return after
	default:
		return time.Now()
	}
}
