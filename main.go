package main

import (
	"github.com/Alieksieiev0/test-task/parser"
	"github.com/Alieksieiev0/test-task/processor"
	"github.com/Alieksieiev0/test-task/source"
	"github.com/Alieksieiev0/test-task/testtask"
)

func main() {
	files := map[string]string{
		"file1.txt": "out1.txt",
		"file2.txt": "out2.txt",
		"file3.txt": "out3.txt",
	}
	reader := testtask.NewAppParserFactory().Create(
		processor.NewSequentialProcessor[string, testtask.ErrorProneEntry[string]](),
		source.NewFileContentFactory().Create,
		parser.NewJsonParser[*testtask.MessageEntry](),
	)

	writer := testtask.NewAppWriterFactory().Create(
		processor.NewSequentialProcessor[string, testtask.Entry[error]](),
		source.NewAsyncFileContentFactory().Create,
	)

	app := testtask.NewAppFactory().Create(
		files,
		source.NewFileCollection(files),
		*reader,
		*writer,
	)

	app.Run()
}
