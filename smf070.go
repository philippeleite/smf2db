package main

import (
	"database/sql"
	"sync"
)

// Type 70 (RMF processor activity). smf070 holds the header and RMF product section of every record.
// Subtype 1 adds the CPU control (smf070_cpc), CPU data (smf070_cpu), ASID data area (smf070_asid),
// PR/SM partition (smf070_lpar) and logical processor (smf070_lcpu), and CPU identification
// (smf070_cpuid) sections; subtype 2 the ICSF services section (smf070_icsf). Not decoded, as they are
// not in the sample data: subtype 1 logical core and tenant resource group sections, subtype 2 CCA,
// accelerator and PKCS11 coprocessor sections.
//
// The field tables below were generated from the manual's field rows (pdftotext -raw); where the manual's
// decimal and hex offsets disagree, the one that continues the sequence was used.

var (
	pStmt070, pStmt070Cpc, pStmt070Cpu, pStmt070Asid, pStmt070Lpar *sql.Stmt
	pStmt070Lcpu, pStmt070CpuID, pStmt070Icsf                      *sql.Stmt
	once070                                                        sync.Once
	count070                                                       int
)

func smf070Init() {
	pStmt070 = createTable("smf070", rmfColumns("SMF70"))
	pStmt070Cpc = createTable("smf070_cpc", "ID INT NOT NULL, "+smf070CpcFields.columns())
	pStmt070Cpu = createTable("smf070_cpu", "ID INT NOT NULL, "+smf070CpuFields.columns())
	pStmt070Asid = createTable("smf070_asid", "ID INT NOT NULL, "+smf070AsidFields.columns())
	pStmt070Lpar = createTable("smf070_lpar", "ID INT NOT NULL, "+smf070LparFields.columns())
	pStmt070Lcpu = createTable("smf070_lcpu", "ID INT NOT NULL, SMF70LPN INT, "+smf070LcpuFields.columns())
	pStmt070CpuID = createTable("smf070_cpuid", "ID INT NOT NULL, SMF70CIX INT, "+smf070CpuIDFields.columns())
	pStmt070Icsf = createTable("smf070_icsf", "ID INT NOT NULL, "+smf070IcsfFields.columns())
}

func smf070Parse(rec []byte) {
	once070.Do(smf070Init)

	s := section(rec)
	count070 += 1
	id := count070
	insert(pStmt070, rmfRow(rec, id, "SMF70")...)

	// each calls f for each section located by the triplet at rec[at:], with its 1-based position
	each := func(at int, f func(k int, sec section)) {
		off, length, n := tripletAt(rec, at)
		for k := 0; k < n; k++ {
			f(k+1, s.sub(off+k*length, length))
		}
	}
	row := func(stmt *sql.Stmt, f fields, lead ...interface{}) func(int, section) {
		return func(_ int, sec section) {
			insert(stmt, append(append([]interface{}{id}, lead...), f.values(sec)...)...)
		}
	}

	switch s.num(18, 2) {
	case 1:
		each(32, row(pStmt070Cpc, smf070CpcFields))
		each(40, row(pStmt070Cpu, smf070CpuFields))
		each(48, row(pStmt070Asid, smf070AsidFields))
		// Each partition is followed, in the logical processor sections, by its SMF70BDN processors.
		var owners []interface{}
		each(56, func(k int, sec section) {
			for i := 0; i < sec.num(10, 2); i++ {
				owners = append(owners, sec.u8(8)) // SMF70LPN
			}
			row(pStmt070Lpar, smf070LparFields)(k, sec)
		})
		each(64, func(k int, sec section) {
			var owner interface{}
			if k <= len(owners) {
				owner = owners[k-1]
			}
			row(pStmt070Lcpu, smf070LcpuFields, owner)(k, sec)
		})
		each(72, func(k int, sec section) { row(pStmt070CpuID, smf070CpuIDFields, k)(k, sec) })
	case 2:
		each(48, row(pStmt070Icsf, smf070IcsfFields))
	}
}

