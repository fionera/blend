package block

import (
	"log"
	"strings"
)

// Code represents a rough type description of a block.
type Code string

func (c Code) String() string {
	cs := string(c)
	if idx := strings.IndexByte(cs, 0x00); idx != -1 {
		cs = cs[:idx]
	}
	return cs
}

// TODO: use codegen for this
func parseCode(code []byte) Code {
	switch string(code) {
	case CodeAR:
		return CodeAR
	case CodeBR:
		return CodeBR
	case CodeCA:
		return CodeCA
	case CodeDATA:
		return CodeDATA
	case CodeDNA1:
		return CodeDNA1
	case CodeENDB:
		return CodeENDB
	case CodeGLOB:
		return CodeGLOB
	case CodeIM:
		return CodeIM
	case CodeLA:
		return CodeLA
	case CodeLS:
		return CodeLS
	case CodeMA:
		return CodeMA
	case CodeME:
		return CodeME
	case CodeOB:
		return CodeOB
	case CodeREND:
		return CodeREND
	case CodeSC:
		return CodeSC
	case CodeSN:
		return CodeSN
	case CodeSR:
		return CodeSR
	case CodeTE:
		return CodeTE
	case CodeTEST:
		return CodeTEST
	case CodeTX:
		return CodeTX
	case CodeWM:
		return CodeWM
	case CodeWO:
		return CodeWO
	case CodeAC:
		return CodeAC
	case CodeNT:
		return CodeNT
	case CodeSO:
		return CodeSO
	case CodeGR:
		return CodeGR
	case CodePL:
		return CodePL
	case CodeWS:
		return CodeWS
	case CodeVF:
		return CodeVF
	case CodeLI:
		return CodeLI
	case CodeID:
		return CodeID
	case CodeCU:
		return CodeCU
	default:
		log.Printf("block code not implemented:  %q", code)
	}

	return Code(code)
}

// Block codes.
const (
	CodeAR   = "AR\x00\x00"
	CodeBR   = "BR\x00\x00"
	CodeCA   = "CA\x00\x00"
	CodeDATA = "DATA"
	CodeDNA1 = "DNA1"
	CodeENDB = "ENDB"
	CodeGLOB = "GLOB"
	CodeIM   = "IM\x00\x00"
	CodeLA   = "LA\x00\x00"
	CodeLS   = "LS\x00\x00"
	CodeMA   = "MA\x00\x00"
	CodeME   = "ME\x00\x00"
	CodeOB   = "OB\x00\x00"
	CodeREND = "REND"
	CodeSC   = "SC\x00\x00"
	CodeSN   = "SN\x00\x00"
	CodeSR   = "SR\x00\x00"
	CodeTE   = "TE\x00\x00"
	CodeTEST = "TEST"
	CodeTX   = "TX\x00\x00"
	CodeWM   = "WM\x00\x00"
	CodeWO   = "WO\x00\x00"
	CodeAC   = "AC\x00\x00"
	CodeNT   = "NT\x00\x00"
	CodeSO   = "SO\x00\x00"
	CodeGR   = "GR\x00\x00"
	CodePL   = "PL\x00\x00"
	CodeWS   = "WS\x00\x00"
	CodeVF   = "VF\x00\x00"
	CodeLI   = "LI\x00\x00"
	CodeID   = "ID\x00\x00"
	CodeCU   = "CU\x00\x00"
)
