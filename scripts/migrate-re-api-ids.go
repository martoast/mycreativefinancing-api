package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// ============================================================================
// CONFIGURATION - Change these values as needed
// ============================================================================
const (
	REAL_ESTATE_API_KEY = "URCREATIVESERVICESPRDATAPRODUCTION-d849-76f1-94d4-2a4c626d9bfe"
	REAL_ESTATE_API_URL = "https://api.realestateapi.com/v2/AutoComplete"

	// PRODUCTION API URL - Change this to your production endpoint
	API_BASE_URL = "https://shark-app-gfe6f.ondigitalocean.app"

	// FORCE_UPDATE - Set to true to re-fetch IDs even if properties already have them
	// Set to false to only update properties with null re_api_id
	FORCE_UPDATE = false
)

// ============================================================================
// Logger with timestamps
// ============================================================================
type Logger struct {
	*log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.Printf("[%s] INFO: %s", timestamp, fmt.Sprintf(format, v...))
}

func (l *Logger) Error(format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.Printf("[%s] ERROR: %s", timestamp, fmt.Sprintf(format, v...))
}

func (l *Logger) Success(format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.Printf("[%s] SUCCESS: %s", timestamp, fmt.Sprintf(format, v...))
}

func (l *Logger) Warning(format string, v ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	l.Printf("[%s] WARNING: %s", timestamp, fmt.Sprintf(format, v...))
}

var logger = NewLogger()

// ============================================================================
// Data Structures
// ============================================================================
type Property struct {
	ID                     uint     `json:"ID"`
	CreatedAt              string   `json:"CreatedAt,omitempty"`
	UpdatedAt              string   `json:"UpdatedAt,omitempty"`
	DeletedAt              *string  `json:"DeletedAt,omitempty"`
	Address                string   `json:"address"`
	Price                  *float64 `json:"price"`
	Description            *string  `json:"description"`
	Images                 *string  `json:"images"`
	Sold                   *bool    `json:"sold"`
	Bedrooms               *int     `json:"bedrooms"`
	Bathrooms              *float64 `json:"bathrooms"`
	RentZestimate          *float64 `json:"rent_zestimate"`
	Zestimate              *float64 `json:"zestimate"`
	PropertyType           *string  `json:"property_type"`
	Zoning                 *string  `json:"zoning"`
	YearBuilt              *int     `json:"year_built"`
	LotSize                *int     `json:"lot_size"`
	PricePerSquareFoot     *float64 `json:"price_per_square_foot"`
	LivingArea             *int     `json:"living_area"`
	PurchasePrice          *float64 `json:"purchase_price"`
	BalanceToClose         *float64 `json:"balance_to_close"`
	MonthlyHoldingCost     *float64 `json:"monthly_holding_cost"`
	InterestRate           *float64 `json:"interest_rate"`
	NearbyHospitals        *string  `json:"nearby_hospitals"`
	NearbySchools          *string  `json:"nearby_schools"`
	NearbyHomes            *string  `json:"nearby_homes"`
	PriceHistory           *string  `json:"price_history"`
	TaxHistory             *string  `json:"tax_history"`
	ContactRecipients      *string  `json:"contact_recipients"`
	MonthlyHoaFee          *int     `json:"monthly_hoa_fee"`
	TransactionDocumentUrl string   `json:"transaction_document_url"`
	BenefitSheetUrl        string   `json:"benefit_sheet_url"`
	Escrow                 *float64 `json:"escrow"`
	DealHolder             *string  `json:"deal_holder"`
	DealHolderPhone        *string  `json:"deal_holder_phone"`
	DealHolderEmail        *string  `json:"deal_holder_email"`
	AssignmentFee          *float64 `json:"assignment_fee"`
	InHouseDeal            *bool    `json:"in_house_deal"`
	RentalRestriction      *bool    `json:"rental_restriction"`
	PriceBreakDown         *string  `json:"price_breakdown"`
	AdditionalBenefits     *string  `json:"additional_benefits"`
	CreatedBy              *string  `json:"created_by"`
	ReApiId                *string  `json:"re_api_id"`
}

type PropertiesResponse struct {
	Properties []Property `json:"properties"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
}

type RealEstateAPIRequest struct {
	Search      string   `json:"search"`
	SearchTypes []string `json:"search_types"`
}

type RealEstateAPIResponse struct {
	Data []struct {
		ID         string  `json:"id"`
		Address    string  `json:"address"`
		SearchType string  `json:"searchType"`
		City       string  `json:"city"`
		State      string  `json:"state"`
		Zip        string  `json:"zip"`
		Latitude   float64 `json:"latitude"`
		Longitude  float64 `json:"longitude"`
	} `json:"data"`
	TotalResults  int    `json:"totalResults"`
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage"`
}

// ============================================================================
// Helper Functions
// ============================================================================

// Parse address to get just street + city for better API results
func parseAddressForSearch(fullAddress string) string {
	parts := strings.Split(fullAddress, ",")
	if len(parts) >= 2 {
		return strings.TrimSpace(parts[0]) + ", " + strings.TrimSpace(parts[1])
	}
	return fullAddress
}

