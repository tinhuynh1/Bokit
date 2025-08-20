package i18n

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

var (
	messages = make(map[string]map[string]string)
	once     sync.Once
)

func loadMessages() {
	basePath := "internal/i18n"

	locales := []string{"en", "vi"}

	for _, locale := range locales {
		filePath := filepath.Join(basePath, locale+".json")
		data, err := os.ReadFile(filePath)
		if err != nil {
			panic("failed to load i18n file: " + err.Error())
		}

		var msg map[string]string
		if err := json.Unmarshal(data, &msg); err != nil {
			panic("invalid i18n json: " + err.Error())
		}
		messages[locale] = msg
	}
}

func T(lang, key string) string {
	once.Do(loadMessages)

	if m, ok := messages[lang]; ok {
		if v, exists := m[key]; exists {
			return v
		}
	}
	return key
}
