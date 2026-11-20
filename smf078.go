package main

import (
	"database/sql"
	"sync"
)

// Type 78 (RMF virtual storage and I/O queuing). smf078 holds the header and RMF product section of every
// record. Subtype 2 adds the common storage section (smf078_common); subtype 3 the I/O queuing global
// (smf078_ioqglob), configuration control (smf078_cfgctl), queuing data (smf078_ioq) and HyperPAV/SuperPAV
// (smf078_hpav) sections, all with ID and SEQ (see rmfSection). Two of them carry their own triplet, whose
// offset is relative to the start of the section: the IOP initiative queue sections of the global section
// (smf078_iop) and the configuration data sections of each configuration control section (smf078_cfg).
// Those rows have ID, SEQ (of the enclosing section) and SUBSEQ (position within it).
// Not decoded, as the sample data has none: the subtype 2 private area and private area subpool sections.
// The field tables were generated from the manual with tools/genfields.py.

var pStmt078, pStmt078Iop, pStmt078Cfg *sql.Stmt
var once078 sync.Once
var count078 int

var smf078Sections = []*rmfSection{
	{subtype: 2, at: 32, table: "smf078_common", f: smf078CommonFields},     // virtual storage common storage
	{subtype: 3, at: 32, table: "smf078_cfgctl", f: smf078CfgCtlFields},     // I/O queuing configuration control
	{subtype: 3, at: 40, table: "smf078_ioq", f: smf078IoqFields},           // I/O queuing data
	{subtype: 3, at: 48, table: "smf078_ioqglob", f: smf078IoqGlobalFields}, // I/O queuing global
	{subtype: 3, at: 56, table: "smf078_hpav", f: smf078HpavFields},         // HyperPAV/SuperPAV
}

func smf078Init() {
	pStmt078 = createTable("smf078", rmfColumns("SMF78"))
	rmfCreateSections(smf078Sections)
	pStmt078Iop = createTable("smf078_iop", "ID INT NOT NULL, SEQ INT NOT NULL, SUBSEQ INT NOT NULL, "+smf078IopFields.columns())
	pStmt078Cfg = createTable("smf078_cfg", "ID INT NOT NULL, SEQ INT NOT NULL, SUBSEQ INT NOT NULL, "+smf078CfgFields.columns())
}

func smf078Parse(rec []byte) {
	once078.Do(smf078Init)

	count078 += 1
	insert(pStmt078, rmfRow(rec, count078, "SMF78")...)
	rmfInsertSections(smf078Sections, rec, count078)
	if section(rec).num(18, 2) != 3 {
		return
	}
	// nested walks the sections located by the triplet at offset 4 of each section of the record triplet at.
	nested := func(at int, stmt *sql.Stmt, f fields) {
		off, length, n := tripletAt(rec, at)
		for k := 0; k < n; k++ {
			outer := section(rec).sub(off+k*length, length)
			ioff, ilen, in := outer.num(4, 4), outer.num(8, 2), outer.num(10, 2)
			for j := 0; j < in; j++ {
				insert(stmt, append([]interface{}{count078, k + 1, j + 1}, f.values(outer.sub(ioff+j*ilen, ilen))...)...)
			}
		}
	}
	nested(32, pStmt078Cfg, smf078CfgFields) // R783CPDS/L/N in each configuration control section
	nested(48, pStmt078Iop, smf078IopFields) // R783GIDS/L/N in the global section
}

var smf078CommonFields = fields{
	{"R782PA", 0, 'u', 4},
	{"R782PS", 4, 'u', 4},
	{"R782EPA", 8, 'u', 4},
	{"R782EPS", 12, 'u', 4},
	{"R782CA", 16, 'u', 4},
	{"R782CS", 20, 'u', 4},
	{"R782ECA", 24, 'u', 4},
	{"R782ECS", 28, 'u', 4},
	{"R782FLG", 32, 'u', 1},
	{"R782MLA", 40, 'u', 4},
	{"R782MLS", 44, 'u', 4},
	{"R782EMLA", 48, 'u', 4},
	{"R782EMLS", 52, 'u', 4},
	{"R782FLA", 56, 'u', 4},
	{"R782FLS", 60, 'u', 4},
	{"R782EFLA", 64, 'u', 4},
	{"R782EFLS", 68, 'u', 4},
	{"R782PLA", 72, 'u', 4},
	{"R782PLS", 76, 'u', 4},
	{"R782ELPA", 80, 'u', 4},
	{"R782ELPS", 84, 'u', 4},
	{"R782SA", 88, 'u', 4},
	{"R782SS", 92, 'u', 4},
	{"R782ESA", 96, 'u', 4},
	{"R782ESS", 100, 'u', 4},
	{"R782NA", 104, 'u', 4},
	{"R782NS", 108, 'u', 4},
	{"R782ENA", 112, 'u', 4},
	{"R782ENS", 116, 'u', 4},
	{"R782NL", 120, 'u', 4},
	{"R782ENL", 124, 'u', 4},
	{"R782LPAI", 128, 'u', 4},
	{"R782ELPI", 132, 'u', 4},
	{"R782MR", 136, 'u', 4},
	{"R782EMR", 140, 'u', 4},
	{"R782RUCA", 1724, 'u', 4},
	{"R782RUCS", 1728, 'u', 4},
	{"R782ERUCA", 1732, 'u', 4},
	{"R782ERUCS", 1736, 'u', 4},
}

