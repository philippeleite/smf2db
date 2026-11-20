package main

import "testing"

func TestSectionBounds(t *testing.T) {
	s := section{0x00, 0x01, 0xff, 0xff, 0xd7, 0xc4}
	if got := s.u16(0); got != 1 {
		t.Errorf("u16: got %v", got)
	}
	if got := s.s32(2); got != -10300 {
		t.Errorf("s32: got %v", got)
	}
	if got := s.u32(4); got != nil { // only 2 bytes left: an older, shorter section
		t.Errorf("u32 past end: got %v, want nil", got)
	}
	if got := section(nil).text(0, 8); got != nil {
		t.Errorf("text of absent section: got %v, want nil", got)
	}
	if s.sub(-1, 2) != nil || s.sub(5, 2) != nil {
		t.Errorf("sub out of range must be nil")
	}
}

func TestSmf042Chain(t *testing.T) {
	// Three 8-byte sections at rec offsets 0, 8, 16; the "next" pointers count from the RDW (+4).
	// The last one points backwards, which must stop the walk instead of looping.
	rec := []byte{
		0, 0, 0, 12, 'A', 0, 0, 0,
		0, 0, 0, 20, 'B', 0, 0, 0,
		0, 0, 0, 4, 'C', 0, 0, 0,
	}
	var got []byte
	smf042Chain(rec, 0, 8, func(s section) { got = append(got, s[4]) })
	if string(got) != "ABC" {
		t.Errorf("smf042Chain: got %q, want %q", got, "ABC")
	}
}
