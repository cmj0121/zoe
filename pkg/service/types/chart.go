package types

import (
	"database/sql"
	"time"
)

// Chart is the data type of the chart
type Chart struct {
	Timestamp string
	Count     int64
}

func ChartFromRows(rows *sql.Rows) (*Chart, error) {
	var chart Chart
	var hour int64

	if err := rows.Scan(&chart.Count, &hour); err != nil {
		return nil, err
	}

	chart.Timestamp = time.Unix(hour*3600, 0).UTC().Format("2006-01-02T15")
	return &chart, nil
}
