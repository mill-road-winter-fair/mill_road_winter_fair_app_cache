package listings

import (
	"net/http"

	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/shared"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)

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
