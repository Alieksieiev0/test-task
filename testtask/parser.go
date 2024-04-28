package testtask

import (
	"io"
	"strings"

	"github.com/Alieksieiev0/test-task/iterator"
	"github.com/Alieksieiev0/test-task/reader"
)

type DefaultParser struct {
	processor     Processor[string, ErrorProneEntry[string]]
	sourceFactory FactoryFunc[reader.Reader[string], Source[string]]
	entryFactory  Factory[io.Reader, ErrorProneEntry[string]]
}

func (a *DefaultParser) Parse(
	reader reader.Reader[string],
) iterator.Iterator[ErrorProneEntry[string]] {
	return a.processor.Process(
		a.sourceFactory(reader).Data(),
		func(data string) ErrorProneEntry[string] {
			return a.entryFactory.Create(strings.NewReader(data))
		},
	)
}

func NewDefaultParserFactory() DefaultParserFactory {
	return DefaultParserFactory{}
}

type DefaultParserFactory struct {
}

func (a DefaultParserFactory) Create(
	processor Processor[string, ErrorProneEntry[string]],
	sourceFactory FactoryFunc[reader.Reader[string], Source[string]],
	parser Parser[io.Reader, *MessageEntry],
) ApplicationParser[string, ErrorProneEntry[string]] {
	return &DefaultParser{
		processor:     processor,
		sourceFactory: sourceFactory,
		entryFactory:  NewMessageEntryFactory(parser),
	}
}
