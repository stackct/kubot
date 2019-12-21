package command

type MockCommand struct{}

func (mc MockCommand) Name() string {
	return "Mock"
}

func (mc MockCommand) Execute(output chan string) {
	defer close(output)
	output <- "fin"
}

type MockParser struct {
	MockError   error
	MockCommand Command
}

type MockParserOption func(*MockParser)

func NewMockParser(options ...MockParserOption) *MockParser {
	mp := &MockParser{MockCommand: MockCommand{}, MockError: nil}

	for _, option := range options {
		option(mp)
	}

	return mp
}

func (mp MockParser) Parse(c string) (Command, error) {
	return mp.MockCommand, mp.MockError
}
