package main

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

var p1Stmt030 *sql.Stmt
var p2Stmt030 *sql.Stmt
var p3Stmt030 *sql.Stmt
var p4Stmt030 *sql.Stmt
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
	// numVols := int(rec[131])
	isFirst := header.IsFirst()
	isAdd := header.IsAdditional()
	more := header.More()
	isLast := header.IsLast()
	_, _, _, _ = isFirst, isAdd, more, isLast

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
			int(rec[ioacOff+16]), // SMF30RDR
			int(rec[ioacOff+17]), // SMF30RDT
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
			int(binary.BigEndian.Uint16(rec[procOff+2:procOff+4])),   // SMF30TFL
			int(binary.BigEndian.Uint32(rec[procOff+4:procOff+8])),   // SMF30CPT
			int(binary.BigEndian.Uint32(rec[procOff+8:procOff+12])),  // SMF30CPS
			int(binary.BigEndian.Uint32(rec[procOff+12:procOff+16])), // SMF30ICU
			int(binary.BigEndian.Uint32(rec[procOff+16:procOff+20])), // SMF30ISB
			int(binary.BigEndian.Uint32(rec[procOff+20:procOff+24])), // SMF30JVU
			int(binary.BigEndian.Uint32(rec[procOff+24:procOff+28])), // SMF30IVU
			int(binary.BigEndian.Uint32(rec[procOff+28:procOff+32])), // SMF30JVA
			int(binary.BigEndian.Uint32(rec[procOff+32:procOff+36])), // SMF30IVA
			smftime(rec[procOff+36:procOff+40]),                      // SMF30IST
			smfdate(rec[procOff+40:procOff+44]),                      // SMF30IDT
			int(binary.BigEndian.Uint32(rec[procOff+44:procOff+48])), // SMF30IIP
			int(binary.BigEndian.Uint32(rec[procOff+48:procOff+52])), // SMF30RCT
			int(binary.BigEndian.Uint32(rec[procOff+52:procOff+56])), // SMF30HPT
			int(binary.BigEndian.Uint32(rec[procOff+56:procOff+60])), // SMF30CSC
			int(binary.BigEndian.Uint32(rec[procOff+60:procOff+64])), // SMF30DMI
			int(binary.BigEndian.Uint32(rec[procOff+64:procOff+68])), // SMF30DMO
			int(binary.BigEndian.Uint32(rec[procOff+68:procOff+72])), // SMF30ASR
			int(binary.BigEndian.Uint32(rec[procOff+72:procOff+76])), // SMF30ENC
			int(binary.BigEndian.Uint32(rec[procOff+76:procOff+80])), // SMF30DET
			int(binary.BigEndian.Uint32(rec[procOff+80:procOff+84])), // SMF30CEP
			int(rec[procOff+84]), // SMF30TF2
			int(rec[procOff+85]), // SMF30T32
			int(rec[procOff+86]), // SMF30T33
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
}
