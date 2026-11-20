package main

import "testing"

func TestHfp(t *testing.T) {
	for _, c := range []struct {
		b    []byte
		want float64
	}{
		{[]byte{0x41, 0x10, 0, 0, 0, 0, 0, 0}, 1},
		{[]byte{0x42, 0x64, 0, 0, 0, 0, 0, 0}, 100},
		{[]byte{0xc1, 0x18, 0, 0, 0, 0, 0, 0}, -1.5},
		{[]byte{0x40, 0x80, 0, 0, 0, 0, 0, 0}, 0.5},
		{make([]byte, 8), 0},
	} {
		if got := hfp(c.b); got != c.want {
			t.Errorf("hfp(% X) = %v, want %v", c.b, got, c.want)
		}
	}
}

func TestFieldsValues(t *testing.T) {
	f := fields{
		{"S", 0, 's', 4},
		{"H", 4, 'h', 4},
		{"I", 8, 'i', 4},
		{"P", 12, 'p', 4},
		{"X", 16, 'x', 3},
		{"PAST", 19, 'u', 4},
	}
	if got := (fields{{"T", 0, 'H', 4}}).values(section{0x11, 0x00, 0x25, 0x6f})[0]; got != "0000-00-00 11:00:25.600" {
		t.Errorf("H: got %v", got)
	}
	if got := (fields{{"T", 0, 'T', 4}}).values(section{0x00, 0x5d, 0x77, 0x40})[0]; got != "0000-00-00 17:00:53.760" {
		t.Errorf("T: got %v", got)
	}
	if got := (fields{{"M", 0, 'M', 8}}).values(section{0xd2, 0x72, 0x28, 0xa6, 0x3b, 0x4c, 0x00, 0x00})[0]; got != 0xd27228a63b4c0000>>12 {
		t.Errorf("M: got %v", got)
	}
	wide := fields{{"W", 0, 'w', 16}}.values(section{0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2})
	if wide[0] != 18446744073709551618.0 { // 2^64 + 2
		t.Errorf("wide: got %v", wide[0])
	}
	s := section{
		0xff, 0xff, 0xd7, 0xc4, // -10300
		0x01, 0x33, 0x00, 0x0f, // 13:30:00
		0x29, 0x59, 0x98, 0x6f, // 29:59.986
		0x00, 0x00, 0x10, 0x0f, // 100
		0x01, 0x23, 0x4a, // "01234A"
		0x00, // only 1 byte left for PAST
	}
	v := f.values(s)
	want := []interface{}{-10300, "0000-00-00 13:30:00.000", 1799.986, 100, "01234A", nil}
	for i := range want {
		if v[i] != want[i] {
			t.Errorf("%s: got %v (%T), want %v", f[i].name, v[i], v[i], want[i])
		}
	}
	if got := f.columns(); got != "S INT, H DATETIME, I REAL, P INT, X TEXT, PAST INT" {
		t.Errorf("columns: %q", got)
	}
}
