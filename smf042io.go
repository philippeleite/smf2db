package main

import "database/sql"

// Type 42 subtypes 5 (storage class and volume I/O statistics) and 6 (data set I/O statistics).
// See smf042.go.

var (
	pStmt042Sc                                                     *sql.Stmt // subtype 5
	pStmt042Vol, pStmt042VolVtoc, pStmt042VolVtix, pStmt042VolVvds *sql.Stmt // subtype 5
	pStmt042Job, pStmt042Ds                                        *sql.Stmt // subtype 6
)

// smf042IOMetrics are the consecutive 4-byte fields shared, after their own prefix, by the storage class
// response time section (S42SC, from offset 32) and the VTOC, VTOC index and VVDS component sections
// (S42VD, S42VX, S42VV, from offset 0).
var smf042IOMetrics = []string{"IOR", "IOC", "IOP", "IOD", "IOQ", "ION", "CND", "HIT", "WCN", "WHI",
	"SEQ", "RLC", "ICL", "DA0", "RDD", "RDT", "HRD", "HWR", "R1U", "C1U", "P1U", "D1U", "Q1U", "A1U", "T1U",
	"B1U", "M1U"}

// smf042DsIO are the consecutive 4-byte fields of the subtype 6 data set I/O statistics section, from offset 0.
var smf042DsIO = []string{"S42DSIOR", "S42DSIOC", "S42DSIOP", "S42DSIOD", "S42DSIOQ", "S42DSION", "S42DSCND",
	"S42DSSHTS", "S42DSWCN", "S42DSWHI", "S42DSSEQ", "S42DSRLC", "S42DSICL", "S42DSDA0", "S42DSMXR", "S42DSMXS",
	"S42DSRDD", "S42DSRDT", "S42DSHRD", "S42DSHWR", "S42DSR1U", "S42DSC1U", "S42DSP1U", "S42DSD1U", "S42DSQ1U",
	"S42DSA1U", "S42DST1U", "S42DSB1U", "S42DSM1U", "S42DSIOS", "S42DSCNT", "S42DSRRU", "S42DSRSU",
	"S42SNAvgARDelay", "S42SNARDelayCount"}

func smf042InitIO() {
	pStmt042Sc = createTable("smf042_sc", "ID INT NOT NULL, S42SCRNN TEXT, "+
		columns("S42SC", "INT", smf042IOMetrics)+", S42SCCNT INT, S42SCCRP INT")
	pStmt042Vol = createTable("smf042_vol", `ID INT NOT NULL, S42VTSER TEXT, S42VTADR INT, S42VTFL1 INT,
		S42VTUNC INT`)
	pStmt042VolVtoc = createTable("smf042_vol_vtoc", "ID INT NOT NULL, S42VTSER TEXT, "+
		columns("S42VD", "INT", smf042IOMetrics))
	pStmt042VolVtix = createTable("smf042_vol_vtix", "ID INT NOT NULL, S42VTSER TEXT, "+
		columns("S42VX", "INT", smf042IOMetrics))
	pStmt042VolVvds = createTable("smf042_vol_vvds", "ID INT NOT NULL, S42VTSER TEXT, "+
		columns("S42VV", "INT", smf042IOMetrics))

	pStmt042Job = createTable("smf042_job", `ID INT NOT NULL, S42JDJNM TEXT, S42JDRST DATETIME,
		S42JDRSD DATETIME, S42JDUID TEXT, S42JDCOD INT, S42JDVER INT, S42JDPGN INT, S42JDST1 INT,
		S42JDGMO INT, S42JDWSC TEXT, S42JDWLD TEXT, S42JDSTN TEXT`)
	pStmt042Ds = createTable("smf042_ds", `ID INT NOT NULL, S42DSN TEXT, S42DSTYP INT, S42DSCOD INT,
		S42DSFL1 INT, S42DSVOL TEXT, S42DSDEV INT, S42DSSC TEXT, S42DSBSZ INT, S42DSTRP INT, S42DSENT INT,
		S42DSCMT INT, `+columns("", "INT", smf042DsIO)+`, S42DSMXI INT, S42DSCRP INT, `+
		columns("S42AM", "INT", []string{"SRB", "SRR", "SWB", "SWR", "DRB", "DRR", "DWB", "DWR", "ZRB", "ZRR",
			"ZWB", "ZWR", "RIB", "WIB", "RBD", "WBE", "RCI", "WCI"}))
}

