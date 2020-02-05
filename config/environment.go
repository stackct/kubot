package config

type Environment struct {
	Name    string   `yaml:"name"`
	Release string   `yaml:"release"`
	Users   []string `yaml:"users"`
	Channel string   `yaml:"channel"`
}

func (e Environment) ContainsUser(u string) bool {
	for _, user := range e.Users {
		if user == u {
			return true
		}
	}

	return false
}
