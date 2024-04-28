package processor

import "github.com/Alieksieiev0/test-task/iterator"

func NewSequentialProcessor[Input, Output any]() *SequentialProcessor[Input, Output] {
	return &SequentialProcessor[Input, Output]{}
}

type SequentialProcessor[Input, Output any] struct {
}

func (p *SequentialProcessor[Input, Output]) Process(
	source iterator.Iterator[Input],
	mapF func(Input) Output,
) iterator.Iterator[Output] {
	return iterator.NewCallbackBased(func() *iterator.Step[Output] {
		step := source.Next()
		if step.Err() != nil {
			return iterator.NewStepErr[Output](step.Err())
		}
		return iterator.NewStepVal(mapF(step.Val()))
	})
}
