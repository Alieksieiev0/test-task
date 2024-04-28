package testtask

import (
	"github.com/Alieksieiev0/test-task/iterator"
	"github.com/Alieksieiev0/test-task/writer"
)

type AppWriter struct {
	processor     Processor[string, Entry[error]]
	sourceFactory FactoryFunc[<-chan string, Source[string]]
	entryFactory  Factory[error, Entry[error]]
}

func (w *AppWriter) Writer(
	writer writer.Writer[string],
	data AsyncEntry[string],
) iterator.Iterator[Entry[error]] {
	return w.processor.Process(w.sourceFactory(data.Val()).Data(), func(data string) Entry[error] {
		return w.entryFactory.Create(writer.Write(data))
	})
}

func NewAppWriterFactory() *AppWriterFactory {
	return &AppWriterFactory{}
}

type AppWriterFactory struct {
}

func (a *AppWriterFactory) Create(
	processor Processor[string, Entry[error]],
	sourceFactory FactoryFunc[<-chan string, Source[string]],
) *AppWriter {
	return &AppWriter{
		processor:     processor,
		sourceFactory: sourceFactory,
		entryFactory:  NewErrorEntryFactory(),
	}
}
