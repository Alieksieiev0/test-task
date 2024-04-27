package testtask

import (
	"fmt"
	"io"
	"strings"

	"github.com/Alieksieiev0/test-task/iterator"
)

type App struct {
	fileOpener FileOpener
	fileReader FileReader
}

type FileReader struct {
	processor    Processor[string, Entry]
	factory      Factory[io.Reader, Entry]
	sourceLoader SourceLoader[io.Reader, string]
	operation    Operation[Entry]
}

func (a *App) Run() {
	entries := a.fileOpener.Open()
	for {
		entry := entries.Next()
		if entry.Err() != nil {
			return
		}
		file := entry.Val()
		if file.Err() != nil {
			return
		}
		file.Close()
		if err := a.fileReader.Read(file.Val()); err != nil {
			fmt.Println(err)
		}
	}
}

func (f *FileReader) Read(file io.Reader) error {
	f.sourceLoader.Load(file)
	return f.operation.Run(f.processor.Process(f.sourceLoader.Data(), func(data string) Entry {
		return f.factory.Create(strings.NewReader(data))
	}))
}

type FileOpener struct {
	processor Processor[string, OSEntry[io.Reader]]
	source    Source[string]
	factory   Factory[string, OSEntry[io.Reader]]
}

func (f *FileOpener) Open() iterator.Iterator[OSEntry[io.Reader]] {
	return f.processor.Process(f.source.Data(), func(data string) OSEntry[io.Reader] {
		return f.factory.Create(data)
	})
}