var smf070CpcFields = fields{
	{"SMF70MOD", 0, 'u', 2},
	{"SMF70VER", 2, 'u', 1},
	{"SMF70BNP", 3, 'u', 1},
	{"SMF70INB", 4, 'u', 1},
	{"SMF70STF", 5, 'u', 1},
	{"SMF70GTS", 6, 'u', 2},
	{"SMF70MDL", 8, 'c', 16},
	{"SMF70DSA", 24, 'u', 2},
	{"SMF70IFA", 26, 'u', 2},
	{"SMF70CPA", 28, 'u', 4},
	{"SMF70WLA", 32, 'u', 4},
	{"SMF70LAC", 36, 'u', 4},
	{"SMF70HOF", 40, 's', 8}, // STCK offset
	{"SMF70HWM", 48, 'c', 16},
	{"SMF70SUP", 64, 'u', 2},
	{"SMF70GJT", 66, 't', 8}, // STCK time (manual says EBCDIC)
	{"SMF70POM", 74, 'c', 4},
	{"SMF70CSC", 78, 'c', 16},
	{"SMF70HHF", 94, 'u', 1},
	{"SMF70CR", 95, 'u', 1},
	{"SMF70PMI", 96, 'u', 4},
	{"SMF70PMU", 100, 'u', 4},
	{"SMF70PMW", 104, 'u', 4},
	{"SMF70PMP", 108, 'u', 4},
	{"SMF70PMT", 112, 'u', 2},
	{"SMF70PML", 114, 'u', 2},
	{"SMF70MPC", 116, 'c', 16},
	{"SMF70MTC", 132, 'c', 16},
	{"SMF70MCR", 148, 'u', 4},
	{"SMF70MPR", 152, 'u', 4},
	{"SMF70MTR", 156, 'u', 4},
	{"SMF70ZEP", 160, 'u', 4},
	{"SMF70ZER", 164, 'u', 8},
	{"SMF70ZEE", 172, 'u', 8},
	{"SMF70ZEC", 180, 'u', 8},
	{"SMF70NRM", 188, 'u', 4},
	{"SMF70GAU", 192, 'u', 4},
	{"SMF70ZEI", 196, 'u', 8},
	{"SMF70NCR", 204, 'u', 4},
	{"SMF70NPR", 208, 'u', 4},
	{"SMF70NTR", 212, 'u', 4},
	{"SMF70CAI", 216, 'u', 1},
	{"SMF70CCR", 217, 'u', 1},
	{"SMF70MCP", 218, 'u', 2},
	{"SMF70ICP", 220, 'u', 2},
	{"SMF70CCP", 222, 'u', 2},
	{"SMF70CPA_actual", 224, 'u', 4},
	{"SMF70CPA_scaling_factor", 228, 'u', 4},
	{"SMF70MCF", 232, 'u', 4},
	{"SMF70MCFS", 236, 'u', 4},
	{"SMF70MCFI", 240, 'u', 4},
	{"SMF70CF", 244, 'u', 4},
	{"SMF70CFS", 248, 'u', 4},
	{"SMF70CFI", 252, 'u', 4},
	{"SMF70ATD", 256, 'u', 4},
	{"SMF70ATDS", 260, 'u', 4},
	{"SMF70ATDI", 264, 'u', 4},
	{"SMF70LACM", 268, 'u', 4},
	{"SMF70LACA", 272, 'u', 4},
	{"SMF70LACB", 276, 'u', 4},
	{"SMF70ADJ", 280, 'u', 4},
	{"SMF70LACCR", 284, 'u', 4},
	{"SMF70MaxPU", 288, 'u', 2},
	{"SMF70OS_PRTCT", 290, 'u', 1},
	{"SMF70_IPL_TIME", 324, 't', 8},
	{"SMF70_TRG_M_CNT", 332, 'u', 4},
	{"SMF70CPC_TYPE", 340, 'u', 4},
	{"SMF70MDL_VAR", 344, 'c', 16},
	{"SMF70MVCR", 360, 'u', 4},
	{"SMF70NVCR", 364, 'u', 4},
	{"SMF70ZSU_on_zIIP", 368, 'u', 8},
	{"SMF70ZSU_on_CP", 376, 'u', 8},
	{"SMF70JSU_on_zIIP", 384, 'u', 8},
	{"SMF70JSU_on_CP", 392, 'u', 8},
	{"SMF70CPE_LO", 400, 'u', 2},
	{"SMF70CPE_HI", 402, 'u', 2},
}

