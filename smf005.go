package main

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"sync"
)

var pStmt005 *sql.Stmt
var once005 sync.Once

func smf005Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf005;
	CREATE TABLE smf005 (
		SMF5TME    DATETIME   NOT NULL,
		SMF5DTE    DATETIME   NOT NULL,
		SMF5SID    TEXT       NOT NULL,
		SMF5JBN    TEXT       NOT NULL,
		SMF5RST    DATETIME   NOT NULL,
		SMF5RSD    DATETIME   NOT NULL,
		SMF5UIF    TEXT       NOT NULL,
		SMF5NST    INT        NOT NULL,
		SMF5JIT    DATETIME   NOT NULL,
		SMF5JID    DATETIME   NOT NULL,
		SMF5NCI    INT        NOT NULL,
		SMF5JCC    INT        NOT NULL,
		SMF5JPTY   INT        NOT NULL,
		SMF5RSTT   DATETIME   NOT NULL,
		SMF5RSTD   DATETIME   NOT NULL,
		SMF5JBTI   INT        NOT NULL,
		SMF5TRAN   INT        NOT NULL,
		SMF5RDCL   INT        NOT NULL,
		SMF5RUTY   INT        NOT NULL,
		SMF5JICL   TEXT       NOT NULL,
		SMF5SPK    INT        NOT NULL,
		SMF5SRBT   INT        NOT NULL,
		SMF5TJS    INT        NOT NULL,
		SMF5TTAT   INT        NOT NULL,
		SMF5PRGN   TEXT       NOT NULL,
		SMF5ACTF   INT        NOT NULL,
		SMF5JSAF   INT        NOT NULL,
		SMF5CPUS   INT        NOT NULL,
		SMF5IOCS   INT        NOT NULL,
		SMF5MSOS   INT        NOT NULL,
		SMF5SRBS   INT        NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
		INSERT INTO smf005
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
	pStmt005, err = db.Prepare(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
	}
}

func smf005Parse(rec []byte) {
	once005.Do(smf005Init)

	bcd2asc(rec[10:14])  // SMF5SID
	bcd2asc(rec[14:22])  // SMF5JBN
	bcd2asc(rec[30:38])  // SMF5UIF
	bcd2asc(rec[71:72])  // SMF5JICL
	bcd2asc(rec[93:113]) // SMF5PRGN
	reloc := 93 + int(rec[92])

	var acctFields strings.Builder
	acctFields.WriteString(" ")
	fieldOff := 117
	numFields := int(rec[116])
	for i := 0; i < numFields; i++ {
		fieldLen := int(rec[fieldOff])
		bcd2asc(rec[fieldOff+1 : fieldOff+1+fieldLen])
		acctFields.WriteString(string(rec[fieldOff+1:fieldOff+1+fieldLen]) + " ")
		fieldOff += fieldLen + 1
	}
	_, err = pStmt005.Exec(
		smftime(rec[2:6]),                        // SMF5TME
		smfdate(rec[6:10]),                       // SMF5DTE
		string(rec[10:14]),                       // SMF5SID
		string(rec[14:22]),                       // SMF5JBN
		smftime(rec[22:26]),                      // SMF5RST
		smfdate(rec[26:30]),                      // SMF5RSD
		string(rec[30:38]),                       // SMF5UIF
		int(rec[38]),                             // SMF5NST
		smftime(rec[39:43]),                      // SMF5JIT
		smfdate(rec[43:47]),                      // SMF5JID
		int(binary.BigEndian.Uint32(rec[47:51])), // SMF5NCI
		int(binary.BigEndian.Uint16(rec[51:53])), // SMF5JCC
		int(rec[53]),                             // SMF5JPTY
		smftime(rec[54:58]),                      // SMF5RSTT
		smfdate(rec[58:62]),                      // SMF5RSTD
		int(rec[62]),                             // SMF5JBTI
		int(binary.BigEndian.Uint32(rec[64:68])), // SMF5TRAN
		int(rec[69]),                             // SMF5RDCL
		int(rec[70]),                             // SMF5RUTY
		string(rec[71:72]),                       // SMF5JICL
		int(rec[72]),                             // SMF5SPK
		b3toInt(rec[73:76]),                      // SMF5SRBT
		int(binary.BigEndian.Uint32(rec[76:80])), // SMF5TJS
		int(binary.BigEndian.Uint32(rec[80:84])), // SMF5TTAT
		string(rec[93:113]),                      // SMF5PRGN
		int(rec[116]),                            // SMF5ACTF
		acctFields.String(),                      // SMF5JSAF
		// Start of relocate section
		int(binary.BigEndian.Uint32(rec[reloc+0:reloc+4])),   // SMF5CPUS
		int(binary.BigEndian.Uint32(rec[reloc+4:reloc+8])),   // SMF5IOCS
		int(binary.BigEndian.Uint32(rec[reloc+8:reloc+12])),  // SMF5MSOS
		int(binary.BigEndian.Uint32(rec[reloc+12:reloc+16])), // SMF5SRBS
	)
	if err != nil {
		log.Println(err)
	}
}
