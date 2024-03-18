package reports

import (
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type GetPublicReportDetailsParams struct {
	ID pgtype.UUID `json="id"`
}
type PrivateReportDetailsResponse struct {
	ReportDetails *models.GetPrivateReportDetailsRow
	Events        *[]models.Event
}
type PublicReportDetailsResponse struct {
	ReportDetails *models.GetPublicReportDetailsRow
	Events        *[]models.Event
}

func ReadPrivateReport(b *base.Base, ctx *gin.Context, params *models.GetPrivateReportDetailsParams) (*PrivateReportDetailsResponse, *base.CError) {
	reportDetails, err := b.DB.Queries.GetPrivateReportDetails(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	events, err := b.DB.Queries.GetEventsByReport(ctx, params.ID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	return &PrivateReportDetailsResponse{
		ReportDetails: &reportDetails,
		Events:        &events,
	}, nil
}

func ReadPublicReport(b *base.Base, ctx *gin.Context, params *GetPublicReportDetailsParams) (*PublicReportDetailsResponse, *base.CError) {
	reportDetails, err := b.DB.Queries.GetPublicReportDetails(ctx, params.ID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	events, err := b.DB.Queries.GetEventsByReport(ctx, params.ID)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	return &PublicReportDetailsResponse{
		ReportDetails: &reportDetails,
		Events:        &events,
	}, nil
}
