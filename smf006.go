package main

import (
	"database/sql"
	"encoding/binary"
	"log"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var p1Stmt006 *sql.Stmt
var p2Stmt006 *sql.Stmt
var once006 sync.Once
var count006 int

const smf006SbsIPPrintWay = 9 // SMF6SBS X'0009'

func smf006Init() {
	creatStmt1 := `
	DROP TABLE IF EXISTS smf006;
	CREATE TABLE smf006 (
	    ID         INT        NOT NULL,
		SMF6TME    DATETIME   NOT NULL,
		SMF6DTE    DATETIME   NOT NULL,
		SMF6SID    TEXT       NOT NULL,
		SMF6JBN    TEXT       NOT NULL,
		SMF6RST    DATETIME   NOT NULL,
		SMF6RSD    DATETIME   NOT NULL,
		SMF6UIF    TEXT       NOT NULL,
		SMF6OWC    TEXT       NOT NULL,
		SMF6WST    DATETIME   NOT NULL,
		SMF6WSD    DATETIME   NOT NULL,
		SMF6NLR    INT        NOT NULL,
		SMF6IOE    INT        NOT NULL,
		SMF6NDS    INT        NOT NULL,
		SMF6FMN    TEXT       NOT NULL,
		SMF6PAD1   INT        NOT NULL,
		SMF6SBS    INT        NOT NULL
	);
	`
	_, err := db.Exec(creatStmt1)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt1)
	}

	// IP PrintWay sections, one row per record with SMF6SBS = X'0009'.
	// Columns of optional sections are NULL when SMF6PAD1 says the section is absent.
	creatStmt2 := `
	DROP TABLE IF EXISTS smf006_ipw;
	CREATE TABLE smf006_ipw (
	    ID         INT        NOT NULL,
		SMF6DCI    INT        NOT NULL,
		SMF6INDC   INT        NOT NULL,
		SMF6OUT    TEXT       NOT NULL,
		SMF6PGE    INT        NOT NULL,
		SMF6EFMN   TEXT       ,
		SMF6JBID   TEXT       ,
		SMF6STNM   TEXT       ,
		SMF6PRNM   TEXT       ,
		SMF6DDNM   TEXT       ,
		SMF6USID   TEXT       ,
		SMF6SECS   TEXT       ,
		SMF6PRMD   TEXT       ,
		SMF6DSNM   TEXT       ,
		SMF6OTOK   BLOB       ,
		SMF6SGID   INT        ,
		SMF6IND    INT        ,
		SMF6JDVT   TEXT       ,
		SMF6TU     BLOB       ,
		SMF6BYTE   INT        ,
		SMF6IP1    INT        ,
		SMF6IP2    INT        ,
		SMF6IP3    INT        ,
		SMF6IP4    INT        ,
		SMF6FTL    INT        ,
		SMF6PRTQ   TEXT       ,
		SMF6BYTD   INT        ,
		SMF6URI    TEXT       ,
		SMF6ACCT   TEXT
	);
	`
	_, err = db.Exec(creatStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt2)
	}

	insertStmt1 := `INSERT INTO smf006 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	p1Stmt006, err = db.Prepare(insertStmt1)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt1)
	}

	insertStmt2 := `INSERT INTO smf006_ipw VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	p2Stmt006, err = db.Prepare(insertStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt2)
	}
}

// smf006Parse stores the header shared by every type 6 writer (External Writer, JES2, JES3, PSF,
// IP PrintWay). Only the IP PrintWay sections are decoded; the other writers' sections are skipped.
func smf006Parse(rec []byte) {
	once006.Do(smf006Init)

	count006 += 1
	_, err = p1Stmt006.Exec(
		count006,
		smftime(rec[2:6]),                        // SMF6TME
		smfdate(rec[6:10]),                       // SMF6DTE
		ebcdic(rec[10:14]),                       // SMF6SID
		ebcdic(rec[14:22]),                       // SMF6JBN
		smftime(rec[22:26]),                      // SMF6RST
		smfdate(rec[26:30]),                      // SMF6RSD
		ebcdic(rec[30:38]),                       // SMF6UIF
		ebcdic(rec[38:39]),                       // SMF6OWC
		smftime(rec[39:43]),                      // SMF6WST
		smfdate(rec[43:47]),                      // SMF6WSD
		int(binary.BigEndian.Uint32(rec[47:51])), // SMF6NLR
		int(rec[51]),                             // SMF6IOE
		int(rec[52]),                             // SMF6NDS
		ebcdic(rec[53:57]),                       // SMF6FMN
		int(rec[57]),                             // SMF6PAD1
		int(binary.BigEndian.Uint16(rec[58:60])), // SMF6SBS
	)
	if err != nil {
		log.Println(err)
	}

	if binary.BigEndian.Uint16(rec[58:60]) == smf006SbsIPPrintWay {
		smf006IPPrintWay(rec)
	}
}

