package lib

// Form to comunicate with plugin
type Form = struct {
	Header    string
	Body      string
	Signature string
}

// Ack is callback
type Ack = func(form Form)

// Game interface engine
type Game interface {
	On(name string, ack Ack)
	Init() error
	Start() error
	Update(msg []byte) error
	Pause() error
	Finish() error
	Shutdown() error
}
