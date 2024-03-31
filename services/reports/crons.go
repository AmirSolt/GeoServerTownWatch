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

func LoadCronJobs(b *base.Base, c *cron.Cron) {
	err := c.AddFunc("0 30 * * * *", func() {

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

		errreq := requests.
			URL(fmt.Sprintf("%s/api/notifs/create", b.USER_SERVER_URL)).
			Method(http.MethodPost).
			Header(base.HeaderSecretKeyName, b.SECRET_API_KEY).
			BodyJSON(&reports).
			CheckStatus(http.StatusOK, http.StatusAccepted).
			Fetch(context.Background())
		if errreq != nil {
			sentry.CaptureException(errreq)
			return
		}
	})

	sentry.CaptureException(err)
}
