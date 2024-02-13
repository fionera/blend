package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unsafe"

	"github.com/mewspring/blend"
	"github.com/mewspring/blend/block"
	"github.com/mewspring/blend/block/generic"
	v400 "github.com/mewspring/blend/block/v400"
	"github.com/mewspring/blend/file"
)

func init() {
	flag.Usage = usage
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: inspect FILE.blend")
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	if flag.NArg() != 1 {
		log.Printf("invalid argument count.")
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder, err := file.NewReader(f)
	if err != nil {
		log.Fatal(err)
	}
	defer decoder.Close()

	b, err := blend.Decode(decoder)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", b.Hdr)

	newBlocks := make([]*block.Block, len(b.Blocks))
	copy(newBlocks, b.Blocks)

	dna, err := b.GetDNA()
	if err != nil {
		log.Fatal(err)
	}

	for i, blk := range b.Blocks {
		switch blk.Hdr.Code {
		case block.CodeSC, block.CodeIM, block.CodeMA:
			break
		default:
			continue
		}

		if err := blk.ParseBody(dna); err != nil {
			log.Fatal(err)
		}

		switch body := blk.Body.(type) {
		case *v400.Image:
			path := int8SliceToString(body.Name[:])
			if len(path) == 0 {
				log.Printf("empty name: %+v", int8SliceToString(body.Id.Name[:]))
				continue
			}

			if body.Packedfile.Addr != 0 {
				data := b.OldAddr[body.Packedfile.Addr]
				if err := data.ParseBody(dna); err != nil {
					log.Fatal(err)
				}

				log.Printf("pfHeader: %+v", data.Hdr)

				pf := data.Body.(*v400.PackedFile)
				pfData := b.OldAddr[pf.Data.Addr]

				if err := pfData.ParseBody(dna); err != nil {
					log.Fatal(err)
				}

				log.Printf("pfDataHeader: %+v; ", pfData.Hdr)
				log.Printf("pfData: %T; Size: %d; Seek: %d", pfData.Body, pf.Size, pf.Seek)
			} else {
				log.Printf("MISSING!")

				ab, bb := packFile(path, dna, b, body)
				newBlocks = append(newBlocks[:i], append([]*block.Block{ab, bb}, newBlocks[i:]...)...)
			}

			log.Println(path, body.Packedfile)
		default:
			log.Printf("unhandled: %T", body)
		}
	}

	b.Blocks = newBlocks
	for _, blk := range b.Blocks {
		log.Printf("%+v", blk.Hdr)
	}

	nf, err := os.OpenFile(os.Args[1]+"_new.blend", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer nf.Close()

	if err := blend.Encode(nf, b); err != nil {
		log.Fatal(err)
	}
}

func packFile(path string, dna *block.DNA, b *blend.Blend, body *v400.Image) (*block.Block, *block.Block) {
	open, err := os.Open("/Users/fionera/Downloads/_" + strings.ReplaceAll(path, "//", ""))
	if err != nil {
		log.Fatal(err)
	}

	all, err := io.ReadAll(open)
	if err != nil {
		log.Fatal(err)
	}

	var idx int
	for i, s := range dna.Structs {
		if s.Type == "PackedFile" {
			idx = i
			break
		}
	}

	dataBlock := &block.Block{
		Hdr: block.Header{
			Code:      block.CodeDATA,
			Size:      int64(len(all)),
			OldAddr:   0,
			SDNAIndex: 0, // is set to 0
			Count:     1,
		},
		Body: all,
	}
	dataBlock.Hdr.OldAddr = uint64(uintptr(unsafe.Pointer(&dataBlock))) * 2
	log.Printf("dataBlock: %+v", dataBlock.Hdr)

	pfBlock := &block.Block{
		Hdr: block.Header{
			Code:      block.CodeDATA,
			Size:      int64(8 + b.Hdr.PtrSize),
			OldAddr:   0,
			SDNAIndex: uint32(idx),
			Count:     1,
		},
		Body: v400.PackedFile{
			Size: int32(dataBlock.Hdr.Size),
			Seek: 0,
			Data: generic.BlockPointer[*any](generic.BlockPointer[[]uint8]{
				Addr: dataBlock.Hdr.OldAddr,
			}),
		},
	}
	pfBlock.Hdr.OldAddr = uint64(uintptr(unsafe.Pointer(&pfBlock))) * 2
	log.Printf("pfHeader: %+v", pfBlock.Hdr)

	body.Packedfile.Addr = pfBlock.Hdr.OldAddr
	body.Packedfiles.First.Addr = pfBlock.Hdr.OldAddr
	body.Packedfiles.Last.Addr = pfBlock.Hdr.OldAddr

	return pfBlock, dataBlock
}

func int8SliceToString(s []uint8) string {
	var sb strings.Builder
	for _, v := range s {
		if v == 0 {
			break
		}
		sb.WriteByte(v)
	}
	return sb.String()
}
