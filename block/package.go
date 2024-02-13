// Package block implements parsing of blend file blocks.
//
// One unique feature of blend files is that they contain a full definition of
// every structure used in its file blocks. The structure definitions are stored
// in the DNA block.
//
// All block structure definitions ("struct.go") and the block parsing logic
// ("parse.go") have been generating by parsing the DNA block of
// "testdata/block.blend".
//
// The tool which was used to generate these two files is available through:
//
//	go get github.com/mewspring/blend/cmd/blendef
//
// More complex blend files may contain structures which are not yet defined in
// this package. If so, use blendef to regenerate "struct.go" and "parse.go" for
// the given blend file.
package block