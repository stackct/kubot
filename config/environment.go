package config

type Environment struct {
	Name      string            `yaml:"name"`
	Users     []string          `yaml:"users"`
	Channel   string            `yaml:"channel"`
	Variables map[string]string `yaml:"variables"`
}

func (e Environment) ContainsUser(u string) bool {
	for _, user := range e.Users {
		if user == u {
			return true
		}
	}

	return false
}
