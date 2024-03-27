package reports

import (
	"townwatch/base"
	"townwatch/models"
	"townwatch/services/areas"
	"townwatch/services/events"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type ReportDetailsResponse struct {
	ReportDetails *models.GetReportDetailsRow `json:"report_details"`
	Events        *[]models.Event             `json:"events"`
}

func GetReportDetails(b *base.Base, ctx *gin.Context, reportID string) (*models.GetReportDetailsRow, *utils.CError) {
	// reportIDBytes, err := utils.ParseUUID(reportID)
	// if err != nil {
	// 	eventID := sentry.CaptureException(err)
	// 	return nil, &utils.CError{
	// 		EventID: eventID,
	// 		Message: "Internal Server Error",
	// 		Error:   err,
	// 	}
	// }

	reportDetails, err := b.DB.Queries.GetReportDetails(ctx, reportID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	reportDetails.Area = areas.CensorArea(reportDetails.Area)

	return &reportDetails, nil
}

func GetEventsByReport(b *base.Base, ctx *gin.Context, reportID string, censorEvents bool) (*[]models.Event, *utils.CError) {
	// reportIDBytes, err := utils.ParseUUID(reportID)
	// if err != nil {
	// 	eventID := sentry.CaptureException(err)
	// 	return nil, &utils.CError{
	// 		EventID: eventID,
	// 		Message: "Internal Server Error",
	// 		Error:   err,
	// 	}
	// }

	eventsO, err := b.DB.Queries.GetEventsByReport(ctx, reportID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	cenEvents := eventsO
	if censorEvents {
		cenEvents = events.CensorEvents(eventsO)
	}
	return &cenEvents, nil
}
