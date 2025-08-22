package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
)

// Database connection
var db *sql.DB

// Constants
const (
	CHUNK_SIZE = 10000 // Process 10k rows at a time to manage memory
)

// Allowed tables
var ALLOWED_TABLES = []string{
	"GTPL_108_gT_40E_P_S7_200_Germany",
	"GTPL_109_gT_40E_P_S7_200_Germany",
	"GTPL_110_gT_40E_P_S7_200_Germany",
	"GTPL_111_gT_80E_P_S7_200_Germany",
	"GTPL_112_gT_80E_P_S7_200_Germany",
	"GTPL_113_gT_80E_P_S7_200_Germany",
	"kabomachinedatasmart200",
	"GTPL_114_GT_140E_S7_1200",
	"GTPL_115_GT_180E_S7_1200",
	"GTPL_119_GT_180E_S7_1200",
	"GTPL_120_GT_180E_S7_1200",
	"GTPL_116_GT_240E_S7_1200",
	"GTPL_117_GT_320E_S7_1200",
	"GTPL_121_GT1000T",
	"gtpl_122_s7_1200_01",
	"GTPL_124_GT_450T_S7_1200",
	"GTPL_131_GT_650T_S7_1200",
	"GTPL_132_GT_650T_S7_1200",
}

// Preferred numeric column order
var PREFERRED_NUMERIC_ORDER = []string{
	"T2_1_ambient_temp", "T2_2_ambient_temp", "T2_temp_mean",
	"T1_1_cold_air_temp", "T1_2_cold_air_temp", "T1_temp_mean",
	"T0_1_air_outlet_temp", "T0_2_air_outlet_temp", "T0_temp_mean",
	"TH_1_supply_air_temp", "TH_2_supply_air_temp", "TH_temp_mean",
	"LP_value", "HP_value", "LP_set_point", "HP_set_point",
	"T1_set_point", "TH_T1_set_point", "Compressor_timer", "Delta_set_to_aeration", "Aeration_duration_set",
	"Running_time_hour", "Running_time_minute", "Running_hours", "Running_hours_min",
	"Blower_speed", "Hot_valve_speed", "AHT_vale_speed", "AHT_valve_speed", "Heater_speed", "Cond_fan_speed",
	"Blower_speed_set_in_manual", "Cond_fan_speed_set_in_manual", "Hot_gas_valve_set_in_manual", "AHT_valve_set_in_manual", "Heater_set_in_manual",
	"Fault_code", "FS", "UF", "RHP", "BLWR_pct", "RMR_pct", "CNPR_pct", "AHT_pct", "HCSR_pct",
}

// Pretty header mapping
var PRETTY_HEADER_MAP = map[string]string{
	"id":                           "Record#",
	"created_at":                   "Date & Time (IST)",
	"created_at_date":              "Date",
	"created_at_time":              "Time",
	"T2_1_ambient_temp":            "T2-1 Ambient Temp (°C)",
	"T2_2_ambient_temp":            "T2-2 Ambient Temp (°C)",
	"T1_1_cold_air_temp":           "T1-1 Cold Air Temp (°C)",
	"T1_2_cold_air_temp":           "T1-2 Cold Air Temp (°C)",
	"T0_1_air_outlet_temp":         "T0-1 Air Outlet Temp (°C)",
	"T0_2_air_outlet_temp":         "T0-2 Air Outlet Temp (°C)",
	"TH_1_supply_air_temp":         "TH-1 Supply Air Temp (°C)",
	"TH_2_supply_air_temp":         "TH-2 Supply Air Temp (°C)",
	"T2_temp_mean":                 "T2 Mean Temp (°C)",
	"T1_temp_mean":                 "T1 Mean Temp (°C)",
	"T0_temp_mean":                 "T0 Mean Temp (°C)",
	"TH_temp_mean":                 "TH Mean Temp (°C)",
	"LP_value":                     "LP Value",
	"HP_value":                     "HP Value",
	"T1_set_point":                 "T1 Set Point",
	"TH_T1_set_point":              "TH-T1 Set Point",
	"Compressor_timer":             "Compressor Timer (s)",
	"Delta_set_to_aeration":        "Delta to Aeration",
	"Aeration_duration_set":        "Aeration Duration Set",
	"Running_time_hour":            "Running Hours",
	"Running_time_minute":          "Running Minutes",
	"HP_set_point":                 "HP Set Point",
	"LP_set_point":                 "LP Set Point",
	"Blower_speed":                 "Blower Speed (%)",
	"Hot_valve_speed":              "Hot Valve Speed (%)",
	"AHT_vale_speed":               "AHT Valve Speed (%)",
	"AHT_valve_speed":              "AHT Valve Speed (%)",
	"Heater_speed":                 "Heater Speed (%)",
	"Cond_fan_speed":               "Cond Fan Speed (%)",
	"Blower_speed_set_in_manual":   "Blower Speed (Manual)",
	"Cond_fan_speed_set_in_manual": "Cond Fan Speed (Manual)",
	"Hot_gas_valve_set_in_manual":  "Hot Gas Valve (Manual)",
	"AHT_valve_set_in_manual":      "AHT Valve (Manual)",
	"Heater_set_in_manual":         "Heater (Manual)",
	"Running_hours":                "Total Running Hours",
	"Running_hours_min":            "Total Running Minutes",
	"Fault_code":                   "Fault Code",
	"UF":                           "UF*",
	"RHP":                          "RHP",
	"BLWR_pct":                     "BLWR%",
	"RMR_pct":                      "RMR%",
	"CNPR_pct":                     "CNPR%",
	"AHT_pct":                      "AHT%",
	"HCSR_pct":                     "HCSR%",
	"FS":                           "FS",
	"Faults":                       "Faults",
}

