package reports

import (
	"townwatch/base"
	"townwatch/models"
	"townwatch/services/areas"
	"townwatch/services/events"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type ReportDetailsResponse struct {
	ReportDetails *models.GetReportDetailsRow
	Events        *[]models.Event
}

func ReadReport(b *base.Base, ctx *gin.Context, params *models.GetReportDetailsParams, censorEvents bool) (*ReportDetailsResponse, *base.CError) {
	reportDetails, err := b.DB.Queries.GetReportDetails(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	eventsO, err := b.DB.Queries.GetEventsByReport(ctx, params.ID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	reportDetails.Area = areas.CensorArea(reportDetails.Area)
	cenEvents := eventsO
	if censorEvents {
		cenEvents = events.CensorEvents(eventsO)
	}
	return &ReportDetailsResponse{
		ReportDetails: &reportDetails,
		Events:        &cenEvents,
	}, nil
}
