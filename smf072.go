package main

import (
	"database/sql"
	"sync"
)

// Type 72 (RMF workload activity, storage data and serialization delay). smf072 holds the header and RMF
// product section of every record; each subtype's sections go to child tables (see smf072Sections), each
// row with ID and SEQ, its 1-based position among the sections of that triplet in the record. SEQ is what
// the index fields of the manual refer to (for example R723CRTX, the first response time distribution
// section of a period, or R723RDNX, the first resource delay type name).
//
// The field tables were generated from the manual with tools/genfields.py, except smf072RtdFields: the
// manual gives the 14 response time distribution counters no names.

var pStmt072 *sql.Stmt
var once072 sync.Once
var count072 int

// smf072Section maps a triplet (at its offset in the header) to a child table. Several triplets share a
// table when the manual maps them with one structure; a field in the section tells them apart.
type smf072Section struct {
	subtype, at int
	table       string
	f           fields
	stmt        *sql.Stmt
}

var smf072Sections = []*smf072Section{
	{subtype: 3, at: 32, table: "smf072_wlm", f: smf072WlmFields},       // workload manager control
	{subtype: 3, at: 40, table: "smf072_served", f: smf072ServedFields}, // service class served
	{subtype: 3, at: 48, table: "smf072_rg", f: smf072RgFields},         // resource group
	{subtype: 3, at: 56, table: "smf072_scp", f: smf072ScpFields},       // service/report class period
	{subtype: 3, at: 64, table: "smf072_rtd", f: smf072RtdFields},       // response time distribution
	{subtype: 3, at: 72, table: "smf072_wrm", f: smf072WrmFields},       // work/resource manager state
	{subtype: 3, at: 80, table: "smf072_dn", f: smf072DnFields},         // resource delay type names
	{subtype: 4, at: 32, table: "smf072_stor", f: smf072StorFields},     // service class period (storage)
	{subtype: 4, at: 40, table: "smf072_swap", f: smf072SwapFields},     // swap reason
	{subtype: 5, at: 32, table: "smf072_ser", f: smf072SerFields},       // serialization control
	{subtype: 5, at: 40, table: "smf072_cms", f: smf072CmsFields},       // CMS lock (R725CMTY 1)
	{subtype: 5, at: 48, table: "smf072_cms", f: smf072CmsFields},       // CMS enqueue/dequeue lock (2)
	{subtype: 5, at: 56, table: "smf072_cms", f: smf072CmsFields},       // CMS latch lock (3)
	{subtype: 5, at: 64, table: "smf072_cms", f: smf072CmsFields},       // CMS SMF lock (4)
	{subtype: 5, at: 72, table: "smf072_local", f: smf072LocalFields},   // local lock
	{subtype: 5, at: 80, table: "smf072_cmlo", f: smf072CmloFields},     // CML lock owner
	{subtype: 5, at: 88, table: "smf072_cmlr", f: smf072CmlrFields},     // CML lock requestor
	{subtype: 5, at: 96, table: "smf072_latch", f: smf072LatchFields},   // GRS latch set creator (R725LATY 1)
	{subtype: 5, at: 104, table: "smf072_latch", f: smf072LatchFields},  // GRS latch requestor (2)
	{subtype: 5, at: 112, table: "smf072_enq", f: smf072EnqFields},      // GRS ENQ scope step (R725ENSC 1)
	{subtype: 5, at: 120, table: "smf072_enq", f: smf072EnqFields},      // GRS ENQ scope system (2)
	{subtype: 5, at: 128, table: "smf072_enq", f: smf072EnqFields},      // GRS ENQ scope systems (3)
	{subtype: 5, at: 136, table: "smf072_qscan", f: smf072QscanFields},  // GRS QScan statistics
}

