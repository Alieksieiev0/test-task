package testtask

import "github.com/Alieksieiev0/test-task/iterator"

type App struct {
	parser    iterator.IteratorWrapper[string, AsyncErrorProneEntry[string]]
	writer    iterator.IteratorWrapper[string, func(AsyncErrorProneEntry[string]) AsyncEntry[error]]
	source    KeyValueSource[string, string]
	operation CallbackOperation[AsyncErrorProneEntry[string], AsyncEntry[error]]
}

func (a *App) Run() {
	a.operation.Run(
		a.parser.Wrap(a.source.Keys()),
		a.writer.Wrap(a.source.Values()),
	)
}

func NewAppFactory() *AppFactory {
	return &AppFactory{}
}

type AppFactory struct {
}

func (a *AppFactory) Create(
	parser iterator.IteratorWrapper[string, AsyncErrorProneEntry[string]],
	writer iterator.IteratorWrapper[string, func(AsyncErrorProneEntry[string]) AsyncEntry[error]],
	source KeyValueSource[string, string],
	operation CallbackOperation[AsyncErrorProneEntry[string], AsyncEntry[error]],
) *App {
	return &App{
		parser:    parser,
		writer:    writer,
		source:    source,
		operation: operation,
	}
}
