package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/listings"
	"github.com/MarauderOne/mill_road_winter_fair_app_db_api/shared"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

func main() {

	//Parse the argument flags (like the ones in Heroku's Procfile)
	flag.Parse()

	//Define PostgreSQL connection string
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		//Hardcoded credentials only used in local environemnt
		dsn = "postgres://mrwfadmin:petersfield@localhost/mrwf_db?sslmode=disable"
	}
	glog.Info("Database connection string set:", dsn)

	//Open a connection to the database
	var dbErr error
	db, dbErr := sql.Open("postgres", dsn)
	if dbErr != nil {
		glog.Fatal(dbErr)
	}
	defer db.Close()

	//Test the connection
	dbErr = db.Ping()
	if dbErr != nil {
		glog.Fatal("Failed to connect to the database:", dbErr)
	}

	//Log successful db connection
	glog.Info("Successfully connected to PostgreSQL database:", db)

	//Inject database variable to other packages
	shared.Database = db

	//Create default webserver config
	webServer := gin.Default()

	//API endpoints to handle shop CRUD operations
	webServer.GET("/listings", listings.GetListings)
	
	webServer.PUT("/listing", listings.CreateListing)
	webServer.GET("/listing", listings.GetListing)
	webServer.POST("/listing", listings.UpdateListing)
	webServer.DELETE("/listing", listings.DeleteListing)

	//Define the port
	port := os.Getenv("PORT")
	//Define default port (for local testing)
	if port == "" {
		port = "8080"
	}

	//Run the webserver
	ginErr := webServer.Run(":" + port)
	if ginErr != nil {
		glog.Fatalf("Web server initialisation failed: %v", ginErr)
	}
}
