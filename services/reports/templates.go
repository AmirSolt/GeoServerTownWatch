package reports

import (
	"bytes"
	"context"
	"townwatch/base"
	"townwatch/models"
	templs "townwatch/templates"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
)

func getNotifEmailStr(b *base.Base, ctx context.Context, reports []models.Report) (string, *utils.CError) {

	component := templs.NotifEmail(b.FRONTEND_URL, reports)

	buf := new(bytes.Buffer)
	err := component.Render(ctx, buf)
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