func smf072Init() {
	pStmt072 = createTable("smf072", rmfColumns("SMF72"))
	stmts := map[string]*sql.Stmt{}
	for _, sec := range smf072Sections {
		if stmts[sec.table] == nil {
			stmts[sec.table] = createTable(sec.table, "ID INT NOT NULL, SEQ INT NOT NULL, "+sec.f.columns())
		}
		sec.stmt = stmts[sec.table]
	}
}

func smf072Parse(rec []byte) {
	once072.Do(smf072Init)

	count072 += 1
	insert(pStmt072, rmfRow(rec, count072, "SMF72")...)
	subtype := section(rec).num(18, 2)
	for _, sec := range smf072Sections {
		if sec.subtype != subtype {
			continue
		}
		off, length, n := tripletAt(rec, sec.at)
		for k := 0; k < n; k++ {
			insert(sec.stmt, append([]interface{}{count072, k + 1},
				sec.f.values(section(rec).sub(off+k*length, length))...)...)
		}
	}
}

// smf072RtdFields: counts of completed transactions by response time, as a percentage of the period's
// midpoint (R723RTDM): <=50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 200, 400 and >400%.
var smf072RtdFields = fields{
	{"R723RTD01", 0, 'u', 4}, {"R723RTD02", 4, 'u', 4}, {"R723RTD03", 8, 'u', 4}, {"R723RTD04", 12, 'u', 4},
	{"R723RTD05", 16, 'u', 4}, {"R723RTD06", 20, 'u', 4}, {"R723RTD07", 24, 'u', 4}, {"R723RTD08", 28, 'u', 4},
	{"R723RTD09", 32, 'u', 4}, {"R723RTD10", 36, 'u', 4}, {"R723RTD11", 40, 'u', 4}, {"R723RTD12", 44, 'u', 4},
	{"R723RTD13", 48, 'u', 4}, {"R723RTD14", 52, 'u', 4},
}

var smf072WlmFields = fields{
	{"R723MSCF", 0, 'u', 1},
	{"R723MFLG", 1, 'u', 1},
	{"R723MFL2", 2, 'u', 1},
	{"R723MNSP", 4, 'c', 8},
	{"R723MDSP", 12, 'c', 32},
	{"R723MTPA", 44, 't', 8},
	{"R723MCPU", 52, 'u', 4},
	{"R723MIOC", 56, 'u', 4},
	{"R723MMSO", 60, 'u', 4},
	{"R723MSRB", 64, 'u', 4},
	{"R723MTVL", 68, 'u', 4},
	{"R723MTV#", 72, 'u', 4},
	{"R723MOPT", 76, 'c', 2},
	{"R723MWNM", 80, 'c', 8},
	{"R723MWDE", 88, 'c', 32},
	{"R723MCNM", 120, 'c', 8},
	{"R723MCDE", 128, 'c', 32},
	{"R723MCPG", 160, 'u', 2},
	{"R723MSUB", 162, 'u', 1},
	{"R723MERF", 166, 'c', 6},
	{"R723MADJ", 172, 'u', 4},
	{"R723MIDN", 176, 'c', 8},
	{"R723MIDD", 184, 'c', 32},
	{"R723MTDI", 216, 't', 8},
	{"R723MIDU", 224, 'c', 8},
	{"R723CLSC", 232, 'c', 8},
	{"R723NFFI", 240, 'u', 4},
	{"R723NFFS", 244, 'u', 4},
	{"R723NADJ", 248, 'u', 4},
	{"R723CECA", 252, 'u', 4},
	{"R723MCF", 256, 'u', 4},
	{"R723MCFS", 260, 'u', 4},
	{"R723MCFI", 264, 'u', 4},
	{"R723CPA_actual", 268, 'u', 4},
	{"R723CPA_scaling_factor", 272, 'u', 4},
}

var smf072ServedFields = fields{
	{"R723SCSN", 0, 'c', 8},
	{"R723SCS#", 8, 'u', 4},
}

