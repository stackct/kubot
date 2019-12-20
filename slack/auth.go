package slack

import (
	"kubot/config"
)

func HasAccess(username string, environment string) bool {
	return config.AppConfig.HasAccess(username, environment)
}
