package main

import (
	"fmt"
	"townwatch/services/events/arcgis"
)

func main() {
	// urlStr := "https://services.arcgis.com/S9th0jAJ7bqgIRjw/ArcGIS/rest/services/YTD_CRIME_WM/FeatureServer/0/query"
	urlStr := "https://services8.arcgis.com/lYI034SQcOoxRCR7/arcgis/rest/services/PublicCrimeDataFME/FeatureServer/0/query"
	where := "0=0"
	endpoint := arcgis.NewArcgisQuery().DefaultQueries().QWhere(where).BuildWithURL(urlStr)
	fmt.Println(">>>>", endpoint)
}