var smf072RgFields = fields{
	{"R723GGNM", 0, 'c', 8},
	{"R723GGDE", 8, 'c', 32},
	{"R723GGLT", 40, 'u', 1},
	{"R723GGTF", 41, 'u', 1},
	{"R723GGMN", 44, 'u', 4},
	{"R723GGMX", 48, 'u', 4},
	{"R723GGML", 52, 'u', 4},
	{"R723GGTI", 56, 'c', 8},
	{"R723GGTN", 64, 'c', 32},
	{"R723GGKY", 96, 'c', 64},
}

var smf072ScpFields = fields{
	{"R723CRTX", 0, 'u', 2},
	{"R723CWMX", 2, 'u', 2},
	{"R723CWMN", 4, 'u', 2},
	{"R723CRS1", 6, 'u', 1},
	{"R723CADF", 7, 'u', 1},
	{"R723CPER", 8, 'u', 1},
	{"R723CRTF", 9, 'u', 1},
	{"R723CRGF", 10, 'u', 1},
	{"R723CVAL", 12, 'u', 4},
	{"R723CPCT", 16, 'u', 2},
	{"R723CIMP", 18, 'u', 2},
	{"R723CDUR", 20, 'u', 4},
	{"R723CSRV", 24, 'f', 8},
	{"R723CCPU", 32, 'f', 8},
	{"R723CIOC", 40, 'f', 8},
	{"R723CMSO", 48, 'f', 8},
	{"R723CSRB", 56, 'f', 8},
	{"R723CPIR", 64, 'f', 8},
	{"R723CHPI", 72, 'f', 8},
	{"R723CBPI", 80, 'f', 8},
	{"R723CPIE", 88, 'f', 8},
	{"R723CBPE", 96, 'f', 8},
	{"R723CBKA", 104, 'f', 8},
	{"R723CBKE", 112, 'f', 8},
	{"R723CPRS", 120, 'f', 8},
	{"R723CERS", 128, 'f', 8},
	{"R723CTRR", 136, 'f', 8},
	{"R723CTAT", 144, 'f', 8},
	{"R723CRCT", 152, 'f', 8},
	{"R723CIIT", 160, 'f', 8},
	{"R723CHST", 168, 'f', 8},
	{"R723CSWC", 176, 'u', 4},
	{"R723CCRM", 180, 'u', 4},
	{"R723CRCP", 184, 'u', 4},
	{"R723CARC", 188, 'u', 4},
	{"R723CNCP", 192, 'u', 4},
	{"R723CANC", 196, 'u', 4},
	{"R723CTET", 200, 'f', 8},
	{"R723CXET", 208, 'f', 8},
	{"R723CETS", 216, 'f', 8},
	{"R723CCUS", 224, 'u', 4},
	{"R723CTOT", 228, 'u', 4},
	{"R723CCDE", 232, 'u', 4},
	{"R723CCCA", 236, 'u', 4},
	{"R723CSWI", 240, 'u', 4},
	{"R723CMPL", 244, 'u', 4},
	{"R723CAPR", 248, 'u', 4},
	{"R723CACO", 252, 'u', 4},
	{"R723CAXM", 256, 'u', 4},
	{"R723CVIO", 260, 'u', 4},
	{"R723CHSP", 264, 'u', 4},
	{"R723CCHS", 268, 'u', 4},
	{"R723CUNK", 272, 'u', 4},
	{"R723CIDL", 276, 'u', 4},
	{"R723CPDE", 280, 'u', 4},
	{"R723CPQU", 284, 'u', 4},
	{"R723CSAC", 288, 'u', 4},
	{"R723CSRS", 292, 'f', 8},
	{"R723CSPA", 300, 'f', 8},
	{"R723CSPE", 308, 'f', 8},
	{"R723CICT", 316, 'f', 8},
	{"R723CIWT", 324, 'f', 8},
	{"R723CIDT", 332, 'f', 8},
	{"R723CIRC", 340, 'u', 4},
	{"R723CTOU", 344, 'u', 4},
	{"R723CIOU", 348, 'u', 4},
	{"R723CIOD", 352, 'u', 4},
	{"R723CQ", 356, 'u', 4},
	{"R723CSPV", 360, 'u', 4},
	{"R723CSVI", 364, 'u', 4},
	{"R723CSHS", 368, 'u', 4},
	{"R723CSMP", 372, 'u', 4},
	{"R723CSSW", 376, 'u', 4},
	{"R723CNDI", 380, 'u', 4},
	{"R723CTDQ", 384, 'u', 4},
	{"R723CTSA", 388, 'f', 8},
	{"R723CIOT", 396, 'f', 8},
	{"R723CQDT", 404, 'f', 8},
	{"R723CADT", 412, 'f', 8},
	{"R723CCVT", 420, 'f', 8},
	{"R723CIQT", 428, 'f', 8},
	{"R723CIEA", 436, 'f', 8},
	{"R723CXEA", 444, 'f', 8},
	{"R723CFEA", 452, 'f', 8},
	{"R723CAMU", 460, 'u', 4},
	{"R723CAMD", 464, 'u', 4},
	{"R723APU", 468, 'u', 4},
	{"R723APD", 472, 'u', 4},
	{"R723FQD", 476, 'u', 4},
	{"R723PLSC", 480, 'c', 8},
	{"R723RCOD", 488, 'u', 4},
	{"R723RCOU", 492, 'u', 4},
	{"R723ECTC", 496, 'f', 8},
	{"R723IFAU", 504, 'u', 4},
	{"R723IFCU", 508, 'u', 4},
	{"R723IFAD", 512, 'u', 4},
	{"R723SUPU", 532, 'u', 4},
	{"R723SUCU", 536, 'u', 4},
	{"R723SUPD", 540, 'u', 4},
	{"R723CSUP", 544, 'f', 8},
	{"R723CSUC", 552, 'f', 8},
	{"R723CIFA", 560, 'f', 8},
	{"R723CIFC", 568, 'f', 8},
	{"R723TPDP", 576, 'f', 8},
	{"R723CPDP", 584, 'f', 8},
	{"R723LPDP", 592, 'f', 8},
	{"R723SPDP", 600, 'f', 8},
	{"R723RTDM", 608, 'u', 4},
	{"R723RTDC", 612, 'u', 4},
	{"R723RTDT", 616, 't', 8},
	{"R723TSUCP", 624, 'f', 8},
	{"R723TSUSP", 632, 'f', 8},
	{"R723TSUOCP", 640, 'f', 8},
	{"R723MSUCP", 648, 'f', 8},
	{"R723MSUSP", 656, 'f', 8},
	{"R723MSUOCP", 664, 'f', 8},
	{"R723ASUCP", 672, 'f', 8},
	{"R723ASUSP", 680, 'f', 8},
	{"R723ASUOCP", 688, 'f', 8},
	{"R723BSUCP", 696, 'f', 8},
	{"R723BSUSP", 704, 'f', 8},
	{"R723BSUOCP", 712, 'f', 8},
	{"R723CTETX", 720, 'f', 8},
	{"R723CXETX", 728, 'f', 8},
	{"R723CETSX", 736, 'f', 8},
	{"R723CQDTX", 744, 'f', 8},
	{"R723CADTX", 752, 'f', 8},
	{"R723CCVTX", 760, 'f', 8},
	{"R723CIQTX", 768, 'f', 8},
	{"R723ENCTRXNUM", 792, 'u', 4},
	{"R723ENCTRXCALLS", 796, 'u', 4},
	{"R723ENCTRXET", 800, 'f', 8},
	{"R723ENCTRXETS", 808, 'f', 8},
}

