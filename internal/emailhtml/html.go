package emailhtml

import (
	"fmt"
	"html"
	"strings"

	"github.com/yuioto/fvti-xsgz-sign/internal/i18n"
	"github.com/yuioto/fvti-xsgz-sign/pkg/client"
)

// FormatSignEmailHTML returns an HTML body for sign notification.
func FormatSignEmailHTML(locale, status, action, runAt string, tasks []client.TaskSummary) string {
	var b strings.Builder
	b.WriteString("<html><body>")
	b.WriteString(fmt.Sprintf("<h2>%s：%s</h2>", i18n.T(locale, "notify.email_status_label"), html.EscapeString(status)))
	b.WriteString(fmt.Sprintf("<p><strong>%s：</strong>%s</p>", i18n.T(locale, "notify.email_action_label"), html.EscapeString(action)))
	b.WriteString(fmt.Sprintf("<p><strong>%s：</strong>%s</p>", i18n.T(locale, "notify.email_run_at_label"), html.EscapeString(runAt)))

	b.WriteString(fmt.Sprintf("<h3>%s</h3>", i18n.T(locale, "notify.email_task_summary_title")))
	if tasks == nil {
		b.WriteString(fmt.Sprintf("<p>%s</p>", i18n.T(locale, "notify.email_fetch_failed")))
	} else if len(tasks) == 0 {
		b.WriteString(fmt.Sprintf("<p>%s</p>", i18n.T(locale, "notify.email_no_tasks")))
	} else {
		b.WriteString("<table border=\"1\" cellpadding=\"4\" cellspacing=\"0\">")
		b.WriteString(fmt.Sprintf("<thead><tr><th>%s</th><th>%s</th><th>%s</th><th>%s</th><th>%s</th></tr></thead><tbody>",
			i18n.T(locale, "notify.email_task_name"),
			i18n.T(locale, "notify.email_task_status"),
			i18n.T(locale, "notify.email_task_qd"),
			i18n.T(locale, "notify.email_task_signid"),
			i18n.T(locale, "notify.email_task_qdtime"),
		))

		for _, t := range tasks {
			b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>",
				html.EscapeString(t.Name),
				html.EscapeString(t.Status),
				html.EscapeString(t.QD),
				html.EscapeString(t.SignID),
				html.EscapeString(t.QDTimeText),
			))
		}
		b.WriteString("</tbody></table>")
	}

	b.WriteString("</body></html>")
	return b.String()
}
