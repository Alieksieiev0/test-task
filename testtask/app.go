package testtask

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/Alieksieiev0/test-task/reader"
	"github.com/Alieksieiev0/test-task/writer"
)

type App struct {
	sourcesMap map[string]string
	source     Source[string]
	reader     AppReader
	writer     AppWriter
}

func (a *App) Run() {
	var wg sync.WaitGroup
	inputs := a.source.Data()
	for {
		inputName := inputs.Next()
		if inputName.Err() != nil {
			break
		}
		wg.Add(1)

		readResults := a.reader.AsyncRead(inputName.Val())
		writeResults := a.writer.AsyncWrite(a.sourcesMap[inputName.Val()], readResults)
		go func() {
			for {
				select {
				case err, ok := <-readResults.Err():
					if !ok {
						continue
					}
					fmt.Println("ERROR READING FILE ", inputName.Val(), err)
					wg.Done()
					return
				case err, more := <-writeResults.Val():
					if !more {
						fmt.Println("NO MORE: ", inputName.Val())

						wg.Done()
						return
					}
					fmt.Println("ERROR SAVING FILE: ", inputName.Val(), err)
					wg.Done()
					return
					//cancelFunc()
				}

			}
		}()
	}
	wg.Wait()
}

func NewAppFactory() *AppFactory {
	return &AppFactory{}
}

type AppFactory struct {
}

func (a *AppFactory) Create(
	sourcesMap map[string]string,
	source Source[string],
	reader AppReader,
	writer AppWriter,
) *App {
	return &App{
		sourcesMap: sourcesMap,
		source:     source,
		reader:     reader,
		writer:     writer,
	}
}

type AppReader struct {
	processor         Processor[string, ErrorProneEntry[string]]
	operation         Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]]
	readerFactory     ErrorFactoryFunc[string, reader.Reader[string]]
	sourceFactory     FactoryFunc[reader.Reader[string], Source[string]]
	entryFactory      Factory[io.Reader, ErrorProneEntry[string]]
	asyncEntryFactory PlainFactory[AsyncErrorProneEntry[string]]
}

func (r *AppReader) AsyncRead(data string) AsyncErrorProneEntry[string] {
	asyncEntry := r.asyncEntryFactory.Create()
	go func() {
		reader, err := r.readerFactory(data)
		if err != nil {
			asyncEntry.PassErr(err)
			return
		}
		fmt.Println(data)
		defer reader.Close()
		r.operation.Run(
			r.processor.Process(
				r.sourceFactory(reader).Data(),
				func(data string) ErrorProneEntry[string] {
					return r.entryFactory.Create(strings.NewReader(data))
				},
			),
			asyncEntry,
		)
	}()
	return asyncEntry
}

func NewStringReaderFactory() *StringReaderFactory {
	return &StringReaderFactory{}
}

type StringReaderFactory struct {
}

func (a *StringReaderFactory) Create(
	processor Processor[string, ErrorProneEntry[string]],
	operation Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]],
	parser Parser[io.Reader, *MessageEntry],
	readerFactory ErrorFactoryFunc[string, reader.Reader[string]],
	sourceFactory FactoryFunc[reader.Reader[string], Source[string]],
) *AppReader {
	return &AppReader{
		processor:         processor,
		operation:         operation,
		readerFactory:     readerFactory,
		sourceFactory:     sourceFactory,
		entryFactory:      NewMessageEntryFactory(parser),
		asyncEntryFactory: NewAsyncStringEntryFactory(),
	}
}

type AppWriter struct {
	processor         Processor[string, error]
	operation         Operation[error, AsyncEntry[error]]
	sourceFactory     FactoryFunc[<-chan string, Source[string]]
	writerFactory     ErrorFactoryFunc[string, writer.Writer[string]]
	asyncEntryFactory PlainFactory[AsyncEntry[error]]
}

func (w *AppWriter) AsyncWrite(target string, data AsyncEntry[string]) AsyncEntry[error] {
	asyncEntry := w.asyncEntryFactory.Create()
	go func() {
		writer, err := w.writerFactory(target)
		if err != nil {
			asyncEntry.PassVal(err)
			return
		}
		defer writer.Close()

		w.operation.Run(
			w.processor.Process(w.sourceFactory(data.Val()).Data(), func(data string) error {
				return writer.Write(data)
			}),
			asyncEntry,
		)
	}()
	return asyncEntry
}

func NewStringWriterFactory() *StringWriterFactory {
	return &StringWriterFactory{}
}

type StringWriterFactory struct {
}

func (a *StringWriterFactory) Create(
	processor Processor[string, error],
	sourceFactory FactoryFunc[<-chan string, Source[string]],
	operation Operation[error, AsyncEntry[error]],
	writerFactory ErrorFactoryFunc[string, writer.Writer[string]],
) *AppWriter {
	return &AppWriter{
		processor:         processor,
		sourceFactory:     sourceFactory,
		operation:         operation,
		writerFactory:     writerFactory,
		asyncEntryFactory: NewAsyncErrorEntryFactory(),
	}
}
