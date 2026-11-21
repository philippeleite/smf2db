package main

import (
	"database/sql"
	"sync"
)

// SMF type 99: System Resource Manager decisions (one record per policy interval and subtype).
// smf099 has the header, the product section, and the fields of the sections that occur once per record:
// the subtype 3 class name (SMF99_PNAM) and the subtype 4 device cluster ID (SMF994DEVCLID).
// The subtype self-defining section starts at SMF99DOF; its triplets map to child tables with ID and SEQ
// (see rmfSection). Period sections carry triplets of their own, with offsets from the start of the record:
// their rows have ID, SEQ of the period and SUBSEQ. Plots (paging plot, subtype 3 period plots, I/O plot)
// share smf099_plot, told apart by PLOT (the IBM name prefix: PAGP, PPRP, MPLP, RUAP, SWPP, PASP, IO_PLOT),
// and their one-curve points share smf099_point.
// Not decoded, as the sample data has none: subtypes 5, 7 and 9 to 14; the subtype 1 software licensing
// table and internal use sections; the subtype 2 EWLM class, cross memory delay, queue server, I/O subsystem
// and AI data sections; the subtype 3 queue delay, queue ready user, queue service time and server
// instance plots; the subtype 8 image CPU and SYSH CPU plot sections.
// The field tables were generated from the manual with tools/genfields.py.

var pStmt099, pStmt099Plot, pStmt099Point *sql.Stmt
var once099 sync.Once
var count099 int

var smf099HeaderFields = fields{
	{"SMF99TME", 2, 'T', 4},
	{"SMF99DTE", 6, 'd', 4},
	{"SMF99SID", 10, 'c', 4},
	{"SMF99SSID", 14, 'c', 4},
	{"SMF99TID", 18, 'u', 2},
}

// at is relative to the start of the subtype self-defining section (SMF99DOF).
var smf099Sections = []*rmfSection{
	{subtype: 1, at: 0, table: "smf099_trace", f: smf099TraceFields},
	{subtype: 1, at: 8, table: "smf099_state", f: smf099StateFields},
	{subtype: 1, at: 24, table: "smf099_prty", f: smf099PrtyFields},
	{subtype: 1, at: 32, table: "smf099_rg", f: smf099RgFields},
	{subtype: 1, at: 40, table: "smf099_gr", f: smf099GrFields},
	{subtype: 1, at: 48, table: "smf099_sl", f: smf099SlFields},
	{subtype: 1, at: 64, table: "smf099_prty_zaap", f: smf099PrtyFields},
	{subtype: 1, at: 80, table: "smf099_prty_ziip", f: smf099PrtyFields},
	{subtype: 2, at: 0, table: "smf099_class", f: smf099ClassFields},
	{subtype: 2, at: 8, table: "smf099_period", f: smf099PeriodFields},
	{subtype: 3, at: 8, table: "smf099_s3_period", f: fields{{"SMF993_INTERNAL_CLASS_NAME", 56, 'c', 8}}},
	{subtype: 4, at: 4, table: "smf099_io_prty", f: smf099IoPrtyFields},
	{subtype: 6, at: 0, table: "smf099_s6_period", f: smf099S6PeriodFields},
	{subtype: 8, at: 0, table: "smf099_lpar", f: smf099LparFields},
	{subtype: 8, at: 8, table: "smf099_s8_prty", f: smf099S8PrtyFields},
	{subtype: 8, at: 16, table: "smf099_s8_period", f: smf099S8PeriodFields},
}

// smf099Nested maps a triplet inside each section of a record triplet to a child table.
type smf099Nest struct {
	subtype, at, inner int
	table              string
	f                  fields
	stmt               *sql.Stmt
}

var smf099Nests = []*smf099Nest{
	{subtype: 2, at: 8, inner: 164, table: "smf099_server", f: smf099ServerFields},     // SMF99_PSERV_OF/LN/ON
	{subtype: 2, at: 8, inner: 172, table: "smf099_as_esp", f: smf099EspFields},        // SMF99_PESP_OF/LN/ON
	{subtype: 2, at: 8, inner: 212, table: "smf099_sdata", f: smf099SdataFields},       // SMF99_PSDATA_OF/LN/ON
	{subtype: 6, at: 0, inner: 52, table: "smf099_s6_server", f: smf099S6ServerFields}, // SMF996_SERVER_DATA_OF/LN/ON
}

// smf099S3Plots are the plots of a subtype 3 period with data in the samples, by triplet offset in the
// period self-defining section.
var smf099S3Plots = []struct {
	at   int
	name string
}{{0, "PPRP"}, {8, "MPLP"}, {16, "RUAP"}, {24, "SWPP"}, {32, "PASP"}}

func smf099Init() {
	pStmt099 = createTable("smf099", "ID INT NOT NULL, "+smf099HeaderFields.columns()+", "+smf099ProductFields.columns()+
		", SMF99_PNAM TEXT, SMF994DEVCLID INT")
	rmfCreateSections(smf099Sections)
	for _, n := range smf099Nests {
		n.stmt = createTable(n.table, "ID INT NOT NULL, SEQ INT NOT NULL, SUBSEQ INT NOT NULL, "+n.f.columns())
	}
	pStmt099Plot = createTable("smf099_plot", "ID INT NOT NULL, SEQ INT NOT NULL, PLOT TEXT NOT NULL, "+
		"SMF99_PLOT_PNUM INT, SMF99_PLOT_BW INT, SMF99_PLOT_LSTX INT")
	pStmt099Point = createTable("smf099_point", "ID INT NOT NULL, SEQ INT NOT NULL, PLOT TEXT NOT NULL, PT INT NOT NULL, "+
		"SMF99_PLOT_XVAL INT, SMF99_PLOT_YVAL INT")
}

