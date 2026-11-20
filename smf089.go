package main

import (
	"database/sql"
	"sync"
)

// Type 89 (usage data, used for sub-capacity pricing). smf089 holds the header, the record product section
// and the system ID section of every record; subtype 1 adds the usage data (smf089_usage) and product
// intersection (smf089_isect) sections, subtype 2 the state data sections (smf089_state), with ID and SEQ.
// The self-defining section is not a list of 8-byte triplets (SMF89UDR sits after the usage triplet), so
// each triplet is read at its own offset, and only if it lies within SMF89SDL: older records end the
// section before the tenant resource group triplets. Not decoded, as the sample data has none: the tenant
// resource group and tenant intersection sections.
// The field tables were generated from the manual with tools/genfields.py; times, dates, STCK offsets and
// the "long floating point" CPU times were fixed by hand.

var pStmt089, pStmt089Usage, pStmt089Isect, pStmt089State *sql.Stmt
var once089 sync.Once
var count089 int

var smf089HeaderFields = fields{
	{"SMF89TME", 2, 'T', 4},
	{"SMF89DTE", 6, 'd', 4},
	{"SMF89SID", 10, 'c', 4},
	{"SMF89SSI", 14, 'c', 4},
	{"SMF89STY", 18, 'u', 2},
}

func smf089Init() {
	pStmt089 = createTable("smf089", "ID INT NOT NULL, "+smf089HeaderFields.columns()+", "+
		smf089ProductFields.columns()+", "+smf089SysFields.columns())
	pStmt089Usage = createTable("smf089_usage", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf089UsageFields.columns())
	pStmt089Isect = createTable("smf089_isect", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf089IsectFields.columns())
	pStmt089State = createTable("smf089_state", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf089StateFields.columns())
}

func smf089Parse(rec []byte) {
	once089.Do(smf089Init)

	s := section(rec)
	count089 += 1
	id := count089
	sdl := s.num(20, 4) // SMF89SDL: length of the self-defining section, which starts at rec[24]
	// triplet returns the sections located by the triplet at rec[at:], if it lies within SMF89SDL.
	triplet := func(at int) (off, length, n int) {
		if at+8 > 24+sdl {
			return 0, 0, 0
		}
		return tripletAt(rec, at)
	}
	first := func(at int) section {
		off, length, n := triplet(at)
		if n == 0 {
			return nil
		}
		return s.sub(off, length)
	}
	row := append([]interface{}{id}, smf089HeaderFields.values(s)...)
	row = append(row, smf089ProductFields.values(first(24))...)
	insert(pStmt089, append(row, smf089SysFields.values(first(32))...)...)

	each := func(at int, stmt *sql.Stmt, f fields) {
		off, length, n := triplet(at)
		for k := 0; k < n; k++ {
			insert(stmt, append([]interface{}{id, k + 1}, f.values(s.sub(off+k*length, length))...)...)
		}
	}
	switch s.num(18, 2) {
	case 1:
		each(40, pStmt089Usage, smf089UsageFields) // SMF89UDO/UDL/UDN
		each(52, pStmt089Isect, smf089IsectFields) // SMF89CNO/CNL/CNN
	case 2:
		each(40, pStmt089State, smf089StateFields) // same triplet, state data sections
	}
}

var smf089ProductFields = fields{
	{"SMF89PNM", 0, 'c', 8},
	{"SMF89RVN", 8, 'u', 4},
	{"SMF89OSL", 12, 'c', 8},
	{"SMF89IST", 20, 'T', 4},
	{"SMF89ISD", 24, 'd', 4},
	{"SMF89IET", 28, 'T', 4},
	{"SMF89IED", 32, 'd', 4},
	{"SMF89PFL", 36, 'u', 1},
	{"SMF89HOF", 40, 's', 8}, // STCK offset,
	{"SMF89DTO", 48, 's', 8}, // local time offset (CVTLDTO), STCK units,
	{"SMF89_CoreMode_CP", 56, 'u', 2},
	{"SMF89_CoreMode_zAAP", 58, 'u', 2},
	{"SMF89_CoreMode_zIIP", 60, 'u', 2},
}

