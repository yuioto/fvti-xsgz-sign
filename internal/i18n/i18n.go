package i18n

import (
"embed"
"fmt"
"sync"

"github.com/pelletier/go-toml/v2"
)

//go:embed locales/*.toml
var localeFiles embed.FS

const fallback = "zh-CN"

var (
bundle   map[string]map[string]string
bundleMu sync.RWMutex
initOnce sync.Once
initErr  error
)

// Init loads all locale files exactly once. Subsequent calls return the same result.
func Init() error {
initOnce.Do(func() {
locales := []string{"zh-CN", "en", "ja"}
m := make(map[string]map[string]string, len(locales))
for _, locale := range locales {
data, err := localeFiles.ReadFile(fmt.Sprintf("locales/%s.toml", locale))
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
}
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
// Falls back to zh-CN if locale is empty or not found.
// Supports fmt.Sprintf-style args when additional arguments are provided.
func T(locale, key string, args ...interface{}) string {
if err := Init(); err != nil {
return key
}

if locale == "" {
locale = fallback
}

bundleMu.RLock()
defer bundleMu.RUnlock()

if messages, ok := bundle[locale]; ok {
if val, ok := messages[key]; ok {
if len(args) == 0 {
return val
}
return fmt.Sprintf(val, args...)
}
}

// fallback to zh-CN
if locale != fallback {
if messages, ok := bundle[fallback]; ok {
if val, ok := messages[key]; ok {
if len(args) == 0 {
return val
}
return fmt.Sprintf(val, args...)
}
}
}

return key
}
