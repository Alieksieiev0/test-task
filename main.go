package main

import (
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
	reader := testtask.NewStringReaderFactory().Create(
		processor.NewSequentialProcessor[string, testtask.ErrorProneEntry[string]](),
		operation.NewAsyncReadOperation[testtask.ErrorProneEntry[string], testtask.AsyncErrorProneEntry[string]](),
		parser.NewJsonParser[*testtask.MessageEntry](),
		reader.NewFileBasedFactory().Create('\n'),
		source.NewFileContentFactory().Create,
	)

	writer := testtask.NewStringWriterFactory().Create(
		processor.NewSequentialProcessor[string, error](),
		source.NewAsyncFileContentFactory().Create,
		operation.NewAsyncWriteOperation[error, testtask.AsyncEntry[error]](),
		writer.NewFileBasedFactory().Create('\n'),
	)

	app := testtask.NewAppFactory().Create(
		files,
		source.NewFileCollection(files),
		*reader,
		*writer,
	)

	app.Run()
}
