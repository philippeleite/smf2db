package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

var p1Stmt018 *sql.Stmt
var p2Stmt018 *sql.Stmt
var once018 sync.Once
var count018 int

func smf018Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf018;
	CREATE TABLE smf018 (
	    ID         INT        NOT NULL,
		SMF18TME   DATETIME   NOT NULL,
		SMF18DTE   DATETIME   NOT NULL,
		SMF18SID   TEXT       NOT NULL,
		SMF18JBN   TEXT       NOT NULL,
		SMF18RST   DATETIME   NOT NULL,
		SMF18RSD   DATETIME   NOT NULL,
		SMF18UID   TEXT       NOT NULL,
		SMF18IN1   INT        NOT NULL,
		SMF18CON   INT        NOT NULL,
		SMF18IN2   INT        NOT NULL,
		SMF180DS   TEXT       NOT NULL,
		SMF18NDS   TEXT       NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
	DROP TABLE IF EXISTS smf018_vol;
	CREATE TABLE smf018_vol (
	    ID         INT        NOT NULL,
		SMF18FVL   TEXT       NOT NULL
	);
	`
	_, err = db.Exec(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
	}

	sqlStmt3 := `
		INSERT INTO smf018
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
			?
		);
`
	p1Stmt018, err = db.Prepare(sqlStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt3)
	}

	sqlStmt4 := `
		INSERT INTO smf018_vol
		VALUES (
			?,
			?
		);
`

	p2Stmt018, err = db.Prepare(sqlStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt4)
	}
}

func smf018Parse(rec []byte) {
	once018.Do(smf018Init)

	bcd2asc(rec[10:14])  // SMF18SID
	bcd2asc(rec[14:22])  // SMF18JBN
	bcd2asc(rec[30:38])  // SMF18UID
	bcd2asc(rec[40:84])  // SMF180DS
	bcd2asc(rec[84:128]) // SMF18NDS

	numVols := int(rec[131])

	count018 += 1
	_, err = p1Stmt018.Exec(
		count018,
		smftime(rec[2:6]),   // SMF18TME
		smfdate(rec[6:10]),  // SMF18DTE
		string(rec[10:14]),  // SMF18SID
		string(rec[14:22]),  // SMF18JBN
		smftime(rec[22:26]), // SMF18RST
		smfdate(rec[26:30]), // SMF18RSD
		string(rec[30:38]),  // SMF18UID
		int(rec[38]),        // SMF18IN1
		int(rec[38] >> 7),   // SMF18CON: bit 0 of SMF18IN1, continuation record of a multi-volume rename
		int(rec[39]),        // SMF18IN2
		string(rec[40:84]),  // SMF180DS
		string(rec[84:128]), // SMF18NDS
	)
	if err != nil {
		log.Println(err)
	}

	devOff := 132
	for i := 0; i < numVols; i++ {
		bcd2asc(rec[devOff+2 : devOff+8])
		_, err = p2Stmt018.Exec(
			count018,
			string(rec[devOff+2:devOff+8]), // SMF18FVL
		)
		if err != nil {
			log.Println(err)
		}
		devOff += 8
	}
}