var smf072WrmFields = fields{
	{"R723RTYP", 0, 'c', 4},
	{"R723RFLG", 4, 'u', 1},
	{"R723RESS", 8, 'u', 4},
	{"R723RACT", 12, 'u', 4},
	{"R723RRDY", 16, 'u', 4},
	{"R723RIDL", 20, 'u', 4},
	{"R723RWLO", 24, 'u', 4},
	{"R723RWIO", 28, 'u', 4},
	{"R723RWCO", 32, 'u', 4},
	{"R723RWDS", 36, 'u', 4},
	{"R723RWSL", 40, 'u', 4},
	{"R723RWSN", 44, 'u', 4},
	{"R723RWSS", 48, 'u', 4},
	{"R723RWTM", 52, 'u', 4},
	{"R723RWO", 56, 'u', 4},
	{"R723RWMS", 60, 'u', 4},
	{"R723RSSL", 64, 'u', 4},
	{"R723RSSS", 68, 'u', 4},
	{"R723RSSN", 72, 'u', 4},
	{"R723RWST", 76, 'u', 4},
	{"R723RWRT", 80, 'u', 4},
	{"R723RWWR", 84, 'u', 4},
	{"R723RAPP", 88, 'u', 4},
	{"R723RWNL", 92, 'u', 4},
	{"R723RW01", 96, 'u', 4},
	{"R723RW02", 100, 'u', 4},
	{"R723RW03", 104, 'u', 4},
	{"R723RW04", 108, 'u', 4},
	{"R723RW05", 112, 'u', 4},
	{"R723RW06", 116, 'u', 4},
	{"R723RW07", 120, 'u', 4},
	{"R723RW08", 124, 'u', 4},
	{"R723RW09", 128, 'u', 4},
	{"R723RW10", 132, 'u', 4},
	{"R723RW11", 136, 'u', 4},
	{"R723RW12", 140, 'u', 4},
	{"R723RW13", 144, 'u', 4},
	{"R723RW14", 148, 'u', 4},
	{"R723RW15", 152, 'u', 4},
	{"R723RBPM", 156, 'u', 4},
	{"R723RDNX", 160, 'u', 2},
	{"R723RDNN", 162, 'u', 2},
}

