package testtask

import (
	"fmt"

	"github.com/Alieksieiev0/test-task/iterator"
	"github.com/Alieksieiev0/test-task/reader"
	"github.com/Alieksieiev0/test-task/writer"
)

type AsyncParserWrapper struct {
	parser            ApplicationParser[string, ErrorProneEntry[string]]
	processor         Processor[string, AsyncErrorProneEntry[string]]
	operation         Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]]
	readerFactory     ErrorFactoryFunc[string, reader.Reader[string]]
	asyncEntryFactory PlainFactory[AsyncErrorProneEntry[string]]
}

func (a *AsyncParserWrapper) Wrap(
	iterator iterator.Iterator[string],
) iterator.Iterator[AsyncErrorProneEntry[string]] {
	return a.processor.Process(
		iterator,
		func(data string) AsyncErrorProneEntry[string] {
			results := a.asyncEntryFactory.Create()
			go func() {
				reader, err := a.readerFactory(data)
				if err != nil {
					results.PassErr(err)
					return
				}
				defer reader.Close()
				a.operation.Run(a.parser.Parse(reader), results)
			}()
			fmt.Println(results)
			return results
		},
	)
}

func NewAsyncParserWrapperFactory() AsyncParserWrapperFactory {
	return AsyncParserWrapperFactory{}
}

type AsyncParserWrapperFactory struct {
}

func (a AsyncParserWrapperFactory) Create(
	parser ApplicationParser[string, ErrorProneEntry[string]],
	processor Processor[string, AsyncErrorProneEntry[string]],
	operation Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]],
	readerFactory ErrorFactoryFunc[string, reader.Reader[string]],
) iterator.IteratorWrapper[string, AsyncErrorProneEntry[string]] {
	return &AsyncParserWrapper{
		parser:            parser,
		processor:         processor,
		operation:         operation,
		readerFactory:     readerFactory,
		asyncEntryFactory: NewAsyncStringEntryFactory(),
	}
}

type AsyncWriterWrapper struct {
	writer            ApplicationWriter[string, AsyncEntry[string], Entry[error]]
	processor         Processor[string, func(AsyncErrorProneEntry[string]) AsyncEntry[error]]
	operation         Operation[Entry[error], AsyncEntry[error]]
	writerFactory     ErrorFactoryFunc[string, writer.Writer[string]]
	asyncErrorFactory PlainFactory[AsyncEntry[error]]
}

func (a *AsyncWriterWrapper) Wrap(
	iterator iterator.Iterator[string],
) iterator.Iterator[func(AsyncErrorProneEntry[string]) AsyncEntry[error]] {
	return a.processor.Process(
		iterator,
		func(target string) func(AsyncErrorProneEntry[string]) AsyncEntry[error] {
			return func(input AsyncErrorProneEntry[string]) AsyncEntry[error] {
				results := a.asyncErrorFactory.Create()
				go func() {
					writer, err := a.writerFactory(target)
					if err != nil {
						results.PassVal(err)
						return
					}
					defer writer.Close()
					a.operation.Run(a.writer.Write(writer, input), results)
				}()
				return results
			}
		},
	)
}

func NewAsyncWriterWrapperFactory() *AsyncWriterWrapperFactory {
	return &AsyncWriterWrapperFactory{}
}

type AsyncWriterWrapperFactory struct {
}

func (a *AsyncWriterWrapperFactory) Create(
	writer ApplicationWriter[string, AsyncEntry[string], Entry[error]],
	processor Processor[string, func(AsyncErrorProneEntry[string]) AsyncEntry[error]],
	operation Operation[Entry[error], AsyncEntry[error]],
	writerFactory ErrorFactoryFunc[string, writer.Writer[string]],
) iterator.IteratorWrapper[string, func(AsyncErrorProneEntry[string]) AsyncEntry[error]] {
	return &AsyncWriterWrapper{
		writer:            writer,
		processor:         processor,
		operation:         operation,
		writerFactory:     writerFactory,
		asyncErrorFactory: NewAsyncErrorEntryFactory(),
	}
}