var smf070CpuFields = fields{
	{"SMF70WAT", 0, 'u', 8},
	{"SMF70CID", 8, 'u', 2},
	{"SMF70CNF", 10, 'u', 1},
	{"SMF70SER", 12, 'x', 3}, // 6 hex digits, no sign
	{"SMF70TYP", 15, 'u', 1},
	{"SMF70SLH", 16, 'u', 4},
	{"SMF70TPI", 20, 'u', 4},
	{"SMF70VFS", 24, 'u', 4},
	{"SMF70V", 28, 'u', 1},
	{"SMF70PAT", 32, 'u', 8},
	{"SMF70TCB", 40, 'u', 8},
	{"SMF70SRB", 48, 'u', 8},
	{"SMF70NIO", 56, 'u', 8},
	{"SMF70SIG", 64, 'u', 8},
	{"SMF70WTD", 72, 'u', 8},
	{"SMF70WTS", 80, 'u', 4},
	{"SMF70WTU", 84, 'u', 4},
	{"SMF70WTI", 88, 'u', 4},
}

var smf070AsidFields = fields{
	{"SMF70RMN", 0, 'u', 2},
	{"SMF70RMM", 2, 'u', 2},
	{"SMF70RTT", 4, 'u', 4},
	{"SMF70R00", 8, 'u', 4},
	{"SMF70R01", 12, 'u', 4},
	{"SMF70R02", 16, 'u', 4},
	{"SMF70R03", 20, 'u', 4},
	{"SMF70R04", 24, 'u', 4},
	{"SMF70R05", 28, 'u', 4},
	{"SMF70R06", 32, 'u', 4},
	{"SMF70R07", 36, 'u', 4},
	{"SMF70R08", 40, 'u', 4},
	{"SMF70R09", 44, 'u', 4},
	{"SMF70R10", 48, 'u', 4},
	{"SMF70R11", 52, 'u', 4},
	{"SMF70R12", 56, 'u', 4},
	{"SMF70R13", 60, 'u', 4},
	{"SMF70R14", 64, 'u', 4},
	{"SMF70R15", 68, 'u', 4},
	{"SMF70IMN", 72, 'u', 2},
	{"SMF70IMM", 74, 'u', 2},
	{"SMF70ITT", 76, 'u', 4},
	{"SMF70I00", 80, 'u', 4},
	{"SMF70I01", 84, 'u', 4}, // manual: 80 / 54
	{"SMF70I02", 88, 'u', 4},
	{"SMF70I03", 92, 'u', 4},
	{"SMF70I04", 96, 'u', 4},
	{"SMF70I05", 100, 'u', 4},
	{"SMF70I06", 104, 'u', 4},
	{"SMF70I07", 108, 'u', 4},
	{"SMF70I08", 112, 'u', 4},
	{"SMF70I09", 116, 'u', 4},
	{"SMF70I10", 120, 'u', 4},
	{"SMF70I11", 124, 'u', 4},
	{"SMF70OMN", 128, 'u', 2},
	{"SMF70OMM", 130, 'u', 2},
	{"SMF70OTT", 132, 'u', 4},
	{"SMF70O00", 136, 'u', 4},
	{"SMF70O01", 140, 'u', 4},
	{"SMF70O02", 144, 'u', 4},
	{"SMF70O03", 148, 'u', 4},
	{"SMF70O04", 152, 'u', 4},
	{"SMF70O05", 156, 'u', 4},
	{"SMF70O06", 160, 'u', 4},
	{"SMF70O07", 164, 'u', 4},
	{"SMF70O08", 168, 'u', 4},
	{"SMF70O09", 172, 'u', 4},
	{"SMF70O10", 176, 'u', 4},
	{"SMF70O11", 180, 'u', 4},
	{"SMF70WMN", 184, 'u', 2},
	{"SMF70WMM", 186, 'u', 2},
	{"SMF70WTT", 188, 'u', 4},
	{"SMF70W00", 192, 'u', 4},
	{"SMF70W01", 196, 'u', 4},
	{"SMF70W02", 200, 'u', 4},
	{"SMF70W03", 204, 'u', 4},
	{"SMF70W04", 208, 'u', 4},
	{"SMF70W05", 212, 'u', 4},
	{"SMF70W06", 216, 'u', 4},
	{"SMF70W07", 220, 'u', 4},
	{"SMF70W08", 224, 'u', 4},
	{"SMF70W09", 228, 'u', 4},
	{"SMF70W10", 232, 'u', 4},
	{"SMF70W11", 236, 'u', 4},
	{"SMF70BMN", 240, 'u', 2},
	{"SMF70BMM", 242, 'u', 2},
	{"SMF70BTT", 244, 'u', 4},
	{"SMF70B00", 248, 'u', 4},
	{"SMF70B01", 252, 'u', 4},
	{"SMF70B02", 256, 'u', 4},
	{"SMF70B03", 260, 'u', 4},
	{"SMF70B04", 264, 'u', 4},
	{"SMF70B05", 268, 'u', 4},
	{"SMF70B06", 272, 'u', 4},
	{"SMF70B07", 276, 'u', 4},
	{"SMF70B08", 280, 'u', 4},
	{"SMF70B09", 284, 'u', 4},
	{"SMF70B10", 288, 'u', 4},
	{"SMF70B11", 292, 'u', 4},
	{"SMF70SMN", 296, 'u', 2},
	{"SMF70SMM", 298, 'u', 2},
	{"SMF70STT", 300, 'u', 4},
	{"SMF70S00", 304, 'u', 4},
	{"SMF70S01", 308, 'u', 4},
	{"SMF70S02", 312, 'u', 4},
	{"SMF70S03", 316, 'u', 4},
	{"SMF70S04", 320, 'u', 4},
	{"SMF70S05", 324, 'u', 4},
	{"SMF70S06", 328, 'u', 4},
	{"SMF70S07", 332, 'u', 4},
	{"SMF70S08", 336, 'u', 4},
	{"SMF70S09", 340, 'u', 4},
	{"SMF70S10", 344, 'u', 4},
	{"SMF70S11", 348, 'u', 4},
	{"SMF70TMN", 352, 'u', 2},
	{"SMF70TMM", 354, 'u', 2},
	{"SMF70TTT", 356, 'u', 4},
	{"SMF70T00", 360, 'u', 4},
	{"SMF70T01", 364, 'u', 4},
	{"SMF70T02", 368, 'u', 4},
	{"SMF70T03", 372, 'u', 4},
	{"SMF70T04", 376, 'u', 4},
	{"SMF70T05", 380, 'u', 4},
	{"SMF70T06", 384, 'u', 4},
	{"SMF70T07", 388, 'u', 4},
	{"SMF70T08", 392, 'u', 4},
	{"SMF70T09", 396, 'u', 4},
	{"SMF70T10", 400, 'u', 4},
	{"SMF70T11", 404, 'u', 4},
	{"SMF70LMN", 408, 'u', 2},
	{"SMF70LMM", 410, 'u', 2},
	{"SMF70LTT", 412, 'u', 4},
	{"SMF70L00", 416, 'u', 4},
	{"SMF70L01", 420, 'u', 4},
	{"SMF70L02", 424, 'u', 4},
	{"SMF70L03", 428, 'u', 4},
	{"SMF70L04", 432, 'u', 4},
	{"SMF70L05", 436, 'u', 4},
	{"SMF70L06", 440, 'u', 4},
	{"SMF70L07", 444, 'u', 4},
	{"SMF70L08", 448, 'u', 4},
	{"SMF70L09", 452, 'u', 4},
	{"SMF70L10", 456, 'u', 4},
	{"SMF70L11", 460, 'u', 4},
	{"SMF70AMN", 464, 'u', 2},
	{"SMF70AMM", 466, 'u', 2},
	{"SMF70ATT", 468, 'u', 4},
	{"SMF70A00", 472, 'u', 4},
	{"SMF70A01", 476, 'u', 4},
	{"SMF70A02", 480, 'u', 4},
	{"SMF70A03", 484, 'u', 4},
	{"SMF70A04", 488, 'u', 4},
	{"SMF70A05", 492, 'u', 4},
	{"SMF70A06", 496, 'u', 4},
	{"SMF70A07", 500, 'u', 4},
	{"SMF70A08", 504, 'u', 4},
	{"SMF70A09", 508, 'u', 4},
	{"SMF70A10", 512, 'u', 4},
	{"SMF70A11", 516, 'u', 4},
	{"SMF70PMN", 520, 'u', 2},
	{"SMF70PMM", 522, 'u', 2},
	{"SMF70PTT", 524, 'u', 4},
	{"SMF70P00", 528, 'u', 4},
	{"SMF70P01", 532, 'u', 4},
	{"SMF70P02", 536, 'u', 4},
	{"SMF70P03", 540, 'u', 4},
	{"SMF70P04", 544, 'u', 4},
	{"SMF70P05", 548, 'u', 4},
	{"SMF70P06", 552, 'u', 4},
	{"SMF70P07", 556, 'u', 4}, // manual: 556 / 22A
	{"SMF70P08", 560, 'u', 4},
	{"SMF70P09", 564, 'u', 4},
	{"SMF70P10", 568, 'u', 4},
	{"SMF70P11", 572, 'u', 4}, // manual: 572 / 23A
	{"SMF70XMN", 576, 'u', 2},
	{"SMF70XMM", 578, 'u', 2},
	{"SMF70XTT", 580, 'u', 4},
	{"SMF70X00", 584, 'u', 4},
	{"SMF70X01", 588, 'u', 4},
	{"SMF70X02", 592, 'u', 4},
	{"SMF70X03", 596, 'u', 4},
	{"SMF70X04", 600, 'u', 4},
	{"SMF70X05", 604, 'u', 4},
	{"SMF70X06", 608, 'u', 4},
	{"SMF70X07", 612, 'u', 4},
	{"SMF70X08", 616, 'u', 4},
	{"SMF70X09", 620, 'u', 4},
	{"SMF70X10", 624, 'u', 4},
	{"SMF70X11", 628, 'u', 4},
	{"SMF70Q00", 632, 'u', 4},
	{"SMF70Q01", 636, 'u', 4},
	{"SMF70Q02", 640, 'u', 4},
	{"SMF70Q03", 644, 'u', 4},
	{"SMF70Q04", 648, 'u', 4},
	{"SMF70Q05", 652, 'u', 4},
	{"SMF70Q06", 656, 'u', 4},
	{"SMF70Q07", 660, 'u', 4},
	{"SMF70Q08", 664, 'u', 4},
	{"SMF70Q09", 668, 'u', 4},
	{"SMF70Q10", 672, 'u', 4},
	{"SMF70Q11", 676, 'u', 4},
	{"SMF70Q12", 680, 'u', 4},
	{"SMF70SRM", 684, 'u', 4},
	{"SMF70CMN", 688, 'u', 4},
	{"SMF70CMM", 692, 'u', 4},
	{"SMF70CTT", 696, 'u', 4},
	{"SMF70DMN", 700, 'u', 4},
	{"SMF70DMM", 704, 'u', 4},
	{"SMF70DTT", 708, 'u', 4},
	{"SMF70EMN", 712, 'u', 4},
	{"SMF70EMM", 716, 'u', 4},
	{"SMF70ETT", 720, 'u', 4},
	{"SMF70U00", 724, 'u', 4},
	{"SMF70U01", 728, 'u', 4},
	{"SMF70U02", 732, 'u', 4},
	{"SMF70U03", 736, 'u', 4},
	{"SMF70U04", 740, 'u', 4},
	{"SMF70U05", 744, 'u', 4},
	{"SMF70U06", 748, 'u', 4},
	{"SMF70U07", 752, 'u', 4},
	{"SMF70U08", 756, 'u', 4},
	{"SMF70U09", 760, 'u', 4},
	{"SMF70U10", 764, 'u', 4},
	{"SMF70U11", 768, 'u', 4},
	{"SMF70U12", 772, 'u', 4},
	{"SMF70U13", 776, 'u', 4},
	{"SMF70U14", 780, 'u', 4},
	{"SMF70U15", 784, 'u', 4},
}

