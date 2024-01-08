package protocol

import "errors"

var ErrUnknownCommand = errors.New("unknown command type")

type Message interface {
	message()
	Raw() []byte
}

type Command interface {
	command()
}

type Protocol interface {
	Read(ch chan<- Message) error
	Send(cmd Command) error
}
