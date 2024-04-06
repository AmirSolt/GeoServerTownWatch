package reports

import (
	"bytes"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
)

type NotifEmailParams struct {
	BaseURL      string
	ReportParams []ReportParam
}

type ReportParam struct {
	Index   int
	ID      string
	BaseURL string
}

func getNotifEmailStr(b *base.Base, reports []models.Report) (string, *utils.CError) {

	reportParams := []ReportParam{}

	for i, report := range reports {
		reportParams = append(reportParams, ReportParam{
			Index:   i + 1,
			ID:      report.ID,
			BaseURL: b.FRONTEND_URL,
		})
	}

	notifParams := NotifEmailParams{
		BaseURL:      b.FRONTEND_URL,
		ReportParams: reportParams,
	}

	buf := new(bytes.Buffer)
	err = b.Emails.NotifEmail.Execute(buf, notifParams)
	if err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
		return "", cerr
	}

	return buf.String(), nil
}
