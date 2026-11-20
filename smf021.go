package main

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var pStmt021 *sql.Stmt
var once021 sync.Once
var count021 int

func smf021Init() {
	creatStmt := `
	DROP TABLE IF EXISTS smf021;
	CREATE TABLE smf021 (
	    ID         INT        NOT NULL,
		SMF21TME   DATETIME   NOT NULL,
		SMF21DTE   DATETIME   NOT NULL,
		SMF21SID   TEXT       NOT NULL,
		SMF21VOL   TEXT       NOT NULL,
		SMF21CA    INT        NOT NULL,
		SMF21UCB   INT        NOT NULL,
		SMF21DEV   INT        NOT NULL,
		SMF21TR    INT        NOT NULL,
		SMF21TW    INT        NOT NULL,
		SMF21SIO   INT        NOT NULL,
		SMF21PR    INT        NOT NULL,
		SMF21PW    INT        NOT NULL,
		SMF21NB    INT        NOT NULL,
		SMF21ERG   INT        NOT NULL,
		SMF21CLN   INT        NOT NULL,
		SMF21BLS   INT        NOT NULL,
		SMF21OFL   INT        NOT NULL,
		SMF21TUS   TEXT       NOT NULL,
		SMF21TRF   INT        NOT NULL,
		SMF21TRB   INT        NOT NULL,
		SMF21TWF   INT        NOT NULL,
		SMF21BR    INT        ,
		SMF21BW    INT        ,
		SMF21FL1   INT        ,
		SMF21BRN   INT        ,
		SMF21BWN   INT        ,
		SMF21LST   INT        ,
		SMF21LBS   INT        ,
		SMF21DBR   INT        ,
		SMF21DBW   INT        ,
		SMF21MCR   INT        ,
		SMF21MCW   INT        ,
		SMF21MDR   INT        ,
		SMF21MDW   INT
	);
	`
	_, err := db.Exec(creatStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, creatStmt)
	}

	insertStmt := `INSERT INTO smf021 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	pStmt021, err = db.Prepare(insertStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, insertStmt)
	}
}

func smf021Parse(rec []byte) {
	once021.Do(smf021Init)

	u32 := func(off int) int { return int(binary.BigEndian.Uint32(rec[off : off+4])) }

	// Fields from SMF21BR on are absent in short records. Each later field is valid only when its
	// SMF21FL1 flag is on; otherwise it is stored as NULL.
	var br, bw, fl1, brn, bwn, lst, lbs, dbr, dbw, mcr, mcw, mdr, mdw interface{}
	if len(rec) >= 60 {
		br, bw, fl1 = b3toInt(rec[52:55]), b3toInt(rec[55:58]), int(rec[58])
	}
	if len(rec) >= 100 {
		flags := rec[58]
		if flags&0x80 != 0 { // SMF21NCT
			brn, bwn = u32(60), u32(64)
		}
		if flags&0x40 != 0 { // SMF21LS
			lst = u32(68)
		}
		if flags&0x20 != 0 { // SMF21LB
			lbs = u32(72)
		}
		if flags&0x10 != 0 { // SMF21DBV
			dbr, dbw = u32(76), u32(80)
		}
		if flags&0x08 != 0 { // SMF21MFV
			mcr, mcw, mdr, mdw = u32(84), u32(88), u32(92), u32(96)
		}
	}

	count021 += 1
	_, err = pStmt021.Exec(
		count021,
		smftime(rec[2:6]),                        // SMF21TME
		smfdate(rec[6:10]),                       // SMF21DTE
		ebcdic(rec[10:14]),                       // SMF21SID
		ebcdic(rec[16:22]),                       // SMF21VOL
		int(binary.BigEndian.Uint16(rec[22:24])), // SMF21CA
		u32(24),                                  // SMF21UCB
		int(rec[27]),                             // SMF21DEV
		int(rec[28]),                             // SMF21TR
		int(rec[29]),                             // SMF21TW
		int(binary.BigEndian.Uint16(rec[30:32])), // SMF21SIO
		int(rec[32]),                             // SMF21PR
		int(rec[33]),                             // SMF21PW
		int(rec[34]),                             // SMF21NB
		int(binary.BigEndian.Uint16(rec[35:37])), // SMF21ERG
		int(binary.BigEndian.Uint16(rec[37:39])), // SMF21CLN
		int(binary.BigEndian.Uint16(rec[40:42])), // SMF21BLS
		int(rec[42]),                             // SMF21OFL
		fmt.Sprintf("%X", rec[43:46]),            // SMF21TUS, packed digits without sign
		int(binary.BigEndian.Uint16(rec[46:48])), // SMF21TRF
		int(binary.BigEndian.Uint16(rec[48:50])), // SMF21TRB
		int(binary.BigEndian.Uint16(rec[50:52])), // SMF21TWF
		br,                                       // SMF21BR
		bw,                                       // SMF21BW
		fl1,                                      // SMF21FL1
		brn,                                      // SMF21BRN
		bwn,                                      // SMF21BWN
		lst,                                      // SMF21LST
		lbs,                                      // SMF21LBS
		dbr,                                      // SMF21DBR
		dbw,                                      // SMF21DBW
		mcr,                                      // SMF21MCR
		mcw,                                      // SMF21MCW
		mdr,                                      // SMF21MDR
		mdw,                                      // SMF21MDW
	)
	if err != nil {
		log.Println(err)
	}
}
