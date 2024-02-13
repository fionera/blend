package block

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type readSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
}

type Reader struct {
	PtrSize int
	Order   binary.ByteOrder
	Parser  Parser

	// Pointers is a map from the memory address of a structure (when it was written to
	// disk) to its file block.
	Pointers map[uint64]*Block
}

func NewReader(order binary.ByteOrder, ptrSize int, version int) *Reader {
	r := &Reader{
		PtrSize:  ptrSize,
		Order:    order,
		Pointers: make(map[uint64]*Block),
	}

	s, ok := Versions[version]
	if !ok {
		log.Printf("Warning: Version mismatch: %d not supported.\n", version)
		log.Println("Use blendef [1] to regenerate the block package.")
		r.Parser = Versions[400]
	}
	r.Parser = s

	return r
}

// ReadBlock parses and returns a file block.
func (r *Reader) ReadBlock(src readSeekerAt) (blk *Block, err error) {
	blk = new(Block)
	blk.r = r

	// Parse block header.
	blk.Hdr, err = r.ParseHeader(src)
	if err != nil {
		return nil, fmt.Errorf("parsing header: %v", err)
	}

	v, ok := r.Pointers[blk.Hdr.OldAddr]
	if ok && blk.Hdr != v.Hdr {
		log.Println(fmt.Errorf("Reader.ReadBlock: multiple occurances of struct address %#x", blk.Hdr.OldAddr))
		log.Println(v.Hdr, blk.Hdr)
	}
	r.Pointers[blk.Hdr.OldAddr] = blk

	// Store section reader for block body.
	off, err := src.Seek(blk.Hdr.Size, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	blk.sr = io.NewSectionReader(src, off-blk.Hdr.Size, blk.Hdr.Size)

	return blk, nil
}

type Writer struct {
	PtrSize int
	Order   binary.ByteOrder
}

func (w *Writer) WriteBlock(dst io.Writer, blk *Block) error {
	blk.w = w
	// Parse block header.
	if err := w.WriteHeader(dst, blk.Hdr); err != nil {
		return err
	}

	if blk.Hdr.Code == CodeENDB {
		return nil
	}

	// Untouched body
	if blk.Body == nil {
		if _, err := blk.sr.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed seeking: %v", err)
		}
		_, err := io.Copy(dst, blk.sr)
		return err
	}

	return blk.WriteBody(dst)
}