func smf099Parse(rec []byte) {
	once099.Do(smf099Init)
	count099++
	id, s := count099, section(rec)
	subtype := s.num(18, 2)
	poff, plen, _ := tripletAt(rec, 24) // SMF99POF/PLN/PON
	dof, _, _ := tripletAt(rec, 32)     // SMF99DOF/DLN/DON
	if dof == 0 {
		dof = len(rec) // no data section: every triplet below reads as absent
	}
	var pnam, devclid interface{}
	switch subtype {
	case 3:
		if off, length, n := tripletAt(rec, dof); n > 0 { // SMF993COF/CLN/CON
			pnam = s.sub(off, length).text(0, 8)
		}
	case 4:
		devclid = s.u32(dof)
	}
	insert(pStmt099, append(append(append([]interface{}{id}, smf099HeaderFields.values(s)...),
		smf099ProductFields.values(s.sub(poff, plen))...), pnam, devclid)...)

	for _, sec := range smf099Sections {
		if sec.subtype != subtype {
			continue
		}
		off, length, n := tripletAt(rec, dof+sec.at)
		for k := 0; k < n; k++ {
			insert(sec.stmt, append([]interface{}{id, k + 1}, sec.f.values(s.sub(off+k*length, length))...)...)
		}
	}
	for _, nest := range smf099Nests {
		if nest.subtype != subtype {
			continue
		}
		off, length, n := tripletAt(rec, dof+nest.at)
		for k := 0; k < n; k++ {
			ioff, ilen, in := tripletAt(rec, off+k*length+nest.inner)
			for j := 0; j < in; j++ {
				insert(nest.stmt, append([]interface{}{id, k + 1, j + 1}, nest.f.values(s.sub(ioff+j*ilen, ilen))...)...)
			}
		}
	}

	// plot inserts a plot header and its points; pnum tells whether the header starts with a period number.
	// Its point triplet is ordered offset, number, length.
	plot := func(seq int, name string, p section, pnum bool) {
		var pn interface{}
		d := 0
		if pnum {
			pn, d = p.u32(0), 4
		}
		insert(pStmt099Plot, id, seq, name, pn, p.u32(d), p.u32(d+4))
		off, n, length := p.num(d+12, 4)-4, p.num(d+16, 2), p.num(d+18, 2)
		for j := 0; off >= 0 && j < n; j++ {
			pt := s.sub(off+j*length, length)
			insert(pStmt099Point, id, seq, name, j+1, pt.u32(0), pt.u32(4))
		}
	}
	switch subtype {
	case 1:
		if off, length, n := tripletAt(rec, dof+16); n > 0 { // SMF99PPOF/PPLN/PPON
			plot(1, "PAGP", s.sub(off, length), false)
		}
	case 3:
		off, length, n := tripletAt(rec, dof+8) // SMF993CPOF/CPLN/CPON
		for k := 0; k < n; k++ {
			for _, pl := range smf099S3Plots {
				poff, plen, pn := tripletAt(rec, off+k*length+pl.at)
				for j := 0; j < pn; j++ {
					plot(k+1, pl.name, s.sub(poff+j*plen, plen), true)
				}
			}
		}
	case 4:
		if off, length, n := tripletAt(rec, dof+12); n > 0 { // SMF994IOPLTOF/LN/ON
			plot(1, "IO_PLOT", s.sub(off, length), false)
		}
	}
}

var smf099ProductFields = fields{
	{"SMF99VN2", 0, 'u', 2},
	{"SMF99RVN", 2, 'c', 2}, // EBCDIC digits ("33"), though the manual says binary
	{"SMF99PNM", 4, 'c', 8},
	{"SMF99SLV", 12, 'c', 8},
	{"SMF99SNM", 20, 'c', 8},
	{"SMP99FLG", 28, 'u', 1},
}

var smf099TraceFields = fields{
	{"SMF99_TPID", 0, 'u', 1},
	{"SMF99_TRID", 1, 'u', 1},
	{"SMF99_TCOD", 2, 'u', 2},
	{"SMF99_TJOB", 4, 'c', 8},
	{"SMF99_TLPI", 12, 'u', 4},
	{"SMF99_TSPI", 16, 'u', 4},
	{"SMF99_TGSR", 20, 'u', 4},
	{"SMF99_TRGN", 36, 'c', 8},
	{"SMF99_TCNM", 44, 'c', 8},
	{"SMF99_TPER", 52, 'u', 2},
	{"SMF99_TASID", 54, 'u', 2},
}

