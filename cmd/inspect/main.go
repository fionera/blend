package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mewspring/blend"
	"github.com/mewspring/blend/block"
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

	dna, err := b.GetDNA()
	if err != nil {
		log.Fatal(err)
	}

	for _, blk := range b.Blocks {
		switch blk.Hdr.Code {
		case block.CodeSC, block.CodeIM:
			break
		default:
			continue
		}

		if err := blk.ParseBody(dna); err != nil {
			log.Fatal(err)
		}

		switch body := blk.Body.(type) {
		case *v400.PackedFile:
			log.Println("file", body)
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

				pf := data.Body.(*v400.PackedFile)

				pfData := b.OldAddr[pf.Data.Addr]
				if err := pfData.ParseBody(dna); err != nil {
					log.Fatal(err)
				}

				log.Printf("pfData: %T; Size: %d; Seek: %d", pfData.Body, pf.Size, pf.Seek)
			}

			log.Println(path, body.Packedfile)
		default:
			log.Printf("unhandled: %T", body)
		}
	}
}

func int8SliceToString(s []uint8) string {
	var sb strings.Builder
	for _, v := range s {
		if v != 0 {
			sb.WriteByte(byte(v))
		}
	}
	return sb.String()
}
