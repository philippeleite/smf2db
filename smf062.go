package main

import (
	"database/sql"
	"sync"
)

var pStmt062, pStmt062Vol *sql.Stmt
var once062 sync.Once
var count062 int

func smf062Init() {
	pStmt062 = createTable("smf062", `ID INT NOT NULL, SMF62TME DATETIME NOT NULL, SMF62DTE DATETIME NOT NULL,
		SMF62SID TEXT NOT NULL, SMF62JBN TEXT, SMF62RST DATETIME, SMF62RSD DATETIME, SMF62UIF TEXT,
		SMF62IND INT, SMF62CNM TEXT, SMF62CVS TEXT, SMF62DNM TEXT, SMF62VCT INT,
		SMF62MGT TEXT, SMF62STR TEXT, SMF62DAT TEXT,
		SMF62TIM DATETIME, SMF62DT DATETIME, SMF62MC1 INT, SMF62MC2 INT, SMF62MC3 INT, SMF62MC4 INT,
		SMF62DET INT, SMF62DKL TEXT, SMF62JobID TEXT, SMF62SysplexName TEXT`)
	pStmt062Vol = createTable("smf062_vol", `ID INT NOT NULL, SMF62VSR TEXT, SMF62DTY INT`)
}

// smf062Parse stores a VSAM component or cluster OPEN. After the volume entries come the SMS classes
// (always present, meaningful only when SMF62IND bit 4 is on), then the statistics section.
func smf062Parse(rec []byte) {
	once062.Do(smf062Init)

	s := section(rec)
	ind := s.num(38, 4)
	vct := s.num(136, 2)
	off := 138 + 10*vct

	sms := s.sub(off, 24)
	var mgt, str, dat interface{}
	if ind&0x08000000 != 0 { // bit 4: SMS class information is included
		mgt, str, dat = sms.text(0, 8), sms.text(8, 8), sms.text(16, 8)
	}
	off += 24
	// shortcut: assumes the 66-byte encryption block (SMS62DET, SMS62DKL) is present only for encrypted data
	// sets (bit 5); no sample record has it. Recheck against an encrypted data set's record.
	if ind&0x04000000 != 0 {
		off += 66
	}
	st := section(rec).sub(off, len(rec)-off) // statistics; newer fields are NULL in shorter records

	count062 += 1
	insert(pStmt062,
		count062,
		s.time(2),       // SMF62TME
		s.date(6),       // SMF62DTE
		s.text(10, 4),   // SMF62SID
		s.text(14, 8),   // SMF62JBN
		s.time(22),      // SMF62RST
		s.date(26),      // SMF62RSD
		s.text(30, 8),   // SMF62UIF
		ind,             // SMF62IND
		s.text(42, 44),  // SMF62CNM
		s.text(86, 6),   // SMF62CVS
		s.text(92, 44),  // SMF62DNM
		vct,             // SMF62VCT
		mgt,             // SMF62MGT
		str,             // SMF62STR
		dat,             // SMF62DAT
		st.time(0),      // SMF62TIM
		st.date(4),      // SMF62DT
		st.u8(8),        // SMF62MC1
		st.u8(9),        // SMF62MC2
		st.u8(10),       // SMF62MC3
		st.u8(11),       // SMF62MC4
		st.u16(12),      // SMF62DET
		st.text(14, 64), // SMF62DKL
		st.text(78, 8),  // SMF62JobID
		st.text(86, 8),  // SMF62SysplexName
	)

	for k := 0; k < vct; k++ {
		v := s.sub(138+10*k, 10)
		insert(pStmt062Vol, count062, v.text(0, 6), v.u32(6)) // SMF62VSR, SMF62DTY
	}
}
