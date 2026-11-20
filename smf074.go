package main

import (
	"database/sql"
	"sync"
)

// Type 74 (RMF activity of several resources). smf074 holds the header and RMF product section of every
// record; the sections of subtypes 1 (device activity), 2 (XCF), 3 (OMVS kernel), 5 (cache subsystem) and
// 6 (HFS) go to child tables with ID and SEQ (see rmfSection). Not decoded, as the sample data has none:
// subtypes 4 and 7 to 10, the XCF path sections, the cache device extension sections and the HFS file
// system sections.
//
// The RAID rank/extent pool section (smf074_raid) has two layouts at the same offsets, chosen by R7451FLG:
// 1 is RAID rank data, 2 extent pool and physical storage data. The table has the columns of both, and only
// the record's layout is filled, as for the channel measurement groups of type 73.
// The field tables were generated from the manual with tools/genfields.py.

var pStmt074, pStmt074Raid *sql.Stmt
var once074 sync.Once
var count074 int

var smf074Sections = []*rmfSection{
	{subtype: 1, at: 32, table: "smf074_devctl", f: smf074DevCtlFields},       // device control
	{subtype: 1, at: 40, table: "smf074_dev", f: smf074DevFields},             // device data
	{subtype: 2, at: 32, table: "smf074_xcfctl", f: smf074XcfCtlFields},       // XCF control
	{subtype: 2, at: 40, table: "smf074_xcfsys", f: smf074XcfSysFields},       // XCF system
	{subtype: 2, at: 56, table: "smf074_xcfmbr", f: smf074XcfMbrFields},       // XCF member
	{subtype: 3, at: 32, table: "smf074_omvs", f: smf074OmvsFields},           // OMVS kernel control
	{subtype: 5, at: 32, table: "smf074_cachectl", f: smf074CacheCtlFields},   // cache control
	{subtype: 5, at: 40, table: "smf074_cachedev", f: smf074CacheDevFields},   // cache device
	{subtype: 5, at: 56, table: "smf074_cachestat", f: smf074CacheStatFields}, // cache control unit status
	{subtype: 6, at: 32, table: "smf074_hfsglobal", f: smf074HfsGlobalFields}, // HFS global
	{subtype: 6, at: 40, table: "smf074_hfsbuf", f: smf074HfsBufFields},       // HFS global buffer
}

func smf074Init() {
	pStmt074 = createTable("smf074", rmfColumns("SMF74"))
	rmfCreateSections(smf074Sections)
	pStmt074Raid = createTable("smf074_raid", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf074RaidFields.columns()+
		", "+smf074RaidRankFields.columns()+", "+smf074RaidPoolFields.columns())
}

func smf074Parse(rec []byte) {
	once074.Do(smf074Init)

	count074 += 1
	insert(pStmt074, rmfRow(rec, count074, "SMF74")...)
	rmfInsertSections(smf074Sections, rec, count074)
	if section(rec).num(18, 2) != 5 {
		return
	}
	off, length, n := tripletAt(rec, 64) // RAID rank/extent pool
	for k := 0; k < n; k++ {
		r := section(rec).sub(off+k*length, length)
		var rank, pool section // only the layout named by R7451FLG is decoded; the other stays NULL
		switch r.num(8, 1) {
		case 1:
			rank = r
		case 2:
			pool = r
		}
		row := append([]interface{}{count074, k + 1}, smf074RaidFields.values(r)...)
		row = append(row, smf074RaidRankFields.values(rank)...)
		insert(pStmt074Raid, append(row, smf074RaidPoolFields.values(pool)...)...)
	}
}

var smf074DevCtlFields = fields{
	{"SMF74NXT", 0, 'u', 2},
	{"SMF74TOT", 2, 'u', 2},
	{"SMF74GEN", 4, 'u', 2},
	{"SMF74SUB", 6, 'u', 2},
	{"SMF74DCF", 8, 'u', 1},
	{"SMF74DMS", 9, 'u', 1},
	{"SMF74ENF", 10, 'u', 1},
	{"SMF74SMF", 11, 'u', 1},
	{"SMF74S15", 12, 'u', 4},
	{"SMF74SRC", 16, 'u', 4},
	{"SMF74SRS", 20, 'u', 4},
	{"SMF74TSR", 24, 'u', 2},
	{"SMF74CFL", 26, 'u', 1},
	{"SMF74LSN", 27, 'u', 1},
	{"SMF74TNM", 30, 'c', 44},
	{"SMF74TSF", 74, 'c', 2},
	{"SMF74TOK", 78, 'c', 16},
	// skipped SMF74TDT at 78: redefines an earlier field
	// skipped SMF74TTM at 86: redefines an earlier field
	{"SMF74MCT", 96, 'u', 4},
	{"SMF74TDY", 100, 'c', 10},
}

