package main

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
	"sync"
)

var p1Stmt004 *sql.Stmt
var p2Stmt004 *sql.Stmt
var once004 sync.Once
var count004 int

func smf004Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf004;
	CREATE TABLE smf004 (
	    ID         INT        NOT NULL,
		SMF4TME    DATETIME   NOT NULL,
		SMF4DTE    DATETIME   NOT NULL,
		SMF4SID    TEXT       NOT NULL,
		SMF4JBN    TEXT       NOT NULL,
		SMF4RST    DATETIME   NOT NULL,
		SMF4RSD    DATETIME   NOT NULL,
		SMF4UIF    TEXT       NOT NULL,
		SMF4STN    INT        NOT NULL,
		SMF4SIT    DATETIME   NOT NULL,
		SMF4STID   DATETIME   NOT NULL,
		SMF4NCI    INT        NOT NULL,
		SMF4SCC    INT        NOT NULL,
		SMF4PRTY   INT        NOT NULL,
		SMF4PGMN   TEXT       NOT NULL,
		SMF4STMN   TEXT       NOT NULL,
		SMF4SYST   INT        NOT NULL,
		SMF4HOST   INT        NOT NULL,
		SMF4RSH0   INT        NOT NULL,
		SMF4SPK    INT        NOT NULL,
		SMF4STI    INT        NOT NULL,
		SMF4AST    INT        NOT NULL,
		SMF4PPST   INT        NOT NULL,
		SMF4SRBT   INT        NOT NULL,
		SMF4RIN    INT        NOT NULL,
		SMF4PGIN   INT        NOT NULL,
		SMF4PGOT   INT        NOT NULL,
		SMF4NSW    INT        NOT NULL,
		SMF4PSI    INT        NOT NULL,
		SMF4PSO    INT        NOT NULL,
		SMF4VPI    INT        NOT NULL,
		SMF4VPO    INT        NOT NULL,
		SMF4SST    INT        NOT NULL,
		SMF4ACT    INT        NOT NULL,
		SMF4PGNO   INT        NOT NULL,
		SMF4TRAN   INT        NOT NULL,
		SMF4CPM    INT        NOT NULL,
		SMF4RCLM   INT        NOT NULL,
		SMF4CPGN   INT        NOT NULL,
		SMF4HSPI   INT        NOT NULL,
		SMF4PGST   INT        NOT NULL,
		SMF4PSEC   INT        NOT NULL,
		SMF4LPAI   INT        NOT NULL,
		SMF4HSPO   INT        NOT NULL,
		SMF4CPUS   INT        NOT NULL,
		SMF4IOCS   INT        NOT NULL,
		SMF4MSOS   INT        NOT NULL,
		SMF4SRBS   INT        NOT NULL,
		SMF4SETM   INT        NOT NULL,
		SMF4ACTF   TEXT       NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
	DROP TABLE IF EXISTS smf004_dev;
	CREATE TABLE smf004_dev (
	    ID         INT        NOT NULL,
		SMF4DEVC   INT        NOT NULL,
		SMF4UTYP   INT        NOT NULL,
		SMF4CUAD   INT        NOT NULL,
		SMF4EXCP   INT        NOT NULL
	);
	`
	_, err = db.Exec(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
	}

	sqlStmt3 := `
		INSERT INTO smf004
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
	p1Stmt004, err = db.Prepare(sqlStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt3)
	}

	sqlStmt4 := `
		INSERT INTO smf004_dev
		VALUES (
			?,
			?,
			?,
			?,
			?
		);
`

	p2Stmt004, err = db.Prepare(sqlStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt4)
	}
}

