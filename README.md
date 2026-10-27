# Mill Road Winter Fair App Caching API
Originally this repository was the home for an entire API and database used by the Mill Road Winter Fair App. However it has now been stripped back to a simple caching system. If for some reason the old code is needed, check out the followiong commit: `865addb7e09fe0a647765de87eef22ab104320b1`

## Purpose
The purpose of this application is to act as a wrapper for the Google Sheets API in order to prevent excessive calls by the Mill Road Winter Fair App.

## Setting Up Your Local Environment
1. Install [Git for Windows](https://git-scm.com/downloads/win).

2. Clone this repository to your local environment using `git clone`.

3. Install [Go](https://go.dev/doc/).

### Environment Variables

At present these are the only environment variables needed to run the application. They allow the application to contact the following Google Sheets:
- Prod: https://docs.google.com/spreadsheets/d/1hkx3d4eVw2roFIEDdrYkpT0wwHKBdx7YaZP8vc-Cg2o/view
- Dev: https://docs.google.com/spreadsheets/d/1MO7sVJiuof_8wcy6yOnPE3Nhi65RRK_etgg1AxZM274/view

These should be stored in `.env` like so:
```
# API key used to call Google Sheets API
GOOGLE_SHEETS_API_KEY=*********************

# Google Sheet information for prod endpoint
GOOGLE_SHEET_ID=1hkx3d4eVw2roFIEDdrYkpT0wwHKBdx7YaZP8vc-Cg2o
GOOGLE_SHEET_RANGE=2025!A1:N350

# Google Sheet information for dev endpoint
DEV_GOOGLE_SHEET_ID=1MO7sVJiuof_8wcy6yOnPE3Nhi65RRK_etgg1AxZM274
DEV_GOOGLE_SHEET_RANGE=2025!A1:Q350

# Google Sheet information for /notifications
NOTIFICATIONS_GOOGLE_SHEET_ID=your-production-notifications-sheet-id
NOTIFICATIONS_GOOGLE_SHEET_RANGE=Sheet1!A1:Z350

# Google Sheet information for /notifications-dev
DEV_NOTIFICATIONS_GOOGLE_SHEET_ID=your-development-notifications-sheet-id
DEV_NOTIFICATIONS_GOOGLE_SHEET_RANGE=Sheet1!A1:Z350

# API key used to call our API endpoints (not the Google Sheets API)
OUR_API_KEY=*********************
```
Please request a copy of the Google Sheets API key from the repository's owner.

### Notification endpoints

`GET /notifications` and `GET /dev-notifications` return the raw Google Sheets JSON from separate production and development notification caches. Set the four notification sheet variables above in `.env`, replacing the example IDs and ranges with your own. Both use the existing `GOOGLE_SHEETS_API_KEY`.

All sheets are polled once a minute, with the first fetch one minute after startup. Failed fetches retain the last successful cached response. If a sheet's configuration is missing, its poller logs an error and does not start; restart after configuring it.

Both notification endpoints require `OUR_API_KEY`, supplied in the `X-API-Key` header or the `key` query parameter. A non-empty header takes precedence. Missing or incorrect keys return HTTP 401 without exposing cached data. An authorised request before the cache is populated returns HTTP 500. Existing listing endpoints retain their current permissive key handling.

### Prerequisites
- n/a 

## Other Links
- [Mill Road Winter Fair App code](https://github.com/MarauderOne/mill_road_winter_fair_app)
- [Test Data Spreadsheet](https://docs.google.com/spreadsheets/d/1-Dk_K8tvDJ4C9vSx0OJSEYhvhGrt6IEkabVRP83n0OM/edit?usp=sharing)
- [Prod Data Spreadsheet](https://docs.google.com/spreadsheets/d/1MO7sVJiuof_8wcy6yOnPE3Nhi65RRK_etgg1AxZM274/edit?usp=sharing)
- [Go Documentation](https://go.dev/doc/)
