package command

import (
	"kubot/config"
)

type Context struct {
	Environment config.Environment
	User        string
}