// Fault patterns
var FAULT_PATTERNS = []string{
	"fault", "overheat", "door_open", "short_circuit", "warning", "top", "protection", "not_achieved",
}

// DataRow represents a single row of data
type DataRow map[string]interface{}

// ExportRequest represents the export request parameters
type ExportRequest struct {
	Table    string `form:"table"`
	FromDate string `form:"fromDate"`
	ToDate   string `form:"toDate"`
	All      string `form:"all"`
	Limit    string `form:"limit"`
	Order    string `form:"order"`
}

// ExportResponse represents the export response
type ExportResponse struct {
	Error   string `json:"error,omitempty"`
	Details string `json:"details,omitempty"`
}

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	// Create router
	r := gin.Default()

	// Add CORS middleware
	r.Use(corsMiddleware())

	// Routes
	r.GET("/export", handleExport)
	r.OPTIONS("/export", handleOptions)
	r.GET("/tables", handleTables)
	r.GET("/status", handleStatus)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		// Check database connection
		err := db.Ping()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "unhealthy",
				"database": "disconnected",
				"error":    err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"database":  "connected",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// Get port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(r.Run(":" + port))
}

// Initialize database connection
func initDB() {
	// Get database connection details from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Use provided credentials as defaults
	if dbHost == "" {
		dbHost = "myshaa.com"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUser == "" {
		dbUser = "myshaa_kabu"
	}
	if dbPassword == "" {
		dbPassword = "T-Cyj;f5g1y6"
	}
	if dbName == "" {
		dbName = "myshaa_kabu"
	}

	// Create connection string with additional parameters for better performance
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci&timeout=30s&readTimeout=60s&writeTimeout=60s",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("Database connection established to %s:%s/%s", dbHost, dbPort, dbName)
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		c.Header("Cache-Control", "no-store, max-age=0")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// Handle OPTIONS request
func handleOptions(c *gin.Context) {
	c.Status(http.StatusOK)
}

// Handle tables list request
func handleTables(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"tables":    ALLOWED_TABLES,
		"count":     len(ALLOWED_TABLES),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Handle status request
func handleStatus(c *gin.Context) {
	// Check database connection
	err := db.Ping()
	dbStatus := "connected"
	if err != nil {
		dbStatus = "disconnected"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "running",
		"database":  dbStatus,
		"uptime":    time.Since(time.Now()).String(),
		"timestamp": time.Now().Format(time.RFC3339),
		"version":   "1.0.0",
	})
}

// Handle export request
func handleExport(c *gin.Context) {
	var req ExportRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, ExportResponse{Error: "Invalid request parameters"})
		return
	}

	// Validate table
	if req.Table == "" || !isTableAllowed(req.Table) {
		c.JSON(http.StatusBadRequest, ExportResponse{Error: "Invalid or missing table name"})
		return
	}

	// Set defaults
	if req.All == "" {
		req.All = "true"
	}
	if req.Order == "" {
		req.Order = "desc"
	}

	// Build WHERE clause
	whereClause, params := buildWhereClause(req.FromDate, req.ToDate)

	// Get total count
	totalCount, err := getTotalCount(req.Table, whereClause, params)
	if err != nil {
		log.Printf("Error getting count: %v", err)
		c.JSON(http.StatusInternalServerError, ExportResponse{Error: "Failed to get record count"})
		return
	}

	log.Printf("Total matching records: %d", totalCount)

	// Process data
	processedRows, err := processDataInChunks(req.Table, whereClause, params, req.Order, totalCount, req.All == "true")
	if err != nil {
		log.Printf("Error processing data: %v", err)
		c.JSON(http.StatusInternalServerError, ExportResponse{Error: "Failed to process data"})
		return
	}

	// Create Excel file
	excelBuffer, err := createExcelFile(processedRows, req.All == "true")
	if err != nil {
		log.Printf("Error creating Excel file: %v", err)
		c.JSON(http.StatusInternalServerError, ExportResponse{Error: "Failed to create Excel file"})
		return
	}

	// Generate filename
	filename := fmt.Sprintf("%s_%s_%drecords.xlsx", req.Table, time.Now().Format("2006-01-02"), len(processedRows))

	// Set response headers
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")

	// Send file
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelBuffer)
}

