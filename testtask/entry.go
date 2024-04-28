package testtask

import (
	"encoding/json"
	"io"
	"time"
)

type MessageEntry struct {
	message   string
	timestamp time.Time
	err       error
}

func (m *MessageEntry) Val() string {
	return "[" + m.timestamp.Format(time.DateTime) + "]: " + m.message
}

func (m *MessageEntry) Err() error {
	return m.err
}

func (m *MessageEntry) UnmarshalJSON(data []byte) error {
	var msg struct {
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
	}

	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	m.message = msg.Message
	m.timestamp = msg.Timestamp
	return nil
}

type ErrorEntry struct {
	err error
}

func (e *ErrorEntry) Val() error {
	return e.err
}

type AsyncErrorEntry struct {
	err chan error
}

func (a *AsyncErrorEntry) Val() <-chan error {
	return a.err
}

func (a *AsyncErrorEntry) PassVal(err error) {
	a.err <- err
}

func (a *AsyncErrorEntry) Close() {
	if a.err != nil {
		close(a.err)
	}
}

type AsyncStringEntry struct {
	message chan string
	err     chan error
}

func (a *AsyncStringEntry) Val() <-chan string {
	return a.message
}

func (a *AsyncStringEntry) Err() <-chan error {
	return a.err
}

func (a *AsyncStringEntry) PassVal(message string) {
	a.message <- message
}

func (a *AsyncStringEntry) PassErr(err error) {
	a.err <- err
}

func (a *AsyncStringEntry) Close() {
	if a.err != nil {
		close(a.err)
	}
	if a.message != nil {
		close(a.message)
	}
}

type FileEntry struct {
	readCloser io.ReadCloser
	err        error
}

func (f *FileEntry) Val() io.Reader {
	return f.readCloser
}

func (f *FileEntry) Err() error {
	return f.err
}

func (f *FileEntry) Close() {
	f.readCloser.Close()
}
