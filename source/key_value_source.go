package source

import "github.com/Alieksieiev0/test-task/iterator"

type SortedFilePairs struct {
	keys   *FileCollection
	values *FileCollection
}

func (s *SortedFilePairs) Keys() iterator.Iterator[string] {
	return s.keys.Data()
}

func (s *SortedFilePairs) Values() iterator.Iterator[string] {
	return s.values.Data()
}