// Check if table is allowed
func isTableAllowed(table string) bool {
	for _, allowedTable := range ALLOWED_TABLES {
		if allowedTable == table {
			return true
		}
	}
	return false
}

// Build WHERE clause for date filtering
func buildWhereClause(fromDate, toDate string) (string, []interface{}) {
	var conditions []string
	var params []interface{}

	if fromDate != "" {
		conditions = append(conditions, "created_at >= CONCAT(?, ' 00:00:00')")
		params = append(params, fromDate)
	}

	if toDate != "" {
		conditions = append(conditions, "created_at < DATE_ADD(CONCAT(?, ' 00:00:00'), INTERVAL 1 DAY)")
		params = append(params, toDate)
	}

	if len(conditions) == 0 {
		return "", params
	}

	return " WHERE " + strings.Join(conditions, " AND "), params
}

// Get total count of records
func getTotalCount(table, whereClause string, params []interface{}) (int, error) {
	query := fmt.Sprintf("SELECT COUNT(*) AS cnt FROM `%s`%s", table, whereClause)

	var count int
	err := db.QueryRow(query, params...).Scan(&count)
	return count, err
}

// Process data in chunks
func processDataInChunks(table, whereClause string, params []interface{}, order string, totalCount int, pretty bool) ([]DataRow, error) {
	var allProcessedRows []DataRow
	offset := 0

	// Determine order
	orderClause := "DESC"
	if strings.ToLower(order) == "asc" {
		orderClause = "ASC"
	}

	for offset < totalCount {
		currentChunkSize := CHUNK_SIZE
		if offset+currentChunkSize > totalCount {
			currentChunkSize = totalCount - offset
		}

		// Query chunk
		query := fmt.Sprintf("SELECT * FROM `%s`%s ORDER BY id %s LIMIT ? OFFSET ?", table, whereClause, orderClause)
		chunkParams := append(params, currentChunkSize, offset)

		rows, err := db.Query(query, chunkParams...)
		if err != nil {
			return nil, fmt.Errorf("error querying chunk: %v", err)
		}

		// Process chunk
		chunkRows, err := processChunk(rows, pretty)
		rows.Close()
		if err != nil {
			return nil, fmt.Errorf("error processing chunk: %v", err)
		}

		allProcessedRows = append(allProcessedRows, chunkRows...)
		offset += currentChunkSize

		// Log progress
		if offset%(CHUNK_SIZE*5) == 0 {
			log.Printf("Processed %d / %d records", offset, totalCount)
		}
	}

	return allProcessedRows, nil
}

// Process a single chunk of data
func processChunk(rows *sql.Rows, pretty bool) ([]DataRow, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var chunkRows []DataRow

	for rows.Next() {
		// Create slice to hold values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Convert to map
		row := make(DataRow)
		for i, col := range columns {
			val := values[i]
			if val == nil {
				row[col] = ""
			} else {
				row[col] = val
			}
		}

		// Process row based on format
		if pretty {
			row = processPrettyRow(row)
		} else {
			row = processRawRow(row)
		}

		chunkRows = append(chunkRows, row)
	}

	return chunkRows, nil
}

// Process row for pretty format
func processPrettyRow(row DataRow) DataRow {
	// Normalize created_at
	createdAt := normalizeCreatedAt(row["created_at"])

	processed := DataRow{
		"id":              row["id"],
		"created_at":      createdAt["full"],
		"created_at_date": createdAt["date"],
		"created_at_time": createdAt["time"],
	}

	// Get numeric keys
	var numericKeys []string
	for k, v := range row {
		if k == "id" || k == "created_at" {
			continue
		}
		if toNum(v) != "" {
			numericKeys = append(numericKeys, k)
		}
	}

	// Order numeric columns
	ordered := orderNumericColumns(numericKeys)
	for _, k := range ordered {
		if v, exists := row[k]; exists {
			processed[k] = toNum(v)
		}
	}

	// Process faults
	faults := extractFaults(row)
	processed["Faults"] = faults

	return processed
}

// Process row for raw format
func processRawRow(row DataRow) DataRow {
	createdAt := normalizeCreatedAt(row["created_at"])
	row["created_at"] = createdAt["full"]
	return row
}

