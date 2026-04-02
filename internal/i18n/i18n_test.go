package i18n

import "testing"

func TestT(t *testing.T) {
	if got := T("zh-CN", "notify.success_title"); got != "签到成功" {
		t.Fatalf("expected 签到成功, got %q", got)
	}

	if got := T("en", "notify.success_title"); got != "Sign-in Success" {
		t.Fatalf("expected Sign-in Success, got %q", got)
	}

	if got := T("ja", "notify.success_title"); got != "サインイン成功" {
		t.Fatalf("expected サインイン成功, got %q", got)
	}

	if got := T("zh-CN", "error.load_config", "errinfo"); got != "加载配置失败：errinfo" {
		t.Fatalf("expected 加载配置失败：errinfo, got %q", got)
	}
}
