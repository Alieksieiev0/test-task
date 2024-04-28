package testtask

import (
	"io"
	"time"
)

type MessageEntry struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	err       error
}

func (m *MessageEntry) Val() string {
	return "[" + m.Timestamp.Format(time.DateTime) + "]: " + m.Message
}

func (m *MessageEntry) Err() error {
	return m.err
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
	close(a.err)
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
	close(a.err)
	close(a.message)
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