var smf089SysFields = fields{
	{"SMF89SYN", 0, 'c', 8},
	{"SMF89UST", 8, 'T', 4},
	{"SMF89USD", 12, 'd', 4},
	{"SMF89UET", 16, 'T', 4},
	{"SMF89UED", 20, 'd', 4},
	{"SMF89CMN", 32, 'x', 2},
	{"SMF89CVN", 34, 'u', 1},
	{"SMF89LPI", 35, 'u', 1},
	{"SMF89SER", 36, 'x', 3},
	{"SMF89LP3", 39, 'u', 1},
	{"SMF89RPP", 40, 'u', 4},
	{"SMF89SPN", 44, 'c', 8},
	{"SMF89CPT", 52, 'c', 6},
	{"SMF89CPM", 58, 'c', 3},
	{"SMF89CPS", 61, 'c', 12},
	{"SMF89SIF", 73, 'u', 1},
	{"SMF89MNF", 76, 'c', 16},
	{"SMF89TID", 92, 'c', 4},
	{"SMF89MDL", 96, 'c', 16},
	{"SMF89SQC", 112, 'c', 16},
	{"SMF89POM", 128, 'c', 4},
	{"SMF89CPC", 132, 'u', 4},
	{"SMF89CCC", 136, 'u', 2},
	{"SMF89SCC", 138, 'u', 2},
	{"SMF89MAF", 140, 'x', 30}, // 15 two-byte adjustment factors,
	{"SMF89LPN", 170, 'c', 8},
	{"SMF89_Capacity_Change_Cnt", 178, 'u', 2},
	{"SMF89_RCTPCPUA_Actual", 180, 'u', 4},
	{"SMF89_RCTPCPUA_Nominal", 184, 'u', 4},
	{"SMF89_RCTPCPUA_scaling_factor", 188, 'u', 4},
	{"SMF89_Capacity_Adjustment_Ind", 192, 'u', 1},
	{"SMF89_Capacity_Change_Rsn", 193, 'u', 1},
	{"SMF89_Capacity_Flags", 194, 'u', 1},
	{"SMF89ZNF", 196, 'u', 4},
	{"SMF89SNF", 200, 'u', 4},
	{"SMF89SEQ", 204, 'u', 2},
	{"SMF89SolutionID", 206, 'c', 64},
}

var smf089UsageFields = fields{
	{"SMF89UPO", 0, 'c', 16},
	{"SMF89UPN", 16, 'c', 16},
	{"SMF89UPV", 32, 'c', 8},
	{"SMF89UPQ", 40, 'c', 8},
	{"SMF89UPI", 48, 'c', 8},
	{"SMF89UCT", 56, 'f', 8}, // TCB time, hundredths of a second
	{"SMF89USR", 64, 'f', 8}, // SRB time, hundredths of a second
	{"SMF89UFG", 72, 'u', 1},
	{"SMF89UF2", 73, 'u', 1},
	{"SMF89_BoostInfo", 74, 'u', 1},
	{"SMF89URT", 75, 'u', 1},
	{"SMF89URD", 76, 'x', 8}, // product data, format given by SMF89URT
	{"SMF89UZT", 84, 'f', 8}, // offload engine time, hundredths of a second
	{"SMF89CountAsTrad", 92, 'u', 4},
	{"SMF89CountAsTrg", 96, 'u', 4},
}

var smf089IsectFields = fields{
	{"SMF89CPO", 0, 'c', 16},
	{"SMF89CPN", 16, 'c', 16},
	{"SMF89CPV", 32, 'c', 8},
	{"SMF89CPQ", 40, 'c', 8},
	{"SMF89CPI", 48, 'c', 8},
	{"SMF89IPO", 56, 'c', 16},
	{"SMF89IPN", 72, 'c', 16},
	{"SMF89IPV", 88, 'c', 8},
	{"SMF89IPQ", 96, 'c', 8},
	{"SMF89IPI", 104, 'c', 8},
	{"SMF89CFG", 112, 'u', 1},
	{"SMF89CCT", 120, 'f', 8}, // intersect TCB time, hundredths of a second
	{"SMF89CZT", 128, 'f', 8}, // intersect offload engine time, hundredths of a second
}

var smf089StateFields = fields{
	{"SMF89T2ProdOwner", 0, 'c', 16},
	{"SMF89T2ProdName", 16, 'c', 16},
	{"SMF89T2FeatureName", 32, 'c', 16},
	{"SMF89T2ProdVers", 48, 'c', 2},
	{"SMF89T2ProdRel", 50, 'c', 2},
	{"SMF89T2ProdMod", 52, 'c', 2},
	{"SMF89T2ProdID", 54, 'c', 8},
	{"SMF89T2Flags", 62, 'u', 1},
	{"SMF89T2_BoostInfo", 63, 'u', 1},
	{"SMF89T2NumInstances", 64, 'u', 4},
	{"SMF89T2NumCurrentInstances", 68, 'u', 4},
	{"SMF89T2NumInstancesinTrgs", 72, 'u', 4},
	{"SMF89T2NumNewRegistrations", 76, 'u', 4},
	{"SMF89T2NumNewTrgRegistrations", 80, 'u', 4},
}
