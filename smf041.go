package main

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

var p1Stmt041 *sql.Stmt
var p2Stmt041 *sql.Stmt
var p3Stmt041 *sql.Stmt
var p4Stmt041 *sql.Stmt
var p5Stmt041 *sql.Stmt
var p6Stmt041 *sql.Stmt
var once041 sync.Once
var count041 int

type header041 struct {
	SMF41TME string
	SMF41DTE string
	SMF41SID string
	SMF41SSI string
	SMF41STY int
	SMF41TRP int
	SMF41OPD int
	SMF41LPD int
	SMF41NPD int
	SMF41OD1 int
	SMF41LD1 int
	SMF41ND1 int
	SMF41OD2 int
	SMF41LD2 int
	SMF41ND2 int
	SMF41OD3 int
	SMF41LD3 int
	SMF41ND3 int
	SMF41OD4 int
	SMF41LD4 int
	SMF41ND4 int
}

func smf041Init() {
	creatStmt1 := `
	DROP TABLE IF EXISTS smf041;
	CREATE TABLE smf041 (
	    ID         INT        NOT NULL,
		SMF41TME   DATETIME   NOT NULL,
		SMF41DTE   DATETIME   NOT NULL,
		SMF41SID   TEXT       NOT NULL,
		SMF41SSI   TEXT       NOT NULL,
		SMF41STY   INT        NOT NULL
	);
	`
	_, err := db.Exec(creatStmt1)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt1)
	}

	creatStmt2 := `
	DROP TABLE IF EXISTS smf041_prod;
	CREATE TABLE smf041_prod (
	    ID         INT        NOT NULL,
		SMF41PL    TEXT       NOT NULL,
		SMF41PN    TEXT       NOT NULL
	);
	`
	_, err = db.Exec(creatStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt2)
	}

	creatStmt3 := `
	DROP TABLE IF EXISTS smf041_obja;
	CREATE TABLE smf041_obja (
	    ID         INT        NOT NULL,
		SMF41DDA   TEXT       NOT NULL,
		SMF41AZA   INT        NOT NULL,
		SMF41ATA   INT        NOT NULL,
		SMF41TYA   INT        NOT NULL,
		SMF41AMA   INT        NOT NULL,
		SMF41JBN   TEXT       NOT NULL
	);
	`
	_, err = db.Exec(creatStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt3)
	}

	creatStmt4 := `
	DROP TABLE IF EXISTS smf041_obju;
	CREATE TABLE smf041_obju (
	    ID         INT        NOT NULL,
		SMF41UZU   INT        NOT NULL,
		SMF41UTU   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt4)
	}

	creatStmt5 := `
	DROP TABLE IF EXISTS smf041_obji;
	CREATE TABLE smf041_obji (
	    ID         INT        NOT NULL,
		SMF41BRD   INT        NOT NULL,
		SMF41BWR   INT        NOT NULL,
		SMF41BRR   INT        NOT NULL,
		SMF41NC    INT        NOT NULL,
		SMF41OUC   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt5)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt5)
	}

	creatStmt6 := `
	DROP TABLE IF EXISTS smf041_vlfs;
	CREATE TABLE smf041_vlfs (
	    ID         INT        NOT NULL,
		SMF41CLS   TEXT       NOT NULL,
		SMF41MVT   INT        NOT NULL,
		SMF41USD   INT        NOT NULL,
		SMF41SRC   INT        NOT NULL,
		SMF41FND   INT        NOT NULL,
		SMF41ADD   INT        NOT NULL,
		SMF41DEL   INT        NOT NULL,
		SMF41TRM   INT        NOT NULL,
		SMF41LRG   INT        NOT NULL
	);
	`
	_, err = db.Exec(creatStmt6)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt6)
	}

	insertStmt1 := `
		INSERT INTO smf041
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?
		);
`
	p1Stmt041, err = db.Prepare(insertStmt1)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt1)
	}

	insertStmt2 := `
		INSERT INTO smf041_prod
		VALUES (
			?,
			?,
			?
		);
`

	p2Stmt041, err = db.Prepare(insertStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt2)
	}

	insertStmt3 := `
		INSERT INTO smf041_obja
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

	p3Stmt041, err = db.Prepare(insertStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt3)
	}

	insertStmt4 := `
		INSERT INTO smf041_obju
		VALUES (
			?,
			?,
			?
		);
