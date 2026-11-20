package main

import (
	"database/sql"
	"sync"
)

// Type 92 (z/OS UNIX file system activity). smf092 holds the header, the subsystem section and the
// identification section of every record (subtype 17 has no identification section: those columns are
// NULL). Data sections go to one table per layout, with ID and SEQ (see rmfSection): file open (10),
// file close (11, and 16 for sockets, character special files, pipes and FIFOs, same layout), file delete
// or rename (14) and file access count (17). Not decoded, as the sample data has none: subtypes 1-9, 12,
// 13, 15 and the zFS subtypes 50-59, which also use another identification section layout.
//
// Clock quirk seen in the samples: open times (SMF92OTO, SMF92CTO) are the raw TOD clock, which includes
// leap seconds (26 s ahead of the record time), while the close time SMF92CTC is already UTC. The real
// time a file was open is SMF92CTC - SMF92CTO + the leap-second offset. Values are stored as recorded.

var pStmt092 *sql.Stmt
var once092 sync.Once
var count092 int

var smf092HeaderFields = fields{
	{"SMF92TME", 2, 'T', 4},
	{"SMF92DTE", 6, 'd', 4},
	{"SMF92SID", 10, 'c', 4},
	{"SMF92WID", 14, 'c', 4},
	{"SMF92STP", 18, 'u', 2},
}

// smf092SubsysFields: subsystem section.
var smf092SubsysFields = fields{
	{"SMF92TYP", 0, 'u', 2},
	{"SMF92RVN", 2, 'c', 2},
	{"SMF92PNM", 4, 'c', 8},
	{"SMF92OSL", 12, 'c', 8},
}

// smf092IdentFields: identification section (process layout, used by every subtype below 50).
var smf092IdentFields = fields{
	{"SMF92JBN", 0, 'c', 8},
	{"SMF92RST", 8, 'T', 4},
	{"SMF92RSD", 12, 'd', 4},
	{"SMF92STM", 16, 'c', 8},
	{"SMF92RGD", 24, 'c', 8},
	{"SMF92RUD", 32, 'c', 8},
	{"SMF92UID", 40, 'u', 4},
	{"SMF92GID", 44, 'u', 4},
	{"SMF92PID", 48, 'u', 4},
	{"SMF92PGD", 52, 'u', 4},
	{"SMF92SSD", 56, 'u', 4},
	{"SMF92API", 60, 'u', 4},
	{"SMF92APG", 64, 'u', 4},
	{"SMF92ASG", 68, 'u', 4},
}

var smf092OpenFields = fields{ // subtype 10
	{"SMF92OTO", 0, 't', 8},
	{"SMF92OTY", 8, 'u', 1},
	{"SMF92OFG", 9, 'u', 1},
	{"SMF92OTK", 12, 'u', 4},
	{"SMF92OIN", 16, 'u', 4},
	{"SMF92ODN", 20, 'u', 4},
}

var smf092CloseFields = fields{ // subtypes 11 and 16
	{"SMF92CTO", 0, 't', 8},
	{"SMF92CTC", 8, 't', 8},
	{"SMF92CTY", 16, 'u', 1},
	{"SMF92CFG", 17, 'u', 1},
	{"SMF92CTK", 20, 'u', 4},
	{"SMF92CIN", 24, 'u', 4},
	{"SMF92CDN", 28, 'u', 4},
	{"SMF92CSR", 32, 'u', 4},
	{"SMF92CSW", 36, 'u', 4},
	{"SMF92CDI", 40, 'u', 4},
	{"SMF92CIR", 44, 'u', 4},
	{"SMF92CIW", 48, 'u', 4},
	{"SMF92CBR", 52, 'u', 8},
	{"SMF92CBW", 60, 'u', 8},
	{"SMF92CPN", 68, 'c', 64},
}

var smf092DeleteFields = fields{ // subtype 14
	{"SMF92DFT", 0, 't', 8},
	{"SMF92DTY", 8, 'u', 1},
	{"SMF92DFLG", 9, 'u', 1},
	{"SMF92DIN", 12, 'u', 4},
	{"SMF92DINP", 16, 'u', 4},
	{"SMF92DDN", 20, 'u', 4},
	{"SMF92DFS", 24, 'c', 44},
	{"SMF92DNL", 68, 'u', 4},
	{"SMF92DFN", 72, 'c', 64},
	{"SMF92DNLR", 136, 'u', 4},
	{"SMF92DFNR", 140, 'c', 64}, // printed as SM92DFNR in the manual
}

var smf092CountFields = fields{ // subtype 17
	{"SMF92FAWT", 0, 't', 8},
	{"SMF92FAFT", 8, 'u', 1},
	{"SMF92FAFG", 9, 'u', 1},
	{"SMF92FAIN", 12, 'u', 4},
	{"SMF92ADN", 16, 'u', 4},
	{"SMF92FATI", 20, 'u', 4},
	{"SMF92FAPN", 24, 'c', 64},
}

var smf092Sections = []*rmfSection{
	{subtype: 10, at: 40, table: "smf092_open", f: smf092OpenFields},
	{subtype: 11, at: 40, table: "smf092_close", f: smf092CloseFields},
	{subtype: 16, at: 40, table: "smf092_close", f: smf092CloseFields},
	{subtype: 14, at: 40, table: "smf092_delete", f: smf092DeleteFields},
	{subtype: 17, at: 40, table: "smf092_count", f: smf092CountFields},
}

func smf092Init() {
	pStmt092 = createTable("smf092", "ID INT NOT NULL, "+smf092HeaderFields.columns()+", "+
		smf092SubsysFields.columns()+", "+smf092IdentFields.columns())
	rmfCreateSections(smf092Sections)
}

func smf092Parse(rec []byte) {
	once092.Do(smf092Init)

	s := section(rec)
	count092 += 1
	soff, slen, _ := tripletAt(rec, 24)
	ioff, ilen, in := tripletAt(rec, 32)
	var ident section // absent in subtype 17
	if in > 0 {
		ident = s.sub(ioff, ilen)
	}
	row := append([]interface{}{count092}, smf092HeaderFields.values(s)...)
	row = append(row, smf092SubsysFields.values(s.sub(soff, slen))...)
	insert(pStmt092, append(row, smf092IdentFields.values(ident)...)...)
	rmfInsertSections(smf092Sections, rec, count092)
}
