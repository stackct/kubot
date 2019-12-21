package slack

func HasAccess(email string, environment string) bool {
	return Conf.HasAccess(email, environment)
}