var smf074DevFields = fields{
	{"SMF74NUM", 0, 'u', 2},
	{"SMF74LCU", 2, 'u', 2},
	{"SMF74CNF", 5, 'u', 1},
	{"SMF74SER", 6, 'c', 6},
	{"SMF74TYP", 12, 'u', 4},
	{"SMF74NUX", 16, 'u', 4},
	{"SMF74SSC", 20, 'u', 4},
	{"SMF74MEC", 24, 'u', 4},
	{"SMF74CNN", 28, 'u', 4},
	{"SMF74PEN", 32, 'u', 4},
	{"SMF74ATV", 36, 'u', 4},
	{"SMF74DIS", 40, 'u', 4},
	{"SMF74QUE", 44, 'u', 4},
	{"SMF74UTL", 48, 'u', 4},
	{"SMF74RSV", 52, 'u', 4},
	{"SMF74DSO", 56, 'u', 4},
	{"SMF74ALC", 60, 'u', 4},
	{"SMF74MTP", 64, 'u', 4},
	{"SMF74NRD", 68, 'u', 4},
	{"SMF74COF", 72, 'u', 2},
	{"SMF74ICT", 74, 'u', 2},
	{"SMF74DVB", 76, 'u', 4},
	{"SMF74CLF", 84, 'u', 1},
	{"SMF74SGN", 88, 'c', 8},
	{"SMF74NDA", 96, 'u', 4},
	{"SMF74DEV", 100, 'c', 8},
	{"SMF74CU", 108, 'c', 8},
	{"SMF74CNX", 120, 'u', 1},
	{"SMF74CN2", 121, 'u', 1},
	{"SMF74MTC", 122, 'u', 2},
	{"SMF74DTS", 124, 'u', 1},
	{"SMF74DCT", 125, 'c', 28},
	{"SMF74HPC", 153, 'u', 1},
	{"SMF74NSS", 154, 'u', 2},
	{"SMF74PSM", 156, 'u', 4},
	{"SMF74PCT", 160, 'u', 4},
	{"SMF74CMR", 164, 'u', 4},
	{"SMF74CAP", 168, 'u', 4},
	{"SMF74IDT", 172, 'u', 4},
	{"SMF74CUQ", 176, 'u', 4},
	{"SMF74SCS", 181, 'u', 1},
	{"SMF74NM2", 182, 'u', 2},
	{"SMF74ATD", 184, 'u', 4},
	{"SMF74AGC", 191, 'u', 1},
	{"SMF74AGS", 192, 'u', 4},
	{"SMF74SBR", 200, 'f', 8},
	{"SMF74SBW", 208, 'f', 8},
	{"SMF74SQR", 216, 'f', 8},
	{"SMF74SQW", 224, 'f', 8},
	{"SMF74SPR", 232, 'f', 8},
	{"SMF74SPW", 240, 'f', 8},
	{"SMF74SFTR", 248, 'f', 8},
	{"SMF74SFTW", 256, 'f', 8},
	{"SMF74SLBR", 264, 'f', 4},
	{"SMF74SLBW", 268, 'f', 4},
	{"SMF74SCMR", 272, 'f', 4},
	{"SMF74SNIS", 276, 'f', 4},
	{"SMF74STOR", 280, 'f', 4},
	{"SMF74STOW", 284, 'f', 4},
	{"SMF74SOR", 288, 'f', 4},
	{"SMF74SOW", 292, 'f', 4},
	{"SMF74IOS", 296, 'f', 4},
}

var smf074XcfCtlFields = fields{
	{"R742TSR", 0, 'u', 2},
	{"R742STOT", 4, 'u', 4},
	{"R742SNXT", 8, 'u', 4},
	{"R742PTOT", 12, 'u', 4},
	{"R742PNXT", 16, 'u', 4},
	{"R742MTOT", 20, 'u', 4},
	{"R742MNXT", 24, 'u', 4},
}

