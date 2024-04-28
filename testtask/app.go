package testtask

import (
	"github.com/Alieksieiev0/test-task/reader"
	"github.com/Alieksieiev0/test-task/writer"
)

type App struct {
	source            KeyValueSource[string, string]
	processor         Processor[string, AsyncErrorProneEntry[string]]
	processorWrite    Processor[string, func(AsyncEntry[string]) AsyncEntry[error]]
	operationRead     Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]]
	operationWrite    Operation[Entry[error], AsyncEntry[error]]
	multiOperation    MultiOperation[AsyncErrorProneEntry[string], func(AsyncEntry[string]) AsyncEntry[error]]
	writer            *AppWriter
	parser            *AppParser
	readerFactory     ErrorFactoryFunc[string, reader.Reader[string]]
	asyncErrorFactory PlainFactory[AsyncEntry[error]]
	asyncEntryFactory PlainFactory[AsyncErrorProneEntry[string]]
	writerFactory     ErrorFactoryFunc[string, writer.Writer[string]]
}

func (a *App) NewRun() {
	readResults := a.processor.Process(
		a.source.Keys(),
		func(data string) AsyncErrorProneEntry[string] {
			results := a.asyncEntryFactory.Create()
			go func() {
				reader, err := a.readerFactory(data)
				if err != nil {
					results.PassErr(err)
					return
				}
				defer reader.Close()
				a.operationRead.Run(a.parser.Parse(reader), results)
			}()
			return results
		},
	)
	writeResults := a.processorWrite.Process(
		a.source.Values(),
		func(target string) func(AsyncEntry[string]) AsyncEntry[error] {
			return func(input AsyncEntry[string]) AsyncEntry[error] {
				results := a.asyncErrorFactory.Create()
				go func() {
					writer, err := a.writerFactory(target)
					if err != nil {
						results.PassVal(err)
						return
					}
					defer writer.Close()
					a.operationWrite.Run(a.writer.Writer(writer, input), results)
				}()
				return results
			}
		},
	)
	a.multiOperation.Run(readResults, writeResults)
}

func NewAppFactory() *AppFactory {
	return &AppFactory{}
}

type AppFactory struct {
}

func (a *AppFactory) Create(
	source KeyValueSource[string, string],
	processor Processor[string, AsyncErrorProneEntry[string]],
	processorWrite Processor[string, func(AsyncEntry[string]) AsyncEntry[error]],
	operationRead Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]],
	operationWrite Operation[Entry[error], AsyncEntry[error]],
	multiOperation MultiOperation[AsyncErrorProneEntry[string], func(AsyncEntry[string]) AsyncEntry[error]],
	writer *AppWriter,
	parser *AppParser,
	readerFactory ErrorFactoryFunc[string, reader.Reader[string]],
	asyncErrorFactory PlainFactory[AsyncEntry[error]],
	asyncEntryFactory PlainFactory[AsyncErrorProneEntry[string]],
	writerFactory ErrorFactoryFunc[string, writer.Writer[string]],
) *App {
	return &App{
		source:            source,
		processor:         processor,
		processorWrite:    processorWrite,
		operationRead:     operationRead,
		operationWrite:    operationWrite,
		multiOperation:    multiOperation,
		writer:            writer,
		parser:            parser,
		readerFactory:     readerFactory,
		asyncErrorFactory: asyncErrorFactory,
		asyncEntryFactory: asyncEntryFactory,
		writerFactory:     writerFactory,
	}
}
