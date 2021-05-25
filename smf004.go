package main

import (
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type smf004 struct {
	SMF4TME  [4]byte
	SMF4DTE  [4]byte
	SMF4SID  [4]byte
	SMF4JBN  [8]byte
	SMF4RST  [4]byte
	SMF4RSD  [4]byte
	SMF4UIF  [8]byte
	SMF4STN  byte
	SMF4SIT  [4]byte
	SMF4STID [4]byte
	SMF4NCI  uint32
	SMF4SCC  uint16
	SMF4PRTY byte
	SMF4PGMN [8]byte
	SMF4STMN [8]byte
	SMF4RSV5 [2]byte
	SMF4SYST uint16
	SMF4HOST uint16
	SMF4RV1  [2]byte
	SMF4RSH0 uint32
	SMF4SPK  byte
}

func smf004Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf004;
	CREATE TABLE smf004 (
	    SMF4TME  DATETIME	not null,
	    SMF4DTE  DATETIME	not null,
	    SMF4SID  TEXT 		not null,
	    SMF4JBN  TEXT 		not null,
	    SMF4RST  DATETIME 	not null,
	    SMF4RSD  DATETIME 	not null,
	    SMF4UIF  TEXT 		not null,
	    SMF4STN  INT 		not null,
	    SMF4SIT  DATETIME 	not null,
	    SMF4STID DATETIME	not null,
	    SMF4NCI  INT 		not null,
	    SMF4SCC  INT 		not null,
	    SMF4PRTY INT 		not null,
	    SMF4PGMN TEXT 		not null,
	    SMF4STMN TEXT 		not null,
	    SMF4SYST INT 		not null,
	    SMF4HOST INT 		not null,
	    SMF4RSH0 INT 		not null,
	    SMF4SPK  INT 		not null
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}

func smf004Parse(rec []byte) {
	//var recbuf smf004
	//rdr := bytes.NewReader(rec)
	//err = binary.Read(rdr, binary.BigEndian, &recbuf)
	sqlStmt := `
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
			?
		);
`
	prepStmt, err := db.Prepare(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	defer prepStmt.Close()
	bcd2asc(rec[10:14], 4)
	bcd2asc(rec[14:22], 8)
	bcd2asc(rec[30:38], 8)
	bcd2asc(rec[54:62], 8)
	bcd2asc(rec[62:70], 8)
	_, err = prepStmt.Exec(
		smftime(rec[2:6]),
		smfdate(rec[6:10]),
		string(rec[10:14]),
		string(rec[14:22]),
		smftime(rec[22:26]),
		smfdate(rec[26:30]),
		string(rec[30:38]),
		int(rec[38]),
		smftime(rec[39:43]),
		smfdate(rec[43:47]),
		binary.BigEndian.Uint32(rec[47:51]),
		binary.BigEndian.Uint16(rec[51:53]),
		int(rec[53]),
		string(rec[54:62]),
		string(rec[62:70]),
		binary.BigEndian.Uint16(rec[72:74]),
		binary.BigEndian.Uint16(rec[74:76]),
		binary.BigEndian.Uint32(rec[78:82]),
		int(rec[82]),
	)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
}
