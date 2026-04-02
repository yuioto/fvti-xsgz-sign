package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

var (
	bundle   = map[string]map[string]string{}
	bundleMu sync.RWMutex
	once     sync.Once
	fallback = "zh-CN"
)

func localeFilePath(locale string) string {
	candidates := []string{
		filepath.Join("locales", fmt.Sprintf("%s.toml", locale)),
		filepath.Join("..", "..", "locales", fmt.Sprintf("%s.toml", locale)),
		filepath.Join("..", "..", "..", "locales", fmt.Sprintf("%s.toml", locale)),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return candidates[0]
}

func Init() error {
	var err error
	once.Do(func() {
		locales := []string{"zh-CN", "en", "ja"}
		for _, locale := range locales {
			path := localeFilePath(locale)
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				err = fmt.Errorf("load locale file %s: %w", path, readErr)
				return
			}

			var raw map[string]interface{}
			if unmarshalErr := toml.Unmarshal(data, &raw); unmarshalErr != nil {
				err = fmt.Errorf("parse locale file %s: %w", path, unmarshalErr)
				return
			}

			localized := map[string]string{}
			flatten("", raw, localized)
			bundleMu.Lock()
			bundle[locale] = localized
			bundleMu.Unlock()
		}
	})
	return err
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
		default:
			// ignore unsupported types
		}
	}
}

func T(locale, key string, args ...interface{}) string {
	if err := Init(); err != nil {
		return fmt.Sprintf("[i18n init error: %v]", err)
	}

	if locale == "" {
		locale = fallback
	}

	bundleMu.RLock()
	messages, ok := bundle[locale]
	bundleMu.RUnlock()
	if !ok {
		bundleMu.RLock()
		messages = bundle[fallback]
		bundleMu.RUnlock()
	}

	if val, ok := messages[key]; ok {
		if len(args) == 0 {
			return val
		}
		return fmt.Sprintf(val, args...)
	}

	bundleMu.RLock()
	fallbackMessages := bundle[fallback]
	bundleMu.RUnlock()
	if val, ok := fallbackMessages[key]; ok {
		if len(args) == 0 {
			return val
		}
		return fmt.Sprintf(val, args...)
	}

	return key
}
