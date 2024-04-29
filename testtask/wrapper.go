package testtask

import (
	"context"
	"sync"

	"github.com/Alieksieiev0/test-task/iterator"
	"github.com/Alieksieiev0/test-task/reader"
	"github.com/Alieksieiev0/test-task/writer"
)

type AsyncParserWrapper struct {
	parser            ApplicationParser[string, ErrorProneEntry[string]]
	processor         Processor[string, func(context.Context) AsyncErrorProneEntry[string]]
	operationFactory  FactoryFunc[context.Context, Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]]]
	readerFactory     ErrorFactoryFunc[string, reader.Reader[string]]
	asyncEntryFactory PlainFactory[AsyncErrorProneEntry[string]]
	wg                *sync.WaitGroup
}

func (a *AsyncParserWrapper) Wrap(
	iterator iterator.Iterator[string],
) iterator.Iterator[func(context.Context) AsyncErrorProneEntry[string]] {
	return a.processor.Process(
		iterator,
		func(data string) func(context.Context) AsyncErrorProneEntry[string] {
			return func(ctx context.Context) AsyncErrorProneEntry[string] {
				results := a.asyncEntryFactory.Create()
				a.wg.Add(1)
				go func() {
					defer a.wg.Done()
					reader, err := a.readerFactory(data)
					if err != nil {
						results.PassErr(err)
						results.Close()
						return
					}
					defer reader.Close()
					a.operationFactory(ctx).Run(a.parser.Parse(reader), results)
				}()
				return results
			}
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
	processor Processor[string, func(context.Context) AsyncErrorProneEntry[string]],
	operationFactory FactoryFunc[context.Context, Operation[ErrorProneEntry[string], AsyncErrorProneEntry[string]]],
	readerFactory ErrorFactoryFunc[string, reader.Reader[string]],
	wg *sync.WaitGroup,
) iterator.IteratorWrapper[string, func(context.Context) AsyncErrorProneEntry[string]] {
	return &AsyncParserWrapper{
		parser:            parser,
		processor:         processor,
		operationFactory:  operationFactory,
		readerFactory:     readerFactory,
		asyncEntryFactory: NewAsyncStringEntryFactory(),
		wg:                wg,
	}
}

type AsyncWriterWrapper struct {
	writer            ApplicationWriter[string, AsyncEntry[string], Entry[error]]
	processor         Processor[string, func(context.Context, AsyncErrorProneEntry[string]) AsyncEntry[error]]
	operationFactory  FactoryFunc[context.Context, Operation[Entry[error], AsyncEntry[error]]]
	writerFactory     ErrorFactoryFunc[string, writer.Writer[string]]
	asyncErrorFactory PlainFactory[AsyncEntry[error]]
	wg                *sync.WaitGroup
}

func (a *AsyncWriterWrapper) Wrap(
	iterator iterator.Iterator[string],
) iterator.Iterator[func(context.Context, AsyncErrorProneEntry[string]) AsyncEntry[error]] {
	return a.processor.Process(
		iterator,
		func(target string) func(context.Context, AsyncErrorProneEntry[string]) AsyncEntry[error] {
			return func(ctx context.Context, input AsyncErrorProneEntry[string]) AsyncEntry[error] {
				results := a.asyncErrorFactory.Create()
				a.wg.Add(1)
				go func() {
					defer a.wg.Done()
					writer, err := a.writerFactory(target)
					if err != nil {
						results.PassVal(err)
						results.Close()
						return
					}
					defer writer.Close()
					a.operationFactory(ctx).Run(a.writer.Write(writer, input), results)
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
	processor Processor[string, func(context.Context, AsyncErrorProneEntry[string]) AsyncEntry[error]],
	operationFactory FactoryFunc[context.Context, Operation[Entry[error], AsyncEntry[error]]],
	writerFactory ErrorFactoryFunc[string, writer.Writer[string]],
	wg *sync.WaitGroup,
) iterator.IteratorWrapper[string, func(context.Context, AsyncErrorProneEntry[string]) AsyncEntry[error]] {
	return &AsyncWriterWrapper{
		writer:            writer,
		processor:         processor,
		operationFactory:  operationFactory,
		writerFactory:     writerFactory,
		asyncErrorFactory: NewAsyncErrorEntryFactory(),
		wg:                wg,
	}
}
