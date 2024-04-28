package testtask

import (
	"io"
)

func NewMessageEntryFactory(parser Parser[io.Reader, *MessageEntry]) MessageEntryFactory {
	return MessageEntryFactory{
		parser: parser,
	}
}

type MessageEntryFactory struct {
	parser Parser[io.Reader, *MessageEntry]
}

func (m MessageEntryFactory) Create(r io.Reader) ErrorProneEntry[string] {
	e := &MessageEntry{}
	e, err := m.parser.Parse(r, e)
	if err != nil {
		return &MessageEntry{err: err}
	}
	return e
}

func NewErrorEntryFactory() ErrorEntryFactory {
	return ErrorEntryFactory{}
}

type ErrorEntryFactory struct {
}

func (s ErrorEntryFactory) Create(err error) Entry[error] {
	return &ErrorEntry{err: err}
}

func NewAsyncErrorEntryFactory() AsyncErrorEntryFactory {
	return AsyncErrorEntryFactory{}
}

type AsyncErrorEntryFactory struct {
}

func (s AsyncErrorEntryFactory) Create() AsyncEntry[error] {
	errChan := make(chan error)
	return &AsyncErrorEntry{err: errChan}
}

func NewAsyncStringEntryFactory() AsyncStringEntryFactory {
	return AsyncStringEntryFactory{}
}

type AsyncStringEntryFactory struct {
}

func (s AsyncStringEntryFactory) Create() AsyncErrorProneEntry[string] {
	errChan := make(chan error)
	messageChan := make(chan string)
	return &AsyncStringEntry{err: errChan, message: messageChan}
}
