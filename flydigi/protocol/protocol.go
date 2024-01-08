package protocol

import "errors"

var ErrUnknownCommand = errors.New("unknown command type")
var ErrUnknownMessage = errors.New("unknown message")

type Message interface {
	message()
}

type Command interface {
	command()
}

type Protocol interface {
	Close() error

	Messages() <-chan Message
	Send(cmd Command) error
}