var smf072DnFields = fields{
	{"R723DNST", 0, 'c', 4},
	{"R723DNNU", 4, 'u', 2},
	{"R723DNDE", 6, 'c', 16},
}

var smf072StorFields = fields{
	{"R724PNAM", 0, 'c', 8},
	{"R724PTM", 8, 't', 8},
	{"R724LCNM", 16, 'c', 8},
	{"R724PER#", 24, 'u', 1},
	{"R724USER", 28, 'u', 4},
	{"R724ACTV", 32, 'u', 4},
	{"R724ACTS", 36, 'u', 4},
	{"R724IDLS", 40, 'u', 4},
	{"R724PAGE", 44, 'u', 4},
	{"R724SWAP", 48, 'u', 4},
	{"R724OUTR", 52, 'u', 4},
	{"R724PGIN", 56, 'u', 4},
	{"R724DIVS", 60, 'u', 4},
	{"R724LSSA", 64, 'u', 4},
	{"R724PSSA", 68, 'u', 4},
	{"R724UPRO", 72, 'u', 4},
	{"R724UDEV", 76, 'u', 4},
	{"R724DPRO", 80, 'u', 4},
	{"R724DDEV", 84, 'u', 4},
	{"R724DSTO", 88, 'u', 4},
	{"R724DJES", 92, 'u', 4},
	{"R724DHSM", 96, 'u', 4},
	{"R724DXCF", 100, 'u', 4},
	{"R724DENQ", 104, 'u', 4},
	{"R724DMNT", 108, 'u', 4},
	{"R724DMSG", 112, 'u', 4},
	{"R724UNKN", 116, 'u', 4},
	{"R724VALD", 120, 'u', 4},
	{"R724LSCT", 124, 'u', 4},
	{"R724ESCT", 128, 'u', 4},
	{"R724PSCT", 132, 'u', 4},
	{"R724ACTF", 136, 'f', 8},
	{"R724IDLE", 144, 'f', 8},
	{"R724SLOT", 152, 'f', 8},
	{"R724DIV", 160, 'f', 8},
	{"R724FIX", 168, 'f', 8},
	{"R724LSCF", 176, 'f', 8},
	{"R724LSEF", 184, 'f', 8},
	{"R724PSEF", 192, 'f', 8},
	{"R724VECT", 200, 'f', 8},
	{"R724ET", 208, 'f', 8},
	{"R724QT", 216, 'f', 8},
	{"R724END", 224, 'f', 8},
	{"R724TSV", 232, 'f', 8},
	{"R724VIN", 240, 'f', 8},
	{"R724VLC", 248, 'f', 8},
	{"R724GPI", 256, 'f', 8},
	{"R724ETX", 264, 'f', 8},
	{"R724QTX", 272, 'f', 8},
}

