package i18n

import (
	"testing"
)

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

func TestNormalizeLocale(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"", "zh-CN"},
		{"en-US", "en"},
		{"en_gb", "en"},
		{"ja-JP", "ja"},
		{"zh", "zh-CN"},
		{"zh-TW", "zh-CN"},
		{"unknown", "zh-CN"},
	}

	for _, c := range cases {
		if got := NormalizeLocale(c.input); got != c.expected {
			t.Fatalf("NormalizeLocale(%q) expected %q, got %q", c.input, c.expected, got)
		}
	}
}

func TestSupportedLocales(t *testing.T) {
	locales := GetSupportedLocales()
	if len(locales) == 0 {
		t.Fatal("GetSupportedLocales should return at least one locale")
	}

	if !IsSupportedLocale("zh-CN") {
		t.Fatal("zh-CN should be a supported locale")
	}

	if !IsSupportedLocale("en") {
		t.Fatal("en should be a supported locale")
	}

	if !IsSupportedLocale("ja") {
		t.Fatal("ja should be a supported locale")
	}

	if IsSupportedLocale("unknown") {
		t.Fatal("unknown should not be a supported locale")
	}
}

func TestMissingKeys(t *testing.T) {
	missing := GetMissingKeys("en")
	if len(missing) != 0 {
		t.Fatalf("expected no missing keys for en, got %d", len(missing))
	}
}

func TestMissingKeysForJa(t *testing.T) {
	missing := GetMissingKeys("ja")
	if len(missing) != 0 {
		t.Fatalf("expected no missing keys for ja, got %d", len(missing))
	}
}
