package main

import (
	"database/sql"
	"sync"
)

// Type 84 (JES monitoring facility). The header differs from RMF records: SMF84SBS (2 = JES2, 5 = JES3),
// SMF84SGN, SMF84FL1 and SMF84VER come before the subtype, and the time is binary (the manual says packed).
// smf084 holds the header, the JES product section and the general section of every record. Subtype 21
// (JES2 resource usage) adds its memory area (smf084_mem) and resource (smf084_res) entries, located by the
// two triplets of the JES2 resource usage section, relative to its start; rows have ID and SEQ.
// Not decoded, as the sample data has none: subtypes 1 to 10 (JES3).
// The product and general field tables were generated from the manual with tools/genfields.py.

var pStmt084, pStmt084Mem, pStmt084Res *sql.Stmt
var once084 sync.Once
var count084 int

func smf084Init() {
	pStmt084 = createTable("smf084", `ID INT NOT NULL, SMF84TME DATETIME NOT NULL, SMF84DTE DATETIME NOT NULL,
		SMF84SID TEXT NOT NULL, SMF84SBS INT, SMF84SGN INT, SMF84FL1 INT, SMF84VER INT, SMF84STY INT, `+
		smf084ProductFields.columns()+", "+smf084GeneralFields.columns())
	pStmt084Mem = createTable("smf084_mem", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf084MemFields.columns())
	pStmt084Res = createTable("smf084_res", "ID INT NOT NULL, SEQ INT NOT NULL, "+smf084ResFields.columns())
}

func smf084Parse(rec []byte) {
	once084.Do(smf084Init)

	s := section(rec)
	count084 += 1
	id := count084
	poff, plen, _ := tripletAt(rec, 24)
	goff, glen, _ := tripletAt(rec, 32)
	row := []interface{}{id, s.time(2), s.date(6), s.text(10, 4), s.u16(14), s.u16(16), s.u8(18), s.u8(19), s.u16(20)}
	row = append(row, smf084ProductFields.values(s.sub(poff, plen))...)
	insert(pStmt084, append(row, smf084GeneralFields.values(s.sub(goff, glen))...)...)

	if s.num(20, 2) != 21 {
		return
	}
	off, length, n := tripletAt(rec, 40) // JES2 resource usage section (SMF84JRU)
	for k := 0; k < n; k++ {
		jru := s.sub(off+k*length, length)
		for _, t := range []struct {
			at   int
			stmt *sql.Stmt
			f    fields
		}{{32, pStmt084Mem, smf084MemFields}, {40, pStmt084Res, smf084ResFields}} {
			eoff, elen, en := jru.num(t.at, 4), jru.num(t.at+4, 2), jru.num(t.at+6, 2)
			for j := 0; j < en; j++ {
				insert(t.stmt, append([]interface{}{id, j + 1}, t.f.values(jru.sub(eoff+j*elen, elen))...)...)
			}
		}
	}
}

var smf084ProductFields = fields{
	{"R84MFVER", 0, 'u', 2},
	{"R84PRDNM", 2, 'c', 8},
	{"R84INTST", 10, 'H', 4},
	{"R84SDATE", 14, 'd', 4},
	{"R84INTEN", 18, 'H', 4},
	{"R84EDATE", 22, 'd', 4},
	{"R84INTER", 26, 'p', 4}, // seconds
	{"R84MFCYC", 30, 'p', 4}, // 00sssttF: hundredths of a second
	{"R84SAMPL", 36, 'u', 4},
	{"R84MFCMD", 40, 'c', 80},
	{"R84MVSRL", 120, 'c', 8},
	{"R84JESRL", 128, 'c', 8},
	{"R84CPUM", 136, 'c', 4},
	{"R84RSTO", 140, 'u', 4},
	{"R84CPUNM", 144, 'c', 8},
	{"R84CPUID", 152, 'c', 4},
	{"R84MPNAM", 156, 'c', 8},
	{"R84JPRTY", 166, 'u', 2},
	{"R84JMFMN", 168, 'u', 4},
	{"R84JMFMX", 172, 'u', 4},
	{"R84JMFAV", 176, 'u', 4},
	{"R84MVSMN", 180, 'u', 4},
	{"R84MVSMX", 184, 'u', 4},
	{"R84MVSAV", 188, 'u', 4},
}
var smf084GeneralFields = fields{
	{"R84CPUSC", 0, 'u', 4},
	{"R84NPA", 4, 'u', 4},
	{"R84APA", 8, 'u', 4},
	{"R84NPNA", 12, 'u', 4},
	{"R84APNA", 16, 'u', 4},
	{"R84NNP", 20, 'u', 4},
	{"R84ANP", 24, 'u', 4},
	{"R84NNW", 28, 'u', 4},
	{"R84ANW", 32, 'u', 4},
	{"R84NSLLR", 36, 'u', 4},
	{"R84ASLLR", 40, 'u', 4},
	{"R84NSO", 44, 'u', 4},
	{"R84ASO", 48, 'u', 4},
}

// smf084MemFields: JES2 memory area usage entry (R84MEMJ2).
var smf084MemFields = fields{
	{"R84MEM_NAME", 0, 'c', 12},
	{"R84MEM_REGION", 16, 'u', 8},
	{"R84MEM_USE", 24, 'u', 8},
	{"R84MEM_LOW", 32, 'u', 8},
	{"R84MEM_HIGH", 40, 'u', 8},
	{"R84MEM_AVERAGE", 48, 'u', 8},
}

// smf084ResFields: JES2 resource usage entry (R84RSUJ2).
var smf084ResFields = fields{
	{"R84RSU_NAME", 0, 'c', 8},
	{"R84RSU_LIMIT", 8, 'u', 4},
	{"R84RSU_INUSE", 12, 'u', 4},
	{"R84RSU_LOW", 16, 'u', 4},
	{"R84RSU_HIGH", 20, 'u', 4},
	{"R84RSU_WARN", 24, 'u', 2},
	{"R84RSU_FLG1", 26, 'u', 1},
	{"R84RSU_OVER", 28, 'u', 4},
	{"R84RSU_AVERAGE", 32, 'u', 4},
}
