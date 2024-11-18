package listings

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/golang/glog"
)

// Define function to get all listings from db
func GetListingsFromDb(db *sql.DB) (dbStatus int, details string, rows *sql.Rows, err error) {

	query := "SELECT * FROM listings"
	rows, err = db.Query(query)

	// ...or if some other error occurred
	if err != nil {
		glog.Errorf("Error retrieving listing from database: %v", err)
		return http.StatusInternalServerError, "Database error", rows, err
	}

	// Successfully retrieved the listing
	glog.Infof("Listings successfully retrieved from database: %v", rows)
	return http.StatusOK, fmt.Sprintf("Listings successfully retrieved from database with ID: %v", rows), rows, nil
}

// Define function to add new listing to db
func CreateListingInDb(db *sql.DB, name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime string) (dbStatus int, details string, err error) {
	var listingId int

	query := "INSERT INTO listings (name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id"
	err = db.QueryRow(query, name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime).Scan(&listingId)

	// Check if no rows were returned...
	if err == sql.ErrNoRows {
		glog.Errorf("Error inserting listing to database: %v", err)
		return http.StatusInternalServerError, "Error inserting listing to database", err
	}

	// ...or if some other error occurred
	if err != nil {
		glog.Errorf("Database error: %v", err)
		return http.StatusInternalServerError, "Database error", err
	}

	glog.Infof("Listing added successfully to database with ID: %v", listingId)
	return http.StatusOK, fmt.Sprintf("Listing added successfully to database with ID: %v", listingId), nil
}

// Define function to get listing from db
func GetListingFromDb(db *sql.DB, listingId int) (dbStatus int, details, name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime string, err error) {

	query := "SELECT name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime FROM listings WHERE id = $1"
	err = db.QueryRow(query, listingId).Scan(&name, &displayname, &primarytype, &secondarytype, &tertiarytype, &email, &website, &phone, &latlng, &starttime, &endtime)

	if listingId < 0 {
		// Invalid ID (less than 0), this should be impossible tbh
		return http.StatusBadRequest, "Invalid ID", "", "", "", "", "", "", "", "", "", "", "", err
	}

	// Check if no rows were returned...
	if err == sql.ErrNoRows {
		glog.Errorf("No listing found with ID: %d", listingId)
		return http.StatusBadRequest, fmt.Sprintf("No listing found with ID: %d", listingId), "", "", "", "", "", "", "", "", "", "", "", err
	}

	// ...or if some other error occurred
	if err != nil {
		glog.Errorf("Error retrieving listing from database: %v", err)
		return http.StatusInternalServerError, "Database error", "", "", "", "", "", "", "", "", "", "", "", err
	}

	// Successfully retrieved the listing
	glog.Infof("Listing successfully retrieved from database with ID: %v", listingId)
	return http.StatusOK, fmt.Sprintf("Listing successfully retrieved from database with ID: %v", listingId), name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime, nil
}

func UpdateListingInDb(db *sql.DB, listingId int, name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime string) (dbStatus int, details string, err error) {

	query := "SELECT id FROM listings WHERE id = $1"
	err = db.QueryRow(query, listingId).Scan(&listingId)

	if listingId < 0 {
		// Invalid ID (less than 0), this should be impossible tbh
		return http.StatusBadRequest, "Invalid ID", err
	}

	// Check if no rows were returned...
	if err == sql.ErrNoRows {
		glog.Errorf("No listing found with ID: %d", listingId)
		return http.StatusBadRequest, fmt.Sprintf("No listing found with ID: %v", listingId), err
	}

	// ...or if some other error occurred
	if err != nil {
		glog.Errorf("Error retrieving listing from database: %v", err)
		return http.StatusInternalServerError, "Error retrieving listing from database", err
	}

	query = "UPDATE listings SET name = $1, displayname = $2, primarytype = $3, secondarytype = $4, tertiarytype = $5, email = $6, website = $7, phone = $8, latlng = $9, starttime = $10, endtime = $11 WHERE id = $12"
	_, err = db.Exec(query, name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime, listingId)

	if err != nil {
		glog.Errorf("Error updating listing in database: %v", err)
		return http.StatusInternalServerError, "Error updating listing in database", err
	}

	//Confirm changes have been successful
	//TODO: Actually check the new (updated) listing values from the db against the listing values supplied by the user
	query = "SELECT name, displayname, primarytype, secondarytype, tertiarytype, email, website, phone, latlng, starttime, endtime FROM listings WHERE id = $1"
	err = db.QueryRow(query, listingId).Scan(&name, &displayname, &primarytype, &secondarytype, &tertiarytype, &email, &website, &phone, &latlng, &starttime, &endtime)

	// Check if no rows were returned...
	if err == sql.ErrNoRows {
		glog.Errorf("Error updating listing in database: %v", err)
		return http.StatusInternalServerError, "Error updating listing in database", err
	}

	// ...or if some other error occurred
	if err != nil {
		glog.Errorf("Database error after listing update: %v", err)
		return http.StatusInternalServerError, "Database error after listing update", err
	}

	glog.Infof("Listing successfully updated in database with ID: %v", listingId)
	return http.StatusOK, fmt.Sprintf("Listing successfully updated in database with ID: %v", listingId), nil
}

func DeleteListingFromDb(db *sql.DB, listingId int) (dbStatus int, details string, err error) {

	query := "SELECT id FROM listings WHERE id = $1"
	err = db.QueryRow(query, listingId).Scan(&listingId)

	if listingId < 0 {
		// Invalid ID (less than 0), this should be impossible tbh
		return http.StatusBadRequest, "Invalid ID", err
	}

	// Check if no rows were returned...
	if err == sql.ErrNoRows {
		glog.Errorf("No listing found with ID: %d", listingId)
		return http.StatusBadRequest, fmt.Sprintf("No listing found with ID: %v", listingId), err
	}

	// ...or if some other error occurred
	if err != nil {
		glog.Errorf("Error retrieving listing from database: %v", err)
		return http.StatusInternalServerError, "Error retrieving listing from database", err
	}

	query = "DELETE FROM listings WHERE id = $1"
	_, err = db.Exec(query, listingId)

	if err != nil {
		glog.Errorf("Deleting listing from database failed: %v", err)
		return http.StatusInternalServerError, fmt.Sprintf("Deleting listing from database failed: %v", listingId), err
	}

	glog.Infof("Listing successfully deleted from database with ID: %v", listingId)
	return http.StatusOK, fmt.Sprintf("Listing successfully deleted from database with ID: %v", listingId), nil
}