var smf078IoqGlobalFields = fields{
	{"R783GFLG", 0, 'u', 1},
	{"R783GFLX", 1, 'u', 1},
	{"R783GNTR", 2, 'u', 2},
	{"R783GIDS", 4, 'u', 4},
	{"R783GIDL", 8, 'u', 2},
	{"R783GIDN", 10, 'u', 2},
	{"R783TSR", 16, 'u', 2},
	{"R783TOT", 20, 'u', 4},
	{"R783NXT", 24, 'u', 4},
	{"R783CFL", 28, 'u', 1},
	{"R783CSS", 29, 'u', 1},
	{"R783TNM", 32, 'c', 44},
	{"R783TSF", 76, 'c', 2},
	{"R783TOK", 80, 'c', 16},
	// skipped R783TDT at 80: redefines an earlier field
	// skipped R783TTM at 88: redefines an earlier field
	{"R783TDY", 96, 'c', 10},
}

var smf078IopFields = fields{
	{"R783IQID", 0, 'u', 2},
	{"R783IFLG", 2, 'u', 1},
	{"R783IQSM", 4, 'u', 4},
	{"R783IQCT", 8, 'u', 4},
	{"R783IIPB", 16, 'f', 8},
	{"R783IIPI", 24, 'f', 8},
	{"R783IIFS", 32, 'f', 8},
	{"R783IPII", 40, 'f', 8},
	{"R783ICPB", 48, 'f', 8},
	{"R783IDPB", 56, 'f', 8},
	{"R783ICUB", 64, 'f', 8},
	{"R783IDVB", 72, 'f', 8},
	{"R783ISCB", 80, 'f', 8},
	{"R783IECB", 88, 'f', 8},
}

var smf078CfgCtlFields = fields{
	{"R783ID1", 0, 'u', 2},
	{"R783NTR", 2, 'u', 2},
	{"R783CPDS", 4, 'u', 4},
	{"R783CPDL", 8, 'u', 2},
	{"R783CPDN", 10, 'u', 2},
	{"R783AMGC", 15, 'u', 1},
	{"R783AMGS", 16, 'u', 4},
}

var smf078CfgFields = fields{
	{"R783CPID", 0, 'u', 1},
	{"R783CPST", 1, 'u', 1},
	{"R783CUN", 2, 'u', 2},
	{"R783CU1", 4, 'u', 2},
	{"R783CU2", 6, 'u', 2},
	{"R783CU3", 8, 'u', 2},
	{"R783CU4", 10, 'u', 2},
	{"R783CUB", 12, 'u', 4},
	{"R783PT", 16, 'u', 4},
	{"R783DPB", 24, 'u', 4},
	{"R783CBT", 28, 'u', 4},
	{"R783CMR", 32, 'u', 4},
	{"R783SBS", 36, 'u', 4},
	{"R783CPXF", 44, 'u', 1},
	{"R783CPAT", 45, 'u', 1},
	{"R783CTMW", 48, 'u', 4},
	{"R783CTRD", 52, 'u', 4},
}

var smf078IoqFields = fields{
	{"R783ID2", 0, 'u', 2},
	{"R783DST", 2, 'u', 1},
	{"R783DSTX", 3, 'u', 1},
	{"R783QSM", 4, 'u', 4},
	{"R783QCT", 8, 'u', 4},
	{"R783MCMN", 12, 'u', 2},
	{"R783MCMX", 14, 'u', 2},
	{"R783MCDF", 16, 'u', 2},
	{"R783PTM", 20, 'u', 4},
	{"R783DPBM", 24, 'u', 4},
	{"R783CUBM", 28, 'u', 4},
	{"R783CBTM", 32, 'u', 4},
	{"R783CMRM", 36, 'u', 4},
	{"R783SBSM", 40, 'u', 4},
	{"R783DCTM", 44, 'u', 4},
	{"R783DDTM", 48, 'u', 4},
	{"R783CSST", 52, 'u', 4},
	{"R783HCNT", 56, 'u', 2},
	{"R783HIX", 58, 'u', 2},
	{"R783TMWM", 60, 'u', 4},
	{"R783TRDM", 64, 'u', 4},
}

var smf078HpavFields = fields{
	{"R783HLCU", 0, 'u', 2},
	{"R783HCU", 2, 'u', 2},
	{"R783HNAI", 8, 'f', 4},
	{"R783HTIO", 12, 'f', 4},
	{"R783HAIU", 16, 'u', 4},
	{"R783HCAD", 20, 'u', 4},
	{"R783HIOQ", 24, 'u', 4},
	{"R783XANC", 40, 'f', 4},
	{"R783XAUC", 44, 'f', 4},
	{"R783XNHC", 48, 'f', 4},
	{"R783XABC", 52, 'f', 4},
	{"R783XCBC", 56, 'u', 4},
	{"R783XHBC", 60, 'u', 4},
	{"R783XALC", 64, 'f', 4},
	{"R783XCLC", 68, 'u', 4},
	{"R783XHLC", 72, 'u', 4},
	{"R783XNAG", 76, 'f', 4},
	{"R783XCQD", 80, 'f', 4},
	{"R783XCIU", 84, 'f', 4},
}
