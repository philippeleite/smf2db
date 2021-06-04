package main

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
	//"strings"
)

var p1Stmt014 *sql.Stmt
var p2Stmt014 *sql.Stmt
var once014 sync.Once
var count014 int

func smf014Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf014;
	CREATE TABLE smf014 (
	    ID         INT        NOT NULL,
		SMF14TME   DATETIME   NOT NULL,
		SMF14DTE   DATETIME   NOT NULL,
		SMF14SID   TEXT       NOT NULL,
		SMF14JBN   TEXT       NOT NULL,
		SMF14RST   DATETIME   NOT NULL,
		SMF14RSD   DATETIME   NOT NULL,
		SMF14UID   TEXT       NOT NULL,
		SMF14RIN   INT        NOT NULL,
		SMF14SDC   INT        NOT NULL,
		SMF14NUC   INT        NOT NULL,
		SMF14SUC   INT        NOT NULL,
		SMF14SET   INT        NOT NULL,
		SMF14OPE   DATETIME   NOT NULL,
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
		SMF14TDA   INT        ,
		SMF14NTU   INT        ,
		SMF14NTR   INT        ,
		SMF14NER   INT        ,
		SMF14EDI   INT        ,
		SMF14FG1   INT        ,
		SMF14OPD   DATETIME   NOT NULL,
		SMF14XF1   INT        ,
		SMF14CDL   INT        ,
		SMF14UDL   INT        ,
		SMF14CDS   INT        ,
		SMF14UDS   INT        ,
		SMF14CIS   INT        ,
		SMF14TKL   INT        ,
		SMF14TKN   BLOB       ,
		SMF14CMPTYPE INT      ,
		SMF14MCN   TEXT       ,
		SMF14DCN   TEXT       ,
		SMF14SCN   TEXT       ,
		SMF14SPN   TEXT       ,
		SMF14PGN   TEXT       ,
		SMF14CFG   INT        ,
		SMF14USR   INT        ,
		SMF14TPE   INT        ,
		SMF14LBL   INT        ,
		SMF14LBS   INT        ,        
		SMF14FLG1  INT        ,        
		SMF14DSVER INT        ,        
		SMF14ALI   TEXT       ,
		SMF14DRD   INT        ,        
		SMF14DRDH  INT        ,        
		SMF14MRD   INT        ,        
		SMF14MRDH  INT        ,        
		SMF14MCE   INT        ,        
		SMF14MST   INT        ,        
		SMF14MNC   INT        ,        
		SMF14MCF   INT        ,
		SMF14KL1   TEXT       ,        
		SMF12CD1   TEXT       ,        
		SMF14KL2   TEXT       ,        
		SMF14CD2   TEXT       ,        
		SMF14KET   INT        ,        
		SMF14RFG0  INT        ,        
		SMF14RFG1  INT        ,        
		SMF14RASDATA0 INT     ,        
		SMF14RASDATA1 INT     ,
		SMF14DEF   INT        ,
		SMF14DET   INT        ,        
		SMF14DKL   INT                
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
		INSERT INTO smf014
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
			?,
			?
		);
	`
	p1Stmt014, err = db.Prepare(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
	}

	sqlStmt3 := `
	DROP TABLE IF EXISTS smf014_ucb;
	CREATE TABLE smf014_ucb (
	    ID         INT        NOT NULL,
		SMFUCBDV   INT        NOT NULL,
		SMFSRTEV   TEXT       NOT NULL,
		SMFUCBTY   INT        NOT NULL,
		SMFSRTES   INT        NOT NULL,
		SMF14NEX   INT        NOT NULL,
		SMFEXCP    INT        NOT NULL,
		SMFSRTEF   INT        ,
		SMFSRTEQ   INT        ,
		SMF14NTA   INT        
	);
	`
	_, err = db.Exec(sqlStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt3)
	}

	sqlStmt4 := `
		INSERT INTO smf014_ucb
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
			?
		);
	`
	p2Stmt014, err = db.Prepare(sqlStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt4)
	}
}

func smf014Parse(rec []byte) {
	once014.Do(smf014Init)

	const (
		JFCBOFF = 64
		UCBOFF  = 268
	)

	numUcbs := int(rec[41])
	ucbSize := int(rec[42])

	bcd2asc(rec[10:14])                         // SMF14SID
	bcd2asc(rec[14:22])                         // SMF14JBN
	bcd2asc(rec[30:38])                         // SMF14UID
	bcd2asc(rec[52:60])                         // SMFTIOE5
	bcd2asc(rec[JFCBOFF : JFCBOFF+44])          // SMFDSNM
	bcd2asc(rec[JFCBOFF+0x2c : JFCBOFF+0x2c+8]) // SMFELNM

	// DCB/DEB section (Tape)
	var SMFDCBBL interface{}
	var SMFDSSNO interface{}
	var SMF14TDA interface{}

	//DCB/DEB section (DASD)
	var SMF14NTU interface{}
	var SMF14NTR interface{}
	var SMF14NER interface{}
	var SMF14EDI interface{}
	var SMF14FG1 interface{}

	// Check if DASD dataset
	if rec[38]&0x20 != 0 {
		SMF14NTU = int(binary.BigEndian.Uint32(rec[252:256]))
		SMF14NTR = int(binary.BigEndian.Uint32(rec[256:260]))
		SMF14NER = int(rec[260])
		SMF14EDI = int(rec[261])
		SMF14FG1 = int(rec[262])
		SMFDCBBL = nil
		SMFDSSNO = nil
		SMF14TDA = nil
	} else {
		SMF14NTU = nil
		SMF14NTR = nil
		SMF14NER = nil
		SMF14EDI = nil
		SMF14FG1 = nil
		SMFDCBBL = int(binary.BigEndian.Uint32(rec[252:256]))
		bcd2asc(rec[256:262])
		SMFDSSNO = string(rec[256:262])
		SMF14TDA = int(rec[263])
	}

	var hipSize int
	if rec[39]&0x08 != 0 {
		hipSize = 20
	}
	extOff := UCBOFF + (numUcbs * ucbSize) + hipSize
	extSize := int(binary.BigEndian.Uint16(rec[extOff : extOff+2]))

	// Compressed format dataset section
	var SMF14XF1 interface{} = nil
	var SMF14CDL interface{} = nil
	var SMF14UDL interface{} = nil
	var SMF14CDS interface{} = nil
	var SMF14UDS interface{} = nil
	var SMF14CIS interface{} = nil
	var SMF14TKL interface{} = nil
	var SMF14TKN interface{} = nil
	var SMF14CMPTYPE interface{} = nil
	// SMS section
	var SMF14MCN interface{} = nil
	var SMF14DCN interface{} = nil
	var SMF14SCN interface{} = nil
	// Step information section
	var SMF14SPN interface{} = nil
	var SMF14PGN interface{} = nil
	// ISO/ANSI information section
	var SMF14CFG interface{} = nil
	var SMF14USR interface{} = nil
	var SMF14TPE interface{} = nil
	var SMF14LBL interface{} = nil
	// Additional information section
	var SMF14LBS interface{} = nil
	var SMF14FLG1 interface{} = nil
	var SMF14DSVER interface{} = nil
	var SMF14ALI interface{} = nil
	// PDSE statistics section
	var SMF14DRD interface{} = nil
	var SMF14DRDH interface{} = nil
	var SMF14MRD interface{} = nil
	var SMF14MRDH interface{} = nil
	var SMF14MCE interface{} = nil
	var SMF14MST interface{} = nil
	var SMF14MNC interface{} = nil
	var SMF14MCF interface{} = nil
	// Tape encryption data section
	var SMF14KL1 interface{} = nil
	var SMF14CD1 interface{} = nil
	var SMF14KL2 interface{} = nil
	var SMF14CD2 interface{} = nil
	var SMF14KET interface{} = nil
	// RAS section
	var SMF14RFG0 interface{} = nil
	var SMF14RFG1 interface{} = nil
	var SMF14RASDATA0 interface{} = nil
	var SMF14RASDATA1 interface{} = nil
	// DASD dataset encryption section
	var SMF14DEF interface{} = nil
	var SMF14DET interface{} = nil
	var SMF14DKL interface{} = nil

	secOff := extOff + 2
	for secOff < extOff+extSize {
		secSize := int(binary.BigEndian.Uint16(rec[secOff : secOff+2]))
		secType := int(rec[secOff+3])
		switch secType {
		case 1:
			SMF14XF1 = int(rec[secOff+4])
			SMF14CDL = int(binary.BigEndian.Uint64(rec[secOff+6 : secOff+14]))
			SMF14UDL = int(binary.BigEndian.Uint64(rec[secOff+14 : secOff+22]))
			SMF14CDS = int(binary.BigEndian.Uint64(rec[secOff+22 : secOff+30]))
			SMF14UDS = int(binary.BigEndian.Uint64(rec[secOff+30 : secOff+38]))
			SMF14CIS = int(binary.BigEndian.Uint32(rec[secOff+38 : secOff+42]))
			SMF14TKL = int(binary.BigEndian.Uint16(rec[secOff+42 : secOff+44]))
			SMF14TKN = rec[secOff+44 : secOff+80]
			if secOff+80 < secOff+secSize {
				SMF14CMPTYPE = int(rec[secOff+80])
			}
		case 2:
			bcd2asc(rec[secOff+4 : secOff+12])
			SMF14MCN = string(rec[secOff+4 : secOff+12])
			bcd2asc(rec[secOff+12 : secOff+20])
			SMF14DCN = string(rec[secOff+12 : secOff+20])
			bcd2asc(rec[secOff+20 : secOff+28])
			SMF14SCN = string(rec[secOff+20 : secOff+28])
		case 3:
			bcd2asc(rec[secOff+4 : secOff+12])
			SMF14SPN = string(rec[secOff+4 : secOff+12])
			bcd2asc(rec[secOff+12 : secOff+20])
			SMF14PGN = string(rec[secOff+12 : secOff+20])
		case 4:
			SMF14CFG = int(rec[secOff+4])
			SMF14USR = int(binary.BigEndian.Uint32(rec[secOff+5 : secOff+9]))
			SMF14TPE = int(binary.BigEndian.Uint32(rec[secOff+9 : secOff+13]))
			SMF14LBL = int(binary.BigEndian.Uint32(rec[secOff+13 : secOff+17]))
		case 5:
			if rec[secOff+4]&0x80 != 0 {
				SMF14LBS = int(binary.BigEndian.Uint64(rec[secOff+6 : secOff+14]))
			}
			if rec[secOff+4]&0x40 != 0 {
				SMF14FLG1 = int(rec[secOff+14])
			}
			if rec[secOff+4]&0x20 != 0 {
				SMF14DSVER = int(rec[secOff+16])
			}
			if rec[secOff+4]&0x10 != 0 {
				bcd2asc(rec[secOff+17 : secOff+61])
				SMF14ALI = string(rec[secOff+17 : secOff+61])
			}
		case 6:
			SMF14DRD = int(binary.BigEndian.Uint32(rec[secOff+4 : secOff+8]))
			SMF14DRDH = int(binary.BigEndian.Uint32(rec[secOff+8 : secOff+12]))
			SMF14MRD = int(binary.BigEndian.Uint32(rec[secOff+12 : secOff+16]))
			SMF14MRDH = int(binary.BigEndian.Uint32(rec[secOff+16 : secOff+20]))
			SMF14MCE = int(binary.BigEndian.Uint32(rec[secOff+20 : secOff+24]))
			SMF14MST = int(binary.BigEndian.Uint32(rec[secOff+24 : secOff+28]))
			SMF14MNC = int(binary.BigEndian.Uint32(rec[secOff+28 : secOff+32]))
			SMF14MCF = int(binary.BigEndian.Uint32(rec[secOff+32 : secOff+36]))
		case 7:
			bcd2asc(rec[secOff+4 : secOff+68])
			SMF14KL1 = string(rec[secOff+4 : secOff+68])
			bcd2asc(rec[secOff+68 : secOff+69])
			SMF14CD1 = string(rec[secOff+68 : secOff+69])
			bcd2asc(rec[secOff+69 : secOff+133])
			SMF14KL2 = string(rec[secOff+69 : secOff+133])
			bcd2asc(rec[secOff+133 : secOff+134])
			SMF14CD2 = string(rec[secOff+133 : secOff+134])
			SMF14KET = int(binary.BigEndian.Uint32(rec[secOff+134 : secOff+138]))
		case 8:
			SMF14RFG0 = int(rec[secOff+4])
			SMF14RFG1 = int(rec[secOff+5])
			if rec[secOff+4]&0x80 != 0 {
				SMF14RASDATA0 = int(rec[secOff+8])
			}
			if rec[secOff+4]&0x40 != 0 {
				SMF14RASDATA1 = int(rec[secOff+9])
			}
		case 9:
			SMF14DEF = int(rec[secOff+4])
			SMF14DET = int(binary.BigEndian.Uint16(rec[secOff+5 : secOff+7]))
			bcd2asc(rec[secOff+7 : secOff+71])
			SMF14DKL = string(rec[secOff+7 : secOff+71])
		}
		secOff += secSize
	}

	count014 += 1
	_, err = p1Stmt014.Exec(
		count014,
		smftime(rec[2:6]),                        // SMF14TME
		smfdate(rec[6:10]),                       // SMF14DTE
		string(rec[10:14]),                       // SMF14SID
		string(rec[14:22]),                       // SMF14JBN
		smftime(rec[22:26]),                      // SMF14RST
		smfdate(rec[26:30]),                      // SMF14RSD
		string(rec[30:38]),                       // SMF14UID
		int(binary.BigEndian.Uint16(rec[38:40])), // SMF14RIN
		int(rec[40]),                             // SMF14SDC
		int(rec[41]),                             // SMF14NUC
		int(rec[42]),                             // SMF14SUC
		int(rec[43]),                             // SMF14SET
		smftime(rec[44:48]),                      // SMF14OPE
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
		SMF14TDA,
		SMF14NTU,
		SMF14NTR,
		SMF14NER,
		SMF14EDI,
		SMF14FG1,
		smfdate(rec[264:268]), // SMF14OPD
		SMF14XF1,
		SMF14CDL,
		SMF14UDL,
		SMF14CDS,
		SMF14UDS,
		SMF14CIS,
		SMF14TKL,
		SMF14TKN,
		SMF14CMPTYPE,
		SMF14MCN,
		SMF14DCN,
		SMF14SCN,
		SMF14SPN,
		SMF14PGN,
		SMF14CFG,
		SMF14USR,
		SMF14TPE,
		SMF14LBL,
		SMF14LBS,
		SMF14FLG1,
		SMF14DSVER,
		SMF14ALI,
		SMF14DRD,
		SMF14DRDH,
		SMF14MRD,
		SMF14MRDH,
		SMF14MCE,
		SMF14MST,
		SMF14MNC,
		SMF14MCF,
		SMF14KL1,
		SMF14CD1,
		SMF14KL2,
		SMF14CD2,
		SMF14KET,
		SMF14RFG0,
		SMF14RFG1,
		SMF14RASDATA0,
		SMF14RASDATA1,
		SMF14DEF,
		SMF14DET,
		SMF14DKL,
	)
	if err != nil {
		log.Println(err)
	}

	var SMF14NTA interface{}
	var SMFSRTEF interface{}
	var SMFSRTEQ interface{}

	tempOff := UCBOFF
	for i := 0; i < numUcbs; i++ {
		bcd2asc(rec[tempOff+2 : tempOff+8])
		// Check if DASD dataset
		if rec[38]&0x20 != 0 {
			SMF14NTA = int(binary.BigEndian.Uint32(rec[tempOff+20 : tempOff+24]))
			SMFSRTEF = nil
			SMFSRTEQ = nil
		} else {
			SMF14NTA = nil
			SMFSRTEF = int(binary.BigEndian.Uint16(rec[tempOff+20 : tempOff+22]))
			SMFSRTEQ = int(binary.BigEndian.Uint16(rec[tempOff+22 : tempOff+24]))
		}

		_, err = p2Stmt014.Exec(
			count014,
			int(binary.BigEndian.Uint16(rec[tempOff+0:tempOff+2])),   // SMFUCBDV
			string(rec[tempOff+2:tempOff+8]),                         // SMFSRTEV
			int(binary.BigEndian.Uint32(rec[tempOff+8:tempOff+12])),  // SMFUCBTY
			int(rec[tempOff+12]),                                     // SMFSRTES
			int(rec[tempOff+13]),                                     // SMF14NEX
			int(binary.BigEndian.Uint32(rec[tempOff+16:tempOff+20])), // SMFEXCP
			SMFSRTEF,
			SMFSRTEQ,
			SMF14NTA,
		)
		if err != nil {
			log.Println(err)
		}
		tempOff += ucbSize
	}
}
