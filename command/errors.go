package command

type UnknownCommandError struct {
	err string
}

func NewUnknownCommandError(err string) *UnknownCommandError {
	return &UnknownCommandError{err}
}

func (e UnknownCommandError) Error() string {
	return e.err
}

func NewCommandArgumentError(err string) *CommandArgumentError {
	return &CommandArgumentError{err}
}

type CommandArgumentError struct {
	err string
}

func (e CommandArgumentError) Error() string {
	return e.err
}
