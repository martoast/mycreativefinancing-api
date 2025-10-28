package models

import (
	"strings"

	"api/pkg/config"

	"gorm.io/gorm"
)

var db *gorm.DB

type Property struct {
	gorm.Model
	Address                string   `json:"address"`
	Price                  *float64 `json:"price"`                    // DECIMAL(10, 2)
	Description            *string  `json:"description"`              // TEXT
	Images                 *string  `json:"images"`                   // JSON array of image URLs or a string of delimited URLs
	Sold                   *bool    `json:"sold"`                     // BOOLEAN
	Bedrooms               *int     `json:"bedrooms"`                 // INT
	Bathrooms              *float64 `json:"bathrooms"`                // DECIMAL(3, 1)
	RentZestimate          *float64 `json:"rent_zestimate"`           // DECIMAL(10, 2)
	Zestimate              *float64 `json:"zestimate"`                // DECIMAL(10, 2)
	PropertyType           *string  `json:"property_type"`            // VARCHAR(255)
	Zoning                 *string  `json:"zoning"`                   // VARCHAR(255)
	YearBuilt              *int     `json:"year_built"`               // INT
	LotSize                *int     `json:"lot_size"`                 // INT
	PricePerSquareFoot     *float64 `json:"price_per_square_foot"`    // DECIMAL(10, 2)
	LivingArea             *int     `json:"living_area"`              // INT
	PurchasePrice          *float64 `json:"purchase_price"`           // DECIMAL(10,2)
	BalanceToClose         *float64 `json:"balance_to_close"`         // DECIMAL(10,2)
	MonthlyHoldingCost     *float64 `json:"monthly_holding_cost"`     // DECIMAL(10,2)
	InterestRate           *float64 `json:"interest_rate"`            // DECIMAL(10,2)
	NearbyHospitals        *string  `json:"nearby_hospitals"`         // JSON array of nearby_hospitals
	NearbySchools          *string  `json:"nearby_schools"`           // JSON array of nearby_schools
	NearbyHomes            *string  `json:"nearby_homes"`             // JSON array of nearby_homes
	PriceHistory           *string  `json:"price_history"`            // JSON array of price history
	TaxHistory             *string  `json:"tax_history"`              // JSON array of tax history
	ContactRecipients      *string  `json:"contact_recipients"`       // JSON array of contact_recipients
	MonthlyHoaFee          *int     `json:"monthly_hoa_fee"`          // INT
	TransactionDocumentUrl string   `json:"transaction_document_url"` // string url of transaction document
	BenefitSheetUrl        string   `json:"benefit_sheet_url"`
	Escrow                 *float64 `json:"escrow"` // DECIMAL(10, 2)
	DealHolder             *string  `json:"deal_holder"`
	DealHolderPhone        *string  `json:"deal_holder_phone"`
	DealHolderEmail        *string  `json:"deal_holder_email"`
	AssignmentFee          *float64 `json:"assignment_fee"`
	InHouseDeal            *bool    `json:"in_house_deal"`       // BOOLEAN
	RentalRestriction      *bool    `json:"rental_restriction"`  // BOOLEAN
	PriceBreakDown         *string  `json:"price_breakdown"`     // VARCHAR(255)
	AdditionalBenefits     *string  `json:"additional_benefits"` // VARCHAR(255)
	CreatedBy              *string  `json:"created_by"`          // VARCHAR(255) - can be "user" or "admin"
	ReApiId                *string  `json:"re_api_id"`           // VARCHAR(255) - Real Estate API ID
}

func init() {
	config.Connect()
	db = config.GetDB()

	db.AutoMigrate(&Property{})

	// DeleteProperty(8)

	// Seed data if needed
	var count int64
	db.Model(&Property{}).Count(&count)
	if count == 0 {
		SeedProperties()
	}
}

func GetAllProperties() []Property {
	var Properties []Property
	db.Find(&Properties)
	return Properties
}

func GetPaginatedProperties(limit int, offset int, sold *bool) ([]Property, int64) {
	var properties []Property
	var total int64

	query := db.Model(&Property{})

	if sold != nil {
		query = query.Where("sold = ?", *sold)
	}

	query.Count(&total)
	query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&properties)

	return properties, total
}

