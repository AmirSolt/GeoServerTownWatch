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

type PostReportsBody struct {
	Reports []models.Report
}

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	c.AddFunc("0 30 * * * *", func() {

		reportsRaw, err := b.DB.Queries.CreateGlobalReports(context.Background(), models.CreateGlobalReportsParams{
			FromDate: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Duration(24) * time.Hour).UTC(),
				Valid: true,
			},
			ToDate: pgtype.Timestamptz{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			ScanEventsCountLimit: int32(b.ScanEventCountLimit),
		})
		if err != nil {
			sentry.CaptureException(err)
			return
		}

		// convert to Report
		reports, errconv := base.ConvertArrayInterface[models.Report](reportsRaw)
		if errconv != nil {
			return
		}

		body := PostReportsBody{
			Reports: reports,
		}
		errreq := requests.
			URL(fmt.Sprintf("%s/api/webhooks/reports", b.USER_SERVER_URL)).
			Method(http.MethodPost).
			Header(base.HeaderSecretKeyName, b.SECRET_API_KEY).
			BodyJSON(&body).
			CheckStatus(http.StatusOK, http.StatusAccepted).
			Fetch(context.Background())
		if errreq != nil {
			sentry.CaptureException(errreq)
			return
		}
	})
}
