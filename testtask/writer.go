package testtask

import (
	"github.com/Alieksieiev0/test-task/iterator"
	"github.com/Alieksieiev0/test-task/writer"
)

type DefaultWriter struct {
	processor     Processor[string, Entry[error]]
	sourceFactory FactoryFunc[<-chan string, Source[string]]
	entryFactory  Factory[error, Entry[error]]
}

func (w *DefaultWriter) Write(
	writer writer.Writer[string],
	data AsyncEntry[string],
) iterator.Iterator[Entry[error]] {
	return w.processor.Process(w.sourceFactory(data.Val()).Data(), func(data string) Entry[error] {
		return w.entryFactory.Create(writer.Write(data))
	})
}

func NewDefaultWriterFactory() *DefaultWriterFactory {
	return &DefaultWriterFactory{}
}

type DefaultWriterFactory struct {
}

func (a *DefaultWriterFactory) Create(
	processor Processor[string, Entry[error]],
	sourceFactory FactoryFunc[<-chan string, Source[string]],
) ApplicationWriter[string, AsyncEntry[string], Entry[error]] {
	return &DefaultWriter{
		processor:     processor,
		sourceFactory: sourceFactory,
		entryFactory:  NewErrorEntryFactory(),
	}
}
