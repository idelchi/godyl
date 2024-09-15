package common

import (
	"github.com/spf13/viper"
)

func FromEnvOrFile(keys ...string) (string, error) {
	for _, key := range keys {
		if value := viper.GetString(key); value != "" {
			return value, nil
		}
	}

	return "", nil
}
