package errors

import "github.com/pkg/errors"

func New(message string) error {
	return errors.New(message)
}

func E(args ...interface{}) error {
	e := &Error{}

	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case ErrKind:
			e.Kind = arg
		case ErrPayload:
			e.Payload = arg
		case ErrClientMessage:
			e.ClientMessage = arg
		case error:
			e.Err = arg
		default:
			panic("Bad call to E")
		}
	}

	return e
}

func Kind(err error) ErrKind {
	e, ok := err.(*Error)
	if !ok {
		return KindServerError
	}

	if e.Kind != 0 {
		return e.Kind
	}

	// Keep unwrapping errors until we find one that has a kind.
	return Kind(e.Err)
}

func ClientMessage(err error) ErrClientMessage {
	e, ok := err.(*Error)
	if !ok {
		return "Something went wrong"
	}

	if e.ClientMessage != "" {
		return e.ClientMessage
	}

	// Keep unwrapping errors until we find one that has a client message.
	return ClientMessage(e.Err)
}

func Payload(err error) ErrPayload {
	e, ok := err.(*Error)
	if !ok {
		return ErrPayload{}
	}

	if e.Payload != (ErrPayload{}) {
		return e.Payload
	}

	// Keep unwrapping errors until we find one that has a kind.
	return Payload(e.Err)
}

func Ops(e *Error) []Op {
	res := []Op{e.Op}

	subErr, ok := e.Err.(*Error)
	if !ok {
		return res
	}

	res = append(res, Ops(subErr)...)
	return res
}