var smf072SwapFields = fields{
	{"R724OR1", 0, 'u', 4},
	{"R724OR2", 4, 'u', 4},
	{"R724OR3", 8, 'u', 4},
	{"R724OR4", 12, 'u', 4},
	{"R724OR5", 16, 'u', 4},
	{"R724OR6", 20, 'u', 4},
	{"R724OR7", 24, 'u', 4},
	{"R724OR8", 28, 'u', 4},
	{"R724OR9", 32, 'u', 4},
	{"R724OR10", 36, 'u', 4},
	{"R724OR11", 40, 'u', 4},
	{"R724OR12", 44, 'u', 4},
	{"R724OR13", 48, 'u', 4},
	{"R724OR14", 52, 'u', 4},
	{"R724OR15", 56, 'u', 4},
	{"R724OR16", 60, 'u', 4},
	{"R724OR17", 64, 'u', 4},
	{"R724OR18", 68, 'u', 4},
	{"R724OR7A", 72, 'u', 4},
}

var smf072SerFields = fields{
	{"R725SGMO", 0, 'u', 1},
	{"R725SCMS", 16, 'u', 8},
	{"R725SCMA", 24, 'u', 8},
	{"R725SCMT", 32, 'u', 8},
	{"R725SEDS", 48, 'u', 8},
	{"R725SEDA", 56, 'u', 8},
	{"R725SEDT", 64, 'u', 8},
	{"R725SLAS", 80, 'u', 8},
	{"R725SLAA", 88, 'u', 8},
	{"R725SLAT", 96, 'u', 8},
	{"R725SSMS", 112, 'u', 8},
	{"R725SSMA", 120, 'u', 8},
	{"R725SSMT", 128, 'u', 8},
	{"R725SLOS", 144, 'u', 8},
	{"R725SLOA", 152, 'u', 8},
	{"R725SLOT", 160, 'u', 8},
	{"R725SCLS", 176, 'u', 8},
	{"R725SCLA", 184, 'u', 8},
	{"R725SCLT", 192, 'u', 8},
	{"R725SLRS", 208, 'u', 8},
	{"R725SLRT", 216, 'u', 8},
	{"R725SLRQ", 224, 'w', 16},
	{"R725SSTR", 240, 'u', 8},
	{"R725SSTS", 248, 'u', 8},
	{"R725SSTT", 256, 'u', 8},
	{"R725SSTQ", 272, 'w', 16},
	{"R725SSYR", 288, 'u', 8},
	{"R725SSYS", 296, 'u', 8},
	{"R725SSYT", 304, 'u', 8},
	{"R725SSYQ", 320, 'w', 16},
	{"R725SSSR", 336, 'u', 8},
	{"R725SSSS", 344, 'u', 8},
	{"R725SSST", 352, 'u', 8},
	{"R725SSSQ", 368, 'w', 16},
}

