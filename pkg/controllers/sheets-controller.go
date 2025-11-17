package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// ORDERED_FIELDS defines the exact order of columns in the Google Sheet
var ORDERED_FIELDS = []string{
	"ID", "CreatedAt", "UpdatedAt", "DeletedAt", "address", "price",
	"description", "images", "sold", "bedrooms", "bathrooms",
	"rent_zestimate", "zestimate", "property_type", "zoning",
	"year_built", "lot_size", "price_per_square_foot", "living_area",
	"purchase_price", "balance_to_close", "monthly_holding_cost",
	"interest_rate", "nearby_hospitals", "nearby_schools", "nearby_homes",
	"price_history", "tax_history", "contact_recipients", "monthly_hoa_fee",
	"transaction_document_url", "benefit_sheet_url", "escrow",
	"deal_holder", "in_house_deal", "rental_restriction",
	"price_breakdown", "additional_benefits",
}

// PropertyToSheetRequest represents the request body
type PropertyToSheetRequest struct {
	Property map[string]interface{} `json:"property"`
}

// safeToString converts any value to string safely
func safeToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case float64:
		return fmt.Sprintf("%v", v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		// For arrays and objects, try to marshal to JSON
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			fmt.Printf("Warning: Failed to stringify value: %v\n", err)
			return ""
		}
		return string(jsonBytes)
	}
}

// WritePropertyToSheet handles writing property data to Google Sheets
func WritePropertyToSheet(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req PropertyToSheetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	property := req.Property
	if property == nil || property["address"] == nil {
		http.Error(w, "No property data or address provided", http.StatusBadRequest)
		return
	}

	// Get credentials from environment variable
	credentialsJSON := os.Getenv("GOOGLE_SHEETS_API_CREDENTIALS")
	if credentialsJSON == "" {
		http.Error(w, "Google Sheets credentials not configured", http.StatusInternalServerError)
		return
	}

	// Parse credentials
	config, err := google.JWTConfigFromJSON(
		[]byte(credentialsJSON),
		sheets.SpreadsheetsScope,
	)
	if err != nil {
		fmt.Printf("Error parsing credentials: %v\n", err)
		http.Error(w, "Failed to parse Google Sheets credentials", http.StatusInternalServerError)
		return
	}

	// Create client
	ctx := context.Background()
	client := config.Client(ctx)

	// Create Sheets service
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Printf("Error creating Sheets service: %v\n", err)
		http.Error(w, "Failed to create Sheets service", http.StatusInternalServerError)
		return
	}

	// Format data in correct order
	var values []interface{}
	for _, field := range ORDERED_FIELDS {
		values = append(values, safeToString(property[field]))
	}

	// Prepare the data to append
	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{values},
	}

	// The spreadsheet ID - you might want to make this configurable
	spreadsheetID := "1-KORQc7eQHidXeZ5cVuA5VR73NfeRQ_IKVO8nkkWEWM"

	// Append the data
	resp, err := srv.Spreadsheets.Values.Append(
		spreadsheetID,
		"Sheet1",
		valueRange,
	).ValueInputOption("USER_ENTERED").Do()

	if err != nil {
		fmt.Printf("Error appending to sheet: %v\n", err)

		// Provide more specific error messages
		errorMessage := "Failed to write to Google Sheet"
		statusCode := http.StatusInternalServerError

		if err.Error() == "googleapi: Error 403" {
			errorMessage = "Permission denied. Check spreadsheet ID and credentials."
			statusCode = http.StatusForbidden
		} else if err.Error() == "googleapi: Error 404" {
			errorMessage = "Spreadsheet not found. Check spreadsheet ID."
			statusCode = http.StatusNotFound
		}

		http.Error(w, errorMessage, statusCode)
		return
	}

	// Success response
	response := map[string]interface{}{
		"message":      "Data written to sheet successfully",
		"updatedRange": resp.Updates.UpdatedRange,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
