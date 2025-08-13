package env

import (
	"os"
)

func String(name, defaultValue string) string {
	if value, ok := os.LookupEnv(name); ok && value != "" {
		return value
	}
	return defaultValue
}
