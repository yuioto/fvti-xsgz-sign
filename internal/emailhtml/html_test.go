package emailhtml

import (
"strings"
"testing"

"github.com/yuioto/fvti-xsgz-sign/pkg/client"
)

func TestFormatSignEmailHTML_zhCN(t *testing.T) {
tasks := []client.TaskSummary{{Name: "任务1", Status: "已完成", QD: "10", SignID: "abc", QDTimeText: "08:00"}}
got := FormatSignEmailHTML("zh-CN", "成功", "签到任务：任务1 (123)", "2026-01-01 08:00:00", tasks)
if !strings.Contains(got, "运行状态：成功") {
t.Fatalf("expected status in zh-CN HTML, got: %s", got)
}
if !strings.Contains(got, "任务名称") {
t.Fatalf("expected table header in zh-CN HTML, got: %s", got)
}
}

func TestFormatSignEmailHTML_en(t *testing.T) {
tasks := []client.TaskSummary{{Name: "task1", Status: "done", QD: "10", SignID: "abc", QDTimeText: "08:00"}}
got := FormatSignEmailHTML("en", "Success", "Sign task: task1 (123)", "2026-01-01 08:00:00", tasks)
if !strings.Contains(got, "Status") {
t.Fatalf("expected status label in en HTML, got: %s", got)
}
if !strings.Contains(got, "Task Name") {
t.Fatalf("expected table header in en HTML, got: %s", got)
}
}

func TestFormatSignEmailHTML_ja(t *testing.T) {
tasks := []client.TaskSummary{{Name: "task1", Status: "完了", QD: "10", SignID: "abc", QDTimeText: "08:00"}}
got := FormatSignEmailHTML("ja", "成功", "サインタスク：task1 (123)", "2026-01-01 08:00:00", tasks)
if !strings.Contains(got, "ステータス") {
t.Fatalf("expected status label in ja HTML, got: %s", got)
}
if !strings.Contains(got, "タスク名") {
t.Fatalf("expected table header in ja HTML, got: %s", got)
}
}
