# Mill Road Winter Fair App Caching API
Originally this repository was the home for an entire API and database used by the Mill Road Winter Fair App. However it has now been stripped back to a simple caching system. If for some reason the old code is needed, check out the followiong commit: `865addb7e09fe0a647765de87eef22ab104320b1`

## Purpose
The purpose of this application is to act as a wrapper for the Google Sheets API in order to prevent excessive calls by the Mill Road Winter Fair App.

## Setting Up Your Local Environment
1. Install [Git for Windows](https://git-scm.com/downloads/win).

2. Clone this repository to your local environment using `git clone`.

3. Install [Go](https://go.dev/doc/).

### Environment Variables

At present these are the only environment variables needed to run the application. They allow the application to contact the following Google Sheet: https://docs.google.com/spreadsheets/d/1hkx3d4eVw2roFIEDdrYkpT0wwHKBdx7YaZP8vc-Cg2o/view

These should be stored in `.env` like so:
```
GOOGLE_SHEETS_API_KEY=*********************
GOOGLE_SHEET_ID=1hkx3d4eVw2roFIEDdrYkpT0wwHKBdx7YaZP8vc-Cg2o
GOOGLE_SHEET_RANGE=2025!A1:N350
```
Please request a copy of the API key from the repository's owner.

### Prerequisites
- n/a 

## Other Links
- [Mill Road Winter Fair App code](https://github.com/MarauderOne/mill_road_winter_fair_app)
- [Test Data Spreadsheet](https://docs.google.com/spreadsheets/d/1-Dk_K8tvDJ4C9vSx0OJSEYhvhGrt6IEkabVRP83n0OM/edit?usp=sharing)
- [Prod Data Spreadsheet](https://docs.google.com/spreadsheets/d/1hkx3d4eVw2roFIEDdrYkpT0wwHKBdx7YaZP8vc-Cg2o/edit?usp=sharing)
- [Go Documentation](https://go.dev/doc/)
