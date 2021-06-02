package main

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	//"strings"
)

var p1Stmt015 *sql.Stmt
var p2Stmt015 *sql.Stmt
var once015 sync.Once

func smf015Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf015;
	CREATE TABLE smf015 (
		SMF15TME   DATETIME   NOT NULL,
		SMF15DTE   DATETIME   NOT NULL,
		SMF15SID   TEXT       NOT NULL,
		SMF15JBN   TEXT       NOT NULL,
		SMF15RST   DATETIME   NOT NULL,
		SMF15RSD   DATETIME   NOT NULL,
		SMF15UID   TEXT       NOT NULL,
		SMF15RIN   INT        NOT NULL,
		SMF15SDC   INT        NOT NULL,
		SMF15NUC   INT        NOT NULL,
		SMF15SUC   INT        NOT NULL,
		SMF15SET   INT        NOT NULL,
		SMF15OPE   DATETIME   NOT NULL,
		SMFTIOE1   INT        NOT NULL,
		SMFTIOE2   INT        NOT NULL,
		SMFTIOE3   INT        NOT NULL,
		SMFTIOE4   INT        NOT NULL,
		SMFTIOE5   TEXT       NOT NULL,
		SMFTIOE6   INT        NOT NULL,
		SMFTIOE7   INT        NOT NULL,
		JFCBDSNM   TEXT       NOT NULL,
		JFCBELNM   TEXT       NOT NULL,
		JFCBTSDM   INT        NOT NULL,
		JFCBDSCB   INT        NOT NULL,
		JFCAMSTR   INT        NOT NULL,
		JFCBADBF   INT        NOT NULL,
		JFCNLREC   INT        NOT NULL,
		JFCTDSI1   INT        NOT NULL,
		JFCTDSI2   INT        NOT NULL,
		JFCBLTYP   INT        NOT NULL,
		JFCBFLSQ   INT        NOT NULL,
		JFCBVLSQ   INT        NOT NULL,
		JFCBFLG1   INT        NOT NULL,
		JFCBFLG2   INT        NOT NULL,
		JFCBCRDT   DATETIME   NOT NULL,
		JFCBXPDT   DATETIME   NOT NULL,
		JFCBIND1   INT        NOT NULL,
		JFCBIND2   INT        NOT NULL,
		JFCDSRG1   INT        NOT NULL,
		JFCDSRG2   INT        NOT NULL,
		JFCRECFM   INT        NOT NULL,
		JFCBLKSI   INT        NOT NULL,
		JFCLRECL   INT        NOT NULL,
		JFCBPQTY   INT        NOT NULL,
		JFCBCTRI   INT        NOT NULL,
		JFCBSQTY   INT        NOT NULL,
		JFCBDQTY   INT        NOT NULL,
		SMFDCBOR   INT        NOT NULL,
		SMFDCBRF   INT        NOT NULL,
		SMFDCBMF   INT        NOT NULL,
		SMFDCBFL   INT        NOT NULL,
		SMFDCBOP   INT        NOT NULL,
		SMFDEBFL   INT        NOT NULL,
		SMFDEBOP   INT        NOT NULL,
		SMFDEBVL   INT        NOT NULL,
		SMFDCBBL   INT        ,
		SMFDSSNO   TEXT       ,
		SMF15TDA   INT        ,
		SMF15NTU   INT        ,
		SMF15NTR   INT        ,
		SMF15NER   INT        ,
		SMF15EDI   INT        ,
		SMF15FG1   INT        ,
		SMF15OPD   DATETIME   NOT NULL,
		SMF15XF1   INT        ,
		SMF15CDL   INT        ,
		SMF15UDL   INT        ,
		SMF15CDS   INT        ,
		SMF15UDS   INT        ,
		SMF15CIS   INT        ,
		SMF15TKL   INT        ,
		SMF15TKN   BLOB       ,
		SMF15CMPTYPE INT      ,
		SMF15MCN   TEXT       ,
		SMF15DCN   TEXT       ,
		SMF15SCN   TEXT       ,
		SMF15SPN   TEXT       ,
		SMF15PGN   TEXT       ,
		SMF15CFG   INT        ,
		SMF15USR   INT        ,
		SMF15TPE   INT        ,
		SMF15LBL   INT        ,
		SMF15LBS   INT        ,        
		SMF15FLG1  INT        ,        
		SMF15DSVER INT        ,        
		SMF15ALI   TEXT       ,
		SMF15DRD   INT        ,        
		SMF15DRDH  INT        ,        
		SMF15MRD   INT        ,        
		SMF15MRDH  INT        ,        
		SMF15MCE   INT        ,        
		SMF15MST   INT        ,        
		SMF15MNC   INT        ,        
		SMF15MCF   INT        ,
		SMF15KL1   TEXT       ,        
		SMF15CD1   TEXT       ,        
		SMF15KL2   TEXT       ,        
		SMF15CD2   TEXT       ,        
		SMF15KET   INT        ,        
		SMF15RFG0  INT        ,        
		SMF15RFG1  INT        ,        
		SMF15RASDATA0 INT     ,        
		SMF15RASDATA1 INT     ,
		SMF15DEF   INT        ,
		SMF15DET   INT        ,        
		SMF15DKL   INT                
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
		INSERT INTO smf015
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
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
			?,
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
	p1Stmt015, err = db.Prepare(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
	}

	sqlStmt3 := `
	DROP TABLE IF EXISTS smf015_ucb;
	CREATE TABLE smf015_ucb (
		SMF15TME   DATETIME   NOT NULL,
		SMF15DTE   DATETIME   NOT NULL,
		SMFUCBDV   INT        NOT NULL,
		SMFSRTEV   TEXT       NOT NULL,
		SMFUCBTY   INT        NOT NULL,
		SMFSRTES   INT        NOT NULL,
		SMF15NEX   INT        NOT NULL,
		SMFEXCP    INT        NOT NULL,
		SMFSRTEF   INT        ,
		SMFSRTEQ   INT        ,
		SMF15NTA   INT        
	);
	`
	_, err = db.Exec(sqlStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt3)
	}

	sqlStmt4 := `
		INSERT INTO smf015_ucb
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
	p2Stmt015, err = db.Prepare(sqlStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt4)
	}
}

func smf015Parse(rec []byte) {
	once015.Do(smf015Init)

	const (
		JFCBOFF = 64
		UCBOFF  = 268
	)

	numUcbs := int(rec[41])
	ucbSize := int(rec[42])

	bcd2asc(rec[10:14])                         // SMF15SID
	bcd2asc(rec[14:22])                         // SMF15JBN
	bcd2asc(rec[30:38])                         // SMF15UID
	bcd2asc(rec[52:60])                         // SMFTIOE5
	bcd2asc(rec[JFCBOFF : JFCBOFF+44])          // SMFDSNM
	bcd2asc(rec[JFCBOFF+0x2c : JFCBOFF+0x2c+8]) // SMFELNM

	// DCB/DEB section (Tape)
	var SMFDCBBL interface{}
	var SMFDSSNO interface{}
	var SMF15TDA interface{}

	//DCB/DEB section (DASD)
	var SMF15NTU interface{}
	var SMF15NTR interface{}
	var SMF15NER interface{}
	var SMF15EDI interface{}
	var SMF15FG1 interface{}

	// Check if DASD dataset
	if rec[38]&0x20 != 0 {
		SMF15NTU = int(binary.BigEndian.Uint32(rec[252:256]))
		SMF15NTR = int(binary.BigEndian.Uint32(rec[256:260]))
		SMF15NER = int(rec[260])
		SMF15EDI = int(rec[261])
		SMF15FG1 = int(rec[262])
		SMFDCBBL = nil
		SMFDSSNO = nil
		SMF15TDA = nil
	} else {
		SMF15NTU = nil
		SMF15NTR = nil
		SMF15NER = nil
		SMF15EDI = nil
		SMF15FG1 = nil
		SMFDCBBL = int(binary.BigEndian.Uint32(rec[252:256]))
		bcd2asc(rec[256:262])
		SMFDSSNO = string(rec[256:262])
		SMF15TDA = int(rec[263])
	}

	var hipSize int
	if rec[39]&0x08 != 0 {
		hipSize = 20
	}
	extOff := UCBOFF + (numUcbs * ucbSize) + hipSize
	extSize := int(binary.BigEndian.Uint16(rec[extOff : extOff+2]))

	// Compressed format dataset section
	var SMF15XF1 interface{} = nil
	var SMF15CDL interface{} = nil
	var SMF15UDL interface{} = nil
	var SMF15CDS interface{} = nil
	var SMF15UDS interface{} = nil
	var SMF15CIS interface{} = nil
	var SMF15TKL interface{} = nil
	var SMF15TKN interface{} = nil
	var SMF15CMPTYPE interface{} = nil
	// SMS section
	var SMF15MCN interface{} = nil
	var SMF15DCN interface{} = nil
	var SMF15SCN interface{} = nil
	// Step information section
	var SMF15SPN interface{} = nil
	var SMF15PGN interface{} = nil
	// ISO/ANSI information section
	var SMF15CFG interface{} = nil
	var SMF15USR interface{} = nil
	var SMF15TPE interface{} = nil
	var SMF15LBL interface{} = nil
	// Additional information section
	var SMF15LBS interface{} = nil
	var SMF15FLG1 interface{} = nil
	var SMF15DSVER interface{} = nil
	var SMF15ALI interface{} = nil
	// PDSE statistics section
	var SMF15DRD interface{} = nil
	var SMF15DRDH interface{} = nil
	var SMF15MRD interface{} = nil
	var SMF15MRDH interface{} = nil
	var SMF15MCE interface{} = nil
	var SMF15MST interface{} = nil
	var SMF15MNC interface{} = nil
	var SMF15MCF interface{} = nil
	// Tape encryption data section
	var SMF15KL1 interface{} = nil
	var SMF15CD1 interface{} = nil
	var SMF15KL2 interface{} = nil
	var SMF15CD2 interface{} = nil
	var SMF15KET interface{} = nil
	// RAS section
	var SMF15RFG0 interface{} = nil
	var SMF15RFG1 interface{} = nil
	var SMF15RASDATA0 interface{} = nil
	var SMF15RASDATA1 interface{} = nil
	// DASD dataset encryption section
	var SMF15DEF interface{} = nil
	var SMF15DET interface{} = nil
	var SMF15DKL interface{} = nil

	secOff := extOff + 2
	for secOff < extOff+extSize {
		secSize := int(binary.BigEndian.Uint16(rec[secOff : secOff+2]))
		secType := int(rec[secOff+3])
		switch secType {
		case 1:
			SMF15XF1 = int(rec[secOff+4])
			SMF15CDL = int(binary.BigEndian.Uint64(rec[secOff+6 : secOff+14]))
			SMF15UDL = int(binary.BigEndian.Uint64(rec[secOff+14 : secOff+22]))
			SMF15CDS = int(binary.BigEndian.Uint64(rec[secOff+22 : secOff+30]))
			SMF15UDS = int(binary.BigEndian.Uint64(rec[secOff+30 : secOff+38]))
			SMF15CIS = int(binary.BigEndian.Uint32(rec[secOff+38 : secOff+42]))
			SMF15TKL = int(binary.BigEndian.Uint16(rec[secOff+42 : secOff+44]))
			SMF15TKN = rec[secOff+44 : secOff+80]
			if secOff+80 < secOff+secSize {
				SMF15CMPTYPE = int(rec[secOff+80])
			}
		case 2:
			bcd2asc(rec[secOff+4 : secOff+12])
			SMF15MCN = string(rec[secOff+4 : secOff+12])
			bcd2asc(rec[secOff+12 : secOff+20])
			SMF15DCN = string(rec[secOff+12 : secOff+20])
			bcd2asc(rec[secOff+20 : secOff+28])
			SMF15SCN = string(rec[secOff+20 : secOff+28])
		case 3:
			bcd2asc(rec[secOff+4 : secOff+12])
			SMF15SPN = string(rec[secOff+4 : secOff+12])
			bcd2asc(rec[secOff+12 : secOff+20])
			SMF15PGN = string(rec[secOff+12 : secOff+20])
		case 4:
			SMF15CFG = int(rec[secOff+4])
			SMF15USR = int(binary.BigEndian.Uint32(rec[secOff+5 : secOff+9]))
			SMF15TPE = int(binary.BigEndian.Uint32(rec[secOff+9 : secOff+13]))
			SMF15LBL = int(binary.BigEndian.Uint32(rec[secOff+13 : secOff+17]))
		case 5:
			if rec[secOff+4]&0x80 != 0 {
				SMF15LBS = int(binary.BigEndian.Uint64(rec[secOff+6 : secOff+14]))
			}
			if rec[secOff+4]&0x40 != 0 {
				SMF15FLG1 = int(rec[secOff+14])
			}
			if rec[secOff+4]&0x20 != 0 {
				SMF15DSVER = int(rec[secOff+16])
			}
			if rec[secOff+4]&0x10 != 0 {
				bcd2asc(rec[secOff+17 : secOff+61])
				SMF15ALI = string(rec[secOff+17 : secOff+61])
			}
		case 6:
			SMF15DRD = int(binary.BigEndian.Uint32(rec[secOff+4 : secOff+8]))
			SMF15DRDH = int(binary.BigEndian.Uint32(rec[secOff+8 : secOff+12]))
			SMF15MRD = int(binary.BigEndian.Uint32(rec[secOff+12 : secOff+16]))
			SMF15MRDH = int(binary.BigEndian.Uint32(rec[secOff+16 : secOff+20]))
			SMF15MCE = int(binary.BigEndian.Uint32(rec[secOff+20 : secOff+24]))
			SMF15MST = int(binary.BigEndian.Uint32(rec[secOff+24 : secOff+28]))
			SMF15MNC = int(binary.BigEndian.Uint32(rec[secOff+28 : secOff+32]))
			SMF15MCF = int(binary.BigEndian.Uint32(rec[secOff+32 : secOff+36]))
		case 7:
			bcd2asc(rec[secOff+4 : secOff+68])
			SMF15KL1 = string(rec[secOff+4 : secOff+68])
			bcd2asc(rec[secOff+68 : secOff+69])
			SMF15CD1 = string(rec[secOff+68 : secOff+69])
			bcd2asc(rec[secOff+69 : secOff+133])
			SMF15KL2 = string(rec[secOff+69 : secOff+133])
			bcd2asc(rec[secOff+133 : secOff+134])
			SMF15CD2 = string(rec[secOff+133 : secOff+134])
			SMF15KET = int(binary.BigEndian.Uint32(rec[secOff+134 : secOff+138]))
		case 8:
			SMF15RFG0 = int(rec[secOff+4])
			SMF15RFG1 = int(rec[secOff+5])
			if rec[secOff+4]&0x80 != 0 {
				SMF15RASDATA0 = int(rec[secOff+8])
			}
			if rec[secOff+4]&0x40 != 0 {
				SMF15RASDATA1 = int(rec[secOff+9])
			}
		case 9:
			SMF15DEF = int(rec[secOff+4])
			SMF15DET = int(binary.BigEndian.Uint16(rec[secOff+5 : secOff+7]))
			bcd2asc(rec[secOff+7 : secOff+71])
			SMF15DKL = string(rec[secOff+7 : secOff+71])
		}
		secOff += secSize
	}

	_, err = p1Stmt015.Exec(
		smftime(rec[2:6]),                        // SMF15TME
		smfdate(rec[6:10]),                       // SMF15DTE
		string(rec[10:14]),                       // SMF15SID
		string(rec[14:22]),                       // SMF15JBN
		smftime(rec[22:26]),                      // SMF15RST
		smfdate(rec[26:30]),                      // SMF15RSD
		string(rec[30:38]),                       // SMF15UID
		int(binary.BigEndian.Uint16(rec[38:40])), // SMF15RIN
		int(rec[40]),                             // SMF15SDC
		int(rec[41]),                             // SMF15NUC
		int(rec[42]),                             // SMF15SUC
		int(rec[43]),                             // SMF15SET
		smftime(rec[44:48]),                      // SMF15OPE
		int(rec[48]),                             // SMFTIOE1
		int(rec[49]),                             // SMFTIOE2
		int(rec[50]),                             // SMFTIOE3
		int(rec[51]),                             // SMFTIOE4
		string(rec[52:60]),                       // SMFTIOE5
		b3toInt(rec[60:63]),                      // SMFTIOE6
		int(rec[63]),                             // SMFTIOE7
		string(rec[JFCBOFF:JFCBOFF+44]),          // JFCBDSNM
		string(rec[JFCBOFF+0x2c:JFCBOFF+0x2c+8]),                       // JFCBELNM
		int(rec[JFCBOFF+0x34]),                                         // JFCBTSDM
		b3toInt(rec[JFCBOFF+0x35:JFCBOFF+0x35+3]),                      // JFCBDSCB
		int(binary.BigEndian.Uint16(rec[JFCBOFF+0x3a:JFCBOFF+0x3a+2])), // JFCAMSTR
		int(binary.BigEndian.Uint16(rec[JFCBOFF+0x3c:JFCBOFF+0x3c+2])), // JFCBADBF
		int(binary.BigEndian.Uint16(rec[JFCBOFF+0x3e:JFCBOFF+0x3e+2])), // JFCNLREC
		int(rec[JFCBOFF+0x40]),                                         // JFCTDSI1
		int(rec[JFCBOFF+0x41]),                                         // JFCTDSI2
		int(rec[JFCBOFF+0x42]),                                         // JFCBLTYP
		int(binary.BigEndian.Uint16(rec[JFCBOFF+0x44:JFCBOFF+0x44+2])), // JFCBFLSQ
		int(binary.BigEndian.Uint16(rec[JFCBOFF+0x46:JFCBOFF+0x46+2])), // JFCBVLSQ
		int(rec[JFCBOFF+0x4d]),                                         // JFCBFLG1
		int(rec[JFCBOFF+0x4e]),                                         // JFCBFLG2
		filedate(rec[JFCBOFF+0x50:JFCBOFF+0x50+3]),                     // JFCBCRDT
		filedate(rec[JFCBOFF+0x53:JFCBOFF+0x53+3]),                     // JFCBXPDT
		int(rec[JFCBOFF+0x56]),                                         // JFCBIND1
		int(rec[JFCBOFF+0x57]),                                         // JFCBIND2
		int(rec[JFCBOFF+0x62]),                                         // JFCDSRG1
		int(rec[JFCBOFF+0x63]),                                         // JFCDSRG2
		int(rec[JFCBOFF+0x64]),                                         // JFCRECFM
		int(binary.BigEndian.Uint16(rec[JFCBOFF+0x66:JFCBOFF+0x66+2])), // JFCBLKSI
		int(binary.BigEndian.Uint16(rec[JFCBOFF+0x68:JFCBOFF+0x68+2])), // JFCLRECL
		b3toInt(rec[JFCBOFF+0x98:JFCBOFF+0x98+3]),                      // JFCBPQTY
		int(rec[JFCBOFF+0x9B]),                                         // JFCBCTRI
		b3toInt(rec[JFCBOFF+0x9c:JFCBOFF+0x9c+3]),                      // JFCBSQTY
		b3toInt(rec[JFCBOFF+0xa0:JFCBOFF+0xa0+3]),                      // JFCBDQTY
		int(binary.BigEndian.Uint16(rec[240:242])),                     // SMFDCBOR
		int(rec[242]), // SMFDCBRF
		int(binary.BigEndian.Uint16(rec[243:245])), // SMFDCBMF
		int(rec[245]), // SMFDCBFL
		int(rec[246]), // SMFDCBOP
		int(rec[248]), // SMFDEBFL
		int(rec[249]), // SMFDEBOP
		int(binary.BigEndian.Uint16(rec[250:252])), // SMFDEBVL
		SMFDCBBL,
		SMFDSSNO,
		SMF15TDA,
		SMF15NTU,
		SMF15NTR,
		SMF15NER,
		SMF15EDI,
		SMF15FG1,
		smfdate(rec[264:268]), // SMF15OPD
		SMF15XF1,
		SMF15CDL,
		SMF15UDL,
		SMF15CDS,
		SMF15UDS,
		SMF15CIS,
		SMF15TKL,
		SMF15TKN,
		SMF15CMPTYPE,
		SMF15MCN,
		SMF15DCN,
		SMF15SCN,
		SMF15SPN,
		SMF15PGN,
		SMF15CFG,
		SMF15USR,
		SMF15TPE,
		SMF15LBL,
		SMF15LBS,
		SMF15FLG1,
		SMF15DSVER,
		SMF15ALI,
		SMF15DRD,
		SMF15DRDH,
		SMF15MRD,
		SMF15MRDH,
		SMF15MCE,
		SMF15MST,
		SMF15MNC,
		SMF15MCF,
		SMF15KL1,
		SMF15CD1,
		SMF15KL2,
		SMF15CD2,
		SMF15KET,
		SMF15RFG0,
		SMF15RFG1,
		SMF15RASDATA0,
		SMF15RASDATA1,
		SMF15DEF,
		SMF15DET,
		SMF15DKL,
	)
	if err != nil {
		log.Println(err)
	}

	var SMF15NTA interface{}
	var SMFSRTEF interface{}
	var SMFSRTEQ interface{}

	tempOff := UCBOFF
	for i := 0; i < numUcbs; i++ {
		bcd2asc(rec[tempOff+2 : tempOff+8])
		// Check if DASD dataset
		if rec[38]&0x20 != 0 {
			SMF15NTA = int(binary.BigEndian.Uint32(rec[tempOff+20 : tempOff+24]))
			SMFSRTEF = nil
			SMFSRTEQ = nil
		} else {
			SMF15NTA = nil
			SMFSRTEF = int(binary.BigEndian.Uint16(rec[tempOff+20 : tempOff+22]))
			SMFSRTEQ = int(binary.BigEndian.Uint16(rec[tempOff+22 : tempOff+24]))
		}

		_, err = p2Stmt015.Exec(
			smftime(rec[2:6]),  // SMF4TME
			smfdate(rec[6:10]), // SMF4DTE
			int(binary.BigEndian.Uint16(rec[tempOff+0:tempOff+2])),   // SMFUCBDV
			string(rec[tempOff+2:tempOff+8]),                         // SMFSRTEV
			int(binary.BigEndian.Uint32(rec[tempOff+8:tempOff+12])),  // SMFUCBTY
			int(rec[tempOff+12]),                                     // SMFSRTES
			int(rec[tempOff+13]),                                     // SMF15NEX
			int(binary.BigEndian.Uint32(rec[tempOff+16:tempOff+20])), // SMFEXCP
			SMFSRTEF,
			SMFSRTEQ,
			SMF15NTA,
		)
		if err != nil {
			log.Println(err)
		}
		tempOff += ucbSize
	}
}
