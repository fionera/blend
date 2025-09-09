package generic

// BlockPointer is the memory address of a structure when it was written to disk.
type BlockPointer[Target any] struct {
	Addr uint64 `bin:"ptrSize"` //TODO: This is very very ugly
}

func (p BlockPointer[T]) Data() T {
	var t T
	return t
}

func (p BlockPointer[T]) Valid() bool {
	return p.Addr != 0
}
