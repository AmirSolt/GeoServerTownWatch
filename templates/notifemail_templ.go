// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.648
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "fmt"
import "strconv"
import "townwatch/models"

func NotifEmail(baseURL string, reports []models.Report) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html id=\"__svelte-email\" lang=\"en\"><head><meta httpequiv=\"Content-Type\" content=\"text/html; charset=UTF-8\"></head><div id=\"__svelte-email-preview\" style=\"display:none;overflow:hidden;line-height:1px;opacity:0;max-height:0;max-width:0;\">We have detected incidents near your location and compiled reports. \r<div>\u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff \u200c\u200b\u200d\u200e\u200f\ufeff\r</div></div><table style=\"width:100%;background-color:#ffffff;\" align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\"><tbody><tr style=\"display:grid;grid-auto-columns:minmax(0, 1fr);grid-auto-flow:column;\"><div></div><div style=\"max-width:37.5em;margin:0 auto;padding:20px 0 48px;width:580px;\"><td style=\"display:inline-flex;justify-content:center;align-items:center;\" role=\"presentation\"><img alt=\"Civil Watch\" src=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s/logo.png", baseURL))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates\notifemail.templ`, Line: 46, Col: 49}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" width=\"60\" height=\"60\" style=\"display:block;outline:none;border:none;text-decoration:none;\"><h1 style=\"font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Oxygen-Sans,Ubuntu,Cantarell,Helvetica Neue,sans-serif;font-size:32px;line-height:1.3;font-weight:700;color:#484848;\">Civil Watch\r</h1></td><table style=\"width:100%;\" align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\"><tbody><tr style=\"display:grid;grid-auto-columns:minmax(0, 1fr);grid-auto-flow:column;\"><img alt=\"Map\" src=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%s/map.png", baseURL))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates\notifemail.templ`, Line: 62, Col: 50}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" width=\"200\" height=\"200\" style=\"display:block;outline:none;border:none;text-decoration:none;margin:0 auto;\"></tr></tbody></table><p style=\"font-size:24px;line-height:1.3;margin:16px 0;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Oxygen-Sans,Ubuntu,Cantarell,Helvetica Neue,sans-serif;font-weight:600;color:#484848;text-align:center;\">We have detected incidents near your location and compiled reports.\r</p><p style=\"font-size:14px;line-height:24px;margin:16px 0;\"></p><hr style=\"width:100%;border:none;border-top:1px solid #eaeaea;border-color:#cccccc;margin:20px 0;\"><p style=\"font-size:28px;line-height:1.3;margin:16px 0;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Oxygen-Sans,Ubuntu,Cantarell,Helvetica Neue,sans-serif;font-weight:600;color:#484848;\">Reports:\r</p><div></div><table style=\"width:100%;\" align=\"center\" border=\"0\" cellpadding=\"0\" cellspacing=\"0\" role=\"presentation\"><tbody><tr style=\"display:grid;grid-auto-columns:minmax(0, 1fr);grid-auto-flow:column;\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for i, report := range reports {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a href=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 templ.SafeURL = templ.SafeURL(fmt.Sprintf("%s/reports/%s", baseURL, report.ID))
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var4)))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" target=\"_blank\" style=\"font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Oxygen-Sans,Ubuntu,Cantarell,Helvetica Neue,sans-serif;background-color:#0470dc;border-radius:3px;color:#fff;font-size:18px;text-decoration:none;text-align:center;display:inline-block;width:100%;p-x:0;p-y:19;line-height:100%;max-width:100%;padding:19px 0px;\"><span></span> <span style=\"font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Oxygen-Sans,Ubuntu,Cantarell,Helvetica Neue,sans-serif;background-color:#0470dc;border-radius:3px;color:#fff;font-size:18px;text-decoration:none;text-align:center;display:inline-block;width:100%;p-x:0;p-y:19;max-width:100%;line-height:120%;text-transform:none;mso-padding-alt:0px;mso-text-raise:14.25;\">Report #")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var5 string
			templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(strconv.Itoa(i + 1))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates\notifemail.templ`, Line: 95, Col: 41}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span> <span></span></a><p style=\"font-size:14px;line-height:24px;margin:16px 0;\"></p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</tr></tbody></table><div></div><hr style=\"width:100%;border:none;border-top:1px solid #eaeaea;border-color:#cccccc;margin:20px 0;\"><a href=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var6 templ.SafeURL = templ.SafeURL(baseURL)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var6)))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" target=\"_blank\" style=\"color:#9ca299;text-decoration:none;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Oxygen-Sans,Ubuntu,Cantarell,Helvetica Neue,sans-serif;font-size:14px;margin-bottom:10px;\">Civil\r Watch @2024\r</a><p style=\"font-size:14px;line-height:24px;margin:16px 0;font-family:-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Oxygen-Sans,Ubuntu,Cantarell,Helvetica Neue,sans-serif;color:#9ca299;margin-bottom:10px;\">Report unsafe behavior: support@civilwatch.net\r</p></div><div></div></tr></tbody></table></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
