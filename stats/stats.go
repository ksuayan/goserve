// stats/stats.go
package stats

import (
	"goserve/types"
	"goserve/utils"
	"math"
	"sort"
)

// Function to calculate the 5th and 95th percentiles and average for each time bucket
func SummarizeGlucose(buckets map[string][]int64) []types.PercentileSummary {
	var summary []types.PercentileSummary

	// Sort the buckets by time
	var keys []string
	for key := range buckets {
		keys = append(keys, key)
	}
	// Sort the keys
	sort.Strings(keys)

	// Iterate over the sorted keys
	for _, key := range keys {
		readings := buckets[key]
		if len(readings) > 0 {
			fivePercentile := Percentile(readings, 5)
			ninetyFivePercentile := Percentile(readings, 95)
			average := AverageReading(readings)

			// Sort the readings
			/*
				sortedReadings := make([]int64, len(readings))
				copy(sortedReadings, readings)
				sort.Slice(sortedReadings, func(i, j int) bool {
					return sortedReadings[i] < sortedReadings[j]
				})
				// Readings: sortedReadings,
			*/

			// Append the summary for the current bucket
			summary = append(summary, types.PercentileSummary{
				BucketKey:            key,
				FivePercentile:       fivePercentile,
				NinetyFivePercentile: ninetyFivePercentile,
				Average:              average,
			})
		}
	}
	return summary
}

// AverageReading calculates the average of the readings
func AverageReading(data []int64) float64 {
	if len(data) == 0 {
		return math.NaN() // or some other value indicating no data
	}
	sum := int64(0)
	for _, v := range data {
		sum += v
	}
	return float64(sum) / float64(len(data))
}

// Percentile calculates the nth percentile of a sorted slice of int64
func Percentile(data []int64, n int) float64 {
	if len(data) == 0 {
		return math.NaN() // or some other value indicating no data
	}
	sorted := make([]int64, len(data))
	copy(sorted, data)

	// Sort the data
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := (n * len(sorted)) / 100
	if index >= len(sorted) {
		return float64(sorted[len(sorted)-1])
	}
	return float64(sorted[index])
}

// collates an arrage of GlucoseData into readings
// grouped by a bucketKey (which is a time of day 5 minute interval)
func CollateBuckets(data []types.Glucose) map[string][]int64 {
	// Organize the data into buckets
	buckets := make(map[string][]int64)
	for _, record := range data {
		if record.HistoricGlucoseMgDl.Valid {
			bucketKey := utils.GetBucketKey(record.DeviceTimestamp)
			buckets[bucketKey] = append(buckets[bucketKey], record.HistoricGlucoseMgDl.Int64)
		}
	}
	return buckets
}
