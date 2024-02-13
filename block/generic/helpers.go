package generic

import (
	"encoding/binary"
	"encoding/hex"
	"io"
	"log"
)

func readOneT[T any](r io.Reader, order binary.ByteOrder, ptrSize int) (_ *T, err error) {
	var body T
	return &body, Read(r, order, ptrSize, &body)
}

func readSliceT[T any](r io.Reader, order binary.ByteOrder, ptrSize int, count uint32) (_ []*T, err error) {
	bodies := make([]*T, count)
	for i := range bodies {
		bodies[i], err = readOneT[T](r, order, ptrSize)
		if err != nil {
			return nil, err
		}
	}

	return bodies, err
}

func ReadT[T any](r io.Reader, order binary.ByteOrder, ptrSize int, count uint32) (any, error) {
	if count == 1 {
		return readOneT[T](r, order, ptrSize)
	}
	return readSliceT[T](r, order, ptrSize, count)
}

func EnsureAllRead(r io.Reader, typ string) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if len(buf) == 0 {
		return nil
	}

	log.Printf("%d unread bytes in %q.", len(buf), typ)
	log.Println(hex.Dump(buf))

	return nil
}
