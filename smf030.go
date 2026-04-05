package main

import (
	"database/sql"
	"encoding/binary"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var p1Stmt030 *sql.Stmt
var p2Stmt030 *sql.Stmt
var p3Stmt030 *sql.Stmt
var p4Stmt030 *sql.Stmt
var p5Stmt030 *sql.Stmt
var p6Stmt030 *sql.Stmt
var p7Stmt030 *sql.Stmt
var p8Stmt030 *sql.Stmt
var p9Stmt030 *sql.Stmt
var p10Stmt030 *sql.Stmt
var p11Stmt030 *sql.Stmt
var p12Stmt030 *sql.Stmt
var p13Stmt030 *sql.Stmt
var p14Stmt030 *sql.Stmt
var p15Stmt030 *sql.Stmt
var p16Stmt030 *sql.Stmt
var p17Stmt030 *sql.Stmt
var once030 sync.Once
var count030 int

type header030 struct {
	SMF30TME string
	SMF30DTE string
	SMF30SID string
	SMF30WID string
	SMF30STP int
	SMF30SOF int
	SMF30SLN int
	SMF30SON int
	SMF30IOF int
	SMF30ILN int
	SMF30ION int
	SMF30UOF int
	SMF30ULN int
	SMF30UON int
	SMF30TOF int
	SMF30TLN int
	SMF30TON int
	SMF30COF int
	SMF30CLN int
	SMF30CON int
	SMF30AOF int
	SMF30ALN int
	SMF30AON int
	SMF30ROF int
	SMF30RLN int
	SMF30RON int
	SMF30POF int
	SMF30PLN int
	SMF30PON int
	SMF30OOF int
	SMF30OLN int
	SMF30OON int
	SMF30EOF int
	SMF30ELN int
	SMF30EON int
	SMF30EOR int
	SMF30EOS int
	SMF30DRO int
	SMF30DRL int
	SMF30DRN int
	SMF30ARO int
	SMF30ARL int
	SMF30ARN int
	SMF30OPO int
	SMF30OPL int
	SMF30OPN int
	SMF30OPM int
	SMF30UDO int
	SMF30UDL int
	SMF30UDN int
	SMF30UDS int
	SMF30RMO int
	SMF30RML int
	SMF30RMN int
	SMF30RMS int
	SMF30MOF int
	SMF30MLN int
	SMF30MNO int
	SMF30MOS int
	SMF30CDO int
	SMF30CDL int
	SMF30CDN int
	SMF30USO int
	SMF30USL int
	SMF30USN int
}

type ssys030 struct {
	SMF30TYP int
	SMF30RVN string
	SMF30PNM string
	SMF30OSL string
	SMF30SYN string
	SMF30SYP string
}

type ident030 struct {
	SMF30JBN string
	SMF30PGM string
	SMF30STM string
	SMF30UIF string
	SMF30JNM string
	SMF30STN int
	SMF30CLS string
	SMF30JF1 int
	SMF30PGN int
	SMF30JPT int
	SMF30AST int
	SMF30PPS int
	SMF30SIT string
	SMF30STD string
	SMF30RST string
	SMF30RSD string
	SMF30RET string
	SMF30RED string
	SMF30USR string
	SMF30GRP string
	SMF30RUD string
	SMF30TID string
	SMF30TSN string
	SMF30PSN string
	SMF30CL8 string
	SMF30ISS int
	SMF30IET int
	SMF30SSN int
	SMF30EXN string
	SMF30ASI interface{}
	SMF30COR interface{}
}

func smf030Init() {
	creatStmt1 := `
	DROP TABLE IF EXISTS smf030;
	CREATE TABLE smf030 (
	    ID         INT        NOT NULL,
		SMF30TME   DATETIME   NOT NULL,
		SMF30DTE   DATETIME   NOT NULL,
		SMF30SID   TEXT       NOT NULL,
		SMF30WID   TEXT       NOT NULL,
		SMF30TYP   INT        NOT NULL,
		SMF30RVN   TEXT       NOT NULL,
		SMF30PNM   TEXT       NOT NULL,
		SMF30OSL   TEXT       NOT NULL,
		SMF30SYN   TEXT       NOT NULL,
		SMF30SYP   TEXT       NOT NULL,
		SMF30JBN   TEXT       NOT NULL,
		SMF30PGM   TEXT       NOT NULL,
		SMF30STM   TEXT       NOT NULL,
		SMF30UIF   TEXT       NOT NULL,
		SMF30JNM   TEXT       NOT NULL,
		SMF30STN   INT        NOT NULL,
		SMF30CLS   INT        NOT NULL,
		SMF30JF1   INT        NOT NULL,
		SMF30PGN   INT        NOT NULL,
		SMF30JPT   INT        NOT NULL,
		SMF30AST   INT        NOT NULL,
		SMF30PPS   INT        NOT NULL,
		SMF30SIT   DATETIME   NOT NULL,
		SMF30STD   DATETIME   NOT NULL,
		SMF30RST   DATETIME   NOT NULL,
		SMF30RSD   DATETIME   NOT NULL,
		SMF30RET   DATETIME   NOT NULL,
		SMF30RED   DATETIME   NOT NULL,
		SMF30USR   TEXT       NOT NULL,
		SMF30GRP   TEXT       NOT NULL,
		SMF30RUD   TEXT       NOT NULL,
		SMF30TID   TEXT       NOT NULL,
		SMF30TSN   TEXT       NOT NULL,
		SMF30PSN   TEXT       NOT NULL,
		SMF30CL8   TEXT       NOT NULL,
		SMF30ISS   INT        NOT NULL,
		SMF30IET   INT        NOT NULL,
		SMF30SSN   INT        NOT NULL,
		SMF30EXN   TEXT       NOT NULL,
		SMF30ASI   INT        ,
		SMF30COR   TEXT       
	);
	`
	_, err := db.Exec(creatStmt1)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt1)
	}

	creatStmt2 := `
	DROP TABLE IF EXISTS smf030_ioac;
	CREATE TABLE smf030_ioac (
	    ID         INT        NOT NULL,
		SMF30INP   INT        NOT NULL,
		SMF30TEP   INT        NOT NULL,
		SMF30TPT   INT        NOT NULL,
		SMF30TGT   INT        NOT NULL,
		SMF30RDR   INT        NOT NULL,
		SMF30RDT   INT        NOT NULL,
		SMF30TCN   INT        NOT NULL,
		SMF30DCF   INT        NOT NULL,
		SMF30TRR   INT        NOT NULL,
		SMF30AIC   INT        NOT NULL,
		SMF30AID   INT        NOT NULL,
		SMF30AIW   INT        NOT NULL,
		SMF30AIS   INT        NOT NULL,
		SMF30EIC   INT        NOT NULL,
		SMF30EID   INT        NOT NULL,
		SMF30EIW   INT        NOT NULL,
		SMF30EIS   INT        NOT NULL,
		SMF30TEX   INT        NOT NULL,
		SMF30DAS   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt2)
	}

	creatStmt3 := `
	DROP TABLE IF EXISTS smf030_comp;
	CREATE TABLE smf030_comp (
	    ID         INT        NOT NULL,
		SMF30SCC   INT        NOT NULL,
		SMF30STI   INT        NOT NULL,
		SMF30ARC   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt3)
	}

	creatStmt4 := `
	DROP TABLE IF EXISTS smf030_proc;
	CREATE TABLE smf030_proc (
	    ID         INT        NOT NULL,
		SMF30TFL   INT        NOT NULL,
		SMF30CPT   INT        NOT NULL,
		SMF30CPS   INT        NOT NULL,
		SMF30ICU   INT        NOT NULL,
		SMF30ISB   INT        NOT NULL,
		SMF30JVU   INT        NOT NULL,
		SMF30IVU   INT        NOT NULL,
		SMF30JVA   INT        NOT NULL,
		SMF30IVA   INT        NOT NULL,
		SMF30IST   DATETIME   NOT NULL,
		SMF30IDT   DATETIME   NOT NULL,
		SMF30IIP   INT        NOT NULL,
		SMF30RCT   INT        NOT NULL,
		SMF30HPT   INT        NOT NULL,
		SMF30CSC   INT        NOT NULL,
		SMF30DMI   INT        NOT NULL,
		SMF30DMO   INT        NOT NULL,
		SMF30ASR   INT        NOT NULL,
		SMF30ENC   INT        NOT NULL,
		SMF30DET   INT        NOT NULL,
		SMF30CEP   INT        NOT NULL,
		SMF30TF2   INT        NOT NULL,
		SMF30T32   INT        NOT NULL,
		SMF30T33   INT        NOT NULL,
		SMF30_TIME_ON_IFA INT        NOT NULL,
		SMF30_ENCLAVE_TIME_ON_IFA INT        NOT NULL,
		SMF30_DEP_ENCLAVE_TIME_ON_IFA INT        NOT NULL,
		SMF30_TIME_IFA_ON_CP INT        NOT NULL,
		SMF30_ENCLAVE_TIME_IFA_ON_CP INT        NOT NULL,
		SMF30_DEP_ENCLAVE_TIME_IFA_ON_CP INT        NOT NULL,
		SMF30CEPI  INT        NOT NULL,
		SMF30_TIME_ON_zIIP INT        NOT NULL,
		SMF30_ENCLAVE_TIME_ON_zIIP INT        NOT NULL,
		SMF30_DEPENC_TIME_ON_zIIP INT        NOT NULL,
		SMF30_TIME_zIIP_ON_CP INT        NOT NULL,
		SMF30_ENCLAVE_TIME_zIIP_ON_CP INT        NOT NULL,
		SMF30_DEPENC_TIME_zIIP_ON_CP INT        NOT NULL,
		SMF30_ENCLAVE_TIME_zIIP_QUAL INT        NOT NULL,
		SMF30_DEPENC_TIME_zIIP_QUAL INT        NOT NULL,
		SMF30CRP   INT        NOT NULL,
		SMF30ICU_Step_Term INT        NOT NULL,
		SMF30ICU_Step_Init INT        NOT NULL,
		SMF30ISB_Step_Term INT        NOT NULL,
		SMF30ISB_Step_Init INT        NOT NULL,
		SMF30_Missed_SMF30BLK INT        NOT NULL,
		SMF30_Missed_SMF30DCT INT        NOT NULL,
		SMF30_Highest_Task_CPU_Percent INT        NOT NULL,
		SMF30_Highest_Task_CPU_Program TEXT       NOT NULL
	);
	`
	_, err = db.Exec(creatStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt4)
	}

	creatStmt5 := `
	DROP TABLE IF EXISTS smf030_excp;
	CREATE TABLE smf030_excp (
	    ID         INT        NOT NULL,
		SMF30DEV   INT        NOT NULL,
		SMF30UTP   INT        NOT NULL,
		SMF30CUA   INT        NOT NULL,
		SMF30DDN   TEXT       NOT NULL,
		SMF30BLK   INT        NOT NULL,
		SMF30BSZ   INT        NOT NULL,
		SMF30DCT   INT        NOT NULL,
		SMF30XBS   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt5)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt5)
	}

	creatStmt6 := `
	DROP TABLE IF EXISTS smf030_acct;
	CREATE TABLE smf030_acct (
	    ID         INT        NOT NULL,
		SMF30ACT   TEXT       NOT NULL
	);
	`
	_, err = db.Exec(creatStmt6)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt6)
	}

	creatStmt7 := `
	DROP TABLE IF EXISTS smf030_stor;
	CREATE TABLE smf030_stor (
	    ID         INT        NOT NULL,
		SMF30SFL   INT        NOT NULL,
		SMF30SPK   INT        NOT NULL,
		SMF30PRV   INT        NOT NULL,
		SMF30SYS   INT        NOT NULL,
		SMF30PGI   INT        NOT NULL,
		SMF30PGO   INT        NOT NULL,
		SMF30CPM   INT        NOT NULL,
		SMF30NSW   INT        NOT NULL,
		SMF30PSI   INT        NOT NULL,
		SMF30PSO   INT        NOT NULL,
		SMF30VPI   INT        NOT NULL,
		SMF30VPO   INT        NOT NULL,
		SMF30VPR   INT        NOT NULL,
		SMF30CPI   INT        NOT NULL,
		SMF30HPI   INT        NOT NULL,
		SMF30LPI   INT        NOT NULL,
		SMF30HPO   INT        NOT NULL,
		SMF30PST   INT        NOT NULL,
		SMF30PSC   INT        NOT NULL,
		SMF30RGB   INT        NOT NULL,
		SMF30ERG   INT        NOT NULL,
		SMF30ARB   INT        NOT NULL,
		SMF30EAR   INT        NOT NULL,
		SMF30URB   INT        NOT NULL,
		SMF30EUR   INT        NOT NULL,
		SMF30RGN   INT        NOT NULL,
		SMF30DSV   INT        NOT NULL,
		SMF30PIE   INT        NOT NULL,
		SMF30POE   INT        NOT NULL,
		SMF30BIA   INT        NOT NULL,
		SMF30BOA   INT        NOT NULL,
		SMF30BIE   INT        NOT NULL,
		SMF30BOE   INT        NOT NULL,
		SMF30KIA   INT        NOT NULL,
		SMF30KOA   INT        NOT NULL,
		SMF30KIE   INT        NOT NULL,
		SMF30KOE   INT        NOT NULL,
		SMF30PSF   INT        NOT NULL,
		SMF30PAI   INT        NOT NULL,
		SMF30PEI   INT        NOT NULL,
		SMF30ERS   INT        NOT NULL,
		SMF30MEM   INT        NOT NULL,
		SMF30MES   INT        NOT NULL,
		SMF30HVR   INT        NOT NULL,
		SMF30HVA   INT        NOT NULL,
		SMF30HVO   INT        NOT NULL,
		SMF30HVH   INT        NOT NULL,
		SMF30HSO   INT        NOT NULL,
		SMF30HSH   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt7)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt7)
	}

	creatStmt8 := `
	DROP TABLE IF EXISTS smf030_perf;
	CREATE TABLE smf030_perf (
	    ID         INT        NOT NULL,
		SMF30SRV   INT        NOT NULL,
		SMF30CSU   INT        NOT NULL,
		SMF30SRB   INT        NOT NULL,
		SMF30IO    INT        NOT NULL,
		SMF30MSO   INT        NOT NULL,
		SMF30TAT   INT        NOT NULL,
		SMF30SUS   INT        NOT NULL,
		SMF30RES   INT        NOT NULL,
		SMF30TRS   INT        NOT NULL,
		SMF30WLM   TEXT       NOT NULL,
		SMF30SCN   TEXT       NOT NULL,
		SMF30GRN   TEXT       NOT NULL,
		SMF30RCN   TEXT       NOT NULL,
		SMF30ETA   INT        NOT NULL,
		SMF30ESU   INT        NOT NULL,
		SMF30ETC   INT        NOT NULL,
		SMF30PFL   TEXT       NOT NULL,
		SMF30JQT   INT        NOT NULL,
		SMF30RQT   INT        NOT NULL,
		SMF30HQT   INT        NOT NULL,
		SMF30SQT   INT        NOT NULL,
		SMF30PF1   INT        NOT NULL,
		SMF30PF2   INT        NOT NULL,
		SMF30INV   INT        NOT NULL,
		SMF30ZEP   INT        NOT NULL,
		SMF30JPN   TEXT       NOT NULL,
		SMF30MSC   INT        NOT NULL,
		SMF30CPC   INT        NOT NULL,
		SMF30LOC   INT        NOT NULL,
		SMF30SRC   INT        NOT NULL,
		SMF30ZNF   INT        NOT NULL,
		SMF30SNF   INT        NOT NULL,
		SMF30SRV_L INT        NOT NULL,
		SMF30CSU_L INT        NOT NULL,
		SMF30SRB_L INT        NOT NULL,
		SMF30IO_L  INT        NOT NULL,
		SMF30MSO_L INT        NOT NULL,
		SMF30ESU_L INT        NOT NULL,
		SMF30ACB   INT        NOT NULL,
		SMF30CR    INT        NOT NULL,
		SMF30_Capacity_Change_Cnt  INT        NOT NULL,
		SMF30_RCTPCPUA_Actual INT        NOT NULL,
		SMF30_RCTPCPUA_Nominal INT        NOT NULL,
		SMF30_RCTPCPUA_scaling_factor INT        NOT NULL,
		SMF30_Capacity_Adjustment_Ind INT        NOT NULL,
		SMF30_Capacity_Change_Rsn  INT        NOT NULL,
		SMF30_Capacity_Flags  INT        NOT NULL,
		SMF30_RMCTADJN_Nominal INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt8)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt8)
	}

	creatStmt9 := `
	DROP TABLE IF EXISTS smf030_oper;
	CREATE TABLE smf030_oper (
	    ID         INT        NOT NULL,
		SMF30PDM   INT        NOT NULL,
		SMF30PRD   INT        NOT NULL,
		SMF30PTM   INT        NOT NULL,
		SMF30TPR   INT        NOT NULL,
		SMF30MTM   INT        NOT NULL,
		SMF30MSR   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt9)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt9)
	}

	creatStmt10 := `
	DROP TABLE IF EXISTS smf030_appc;
	CREATE TABLE smf030_appc (
	    ID         INT        NOT NULL,
		SMF30DC    INT        NOT NULL,
		SMF30DCA   INT        NOT NULL,
		SMF30DSC   INT        NOT NULL,
		SMF30DDS   BLOB       NOT NULL,
		SMF30DRC   INT        NOT NULL,
		SMF30DDR   BLOB       NOT NULL,
		SMF30DAC   INT        NOT NULL,
		SMF30DTR   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt10)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt10)
	}

	creatStmt11 := `
	DROP TABLE IF EXISTS smf030_appc_cum;
	CREATE TABLE smf030_appc_cum (
	    ID         INT        NOT NULL,
		SMF30CN    INT        NOT NULL,
		SMF30CNA   INT        NOT NULL,
		SMF30SEN   INT        NOT NULL,
		SMF30DAT   BLOB       NOT NULL,
		SMF30REC   INT        NOT NULL,
		SMF30DAR   BLOB       NOT NULL,
		SMF30TAC   INT        NOT NULL,
		SMF30ATR   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt11)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt11)
	}

	creatStmt12 := `
	DROP TABLE IF EXISTS smf030_unix;
	CREATE TABLE smf030_unix (
	    ID         INT        NOT NULL,
		SMF30OPI   INT        NOT NULL,
		SMF30OPG   INT        NOT NULL,
		SMF30OUI   INT        NOT NULL,
		SMF30OUG   INT        NOT NULL,
		SMF30OSI   INT        NOT NULL,
		SMF30OSC   INT        NOT NULL,
		SMF30OST   INT        NOT NULL,
		SMF30ODR   INT        NOT NULL,
		SMF30OFR   INT        NOT NULL,
		SMF30OFW   INT        NOT NULL,
		SMF30OPR   INT        NOT NULL,
		SMF30OPW   INT        NOT NULL,
		SMF30OLL   INT        NOT NULL,
		SMF30OLP   INT        NOT NULL,
		SMF30OGL   INT        NOT NULL,
		SMF30OGP   INT        NOT NULL,
		SMF30OPP   INT        NOT NULL,
		SMF30OSY   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt12)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt12)
	}

	creatStmt13 := `
	DROP TABLE IF EXISTS smf030_arm;
	CREATE TABLE smf030_arm (
	    ID         INT        NOT NULL,
		SMF30RNM   TEXT       NOT NULL,
		SMF30RTP   TEXT       NOT NULL,
		SMF30RRG   TEXT       NOT NULL,
		SMF30RSN   TEXT       NOT NULL,
		SMF30RGT   DATETIME   NOT NULL,
		SMF30RGD   DATETIME   NOT NULL,
		SMF30RWT   DATETIME   NOT NULL,
		SMF30RWD   DATETIME   NOT NULL,
		SMF30RYT   DATETIME   NOT NULL,
		SMF30RYD   DATETIME   NOT NULL,
		SMF30RTT   DATETIME   NOT NULL,
		SMF30RTD   DATETIME   NOT NULL
	);
	`
	_, err = db.Exec(creatStmt13)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt13)
	}

	creatStmt14 := `
	DROP TABLE IF EXISTS smf030_usag;
	CREATE TABLE smf030_usag (
	    ID         INT        NOT NULL,
		SMF30UPO   TEXT       NOT NULL,
		SMF30UPN   TEXT       NOT NULL,
		SMF30UPV   TEXT       NOT NULL,
		SMF30UPQ   TEXT       NOT NULL,
		SMF30UPI   TEXT       NOT NULL,
		SMF30UCT   INT        NOT NULL,
		SMF30UCS   INT        NOT NULL,
		SMF30URD   BLOB       NOT NULL,
		SMF30UDF   INT        NOT NULL,
		SMF30UFG   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt14)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt14)
	}

	creatStmt15 := `
	DROP TABLE IF EXISTS smf030_mers;
	CREATE TABLE smf030_mers (
	    ID         INT        NOT NULL,
		SMF30MRS   TEXT       NOT NULL,
		SMF30MRA   INT        NOT NULL,
		SMF30MRD   INT        NOT NULL,
		SMF30MRI   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt15)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt15)
	}

	creatStmt16 := `
	DROP TABLE IF EXISTS smf030_cntr;
	CREATE TABLE smf030_cntr (
	    ID         INT        NOT NULL,
		SMF30InstFlgs1MRS              INT        NOT NULL,
		SMF30_Inst_CP_Task             INT        NOT NULL,
		SMF30_Inst_CP_NonPreemptSRB    INT        NOT NULL,
		SMF30_Inst_CP_PreemptSRB       INT        NOT NULL,
		SMF30_Inst_Offload             INT        NOT NULL,
		SMF30_Inst_OffloadOnCP         INT        NOT NULL,
		SMF30_Inst_CP_Enclave          INT        NOT NULL,
		SMF30_Inst_Offload_Enclave     INT        NOT NULL,
		SMF30_Inst_OffloadOnCP_Enclave INT        NOT NULL,
		SMF30_Inst_CP_DepEnc           INT        NOT NULL,
		SMF30_Inst_Offload_DepEnc      INT        NOT NULL,
		SMF30_Inst_OffloadOnCP_DepEnc  INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt16)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt16)
	}

	creatStmt17 := `
	DROP TABLE IF EXISTS smf030_zedc;
	CREATE TABLE smf030_zedc (
	    ID         INT        NOT NULL,
		SMF30_US_ComprReq       INT        NOT NULL,
		SMF30_US_ComprReq_Prob  INT        NOT NULL,
		SMF30_US_QueueTime      INT        NOT NULL,
		SMF30_US_ExecTime       INT        NOT NULL,
		SMF30_US_Def_UncomprIn  INT        NOT NULL,
		SMF30_US_Def_ComprOut   INT        NOT NULL,
		SMF30_US_Inf_ComprIn    INT        NOT NULL,
		SMF30_US_Inf_DecomprOut INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt17)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt17)
	}

	insertStmt1 := `
		INSERT INTO smf030
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
`
	p1Stmt030, err = db.Prepare(insertStmt1)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt1)
	}

	insertStmt2 := `
		INSERT INTO smf030_ioac
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p2Stmt030, err = db.Prepare(insertStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt2)
	}

	insertStmt3 := `
		INSERT INTO smf030_comp
		VALUES (
			?,
			?,
			?,
			?
		);
	`

	p3Stmt030, err = db.Prepare(insertStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt3)
	}

	insertStmt4 := `
		INSERT INTO smf030_proc
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p4Stmt030, err = db.Prepare(insertStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt4)
	}

	insertStmt5 := `
		INSERT INTO smf030_excp
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p5Stmt030, err = db.Prepare(insertStmt5)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt5)
	}

	insertStmt6 := `
		INSERT INTO smf030_acct
		VALUES (
			?,
			?
		);
	`

	p6Stmt030, err = db.Prepare(insertStmt6)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt6)
	}

	insertStmt7 := `
		INSERT INTO smf030_stor
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p7Stmt030, err = db.Prepare(insertStmt7)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt7)
	}

	insertStmt8 := `
		INSERT INTO smf030_perf
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p8Stmt030, err = db.Prepare(insertStmt8)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt8)
	}

	insertStmt9 := `
		INSERT INTO smf030_oper
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p9Stmt030, err = db.Prepare(insertStmt9)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt9)
	}

	insertStmt10 := `
		INSERT INTO smf030_appc
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p10Stmt030, err = db.Prepare(insertStmt10)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt10)
	}

	insertStmt11 := `
		INSERT INTO smf030_appc_cum
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p11Stmt030, err = db.Prepare(insertStmt11)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt11)
	}

	insertStmt12 := `
		INSERT INTO smf030_unix
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p12Stmt030, err = db.Prepare(insertStmt12)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt12)
	}

	insertStmt13 := `
		INSERT INTO smf030_arm
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p13Stmt030, err = db.Prepare(insertStmt13)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt13)
	}

	insertStmt14 := `
		INSERT INTO smf030_usag
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p14Stmt030, err = db.Prepare(insertStmt14)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt14)
	}

	insertStmt15 := `
		INSERT INTO smf030_mers
		VALUES (
			?,
			?,
			?,
			?,
			?
		);
	`

	p15Stmt030, err = db.Prepare(insertStmt15)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt15)
	}

	insertStmt16 := `
		INSERT INTO smf030_cntr
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p16Stmt030, err = db.Prepare(insertStmt16)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt16)
	}

	insertStmt17 := `
		INSERT INTO smf030_zedc
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?
		);
	`

	p17Stmt030, err = db.Prepare(insertStmt17)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt17)
	}
}

func (hdr *header030) IsFirst() bool {
	return (hdr.SMF30AON |
		hdr.SMF30ARN |
		hdr.SMF30CON |
		hdr.SMF30DRN |
		hdr.SMF30OON |
		hdr.SMF30PON |
		hdr.SMF30RON |
		hdr.SMF30TON |
		hdr.SMF30UON) != 0
}

func (hdr *header030) IsAdditional() bool {
	return (hdr.SMF30AON |
		hdr.SMF30ARN |
		hdr.SMF30CON |
		hdr.SMF30DRN |
		hdr.SMF30OON |
		hdr.SMF30PON |
		hdr.SMF30RON |
		hdr.SMF30TON |
		hdr.SMF30UON) == 0
}

func (hdr *header030) More() bool {
	return (hdr.SMF30EOS |
		hdr.SMF30MOS |
		hdr.SMF30OPM |
		hdr.SMF30RMS |
		hdr.SMF30UDS) != 0
}

func (hdr *header030) IsLast() bool {
	return (hdr.SMF30EOS |
		hdr.SMF30MOS |
		hdr.SMF30OPM |
		hdr.SMF30RMS |
		hdr.SMF30UDS) == 0
}

func smf030Parse(rec []byte) {
	once030.Do(smf030Init)

	var header header030
	bcd2asc(rec[10:14]) // SMF30SID
	bcd2asc(rec[14:18]) // SMF30WID
	header.SMF30TME = smftime(rec[2:6])
	header.SMF30DTE = smfdate(rec[6:10])
	header.SMF30SID = string(rec[10:14])
	header.SMF30WID = string(rec[14:18])
	header.SMF30STP = int(binary.BigEndian.Uint16(rec[18:20]))   // SMF30STP
	header.SMF30SOF = int(binary.BigEndian.Uint32(rec[20:24]))   // SMF30SOF
	header.SMF30SLN = int(binary.BigEndian.Uint16(rec[24:26]))   // SMF30SLN
	header.SMF30SON = int(binary.BigEndian.Uint16(rec[26:28]))   // SMF30SON
	header.SMF30IOF = int(binary.BigEndian.Uint32(rec[28:32]))   // SMF30IOF
	header.SMF30ILN = int(binary.BigEndian.Uint16(rec[32:34]))   // SMF30ILN
	header.SMF30ION = int(binary.BigEndian.Uint16(rec[34:36]))   // SMF30ION
	header.SMF30UOF = int(binary.BigEndian.Uint32(rec[36:40]))   // SMF30UOF
	header.SMF30ULN = int(binary.BigEndian.Uint16(rec[40:42]))   // SMF30ULN
	header.SMF30UON = int(binary.BigEndian.Uint16(rec[42:44]))   // SMF30UON
	header.SMF30TOF = int(binary.BigEndian.Uint32(rec[44:48]))   // SMF30TOF
	header.SMF30TLN = int(binary.BigEndian.Uint16(rec[48:50]))   // SMF30TLN
	header.SMF30TON = int(binary.BigEndian.Uint16(rec[50:52]))   // SMF30TON
	header.SMF30COF = int(binary.BigEndian.Uint32(rec[52:56]))   // SMF30COF
	header.SMF30CLN = int(binary.BigEndian.Uint16(rec[56:58]))   // SMF30CLN
	header.SMF30CON = int(binary.BigEndian.Uint16(rec[58:60]))   // SMF30CON
	header.SMF30AOF = int(binary.BigEndian.Uint32(rec[60:64]))   // SMF30AOF
	header.SMF30ALN = int(binary.BigEndian.Uint16(rec[64:66]))   // SMF30ALN
	header.SMF30AON = int(binary.BigEndian.Uint16(rec[66:68]))   // SMF30AON
	header.SMF30ROF = int(binary.BigEndian.Uint32(rec[68:72]))   // SMF30ROF
	header.SMF30RLN = int(binary.BigEndian.Uint16(rec[72:74]))   // SMF30RLN
	header.SMF30RON = int(binary.BigEndian.Uint16(rec[74:76]))   // SMF30RON
	header.SMF30POF = int(binary.BigEndian.Uint32(rec[76:80]))   // SMF30POF
	header.SMF30PLN = int(binary.BigEndian.Uint16(rec[80:82]))   // SMF30PLN
	header.SMF30PON = int(binary.BigEndian.Uint16(rec[82:84]))   // SMF30PON
	header.SMF30OOF = int(binary.BigEndian.Uint32(rec[84:88]))   // SMF30OOF
	header.SMF30OLN = int(binary.BigEndian.Uint16(rec[88:90]))   // SMF30OLN
	header.SMF30OON = int(binary.BigEndian.Uint16(rec[90:92]))   // SMF30OON
	header.SMF30EOF = int(binary.BigEndian.Uint32(rec[92:96]))   // SMF30EOF
	header.SMF30ELN = int(binary.BigEndian.Uint16(rec[96:98]))   // SMF30ELN
	header.SMF30EON = int(binary.BigEndian.Uint16(rec[98:100]))  // SMF30EON
	header.SMF30EOR = int(binary.BigEndian.Uint16(rec[100:102])) // SMF30EOR
	header.SMF30EOS = int(binary.BigEndian.Uint32(rec[104:108])) // SMF30EOS
	header.SMF30DRO = int(binary.BigEndian.Uint32(rec[108:112])) // SMF30DRO
	header.SMF30DRL = int(binary.BigEndian.Uint16(rec[112:114])) // SMF30DRL
	header.SMF30DRN = int(binary.BigEndian.Uint16(rec[114:116])) // SMF30DRN
	header.SMF30ARO = int(binary.BigEndian.Uint32(rec[116:120])) // SMF30ARO
	header.SMF30ARL = int(binary.BigEndian.Uint16(rec[120:122])) // SMF30ARL
	header.SMF30ARN = int(binary.BigEndian.Uint16(rec[122:124])) // SMF30ARN
	header.SMF30OPO = int(binary.BigEndian.Uint32(rec[124:128])) // SMF30OPO
	header.SMF30OPL = int(binary.BigEndian.Uint16(rec[128:130])) // SMF30OPL
	header.SMF30OPN = int(binary.BigEndian.Uint16(rec[130:132])) // SMF30OPN
	header.SMF30OPM = int(binary.BigEndian.Uint32(rec[132:136])) // SMF30OPM
	header.SMF30UDO = int(binary.BigEndian.Uint32(rec[136:140])) // SMF30UDO
	header.SMF30UDL = int(binary.BigEndian.Uint16(rec[140:142])) // SMF30UDL
	header.SMF30UDN = int(binary.BigEndian.Uint16(rec[142:144])) // SMF30UDN
	header.SMF30UDS = int(binary.BigEndian.Uint32(rec[144:148])) // SMF30UDS
	header.SMF30RMO = int(binary.BigEndian.Uint32(rec[148:152])) // SMF30RMO
	header.SMF30RML = int(binary.BigEndian.Uint16(rec[152:154])) // SMF30RML
	header.SMF30RMN = int(binary.BigEndian.Uint16(rec[154:156])) // SMF30RMN
	header.SMF30RMS = int(binary.BigEndian.Uint32(rec[156:160])) // SMF30RMS
	header.SMF30MOF = int(binary.BigEndian.Uint32(rec[160:164])) // SMF30MOF
	header.SMF30MLN = int(binary.BigEndian.Uint16(rec[164:166])) // SMF30MLN
	header.SMF30MNO = int(binary.BigEndian.Uint16(rec[166:168])) // SMF30MNO
	header.SMF30MOS = int(binary.BigEndian.Uint32(rec[168:172])) // SMF30MOS
	if 172 < header.SMF30SOF-4 {
		header.SMF30CDO = int(binary.BigEndian.Uint32(rec[172:176])) // SMF30CDO
		header.SMF30CDL = int(binary.BigEndian.Uint16(rec[176:178])) // SMF30CDL
		header.SMF30CDN = int(binary.BigEndian.Uint16(rec[178:180])) // SMF30CDN
	} else {
		header.SMF30CDO = 0
		header.SMF30CDL = 0
		header.SMF30CDN = 0
	}
	if 180 < header.SMF30SOF-4 {
		header.SMF30USO = int(binary.BigEndian.Uint32(rec[180:184])) // SMF30USO
		header.SMF30USL = int(binary.BigEndian.Uint16(rec[184:186])) // SMF30USL
		header.SMF30USN = int(binary.BigEndian.Uint16(rec[186:188])) // SMF30USN
	} else {
		header.SMF30USO = 0
		header.SMF30USL = 0
		header.SMF30USN = 0
	}
	//isFirst := header.IsFirst()
	//isAdd := header.IsAdditional()
	//more := header.More()
	//isLast := header.IsLast()
	//_, _, _, _ = isFirst, isAdd, more, isLast

	// Subsystem section
	ssysOff := header.SMF30SOF - 4
	var ssys ssys030
	bcd2asc(rec[ssysOff+4 : ssysOff+6])                                      // SMF30RVN
	bcd2asc(rec[ssysOff+6 : ssysOff+14])                                     // SMF30PNM
	bcd2asc(rec[ssysOff+14 : ssysOff+22])                                    // SMF30OSL
	bcd2asc(rec[ssysOff+22 : ssysOff+30])                                    // SMF30SYN
	bcd2asc(rec[ssysOff+30 : ssysOff+38])                                    // SMF30SYP
	ssys.SMF30TYP = int(binary.BigEndian.Uint16(rec[ssysOff+0 : ssysOff+2])) // SMF30TYP
	ssys.SMF30RVN = string(rec[ssysOff+4 : ssysOff+6])                       // SMF30RVN
	ssys.SMF30PNM = string(rec[ssysOff+6 : ssysOff+14])                      // SMF30PNM
	ssys.SMF30OSL = string(rec[ssysOff+14 : ssysOff+22])                     // SMF30OSL
	ssys.SMF30SYN = string(rec[ssysOff+22 : ssysOff+30])                     // SMF30SYN
	ssys.SMF30SYP = string(rec[ssysOff+30 : ssysOff+38])                     // SMF30SYP

	// Identification section
	idOff := header.SMF30IOF - 4
	var ident ident030
	bcd2asc(rec[idOff+0 : idOff+8])                                           // SMF30JBN
	bcd2asc(rec[idOff+8 : idOff+16])                                          // SMF30PGM
	bcd2asc(rec[idOff+16 : idOff+24])                                         // SMF30STM
	bcd2asc(rec[idOff+24 : idOff+32])                                         // SMF30UIF
	bcd2asc(rec[idOff+32 : idOff+40])                                         // SMF30JNM
	bcd2asc(rec[idOff+42 : idOff+43])                                         // SMF30JNM
	bcd2asc(rec[idOff+80 : idOff+100])                                        // SMF30USR
	bcd2asc(rec[idOff+100 : idOff+108])                                       // SMF30GRP
	bcd2asc(rec[idOff+108 : idOff+116])                                       // SMF30RUD
	bcd2asc(rec[idOff+116 : idOff+124])                                       // SMF30TID
	bcd2asc(rec[idOff+124 : idOff+132])                                       // SMF30TSN
	bcd2asc(rec[idOff+132 : idOff+140])                                       // SMF30PSN
	bcd2asc(rec[idOff+140 : idOff+148])                                       // SMF30CL8
	bcd2asc(rec[idOff+168 : idOff+184])                                       // SMF30EXN
	bcd2asc(rec[idOff+186 : idOff+250])                                       // SMF30COR
	ident.SMF30JBN = string(rec[idOff+0 : idOff+8])                           // SMF30JBN
	ident.SMF30PGM = string(rec[idOff+8 : idOff+16])                          // SMF30PGM
	ident.SMF30STM = string(rec[idOff+16 : idOff+24])                         // SMF30STM
	ident.SMF30UIF = string(rec[idOff+24 : idOff+32])                         // SMF30UIF
	ident.SMF30JNM = string(rec[idOff+32 : idOff+40])                         // SMF30JNM
	ident.SMF30STN = int(binary.BigEndian.Uint16(rec[idOff+40 : idOff+42]))   // SMF30STN
	ident.SMF30CLS = string(rec[idOff+42 : idOff+43])                         // SMF30CLS
	ident.SMF30JF1 = int(rec[idOff+43])                                       // SMF30JF1
	ident.SMF30PGN = int(binary.BigEndian.Uint16(rec[idOff+44 : idOff+46]))   // SMF30PGN
	ident.SMF30JPT = int(binary.BigEndian.Uint16(rec[idOff+46 : idOff+48]))   // SMF30JPT
	ident.SMF30AST = int(binary.BigEndian.Uint32(rec[idOff+48 : idOff+52]))   // SMF30AST
	ident.SMF30PPS = int(binary.BigEndian.Uint32(rec[idOff+52 : idOff+56]))   // SMF30PPS
	ident.SMF30SIT = smftime(rec[idOff+56 : idOff+60])                        // SMF30SIT
	ident.SMF30STD = smfdate(rec[idOff+60 : idOff+64])                        // SMF30STD
	ident.SMF30RST = smftime(rec[idOff+64 : idOff+68])                        // SMF30RST
	ident.SMF30RSD = smfdate(rec[idOff+68 : idOff+72])                        // SMF30RSD
	ident.SMF30RET = smftime(rec[idOff+72 : idOff+76])                        // SMF30RET
	ident.SMF30RED = smfdate(rec[idOff+76 : idOff+80])                        // SMF30RED
	ident.SMF30USR = string(rec[idOff+80 : idOff+100])                        // SMF30USR
	ident.SMF30GRP = string(rec[idOff+100 : idOff+108])                       // SMF30GRP
	ident.SMF30RUD = string(rec[idOff+108 : idOff+116])                       // SMF30RUD
	ident.SMF30TID = string(rec[idOff+116 : idOff+124])                       // SMF30TID
	ident.SMF30TSN = string(rec[idOff+124 : idOff+132])                       // SMF30TSN
	ident.SMF30PSN = string(rec[idOff+132 : idOff+140])                       // SMF30PSN
	ident.SMF30CL8 = string(rec[idOff+140 : idOff+148])                       // SMF30CL8
	ident.SMF30ISS = int(binary.BigEndian.Uint64(rec[idOff+148 : idOff+156])) // SMF30ISS
	ident.SMF30IET = int(binary.BigEndian.Uint64(rec[idOff+156 : idOff+164])) // SMF30IET
	ident.SMF30SSN = int(binary.BigEndian.Uint32(rec[idOff+164 : idOff+168])) // SMF30SSN
	ident.SMF30EXN = string(rec[idOff+168 : idOff+184])                       // SMF30EXN
	if 184 < header.SMF30ILN {
		ident.SMF30ASI = int(binary.BigEndian.Uint16(rec[idOff+184 : idOff+186])) // SMF30ASI
	} else {
		ident.SMF30ASI = nil
	}
	if 186 < header.SMF30ILN {
		ident.SMF30COR = string(rec[idOff+186 : idOff+250]) // SMF30COR
	} else {
		ident.SMF30COR = nil
	}

	count030 += 1
	_, err = p1Stmt030.Exec(
		count030,
		header.SMF30TME,
		header.SMF30DTE,
		header.SMF30SID,
		header.SMF30WID,
		ssys.SMF30TYP,
		ssys.SMF30RVN,
		ssys.SMF30PNM,
		ssys.SMF30OSL,
		ssys.SMF30SYN,
		ssys.SMF30SYP,
		ident.SMF30JBN,
		ident.SMF30PGM,
		ident.SMF30STM,
		ident.SMF30UIF,
		ident.SMF30JNM,
		ident.SMF30STN,
		ident.SMF30CLS,
		ident.SMF30JF1,
		ident.SMF30PGN,
		ident.SMF30JPT,
		ident.SMF30AST,
		ident.SMF30PPS,
		ident.SMF30SIT,
		ident.SMF30STD,
		ident.SMF30RST,
		ident.SMF30RSD,
		ident.SMF30RET,
		ident.SMF30RED,
		ident.SMF30USR,
		ident.SMF30GRP,
		ident.SMF30RUD,
		ident.SMF30TID,
		ident.SMF30TSN,
		ident.SMF30PSN,
		ident.SMF30CL8,
		ident.SMF30ISS,
		ident.SMF30IET,
		ident.SMF30SSN,
		ident.SMF30EXN,
		ident.SMF30ASI,
		ident.SMF30COR,
	)

	if err != nil {
		log.Println(err)
	}

	// I/O activity section
	if header.SMF30UON > 0 {
		ioacOff := header.SMF30UOF - 4
		_, err = p2Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint32(rec[ioacOff+0:ioacOff+4])),   // SMF30INP
			int(binary.BigEndian.Uint32(rec[ioacOff+4:ioacOff+8])),   // SMF30TEP
			int(binary.BigEndian.Uint32(rec[ioacOff+8:ioacOff+12])),  // SMF30TPT
			int(binary.BigEndian.Uint32(rec[ioacOff+12:ioacOff+16])), // SMF30TGT
			int(rec[ioacOff+16]),                                     // SMF30RDR
			int(rec[ioacOff+17]),                                     // SMF30RDT
			int(binary.BigEndian.Uint32(rec[ioacOff+18:ioacOff+22])), // SMF30TCN
			int(binary.BigEndian.Uint32(rec[ioacOff+22:ioacOff+26])), // SMF30DCF
			int(binary.BigEndian.Uint32(rec[ioacOff+28:ioacOff+32])), // SMF30TRR
			int(binary.BigEndian.Uint32(rec[ioacOff+32:ioacOff+36])), // SMF30AIC
			int(binary.BigEndian.Uint32(rec[ioacOff+36:ioacOff+40])), // SMF30AID
			int(binary.BigEndian.Uint32(rec[ioacOff+40:ioacOff+44])), // SMF30AIW
			int(binary.BigEndian.Uint32(rec[ioacOff+44:ioacOff+48])), // SMF30AIS
			int(binary.BigEndian.Uint32(rec[ioacOff+48:ioacOff+52])), // SMF30EIC
			int(binary.BigEndian.Uint32(rec[ioacOff+52:ioacOff+56])), // SMF30EID
			int(binary.BigEndian.Uint32(rec[ioacOff+56:ioacOff+60])), // SMF30EIW
			int(binary.BigEndian.Uint32(rec[ioacOff+60:ioacOff+64])), // SMF30EIS
			int(binary.BigEndian.Uint64(rec[ioacOff+64:ioacOff+72])), // SMF30TEX
			int(binary.BigEndian.Uint32(rec[ioacOff+72:ioacOff+76])), // SMF30DAS
		)
		if err != nil {
			log.Println(err)
		}
	}

	// Completion section
	if header.SMF30TON > 0 {
		compOff := header.SMF30TOF - 4
		_, err = p3Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint16(rec[compOff+0:compOff+2])), // SMF30SCC
			int(binary.BigEndian.Uint16(rec[compOff+2:compOff+4])), // SMF30STI
			int(binary.BigEndian.Uint32(rec[compOff+4:compOff+8])), // SMF30ARC
		)
		if err != nil {
			log.Println(err)
		}
	}

	// Processor accounting section
	if header.SMF30CON > 0 {
		procOff := header.SMF30COF - 4
		bcd2asc(rec[procOff+178 : procOff+186])
		_, err = p4Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint16(rec[procOff+2:procOff+4])),     // SMF30TFL
			int(binary.BigEndian.Uint32(rec[procOff+4:procOff+8])),     // SMF30CPT
			int(binary.BigEndian.Uint32(rec[procOff+8:procOff+12])),    // SMF30CPS
			int(binary.BigEndian.Uint32(rec[procOff+12:procOff+16])),   // SMF30ICU
			int(binary.BigEndian.Uint32(rec[procOff+16:procOff+20])),   // SMF30ISB
			int(binary.BigEndian.Uint32(rec[procOff+20:procOff+24])),   // SMF30JVU
			int(binary.BigEndian.Uint32(rec[procOff+24:procOff+28])),   // SMF30IVU
			int(binary.BigEndian.Uint32(rec[procOff+28:procOff+32])),   // SMF30JVA
			int(binary.BigEndian.Uint32(rec[procOff+32:procOff+36])),   // SMF30IVA
			smftime(rec[procOff+36:procOff+40]),                        // SMF30IST
			smfdate(rec[procOff+40:procOff+44]),                        // SMF30IDT
			int(binary.BigEndian.Uint32(rec[procOff+44:procOff+48])),   // SMF30IIP
			int(binary.BigEndian.Uint32(rec[procOff+48:procOff+52])),   // SMF30RCT
			int(binary.BigEndian.Uint32(rec[procOff+52:procOff+56])),   // SMF30HPT
			int(binary.BigEndian.Uint32(rec[procOff+56:procOff+60])),   // SMF30CSC
			int(binary.BigEndian.Uint32(rec[procOff+60:procOff+64])),   // SMF30DMI
			int(binary.BigEndian.Uint32(rec[procOff+64:procOff+68])),   // SMF30DMO
			int(binary.BigEndian.Uint32(rec[procOff+68:procOff+72])),   // SMF30ASR
			int(binary.BigEndian.Uint32(rec[procOff+72:procOff+76])),   // SMF30ENC
			int(binary.BigEndian.Uint32(rec[procOff+76:procOff+80])),   // SMF30DET
			int(binary.BigEndian.Uint32(rec[procOff+80:procOff+84])),   // SMF30CEP
			int(rec[procOff+84]),                                       // SMF30TF2
			int(rec[procOff+85]),                                       // SMF30T32
			int(rec[procOff+86]),                                       // SMF30T33
			int(binary.BigEndian.Uint32(rec[procOff+88:procOff+92])),   // SMF30_TIME_ON_IFA
			int(binary.BigEndian.Uint32(rec[procOff+92:procOff+96])),   // SMF30_ENCLAVE_TIME_ON_IFA
			int(binary.BigEndian.Uint32(rec[procOff+96:procOff+100])),  // SMF30_DEP_ENCLAVE_TIME_ON_IFA
			int(binary.BigEndian.Uint32(rec[procOff+100:procOff+104])), // SMF30_TIME_IFA_ON_CP
			int(binary.BigEndian.Uint32(rec[procOff+104:procOff+108])), // SMF30_ENCLAVE_TIME_IFA_ON_CP
			int(binary.BigEndian.Uint32(rec[procOff+108:procOff+112])), // SMF30_DEP_ENCLAVE_TIME_IFA_ON_CP
			int(binary.BigEndian.Uint32(rec[procOff+112:procOff+116])), // SMF30CEPI
			int(binary.BigEndian.Uint32(rec[procOff+116:procOff+120])), // SMF30_TIME_ON_zIIP
			int(binary.BigEndian.Uint32(rec[procOff+120:procOff+124])), // SMF30_ENCLAVE_TIME_ON_zIIP
			int(binary.BigEndian.Uint32(rec[procOff+124:procOff+128])), // SMF30_DEPENC_TIME_ON_zIIP
			int(binary.BigEndian.Uint32(rec[procOff+128:procOff+132])), // SMF30_TIME_zIIP_ON_CP
			int(binary.BigEndian.Uint32(rec[procOff+132:procOff+136])), // SMF30_ENCLAVE_TIME_zIIP_ON_CP
			int(binary.BigEndian.Uint32(rec[procOff+136:procOff+140])), // SMF30_DEPENC_TIME_zIIP_ON_CP
			int(binary.BigEndian.Uint32(rec[procOff+140:procOff+144])), // SMF30_ENCLAVE_TIME_zIIP_QUAL
			int(binary.BigEndian.Uint32(rec[procOff+144:procOff+148])), // SMF30_DEPENC_TIME_zIIP_QUAL
			int(binary.BigEndian.Uint32(rec[procOff+148:procOff+152])), // SMF30CRP
			int(binary.BigEndian.Uint32(rec[procOff+152:procOff+156])), // SMF30ICU_Step_Term
			int(binary.BigEndian.Uint32(rec[procOff+156:procOff+160])), // SMF30ICU_Step_Init
			int(binary.BigEndian.Uint32(rec[procOff+160:procOff+164])), // SMF30ISB_Step_Term
			int(binary.BigEndian.Uint32(rec[procOff+164:procOff+168])), // SMF30ISB_Step_Init
			int(binary.BigEndian.Uint32(rec[procOff+168:procOff+172])), // SMF30_Missed_SMF30BLK
			int(binary.BigEndian.Uint32(rec[procOff+172:procOff+176])), // SMF30_Missed_SMF30DCT
			int(binary.BigEndian.Uint16(rec[procOff+176:procOff+178])), // SMF30_Highest_Task_CPU_Percent
			string(rec[procOff+178:procOff+186]),                       // SMF30_Highest_Task_CPU_Program
		)
		if err != nil {
			log.Println(err)
		}
	}

	// EXCP section
	excpOff := header.SMF30EOF - 4
	for i := 0; i < header.SMF30EON; i++ {
		bcd2asc(rec[excpOff+4 : excpOff+12]) // SMF30DDN
		_, err = p5Stmt030.Exec(
			count030,
			int(rec[excpOff+0]),                                      // SMF30DEV
			int(rec[excpOff+1]),                                      // SMF30UTP
			int(binary.BigEndian.Uint16(rec[excpOff+2:excpOff+4])),   // SMF30CUA
			string(rec[excpOff+4:excpOff+12]),                        // SMF30DDN
			int(binary.BigEndian.Uint32(rec[excpOff+12:excpOff+16])), // SMF30BLK
			int(binary.BigEndian.Uint16(rec[excpOff+16:excpOff+18])), // SMF30BSZ
			int(binary.BigEndian.Uint32(rec[excpOff+18:excpOff+22])), // SMF30DCT
			int(binary.BigEndian.Uint64(rec[excpOff+22:excpOff+30])), // SMF30XBS
		)
		if err != nil {
			log.Println(err)
		}
		excpOff += header.SMF30ELN
	}

	// Account section
	acctOff := header.SMF30AOF - 4
	for i := 0; i < header.SMF30AON; i++ {
		acctLen := int(rec[acctOff])
		fieldOff := acctOff + 1
		bcd2asc(rec[fieldOff : fieldOff+acctLen])
		_, err = p6Stmt030.Exec(
			count030,
			string(rec[fieldOff:fieldOff+acctLen]),
		)
		if err != nil {
			log.Println(err)
		}
		acctOff += header.SMF30ALN
	}

	// Storage and Paging section
	if header.SMF30RON > 0 {
		storOff := header.SMF30ROF - 4
		_, err = p7Stmt030.Exec(
			count030,
			int(rec[storOff+2]),                                        // SMF30SFL
			int(rec[storOff+3]),                                        // SMF30SPK
			int(binary.BigEndian.Uint16(rec[storOff+4:storOff+6])),     // SMF30PRV
			int(binary.BigEndian.Uint16(rec[storOff+6:storOff+8])),     // SMF30SYS
			int(binary.BigEndian.Uint32(rec[storOff+8:storOff+12])),    // SMF30PGI
			int(binary.BigEndian.Uint32(rec[storOff+12:storOff+16])),   // SMF30PGO
			int(binary.BigEndian.Uint32(rec[storOff+16:storOff+20])),   // SMF30CPM
			int(binary.BigEndian.Uint32(rec[storOff+20:storOff+24])),   // SMF30NSW
			int(binary.BigEndian.Uint32(rec[storOff+24:storOff+28])),   // SMF30PSI
			int(binary.BigEndian.Uint32(rec[storOff+28:storOff+32])),   // SMF30PSO
			int(binary.BigEndian.Uint32(rec[storOff+32:storOff+36])),   // SMF30VPI
			int(binary.BigEndian.Uint32(rec[storOff+36:storOff+40])),   // SMF30VPO
			int(binary.BigEndian.Uint32(rec[storOff+40:storOff+44])),   // SMF30VPR
			int(binary.BigEndian.Uint32(rec[storOff+44:storOff+48])),   // SMF30CPI
			int(binary.BigEndian.Uint32(rec[storOff+48:storOff+52])),   // SMF30HPI
			int(binary.BigEndian.Uint32(rec[storOff+52:storOff+56])),   // SMF30LPI
			int(binary.BigEndian.Uint32(rec[storOff+56:storOff+60])),   // SMF30HPO
			int(binary.BigEndian.Uint32(rec[storOff+60:storOff+64])),   // SMF30PST
			int(binary.BigEndian.Uint64(rec[storOff+64:storOff+72])),   // SMF30PSC
			int(binary.BigEndian.Uint32(rec[storOff+72:storOff+76])),   // SMF30RGB
			int(binary.BigEndian.Uint32(rec[storOff+76:storOff+80])),   // SMF30ERG
			int(binary.BigEndian.Uint32(rec[storOff+80:storOff+84])),   // SMF30ARB
			int(binary.BigEndian.Uint32(rec[storOff+84:storOff+88])),   // SMF30EAR
			int(binary.BigEndian.Uint32(rec[storOff+88:storOff+92])),   // SMF30URB
			int(binary.BigEndian.Uint32(rec[storOff+92:storOff+96])),   // SMF30EUR
			int(binary.BigEndian.Uint32(rec[storOff+96:storOff+100])),  // SMF30RGN
			int(binary.BigEndian.Uint32(rec[storOff+100:storOff+104])), // SMF30DSV
			int(binary.BigEndian.Uint32(rec[storOff+104:storOff+108])), // SMF30PIE
			int(binary.BigEndian.Uint32(rec[storOff+108:storOff+112])), // SMF30POE
			int(binary.BigEndian.Uint32(rec[storOff+112:storOff+116])), // SMF30BIA
			int(binary.BigEndian.Uint32(rec[storOff+116:storOff+120])), // SMF30BOA
			int(binary.BigEndian.Uint32(rec[storOff+120:storOff+124])), // SMF30BIE
			int(binary.BigEndian.Uint32(rec[storOff+124:storOff+128])), // SMF30BOE
			int(binary.BigEndian.Uint32(rec[storOff+128:storOff+132])), // SMF30KIA
			int(binary.BigEndian.Uint32(rec[storOff+132:storOff+136])), // SMF30KOA
			int(binary.BigEndian.Uint32(rec[storOff+136:storOff+140])), // SMF30KIE
			int(binary.BigEndian.Uint32(rec[storOff+140:storOff+144])), // SMF30KOE
			int(binary.BigEndian.Uint64(rec[storOff+144:storOff+152])), // SMF30PSF
			int(binary.BigEndian.Uint32(rec[storOff+152:storOff+156])), // SMF30PAI
			int(binary.BigEndian.Uint32(rec[storOff+156:storOff+160])), // SMF30PEI
			int(binary.BigEndian.Uint64(rec[storOff+160:storOff+168])), // SMF30ERS
			int(binary.BigEndian.Uint64(rec[storOff+168:storOff+176])), // SMF30MEM
			int(rec[storOff+176]),                                      // SMF30MES
			int(binary.BigEndian.Uint64(rec[storOff+184:storOff+192])), // SMF30HVR
			int(binary.BigEndian.Uint64(rec[storOff+192:storOff+200])), // SMF30HVA
			int(binary.BigEndian.Uint64(rec[storOff+200:storOff+208])), // SMF30HVO
			int(binary.BigEndian.Uint64(rec[storOff+208:storOff+216])), // SMF30HVH
			int(binary.BigEndian.Uint64(rec[storOff+216:storOff+224])), // SMF30HSO
			int(binary.BigEndian.Uint64(rec[storOff+224:storOff+232])), // SMF30HSH
		)
		if err != nil {
			log.Println(err)
		}
	}

	// Performance section
	perfOff := header.SMF30POF - 4
	for i := 0; i < header.SMF30PON; i++ {
		bcd2asc(rec[perfOff+36 : perfOff+44])   // SMF30WLM
		bcd2asc(rec[perfOff+44 : perfOff+52])   // SMF30SCN
		bcd2asc(rec[perfOff+52 : perfOff+60])   // SMF30GRN
		bcd2asc(rec[perfOff+60 : perfOff+68])   // SMF30RCN
		bcd2asc(rec[perfOff+80 : perfOff+96])   // SMF30PFL
		bcd2asc(rec[perfOff+116 : perfOff+124]) // SMF30JPN
		_, err = p8Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint32(rec[perfOff+0:perfOff+4])),     // SMF30SRV
			int(binary.BigEndian.Uint32(rec[perfOff+4:perfOff+8])),     // SMF30CSU
			int(binary.BigEndian.Uint32(rec[perfOff+8:perfOff+12])),    // SMF30SRB
			int(binary.BigEndian.Uint32(rec[perfOff+12:perfOff+16])),   // SMF30IO
			int(binary.BigEndian.Uint32(rec[perfOff+16:perfOff+20])),   // SMF30MSO
			int(binary.BigEndian.Uint32(rec[perfOff+20:perfOff+24])),   // SMF30TAT
			int(binary.BigEndian.Uint32(rec[perfOff+24:perfOff+28])),   // SMF30SUS
			int(binary.BigEndian.Uint32(rec[perfOff+28:perfOff+32])),   // SMF30RES
			int(binary.BigEndian.Uint32(rec[perfOff+32:perfOff+36])),   // SMF30TRS
			string(rec[perfOff+36:perfOff+44]),                         // SMF30WLM
			string(rec[perfOff+44:perfOff+52]),                         // SMF30SCN
			string(rec[perfOff+52:perfOff+60]),                         // SMF30GRN
			string(rec[perfOff+60:perfOff+68]),                         // SMF30RCN
			int(binary.BigEndian.Uint32(rec[perfOff+68:perfOff+72])),   // SMF30ETA
			int(binary.BigEndian.Uint32(rec[perfOff+72:perfOff+76])),   // SMF30ESU
			int(binary.BigEndian.Uint32(rec[perfOff+76:perfOff+80])),   // SMF30ETC
			string(rec[perfOff+80:perfOff+96]),                         // SMF30PFL
			int(binary.BigEndian.Uint32(rec[perfOff+96:perfOff+100])),  // SMF30JQT
			int(binary.BigEndian.Uint32(rec[perfOff+100:perfOff+104])), // SMF30RQT
			int(binary.BigEndian.Uint32(rec[perfOff+104:perfOff+108])), // SMF30HQT
			int(binary.BigEndian.Uint32(rec[perfOff+108:perfOff+112])), // SMF30SQT
			int(rec[perfOff+112]),                                      // SMF30PF1
			int(rec[perfOff+113]),                                      // SMF30PF2
			int(rec[perfOff+114]),                                      // SMF30INV
			int(rec[perfOff+115]),                                      // SMF30ZEP
			string(rec[perfOff+116:perfOff+124]),                       // SMF30JPN
			int(binary.BigEndian.Uint32(rec[perfOff+124:perfOff+128])), // SMF30MSC
			int(binary.BigEndian.Uint16(rec[perfOff+128:perfOff+130])), // SMF30CPC
			int(binary.BigEndian.Uint16(rec[perfOff+130:perfOff+132])), // SMF30LOC
			int(binary.BigEndian.Uint16(rec[perfOff+132:perfOff+134])), // SMF30SRC
			int(binary.BigEndian.Uint16(rec[perfOff+134:perfOff+136])), // SMF30ZNF
			int(binary.BigEndian.Uint16(rec[perfOff+136:perfOff+138])), // SMF30SNF
			int(binary.BigEndian.Uint64(rec[perfOff+144:perfOff+152])), // SMF30SRV_L
			int(binary.BigEndian.Uint64(rec[perfOff+152:perfOff+160])), // SMF30CSU_L
			int(binary.BigEndian.Uint64(rec[perfOff+160:perfOff+168])), // SMF30SRB_L
			int(binary.BigEndian.Uint64(rec[perfOff+168:perfOff+176])), // SMF30IO_L
			int(binary.BigEndian.Uint64(rec[perfOff+176:perfOff+184])), // SMF30MSO_L
			int(binary.BigEndian.Uint64(rec[perfOff+184:perfOff+192])), // SMF30ESU_L
			int(rec[perfOff+192]),                                      // SMF30ACB
			int(rec[perfOff+193]),                                      // SMF30CR
			int(binary.BigEndian.Uint16(rec[perfOff+194:perfOff+196])), // SMF30_Capacity_Change_Cnt
			int(binary.BigEndian.Uint32(rec[perfOff+196:perfOff+200])), // SMF30_RCTPCPUA_Actual
			int(binary.BigEndian.Uint32(rec[perfOff+200:perfOff+204])), // SMF30_RCTPCPUA_Nominal
			int(binary.BigEndian.Uint32(rec[perfOff+204:perfOff+208])), // SMF30_RCTPCPUA_scaling_factor
			int(rec[perfOff+208]),                                      // SMF30_Capacity_Adjustment_Ind
			int(rec[perfOff+209]),                                      // SMF30_Capacity_Change_Rsn
			int(rec[perfOff+210]),                                      // SMF30_Capacity_Flags
			int(binary.BigEndian.Uint32(rec[perfOff+212:perfOff+216])), // SMF30_RMCTADJN_Nominal
		)
		if err != nil {
			log.Println(err)
		}
		perfOff += header.SMF30PLN
	}

	// Operator section
	operOff := header.SMF30OOF - 4
	for i := 0; i < header.SMF30OON; i++ {
		_, err = p9Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint32(rec[operOff+0:operOff+4])),   // SMF30PDM
			int(binary.BigEndian.Uint32(rec[operOff+4:operOff+8])),   // SMF30PRD
			int(binary.BigEndian.Uint32(rec[operOff+8:operOff+12])),  // SMF30PTM
			int(binary.BigEndian.Uint32(rec[operOff+12:operOff+16])), // SMF30TPR
			int(binary.BigEndian.Uint32(rec[operOff+16:operOff+20])), // SMF30MTM
			int(binary.BigEndian.Uint32(rec[operOff+20:operOff+24])), // SMF30MSR
		)
		if err != nil {
			log.Println(err)
		}
		operOff += header.SMF30OLN
	}

	// APPC resource section
	// TODO: Convert fields SMF30DDS and SMF30DDR from HFP to BFP
	if header.SMF30DRN > 0 {
		appcOff := header.SMF30DRO - 4
		_, err = p10Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint32(rec[appcOff+0:appcOff+4])),   // SMF30DC
			int(binary.BigEndian.Uint32(rec[appcOff+4:appcOff+8])),   // SMF30DCA
			int(binary.BigEndian.Uint32(rec[appcOff+8:appcOff+12])),  // SMF30DSC
			rec[appcOff+12:appcOff+20],                               // SMF30DDS
			int(binary.BigEndian.Uint32(rec[appcOff+20:appcOff+24])), // SMF30DRC
			rec[appcOff+24:appcOff+32],                               // SMF30DDR
			int(binary.BigEndian.Uint32(rec[appcOff+32:appcOff+36])), // SMF30DAC
			int(binary.BigEndian.Uint32(rec[appcOff+36:appcOff+40])), // SMF30DTR
		)
		if err != nil {
			log.Println(err)
		}
	}

	// APPC cumulative resource section
	// TODO: Convert fields SMF30DAT and SMF30DAR from HFP to BFP
	if header.SMF30ARN > 0 {
		appcOff := header.SMF30ARO - 4
		_, err = p11Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint32(rec[appcOff+0:appcOff+4])),   // SMF30CN
			int(binary.BigEndian.Uint32(rec[appcOff+4:appcOff+8])),   // SMF30CNA
			int(binary.BigEndian.Uint32(rec[appcOff+8:appcOff+12])),  // SMF30SEN
			rec[appcOff+12:appcOff+20],                               // SMF30DAT
			int(binary.BigEndian.Uint32(rec[appcOff+20:appcOff+24])), // SMF30REC
			rec[appcOff+24:appcOff+32],                               // SMF30DAR
			int(binary.BigEndian.Uint32(rec[appcOff+32:appcOff+36])), // SMF30TAC
			int(binary.BigEndian.Uint32(rec[appcOff+36:appcOff+40])), // SMF30ATR
		)
		if err != nil {
			log.Println(err)
		}
	}

	// UNIX process section
	unixOff := header.SMF30OPO - 4
	for i := 0; i < header.SMF30OPN; i++ {
		_, err = p12Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint32(rec[unixOff+0:unixOff+4])),   // SMF30OPI
			int(binary.BigEndian.Uint32(rec[unixOff+4:unixOff+8])),   // SMF30OPG
			int(binary.BigEndian.Uint32(rec[unixOff+8:unixOff+12])),  // SMF30OUI
			int(binary.BigEndian.Uint32(rec[unixOff+12:unixOff+16])), // SMF30OUG
			int(binary.BigEndian.Uint32(rec[unixOff+16:unixOff+20])), // SMF30OSI
			int(binary.BigEndian.Uint32(rec[unixOff+20:unixOff+24])), // SMF30OSC
			int(binary.BigEndian.Uint32(rec[unixOff+24:unixOff+28])), // SMF30OST
			int(binary.BigEndian.Uint32(rec[unixOff+28:unixOff+32])), // SMF30ODR
			int(binary.BigEndian.Uint32(rec[unixOff+32:unixOff+36])), // SMF30OFR
			int(binary.BigEndian.Uint32(rec[unixOff+36:unixOff+40])), // SMF30OFW
			int(binary.BigEndian.Uint32(rec[unixOff+40:unixOff+44])), // SMF30OPR
			int(binary.BigEndian.Uint32(rec[unixOff+44:unixOff+48])), // SMF30OPW
			int(binary.BigEndian.Uint32(rec[unixOff+56:unixOff+60])), // SMF30OLL
			int(binary.BigEndian.Uint32(rec[unixOff+60:unixOff+64])), // SMF30OLP
			int(binary.BigEndian.Uint32(rec[unixOff+64:unixOff+68])), // SMF30OGL
			int(binary.BigEndian.Uint32(rec[unixOff+68:unixOff+72])), // SMF30OGP
			int(binary.BigEndian.Uint32(rec[unixOff+72:unixOff+76])), // SMF30OPP
			int(binary.BigEndian.Uint32(rec[unixOff+92:unixOff+96])), // SMF30OSY
		)
		if err != nil {
			log.Println(err)
		}
		unixOff += header.SMF30OPL
	}

	// ARM section
	armOff := header.SMF30RMO - 4
	for i := 0; i < header.SMF30RMN; i++ {
		bcd2asc(rec[armOff+0 : armOff+16])  // SMF30RNM
		bcd2asc(rec[armOff+16 : armOff+24]) // SMF30RTP
		bcd2asc(rec[armOff+24 : armOff+40]) // SMF30RRG
		bcd2asc(rec[armOff+40 : armOff+48]) // SMF30RSN
		_, err = p13Stmt030.Exec(
			count030,
			string(rec[armOff+0:armOff+16]),   // SMF30RNM
			string(rec[armOff+16:armOff+24]),  // SMF30RTP
			string(rec[armOff+24:armOff+40]),  // SMF30RRG
			string(rec[armOff+40:armOff+48]),  // SMF30RSN
			smftime(rec[armOff+48:armOff+52]), // SMF30RGT
			smfdate(rec[armOff+52:armOff+56]), // SMF30RGD
			smftime(rec[armOff+56:armOff+60]), // SMF30RWT
			smfdate(rec[armOff+60:armOff+64]), // SMF30RWD
			smftime(rec[armOff+64:armOff+68]), // SMF30RYT
			smfdate(rec[armOff+68:armOff+72]), // SMF30RYD
			smftime(rec[armOff+72:armOff+76]), // SMF30RTT
			smfdate(rec[armOff+76:armOff+80]), // SMF30RTD
		)
		if err != nil {
			log.Println(err)
		}
		armOff += header.SMF30RML
	}

	// Usage data section
	usagOff := header.SMF30UDO - 4
	for i := 0; i < header.SMF30UDN; i++ {
		bcd2asc(rec[usagOff+0 : usagOff+16])  // SMF30UPO
		bcd2asc(rec[usagOff+16 : usagOff+32]) // SMF30UPN
		bcd2asc(rec[usagOff+32 : usagOff+40]) // SMF30UPV
		bcd2asc(rec[usagOff+40 : usagOff+48]) // SMF30UPQ
		bcd2asc(rec[usagOff+48 : usagOff+56]) // SMF30UPI
		_, err = p14Stmt030.Exec(
			count030,
			string(rec[usagOff+0:usagOff+16]),                        // SMF30UPO
			string(rec[usagOff+16:usagOff+32]),                       // SMF30UPN
			string(rec[usagOff+32:usagOff+40]),                       // SMF30UPV
			string(rec[usagOff+40:usagOff+48]),                       // SMF30UPQ
			string(rec[usagOff+48:usagOff+56]),                       // SMF30UPI
			int(binary.BigEndian.Uint32(rec[usagOff+56:usagOff+60])), // SMF30UCT
			int(binary.BigEndian.Uint32(rec[usagOff+60:usagOff+64])), // SMF30UCS
			rec[usagOff+64:usagOff+72],                               // SMF30URD
			int(rec[usagOff+72]),                                     // SMF30UDF
			int(rec[usagOff+73]),                                     // SMF30UFG
		)
		if err != nil {
			log.Println(err)
		}
		usagOff += header.SMF30UDL
	}

	// Multisystem enclave remote system data section
	mersOff := header.SMF30MOF - 4
	for i := 0; i < header.SMF30MNO; i++ {
		bcd2asc(rec[mersOff+0 : mersOff+8]) // SMF30MRS
		_, err = p15Stmt030.Exec(
			count030,
			string(rec[mersOff+0:mersOff+8]),                         // SMF30MRS
			int(binary.BigEndian.Uint32(rec[mersOff+8:mersOff+12])),  // SMF30MRA
			int(binary.BigEndian.Uint32(rec[mersOff+12:mersOff+16])), // SMF30MRD
			int(binary.BigEndian.Uint32(rec[mersOff+16:mersOff+20])), // SMF30MRI
		)
		if err != nil {
			log.Println(err)
		}
		mersOff += header.SMF30MLN
	}

	// Counter data section
	cntrOff := header.SMF30CDO - 4
	for i := 0; i < header.SMF30CDN; i++ {
		_, err = p16Stmt030.Exec(
			count030,
			int(rec[cntrOff+0]),                                      // SMF30InstFlgs1MRS
			int(binary.BigEndian.Uint64(rec[cntrOff+8:cntrOff+16])),  // SMF30_Inst_CP_Task
			int(binary.BigEndian.Uint64(rec[cntrOff+16:cntrOff+24])), // SMF30_Inst_CP_NonPreemptSRB
			int(binary.BigEndian.Uint64(rec[cntrOff+24:cntrOff+32])), // SMF30_Inst_CP_PreemptSRB
			int(binary.BigEndian.Uint64(rec[cntrOff+32:cntrOff+40])), // SMF30_Inst_Offload
			int(binary.BigEndian.Uint64(rec[cntrOff+40:cntrOff+48])), // SMF30_Inst_OffloadOnCP
			int(binary.BigEndian.Uint64(rec[cntrOff+48:cntrOff+56])), // SMF30_Inst_CP_Enclave
			int(binary.BigEndian.Uint64(rec[cntrOff+56:cntrOff+64])), // SMF30_Inst_Offload_Enclave
			int(binary.BigEndian.Uint64(rec[cntrOff+64:cntrOff+72])), // SMF30_Inst_OffloadOnCP_Enclave
			int(binary.BigEndian.Uint64(rec[cntrOff+72:cntrOff+80])), // SMF30_Inst_CP_DepEnc
			int(binary.BigEndian.Uint64(rec[cntrOff+80:cntrOff+88])), // SMF30_Inst_Offload_DepEnc
			int(binary.BigEndian.Uint64(rec[cntrOff+88:cntrOff+96])), // SMF30_Inst_OffloadOnCP_DepEnc
		)
		if err != nil {
			log.Println(err)
		}
		cntrOff += header.SMF30CDL
	}

	// zEDC data section
	zedcOff := header.SMF30USO - 4
	for i := 0; i < header.SMF30USN; i++ {
		_, err = p17Stmt030.Exec(
			count030,
			int(binary.BigEndian.Uint64(rec[zedcOff+0:zedcOff+8])),   // SMF30_US_ComprReq
			int(binary.BigEndian.Uint64(rec[zedcOff+8:zedcOff+16])),  // SMF30_US_ComprReq_Prob
			int(binary.BigEndian.Uint64(rec[zedcOff+16:zedcOff+24])), // SMF30_US_QueueTime
			int(binary.BigEndian.Uint64(rec[zedcOff+24:zedcOff+32])), // SMF30_US_ExecTime
			int(binary.BigEndian.Uint64(rec[zedcOff+32:zedcOff+40])), // SMF30_US_Def_UncomprIn
			int(binary.BigEndian.Uint64(rec[zedcOff+40:zedcOff+48])), // SMF30_US_Def_ComprOut
			int(binary.BigEndian.Uint64(rec[zedcOff+48:zedcOff+56])), // SMF30_US_Inf_ComprIn
			int(binary.BigEndian.Uint64(rec[zedcOff+56:zedcOff+64])), // SMF30_US_Inf_DecomprOut
		)
		if err != nil {
			log.Println(err)
		}
		zedcOff += header.SMF30USL
	}
}
