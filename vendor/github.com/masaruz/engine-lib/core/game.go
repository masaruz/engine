package core

// Ack is callback
type Ack = func(form string)

// Game interface engine
type Game interface {
	Init() error
	Start() error
	Update(msg []byte, ack Ack) error
	GetState() ([]byte, error)
}