// Fetch Real Estate API ID for an address
func fetchRealEstateApiId(address string) (string, error) {
	searchAddress := parseAddressForSearch(address)
	logger.Info("  → Searching Real Estate API: \"%s\"", searchAddress)

	reqBody := RealEstateAPIRequest{
		Search:      searchAddress,
		SearchTypes: []string{"A"}, // Only full addresses
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", REAL_ESTATE_API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", REAL_ESTATE_API_KEY)
	req.Header.Set("x-user-id", "migration-production-script")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResp RealEstateAPIResponse
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	if len(apiResp.Data) > 0 {
		// Find exact address match (searchType: 'A')
		for _, item := range apiResp.Data {
			if item.SearchType == "A" && item.ID != "" {
				logger.Success("  ✓ Found ID: %s", item.ID)
				return item.ID, nil
			}
		}
		// Fallback to first result
		if apiResp.Data[0].ID != "" {
			logger.Success("  ✓ Found ID (first result): %s", apiResp.Data[0].ID)
			return apiResp.Data[0].ID, nil
		}
	}

	return "", fmt.Errorf("no ID found")
}

// Fetch all properties from API
func fetchAllProperties() ([]Property, error) {
	url := fmt.Sprintf("%s/properties?page=1&pageSize=200", API_BASE_URL)
	logger.Info("Fetching properties from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch properties: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var propertiesResp PropertiesResponse
	if err := json.Unmarshal(bodyBytes, &propertiesResp); err != nil {
		return nil, fmt.Errorf("failed to decode properties: %v", err)
	}

	return propertiesResp.Properties, nil
}

// Update property with new re_api_id
func updateProperty(property Property, reApiId string) error {
	url := fmt.Sprintf("%s/properties/%d", API_BASE_URL, property.ID)

	oldId := "null"
	if property.ReApiId != nil && *property.ReApiId != "" {
		oldId = *property.ReApiId
	}
	property.ReApiId = &reApiId

	if oldId != "null" {
		logger.Info("  → Updating: %s → %s", oldId, reApiId)
	} else {
		logger.Info("  → Setting: %s", reApiId)
	}

	jsonData, err := json.Marshal(property)
	if err != nil {
		return fmt.Errorf("failed to marshal property: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Verify update
	var updatedProperty Property
	if err := json.Unmarshal(bodyBytes, &updatedProperty); err == nil {
		if updatedProperty.ReApiId != nil && *updatedProperty.ReApiId == reApiId {
			logger.Success("  ✓ Updated successfully")
			return nil
		}
	}

	return fmt.Errorf("update verification failed")
}

// ============================================================================
// Main Function
// ============================================================================
func main() {
	logger.Info("🚀 Starting Real Estate API ID PRODUCTION Migration")
	logger.Info("Configuration:")
	logger.Info("  - API: %s", API_BASE_URL)
	logger.Info("  - Force Update: %v", FORCE_UPDATE)
	fmt.Println()

	// Fetch all properties
	logger.Info("📥 Fetching properties...")
	properties, err := fetchAllProperties()
	if err != nil {
		logger.Error("FATAL: %v", err)
		os.Exit(1)
	}

	if len(properties) == 0 {
		logger.Info("No properties found. Exiting.")
		return
	}

	logger.Success("Found %d properties", len(properties))
	fmt.Println()

	// Stats
	successful := 0
	failed := 0
	skipped := 0
	failedProperties := []string{}

	// Process each property
	for i, property := range properties {
		fmt.Println(strings.Repeat("-", 80))
		logger.Info("[%d/%d] Property ID %d: %s", i+1, len(properties), property.ID, property.Address)

		// Skip if empty address
		if property.Address == "" {
			logger.Warning("  ⊘ Empty address - Skipping")
			skipped++
			continue
		}

		// Check if already has re_api_id
		if property.ReApiId != nil && *property.ReApiId != "" {
			if FORCE_UPDATE {
				logger.Info("  ⚠️  Has ID: %s - Force updating...", *property.ReApiId)
			} else {
				logger.Info("  ⊘ Already has ID: %s - Skipping", *property.ReApiId)
				skipped++
				continue
			}
		}

		// Fetch Real Estate API ID
		reApiId, err := fetchRealEstateApiId(property.Address)
		if err != nil {
			logger.Error("  ✗ Failed: %v", err)
			failed++
			failedProperties = append(failedProperties, fmt.Sprintf("ID %d: %s - %v", property.ID, property.Address, err))

			// Rate limit delay even on failure
			if i < len(properties)-1 {
				time.Sleep(1 * time.Second)
			}
			continue
		}

		// Update the property
		if err := updateProperty(property, reApiId); err != nil {
			logger.Error("  ✗ Update failed: %v", err)
			failed++
			failedProperties = append(failedProperties, fmt.Sprintf("ID %d: %s - update failed", property.ID, property.Address))
		} else {
			successful++
		}

		// Rate limiting: 1 second between requests
		if i < len(properties)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	// Final Summary
	fmt.Println()
	fmt.Println(strings.Repeat("=", 80))
	logger.Info("📊 MIGRATION COMPLETE")
	fmt.Println(strings.Repeat("=", 80))
	logger.Info("Total: %d properties", len(properties))
	logger.Success("✓ Successfully updated: %d", successful)
	logger.Error("✗ Failed: %d", failed)
	logger.Info("⊘ Skipped: %d", skipped)

	if len(failedProperties) > 0 {
		fmt.Println()
		logger.Error("Failed Properties:")
		for _, prop := range failedProperties {
			logger.Error("  - %s", prop)
		}
	}

	fmt.Println(strings.Repeat("=", 80))

	if failed > 0 {
		fmt.Println()
		logger.Warning("⚠️  Completed with %d errors", failed)
		os.Exit(1)
	} else {
		fmt.Println()
		logger.Success("✅ All properties migrated successfully!")
	}
}
