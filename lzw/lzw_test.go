package lzw_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/d-jo/lzw-go/lzw"
)

func TestEncode(t *testing.T) {
	input := "itty bitty bit bin"

	outBuffer := bytes.NewBuffer(nil)

	lzwObj := lzw.NewLZW()

	//lzwObj.InitDict(input)

	lzwObj.Encode(outBuffer, input)

	d := lzwObj.GetDict()
	t.Log("Dict")
	t.Logf("%+v", d)

	inBuffer := make([]byte, 2)
	for {
		outBuffer.Read(inBuffer)

		translated := binary.LittleEndian.Uint16(inBuffer)
		t.Logf("%03d \t %s", translated, d[translated])

		if outBuffer.Len() == 0 {
			break
		}

	}

	//t.Logf("Encoded: %+v", outBuffer)

	t.Fail()
}
