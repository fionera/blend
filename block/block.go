package block

import (
	"fmt"
	"io"
	"log"

	"github.com/mewspring/blend/block/generic"
	v400 "github.com/mewspring/blend/block/v400"
)

// A Block contains a header and a type dependent body.
type Block struct {
	Hdr  Header
	Body any

	sr *io.SectionReader
	r  *Reader
	w  *Writer
}

// ParseBody parses the block body and stores it in blk.Body. It is safe to call
// ParseBody multiple times on the same block.
func (blk *Block) ParseBody(dna *DNA) (err error) {
	if blk.Body != nil {
		// Body has already been parsed.
		return nil
	}

	index := blk.Hdr.SDNAIndex
	if index == 0 {
		// Parse based on block code.
		switch blk.Hdr.Code {
		case CodeDATA:
			blk.Body, err = io.ReadAll(blk.sr)
			if err != nil {
				return err
			}
		case CodeDNA1:
			blk.Body, err = ParseDNA(blk.sr, blk.r.Order)
			if err != nil {
				return err
			}
		case CodeREND, CodeTEST:
			/// TODO: implement specific block body parsing for REND and TEST.
			blk.Body, err = io.ReadAll(blk.sr)
			if err != nil {
				return err
			}
		default:
			err = fmt.Errorf("Block.ParseBody: parsing of %q not yet implemented", blk.Hdr.Code)
		}

		return
	}

	// Parse based on SDNA index.
	typ := dna.Structs[index].Type
	blk.Body, err = blk.r.Parser.ParseStructure(blk.sr, blk.r.Order, blk.r.PtrSize, typ, blk.Hdr.Count)
	return
}

func (blk *Block) WriteBody(dst io.Writer) error {
	if blk.Body == nil {
		return fmt.Errorf("nil body cant be written")
	}

	index := blk.Hdr.SDNAIndex
	if index == 0 {
		// Parse based on block code.
		switch blk.Hdr.Code {
		case CodeDATA:
			_, err := dst.Write(blk.Body.([]byte))
			if err != nil {
				return err
			}
		case CodeDNA1:
			// TODO: Re Encode DNA?
			if _, err := blk.sr.Seek(0, io.SeekStart); err != nil {
				return fmt.Errorf("failed seeking: %v", err)
			}

			_, err := io.Copy(dst, blk.sr)
			if err != nil {
				return err
			}
		case CodeREND, CodeTEST:
			/// TODO: implement specific block body writing for REND and TEST.
			_, err := dst.Write(blk.Body.([]byte))
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("Block.ParseBody: writing of %q not yet implemented", blk.Hdr.Code)
		}

		return nil
	}

	if img, ok := blk.Body.(v400.Image); ok {
		log.Println(img.Packedfile)
	}
	return generic.Write(dst, blk.w.Order, blk.w.PtrSize, blk.Body)
}
