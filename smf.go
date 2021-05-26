package main

import "encoding/binary"

var parseTable = map[int]func([]byte){
	4: smf004Parse,
}

func b3toInt(s []byte) int {
	var t = make([]byte, 4)
	copy(t[1:], s)
	return int(binary.BigEndian.Uint32(t))
}

func smfInit() {
	smf004Init()
}

func smfParse(rec []byte) {
	if f, ok := parseTable[int(rec[1])]; ok {
		f(rec)
	}
}
