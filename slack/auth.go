package slack

import "kubot/config"

func HasAccess(email string, environment string) bool {
	return config.Conf.HasAccess(email, environment)
}
