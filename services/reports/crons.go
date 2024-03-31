package reports

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"townwatch/base"
	"townwatch/models"

	"github.com/carlmjohnson/requests"
	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/robfig/cron"
)

type NotifCreateManyParams struct {
	params []NotifCreateParams
}

type NotifCreateParams struct {
	UserID   string `db:"user_id" json:"user_id"`
	Subject  string `db:"subject" json:"subject"`
	BodyHTML string `db:"body_html" json:"body_html"`
}

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	err := c.AddFunc("0 30 * * * *", func() {

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
			return
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

		errreq := requests.
			URL(fmt.Sprintf("%s/api/notifs/create", b.USER_SERVER_URL)).
			Method(http.MethodPost).
			Header(base.HeaderSecretKeyName, b.SECRET_API_KEY).
			BodyJSON(&NotifCreateManyParams{params: params}).
			CheckStatus(http.StatusOK, http.StatusAccepted).
			Fetch(context.Background())
		if errreq != nil {
			sentry.CaptureException(errreq)
			return
		}
	})

	if err != nil {
		sentry.CaptureException(err)
	}
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
