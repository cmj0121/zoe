package types

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/cmj0121/zoe/pkg/database"
)

type Report struct {
	Value string `json:"value"`
	Count int    `json:"count"`
}

func ReportFromRow(rows *sql.Rows) (*Report, error) {
	var report Report

	err := rows.Scan(&report.Count, &report.Value)
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func DailyPopularMessages(ctx context.Context, field string, count int) []*Report {
	sess := database.Session()
	today := time.Now().Truncate(24 * time.Hour).Add(-24 * time.Hour).Format("2006-01-02")

	stmt := fmt.Sprintf(`
		SELECT
			COUNT(message.%[1]v) AS count,
			message.%[1]v AS value
		FROM message
		WHERE
			DATE(message.created_at) = ? AND message.%[1]v IS NOT NULL
		GROUP BY message.%[1]v
		ORDER BY count DESC
		LIMIT ?
	`, field)

	rows, err := sess.QueryContext(ctx, stmt, today, count)
	if err != nil {
		log.Warn().Err(err).Str("field", field).Msg("failed to query the popular messages")
		return nil
	}

	var reports []*Report
	for rows.Next() {
		report, err := ReportFromRow(rows)
		if err != nil {
			log.Warn().Err(err).Msg("failed to parse the popular message")
			continue
		}

		reports = append(reports, report)
	}

	return reports
}
