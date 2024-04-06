package base

import (
	"fmt"
	"text/template"

	"github.com/getsentry/sentry-go"
)

type Emails struct {
	NotifEmail *template.Template
}

func (base *Base) LoadEmails() {
	base.Emails.NotifEmail = loadNotifEmail()
}

func loadNotifEmail() *template.Template {
	fileName := "report_email.tmpl"
	filePath := fmt.Sprintf("./services/reports/templates/%s", fileName)
	tmpl, err := template.New(fileName).ParseFiles(filePath)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}

	return tmpl
}