var smf070LparFields = fields{
	{"SMF70LPM", 0, 'c', 8},
	{"SMF70LPN", 8, 'u', 1},
	{"SMF70PFG", 9, 'u', 1},
	{"SMF70BDN", 10, 'u', 2},
	{"SMF70BDS", 12, 'u', 4},
	{"SMF70BDA", 16, 'u', 4},
	{"SMF70SPN", 20, 'c', 8},
	{"SMF70STN", 28, 'c', 8},
	{"SMF70CSF", 40, 'u', 4},
	{"SMF70ESF", 48, 'u', 4},
	{"SMF70MSU", 52, 'u', 4},
	{"SMF70PFL", 56, 'u', 2},
	{"SMF70UPI", 58, 'u', 1},
	{"SMF70MTID", 59, 'u', 1},
	{"SMF70GNM", 60, 'c', 8},
	{"SMF70GMU", 68, 'u', 4},
	{"SMF70HWGr_Name", 72, 'c', 8},
	{"SMF70_BoostInfo", 80, 'u', 1},
}

var smf070LcpuFields = fields{
	{"SMF70PDT", 0, 'u', 8},
	{"SMF70VPA", 8, 'u', 2},
	{"SMF70BPS", 10, 'u', 2},
	{"SMF70VPF", 12, 'u', 1},
	{"SMF70POF", 13, 'u', 1},
	{"SMF70CIX", 14, 'u', 2},
	{"SMF70EDT", 16, 'u', 8},
	{"SMF70ACS", 24, 'u', 4},
	{"SMF70MIS", 28, 'u', 2},
	{"SMF70MAS", 30, 'u', 2},
	{"SMF70NSI", 32, 'u', 4},
	{"SMF70NSA", 36, 'u', 4},
	{"SMF70ONT", 40, 'u', 8},
	{"SMF70WST", 48, 'u', 8},
	{"SMF70PMA", 56, 's', 4},
	{"SMF70NSW", 60, 'u', 4},
	{"SMF70POW", 64, 'u', 4},
	{"SMF70NCA", 68, 'u', 4},
	{"SMF70HW_Cap_Limit", 72, 'u', 4},
	{"SMF70HWGr_Cap_Limit", 76, 'u', 4},
	{"SMF70MTIT", 80, 'u', 8},
	{"SMF70LPF", 88, 'u', 1},
	{"SMF70MaxNL", 89, 'u', 1},
	{"SMF70CordL1", 90, 'u', 1},
	{"SMF70CordL2", 91, 'u', 1},
	{"SMF70CordL3", 92, 'u', 1},
	{"SMF70CordL4", 93, 'u', 1},
	{"SMF70CordL5", 94, 'u', 1},
	{"SMF70CordL6", 95, 'u', 1},
}

