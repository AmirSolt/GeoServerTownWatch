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

func GetReportsByUser(b *base.Base, ctx *gin.Context, userID string) (*[]models.Report, *utils.CError) {
	reports, err := b.DB.Queries.GetReportsByUser(ctx, userID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return &reports, nil
}

func GetReportsByArea(b *base.Base, ctx *gin.Context, params *models.GetReportsByAreaParams) (*[]models.Report, *utils.CError) {
	reports, err := b.DB.Queries.GetReportsByArea(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return &reports, nil
}
