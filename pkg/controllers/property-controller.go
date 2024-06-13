package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"api/pkg/models"
	"api/pkg/utils"

	"github.com/gorilla/mux"
)

var NewProperty models.Property

func GetProperties(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for page and pageSize
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	soldParam := r.URL.Query().Get("sold")

	// Default values if not provided
	page := 1
	pageSize := 10
	var sold *bool

	var err error
	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}

	if pageSizeStr != "" {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			http.Error(w, "Invalid pageSize parameter", http.StatusBadRequest)
			return
		}
	}

	if soldParam != "" {
		soldValue, err := strconv.ParseBool(soldParam)
		if err != nil {
			http.Error(w, "Invalid sold parameter", http.StatusBadRequest)
			return
		}
		sold = &soldValue
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Get paginated properties from the model
	newProperties, total := models.GetPaginatedProperties(pageSize, offset, sold)

	// Create response with properties and total count
	response := map[string]interface{}{
		"properties": newProperties,
		"total":      total,
		"page":       page,
		"pageSize":   pageSize,
	}

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode properties as JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetPropertyById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PropertyId := vars["PropertyId"]
	ID, err := strconv.ParseInt(PropertyId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing")
	}
	PropertyDetails, _ := models.GetPropertyById(ID)
	res, _ := json.Marshal(PropertyDetails)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateProperty(w http.ResponseWriter, r *http.Request) {
	PropertyModel := &models.Property{}
	utils.ParseBody(r, PropertyModel)
	b := PropertyModel.CreateProperty()
	res, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	PropertyId := vars["PropertyId"]
	ID, err := strconv.ParseInt(PropertyId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing")
		http.Error(w, "Invalid property ID", http.StatusBadRequest)
		return
	}
	property := models.DeleteProperty(ID)
	res, err := json.Marshal(property)
	if err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	var updateProperty = &models.Property{}
	utils.ParseBody(r, updateProperty)
	vars := mux.Vars(r)
	PropertyId := vars["PropertyId"]
	ID, err := strconv.ParseInt(PropertyId, 10, 64)
	if err != nil {
		fmt.Println("error while parsing")
	}
	propertyDetails, db := models.GetPropertyById(ID)

	propertyDetails.Price = updateProperty.Price
	propertyDetails.Description = updateProperty.Description
	propertyDetails.Images = updateProperty.Images
	propertyDetails.Sold = updateProperty.Sold
	propertyDetails.Bedrooms = updateProperty.Bedrooms
	propertyDetails.Bathrooms = updateProperty.Bathrooms
	propertyDetails.RentZestimate = updateProperty.RentZestimate
	propertyDetails.Zestimate = updateProperty.Zestimate
	propertyDetails.PropertyType = updateProperty.PropertyType
	propertyDetails.Zoning = updateProperty.Zoning
	propertyDetails.YearBuilt = updateProperty.YearBuilt
	propertyDetails.LotSize = updateProperty.LotSize
	propertyDetails.PricePerSquareFoot = updateProperty.PricePerSquareFoot
	propertyDetails.LivingArea = updateProperty.LivingArea
	propertyDetails.PurchasePrice = updateProperty.PurchasePrice
	propertyDetails.BalanceToClose = updateProperty.BalanceToClose
	propertyDetails.MonthlyHoldingCost = updateProperty.MonthlyHoldingCost
	propertyDetails.InterestRate = updateProperty.InterestRate
	propertyDetails.NearbyHospitals = updateProperty.NearbyHospitals
	propertyDetails.NearbySchools = updateProperty.NearbySchools
	propertyDetails.NearbyHomes = updateProperty.NearbyHomes
	propertyDetails.PriceHistory = updateProperty.PriceHistory
	propertyDetails.TaxHistory = updateProperty.TaxHistory
	propertyDetails.MonthlyHoaFee = updateProperty.MonthlyHoaFee
	propertyDetails.ContactRecipients = updateProperty.ContactRecipients

	db.Save(&propertyDetails)
	res, err := json.Marshal(propertyDetails)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
