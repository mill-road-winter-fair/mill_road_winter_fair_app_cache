# Mill Road Winter Fair App Caching API
Originally this repository was the home for an entire API and database used by the Mill Road Winter Fair App. However it has now been stripped back to a simple caching system. If for some reason the old code is needed, check out the followiong commit: `865addb7e09fe0a647765de87eef22ab104320b1`

## Purpose
The purpose of this application is to act as a wrapper for the Google Sheets API in order to prevent excessive calls by the Mill Road Winter Fair App.

### Potential Future Development Ideas
- Add validation which removes incomplete records.

## Setting Up Your Local Environment

### Environment Variables

At present these are the only environment variables needed to run the application. They allow the application contact the following Google Sheet: https://docs.google.com/spreadsheets/d/1hkx3d4eVw2roFIEDdrYkpT0wwHKBdx7YaZP8vc-Cg2o/view

These should be stored in `.env` like so:
```
GOOGLE_SHEETS_API_KEY=*********************
GOOGLE_SHEET_ID=1hkx3d4eVw2roFIEDdrYkpT0wwHKBdx7YaZP8vc-Cg2o
GOOGLE_SHEET_RANGE=A1:L200
```
Please request a copy of the API key from the repository's owner.


### Prerequisites
- n/a 

## Postman API Documentation

Example API use cases can be found on Postman, [here](https://orange-crater-235389.postman.co/workspace/Mill-Road-Winter-Fair-App-API~2f3cf4fe-aa38-46a2-bd71-d2b5fa5e76fb/overview).

## Other Links
- Mill Road Winter Fair App: https://github.com/MarauderOne/mill_road_winter_fair_app
