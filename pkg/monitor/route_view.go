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
func (m *Monitor) queryMessages(filter string, args ...any) ([]*types.Message, error) {
	stmt := fmt.Sprintf(`
		SELECT
			message.service,
			message.client_ip,
			message.username,
			message.password,
			message.command,
			message.created_at
		FROM message
		WHERE %v
		ORDER BY created_at DESC
		LIMIT ?
	`, filter)

	// limit the query result
	args = append(args, 20)

	rows, err := m.Query(stmt, args...)
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

func (m *Monitor) queryGroupBy(field, filter string, args ...any) ([]*types.GroupBy, error) {
	stmt := fmt.Sprintf(`
		SELECT
			message.%[1]v,
			COUNT(message.%[1]v) AS count,
			MAX(message.created_at) AS last_seen
		FROM message
		WHERE %[2]v
		GROUP BY message.%[1]v
		ORDER BY count DESC
		LIMIT ?
	`, field, filter)

	// limit the query result
	args = append(args, 20)

	rows, err := m.Query(stmt, args...)
	if err != nil {
		log.Warn().Err(err).Msg("failed to query the group_by")
		return nil, err
	}

	var group_by []*types.GroupBy
	for rows.Next() {
		switch gb, err := types.GroupByFromRows(rows); err {
		case nil:
			group_by = append(group_by, gb)
		default:
			log.Warn().Err(err).Msg("failed to parse the group_by")
			continue
		}
	}

	return group_by, nil
}

// Show the index page of the monitor service
func (m *Monitor) index(c *gin.Context) {
	now_ns := time.Now().UnixNano()
	args := []any{now_ns}

	filter := "created_at < ?"

	// process the query filter
	if client_ip, ok := c.GetQuery("client_ip"); ok {
		filter += " AND client_ip = ?"
		args = append(args, client_ip)
	}

	if username, ok := c.GetQuery("username"); ok {
		filter += " AND username = ?"
		args = append(args, username)
	}

	if password, ok := c.GetQuery("password"); ok {
		filter += " AND password = ?"
		args = append(args, password)
	}

	if command, ok := c.GetQuery("command"); ok {
		filter += " AND command = ?"
		args = append(args, command)
	}

	switch messages, err := m.queryMessages(filter, args...); err {
	case nil:
		c.HTML(http.StatusOK, "index.htm", gin.H{
			"year":     time.Now().Year(),
			"messages": messages,
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Show the group-by page of the monitor service
func (m *Monitor) group_by(c *gin.Context) {
	var field string
	switch field = c.Param("field"); field {
	case "client_ip", "username", "password", "command":
	default:
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusNotFound, "404 page not found")
		return
	}

	duration := time.Duration(time.Hour * 24 * 30)
	now_ns := time.Now().Add(-duration).UnixNano()
	filter := "created_at > ?"
	args := []any{now_ns}

	switch group_by, err := m.queryGroupBy(field, filter, args...); err {
	case nil:
		c.HTML(http.StatusOK, "group_by.htm", gin.H{
			"year":     time.Now().Year(),
			"fields":   []string{"client_ip", "username", "password", "command"},
			"field":    field,
			"group_by": group_by,
			"duration": fmt.Sprintf("%v days", duration.Hours()/24),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

// Show the chart page of the monitor service
func (m *Monitor) chart(c *gin.Context) {
	var field string

	stmt := `
		SELECT
			COUNT(*),
			created_at / 1000000000 / 60 / 60 AS hour
		FROM message
		WHERE service = ?
		GROUP BY hour
		ORDER BY hour DESC;
	`

	switch field = c.Param("field"); field {
	case "ssh", "form", "shell":
	default:
		c.Header("Content-Type", "text/plain")
		c.String(http.StatusNotFound, "404 page not found")
		return
	}

	rows, err := m.Query(stmt, field)
	if err != nil {
		log.Warn().Err(err).Msg("failed to query the chart data")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	var charts []*types.Chart
	for rows.Next() {
		switch chart, err := types.ChartFromRows(rows); err {
		case nil:
			charts = append(charts, chart)
		default:
			log.Warn().Err(err).Msg("failed to parse the chart data")
			continue
		}
	}

	c.HTML(http.StatusOK, "chart.htm", gin.H{
		"year":   time.Now().Year(),
		"fields": []string{"ssh", "form", "shell"},
		"charts": charts,
	})
}

// Get the static file from the embed.FS
func (m *Monitor) static(c *gin.Context) {
	filepath := fmt.Sprintf("web/static/%v", c.Param("filepath"))
	c.FileFromFS(filepath, http.FS(static))
}
