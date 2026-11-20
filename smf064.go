package main

import (
	"database/sql"
	"sync"
)

var pStmt064, pStmt064Ext *sql.Stmt
var once064 sync.Once
var count064 int

// smf064Counts are the consecutive 4-byte fields of the statistics section from offset 4: 11 totals at OPEN,
// 11 changes since OPEN (signed: they go negative, e.g. fewer free CIs), then block, CI and record sizes.
var smf064Counts = []string{"NIL", "NEX", "NLR", "NDE", "NIN", "NUP", "NRE", "NFS", "NCS", "NAS", "NEP",
	"DIL", "DEX", "DLR", "DDE", "DIN", "DUP", "DRE", "DFS", "DCS", "DAS", "DEP", "DBS", "DCI", "DLS"}

func smf064Init() {
	pStmt064 = createTable("smf064", `ID INT NOT NULL, SMF64TME DATETIME NOT NULL, SMF64DTE DATETIME NOT NULL,
		SMF64SID TEXT NOT NULL, SMF64JBN TEXT, SMF64RST DATETIME, SMF64RSD DATETIME, SMF64UIF TEXT,
		SMF64RIN INT, SMF64DTY INT, SMF64CNM TEXT, SMF64DNM TEXT, SMF64NTR INT, SMF64CHR INT, `+
		columns("SMF64", "INT", smf064Counts)+`, SMF64DKL INT, SMF64DDN TEXT, SMF64STR INT, SMF64BNO INT,
		SMF64BSP INT, SMF64BFD INT, SMF64BFI INT, SMF64CLN TEXT, SMF64PLH INT,
		SMF64MC1 INT, SMF64MC2 INT, SMF64MC3 INT, SMF64MC4 INT, SMF64SMB INT, SMF64RSC INT,
		SMF64SIO INT, SMF64HIT INT, SMF64WTS INT, SMF64MIS INT, SMF64IOS INT,
		SMF64SDS INT, SMF64CDS INT, SMF64CSS INT, SMF64CCS INT, SMF64DTK BLOB, SMF64TRK INT,
		SMF64BMH INT, SMF64CFH INT, SMF64RIO INT, SMF64TIM DATETIME, SMF64DT DATETIME,
		SMF64FG1 INT, SMF64FD1 INT, SMF64FD2 INT, SMF64DAU INT, SMF64RLM INT, SMF64NTA INT,
		SMF64JobID TEXT, SMF64SysplexName TEXT`)
	pStmt064Ext = createTable("smf064_ext", `ID INT NOT NULL, SMF64FCC INT, SMF64TCC INT, SMF64VSN TEXT,
		SMF64CUU INT, SMF64IND TEXT, SMF64UTY INT`)
}

// smf064Parse stores a VSAM component or cluster status record (CLOSE, EOV, out of space). The extent
// entries follow SMF64ESL; the statistics section follows them and carries its own length. Records with
// SMF64RIN X'01' have no extents but may end with unused bytes, which are ignored.
func smf064Parse(rec []byte) {
	once064.Do(smf064Init)

	s := section(rec)
	esl := s.num(134, 2)
	statOff := 136 + esl
	st := s.sub(statOff, s.num(statOff, 4))
	if st == nil { // length beyond the record: keep what is there, missing fields become NULL
		st = s.sub(statOff, len(rec)-statOff)
	}

	row := []interface{}{count064 + 1,
		s.time(2),      // SMF64TME
		s.date(6),      // SMF64DTE
		s.text(10, 4),  // SMF64SID
		s.text(14, 8),  // SMF64JBN
		s.time(22),     // SMF64RST
		s.date(26),     // SMF64RSD
		s.text(30, 8),  // SMF64UIF
		s.u8(38),       // SMF64RIN
		s.u8(39),       // SMF64DTY
		s.text(40, 44), // SMF64CNM
		s.text(84, 44), // SMF64DNM
		s.u16(128),     // SMF64NTR
		s.u32(130),     // SMF64CHR (or SMF64CHC); documented as EBCDIC, but it is a binary RBA/CI
	}
	for i := range smf064Counts {
		if i >= 11 && i < 22 {
			row = append(row, st.s32(4+4*i))
		} else {
			row = append(row, st.u32(4+4*i))
		}
	}
	row = append(row,
		st.u16(104),      // SMF64DKL
		st.text(106, 8),  // SMF64DDN
		st.u8(114),       // SMF64STR
		st.u8(115),       // SMF64BNO
		st.u32(116),      // SMF64BSP
		st.u16(120),      // SMF64BFD
		st.u16(122),      // SMF64BFI
		st.text(124, 44), // SMF64CLN
		st.u16(168),      // SMF64PLH
		st.u8(170),       // SMF64MC1
		st.u8(171),       // SMF64MC2
		st.u8(172),       // SMF64MC3
		st.u8(173),       // SMF64MC4
		st.u8(174),       // SMF64SMB
		st.u8(175),       // SMF64RSC
		st.u32(176),      // SMF64SIO
		st.u32(180),      // SMF64HIT
		st.u32(184),      // SMF64WTS
		st.u32(188),      // SMF64MIS
		st.u32(192),      // SMF64IOS
		st.u64(196),      // SMF64SDS
		st.u64(204),      // SMF64CDS
		st.u64(212),      // SMF64CSS
		st.u64(220),      // SMF64CCS
		st.blob(228, 36), // SMF64DTK
		st.u32(264),      // SMF64TRK
		st.u32(268),      // SMF64BMH
		st.u32(272),      // SMF64CFH
		st.u32(276),      // SMF64RIO
		st.time(280),     // SMF64TIM
		st.date(284),     // SMF64DT
		st.u8(288),       // SMF64FG1
		st.u8(289),       // SMF64FD1
		st.u8(290),       // SMF64FD2
		st.s32(292),      // SMF64DAU (a change, signed)
		st.u32(296),      // SMF64RLM
		st.u32(300),      // SMF64NTA
		st.text(304, 8),  // SMF64JobID
		st.text(312, 8),  // SMF64SysplexName
	)
	count064 += 1
	insert(pStmt064, row...)

	for off := 136; off+26 <= statOff; off += 26 {
		e := s.sub(off, 26)
		insert(pStmt064Ext, count064, e.u32(0), e.u32(4), e.text(8, 6), e.u16(14), e.text(16, 2), e.u32(18))
	}
}