`

	p4Stmt041, err = db.Prepare(insertStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt4)
	}

	insertStmt5 := `
		INSERT INTO smf041_obji
		VALUES (
			?,
			?,
			?,
			?,
			?,
			?
		);
`

	p5Stmt041, err = db.Prepare(insertStmt5)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt5)
	}

	insertStmt6 := `
		INSERT INTO smf041_vlfs
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

	p6Stmt041, err = db.Prepare(insertStmt6)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt6)
	}
}

func smf041Parse(rec []byte) {
	once041.Do(smf041Init)

	bcd2asc(rec[10:14]) // SMF41SID
	bcd2asc(rec[14:18]) // SMF41SSI

	var header = header041{
		SMF41TME: smftime(rec[2:6]),                        // SMF41TME
		SMF41DTE: smfdate(rec[6:10]),                       // SMF41DTE
		SMF41SID: string(rec[10:14]),                       // SMF41SID
		SMF41SSI: string(rec[14:18]),                       // SMF41SSI
		SMF41STY: int(binary.BigEndian.Uint16(rec[18:20])), // SMF41STY
		SMF41TRP: int(binary.BigEndian.Uint16(rec[20:22])), // SMF41TRP
		SMF41OPD: int(binary.BigEndian.Uint32(rec[24:28])), // SMF41OPD
		SMF41LPD: int(binary.BigEndian.Uint16(rec[28:30])), // SMF41LPD
		SMF41NPD: int(binary.BigEndian.Uint16(rec[30:32])), // SMF41NPD
		SMF41OD1: int(binary.BigEndian.Uint32(rec[32:36])), // SMF41OD1
		SMF41LD1: int(binary.BigEndian.Uint16(rec[36:38])), // SMF41LD1
		SMF41ND1: int(binary.BigEndian.Uint16(rec[38:40])), // SMF41ND1
		SMF41OD2: int(binary.BigEndian.Uint32(rec[40:44])), // SMF41OD2
		SMF41LD2: int(binary.BigEndian.Uint16(rec[44:46])), // SMF41LD2
		SMF41ND2: int(binary.BigEndian.Uint16(rec[46:48])), // SMF41ND2
		SMF41OD3: int(binary.BigEndian.Uint32(rec[48:52])), // SMF41OD3
		SMF41LD3: int(binary.BigEndian.Uint16(rec[52:54])), // SMF41LD3
		SMF41ND3: int(binary.BigEndian.Uint16(rec[54:56])), // SMF41ND3
		SMF41OD4: int(binary.BigEndian.Uint32(rec[56:60])), // SMF41OD4
		SMF41LD4: int(binary.BigEndian.Uint16(rec[60:62])), // SMF41LD4
		SMF41ND4: int(binary.BigEndian.Uint16(rec[62:64])), // SMF41ND4
	}

	count041 += 1
	_, err = p1Stmt041.Exec(
		count041,
		header.SMF41DTE,
		header.SMF41TME,
		header.SMF41SID,
		header.SMF41SSI,
		header.SMF41STY,
	)
	if err != nil {
		log.Println(err)
	}

	// Product section
	prodOff := header.SMF41OPD - 4
	for i := 0; i < header.SMF41NPD; i++ {
		bcd2asc(rec[prodOff+0 : prodOff+8])
		bcd2asc(rec[prodOff+8 : prodOff+24])
		_, err = p2Stmt041.Exec(
			count041,
			string(rec[prodOff+0:prodOff+8]),  // SMF41PL
			string(rec[prodOff+8:prodOff+24]), // SMF41PN
		)
		if err != nil {
			log.Println(err)
		}
		prodOff += header.SMF41LPD
	}

	// Object ACCESS data section
	objaOff := header.SMF41OD1 - 4
	for i := 0; i < header.SMF41ND1; i++ {
		bcd2asc(rec[objaOff+0 : objaOff+8])   // SMF41DDA
		bcd2asc(rec[objaOff+18 : objaOff+26]) // SMF41JBN
		_, err = p3Stmt041.Exec(
			count041,
			string(rec[objaOff+0:objaOff+8]), // SMF41DDA
			int(binary.BigEndian.Uint32(rec[objaOff+8:objaOff+12])),  // SMF41AZA
			int(binary.BigEndian.Uint32(rec[objaOff+12:objaOff+16])), // SMF41ATA
			int(rec[objaOff+16]),               // SMF41TYA
			int(rec[objaOff+17]),               // SMF41AMA
			string(rec[objaOff+18:objaOff+26]), // SMF41JBN
		)
		if err != nil {
			log.Println(err)
		}
		objaOff += header.SMF41LD1
	}

	// Object UNACCESS data section
	objuOff := header.SMF41OD2 - 4
	for i := 0; i < header.SMF41ND2; i++ {
		_, err = p4Stmt041.Exec(
			count041,
			int(binary.BigEndian.Uint32(rec[objaOff+0:objaOff+4])), // SMF41UZU
			int(binary.BigEndian.Uint32(rec[objaOff+4:objaOff+8])), // SMF41UTU
		)
		if err != nil {
			log.Println(err)
		}
		objuOff += header.SMF41LD2
	}

	// Object I/O activity section
	objiOff := header.SMF41OD3 - 4
	for i := 0; i < header.SMF41ND3; i++ {
		_, err = p5Stmt041.Exec(
			count041,
			int(binary.BigEndian.Uint32(rec[objiOff+0:objiOff+4])),   // SMF41BRD
			int(binary.BigEndian.Uint32(rec[objiOff+4:objiOff+8])),   // SMF41BWR
			int(binary.BigEndian.Uint32(rec[objiOff+8:objiOff+12])),  // SMF41BRR
			int(binary.BigEndian.Uint32(rec[objiOff+12:objiOff+16])), // SMF41NC
			int(binary.BigEndian.Uint32(rec[objiOff+16:objiOff+20])), // SMF41OUC
		)
		if err != nil {
			log.Println(err)
		}
		objiOff += header.SMF41LD3
	}

	// VLF statistics section
	vlfsOff := header.SMF41OD4 - 4
	for i := 0; i < header.SMF41ND4; i++ {
		bcd2asc(rec[vlfsOff+0 : vlfsOff+8]) // SMF41CLS
		_, err = p6Stmt041.Exec(
			count041,
			string(rec[vlfsOff+0:vlfsOff+8]), // SMF41CLS
			int(binary.BigEndian.Uint32(rec[vlfsOff+8:vlfsOff+12])),  // SMF41MVT
			int(binary.BigEndian.Uint32(rec[vlfsOff+12:vlfsOff+16])), // SMF41USD
			int(binary.BigEndian.Uint32(rec[vlfsOff+16:vlfsOff+20])), // SMF41SRC
			int(binary.BigEndian.Uint32(rec[vlfsOff+20:vlfsOff+24])), // SMF41FND
			int(binary.BigEndian.Uint32(rec[vlfsOff+24:vlfsOff+28])), // SMF41ADD
			int(binary.BigEndian.Uint32(rec[vlfsOff+28:vlfsOff+32])), // SMF41DEL
			int(binary.BigEndian.Uint32(rec[vlfsOff+32:vlfsOff+36])), // SMF41TRM
			int(binary.BigEndian.Uint32(rec[vlfsOff+36:vlfsOff+40])), // SMF41LRG
		)
		if err != nil {
			log.Println(err)
		}
		vlfsOff += header.SMF41LD4
	}
}
