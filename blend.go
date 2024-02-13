package blend

import (
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/mewspring/blend/block"
	"github.com/mewspring/blend/file"
)

// Blend represents the information contained within a blend file. It contains a
// file header and a slice of file blocks.
type Blend struct {
	Hdr     Header
	Blocks  []*block.Block
	OldAddr map[uint64]*block.Block
}

func Decode(d *file.Reader) (*Blend, error) {
	b := &Blend{OldAddr: make(map[uint64]*block.Block)}

	// Parse file header.
	var err error
	b.Hdr, err = ReadHeader(d)
	if err != nil {
		return nil, fmt.Errorf("reading header: %v", err)
	}

	blkReader := block.NewReader(b.Hdr.Order, b.Hdr.PtrSize, b.Hdr.Ver)
	// Parse file blocks.
	for {
		blk, err := blkReader.ReadBlock(d)
		if err != nil {
			return nil, fmt.Errorf("reading block: %v", err)
		}

		if blk.Hdr.Code == block.CodeENDB {
			break
		}
		b.OldAddr[blk.Hdr.OldAddr] = blk

		b.Blocks = append(b.Blocks, blk)
	}

	return b, nil
}

func Encode(dst io.Writer, b *Blend) error {
	if err := WriteHeader(dst, b.Hdr); err != nil {
		return err
	}

	w := &block.Writer{
		PtrSize: b.Hdr.PtrSize,
		Order:   b.Hdr.Order,
	}

	for _, blk := range b.Blocks {
		if err := w.WriteBlock(dst, blk); err != nil {
			return err
		}
	}

	return w.WriteBlock(dst, &block.Block{
		Hdr: block.Header{
			Code:  block.CodeENDB,
			Size:  16,
			Count: 1,
		},
		Body: nil,
	})
}

// GetDNA locates, parses and returns the DNA block.
func (b *Blend) GetDNA() (dna *block.DNA, err error) {
	for _, blk := range b.Blocks {
		dna, ok := blk.Body.(*block.DNA)
		if ok {
			// DNA block already parsed.
			return dna, nil

		}
		if blk.Hdr.Code == block.CodeDNA1 {
			log.Printf("%+v", blk.Hdr)
			err := blk.ParseBody(nil)
			if err != nil {
				return nil, err
			}
			return blk.Body.(*block.DNA), nil
		}
	}
	return nil, errors.New("Blend.GetDNA: unable to locate DNA block")
}
