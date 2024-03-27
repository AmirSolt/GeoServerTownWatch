package base

import (
	"fmt"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
)

func ConvertArrayInterface[T any](intfSlice []interface{}) ([]T, *utils.CError) {
	var reportSlice []T
	for _, element := range intfSlice {
		report, ok := element.(T)
		if !ok {
			err := fmt.Errorf("type conversion failed. element: %v", element)
			eventID := sentry.CaptureException(err)

			return nil, &utils.CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
		}
		reportSlice = append(reportSlice, report)
	}
	return reportSlice, nil
}
