package arcgis

import (
	"fmt"
	"net/url"
	"time"
)

func convertToArcgisQueryTime(time time.Time) string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d\n",
		time.Year(),
		time.Month(),
		time.Day(),
		time.Hour(),
		time.Hour(),
		time.Second())
}

// ==============================================================
// Query Builder

type ArcgisQuery struct {
	QueryMap map[string]string
}

func NewArcgisQuery() *ArcgisQuery {
	return &ArcgisQuery{
		QueryMap: map[string]string{},
	}
}

func (query *ArcgisQuery) DefaultQueries() *ArcgisQuery {
	query.QInSR("4326")
	query.QOutFields("*")
	query.QReturnGeometry("true")
	query.QFormat("pjson")
	return query
}
func (query *ArcgisQuery) QWhere(value string) *ArcgisQuery {
	query.QueryMap["where"] = value
	return query
}
func (query *ArcgisQuery) QInSR(value string) *ArcgisQuery {
	query.QueryMap["inSR"] = value
	return query
}

func (query *ArcgisQuery) QOutFields(value string) *ArcgisQuery {
	query.QueryMap["outFields"] = value
	return query
}

func (query *ArcgisQuery) QReturnGeometry(value string) *ArcgisQuery {
	query.QueryMap["returnGeometry"] = value
	return query
}

func (query *ArcgisQuery) QFormat(value string) *ArcgisQuery {
	query.QueryMap["f"] = value
	return query
}

func (query *ArcgisQuery) Build() string {
	values := url.Values{}
	for key, value := range query.QueryMap {
		values.Add(key, value)
	}
	return values.Encode()
}

func (query *ArcgisQuery) BuildWithURL(url string) string {
	queryStr := query.Build()
	return fmt.Sprintf("%s?%s", url, queryStr)
}
