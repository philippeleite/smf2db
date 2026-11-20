package main

import (
	"database/sql"
	"sync"
)

// Types 61 (catalog define), 65 (catalog delete) and 66 (catalog alter) share the IFASMFI6 layout:
// job and catalog entry identification, the raw catalog record (length in its first two bytes) and,
// in newer versions, a 16-byte STCKE timestamp after it. Type 65 adds the function indicator SMFnnFNC;
// type 66 adds it and the new entry name SMFnnNNM. (Type 60 also uses IFASMFI6, but ends differently.)

var pStmt061, pStmt065, pStmt066 *sql.Stmt
var once061, once065, once066 sync.Once
var count061, count065, count066 int

func smf061Parse(rec []byte) {
	once061.Do(func() { pStmt061 = createTable("smf061", smfCatalogColumns("SMF61", false, false)) })
	count061 += 1
	insert(pStmt061, smfCatalogRow(rec, count061, false, false)...)
}

func smf065Parse(rec []byte) {
	once065.Do(func() { pStmt065 = createTable("smf065", smfCatalogColumns("SMF65", true, false)) })
	count065 += 1
	insert(pStmt065, smfCatalogRow(rec, count065, true, false)...)
}

func smf066Parse(rec []byte) {
	once066.Do(func() { pStmt066 = createTable("smf066", smfCatalogColumns("SMF66", true, true)) })
	count066 += 1
	insert(pStmt066, smfCatalogRow(rec, count066, true, true)...)
}

// smfCatalogColumns returns the column definitions matching smfCatalogRow.
func smfCatalogColumns(p string, fnc, nnm bool) string {
	cols := "ID INT NOT NULL, " + p + "TME DATETIME NOT NULL, " + p + "DTE DATETIME NOT NULL, " +
		p + "CPU TEXT NOT NULL, " + p + "SUB TEXT NOT NULL, " + p + "VER TEXT, " + p + "PNM TEXT, " +
		p + "JNM TEXT, " + p + "RST DATETIME, " + p + "RDT DATETIME, " + p + "UID TEXT, "
	if fnc {
		cols += p + "FNC TEXT, "
	}
	cols += p + "CNM TEXT, " + p + "TYP TEXT, " + p + "ENM TEXT, "
	if nnm {
		cols += p + "NNM TEXT, "
	}
	return cols + p + "CRC BLOB, " + p + "ETS DATETIME"
}

// smfCatalogRow decodes a type 61, 65 or 66 record. fnc includes the function indicator (offset 74 in the
// manual), nnm the new entry name (offset 164, reserved in types 61 and 65).
func smfCatalogRow(rec []byte, id int, fnc, nnm bool) []interface{} {
	s := section(rec)
	crcLen := s.num(204, 2)
	row := []interface{}{id,
		s.time(2),     // TME
		s.date(6),     // DTE
		s.text(10, 4), // CPU
		s.text(18, 2), // SUB: IN, UP or DE
		s.text(36, 2), // VER
		s.text(38, 8), // PNM
		s.text(46, 8), // JNM
		s.time(54),    // RST
		s.date(58),    // RDT
		s.text(62, 8), // UID
	}
	if fnc {
		row = append(row, s.text(70, 1)) // FNC: type 65 S = scratched, U = catalog only; type 66 R = renamed
	}
	row = append(row,
		s.text(71, 44),  // CNM
		s.text(115, 1),  // TYP, same values as SMF60TYP
		s.text(116, 44), // ENM
	)
	if nnm {
		row = append(row, s.text(160, 44)) // NNM
	}
	var ets interface{} // ETS: STCKE timestamp after the catalog record; absent in older versions
	if b := s.sub(204+crcLen, 16); b != nil {
		ets = stcketime(b)
	}
	return append(row, s.blob(204, crcLen), ets) // CRC, ETS
}
