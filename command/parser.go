package command

type CommandParser interface {
	Parse(input string) (Command, error)
}