var smf072CmsFields = fields{
	{"R725CMJN", 0, 'c', 8},
	{"R725CMSP", 11, 'u', 1},
	{"R725CMAS", 12, 'u', 2},
	{"R725CMST", 16, 'c', 8},
	{"R725CMSN", 24, 'c', 8},
	{"R725CMTY", 32, 'u', 1},
	{"R725CMSU", 40, 'u', 8},
	{"R725CMAL", 48, 'u', 8},
	{"R725CMTI", 56, 'u', 8},
}

var smf072LocalFields = fields{
	{"R725LOJN", 0, 'c', 8},
	{"R725LOSP", 11, 'u', 1},
	{"R725LOAS", 12, 'u', 2},
	{"R725LOST", 16, 'c', 8},
	{"R725LOSN", 24, 'c', 8},
	{"R725LOSU", 40, 'u', 8},
	{"R725LOAL", 48, 'u', 8},
	{"R725LOTI", 56, 'u', 8},
	{"R725LCSU", 64, 'u', 8},
	{"R725LCAL", 72, 'u', 8},
	{"R725LCTI", 80, 'u', 8},
}

var smf072CmloFields = fields{
	{"R725COJN", 0, 'c', 8},
	{"R725COSP", 11, 'u', 1},
	{"R725COAS", 12, 'u', 2},
	{"R725COST", 16, 'c', 8},
	{"R725COSN", 24, 'c', 8},
	{"R725COSU", 40, 'u', 8},
	{"R725COAL", 48, 'u', 8},
	{"R725COTI", 56, 'u', 8},
	{"R725CLSU", 64, 'u', 8},
	{"R725CLAL", 72, 'u', 8},
	{"R725CLTI", 80, 'u', 8},
}

var smf072CmlrFields = fields{
	{"R725CRJN", 0, 'c', 8},
	{"R725CRSP", 11, 'u', 1},
	{"R725CRAS", 12, 'u', 2},
	{"R725CRST", 16, 'c', 8},
	{"R725CRSN", 24, 'c', 8},
	{"R725CRSU", 40, 'u', 8},
	{"R725CRAL", 48, 'u', 8},
	{"R725CRTI", 56, 'u', 8},
}

var smf072LatchFields = fields{
	{"R725LAJN", 0, 'c', 8},
	{"R725LASP", 11, 'u', 1},
	{"R725LAAS", 12, 'u', 2},
	{"R725LAST", 16, 'c', 8},
	{"R725LASN", 24, 'c', 8},
	{"R725LATY", 32, 'u', 1},
	{"R725LASU", 40, 'u', 8},
	{"R725LATI", 48, 'u', 8},
	{"R725LASQ", 56, 'w', 16},
}

var smf072EnqFields = fields{
	{"R725ENJN", 0, 'c', 8},
	{"R725ENSP", 11, 'u', 1},
	{"R725ENAS", 12, 'u', 2},
	{"R725ENST", 16, 'c', 8},
	{"R725ENSN", 24, 'c', 8},
	{"R725ENSC", 32, 'u', 1},
	{"R725ENRC", 40, 'u', 8},
	{"R725ENSU", 48, 'u', 8},
	{"R725ENTI", 56, 'u', 8},
	{"R725ENSQ", 64, 'w', 16},
}

var smf072QscanFields = fields{
	{"R725QSJN", 0, 'c', 8},
	{"R725QSSP", 11, 'u', 1},
	{"R725QSAS", 12, 'u', 2},
	{"R725QSST", 16, 'c', 8},
	{"R725QSSN", 24, 'c', 8},
	{"R725QSRC", 40, 'u', 8},
	{"R725QSSC", 48, 'u', 8},
	{"R725QSRR", 56, 'u', 8},
	{"R725QSRQ", 64, 'w', 16},
	{"R725QSTI", 80, 'u', 8},
	{"R725QSTQ", 88, 'w', 16},
}
