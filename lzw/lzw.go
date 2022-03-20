package lzw

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

const (
	START_CODE      = uint16(256)
	STOP_CODE       = uint16(257)
	NEW_DICT_OFFSET = uint16(258)
)

type LZW struct {
	codeDict        map[uint16]string
	reverseCodeDict map[string]uint16
	currOffset      uint16
	mux             sync.RWMutex
}

func NewLZW() *LZW {
	lzw := &LZW{
		codeDict:        make(map[uint16]string),
		reverseCodeDict: make(map[string]uint16),
		currOffset:      0,
	}

	lzw.Put(START_CODE, "START")
	lzw.Put(STOP_CODE, "STOP")

	return lzw
}

func (lzw *LZW) GetDict() map[uint16]string {
	return lzw.codeDict
}

/*
	Associates a code with a string in the codeDict
	and reverseCodeDict
*/
func (lzw *LZW) Put(code uint16, str string) {
	lzw.mux.Lock()
	defer lzw.mux.Unlock()
	lzw.codeDict[code] = str
	lzw.reverseCodeDict[str] = code
}

/*
	Gets the string associated with the code from the dict
*/
func (lzw *LZW) Get(code uint16) (string, bool) {
	lzw.mux.RLock()
	defer lzw.mux.RUnlock()
	str, ok := lzw.codeDict[code]
	return str, ok
}

/*
	Gets the code associated with a string from the reverse dict
*/
func (lzw *LZW) ReverseGet(str string) (uint16, bool) {
	lzw.mux.RLock()
	defer lzw.mux.RUnlock()
	code, ok := lzw.reverseCodeDict[str]
	return code, ok
}

/*
	Gets the next available code for a string and returns the value
*/
func (lzw *LZW) GetCodeForStr(str string) uint16 {
	if len(str) == 1 {
		return uint16(rune(str[0]))
	}
	retCode := NEW_DICT_OFFSET + lzw.currOffset
	lzw.currOffset += 1
	return retCode
}

func (lzw *LZW) InitDict(str string) {
	for _, c := range str {
		amt := uint16(rune(c))
		fmt.Printf("init: %s=%d\n", string(c), amt)
		lzw.Put(amt, string(c))
	}

	lzw.Put(START_CODE, "START")
	lzw.Put(STOP_CODE, "STOP")
}

func (lzw *LZW) Encode(w io.Writer, data string) {
	// clear dict

	// create new entry
	var newEntry string = ""
	var buffer []byte = make([]byte, 2)

	// send start code
	binary.LittleEndian.PutUint16(buffer, START_CODE)
	w.Write(buffer)

	for _, currChar := range data {
		//currChar := data[i]

		lzw.Put(uint16(rune(currChar)), string(currChar))

		testNewEntry := newEntry + string(currChar)

		_, ok := lzw.ReverseGet(testNewEntry)
		fmt.Printf("%s === %s === %v\n", newEntry, testNewEntry, ok)

		if !ok {
			// send code for new entry
			//newEntryCode := lzw.GetCodeForStr(newEntry)
			newEntryCode, _ := lzw.ReverseGet(newEntry)
			//if !ok {
			//	newEntryCode = lzw.GetCodeForStr(string(currChar))
			//}
			fmt.Printf("new code: %d\n", newEntryCode)
			binary.LittleEndian.PutUint16(buffer, newEntryCode)
			w.Write(buffer)

			// add new entry appended with char as new dict entry
			//lzw.codeDict[newEntryCode] = testNewEntry
			// this here, newEntryCode if for newEntry, not testNewEntry
			// need to get the code for testNewEntry
			testEntryNewCode := lzw.GetCodeForStr(testNewEntry)
			lzw.Put(testEntryNewCode, testNewEntry)

			// reset new entry
			newEntry = ""
		}

		// append char to new entry
		newEntry += string(currChar)
	}

	c, ok := lzw.ReverseGet(newEntry)

	if !ok {
		// send code for new entry
		newEntryCode := lzw.GetCodeForStr(newEntry)
		binary.LittleEndian.PutUint16(buffer, newEntryCode)
		w.Write(buffer)
	} else {
		// send existing code
		binary.LittleEndian.PutUint16(buffer, c)
		w.Write(buffer)
	}

	// send stop
	binary.LittleEndian.PutUint16(buffer, STOP_CODE)
	w.Write(buffer)
}

func (lzw *LZW) Decode(r io.Reader) (string, error) {
	return "", nil
}
