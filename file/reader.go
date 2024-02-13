package file

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	seekable "github.com/SaveTheRbtz/zstd-seekable-format-go"
	"github.com/klauspost/compress/zstd"
	"go.uber.org/multierr"
)

type readSeekerAt interface {
	io.ReadSeeker
	io.ReaderAt
}

type Reader struct {
	readSeekerAt

	zstdCloser       func()
	zstdSeekerCloser func() error
}

func (d *Reader) Close() error {
	var mErr error
	if d.zstdSeekerCloser != nil {
		mErr = multierr.Append(mErr, d.zstdSeekerCloser())
	}

	if d.zstdCloser != nil {
		d.zstdCloser()
	}

	return mErr
}

func NewReader(src readSeekerAt) (*Reader, error) {
	header := make([]byte, 12)
	if _, err := io.ReadFull(src, header); err != nil {
		return nil, err
	}

	/* Rewind the file after reading the header. */
	if _, err := src.Seek(io.SeekStart, 0); err != nil {
		return nil, err
	}

	var r Reader

	// File identifier.
	magic := header[0:7]
	switch {
	case bytes.Equal(magic, []byte("BLENDER")):
		// nothing, just continue using
		r.readSeekerAt = src
	case magicIsGZIP(magic):
		return nil, fmt.Errorf("gzip compression is not implemented")

	case magicIsZSTD(magic):
		dec, err := zstd.NewReader(nil)
		if err != nil {
			return nil, err
		}
		r.zstdCloser = dec.Close

		seeker, err := seekable.NewReader(src, dec)
		if err != nil {
			return nil, err
		}
		r.zstdSeekerCloser = seeker.Close
		r.readSeekerAt = seeker
	}

	return &r, nil
}

func magicIsGZIP(header []byte) bool {
	/* GZIP itself starts with the magic bytes 0x1f 0x8b.
	 * The third byte indicates the compression method, which is 0x08 for DEFLATE. */
	return header[0] == 0x1f && header[1] == 0x8b && header[2] == 0x08
}

func magicIsZSTD(header []byte) bool {
	/* ZSTD files consist of concatenated frames, each either a ZSTD frame or a skippable frame.
	 * Both types of frames start with a magic number: `0xFD2FB528` for ZSTD frames and `0x184D2A5`
	 * for skippable frames, with the * being anything from 0 to F.
	 *
	 * To check whether a file is ZSTD-compressed, we just check whether the first frame matches
	 * either. Seeking through the file until a ZSTD frame is found would make things more
	 * complicated and the probability of a false positive is rather low anyways.
	 *
	 * Note that LZ4 uses a compatible format, so even though its compressed frames have a
	 * different magic number, a valid LZ4 file might also start with a skippable frame matching
	 * the second check here.
	 *
	 * For more details, see https://github.com/facebook/zstd/blob/dev/doc/zstd_compression_format.md
	 */

	magic := binary.LittleEndian.Uint32(header)
	if magic == 0xFD2FB528 {
		return true
	}
	if (magic >> 4) == 0x184D2A5 {
		return true
	}
	return false
}
