package block

import (
	"encoding/binary"
	"io"

	v305 "github.com/mewspring/blend/block/v305"
	v400 "github.com/mewspring/blend/block/v400"
	v401 "github.com/mewspring/blend/block/v401"
)

type Parser struct {
	ParseStructure func(r io.Reader, order binary.ByteOrder, ptrSize int, typ string, count uint32) (body any, err error)
}

var Versions = map[int]Parser{
	v305.BlenderVer: {
		ParseStructure: v305.ParseStructure,
	},
	v400.BlenderVer: {
		ParseStructure: v400.ParseStructure,
	},
	v401.BlenderVer: {
		ParseStructure: v401.ParseStructure,
	},
}
