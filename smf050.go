package main

import (
	"database/sql"
	"sync"
)

// Type 50 (VTAM tuning statistics). The layouts come from z/OS Communications Server: SNA Network
// Implementation Guide (SC27-3672, data/istimp0_v3r1.pdf, "Tuning statistics file report" tables), not from
// the SMF manual. The format depends on the resource; the manual gives the fields no IBM names, so columns
// are named SMF50_ after the IST message keywords or the descriptions.
//
// smf050 holds what every format shares: header, resource (line/TRLE) name, extension length, attachment
// type and version (X'F0F3' for TCP connections, X'F0F2' for MPC, X'F0F4' for SNA controllers), plus the
// rest of the record as a BLOB. Version F0F3 records are decoded into smf050_tcp. MPC records (F0F2) are
// not decoded: they can follow either the XCF or the channel table, the sample records are longer than
// both and all their counters are zero, so the layout could not be checked. SNA controller, channel-to-
// channel and RoCE records are not in the sample data.

var pStmt050, pStmt050Tcp *sql.Stmt
var once050 sync.Once
var count050 int

var smf050Fields = fields{
	{"SMF50TME", 2, 'T', 4},
	{"SMF50DTE", 6, 'd', 4},
	{"SMF50SID", 10, 'c', 4},
	{"SMF50_NAME", 14, 'c', 8}, // controller, line or TRLE name
	{"SMF50_EXTLEN", 58, 'u', 1},
	{"SMF50_ATTACH", 59, 'u', 1}, // MPC: X'01' group, X'02' subchannel, X'03'/X'04' queue records
	{"SMF50_VERSION", 60, 'c', 2},
}

// smf050TcpFields: Table 52, TCP connections (manual offset - 4).
var smf050TcpFields = fields{
	{"SMF50_WRITECP", 26, 'u', 4}, // write channel programs (not CDLC)
	{"SMF50_READCP", 30, 'u', 4},  // read channel programs (not CDLC)
	{"SMF50_ATTN", 34, 'u', 4},
	{"SMF50_OUTMAX", 38, 'u', 4}, // largest outbound packet sent
	{"SMF50_INMAX", 42, 'u', 4},  // largest packet received
	{"SMF50_INPACKET", 78, 'u', 4},
	{"SMF50_OTPACKET", 82, 'u', 4},
	{"SMF50_INBYTE", 86, 'u', 4},
	{"SMF50_INBYTEO", 90, 'u', 4}, // overflow
	{"SMF50_OTBYTE", 94, 'u', 4},
	{"SMF50_OTBYTEO", 98, 'u', 4}, // overflow
	{"SMF50_LEGACYTYPE", 102, 'u', 1},
	{"SMF50_PCI", 106, 'u', 4}, // CLAW only, from here on
	{"SMF50_READCCW", 110, 'u', 4},
	{"SMF50_WRITECCW", 114, 'u', 4},
	{"SMF50_APPENDS", 118, 'u', 4},
	{"SMF50_RWSIO", 122, 'u', 4},
}

func smf050Init() {
	pStmt050 = createTable("smf050", "ID INT NOT NULL, "+smf050Fields.columns()+", SMF50_BODY BLOB")
	pStmt050Tcp = createTable("smf050_tcp", "ID INT NOT NULL, "+smf050TcpFields.columns())
}

func smf050Parse(rec []byte) {
	once050.Do(smf050Init)

	s := section(rec)
	count050 += 1
	v := smf050Fields.values(s) // converts the name and version in place; the body BLOB excludes them
	insert(pStmt050, append(append([]interface{}{count050}, v...), s.blob(22, len(rec)-22))...)
	if v[6] == "03" {
		insert(pStmt050Tcp, append([]interface{}{count050}, smf050TcpFields.values(s)...)...)
	}
}
