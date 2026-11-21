package main

import (
	"database/sql"
	"sync"
)

// Types 2 (dump header) and 3 (dump trailer), written by the SMF dump program around the records it
// dumps. They share one layout: the standard header only, when SMF2FLG/SMF3FLG bit 1 is off. Not decoded,
// as the sample data has none: type 2 subtypes 1 and 2 (signature records, with bit 1 on).

var pStmt002, pStmt003 *sql.Stmt
var once002, once003 sync.Once
var count002, count003 int

func smfDumpFields(p string) fields {
	return fields{
		{p + "FLG", 0, 'u', 1},
		{p + "TME", 2, 'T', 4},
		{p + "DTE", 6, 'd', 4},
		{p + "SID", 10, 'c', 4},
	}
}

var smf002Fields, smf003Fields = smfDumpFields("SMF2"), smfDumpFields("SMF3")

func smf002Parse(rec []byte) {
	once002.Do(func() { pStmt002 = createTable("smf002", "ID INT NOT NULL, "+smf002Fields.columns()) })
	count002 += 1
	insert(pStmt002, append([]interface{}{count002}, smf002Fields.values(section(rec))...)...)
}

func smf003Parse(rec []byte) {
	once003.Do(func() { pStmt003 = createTable("smf003", "ID INT NOT NULL, "+smf003Fields.columns()) })
	count003 += 1
	insert(pStmt003, append([]interface{}{count003}, smf003Fields.values(section(rec))...)...)
}
