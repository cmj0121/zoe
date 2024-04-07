package monitor

import (
	"database/sql"
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

// Do query the from SQL statement and return the result
func (m *Monitor) query(query string, args ...any) (*sql.Rows, error) {
	stmt, err := m.Prepare(query)

	if err != nil {
		log.Error().Err(err).Str("query", query).Msg("failed to prepare the statement")
		return nil, err
	}

	return stmt.Query(args...)
}

// List all the records with pagination
func (m *Monitor) listMessages(c *gin.Context) {
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
	rows, err := m.query(stmt, m.PageBefore(c).UnixNano(), m.PageSize(c))
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

// List the number of records per remote IP
func (m *Monitor) groupByRemote(c *gin.Context) {
	stmt := `
		SELECT
			COUNT(message.client_ip) AS count,
			message.client_ip
		FROM message
		GROUP BY message.client_ip
		ORDER BY count DESC
	`

	rows, err := m.query(stmt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type record struct {
		Count  int
		Remote string
	}
	records := []record{}

	for rows.Next() {
		var r record

		if err := rows.Scan(&r.Count, &r.Remote); err != nil {
			log.Error().Err(err).Msg("failed to scan the row")
			continue
		}

		records = append(records, r)
	}

	c.JSON(http.StatusOK, records)
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
func (m *Monitor) PageBefore(c *gin.Context) time.Time {
	before := c.DefaultQuery("before", "")
	switch before {
	case "":
		return time.Now()
	default:
		switch ns, err := strconv.Atoi(before); err {
		case nil:
			return time.Unix(int64(ns)/1e9, int64(ns)%1e9)
		default:
			return time.Now()
		}
	}
}
