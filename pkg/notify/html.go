package notify

import (
	"fmt"
	"html"
	"strings"

	"github.com/yuioto/fvti-xsgz-sign/pkg/client"
)

// FormatSignEmailHTML returns an HTML body for sign notification.
func FormatSignEmailHTML(status, action, runAt string, tasks []client.TaskSummary) string {
	var b strings.Builder
	b.WriteString("<html><body>")
	b.WriteString(fmt.Sprintf("<h2>运行状态：%s</h2>", html.EscapeString(status)))
	b.WriteString(fmt.Sprintf("<p><strong>本次运行做了什么：</strong>%s</p>", html.EscapeString(action)))
	b.WriteString(fmt.Sprintf("<p><strong>本次运行的UTC+8时间：</strong>%s</p>", html.EscapeString(runAt)))

	b.WriteString("<h3>签到列表状态</h3>")
	if tasks == nil {
		b.WriteString("<p>获取执行结果失败。</p>")
	} else if len(tasks) == 0 {
		b.WriteString("<p>当前没有签到任务。</p>")
	} else {
		b.WriteString("<table border=\"1\" cellpadding=\"4\" cellspacing=\"0\">")
		b.WriteString("<thead><tr><th>任务名称</th><th>签到情况</th><th>QD</th><th>SignID</th><th>QDTimeText</th></tr></thead><tbody>")
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
