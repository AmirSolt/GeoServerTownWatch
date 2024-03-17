package base

import (
	"fmt"

	"github.com/getsentry/sentry-go"
)

func ConvertArrayInterface[T any](intfSlice []interface{}) ([]T, *CError) {
	var reportSlice []T
	for _, element := range intfSlice {
		report, ok := element.(T)
		if !ok {
			err := fmt.Errorf("type conversion failed. element: %v", element)
			eventID := sentry.CaptureException(err)
			return nil, &CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
		}
		reportSlice = append(reportSlice, report)
	}
	return reportSlice, nil
}
