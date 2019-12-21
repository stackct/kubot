package command

type CommandParser interface {
	Parse(string) (Command, error)
}
