package block

import (
	"io"
)

// Header contains information about the block's type and size.
type Header struct {
	// Code provides a rough type description of the block.
	Code Code
	// Total length of the data after the block header.
	Size int64
	// Memory address of the structure when it was written to disk.
	OldAddr uint64
	// Index in the Structure DNA.
	SDNAIndex uint32
	// Number of structures located in this block.
	Count uint32
}

// ParseHeader parses and returns a file block header.
//
// Example file block header:
//
//	44 41 54 41  E0 00 00 00  88 5E 9D 04  00 00 00 00    DATA.....^......
//	F8 00 00 00  0E 00 00 00                              ........
//
//	//   0-3   block code   ("DATA")
//	//   4-7   size         (0x000000E0 = 224)
//	//  8-15   old addr     (0x00000000049D5E88) // size depends on PtrSize.
//	// 16-19   sdna index   (0x000000F8 = 248)
//	// 20-23   count        (0x0000000E = 14)
const headerSizeWithoutPtr = 16

func (r *Reader) ParseHeader(src readSeekerAt) (hdr Header, _ error) {
	// Block code.
	header := make([]byte, headerSizeWithoutPtr+r.PtrSize)
	if _, err := io.ReadFull(src, header); err != nil {
		return hdr, err
	}

	var offset int
	hdr.Code = parseCode(header[:4])
	offset += 4

	// Block size.
	hdr.Size = int64(r.Order.Uint32(header[offset:]))
	offset += 4

	// Old memory address.
	switch r.PtrSize {
	case 4:
		hdr.OldAddr = uint64(r.Order.Uint32(header[offset:]))
	case 8:
		hdr.OldAddr = r.Order.Uint64(header[offset:])
	}
	offset += r.PtrSize

	// SDNA index.
	hdr.SDNAIndex = r.Order.Uint32(header[offset:])
	offset += 4

	// Structure count.
	hdr.Count = r.Order.Uint32(header[offset:])
	offset += 4

	return hdr, nil
}

func (w *Writer) WriteHeader(dst io.Writer, hdr Header) error {
	// Block code.
	header := make([]byte, headerSizeWithoutPtr+w.PtrSize)

	var offset int
	copy(header[:4], hdr.Code)
	offset += 4

	// Block size.
	w.Order.PutUint32(header[offset:], uint32(hdr.Size))
	offset += 4

	// Old memory address.
	switch w.PtrSize {
	case 4:
		w.Order.PutUint32(header[offset:], uint32(hdr.OldAddr))
	case 8:
		w.Order.PutUint64(header[offset:], hdr.OldAddr)
	}
	offset += w.PtrSize

	// SDNA index.
	w.Order.PutUint32(header[offset:], hdr.SDNAIndex)
	offset += 4

	// Structure count.
	w.Order.PutUint32(header[offset:], hdr.Count)
	offset += 4

	_, err := dst.Write(header)
	return err
}