var smf074XcfSysFields = fields{
	{"R742SNME", 0, 'c', 8},
	{"R742SSTF", 8, 'u', 1},
	{"R742SDIR", 9, 'u', 1},
	{"R742SPTH", 12, 'u', 4},
	{"R742SBSY", 16, 'u', 4},
	{"R742SNOP", 20, 'u', 4},
	{"R742SMXB", 24, 'u', 4},
	{"R742SBIG", 28, 'u', 4},
	{"R742SFIT", 32, 'u', 4},
	{"R742SSML", 36, 'u', 4},
	{"R742SOVR", 40, 'u', 4},
	{"R742STCL", 44, 'u', 4},
	{"R742STCN", 48, 'c', 8},
}

var smf074XcfMbrFields = fields{
	{"R742MSYS", 0, 'c', 8},
	{"R742MGRP", 8, 'c', 8},
	{"R742MMEM", 16, 'c', 16},
	{"R742MSTF", 32, 'u', 1},
	{"R742MST1", 33, 'u', 1},
	{"R742MST2", 34, 'u', 1},
	{"R742MSNT", 36, 'u', 4},
	{"R742MRCV", 40, 'u', 4},
	{"R742MINT", 44, 'u', 4},
	{"R742MJOB", 48, 'c', 8},
}

var smf074OmvsFields = fields{
	{"R743CYCU", 0, 'u', 4},
	{"R743CYCT", 4, 'u', 4},
	{"R743FLG", 8, 'u', 4},
	{"R743SYSC", 12, 'f', 4},
	{"R743SCMN", 16, 'u', 4},
	{"R743SCMX", 20, 'u', 4},
	{"R743CPU", 24, 'f', 4},
	{"R743CTMN", 28, 'u', 4},
	{"R743CTMX", 32, 'u', 4},
	{"R743OPR", 36, 'f', 4},
	{"R743OPMN", 40, 'u', 4},
	{"R743OPMX", 44, 'u', 4},
	{"R743OUS", 48, 'f', 4},
	{"R743OUMN", 52, 'u', 4},
	{"R743OUMX", 56, 'u', 4},
	{"R743OPRU", 60, 'f', 4},
	{"R743ORMN", 64, 'u', 4},
	{"R743ORMX", 68, 'u', 4},
	{"R743MAXP", 72, 'u', 2},
	{"R743MAXU", 74, 'u', 2},
	{"R743MXPU", 76, 'u', 2},
	{"R743CURP", 80, 'f', 4},
	{"R743CPMN", 84, 'u', 2},
	{"R743CPMX", 86, 'u', 2},
	{"R743CURU", 88, 'f', 4},
	{"R743CUMN", 92, 'u', 2},
	{"R743CUMX", 94, 'u', 2},
	{"R743MMSG", 96, 'u', 4},
	{"R743MSEM", 100, 'u', 4},
	{"R743MSHM", 104, 'u', 4},
	{"R743MSPG", 108, 'u', 4},
	{"R743CMSG", 112, 'f', 4},
	{"R743CMMN", 116, 'u', 4},
	{"R743CMMX", 120, 'u', 4},
	{"R743CSEM", 124, 'f', 4},
	{"R743CSMN", 128, 'u', 4},
	{"R743CSMX", 132, 'u', 4},
	{"R743CSHM", 136, 'f', 4},
	{"R743CHMN", 140, 'u', 4},
	{"R743CHMX", 144, 'u', 4},
	{"R743CSPG", 148, 'f', 4},
	{"R743CGMN", 152, 'u', 4},
	{"R743CGMX", 156, 'u', 4},
	{"R743OMSG", 160, 'f', 4},
	{"R743OMMN", 164, 'u', 4},
	{"R743OMMX", 168, 'u', 4},
	{"R743OSEM", 172, 'f', 4},
	{"R743OSMN", 176, 'u', 4},
	{"R743OSMX", 180, 'u', 4},
	{"R743OSHM", 184, 'f', 4},
	{"R743OHMN", 188, 'u', 4},
	{"R743OHMX", 192, 'u', 4},
	{"R743OSPG", 196, 'f', 4},
	{"R743OGMN", 200, 'u', 4},
	{"R743OGMX", 204, 'u', 4},
	{"R743MMAP", 208, 'u', 4},
	{"R743CMAP", 212, 'f', 4},
	{"R743CAMN", 216, 'u', 4},
	{"R743CAMX", 220, 'u', 4},
	{"R743OMAP", 224, 'f', 4},
	{"R743OAMN", 228, 'u', 4},
	{"R743OAMX", 232, 'u', 4},
	{"R743MPAG", 236, 'u', 4},
	{"R743CPAG", 240, 'f', 4},
	{"R743CXMN", 244, 'u', 4},
	{"R743CXMX", 248, 'u', 4},
	{"R743OPAG", 252, 'f', 4},
	{"R743OXMN", 256, 'u', 4},
	{"R743OXMX", 260, 'u', 4},
	{"R743MSLR", 264, 'u', 4},
	{"R743CSLR", 268, 'f', 4},
	{"R743CLMN", 272, 'u', 4},
	{"R743CLMX", 276, 'u', 4},
	{"R743OSLR", 280, 'f', 4},
	{"R743OLMN", 284, 'u', 4},
	{"R743OLMX", 288, 'u', 4},
	{"R743MQDS", 292, 'u', 4},
	{"R743OQDS", 296, 'f', 4},
	{"R743OQMN", 300, 'u', 4},
	{"R743OQMX", 304, 'u', 4},
}

