package main

import (
	"database/sql"
	"sync"
)

// Type 90 (system status: SET, SETSMF, SWITCH SMF, VARY WLM, RESET and other commands). The header has no
// subtype field (offset 18 is reserved): the subtype is SMF90TID, the first field of the product section.
// smf090 holds the header and product section of every record; the subtype data section, located by the
// triplet at manual offset 28, goes to a table per subtype: SWITCH SMF (6, smf090_switch), VARY WLM (24,
// smf090_wlm) and RESET (30, smf090_reset). Subtype 24 keeps the service policy as a BLOB: it is mapped by
// IWMSVPOL (z/OS MVS Data Areas, not decoded) and may span several records (SMF9024S of SMF9024N, same
// SMF9024I). Not decoded, as the sample data has none: every other subtype.

var pStmt090, pStmt090Switch, pStmt090Wlm, pStmt090Reset *sql.Stmt
var once090 sync.Once
var count090 int

var smf090HeaderFields = fields{
	{"SMF90TME", 2, 'T', 4},
	{"SMF90DTE", 6, 'd', 4},
	{"SMF90SID", 10, 'c', 4},
}

var smf090ProductFields = fields{
	{"SMF90TID", 0, 'u', 2}, // subtype
	{"SMF90RVN", 2, 'c', 2},
	{"SMF90PNM", 4, 'c', 8},
}

var smf090SwitchFields = fields{ // subtype 6
	{"SMF90SWO", 0, 'c', 10},
	{"SMF90SWN", 10, 'c', 10},
	{"SMF90IT", 20, 'T', 4},
	{"SMF90ID", 24, 'd', 4},
	{"SMF90SOD", 28, 'c', 44},
	{"SMF90SND", 72, 'c', 44},
}

var smf090WlmFields = fields{ // subtype 24
	{"SMF9024S", 0, 'u', 4},
	{"SMF9024N", 4, 'u', 4},
	{"SMF9024I", 8, 'u', 4},
}

var smf090ResetFields = fields{ // subtype 30
	{"SMF90T30_JOBNAME", 0, 'c', 8},
	{"SMF90T30_JOBID", 8, 'c', 8},
	{"SMF90T30_ENTRY_TIME", 16, 'T', 4},
	{"SMF90T30_ENTRY_DATE", 20, 'd', 4},
	{"SMF90T30_OPERATOR", 24, 'c', 8},
	{"SMF90T30_FLAGS", 32, 'u', 1},
	{"SMF90T30_FLAGS2", 33, 'u', 1},
	{"SMF90T30_OLDSRV", 36, 'c', 8},
	{"SMF90T30_NEWSRV", 44, 'c', 8},
	{"SMF90T30_OLDPGN", 52, 'u', 2},
	{"SMF90T30_NEWPGN", 54, 'u', 2},
	{"SMF90T30_EnclaveOwner", 56, 'c', 8},
}

func smf090Init() {
	pStmt090 = createTable("smf090", "ID INT NOT NULL, "+smf090HeaderFields.columns()+", "+smf090ProductFields.columns())
	pStmt090Switch = createTable("smf090_switch", "ID INT NOT NULL, "+smf090SwitchFields.columns())
	pStmt090Wlm = createTable("smf090_wlm", "ID INT NOT NULL, "+smf090WlmFields.columns()+", SMF9024A BLOB")
	pStmt090Reset = createTable("smf090_reset", "ID INT NOT NULL, "+smf090ResetFields.columns())
}

func smf090Parse(rec []byte) {
	once090.Do(smf090Init)

	s := section(rec)
	count090 += 1
	id := count090
	poff, plen, pn := tripletAt(rec, 16)
	var prod section
	if pn > 0 {
		prod = s.sub(poff, plen)
	}
	insert(pStmt090, append(append([]interface{}{id}, smf090HeaderFields.values(s)...), smf090ProductFields.values(prod)...)...)

	doff, dlen, dn := tripletAt(rec, 24)
	if dn == 0 {
		return
	}
	d := s.sub(doff, dlen)
	switch prod.num(0, 2) {
	case 6:
		insert(pStmt090Switch, append([]interface{}{id}, smf090SwitchFields.values(d)...)...)
	case 24:
		insert(pStmt090Wlm, append(append([]interface{}{id}, smf090WlmFields.values(d)...), d.blob(12, len(d)-12))...)
	case 30:
		insert(pStmt090Reset, append([]interface{}{id}, smf090ResetFields.values(d)...)...)
	}
}
