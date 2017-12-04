package game

type Runnable interface {
	Run() error
	CleanUp() error
	SendMessage(message interface{}) ([]byte, error)
	SendMessageNoResponse(message interface{}) error
	Stderr() string
}