// Normalize created_at timestamp
func normalizeCreatedAt(raw interface{}) map[string]string {
	result := map[string]string{
		"full": "",
		"date": "",
		"time": "",
	}

	if raw == nil {
		return result
	}

	// Convert to string
	var timeStr string
	switch v := raw.(type) {
	case time.Time:
		timeStr = v.Format("2006-01-02 15:04:05")
	case string:
		timeStr = strings.Replace(v, "T", " ", 1)
		if len(timeStr) > 19 {
			timeStr = timeStr[:19]
		}
	default:
		return result
	}

	if timeStr == "" {
		return result
	}

	parts := strings.Split(timeStr, " ")
	if len(parts) >= 2 {
		result["full"] = timeStr
		result["date"] = parts[0]
		result["time"] = parts[1]
	} else if len(parts) == 1 {
		result["full"] = parts[0]
		result["date"] = parts[0]
	}

	return result
}

// Order numeric columns based on preference
func orderNumericColumns(numericKeys []string) []string {
	ordered := make([]string, 0, len(numericKeys))

	// Add preferred columns first
	for _, preferred := range PREFERRED_NUMERIC_ORDER {
		for _, key := range numericKeys {
			if key == preferred {
				ordered = append(ordered, key)
				break
			}
		}
	}

	// Add remaining columns
	for _, key := range numericKeys {
		found := false
		for _, orderedKey := range ordered {
			if key == orderedKey {
				found = true
				break
			}
		}
		if !found {
			ordered = append(ordered, key)
		}
	}

	return ordered
}

// Extract fault information
func extractFaults(row DataRow) string {
	var faults []string

	for k, v := range row {
		if !looksLikeFaultKey(k) {
			continue
		}
		if isTrueish(v) {
			faults = append(faults, strings.ReplaceAll(k, "_", " "))
		}
	}

	return strings.Join(faults, ", ")
}

// Check if key looks like a fault key
func looksLikeFaultKey(key string) bool {
	keyLower := strings.ToLower(key)
	for _, pattern := range FAULT_PATTERNS {
		if strings.Contains(keyLower, pattern) {
			return true
		}
	}
	return false
}

// Check if value is true-ish
func isTrueish(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		return strings.ToLower(val) == "true" || val == "1"
	case int:
		return val == 1
	case float64:
		return val == 1
	default:
		return false
	}
}

// Convert value to number
func toNum(v interface{}) interface{} {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		if val == "" {
			return ""
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
		return ""
	case int, int64, float64:
		return val
	default:
		return ""
	}
}

// Create Excel file
func createExcelFile(rows []DataRow, pretty bool) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()

	if len(rows) == 0 {
		// Create empty sheet
		f.SetCellValue("Sheet1", "A1", "id")
		f.SetCellValue("Sheet1", "B1", "created_at")
		f.SetCellValue("Sheet1", "A2", "No records found for selected criteria")
	} else {
		// Get headers
		var headers []string
		if pretty {
			headers = getPrettyHeaders(rows)
		} else {
			headers = getRawHeaders(rows)
		}

		// Write headers
		for i, header := range headers {
			prettyHeader := PRETTY_HEADER_MAP[header]
			if prettyHeader == "" {
				prettyHeader = header
			}
			cell := fmt.Sprintf("%c1", 'A'+i)
			f.SetCellValue("Sheet1", cell, prettyHeader)
		}

		// Write data
		for rowIdx, row := range rows {
			for colIdx, header := range headers {
				if value, exists := row[header]; exists {
					cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+2)
					f.SetCellValue("Sheet1", cell, value)
				}
			}
		}

		// Auto-fit columns
		for i := range headers {
			col := string(rune('A' + i))
			f.SetColWidth("Sheet1", col, col, 15)
		}
	}

	// Rename sheet
	f.SetSheetName("Sheet1", "Data")

	// Write to buffer
	return f.WriteToBuffer()
}

// Get headers for pretty format
func getPrettyHeaders(rows []DataRow) []string {
	if len(rows) == 0 {
		return []string{"id", "created_at", "created_at_date", "created_at_time"}
	}

	fixed := []string{"id", "created_at", "created_at_date", "created_at_time"}

	// Get all numeric keys
	allKeys := make(map[string]bool)
	for _, row := range rows {
		for k := range row {
			if k != "id" && k != "created_at" && k != "created_at_date" && k != "created_at_time" && k != "Faults" {
				allKeys[k] = true
			}
		}
	}

	// Convert to slice
	var dynamic []string
	for k := range allKeys {
		dynamic = append(dynamic, k)
	}

	// Order dynamic keys
	orderedDynamic := orderNumericColumns(dynamic)

	return append(fixed, append(orderedDynamic, "Faults")...)
}

// Get headers for raw format
func getRawHeaders(rows []DataRow) []string {
	if len(rows) == 0 {
		return []string{"id", "created_at"}
	}

	headers := make([]string, 0, len(rows[0]))
	for k := range rows[0] {
		headers = append(headers, k)
	}
	return headers
}
