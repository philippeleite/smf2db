package main

import (
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
	//"strings"
)

func smf014Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf014;
	CREATE TABLE smf014 (
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
		SMFDCBBL   INT        NOT NULL,
		SMFDSSNO   TEXT       NOT NULL,
		SMF14TDA   INT        NOT NULL,
		SMF14NTU   INT        NOT NULL,
		SMF14NTR   INT        NOT NULL,
		SMF14NER   INT        NOT NULL,
		SMF14EDI   INT        NOT NULL,
		SMF14FG1   INT        NOT NULL,
		SMF14OPD   DATETIME   NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}

func smf014Parse(rec []byte) {
	sqlStmt := `
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
			?
		);
	`
	prepStmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	defer prepStmt.Close()

	jfcboff := 64
	bcd2asc(rec[10:14])                         // SMF14SID
	bcd2asc(rec[14:22])                         // SMF14JBN
	bcd2asc(rec[30:38])                         // SMF14UID
	bcd2asc(rec[52:60])                         // SMFTIOE5
	bcd2asc(rec[jfcboff : jfcboff+44])          // SMFDSNM
	bcd2asc(rec[jfcboff+0x2c : jfcboff+0x2c+8]) // SMFELNM
	// reloc := 93 + int(rec[92])

	//var acctFields strings.Builder
	//acctFields.WriteString(" ")
	//fieldOff := 117
	//numFields := int(rec[116])
	//for i := 0; i < numFields; i++ {
	//	fieldLen := int(rec[fieldOff])
	//	bcd2asc(rec[fieldOff+1 : fieldOff+1+fieldLen])
	//	acctFields.WriteString(string(rec[fieldOff+1:fieldOff+1+fieldLen]) + " ")
	//	fieldOff += fieldLen + 1
	//}

	// DCB/DEB section (Tape)
	var SMFDCBBL int
	var SMFDSSNO string
	var SMF14TDA int

	//DCB/DEB section (DASD)
	var SMF14NTU int
	var SMF14NTR int
	var SMF14NER int
	var SMF14EDI int
	var SMF14FG1 int

	// Check if DASD dataset
	if rec[38]&0x20 != 0 {
		SMF14NTU = int(binary.BigEndian.Uint32(rec[252:256]))
		SMF14NTR = int(binary.BigEndian.Uint32(rec[256:260]))
		SMF14NER = int(rec[260])
		SMF14EDI = int(rec[261])
		SMF14FG1 = int(rec[262])
	} else {
		SMFDCBBL = int(binary.BigEndian.Uint32(rec[252:256]))
		bcd2asc(rec[256:262])
		SMFDSSNO = string(rec[256:262])
		SMF14TDA = int(rec[263])
	}
	_, err = prepStmt.Exec(
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
		string(rec[jfcboff:jfcboff+44]),          // JFCBDSNM
		string(rec[jfcboff+0x2c:jfcboff+0x2c+8]),                       // JFCBELNM
		int(rec[jfcboff+0x34]),                                         // JFCBTSDM
		b3toInt(rec[jfcboff+0x35:jfcboff+0x35+3]),                      // JFCBDSCB
		int(binary.BigEndian.Uint16(rec[jfcboff+0x3a:jfcboff+0x3a+2])), // JFCAMSTR
		int(binary.BigEndian.Uint16(rec[jfcboff+0x3c:jfcboff+0x3c+2])), // JFCBADBF
		int(binary.BigEndian.Uint16(rec[jfcboff+0x3e:jfcboff+0x3e+2])), // JFCNLREC
		int(rec[jfcboff+0x40]),                                         // JFCTDSI1
		int(rec[jfcboff+0x41]),                                         // JFCTDSI2
		int(rec[jfcboff+0x42]),                                         // JFCBLTYP
		int(binary.BigEndian.Uint16(rec[jfcboff+0x44:jfcboff+0x44+2])), // JFCBFLSQ
		int(binary.BigEndian.Uint16(rec[jfcboff+0x46:jfcboff+0x46+2])), // JFCBVLSQ
		int(rec[jfcboff+0x4d]),                                         // JFCBFLG1
		int(rec[jfcboff+0x4e]),                                         // JFCBFLG2
		filedate(rec[jfcboff+0x50:jfcboff+0x50+3]),                     // JFCBCRDT
		filedate(rec[jfcboff+0x53:jfcboff+0x53+3]),                     // JFCBXPDT
		int(rec[jfcboff+0x56]),                                         // JFCBIND1
		int(rec[jfcboff+0x57]),                                         // JFCBIND2
		int(rec[jfcboff+0x62]),                                         // JFCDSRG1
		int(rec[jfcboff+0x63]),                                         // JFCDSRG2
		int(rec[jfcboff+0x64]),                                         // JFCRECFM
		int(binary.BigEndian.Uint16(rec[jfcboff+0x66:jfcboff+0x66+2])), // JFCBLKSI
		int(binary.BigEndian.Uint16(rec[jfcboff+0x68:jfcboff+0x68+2])), // JFCLRECL
		b3toInt(rec[jfcboff+0x98:jfcboff+0x98+3]),                      // JFCBPQTY
		int(rec[jfcboff+0x9B]),                                         // JFCBCTRI
		b3toInt(rec[jfcboff+0x9c:jfcboff+0x9c+3]),                      // JFCBSQTY
		b3toInt(rec[jfcboff+0xa0:jfcboff+0xa0+3]),                      // JFCBDQTY
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
		// Start of relocate section
	)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}
