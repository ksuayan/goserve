package utils

import (
	"encoding/json"
	"goserve/types"
	"os"
	"time"
)

func GetBucketKey(timestamp time.Time) string {
	hour := timestamp.Hour()
	minute := timestamp.Minute()

	// Round down the minute to the nearest 5-minute interval
	minuteBucket := (minute / 5) * 5

	// Create a new time.Time for the start of the bucket
	bucketStartTime := time.Date(
		timestamp.Year(), 
		timestamp.Month(), 
		timestamp.Day(), 
		hour, minuteBucket, 0, 0, timestamp.Location())
	return bucketStartTime.Format("15:04:05") // Format to "HH:MM:SS"
}

func LoadConfig(file string) (*types.Config, error) {
	config := &types.Config{}
	// Read file content
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	// Unmarshal JSON data into config struct
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
