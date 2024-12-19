package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"goserve/stats"
	"goserve/types"
	"goserve/utils"

	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase(config *types.Config) {
	// for PRODUCTION code credentials should be externalized source tree of course
	// dsn := "root:admin@tcp(127.0.0.1:3306)/freestyle?charset=utf8mb4&parseTime=True&loc=Local"

	// open DB connection for GORM
	database, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second, // Log queries slower than this threshold
				LogLevel:                  logger.Info, // Log level (Info will log the SQL queries)
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Enable color output for terminals
			},
		),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Failed to get sqlDB instance: ", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(time.Minute * 5)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = database
}

func GetGlucoseRecords(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fromDate := params["fromDate"] + " 00:00:00"
	toDate := params["toDate"] + " 23:59:59"
	localtime := time.Local

	// Parse dates to time.Time
	from, err := time.ParseInLocation("2006-01-02 15:04:05", fromDate, localtime) // Assuming dates are in "YYYY-MM-DD" format
	if err != nil {
		http.Error(w, "Invalid fromDate format", http.StatusBadRequest)
		return
	}
	to, err := time.ParseInLocation("2006-01-02 15:04:05", toDate, localtime)
	if err != nil {
		http.Error(w, "Invalid toDate format", http.StatusBadRequest)
		return
	}

	var results []types.Glucose

	if err := DB.Table("glucose").
		Where("device_timestamp BETWEEN ? AND ? AND record_type = ?", from, to, 0).
		Select("device_timestamp, historic_glucose_mg_dl").
		Find(&results).Error; err != nil {
		log.Println("Error querying glucose:", err)
		http.Error(w, "Error retrieving records", http.StatusInternalServerError)
		return
	}

	// Store each reading into time interval buckets
	buckets := stats.CollateBuckets(results)
	// Generate the statistical summary for each bucket
	summary := stats.SummarizeGlucose(buckets)

	// Send the results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

func GetRawData(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fromDate := params["fromDate"] + " 00:00:00"
	toDate := params["toDate"] + " 23:59:59"
	localtime := time.Local

	// Parse dates to time.Time
	from, err := time.ParseInLocation("2006-01-02 15:04:05", fromDate, localtime) // Assuming dates are in "YYYY-MM-DD" format
	if err != nil {
		http.Error(w, "Invalid fromDate format", http.StatusBadRequest)
		return
	}
	to, err := time.ParseInLocation("2006-01-02 15:04:05", toDate, localtime)
	if err != nil {
		http.Error(w, "Invalid toDate format", http.StatusBadRequest)
		return
	}

	var results []types.Glucose
	if err := DB.Table("glucose").
		Where("device_timestamp BETWEEN ? AND ? AND record_type = ?", from, to, 0).
		Select("device_timestamp, historic_glucose_mg_dl").
		Find(&results).Error; err != nil {
		log.Println("Error querying glucose:", err)
		http.Error(w, "Error retrieving records", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	var response []types.GlucoseResponse
	for _, result := range results {
		// Only include valid glucose readings
		if result.HistoricGlucoseMgDl.Valid {
			response = append(response, types.GlucoseResponse{
				DeviceTimestamp:     result.DeviceTimestamp.Format(time.RFC3339),
				HistoricGlucoseMgDl: result.HistoricGlucoseMgDl.Int64, // Extract Int64 value
			})
		}
	}

	// Send the results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}


func GetSerialNumbers() ([]types.GlucoseSummary, error) {
	var summaries []types.GlucoseSummary

	// Custom SQL query to get first and last timestamps for each serial_number
	query := `
			SELECT serial_number,
						 MIN(device_timestamp) AS first_timestamp,
						 MAX(device_timestamp) AS last_timestamp
			FROM glucose
			GROUP BY serial_number
	`

	// Execute the query and scan the results into the summaries slice
	if err := DB.Raw(query).Scan(&summaries).Error; err != nil {
			return nil, err
	}

	return summaries, nil
}

func GetSerialNumbersHandler(w http.ResponseWriter, r *http.Request) {
	summaries, err := GetSerialNumbers()
	if err != nil {
			http.Error(w, "Error retrieving glucose serial numbers", http.StatusInternalServerError)
			return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}


func main() {

	// Load configuration from file
	config, err := utils.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the database
	ConnectDatabase(config)
	// Automatically migrate your schema
	DB.AutoMigrate(&types.Glucose{})
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/api/glucose/{fromDate}/{toDate}", GetGlucoseRecords).Methods("GET")
	r.HandleFunc("/api/raw/{fromDate}/{toDate}", GetRawData).Methods("GET")
	r.HandleFunc("/api/serials", GetSerialNumbersHandler).Methods("GET")

	// Serve static files from the React build directory
	staticDir := "./react-app/build"
	fs := http.FileServer(http.Dir(staticDir))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Catch-all route to serve React app (client-side routing)
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, staticDir+"/index.html")
	})

	// Start the server
	log.Fatal(http.ListenAndServe(":8000", r))
}