func smf004Parse(rec []byte) {
	once004.Do(smf004Init)

	bcd2asc(rec[10:14]) // SMF4SID
	bcd2asc(rec[14:22]) // SMF4JBN
	bcd2asc(rec[30:38]) // SMF4UIF
	bcd2asc(rec[54:62]) // SMF4PGMN
	bcd2asc(rec[62:70]) // SMF4STMN

	reloc := int(binary.BigEndian.Uint16(rec[100:102]))
	numDevs := (int(binary.BigEndian.Uint16(rec[102:104])) - 2) / 8

	acctOff := 104 + (numDevs * 8)

	var acctFields strings.Builder
	acctFields.WriteString(" ")
	numFields := int(rec[acctOff+4])
	fieldOff := acctOff + 5
	for i := 0; i < numFields; i++ {
		fieldLen := int(rec[fieldOff])
		bcd2asc(rec[fieldOff+1 : fieldOff+1+fieldLen])
		acctFields.WriteString(string(rec[fieldOff+1:fieldOff+1+fieldLen]) + " ")
		fieldOff += fieldLen + 1
	}

	count004 += 1
	_, err = p1Stmt004.Exec(
		count004,
		smftime(rec[2:6]),                         // SMF4TME
		smfdate(rec[6:10]),                        // SMF4DTE
		string(rec[10:14]),                        // SMF4SID
		string(rec[14:22]),                        // SMF4JBN
		smftime(rec[22:26]),                       // SMF4RST
		smfdate(rec[26:30]),                       // SMF4RSD
		string(rec[30:38]),                        // SMF4UIF
		int(rec[38]),                              // SMF4STN
		smftime(rec[39:43]),                       // SMF4SIT
		smfdate(rec[43:47]),                       // SMF4STID
		int(binary.BigEndian.Uint32(rec[47:51])),  // SMF4NCI
		int(binary.BigEndian.Uint16(rec[51:53])),  // SMF4SCC
		int(rec[53]),                              // SMF4PRTY
		string(rec[54:62]),                        // SMF4PGMN
		string(rec[62:70]),                        // SMF4STMN
		int(binary.BigEndian.Uint16(rec[72:74])),  // SMF4SYST
		int(binary.BigEndian.Uint16(rec[74:76])),  // SMF4HOST
		int(binary.BigEndian.Uint32(rec[78:82])),  // SMF4RSH0
		int(rec[82]),                              // SMF4SPK
		int(rec[83]),                              // SMF4STI
		int(binary.BigEndian.Uint32(rec[86:90])),  // SMF4AST
		int(binary.BigEndian.Uint32(rec[90:94])),  // SMF4PPST
		b3toInt(rec[95:98]),                       // SMF4SRBT
		int(binary.BigEndian.Uint16(rec[98:100])), // SMF4RIN
		// Start of Relocation section
		int(binary.BigEndian.Uint32(rec[reloc+0:reloc+4])),   // SMF4PGIN
		int(binary.BigEndian.Uint32(rec[reloc+4:reloc+8])),   // SMF4PGOT
		int(binary.BigEndian.Uint32(rec[reloc+8:reloc+12])),  // SMF4NSW
		int(binary.BigEndian.Uint32(rec[reloc+12:reloc+16])), // SMF4PSI
		int(binary.BigEndian.Uint32(rec[reloc+16:reloc+20])), // SMF4PSO
		int(binary.BigEndian.Uint32(rec[reloc+20:reloc+24])), // SMF4VPI
		int(binary.BigEndian.Uint32(rec[reloc+24:reloc+28])), // SMF4VPO
		int(binary.BigEndian.Uint32(rec[reloc+28:reloc+32])), // SMF4SST
		int(binary.BigEndian.Uint32(rec[reloc+32:reloc+36])), // SMF4ACT
		int(binary.BigEndian.Uint16(rec[reloc+36:reloc+38])), // SMF4PGNO
		int(binary.BigEndian.Uint32(rec[reloc+38:reloc+42])), // SMF4TRAN
		int(binary.BigEndian.Uint32(rec[reloc+42:reloc+46])), // SMF4CPM
		int(binary.BigEndian.Uint32(rec[reloc+46:reloc+50])), // SMF4RCLM
		int(binary.BigEndian.Uint32(rec[reloc+50:reloc+54])), // SMF4CPGN
		int(binary.BigEndian.Uint32(rec[reloc+54:reloc+58])), // SMF4HSPI
		int(binary.BigEndian.Uint32(rec[reloc+58:reloc+62])), // SMF4PGST
		int(binary.BigEndian.Uint64(rec[reloc+62:reloc+70])), // SMF4PSEC
		int(binary.BigEndian.Uint32(rec[reloc+70:reloc+74])), // SMF4LPAI
		int(binary.BigEndian.Uint32(rec[reloc+74:reloc+78])), // SMF4HSPO
		int(binary.BigEndian.Uint32(rec[reloc+78:reloc+82])), // SMF4CPUS
		int(binary.BigEndian.Uint32(rec[reloc+82:reloc+86])), // SMF4IOCS
		int(binary.BigEndian.Uint32(rec[reloc+86:reloc+90])), // SMF4MSOS
		int(binary.BigEndian.Uint32(rec[reloc+90:reloc+94])), // SMF4SRBS
		b3toInt(rec[acctOff+1:acctOff+4]),                    // SMF4SETM
		acctFields.String(),                                  // SMF4ACTF
	)
	if err != nil {
		log.Println(err)
	}

	devOff := 104
	for i := 0; i < numDevs; i++ {
		_, err = p2Stmt004.Exec(
			count004,
			int(rec[devOff+0]), // SMF4DEVC
			int(rec[devOff+1]), // SMF4UTYP
			int(binary.BigEndian.Uint16(rec[devOff+2:devOff+4])), // SMF4CUAD
			int(binary.BigEndian.Uint32(rec[devOff+4:devOff+8])), // SMF4EXCP
		)
		if err != nil {
			log.Println(err)
		}
		devOff += 8
	}
}
