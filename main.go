package main

import (
	"context"
	"math/rand"
	"sync"
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
	var (
		files = map[string]string{
			"file1.txt": "out1.txt",
			"file2.txt": "out2.txt",
			"file3.txt": "out3.txt",
		}
		wg sync.WaitGroup
	)

	parser := testtask.NewAsyncParserWrapperFactory().Create(
		testtask.NewDefaultParserFactory().Create(
			processor.NewSequentialProcessor[string, testtask.ErrorProneEntry[string]](),
			source.NewFileContentFactory().Create,
			parser.NewJsonParser[*testtask.MessageEntry](),
		),
		processor.NewSequentialProcessor[string, func(context.Context) testtask.AsyncErrorProneEntry[string]](),
		operation.NewAsyncReadOperationFactory().
			Create(func() time.Duration {
				return time.Duration(rand.Intn(5000) * int(time.Millisecond))
			}),
		reader.NewFileBasedFactory().Create('\n'),
		&wg,
	)

	writer := testtask.NewAsyncWriterWrapperFactory().Create(
		testtask.NewDefaultWriterFactory().Create(
			processor.NewSequentialProcessor[string, testtask.Entry[error]](),
			source.NewAsyncFileContentFactory().Create,
		),
		processor.NewSequentialProcessor[string, func(context.Context, testtask.AsyncErrorProneEntry[string]) testtask.AsyncEntry[error]](),
		operation.NewAsyncWriteOperationFactory().Create(),
		writer.NewFileBasedFactory().Create('\n'),
		&wg,
	)

	app := testtask.NewAppFactory().Create(
		parser,
		writer,
		source.NewSortedFilePairsFactory().CreateFromMap(files),
		operation.NewAsyncFileDecoder[testtask.AsyncErrorProneEntry[string], testtask.AsyncEntry[error]](
			&wg,
		),
	)

	app.Run()

}