var smf070CpuIDFields = fields{
	{"SMF70CIN", 0, 'c', 16},
	{"SMF70CTN", 16, 'u', 2},
	{"SMF70CAN", 20, 'u', 4},
}

var smf070IcsfFields = fields{
	{"R702SNEC", 0, 'f', 8},
	{"R702SNEB", 8, 'f', 8},
	{"R702SNEI", 16, 'f', 8},
	{"R702TNEC", 24, 'f', 8},
	{"R702TNEB", 32, 'f', 8},
	{"R702TNEI", 40, 'f', 8},
	{"R702SNDC", 48, 'f', 8},
	{"R702SNDB", 56, 'f', 8},
	{"R702SNDI", 64, 'f', 8},
	{"R702TNDC", 72, 'f', 8},
	{"R702TNDB", 80, 'f', 8},
	{"R702TNDI", 88, 'f', 8},
	{"R702NMGC", 96, 'f', 8},
	{"R702NMGB", 104, 'f', 8},
	{"R702NMGI", 112, 'f', 8},
	{"R702NMVC", 120, 'f', 8},
	{"R702NMVB", 128, 'f', 8},
	{"R702NMVI", 136, 'f', 8},
	{"R702NHAC", 144, 'f', 8},
	{"R702NHAB", 152, 'f', 8},
	{"R702NHAI", 160, 'f', 8},
	{"R702NPTC", 168, 'f', 8},
	{"R702NPVC", 176, 'f', 8},
	{"R702NH2C", 184, 'f', 8},
	{"R702NH2B", 192, 'f', 8},
	{"R702NH2I", 200, 'f', 8},
	{"R702NH5C", 208, 'f', 8},
	{"R702NH5B", 216, 'f', 8},
	{"R702NH5I", 224, 'f', 8},
	{"R702CDLV", 232, 'u', 4},
	{"R702AESC", 236, 'f', 8},
	{"R702AESB", 244, 'f', 8},
	{"R702AESI", 252, 'f', 8},
	{"R702ASDC", 260, 'f', 8},
	{"R702ASDB", 268, 'f', 8},
	{"R702ASDI", 276, 'f', 8},
	{"R702DRGC", 284, 'f', 8},
	{"R702DRVC", 308, 'f', 8},
	{"R702DEGC", 332, 'f', 8},
	{"R702DEVC", 356, 'f', 8},
	{"R702AMGC", 380, 'f', 8},
	{"R702AMGB", 388, 'f', 8},
	{"R702AMGI", 396, 'f', 8},
	{"R702AMVC", 404, 'f', 8},
	{"R702AMVB", 412, 'f', 8},
	{"R702AMVI", 420, 'f', 8},
	{"R702FPEC", 428, 'f', 8},
	{"R702FPEB", 436, 'f', 8},
	{"R702FPEI", 444, 'f', 8},
	{"R702FPDC", 452, 'f', 8},
	{"R702FPDB", 460, 'f', 8},
	{"R702FPDI", 468, 'f', 8},
	{"R702FPTC", 476, 'f', 8},
	{"R702FPTB", 484, 'f', 8},
	{"R702FPTI", 492, 'f', 8},
	{"R702FXEC", 500, 'f', 8},
	{"R702FXEB", 508, 'f', 8},
	{"R702FXEI", 516, 'f', 8},
	{"R702FXDC", 524, 'f', 8},
	{"R702FXDB", 532, 'f', 8},
	{"R702FXDI", 540, 'f', 8},
	{"R702FXTC", 548, 'f', 8},
	{"R702FXTB", 556, 'f', 8},
	{"R702FXTI", 564, 'f', 8},
	{"R702DQGC", 572, 'f', 8},
	{"R702DQVC", 596, 'f', 8},
}
