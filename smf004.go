package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
)

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
	    SMF4SID  DATETIME 	not null,
	    SMF4STID DATETIME	not null,
	    SMF4NCI  INT 		not null,
	    SMF4SCC  INT 		not null,
	    SMF4PRTY INT 		not null,
	    SMF4PGMN TEXT 		not null,
	    SMF4STMN TEXT 		not null,
	    SMF4SYST INT 		not null,
	    SMF4HOST INT 		not null,
	    SMF4RSH0 INT 		not null,
	    SMF4SPK  INT 		not null,
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func smf004Parse(rec []byte) {

}
