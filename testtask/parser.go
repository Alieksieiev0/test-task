package testtask

import (
	"io"
	"strings"

	"github.com/Alieksieiev0/test-task/iterator"
	"github.com/Alieksieiev0/test-task/reader"
)

type AppParser struct {
	processor     Processor[string, ErrorProneEntry[string]]
	sourceFactory FactoryFunc[reader.Reader[string], Source[string]]
	entryFactory  Factory[io.Reader, ErrorProneEntry[string]]
}

func (a *AppParser) Parse(
	reader reader.Reader[string],
) iterator.Iterator[ErrorProneEntry[string]] {
	return a.processor.Process(
		a.sourceFactory(reader).Data(),
		func(data string) ErrorProneEntry[string] {
			return a.entryFactory.Create(strings.NewReader(data))
		},
	)
}

func NewAppParserFactory() AppParserFactory {
	return AppParserFactory{}
}

type AppParserFactory struct {
}

func (a AppParserFactory) Create(
	processor Processor[string, ErrorProneEntry[string]],
	sourceFactory FactoryFunc[reader.Reader[string], Source[string]],
	parser Parser[io.Reader, *MessageEntry],
) *AppParser {
	return &AppParser{
		processor:     processor,
		sourceFactory: sourceFactory,
		entryFactory:  NewMessageEntryFactory(parser),
	}
}
