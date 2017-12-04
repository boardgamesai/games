package game

type RunnablePlayerMock struct{}

func (p *RunnablePlayerMock) Run() error {
	return nil
}

func (p *RunnablePlayerMock) CleanUp() error {
	return nil
}

func (p *RunnablePlayerMock) SendMessage(message interface{}) ([]byte, error) {
	return []byte{}, nil
}

func (p *RunnablePlayerMock) SendMessageNoResponse(message interface{}) error {
	return nil
}

func (p *RunnablePlayerMock) Stderr() string {
	return ""
}
