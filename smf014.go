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
		SMFDSNM    TEXT       NOT NULL,
		SMFELNM    TEXT       NOT NULL,
		SMFDCBOR   INT        NOT NULL,
		SMFDCBRF   INT        NOT NULL,
		SMFDCBMF   INT        NOT NULL,
		SMFDCBFL   INT        NOT NULL,
		SMFDCBOP   INT        NOT NULL,
		SMFDEBFL   INT        NOT NULL,
		SMFDEBOP   INT        NOT NULL,
		SMFDEBVL   INT        NOT NULL,
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
		string(rec[jfcboff:jfcboff+44]),          // SMFDSNM
		string(rec[jfcboff+0x2c:jfcboff+0x2c+8]),   // SMFELNM
		int(binary.BigEndian.Uint16(rec[240:242])), // SMFDCBOR
		int(rec[242]), // SMFDCBRF
		int(binary.BigEndian.Uint16(rec[243:245])), // SMFDCBMF
		int(rec[245]), // SMFDCBFL
		int(rec[246]), // SMFDCBOP
		int(rec[248]), // SMFDEBFL
		int(rec[249]), // SMFDEBOP
		int(binary.BigEndian.Uint16(rec[250:252])), // SMFDEBVL
		smfdate(rec[264:268]),                      // SMF14OPD
		// Start of relocate section
	)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}