func SearchProperties(query string, limit int, offset int, sold *bool) ([]Property, int64) {
	var properties []Property
	var total int64

	// Prepare the base query
	dbQuery := db.Model(&Property{})

	// Apply sold filter if provided
	if sold != nil {
		dbQuery = dbQuery.Where("sold = ?", *sold)
	}

	// Search across address field
	searchPattern := "%" + strings.ToLower(query) + "%"
	dbQuery = dbQuery.Where("LOWER(address) LIKE ?", searchPattern)

	// Get total count
	dbQuery.Count(&total)

	// Get paginated results
	dbQuery.Order("created_at DESC").Limit(limit).Offset(offset).Find(&properties)

	return properties, total
}

func GetPropertyById(ID int64) (*Property, *gorm.DB) {
	var getProperty Property
	db := db.Where("ID=?", ID).Find(&getProperty)
	return &getProperty, db
}

func (b *Property) CreateProperty() *Property {
	db.Create(&b)
	return b

}

func DeleteProperty(ID int64) Property {
	var property Property
	db.Where("ID = ?", ID).Delete(&property)
	return property
}

func SeedProperties() {
	// Example set of properties to seed
	properties := []Property{
		{
			Address:                "4949 Corrado Ave, Ave Maria, FL 34142",
			Price:                  newFloat64(300000),
			Description:            newString("Beautiful family home in a quiet neighborhood."),
			Images:                 newString("[\"https://static.tildacdn.com/stor3630-6334-4663-b532-393032356238/65960768.jpg\", \"https://static.tildacdn.com/stor3663-3339-4534-b332-393563363363/61347039.jpg\"]"),
			Sold:                   newBool(false),
			Bedrooms:               newInt(3),
			Bathrooms:              newFloat64(2.5),
			RentZestimate:          newFloat64(2500),
			Zestimate:              newFloat64(300000),
			PropertyType:           newString("Single Family"),
			Zoning:                 newString("R-1:SINGLE FAM-RES"),
			YearBuilt:              newInt(1990),
			LotSize:                newInt(5000),
			LivingArea:             newInt(3000),
			PricePerSquareFoot:     newFloat64(300),
			PurchasePrice:          newFloat64(300000),
			BalanceToClose:         newFloat64(10000),
			MonthlyHoldingCost:     newFloat64(5000),
			InterestRate:           newFloat64(300),
			NearbyHospitals:        newString("[\"Hospital A\", \"Hospital B\"]"),
			NearbySchools:          newString("[\"School A\", \"School B\"]"),
			NearbyHomes:            newString("[\"Home A\", \"Home B\"]"),
			PriceHistory:           newString("[{\"date\": \"2022-01-01\", \"price\": 295000}, {\"date\": \"2023-01-01\", \"price\": 300000}]"),
			TaxHistory:             newString("[{\"year\": 2022, \"tax\": 3500}, {\"year\": 2023, \"tax\": 3600}]"),
			ContactRecipients:      newString("[{\"agent_reason\":1,\"zpro\":null,\"recent_sales\":0,\"review_count\":8,\"display_name\":\"Elizabeth Jimenez\",\"zuid\":\"X1-ZU12a3ye9stjcw9_26nu5\",\"rating_average\":5,\"badge_type\":\"Premier Agent\",\"phone\":{\"prefix\":\"484\",\"areacode\":\"424\",\"number\":\"9901\"},\"image_url\":\"https://photos.zillowstatic.com/fp/a9702d055054a53bd296d7175519fb29-h_n.jpg\"}]"),
			MonthlyHoaFee:          newInt(1000),
			TransactionDocumentUrl: "https://docs.google.com/spreadsheets/d/1-Ot5O9Fh7mOVQa5SJieBGrU9rIaItGVyZmEXwz4aAJY/edit?gid=0#gid=0",
			PriceBreakDown:         newString("This is the price breakdown"),
			AdditionalBenefits:     newString("This is the additional benefits"),
			CreatedBy:              newString("admin"),
			ReApiId:                newString("4574363"),
		},
	}

	for _, property := range properties {
		db.Create(&property)
	}
}

func newFloat64(v float64) *float64 { return &v }
func newString(s string) *string    { return &s }
func newInt(i int) *int             { return &i }
func newBool(b bool) *bool          { return &b }