var smf074CacheCtlFields = fields{
	{"R745CLVL", 0, 'u', 1},
	{"R745CMDL", 1, 'u', 1},
	{"R745CCNT", 2, 'u', 1},
	{"R745CUID", 3, 'c', 1},
	{"R745CSC", 4, 'u', 1},
	{"R745CAE", 5, 'u', 3},
	{"R745CRTN", 8, 'u', 2},
	{"R745CIOC", 10, 'u', 1},
	{"R745CINT", 12, 'u', 4},
	{"R745CCMT", 16, 'c', 28},
	{"R745CFDV", 45, 'u', 3},
}

var smf074CacheDevFields = fields{
	{"R745DVOL", 0, 'c', 6},
	{"R745DFL4", 6, 'u', 1},
	{"R745DCID", 7, 'c', 1},
	// skipped R745DCCU at 7: redefines an earlier field
	{"R745DUNT", 8, 'u', 3},
	{"R745DEVN", 12, 'u', 2},
	{"R745DFLG", 16, 'u', 1},
	{"R745DVID", 17, 'u', 1},
	{"R745DVS1", 18, 'u', 1},
	{"R745DVS2", 19, 'u', 1},
	{"R745DRCR", 20, 'f', 4},
	{"R745DCRH", 24, 'f', 4},
	{"R745DWRC", 28, 'f', 4},
	{"R745DWCH", 32, 'f', 4},
	{"R745DRSR", 36, 'f', 4},
	{"R745DRSH", 40, 'f', 4},
	{"R745DWSR", 44, 'f', 4},
	{"R745DWSH", 48, 'f', 4},
	{"R745DRNR", 52, 'f', 4},
	{"R745DNRH", 56, 'f', 4},
	{"R745DWNR", 60, 'f', 4},
	{"R745DWNH", 64, 'f', 4},
	{"R745DICL", 68, 'f', 4},
	{"R745DBCR", 72, 'f', 4},
	{"R745DTC", 76, 'f', 4},
	{"R745DNTD", 80, 'f', 4},
	{"R745DCTD", 84, 'f', 4},
	{"R745DFWB", 88, 'f', 4},
	{"R745DFWC", 92, 'f', 4},
	{"R745DFWS", 96, 'f', 4},
	{"R745DCRM", 100, 'f', 4},
	{"R745DSG2", 104, 'u', 1},
	{"R745DSID", 106, 'u', 2},
	{"R745DCWP", 108, 'f', 4},
	{"R745DKDW", 112, 'f', 4},
	{"R745DKDH", 116, 'f', 4},
	{"R745DFWR", 120, 'f', 4},
}

