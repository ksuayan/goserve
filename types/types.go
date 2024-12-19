package types

import (
	"database/sql"
	"time"
)

// PercentileSummary holds the percentiles, average, and the readings for a specific bucket
type PercentileSummary struct {
	BucketKey            string  `json:"timeslot"`
	Readings             []int64 `json:"readings"`
	FivePercentile       float64 `json:"pct_05"`
	NinetyFivePercentile float64 `json:"pct_95"`
	Average              float64 `json:"average"`
}

type Config struct {
	DSN string `json:"dsn"`
}

type Glucose struct {
	DeviceTimestamp     time.Time     `json:"device_timestamp"`
	HistoricGlucoseMgDl sql.NullInt64 `json:"historic_glucose_mg_dl"`
	RecordType          int           `json:"-"`
}

type GlucoseResponse struct {
	DeviceTimestamp     string `json:"timestamp"`
	HistoricGlucoseMgDl int64  `json:"glucose"`
}

type GlucoseSummary struct {
	SerialNumber   string    `json:"serial_number"`
	FirstTimestamp time.Time `json:"first_timestamp"`
	LastTimestamp  time.Time `json:"last_timestamp"`
}