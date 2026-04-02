package i18n

import (
	"embed"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

//go:embed locales/*.toml
var localeFiles embed.FS

const fallback = "zh-CN"

var (
	bundle           map[string]map[string]string
	bundleMu         sync.RWMutex
	initOnce         sync.Once
	initErr          error
	supportedLocales []string
)

// Init loads all locale files exactly once. Subsequent calls return the same result.
func Init() error {
	initOnce.Do(func() {
		entries, err := localeFiles.ReadDir("locales")
		if err != nil {
			initErr = fmt.Errorf("read locales: %w", err)
			return
		}

		m := make(map[string]map[string]string)
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".toml") {
				continue
			}
			locale := strings.TrimSuffix(entry.Name(), ".toml")
			data, err := localeFiles.ReadFile("locales/" + entry.Name())
			if err != nil {
				initErr = fmt.Errorf("load locale %s: %w", locale, err)
				return
			}

			var raw map[string]interface{}
			if err := toml.Unmarshal(data, &raw); err != nil {
				initErr = fmt.Errorf("parse locale %s: %w", locale, err)
				return
			}

			flat := make(map[string]string)
			flatten("", raw, flat)
			m[locale] = flat
			supportedLocales = append(supportedLocales, locale)
		}

		if len(supportedLocales) == 0 {
			initErr = fmt.Errorf("no locale files found")
			return
		}

		sort.Strings(supportedLocales)

		bundleMu.Lock()
		bundle = m
		bundleMu.Unlock()
	})
	return initErr
}

func flatten(prefix string, src map[string]interface{}, dst map[string]string) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case string:
			dst[key] = val
		case map[string]interface{}:
			flatten(key, val, dst)
		}
	}
}

// T returns the localized string for key in the given locale.
//
// It will normalize the locale string (e.g. en-US -> en), then attempt to
// look up the translation in the loaded locale bundle.
// If the locale key is not found, it falls back to zh-CN.
// If the key is still missing, it returns the key itself.
//
// T 会对 locale 做标准化（例如 en-US -> en），再在已加载的语言包中查找翻译。
// 如果在当前 locale 找不到则回退到 zh-CN，再找不到时返回 key 本身。
// 支持额外的 fmt.Sprintf 参数，允许模板格式化。
func T(locale, key string, args ...interface{}) string {
	if err := Init(); err != nil {
		return key
	}

	locale = NormalizeLocale(locale)

	bundleMu.RLock()
	defer bundleMu.RUnlock()

	try := func(loc string) (string, bool) {
		if messages, ok := bundle[loc]; ok {
			if val, ok := messages[key]; ok {
				if len(args) == 0 {
					return val, true
				}
				return fmt.Sprintf(val, args...), true
			}
		}
		return "", false
	}

	if val, ok := try(locale); ok {
		return val
	}

	if locale != fallback {
		if val, ok := try(fallback); ok {
			return val
		}
	}

	return key
}

// GetSupportedLocales returns a copy of locales that are currently loaded.
// It is safe for callers to mutate the returned slice.
//
// 返回当前已加载的 locale 列表副本，防止外部修改内部 state。
func GetSupportedLocales() []string {
	if err := Init(); err != nil {
		return nil
	}

	bundleMu.RLock()
	defer bundleMu.RUnlock()

	out := make([]string, len(supportedLocales))
	copy(out, supportedLocales)
	return out
}

// IsSupportedLocale checks whether the locale identifier is supported.
//
// 通过不区分大小写匹配来判断 locale 是否在当前支持列表中。
func IsSupportedLocale(locale string) bool {
	_, ok := findSupportedLocale(locale)
	return ok
}

// NormalizeLocale normalizes a locale string to a supported canonical locale.
//
// 例如：
//   "" -> "zh-CN"
//   "en-US" -> "en"
//   "zh-TW" -> "zh-CN"
//   "UNKNOWN" -> "zh-CN"（fallback）
//
// 针对常见的 locale 变化（下划线、大小写、国家/地区后缀）做自动映射。
func NormalizeLocale(locale string) string {
	if locale == "" {
		return fallback
	}
	locale = strings.TrimSpace(locale)
	if locale == "" {
		return fallback
	}
	locale = strings.ReplaceAll(locale, "_", "-")
	localeLower := strings.ToLower(locale)

	switch localeLower {
	case "zh", "zh-cn", "zh-hans", "zh-hans-cn":
		return "zh-CN"
	case "zh-tw", "zh-hant", "zh-hk":
		return "zh-CN"
	case "en", "en-us", "en-gb", "en-au", "en-ca":
		return "en"
	case "ja", "ja-jp":
		return "ja"
	}

	if canonical, ok := findSupportedLocale(localeLower); ok {
		return canonical
	}

	if strings.Contains(localeLower, "-") {
		parts := strings.SplitN(localeLower, "-", 2)
		if canonical, ok := findSupportedLocale(parts[0]); ok {
			return canonical
		}
	}

	return fallback
}

// findSupportedLocale returns a supported locale key that matches the input, case-insensitive.
//
// 为 NormalizeLocale 和 IsSupportedLocale 提供内部匹配功能。
func findSupportedLocale(locale string) (string, bool) {
	if err := Init(); err != nil {
		return "", false
	}
	locale = strings.TrimSpace(locale)
	if locale == "" {
		return "", false
	}

	for _, sup := range supportedLocales {
		if strings.EqualFold(sup, locale) {
			return sup, true
		}
	}

	return "", false
}

// GetMissingKeys returns all translation keys present in fallback locale but missing in the requested locale.
//
// 用于测试和检查翻译覆盖率，方便补全各 locale 的翻译。
func GetMissingKeys(locale string) []string {
	if err := Init(); err != nil {
		return nil
	}

	locale = NormalizeLocale(locale)

	bundleMu.RLock()
	defer bundleMu.RUnlock()

	fallbackMessages, ok := bundle[fallback]
	if !ok {
		return nil
	}

	localeMessages, ok := bundle[locale]
	if !ok {
		localeMessages = map[string]string{}
	}

	var missing []string
	for key := range fallbackMessages {
		if _, ok := localeMessages[key]; !ok {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	return missing
}