var smf074CacheStatFields = fields{
	{"R745SVOL", 0, 'c', 6},
	{"R745SUNT", 8, 'u', 3},
	{"R745SDEV", 12, 'u', 2},
	{"R745SLN", 14, 'u', 2},
	{"R745SFT", 16, 'u', 1},
	{"R745SDID", 17, 'u', 1},
	{"R745SNAD", 18, 'u', 1},
	{"R745SNSS", 19, 'u', 1},
	{"R745SCS", 20, 'u', 1},
	{"R745SVSS", 21, 'u', 1},
	{"R745SCLN", 22, 'u', 2},
	{"R745SCSF", 24, 'u', 2},
	{"R745SCNF", 26, 'u', 4},
	{"R745SAVL", 30, 'u', 4},
	{"R745SPIN", 34, 'u', 4},
	{"R745SOFF", 38, 'u', 4},
	{"R745SDS1", 42, 'u', 1},
	{"R745SDS2", 43, 'u', 1},
	{"R745SCNV", 44, 'u', 4},
	{"R745SPND", 48, 'u', 4},
	{"R745SG2", 52, 'u', 1},
	{"R745SGL", 53, 'u', 1},
	{"R745SSID", 54, 'u', 2},
}

var smf074RaidFields = fields{ // fields common to both layouts
	{"R7451DVN", 0, 'u', 2},
	{"R7451INC", 2, 'u', 1},
	{"R7451SCS", 3, 'u', 1},
	{"R7451RSV", 4, 'f', 4},
	{"R7451FLG", 8, 'u', 1},
	{"R7451AID", 9, 'u', 1},
	{"R7451RMR", 32, 'f', 4},
	{"R7451XSF", 36, 'f', 4},
	{"R7451XCW", 40, 'f', 4},
	{"R7451TSP", 44, 'f', 4},
	{"R7451NVS", 48, 'f', 4},
	{"R7451CT1", 60, 'f', 4},
	{"R7451CT2", 64, 'f', 4},
	{"R7451CT3", 68, 'f', 4},
	{"R7451CT4", 72, 'f', 4},
	{"R7451CT5", 76, 'f', 4},
	{"R7451CT6", 80, 'f', 4},
	{"R7451ZHL", 88, 'f', 4},
	{"R7451ZHH", 92, 'f', 4},
	{"R7451GSF", 96, 'f', 4},
	{"R7451GSS", 100, 'f', 4},
	{"R7451SRR", 104, 'f', 4},
	{"R7451SRH", 108, 'f', 4},
	{"R7451SWR", 112, 'f', 4},
	{"R7451SWH", 116, 'f', 4},
}

var smf074RaidRankFields = fields{ // R7451FLG 1: RAID rank data
	{"R7451RID", 10, 'u', 2},
	{"R7451HDD", 12, 'u', 1},
	{"R7451RTY", 13, 'u', 1},
	{"R7451HSS", 14, 'u', 2},
	{"R7451RRQ", 16, 'f', 4},
	{"R7451WRQ", 20, 'f', 4},
	{"R7451SR", 24, 'f', 4},
	{"R7451SW", 28, 'f', 4},
	{"R7451RRT", 52, 'f', 4},
	{"R7451WRT", 56, 'f', 4},
}

var smf074RaidPoolFields = fields{ // R7451FLG 2: extent pool and physical storage data
	{"R7451XID", 10, 'u', 2},
	{"R7452XTY", 12, 'u', 1},
	{"R7452XFL", 13, 'u', 1},
	{"R7452PRO", 16, 'f', 4},
	{"R7452PWO", 20, 'f', 4},
	{"R7452PBR", 24, 'f', 4},
	{"R7452PBW", 28, 'f', 4},
	{"R7452PRT", 52, 'f', 4},
	{"R7452PWT", 56, 'f', 4},
}

var smf074HfsGlobalFields = fields{
	{"R746GMXV", 0, 'u', 4},
	{"R746GUSV", 4, 'u', 4},
	{"R746GMNF", 8, 'u', 4},
	{"R746GUSF", 12, 'u', 4},
	{"R746GMC", 16, 'f', 8},
	{"R746GMNC", 24, 'f', 8},
	{"R746G1C", 32, 'f', 8},
	{"R746G1NC", 40, 'f', 8},
	{"R746GLRC", 48, 'u', 4},
	{"R746GLRS", 52, 'u', 4},
	{"R746GSRC", 56, 'u', 4},
	{"R746GSRS", 60, 'u', 4},
	{"R746GSFL", 64, 'u', 1},
}

var smf074HfsBufFields = fields{
	{"R746GSB", 0, 'u', 2},
	{"R746GNDS", 2, 'u', 2},
	{"R746GSBP", 4, 'u', 4},
	{"R746GSBF", 8, 'u', 4},
	{"R746GBF", 16, 'f', 8},
	{"R746GBNF", 24, 'f', 8},
}
