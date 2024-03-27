package reports

import (
	"townwatch/base"
	"townwatch/models"
	"townwatch/services/areas"
	"townwatch/services/events"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReportDetailsResponse struct {
	ReportDetails *models.GetReportDetailsRow `json:"report_details"`
	Events        *[]models.Event             `json:"events"`
}

func GetReportDetails(b *base.Base, ctx *gin.Context, reportID string) (*models.GetReportDetailsRow, *utils.CError) {
	var byteArray [16]byte
	copy(byteArray[:], reportID)
	reportDetails, err := b.DB.Queries.GetReportDetails(ctx, pgtype.UUID{Bytes: byteArray, Valid: true})
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
	var byteArray [16]byte
	copy(byteArray[:], reportID)

	eventsO, err := b.DB.Queries.GetEventsByReport(ctx, pgtype.UUID{Bytes: byteArray, Valid: true})
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
