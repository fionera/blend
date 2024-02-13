// Package blend implements parsing of Blender files.
package blend

import (
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
)

// A Header is present at the beginning of each blend file.
// Example file header:
//
//	"BLENDER_V100"
//	//  0-6   magic ("BLENDER")
//	//    7   pointer size ("_" or "-")
//	//    8   endianness ("V" or "v")
//	// 9-11   version ("100")
type Header struct {
	// Pointer size.
	PtrSize int
	// Byte order.
	Order binary.ByteOrder
	// Blender version.
	Ver int
}

const headerSize = 12
const headerMagic = "BLENDER"

// ReadHeader parses and returns the header of a blend file.
func ReadHeader(r io.Reader) (hdr Header, err error) {
	// create a section reader for the header area to
	// allow multiple calls of this function. Totally
	// unneeded but who cares lol.
	//sr := io.NewSectionReader(d, 0, headerSize)

	var buf [headerSize]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return hdr, err
	}

	// File identifier.
	if magic := buf[0:7]; string(magic) != headerMagic {
		return hdr, fmt.Errorf("invalid file identifier: %q", magic)
	}

	// Pointer size.
	switch size := buf[7]; size {
	case '_':
		// _ = 4 byte pointer
		hdr.PtrSize = 4
	case '-':
		// - = 8 byte pointer
		hdr.PtrSize = 8
	default:
		return hdr, fmt.Errorf("invalid pointer size character: %q", size)
	}

	// Byte order.
	switch order := buf[8]; order {
	case 'v':
		// v = little endian
		hdr.Order = binary.LittleEndian
	case 'V':
		// V = big endian
		hdr.Order = binary.BigEndian
	default:
		return hdr, fmt.Errorf("invalid byte order character: %q", order)
	}

	// Version.
	hdr.Ver, err = strconv.Atoi(string(buf[9:12]))
	if err != nil {
		return hdr, fmt.Errorf("invalid version: %s", err)
	}

	return
}

func WriteHeader(w io.Writer, hdr Header) error {
	var buf [headerSize]byte

	// File identifier.
	copy(buf[0:7], headerMagic)

	// Pointer size.
	switch hdr.PtrSize {
	case 4:
		// _ = 4 byte pointer
		buf[7] = '_'
	case 8:
		// - = 8 byte pointer
		buf[7] = '-'
	default:
		return fmt.Errorf("invalid pointer size: %q", hdr.PtrSize)
	}

	// Byte order.
	switch hdr.Order {
	case binary.LittleEndian:
		// v = little endian
		buf[8] = 'v'
	case binary.BigEndian:
		// V = big endian
		buf[8] = 'V'
	default:
		return fmt.Errorf("invalid byte order: %q", hdr.Order)
	}

	// Version.
	copy(buf[9:12], strconv.Itoa(hdr.Ver))

	_, err := w.Write(buf[:])
	return err
}
