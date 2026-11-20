package main

import "testing"

func TestSmfParseSkipsMalformedRecord(t *testing.T) {
	called := false
	parseTable[255] = func(rec []byte) {
		called = true
		_ = rec[100] // simulate a parser reading past the end of a truncated record
	}
	defer delete(parseTable, 255)

	smfParse([]byte{0x00, 0xff, 0x01}) // must not panic
	smfParse([]byte{0x00})             // too short to hold a type, must not panic
	if !called {
		t.Errorf("parser was not called")
	}
}