// smf042Chain calls f for each section of a linked list: first is the offset of the first section within
// rec (negative when there is none), and each section starts with the 4-byte offset (from the RDW) of
// the next one, 0 for the last.
func smf042Chain(rec []byte, first, length int, f func(s section)) {
	for off := first; off >= 0; {
		s := section(rec).sub(off, length)
		if s == nil {
			return
		}
		f(s)
		next := s.num(0, 4) - 4
		if next <= off { // the list only moves forward; anything else would loop forever
			return
		}
		off = next
	}
}

// Subtype 5: storage class response time sections, and the chain of volume header sections with their
// VTOC, VTOC index and VVDS statistics. Volume metrics, system I/O and background activity sections are
// not decoded: the sample records are too old to have them.
func smf042s05(rec []byte, id int) {
	smf042Each(rec, 0, func(s section) {
		insert(pStmt042Sc, append(append([]interface{}{id, s.text(2, s.num(0, 2))},
			s.u32s(32, len(smf042IOMetrics))...), s.u32(148), s.u64(152))...)
	})

	first, length, n := smf042Triplet(rec, 1)
	if n == 0 {
		return
	}
	smf042Chain(rec, first, length, func(v section) {
		ser := v.text(4, 6)
		insert(pStmt042Vol, id, ser, v.u16(10), v.u8(12), v.u32(20))
		for _, c := range []struct {
			at   int
			stmt *sql.Stmt
		}{{24, pStmt042VolVtoc}, {32, pStmt042VolVtix}, {40, pStmt042VolVvds}} {
			off, clen := v.num(c.at, 4), v.num(c.at+4, 2)
			if off == 0 || clen == 0 {
				continue
			}
			comp := section(rec).sub(off-4, clen)
			insert(c.stmt, append([]interface{}{id, ser}, comp.u32s(0, len(smf042IOMetrics))...)...)
		}
	})
}

// Subtype 6: job header, and the chain of data set header sections with their I/O and access method
// statistics in the same row. The synchronous I/O sections are not decoded.
func smf042s06(rec []byte, id int) {
	j := smf042First(rec, 0)
	if j == nil {
		return
	}
	insert(pStmt042Job, id, j.text(0, 8), j.time(8), j.date(12), j.text(16, 8), j.u8(30), j.u8(31), j.u16(32),
		j.u8(38), j.s32(40), j.text(44, 8), j.text(52, 8), j.text(64, 8))

	ioLen, amLen := j.num(34, 2), j.num(36, 2)
	smf042Chain(rec, j.num(24, 4)-4, j.num(28, 2), func(d section) {
		var io, am section // stay nil (all NULL) when the data set has no such section
		if off := d.num(52, 4); off > 0 {
			io = section(rec).sub(off-4, ioLen)
		}
		if off := d.num(56, 4); off > 0 {
			am = section(rec).sub(off-4, amLen)
		}
		row := []interface{}{id, d.text(4, 44), d.u8(48), d.u8(49), d.u8(50), d.text(60, 6), d.u16(66),
			d.text(68, 8), d.u32(76), d.u16(80), d.u16(82), d.u8(84)}
		row = append(row, io.u32s(0, len(smf042DsIO))...)
		row = append(row, io.u16(140), io.u64(144))
		row = append(row, am.u32s(0, 12)...)
		row = append(row, am.u64(48), am.u64(56), am.u64(64), am.u64(72), am.u32(80), am.u32(84))
		insert(pStmt042Ds, row...)
	})
}