// smf006IPPrintWay stores the IP PrintWay sections. The I/O data section starts right after the header;
// the optional sections follow it in SMF6PAD1 bit order, each one starting with its own 2-byte length.
func smf006IPPrintWay(rec []byte) {
	const ioOff = 60
	pad1 := rec[57]
	off := ioOff + int(binary.BigEndian.Uint16(rec[ioOff:ioOff+2])) // SMF6LN1

	// Common section
	var efmn, jbid, stnm, prnm, ddnm, usid, secs, prmd, dsnm, otok interface{}
	if pad1&0x40 != 0 {
		sec := rec[off : off+int(binary.BigEndian.Uint16(rec[off:off+2]))] // SMF6LN3
		efmn = ebcdic(sec[6:14])
		jbid = ebcdic(sec[30:38])
		if len(sec) >= 162 {
			stnm = ebcdic(sec[38:46])
			prnm = ebcdic(sec[46:54])
			ddnm = ebcdic(sec[54:62])
			usid = ebcdic(sec[62:70])
			secs = ebcdic(sec[70:78])
			prmd = ebcdic(sec[78:86])
			dsnm = ebcdic(sec[86:139])
			otok = sec[142:162] // documented as EBCDIC, but real records carry binary bytes after the job ID
		}
		off += len(sec)
	}

	// Enhanced SYSOUT support (ESS) section
	var sgid, ind, jdvt, tu interface{}
	if pad1&0x10 != 0 {
		sec := rec[off : off+int(binary.BigEndian.Uint16(rec[off:off+2]))] // SMF6LN5
		sgid = int(binary.BigEndian.Uint32(sec[2:6]))
		ind = int(sec[6])
		jdvt = ebcdic(sec[8:16])
		tu = sec[18 : 18+int(binary.BigEndian.Uint16(sec[16:18]))] // SMF6TUL, raw SWBTU mapped by IEFSJPFX
		off += len(sec)
	}

	// File transfer section
	var byteCnt, ip1, ip2, ip3, ip4, ftl, prtq, bytd, uri, acct interface{}
	if pad1&0x08 != 0 {
		sec := rec[off : off+int(binary.BigEndian.Uint16(rec[off:off+2]))] // SMF6LN6
		byteCnt = int(binary.BigEndian.Uint32(sec[2:6]))
		ip1, ip2, ip3, ip4 = int(sec[6]), int(sec[7]), int(sec[8]), int(sec[9])
		ftl = int(sec[10])
		prtq = ebcdic(sec[24 : 24+int(binary.BigEndian.Uint16(sec[22:24]))]) // SMF6PQLN
		if sec[10] >= 1 {                                                    // extended mode
			bytd = int(binary.BigEndian.Uint64(sec[48:56]))
			uri = ebcdic(sec[74 : 74+int(binary.BigEndian.Uint16(sec[72:74]))]) // SMF6URIL
		}
		if sec[10] == 3 { // extended mode with accounting fields
			if actl := int(binary.BigEndian.Uint16(sec[332:334])); actl > 0 {
				acct = smf006Acct(sec[334 : 334+actl])
			}
		}
	}

	_, err = p2Stmt006.Exec(
		count006,
		int(rec[ioOff+2]),             // SMF6DCI
		int(rec[ioOff+3]),             // SMF6INDC
		ebcdic(rec[ioOff+8:ioOff+16]), // SMF6OUT
		int(binary.BigEndian.Uint32(rec[ioOff+24:ioOff+28])), // SMF6PGE
		efmn,    // SMF6EFMN
		jbid,    // SMF6JBID
		stnm,    // SMF6STNM
		prnm,    // SMF6PRNM
		ddnm,    // SMF6DDNM
		usid,    // SMF6USID
		secs,    // SMF6SECS
		prmd,    // SMF6PRMD
		dsnm,    // SMF6DSNM
		otok,    // SMF6OTOK
		sgid,    // SMF6SGID
		ind,     // SMF6IND
		jdvt,    // SMF6JDVT
		tu,      // SMF6TU
		byteCnt, // SMF6BYTE
		ip1,     // SMF6IP1
		ip2,     // SMF6IP2
		ip3,     // SMF6IP3
		ip4,     // SMF6IP4
		ftl,     // SMF6FTL
		prtq,    // SMF6PRTQ
		bytd,    // SMF6BYTD
		uri,     // SMF6URI
		acct,    // SMF6ACCT
	)
	if err != nil {
		log.Println(err)
	}
}

// smf006Acct decodes SMF6ACCT: a count byte followed by (length byte, EBCDIC string) pairs.
// Same space-separated form as SMF5JSAF.
func smf006Acct(acct []byte) string {
	var b strings.Builder
	b.WriteString(" ")
	off := 1
	for i := 0; i < int(acct[0]); i++ {
		n := int(acct[off])
		b.WriteString(ebcdic(acct[off+1:off+1+n]) + " ")
		off += n + 1
	}
	return b.String()
}
