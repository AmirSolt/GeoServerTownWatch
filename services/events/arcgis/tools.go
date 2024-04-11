package arcgis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
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

func fetchArcgis[T any](endpoint string) (*ArcgisResponse[T], *utils.CError) {

	resp, err := http.Get(endpoint)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	defer resp.Body.Close()

	var response ArcgisResponse[T]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	if len(response.Features) == 0 {
		sentry.CaptureMessage(fmt.Sprintf("Toronto Arcgis Response: Feature Len is 0 | URL: %s", endpoint))
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	vErr := validate.Struct(response)
	if vErr != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return &response, nil
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
	query.QOutSR("4326")
	query.QOutFields("*")
	query.QReturnGeometry("true")
	query.QFormat("pjson")
	return query
}
func (query *ArcgisQuery) QWhere(value string) *ArcgisQuery {
	query.QueryMap["where"] = value
	return query
}
func (query *ArcgisQuery) QOutSR(value string) *ArcgisQuery {
	query.QueryMap["outSR"] = value
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
