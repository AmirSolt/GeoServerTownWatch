package reports

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/services/areas"
	"townwatch/services/events"
	"townwatch/utils"

	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

// ==========================================================

type NotifCreateManyParams struct {
	params []NotifCreateParams
}

type NotifCreateParams struct {
	UserID   string `db:"user_id" json:"user_id"`
	Subject  string `db:"subject" json:"subject"`
	BodyHTML string `db:"body_html" json:"body_html"`
}

func CreateGlobalReports(b *base.Base) (*[]models.Report, *utils.CError) {

	sentry.CaptureMessage(fmt.Sprintf("Reports cron started at: %s", time.Now().Format(time.RFC1123)))

	reports, err := b.DB.Queries.CreateGlobalReports(context.Background(), models.CreateGlobalReportsParams{
		FromDate: pgtype.Timestamptz{
			Time:  time.Now().Add(-time.Duration(24) * time.Hour).UTC(),
			Valid: true,
		},
		ToDate: pgtype.Timestamptz{
			Time:  time.Now().UTC(),
			Valid: true,
		},
		EventsLimit: int32(b.ScanEventCountLimit),
	})
	if err != nil {
		return nil, err
	}

	aggUserReports := aggregateReportsByUser(reports)

	var params []NotifCreateParams
	for _, aggReports := range aggUserReports {
		params = append(params, NotifCreateParams{
			UserID:   aggReports[0].UserID,
			Subject:  "Test Subject",
			BodyHTML: fmt.Sprintf("Test Reports %+v", aggReports),
		})
	}

	errreq := createNotifsOnUserServer(b, params)
	if errreq != nil {
		eventID := sentry.CaptureException(errreq)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   errreq,
		}
	}

	return &reports, nil
}

func createNotifsOnUserServer(b *base.Base, params []NotifCreateParams) error {
	return requests.
		URL(fmt.Sprintf("%s/api/notifs/create/many", b.USER_SERVER_URL)).
		Method(http.MethodPost).
		Header(base.HeaderSecretKeyName, b.SECRET_API_KEY).
		BodyJSON(&NotifCreateManyParams{params: params}).
		CheckStatus(http.StatusOK, http.StatusAccepted).
		Fetch(context.Background())
}

func aggregateReportsByUser(reports []models.Report) [][]models.Report {
	reportsMap := make(map[string][]models.Report)
	for _, report := range reports {
		reportsMap[report.UserID] = append(reportsMap[report.UserID], report)
	}
	aggregatedReports := make([][]models.Report, 0, len(reportsMap))
	for _, userReports := range reportsMap {
		aggregatedReports = append(aggregatedReports, userReports)
	}

	return aggregatedReports
}
