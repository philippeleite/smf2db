package main

import (
	"database/sql"
	"sync"
)

// Type 73 (RMF channel path activity). smf073 holds the header and RMF product section, smf073_ctl the
// channel path control section, smf073_chp one row per channel path data section and smf073_ext the
// extended channel path data sections. Child rows have ID and SEQ (1-based position in the record).
//
// The 48-byte channel measurement data (SMF73CCM, offset 24 of the channel path section) has a different
// layout for each measurement group (SMF73CMG 1, 2 or 3): smf073_chp has the columns of all three, and only
// the group of the channel path is filled. The fields documented as "float" are IBM hexadecimal floating
// point: in the samples, a HiperSockets data unit size decodes to 4096 as HFP (576 as IEEE).
// A record has a section for every possible channel path ID; SMF73FG3 bit 6 (X'02') marks real entries.
// The field tables were generated from the manual with tools/genfields.py.

var pStmt073, pStmt073Ctl, pStmt073Chp, pStmt073Ext *sql.Stmt
var once073 sync.Once
var count073 int

func smf073Init() {
	pStmt073 = createTable("smf073", rmfColumns("SMF73"))
	pStmt073Ctl = createTable("smf073_ctl", "ID INT NOT NULL, "+smf073CtlFields.columns())
	pStmt073Chp = createTable("smf073_chp", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf073ChpFields.columns()+", "+
		smf073G1Fields.columns()+", "+smf073G2Fields.columns()+", "+smf073G3Fields.columns())
	pStmt073Ext = createTable("smf073_ext", "ID INT NOT NULL, SEQ INT NOT NULL, SMF73ECP INT, "+smf073EdtFields.columns())
}

func smf073Parse(rec []byte) {
	once073.Do(smf073Init)

	count073 += 1
	id := count073
	insert(pStmt073, rmfRow(rec, id, "SMF73")...)

	s := section(rec)
	if off, length, n := tripletAt(rec, 32); n > 0 {
		insert(pStmt073Ctl, append([]interface{}{id}, smf073CtlFields.values(s.sub(off, length))...)...)
	}
	off, length, n := tripletAt(rec, 40)
	for k := 0; k < n; k++ {
		c := s.sub(off+k*length, length)
		ccm := c.sub(24, 48)
		var g1, g2, g3 section // only the channel's measurement group is decoded; the others stay NULL
		switch c.num(22, 1) {  // SMF73CMG
		case 1:
			g1 = ccm
		case 2:
			g2 = ccm
		case 3:
			g3 = ccm
		}
		row := append([]interface{}{id, k + 1}, smf073ChpFields.values(c)...)
		row = append(row, smf073G1Fields.values(g1)...)
		row = append(row, smf073G2Fields.values(g2)...)
		insert(pStmt073Chp, append(row, smf073G3Fields.values(g3)...)...)
	}
	off, length, n = tripletAt(rec, 48)
	for k := 0; k < n; k++ {
		e := s.sub(off+k*length, length)
		insert(pStmt073Ext, append([]interface{}{id, k + 1, e.u8(0)}, smf073EdtFields.values(e.sub(4, 32))...)...)
	}
}

var smf073CtlFields = fields{
	{"SMF73SMP", 0, 'u', 4},
	{"SMF73CFL", 4, 'u', 1},
	{"SMF73SFL", 5, 'u', 1},
	{"SMF73TNM", 8, 'c', 44},
	{"SMF73TSF", 52, 'c', 2},
	{"SMF73TOK", 56, 'c', 16},
	// skipped SMF73TDT at 56: redefines an earlier field
	// skipped SMF73TTM at 64: redefines an earlier field
	{"SMF73CRC", 72, 'u', 4},
	{"SMF73CSC", 76, 'u', 4},
	{"SMF73TDY", 80, 'c', 10},
	{"SMF73CMI", 90, 'u', 1},
	{"SMF73CSS", 91, 'c', 1},
}

var smf073ChpFields = fields{
	{"SMF73PID", 0, 'u', 1},
	{"SMF73FG2", 1, 'u', 1},
	{"SMF73FG3", 2, 'u', 1},
	{"SMF73FG4", 3, 'u', 1},
	{"SMF73BSY", 4, 'u', 4},
	{"SMF73PBY", 8, 'u', 4},
	{"SMF73PTI", 12, 'u', 4},
	{"SMF73CPD", 16, 'u', 1},
	{"SMF73ACR", 17, 'c', 5},
	{"SMF73CMG", 22, 'u', 1},
	{"SMF73FG5", 23, 'u', 1},
	{"SMF73CPP", 72, 'u', 1},
	{"SMF73GEN", 73, 'u', 1},
	{"SMF73EIX", 74, 'u', 2},
	{"SMF73SPD", 76, 'u', 2},
	{"SMF73MSC", 78, 'u', 1},
	{"SMF73NT1", 80, 'c', 16},
	{"SMF73NT2", 96, 'c', 16},
}

var smf073G1Fields = fields{
	{"SMF73TUT", 0, 'f', 4},
	{"SMF73PUT", 4, 'f', 4},
}
var smf073G2Fields = fields{
	{"SMF73MBC", 0, 'f', 4},
	{"SMF73MCU", 4, 'f', 4},
	{"SMF73MWU", 8, 'f', 4},
	{"SMF73MRU", 12, 'f', 4},
	{"SMF73US", 16, 'f', 4},
	{"SMF73TBC", 20, 'f', 4},
	{"SMF73TUC", 24, 'f', 4},
	{"SMF73PUC", 28, 'f', 4},
	{"SMF73TWU", 32, 'f', 4},
	{"SMF73PWU", 36, 'f', 4},
	{"SMF73TRU", 40, 'f', 4},
	{"SMF73PRU", 44, 'f', 4},
}
var smf073G3Fields = fields{
	{"SMF73PDU", 0, 'f', 4},
	{"SMF73TDU", 4, 'f', 4},
	{"SMF73PUM", 8, 'f', 4},
	{"SMF73TUM", 12, 'f', 4},
	{"SMF73PMS", 20, 'f', 4},
	{"SMF73TMS", 24, 'f', 4},
	{"SMF73PUS", 28, 'f', 4},
	{"SMF73PUB", 32, 'f', 4},
	{"SMF73TUB", 36, 'f', 4},
	{"SMF73PDS", 40, 'f', 4},
	{"SMF73TDS", 44, 'f', 4},
}
var smf073EdtFields = fields{
	{"SMF73EOC", 0, 'f', 4},
	{"SMF73EOD", 4, 'f', 4},
	{"SMF73EOS", 8, 'f', 8},
	{"SMF73ETC", 16, 'f', 4},
	{"SMF73ETD", 20, 'f', 4},
	{"SMF73ETS", 24, 'f', 8},
}
