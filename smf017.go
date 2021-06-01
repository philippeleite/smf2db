package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var p1Stmt017 *sql.Stmt
var p2Stmt017 *sql.Stmt

func smf017Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf017;
	CREATE TABLE smf017 (
		SMF17TME   DATETIME   NOT NULL,
		SMF17DTE   DATETIME   NOT NULL,
		SMF17SID   TEXT       NOT NULL,
		SMF17JBN   TEXT       NOT NULL,
		SMF17RST   DATETIME   NOT NULL,
		SMF17RSD   DATETIME   NOT NULL,
		SMF17UID   TEXT       NOT NULL,
		SMF17DSN   TEXT       NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
	DROP TABLE IF EXISTS smf017_vol;
	CREATE TABLE smf017_vol (
		SMF17TME   DATETIME   NOT NULL,
		SMF17DTE   DATETIME   NOT NULL,
		SMF17FVL   TEXT       NOT NULL
	);
	`
	_, err = db.Exec(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
	}

	sqlStmt3 := `
		INSERT INTO smf017
		VALUES (
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
	p1Stmt017, err = db.Prepare(sqlStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt3)
	}

	sqlStmt4 := `
		INSERT INTO smf017_vol
		VALUES (
			?,
			?,
			?
		);
`

	p2Stmt017, err = db.Prepare(sqlStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt4)
	}
}

func smf017Parse(rec []byte) {

	bcd2asc(rec[10:14]) // SMF17SID
	bcd2asc(rec[14:22]) // SMF17JBN
	bcd2asc(rec[30:38]) // SMF17UID
	bcd2asc(rec[40:84]) // SMF17DSN

	numVols := int(rec[87])

	_, err = p1Stmt017.Exec(
		smftime(rec[2:6]),   // SMF17TME
		smfdate(rec[6:10]),  // SMF17DTE
		string(rec[10:14]),  // SMF17SID
		string(rec[14:22]),  // SMF17JBN
		smftime(rec[22:26]), // SMF17RST
		smfdate(rec[26:30]), // SMF17RSD
		string(rec[30:38]),  // SMF17UID
		string(rec[40:84]),  // SMF17DSN
	)
	if err != nil {
		log.Println(err)
	}

	devOff := 88
	for i := 0; i < numVols; i++ {
		bcd2asc(rec[devOff+2 : devOff+8])
		_, err = p2Stmt017.Exec(
			smftime(rec[2:6]),              // SMF17TME
			smfdate(rec[6:10]),             // SMF17DTE
			string(rec[devOff+2:devOff+8]), // SMF17FVL
		)
		if err != nil {
			log.Println(err)
		}
		devOff += 8
	}
}