var smf099StateFields = fields{
	{"SMF99_CPUA", 0, 'u', 2},
	{"SMF99_UMP", 2, 'u', 2},
	{"SMF99_UIC1", 4, 'u', 4},
	{"SMF99_UIC2", 8, 'u', 4},
	{"SMF99_UIC3", 12, 'u', 4},
	{"SMF99_UIC4", 16, 'u', 4},
	{"SMF99_EUIC1", 20, 'u', 4},
	{"SMF99_EUIC2", 24, 'u', 4},
	{"SMF99_EUIC3", 28, 'u', 4},
	{"SMF99_EUIC4", 32, 'u', 4},
	{"SMF99_FRV1", 36, 'u', 2},
	{"SMF99_FRV2", 38, 'u', 2},
	{"SMF99_FRV3", 40, 'u', 2},
	{"SMF99_ESTB1", 42, 'u', 2},
	{"SMF99_ESTB2", 44, 'u', 2},
	{"SMF99_ESTB3", 46, 'u', 2},
	{"SMF99_W2MIG", 48, 'u', 4},
	{"SMF99_PTAVAIL", 52, 'u', 4},
	{"SMF99_SHORT_FLAGS", 56, 'u', 1},
	{"SMF99_STATUS_FLAGS", 57, 'u', 1},
	{"SMF99_TOTAL_PAG_COST", 58, 'u', 2},
	{"SMF99_CPPS", 60, 'u', 4},
	{"SMF99_ILSU_ARRAY0", 64, 'u', 4},
	{"SMF99_ILSU_ARRAY1", 68, 'u', 4},
	{"SMF99_ILSU_ARRAY2", 72, 'u', 4},
	{"SMF99_ILSU_ARRAY3", 76, 'u', 4},
	{"SMF99_ILSU_ARRAY4", 80, 'u', 4},
	{"SMF99_ILSU_ARRAY5", 84, 'u', 4},
	{"SMF99_ILSU_ARRAY6", 88, 'u', 4},
	{"SMF99_ILSU_ARRAY7", 92, 'u', 4},
	{"SMF99_SUIC1", 96, 'u', 4},
	{"SMF99_SUIC2", 100, 'u', 4},
	{"SMF99_SUIC3", 104, 'u', 4},
	{"SMF99_SUIC4", 108, 'u', 4},
	{"SMF99_SEUC1", 112, 'u', 4},
	{"SMF99_SEUC2", 116, 'u', 4},
	{"SMF99_SEUC3", 120, 'u', 4},
	{"SMF99_SEUC4", 124, 'u', 4},
	{"SMF99_STWSS", 128, 'u', 4},
	{"SMF99_NUM_EXT_SC", 132, 'u', 4},
	{"SMF99_DEFAULT_IO_VELOCITY", 136, 'u', 4},
	{"SMF99_SU_IFACTOR", 140, 'u', 4},
	{"SMF99_StgCrit_Hsk_Skip_Clock0", 144, 'u', 2},
	{"SMF99_StgCrit_Hsk_Skip_Clock1", 146, 'u', 2},
	{"SMF99_StgCrit_Hsk_Skip_Clock2", 148, 'u', 2},
	{"SMF99_StgCrit_Hsk_Skip_Clock3", 150, 'u', 2},
	{"SMF99_StgCrit_Hsk_Skip_Clock4", 152, 'u', 2},
	{"SMF99_StgCrit_Hsk_Skip_Clock5", 154, 'u', 2},
	{"SMF99_StgCrit_Hsk_Skip_Clock6", 156, 'u', 2},
	{"SMF99_LS_DISC", 160, 'u', 4},
	{"SMF99_CAPWS", 164, 'u', 4},
	{"SMF99_SECWS", 168, 'u', 4},
	{"SMF99_PGINS", 172, 'u', 4},
	{"SMF99_IFA_NORMALIZATION", 176, 'u', 4},
	{"SMF99_CPUS_ONLINE", 180, 'u', 2},
	{"SMF99_IFAS_ONLINE", 182, 'u', 2},
	{"SMF99_IFAA", 184, 'u', 2},
	{"SMF99_PROA", 186, 'u', 2},
	{"SMF99_IFA_FLAGS", 188, 'u', 1},
	{"SMF99_FREE_LPAR_CAPACITY_WT_RELATED", 192, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_GUARANTEED", 196, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_CEC_RELATED", 200, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_LCP_CONFIG", 204, 'u', 4},
	{"SMF99_ITAVAIL", 208, 'u', 4},
	{"SMF99_SUP_NORMALIZATION", 212, 'u', 4},
	{"SMF99_SUPS_ONLINE", 216, 'u', 2},
	{"SMF99_SUPA", 218, 'u', 2},
	{"SMF99_GUARANTED_IMAGE_CAPACITY", 220, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY0", 224, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY1", 228, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY2", 232, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY3", 236, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY4", 240, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY5", 244, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY6", 248, 'u', 4},
	{"SMF99_ILSU_ZAAP_ARRAY7", 252, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY0", 256, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY1", 260, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY2", 264, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY3", 268, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY4", 272, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY5", 276, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY6", 280, 'u', 4},
	{"SMF99_ILSU_ZIIP_ARRAY7", 284, 'u', 4},
	{"SMF99_CCTINTHD", 288, 'u', 2},
	{"SMF99_CCTTRPCT", 290, 'u', 2},
	{"SMF99_CCTTRATE", 292, 'u', 4},
	{"SMF99_CCCTTSH", 296, 'u', 4},
	{"SMF99_CCTRC100", 300, 'u', 4},
	{"SMF99_CCTRCDSP", 304, 'u', 4},
	{"SMF99_CCTRCUSE", 308, 'u', 4},
	{"SMF99_CCTRCWTR", 312, 'u', 4},
	{"SMF99_CCCITTSH", 316, 'u', 4},
	{"SMF99_STAVAIL", 320, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_WT_RELATED_ZAAP", 324, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_GUARANTEED_ZAAP", 328, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_CEC_RELATD_ZAAP", 332, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_LCP_CONFIG_ZAAP", 336, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_WT_RELATED_ZIIP", 340, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_GUARANTEED_ZIIP", 344, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_CEC_RELATD_ZIIP", 348, 'u', 4},
	{"SMF99_FREE_LPAR_CAPACITY_LCP_CONFIG_ZIIP", 352, 'u', 4},
	{"SMF99_SVTWTSS", 356, 'u', 4},
	{"SMF99_BOOSTINFO", 644, 'u', 1},
}

var smf099PrtyFields = fields{
	{"SMF99_PTPRTY", 0, 'u', 2},
	{"SMF99_PTNP", 2, 'u', 2},
	{"SMF99_PTIMDP", 4, 'u', 4},
	{"SMF99_PTPMDP", 8, 'u', 4},
	{"SMF99_PTCPUU", 12, 'u', 4},
	{"SMF99_PTCPUD", 16, 'u', 4},
	{"SMF99_PTW2UR", 20, 'u', 4},
	{"SMF99_PTAPU", 24, 'u', 4},
	{"SMF99_PTPPU", 28, 'u', 4},
	{"SMF99_PTACMD", 32, 'u', 4},
	{"SMF99_PTIMAXD", 40, 'u', 4},
	{"SMF99_PTWMAXD", 44, 'u', 4},
	{"SMF99_PTIAMTW", 48, 'u', 4},
	{"SMF99_PTWAMTW", 52, 'u', 4},
	{"SMF99_PTSCPUU", 56, 'u', 4},
	{"SMF99_PTSCPUD", 60, 'u', 4},
}

var smf099RgFields = fields{
	{"SMF99_RGNAME", 0, 'c', 8},
	{"SMF99_MIN_SR", 8, 'u', 4},
	{"SMF99_MAX_SR", 12, 'u', 4},
	{"SMF99_ACT_SR", 16, 'u', 4},
	{"SMF99_SPAS", 20, 'u', 4},
	{"SMF99_SLICES", 24, 'u', 2},
	{"SMF99_RHELPCNT0", 26, 'u', 2},
	{"SMF99_RHELPCNT1", 28, 'u', 2},
	{"SMF99_RHELPCNT2", 30, 'u', 2},
	{"SMF99_RHELPCNT3", 32, 'u', 2},
	{"SMF99_RHELPCNT4", 34, 'u', 2},
	{"SMF99_RHELPCNT5", 36, 'u', 2},
	{"SMF99_RHELPCNT6", 38, 'u', 2},
	{"SMF99_LHELP_FLGS", 40, 'u', 1},
	{"SMF99_RG_FLAGS", 41, 'u', 1},
	{"SMF_RG_PERC_MIN", 44, 'u', 4},
	{"SMF_RG_PERC_MAX", 48, 'u', 4},
	{"SMF99_RG_MEM_LIMIT", 56, 'u', 4},
	{"SMF99_RG_LACS", 108, 'u', 4},
	{"SMF99_RG_SUsIFA", 112, 'u', 8},
	{"SMF99_RG_SUsSUP", 120, 'u', 8},
	{"SMF99_RG_MEMSMPCNT", 132, 'u', 4},
	{"SMF99_RG_FRAMECNT", 136, 'u', 8},
}

var smf099GrFields = fields{
	{"SMF99_GR_SYSNAME", 0, 'c', 8},
	{"SMF99_GR_TSO_SESSIONS_ROUTED", 8, 'u', 4},
	{"SMF99_GR_NONTSO_SESSIONS_ROUTED", 12, 'u', 4},
	{"SMF99_GR_TSO_AVG_COST", 16, 'u', 4},
	{"SMF99_GR_TSO_PI", 20, 'u', 4},
	{"SMF99_GR_FLAGS", 24, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE0", 28, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE1", 32, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE2", 36, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE3", 40, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE4", 44, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE5", 48, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE6", 52, 'u', 4},
	{"SMF99_GR_SERVICE_BY_IMPORTANCE7", 56, 'u', 4},
}

var smf099SlFields = fields{
	{"SMF99_SLConfigFlags", 0, 'u', 1},
	{"SMF99_SLStateFlags", 1, 'u', 1},
	{"SMF99_SLImgCapacity", 4, 'u', 4},
	{"SMF99_SLCecCapacity", 8, 'u', 4},
	{"SMF99_SLCecCpuCount", 12, 'u', 2},
	{"SMF99_SLLogicalCpuCount", 14, 'u', 2},
	{"SMF99_SLCecServiceUnitsPerSecToShare", 16, 'u', 4},
	{"SMF99_SLImgMsuAtCurrentWeight", 20, 'u', 4},
	{"SMF99_SLAvgMsu", 28, 'u', 4},
	{"SMF99_SLAvgMsuCapped", 32, 'u', 4},
	{"SMF99_SLAvgMsuUncapped", 36, 'u', 4},
	{"SMF99_SLIntervalService", 40, 'u', 4},
	{"SMF99_SLIntervalTime", 44, 'u', 4},
	{"SMF99_SLRollInterval", 52, 'u', 2},
	{"SMF99_SLServiceTableIntervals", 54, 'u', 2},
	{"SMF99_SLIntervalsToCap", 56, 'u', 2},
	{"SMF99_SLIntervalsToUncap", 58, 'u', 2},
	{"SMF99_SLPatternIntervalCount", 60, 'u', 2},
	{"SMF99_SL_Query_Response_Code", 64, 'u', 4},
	{"SMF99_SL_Setcap_Response_Code", 68, 'u', 4},
	{"SMF99_SLUnusedScaling", 80, 'u', 1},
	{"SMF99_HardwareGroupname", 84, 'c', 8},
	{"SMF99_HardwareGroupCpuTypeCapCp", 92, 'u', 4},
	{"SMF99_HardwareGroupCpuTypeCapZaap", 96, 'u', 4},
	{"SMF99_HardwareGroupCpuTypeCapZiip", 100, 'u', 4},
	{"SMF99_RTCapLeadTime", 104, 'u', 2},
	{"SMF99_Time_To_Cap", 106, 'u', 2},
	{"SMF99_Time_To_Cap_Group", 108, 'u', 2},
	{"SMF99_SUS_ZIIP", 184, 'u', 8},
	{"SMF99_SUS_ZIIP_ON_CP", 192, 'u', 8},
	{"SMF99_SUS_JAVA_ON_ZIIP", 200, 'u', 8},
	{"SMF99_SUS_JAVA_ON_CP", 208, 'u', 8},
}

var smf099ClassFields = fields{
	{"SMF99_CNAM", 0, 'c', 8},
	{"SMF99_CGRN", 8, 'c', 8},
	{"SMF99_CNUMP", 16, 'u', 2},
	{"SMF99_CINDEX", 18, 'u', 2},
}

var smf099PeriodFields = fields{
	{"SMF99_PCNM", 0, 'c', 8},
	{"SMF99_PNUM", 8, 'u', 2},
	{"SMF99_PGOALTYP", 10, 'u', 1},
	{"SMF99_PGOALVAL", 12, 'u', 4},
	{"SMF99_PIMPOR", 16, 'u', 2},
	{"SMF99_PBDP", 18, 'u', 1},
	{"SMF99_PMPLI", 20, 'u', 2},
	{"SMF99_PMPLO", 22, 'u', 2},
	{"SMF99_PAMTA", 24, 'u', 4},
	{"SMF99_PRUA", 28, 'u', 4},
	{"SMF99_PLRUA", 32, 'u', 4},
	{"SMF99_PPSPT", 36, 'u', 4},
	{"SMF99_PPSITAR", 40, 'u', 4},
	{"SMF99_PESPOL", 44, 'u', 1},
	{"SMF99_PESVIO", 45, 'u', 1},
	{"SMF99_PESHSP", 46, 'u', 1},
	{"SMF99_PESSWAP", 47, 'u', 1},
	{"SMF99_PPROT", 48, 'u', 2},
	{"SMF99_PLRU", 50, 'u', 2},
	{"SMF99_PSPAV", 52, 'u', 2},
	{"SMF99_PVIOL", 54, 'u', 2},
	{"SMF99_PVIOS", 56, 'u', 2},
	{"SMF99_PHSPL", 58, 'u', 2},
	{"SMF99_PHSPS", 60, 'u', 2},
	{"SMF99_PESCS", 62, 'u', 2},
	{"SMF99_PLPI", 64, 'u', 4},
	{"SMF99_PSPI", 68, 'u', 4},
	{"SMF99_PSERV", 72, 'u', 4},
	{"SMF99_PMDP", 76, 'u', 4},
	{"SMF99_PLCPUU", 80, 'u', 4},
	{"SMF99_PLCPUD", 84, 'u', 4},
	{"SMF99_PMTTWA", 88, 'u', 4},
	{"SMF99_PADP", 92, 'u', 4},
	{"SMF99_PASERC", 96, 'u', 4},
	{"SMF99_PPRSER", 100, 'u', 4},
	{"SMF99_PIDLE", 104, 'u', 4},
	{"SMF99_POTHR", 108, 'u', 4},
	{"SMF99_PCPUU", 112, 'u', 4},
	{"SMF99_PCPUD", 116, 'u', 4},
	{"SMF99_PAUXP", 120, 'u', 4},
	{"SMF99_PAUXC", 124, 'u', 4},
	{"SMF99_PVIO", 128, 'u', 4},
	{"SMF99_PHSS", 132, 'u', 4},
	{"SMF99_PHSC", 136, 'u', 4},
	{"SMF99_PASWP", 140, 'u', 4},
	{"SMF99_PMPLD", 144, 'u', 4},
	{"SMF99_PCAPD", 148, 'u', 4},
	{"SMF99_PXMO", 152, 'u', 4},
	{"SMF99_PXMEM_OF", 156, 'u', 4},
	{"SMF99_PXMEM_LN", 160, 'u', 2},
	{"SMF99_PXMEM_ON", 162, 'u', 2},
	{"SMF99_PSERV_OF", 164, 'u', 4},
	{"SMF99_PSERV_LN", 168, 'u', 2},
	{"SMF99_PSERV_ON", 170, 'u', 2},
	{"SMF99_PESP_OF", 172, 'u', 4},
	{"SMF99_PESP_LN", 176, 'u', 2},
	{"SMF99_PESP_ON", 178, 'u', 2},
	{"SMF99_PCDCLOCK", 180, 'u', 2},
	{"SMF99_PNH", 182, 'u', 1},
	{"SMF99_PRTP", 183, 'u', 1},
	{"SMF99_PAUXS", 184, 'u', 4},
	{"SMF99_PIOU", 188, 'u', 4},
	{"SMF99_PIOD", 192, 'u', 4},
	{"SMF99_PIO_MDP", 196, 'u', 4},
	{"SMF99_PIODP", 200, 'u', 1},
	{"SMF99_FLAGS", 201, 'u', 1},
	{"SMF99_Per_IO_Mgmt_Support_Data", 202, 'u', 2},
	{"SMF99_PDEVCL", 204, 'u', 4},
	{"SMF99_PSERVER_TYPE", 208, 'u', 4},
	{"SMF99_PSDATA_OF", 212, 'u', 4},
	{"SMF99_PSDATA_LN", 216, 'u', 2},
	{"SMF99_PSDATA_ON", 218, 'u', 2},
	{"SMF99_PQDATA_OF", 220, 'u', 4},
	{"SMF99_PQDATA_LN", 224, 'u', 2},
	{"SMF99_PQDATA_ON", 226, 'u', 2},
	{"SMF99_PAVG_SIZE", 228, 'u', 4},
	{"SMF99_PGRN", 232, 'c', 8},
	{"SMF99_PSYS_CPUU", 240, 'u', 4},
	{"SMF99_PSYS_NONIDLE", 244, 'u', 4},
	{"SMF99_PSYS_IDLE", 248, 'u', 4},
	{"SMF99_PSYS_OTHER", 252, 'u', 4},
	{"SMF99_IOSUBSAMOF", 256, 'u', 4},
	{"SMF99_IOSUBSAMLN", 260, 'u', 2},
	{"SMF99_IOSUBSAMON", 262, 'u', 2},
	{"SMF99_SPMDP", 264, 'u', 4},
	{"SMF99_AVG_NUM_TASKS", 268, 'u', 4},
	{"SMF99_PBPD", 272, 'u', 4},
	{"SMF99_SWCT", 276, 'u', 4},
	{"SMF99_FLAGS2", 280, 'u', 2},
	{"SMF99_NUM_SAMP_HIST_ROWS_USED", 282, 'u', 2},
	{"SMF99_CADP", 284, 'u', 4},
	{"SMF99_SBCPUU", 288, 'u', 4},
	{"SMF99_SBCPUD", 292, 'u', 4},
	{"SMF99_PSYS_IO_DLY", 296, 'u', 4},
	{"SMF99_PSYS_NON_IO_DLY", 300, 'u', 4},
	{"SMF99_PIFAU", 304, 'u', 4},
	{"SMF99_PIFAD", 308, 'u', 4},
	{"SMF99_PISERV", 312, 'u', 4},
	{"SMF99_PIMDP", 316, 'u', 4},
	{"SMF99_PIMTTWA", 320, 'u', 4},
	{"SFM99_PIADP", 324, 'u', 4},
	{"SMF99_PIASERC", 328, 'u', 4},
	{"SMF99_PIPRSER", 332, 'u', 4},
	{"SMF99_ICADP", 336, 'u', 4},
	{"SMF99_PIFAONCP", 340, 'u', 4},
	{"SMF99_PLIFAU", 344, 'u', 4},
	{"SMF99_PLIFAD", 348, 'u', 4},
	{"SMF99_PSUPU", 352, 'u', 4},
	{"SMF99_PSUPD", 356, 'u', 4},
	{"SMF99_PSUPONCP", 360, 'u', 4},
	{"SMF99_PLSUPU", 364, 'u', 4},
	{"SMF99_PLSUPD", 368, 'u', 4},
	{"SMF99_PSSERV", 372, 'u', 4},
	{"SMF99_Time_at_PDP_Using", 376, 'u', 4},
	{"SMF99_Time_at_PDP", 380, 'u', 4},
	{"SMF99_EWLM_LOCAL_PI", 384, 'u', 4},
	{"SMF99_EWLM_GLOBAL_PI", 388, 'u', 4},
	{"SMF99_PSMDP", 392, 'u', 4},
	{"SMF99_PSMTTWA", 396, 'u', 4},
	{"SFM99_PSADP", 400, 'u', 4},
	{"SMF99_PSASERC", 404, 'u', 4},
	{"SMF99_PSPRSER", 408, 'u', 4},
	{"SMF99_SCADP", 412, 'u', 4},
	{"SMF99_HdLockPromotion_Time_at_PDP", 416, 'u', 4},
	{"SMF99_HdLock_Time_at_PDP_Using", 420, 'u', 4},
	{"SMF99_PNS_PSERV", 424, 'u', 4},
	{"SMF99_PNS_IPSERV", 428, 'u', 4},
	{"SMF99_PNS_SPSERV", 432, 'u', 4},
	{"SMF99_VARTIME_AT_PDP", 436, 'u', 4},
	{"SMF99_VARTIME_AT_PDP_USING", 440, 'u', 4},
	{"SMF99_VARWEIGHTED_TIME_AT_PDP", 444, 'u', 4},
	{"SMF99_RT_DISTRI_MID_POINT", 448, 'u', 4},
	{"SMF99_RT_DISTRI_TIME_STAMP", 452, 't', 8}, // STCK, though the manual says EBCDIC
	{"SMF99_RT_DISTRI_AVG_RESP_TIME", 460, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_ROWS", 464, 'u', 4},
	{"SMF99_RT_DISTRI_SUM_TRANS", 468, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS01", 472, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS02", 476, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS03", 480, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS04", 484, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS05", 488, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS06", 492, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS07", 496, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS08", 500, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS09", 504, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS10", 508, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS11", 512, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS12", 516, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS13", 520, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS14", 524, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS15", 528, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS16", 532, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS17", 536, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS18", 540, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS19", 544, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS20", 548, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS21", 552, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS22", 556, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS23", 560, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS24", 564, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS25", 568, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS26", 572, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS27", 576, 'u', 4},
	{"SMF99_RT_DISTRI_NUM_TRANS28", 580, 'u', 4},
	{"SMF99_RT_DISTRI_RUNNING_COUNT", 584, 'u', 4},
	{"SMF99_RT_DISTRI_INT_SUM_TRANS", 588, 'u', 4},
	{"SMF99_GAV_ARRAY", 592, 'x', 6},
	{"SMF99_VELO_USING_STATES", 600, 'u', 4},
	{"SMF99_USING_IOSM_CURR", 604, 'u', 4},
	{"SMF99_DELAY_IOSM_CURR", 608, 'u', 4},
	{"SMF99_DISC_IOSM_CURR", 612, 'u', 4},
	{"SMF99_DISC_IOSM_ACCUM", 616, 'u', 4},
	{"SMF99_CUQT_IOSM_CURR", 620, 'u', 4},
	{"SMF99_CUQT_IOSM_ACCUM", 624, 'u', 4},
	{"Smf99_Thro_IOSM_curr", 628, 'u', 4},
	{"Smf99_Thro_IOSM_accum", 632, 'u', 4},
	{"Smf99_CntD_IOSM_curr", 636, 'u', 4},
	{"Smf99_CntD_IOSM_accum", 640, 'u', 4},
	{"Smf99_Pend_IOSM_curr", 644, 'u', 4},
	{"Smf99_Pend_IOSM_accum", 648, 'u', 4},
	// SMF99_SpecRpt(4): all, mobile, category A and category B transactions
	{"SMF99_SpecCP_All", 656, 'u', 8},
	{"SMF99_SpecOffload_All", 664, 'u', 8},
	{"SMF99_SpecOffloadOnCP_All", 672, 'u', 8},
	{"SMF99_SpecCP_Mobile", 680, 'u', 8},
	{"SMF99_SpecOffload_Mobile", 688, 'u', 8},
	{"SMF99_SpecOffloadOnCP_Mobile", 696, 'u', 8},
	{"SMF99_SpecCP_CatA", 704, 'u', 8},
	{"SMF99_SpecOffload_CatA", 712, 'u', 8},
	{"SMF99_SpecOffloadOnCP_CatA", 720, 'u', 8},
	{"SMF99_SpecCP_CatB", 728, 'u', 8},
	{"SMF99_SpecOffload_CatB", 736, 'u', 8},
	{"SMF99_SpecOffloadOnCP_CatB", 744, 'u', 8},
	{"SMF99_TRG_NAME", 768, 'c', 8},
	{"SMF99_AI_OF", 776, 'u', 4},
	{"SMF99_AI_LN", 780, 'u', 2},
	{"SMF99_AI_ON", 782, 'u', 2},
}

var smf099ServerFields = fields{
	{"SMF99_SERVER_CNM", 0, 'c', 8},
	{"SMF99_SERVER_PNUM", 8, 'u', 4},
	{"SMF99_SERVER_OBS", 12, 'u', 4},
}

var smf099SdataFields = fields{
	{"SMF99_SDATA_WQDEL", 0, 'u', 4},
	{"SMF99_SDATA_ENC_AUXP", 4, 'u', 4},
	{"SMF99_SDATA_ENC_VIO", 8, 'u', 4},
	{"SMF99_SDATA_ENC_HSP", 12, 'u', 4},
	{"SMF99_SDATA_ENC_MPLD", 16, 'u', 4},
	{"SMF99_SDATA_ENC_ASWP", 20, 'u', 4},
	{"SMF99_SDATA_SERVER_CLASS_NAME", 24, 'c', 8},
	{"SMF99_SDATA_SERVER_TYPE", 32, 'u', 1},
	{"SMF99_SDATA_SUBSYS_TYPE", 36, 'c', 4},
	{"SMF99_SDATA_SUBSYS_NAME", 40, 'c', 8},
}

var smf099EspFields = fields{
	{"SMF99_AS_ESP_ANAM", 0, 'c', 8},
	{"SMF99_AS_ESP_AP", 8, 'u', 1},
	{"SMF99_AS_ESP_VP", 9, 'u', 1},
	{"SMF99_AS_ESP_HP", 10, 'u', 1},
	{"SMF99_AS_ESP_ASID", 11, 'u', 2},
	{"SMF99_AS_ESP_FLAGS", 13, 'u', 1},
	{"SMF99_AS_ESP_CS_FMCT", 16, 'u', 4},
	{"SMF99_AS_ESP_ES_FMCT", 20, 'u', 4},
	{"SMF99_AS_ESP_PPS_TAR", 24, 'u', 4},
	{"SMF99_AS_CPSRP_SAMP", 28, 'u', 2},
	{"SMP99_AS_CPSRP_CUR_FP_SAMP", 30, 'u', 2},
	{"SMP99_AS_CPSRP_PREV_FP_SAMP", 32, 'u', 2},
	{"SMF99_AS_HealthInd", 34, 'u', 1},
	{"SMF99_AS_TOTAL_SERVICE", 36, 'u', 4},
	{"SMF99_AS_CPU_SERVICE", 40, 'u', 4},
	{"SMF99_AS_SRB_SERVICE", 44, 'u', 4},
	{"SMF99_AS_MSO_SERVICE", 48, 'u', 4},
	{"SMF99_AS_TRN_SERVICE", 52, 'u', 4},
	{"SMF99_AS_IO_SERVICE", 56, 'u', 4},
	{"SMF99_AS_DISP_COUNT", 60, 'u', 2},
	{"SMF99_AS_IFA_SERVICE", 64, 'u', 8},
	{"SMF99_AS_IFACP_SERVICE", 72, 'u', 8},
	{"SMF99_AS_SUP_SERVICE", 80, 'u', 8},
	{"SMF99_AS_SUPCP_SERVICE", 88, 'u', 8},
	{"SMF99_AS_PB_SERVICE", 96, 'u', 8},
	{"SMF99_AS_PB_OFFLOAD_SERVICE", 104, 'u', 8},
	{"SMF99_AS_PB_OFFLOADONCP_SERVICE", 112, 'u', 8},
	{"SMF99_AS_ENCLAVE_TIME", 120, 'u', 4},
	{"SMF99_AS_ENCLAVE_CPU_SERVICE", 124, 'u', 4},
	{"SMF99_AS_ENCLAVE_IFA_TIME", 128, 'u', 8},
	{"SMF99_AS_ENCLAVE_IFACP_TIME", 136, 'u', 8},
	{"SMF99_AS_ENCLAVE_SUP_TIME", 144, 'u', 8},
	{"SMF99_AS_ENCLAVE_SUPCP_TIME", 152, 'u', 8},
	{"SMF99_AS_BA_BrkLoc", 160, 'x', 18},
	{"SMF99_AS_BA_Loc", 178, 'x', 18},
	{"SMF99_AS_BA_Flag_Area", 196, 'u', 3},
	{"SMF_AS_BA_Mem_Score", 200, 'u', 8},
	{"SMF_AS_MP_Actual_Pool", 208, 'x', 8},
	{"SMF99_AS_MP_Requested_Pool", 216, 'x', 8},
	{"SMF99_AS_TRC", 224, 'c', 8},
	{"SMF99_AS_TRG", 232, 'c', 8},
}

var smf099IoPrtyFields = fields{
	{"SMF99_IPTPRTY", 0, 'u', 2},
	{"SMF99_IPTNP", 2, 'u', 2},
	{"SMF99_IPTIMDP", 4, 'u', 4},
	{"SMF99_IPTPMDP", 8, 'u', 4},
	{"SMF99_IPTW2UR", 12, 'u', 4},
	{"SMF99_IPTIMAXD", 16, 'u', 4},
	{"SMF99_IPTWMAXD", 20, 'u', 4},
}

var smf099S6PeriodFields = fields{
	{"SMF996_ECLASS_NAME", 0, 'c', 8},
	{"SMF996_PER_NUM", 8, 'u', 2},
	{"SMF996_GOALTYPE", 10, 'u', 1},
	{"SMF996_PERCENTILE", 11, 'u', 1},
	{"SMF996_ICLASS_NAME", 12, 'c', 8},
	{"SMF996_GOALVAL", 20, 'u', 4},
	{"SMF996_IMPOR", 24, 'u', 2},
	{"SMF996_DP", 26, 'u', 1},
	{"SMF996_IODP", 27, 'u', 1},
	{"SMF996_MPLI", 28, 'u', 2},
	{"SMF996_MPLO", 30, 'u', 2},
	{"SMF996_RUA", 32, 'u', 4},
	{"SMF996_PSPT", 36, 'u', 4},
	{"SMF996_PSITAR", 40, 'u', 4},
	{"SMF996_LOCAL_PI", 44, 'u', 4},
	{"SMF996_SYSPLEX_PI", 48, 'u', 4},
	{"SMF996_SERVER_DATA_OF", 52, 'u', 4},
	{"SMF996_SERVER_DATA_LN", 56, 'u', 2},
	{"SMF996_SERVER_DATA_ON", 58, 'u', 2},
	{"SMF996_PSERV", 60, 'u', 4},
	{"SMF996_PISERV", 64, 'u', 4},
	{"SMF996_PSSERV", 68, 'u', 4},
	{"SMF996_Time_at_PDP_Using", 72, 'u', 4},
	{"SMF996_Time_at_PDP", 76, 'u', 4},
	{"SMF996_FLAGS", 80, 'u', 1},
	{"SMF996_EWLM_LOCAL_PI", 84, 'u', 4},
	{"SMF996_EWLM_GLOBAL_PI", 88, 'u', 4},
}

var smf099S6ServerFields = fields{
	{"SMF99_S6_SERVER_CLASS_NAME", 0, 'c', 8},
	{"SMF99_S6_SERVER_PER_NUM", 8, 'u', 4},
}

var smf099LparFields = fields{
	{"SMF998_ImageSystemName", 0, 'c', 8},
	{"SMF998_ControlFlag1", 8, 'u', 1},
	{"SMF998_SystemNumber", 9, 'u', 1},
	{"SMF998_LastSetTime", 12, 't', 8},
	{"SMF998_TotalWeight", 20, 'u', 4},
	{"SMF998_NumberOfSharedPhysicalCPUs", 24, 'u', 1},
	{"SMF998_ImageID", 28, 'u', 1},
	{"SMF998_ImageFlags", 29, 'u', 1},
	{"SMF998_NumberCpusActive", 32, 'u', 2},
	{"SMF998_AverageCpuUtilization", 34, 'u', 2},
	{"SMF998_ImageInitialWeight", 36, 'u', 2},
	{"SMF998_ImageCurrentWeight", 38, 'u', 2},
	{"SMF998_ImageMinimumWeight", 40, 'u', 2},
	{"SMF998_ImageMaximumWeight", 42, 'u', 2},
	{"SMF998_Pro_Time_Avail", 44, 'u', 4},
	{"SMF998_Service_Units_Per_Second", 48, 'u', 4},
	{"SMF998_SoftCapMsu", 68, 'u', 4},
	{"SMF998_PricingManagementWeight", 76, 'u', 4},
}

var smf099S8PrtyFields = fields{
	{"SMF998_PTPRTY", 0, 'u', 2},
	{"SMF998_PTNP", 2, 'u', 2},
	{"SMF998_PTIMDP", 4, 'u', 4},
	{"SMF998_PTPMDP", 8, 'u', 4},
	{"SMF998_PTCPUU", 12, 'u', 4},
	{"SMF998_PTCPUD", 16, 'u', 4},
	{"SMF998_PTW2UR", 20, 'u', 4},
	{"SMF998_PTAPU", 24, 'u', 4},
	{"SMF998_PTPPU", 28, 'u', 4},
	{"SMF998_PTACMD", 32, 'u', 4},
	{"SMF998_PTIMAXD", 40, 'u', 4},
	{"SMF998_PTWMAXD", 44, 'u', 4},
	{"SMF998_PTIAMTW", 48, 'u', 4},
	{"SMF998_PTWAMTW", 52, 'u', 4},
}

var smf099S8PeriodFields = fields{
	{"SMF998_Service_Class_Name", 0, 'c', 8},
	{"SMF998_Importance", 8, 'u', 2},
	{"SMF998_Period_Number", 10, 'u', 2},
	{"SMF998_Dasd_Io_Delay_Sample", 16, 'u', 4},
	{"SMF998_Non_Idle_Samp", 20, 'u', 4},
	{"SMF998_Cpu_Using_Sample", 24, 'u', 4},
	{"SMF998_Cpu_Delay_Sample", 28, 'u', 4},
	{"SMF998_Wlm_Queue_Delay_Sample", 32, 'u', 4},
	{"SMF998_Dasd_Io_Using_Sample", 36, 'u', 4},
	{"SMF998_Max_Dem_Per", 40, 'u', 4},
	{"SMF998_Pi_Achieved", 44, 'u', 4},
	{"SMF998_Sysplex_Pi_Achieved", 48, 'u', 4},
	{"SMF998_Service_Ow", 52, 'u', 4},
	{"SMF998_Mtw_Adj", 56, 'u', 4},
	{"SMF998_Base_Priority", 60, 'u', 2},
	{"SMF998_Cap_Num_Slices", 62, 'u', 2},
	{"SMF998_Work_Pro_Used", 64, 'u', 4},
	{"SMF998_Current_Ach_Dem_Per", 68, 'u', 4},
	{"SMF998_Ach_Dem_Per", 72, 'u', 4},
	{"SMF998_Old_Work_Pro_Used", 76, 'u', 4},
	{"SMF998_Proj_Pi_Com", 80, 'u', 4},
	{"SMF998_Using_Delta", 84, 'u', 4},
	{"SMF998_LparMgmt_Delay_Delta", 88, 'u', 4},
	{"SMF998_Cpu_Cap_Delay_Sample", 92, 'u', 4},
	{"SMF998_Iosub_Samples_Data_Offset", 96, 'u', 4},
	{"SMF998_Iosub_Samples_Data_Length", 100, 'u', 2},
	{"SMF998_Iosub_Samples_Data_number", 102, 'u', 2},
	{"SMF998_Sysplex_Proj_Pi_Com", 104, 'u', 4},
	{"SMF998_PC_CSS_NUMBER", 108, 'u', 1},
}
