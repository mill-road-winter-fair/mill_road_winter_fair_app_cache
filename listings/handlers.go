package listings

import (
	"net/http"

	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/shared"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

func GetListings(c *gin.Context) {

	//Call the DB function
	glog.Info("Calling GetListingsFromDb function")
	dbStatus, details, rows, err := GetListingsFromDb(shared.Database)

	if err != nil {
		//Return the status code and body from the function
		c.JSON(dbStatus, gin.H{"status": details, "error": err.Error()})
		return
	}

	//Ensure rows are closed after processing
	defer rows.Close()

	var rowsJson []shared.ListingData
	for rows.Next() {
		var row shared.ListingData
		//Map database row values to the struct fields
		if err := rows.Scan(&row.Id, &row.Name, &row.DisplayName, &row.PrimaryType, &row.SecondaryType, &row.TertiaryType, &row.Email, &row.Website, &row.Phone, &row.PlusCode, &row.StartTime, &row.EndTime); err != nil {
			glog.Fatalf("Scanning rows failed: %v", err)
			c.JSON(500, gin.H{"status": "error", "error": "Failed to scan row data"})
			return
		}
		rowsJson = append(rowsJson, row)
	}

	//Return successful response with the processed JSON data
	c.JSON(dbStatus, rowsJson)
	return
}

func CreateListing(c *gin.Context) {
	//Define variable for user's input
	var userInput shared.ListingData

	//Ensure input is JSON
	if err := c.ShouldBindJSON(&userInput); err != nil {
		glog.Errorf("Unable to bind listing JSON from page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	//Call the DB function
	glog.Infof("Calling createListingInDb function with userInput: %v", userInput)
	dbStatus, details, err := CreateListingInDb(shared.Database, userInput.Name, userInput.DisplayName, userInput.PrimaryType, userInput.SecondaryType, userInput.TertiaryType, userInput.Email, userInput.Website, userInput.Phone, userInput.PlusCode, userInput.StartTime, userInput.EndTime)

	if err != nil {
		//Return the status code and body from the function
		c.JSON(dbStatus, gin.H{"status": details, "error": err})
		return
	}

	//Return the status code and body from the function
	c.JSON(dbStatus, gin.H{"status": details})
	return
}

func GetListing(c *gin.Context) {
	//Define variable for user's input
	var userInput shared.ListingData

	//Ensure input is JSON
	if err := c.ShouldBindJSON(&userInput); err != nil {
		glog.Errorf("Unable to bind listing JSON from page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	//Call the DB function
	glog.Infof("Calling getListingFromDb function with userInput: %v", userInput)
	dbStatus, details, name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, pluscode, starttime, endtime, err := GetListingFromDb(shared.Database, userInput.Id)

	if err != nil {
		//Return the status code and body from the function
		c.JSON(dbStatus, gin.H{"status": details, "error": err})
		return
	}

	//Return successful response
	c.JSON(dbStatus, gin.H{"id": userInput.Id, "name": name, "displayName": displayname, "primarType": primarytype, "secondaryType": secondarytype, "tertiaryType": tertiarytype, "email": email, "website": website, "phone": phone, "plusCode": pluscode, "startTime": starttime, "endTime": endtime})
	return
}

func UpdateListing(c *gin.Context) {
	//Define variable for user's input
	var userInput shared.ListingData

	//Ensure input is JSON
	if err := c.ShouldBindJSON(&userInput); err != nil {
		glog.Errorf("Unable to bind listing JSON from page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	//Call the DB function
	glog.Infof("Calling getListingFromDb function with userInput: %v", userInput)
	dbStatus, details, err := UpdateListingInDb(shared.Database, userInput.Id, userInput.Name, userInput.DisplayName, userInput.PrimaryType, userInput.SecondaryType, userInput.TertiaryType, userInput.Email, userInput.Website, userInput.Phone, userInput.PlusCode, userInput.StartTime, userInput.EndTime)

	if err != nil {
		//Return the status code and body from the function
		c.JSON(dbStatus, gin.H{"status": details, "error": err})
		return
	}

	//Return the status code and body from the function
	c.JSON(dbStatus, gin.H{"status": details})
	return
}

func DeleteListing(c *gin.Context) {
	//Define variable for user's input
	var userInput shared.ListingData

	//Ensure input is JSON
	if err := c.ShouldBindJSON(&userInput); err != nil {
		glog.Errorf("Unable to bind listing JSON from page: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	//Call the DB function
	glog.Infof("Calling getListingFromDb function with userInput: %v", userInput)
	dbStatus, details, err := DeleteListingFromDb(shared.Database, userInput.Id)

	if err != nil {
		//Return the status code and body from the function
		c.JSON(dbStatus, gin.H{"status": details, "error": err})
		return
	}

	//Return the status code and body from the function
	c.JSON(dbStatus, gin.H{"status": details})
	return
}
