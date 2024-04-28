package main

import (
	"time"

	"github.com/Alieksieiev0/test-task/operation"
	"github.com/Alieksieiev0/test-task/parser"
	"github.com/Alieksieiev0/test-task/processor"
	"github.com/Alieksieiev0/test-task/reader"
	"github.com/Alieksieiev0/test-task/source"
	"github.com/Alieksieiev0/test-task/testtask"
	"github.com/Alieksieiev0/test-task/writer"
)

func main() {
	files := map[string]string{
		"file1.txt": "out1.txt",
		"file2.txt": "out2.txt",
		"file3.txt": "out3.txt",
	}
	parser := testtask.NewAsyncParserWrapperFactory().Create(
		testtask.NewDefaultParserFactory().Create(
			processor.NewSequentialProcessor[string, testtask.ErrorProneEntry[string]](),
			source.NewFileContentFactory().Create,
			parser.NewJsonParser[*testtask.MessageEntry](),
		),
		processor.NewSequentialProcessor[string, testtask.AsyncErrorProneEntry[string]](),
		operation.NewAsyncReadOperation[testtask.ErrorProneEntry[string], testtask.AsyncErrorProneEntry[string]](
			time.Millisecond*5000,
		),
		reader.NewFileBasedFactory().Create('\n'),
	)

	writer := testtask.NewAsyncWriterWrapperFactory().Create(
		testtask.NewDefaultWriterFactory().Create(
			processor.NewSequentialProcessor[string, testtask.Entry[error]](),
			source.NewAsyncFileContentFactory().Create,
		),
		processor.NewSequentialProcessor[string, func(testtask.AsyncErrorProneEntry[string]) testtask.AsyncEntry[error]](),
		operation.NewAsyncWriteOperation[testtask.Entry[error], testtask.AsyncEntry[error]](),
		writer.NewFileBasedFactory().Create('\n'),
	)

	app := testtask.NewAppFactory().Create(
		parser,
		writer,
		source.NewSortedFilePairsFactory().CreateFromMap(files),
		operation.NewAsyncFileDecoder[testtask.AsyncErrorProneEntry[string], testtask.AsyncEntry[error]](),
	)

	app.Run()
}
