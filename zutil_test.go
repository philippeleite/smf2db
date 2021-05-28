package main

import (
	"bytes"
	"testing"
)

func TestSmfTime(t *testing.T) {
	var ttime = []byte{0x00, 0x5d, 0x77, 0x40}
	stime := smftime(ttime)
	if stime != "0000-00-00 17:00:53.760" {
		t.Errorf("Time error")
	}
}

func TestSmfDate(t *testing.T) {
	var tdate = []byte{0x01, 0x16, 0x06, 0x0f}
	sdate := smfdate(tdate)
	if sdate != "2016-02-29 00:00:00.000" {
		t.Errorf("Date error")
	}
}

func TestFileDate(t *testing.T) {
	tdate := []byte{0x74, 0x00, 0x60}
	sdate := filedate(tdate)
	if sdate != "2016-02-29 00:00:00.000" {
		t.Errorf("Date error")
	}
	tdate = []byte{0x00, 0x00, 0x00}
	sdate = filedate(tdate)
	if sdate != "0000-00-00 00:00:00.000" {
		t.Errorf("Date error")
	}
}

func TestAsc2bcd(t *testing.T) {
	var tstr = []byte("ABCD")
	asc2bcd(tstr)
	if !bytes.Equal(tstr, []byte{0xc1, 0xc2, 0xc3, 0xc4}) {
		t.Errorf("asc2bcd error")
	}
}

func TestBcd2asc(t *testing.T) {
	var tstr = []byte{0xc1, 0xc2, 0xc3, 0xc4}
	bcd2asc(tstr)
	if !bytes.Equal(tstr, []byte("ABCD")) {
		t.Errorf("bcd2asc error")
	}
}
