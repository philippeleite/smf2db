package main

import "database/sql"

// Type 42 subtypes 1, 2, 3, 4, 9, 21, 24 and 27. See smf042.go.

var (
	pStmt042Bmf, pStmt042BmfSc      *sql.Stmt // subtype 1
	pStmt042Cu, pStmt042CuVol       *sql.Stmt // subtype 2
	pStmt042Ea                      *sql.Stmt // subtype 3
	pStmt042Ext                     *sql.Stmt // subtype 4
	pStmt042X37                     *sql.Stmt // subtype 9
	pStmt042Mdel, pStmt042MdelAlias *sql.Stmt // subtype 21
	pStmt042Madd, pStmt042MaddAlias *sql.Stmt // subtype 24
	pStmt042Vtoc                    *sql.Stmt // subtype 27
)

func smf042InitMisc() {
	pStmt042Bmf = createTable("smf042_bmf", `ID INT NOT NULL, SMF42TNA INT, SMF42TMT INT, SMF42TRT INT,
		SMF42TRH INT, SMF42TDT INT, SMF42TDH INT, SMF42BUF INT, SMF42BMX INT, SMF42LRU INT, SMF42UIC INT`)
	pStmt042BmfSc = createTable("smf042_bmf_sc", `ID INT NOT NULL, SMF42PNN TEXT, SMF42SRT INT, SMF42SRH INT,
		SMF42SDT INT, SMF42SDH INT`)

	pStmt042Cu = createTable("smf042_cu", `ID INT NOT NULL, SMF42SCS INT, SMF42NCS INT, SMF42CID INT,
		SMF42CSS INT, SMF42SSA INT, SMF42SAP INT, SMF42SSU INT, SMF42NSZ INT, SMF42SPR INT,
		SMF42LCT INT, SMF42LFW INT, SMF42LRH INT, SMF42LWM INT, SMF42LYY INT, SMF42LDD INT, SMF42LTM INT,
		SMF42CCT INT, SMF42CFW INT, SMF42CRH INT, SMF42CWM INT, SMF42CYY INT, SMF42CDD INT, SMF42CTM INT,
		SMF42IHR INT, SMF42IFW INT`)
	pStmt042CuVol = createTable("smf042_cu_vol", `ID INT NOT NULL, SMF42VOL TEXT, SMF42DEV INT,
		SMF42DB1 INT, SMF42DB2 INT`)

	pStmt042Ea = createTable("smf042_ea", `ID INT NOT NULL, SMF42EAC TEXT, SMF42ERC INT, SMF42ERS INT,
		SMF42EUA INT, SMF42EVO TEXT, SMF42EOS INT, SMF42ENS INT, SMF42ETY TEXT, SMF42ENM TEXT,
		SMF42EVL TEXT, SMF42ESY TEXT, SMF42EST TEXT, SMF42ESD TEXT, SMF42EAD TEXT`)

	pStmt042Ext = createTable("smf042_ext", `ID INT NOT NULL, S42EXID INT, S42EXRQS TEXT, S42EXJNM TEXT,
		S42EXJNO TEXT, S42EXSTM DATETIME, S42EXETM DATETIME, S42EXTS TEXT, S42EVLRT INT, S42EVLWT INT`)

	pStmt042X37 = createTable("smf042_x37", `ID INT NOT NULL, S42ASYID TEXT, S42JOBN TEXT, S42RDST DATETIME,
		S42RDSD DATETIME, S42AAUID TEXT, S42ASTPN INT, S42FLAGS INT, S42DSORG INT, S42ADISP INT,
		S42DSNME TEXT, S42VOLSR TEXT, S42UCBTP INT, S42NEXT INT, S42TNTRK INT, S42ASSAT INT, S42ADRLH INT,
		S42MCNME TEXT, S42SCNME TEXT, S42DCNME TEXT`)

	pStmt042Mdel = createTable("smf042_mdel", `ID INT NOT NULL, SMF42LJB TEXT, SMF42LST TEXT, SMF42LPR TEXT,
		SMF42LDS TEXT, SMF42LVS TEXT, SMF42LFL INT, SMF42LMN TEXT, SMF42LUI BLOB`)
	pStmt042MdelAlias = createTable("smf042_mdel_alias", `ID INT NOT NULL, SMF42LAN TEXT`)

	pStmt042Madd = createTable("smf042_madd", `ID INT NOT NULL, SMF42PJB TEXT, SMF42PST TEXT, SMF42PPR TEXT,
		SMF42PDS TEXT, SMF42PVS TEXT, SMF42PF1 INT, SMF42PMN TEXT, SMF42PUI BLOB`)
	pStmt042MaddAlias = createTable("smf042_madd_alias", `ID INT NOT NULL, SMF42PAN TEXT`)

	pStmt042Vtoc = createTable("smf042_vtoc", `ID INT NOT NULL, SMF42RJOB TEXT, SMF42RJNO TEXT,
		SMF42RSTN TEXT, SMF42RPRN TEXT, SMF42RVOL TEXT, SMF42RDEV INT, SMF42RACT TEXT, SMF42RIND INT,
		SMF42RDS1 INT, SMF42RSEEK INT, SMF42RSRCH BLOB, SMF42RCMDS BLOB, SMF42RUPSW INT, SMF42RUTOK BLOB,
		SMF42RDSCB_OLD BLOB, SMF42RDSCB_NEW BLOB`)
}

// Subtype 1: BMF totals and storage class summary.
func smf042s01(rec []byte, id int) {
	if s := smf042First(rec, 0); s != nil {
		insert(pStmt042Bmf, append([]interface{}{id}, append(s.u32s(0, 8), s.u16(32), s.u16(34))...)...)
	}
	smf042Each(rec, 1, func(s section) {
		insert(pStmt042BmfSc, id, s.text(2, s.num(0, 2)), s.u32(32), s.u32(36), s.u32(40), s.u32(44))
	})
}

// Subtype 2: control unit cache and SMS volume status.
func smf042s02(rec []byte, id int) {
	smf042Each(rec, 0, func(s section) {
		insert(pStmt042Cu, id, s.u8(0), s.u8(1), s.u16(2),
			s.u32(4), s.u32(8), s.u32(12), s.u32(16), s.u32(20), s.u32(24),
			s.u32(32), s.u32(36), s.u32(40), s.u32(44), s.u16(48), s.u16(50), s.u32(52),
			s.u32(56), s.u32(60), s.u32(64), s.u32(68), s.u16(72), s.u16(74), s.u32(76),
			s.u16(80), s.u16(82))
	})
	// SMF42VLL is the total length of all volume sections, not the length of one.
	off, total, n := smf042Triplet(rec, 1)
	if n == 0 {
		return
	}
	size := total / n
	for k := 0; k < n; k++ {
		s := section(rec).sub(off+k*size, size)
		insert(pStmt042CuVol, id, s.text(0, 6), s.u32(8), s.u8(12), s.u8(13))
	}
}

// Subtype 3: SMS configuration change (event audit).
func smf042s03(rec []byte, id int) {
	smf042Each(rec, 0, func(s section) {
		insert(pStmt042Ea, id, s.text(0, 8), s.u32(8), s.u32(12), s.u32(16), s.text(20, 6),
			s.u8(28), s.u8(29), s.text(32, 8), s.text(42, s.num(40, 2)), s.text(72, 6),
			s.text(80, 8), s.text(88, 12), s.text(104, 44), s.text(148, 44))
	})
}

// Subtype 4: System Data Mover statistics. Only the extended format data set (EXT) sections are decoded;
// the concurrent copy and virtual concurrent copy sessions are not in the sample data.
func smf042s04(rec []byte, id int) {
	smf042Each(rec, 1, func(s section) {
		insert(pStmt042Ext, id, s.u32(0), s.text(4, 3), s.text(8, 8), s.text(16, 8), s.tod(24), s.tod(32),
			s.text(40, 1), s.u64(44), s.u64(52))
	})
}

// Subtype 9: B37/D37/E37 abend, with the SMS section when the data set is SMS managed.
func smf042s09(rec []byte, id int) {
	s := smf042First(rec, 0)
	if s == nil {
		return
	}
	sms := smf042First(rec, 1)
	insert(pStmt042X37, id, s.text(0, 4), s.text(4, 8), s.time(12), s.date(16), s.text(20, 8), s.u8(28), s.u8(29),
		s.u16(34), s.u8(36), s.text(37, 44), s.text(81, 6), s.u32(87), s.u8(91), s.u32(92), s.u32(96), s.u24(100),
		sms.text(0, 30), sms.text(30, 30), sms.text(60, 30))
}

// Subtype 21: PDS/PDSE member delete.
func smf042s21(rec []byte, id int) {
	s := smf042First(rec, 0)
	if s == nil {
		return
	}
	ui := smf042First(rec, 2)
	insert(pStmt042Mdel, id, s.text(0, 8), s.text(8, 8), s.text(16, 8), s.text(24, 44), s.text(68, 6),
		s.u32(76), s.text(80, s.num(74, 2)), ui.blob(0, 80)) // SMF42LUI: RACF token (ICHRUTKN), binary
	smf042Aliases(rec, 1, func(name interface{}) { insert(pStmt042MdelAlias, id, name) })
}

// Subtype 24: PDS/PDSE member add/replace.
func smf042s24(rec []byte, id int) {
	s := smf042First(rec, 0)
	if s == nil {
		return
	}
	ui := smf042First(rec, 2)
	insert(pStmt042Madd, id, s.text(0, 8), s.text(8, 8), s.text(16, 8), s.text(24, 44), s.text(68, 6),
		s.u8(76), s.text(80, s.num(74, 2)), ui.blob(0, 80)) // SMF42PUI: RACF token (ICHRUTKN), binary
	smf042Aliases(rec, 1, func(name interface{}) { insert(pStmt042MaddAlias, id, name) })
}

// smf042Aliases walks an alias names section: a 2-byte count, then (2-byte length, EBCDIC name) entries.
func smf042Aliases(rec []byte, i int, f func(name interface{})) {
	s := smf042First(rec, i)
	off := 2
	for k := 0; k < s.num(0, 2); k++ {
		n := s.num(off, 2)
		f(s.text(off+2, n))
		off += 2 + n
	}
}

// Subtype 27: VTOC audit log, with the old and new DSCB when present.
func smf042s27(rec []byte, id int) {
	s := smf042First(rec, 0)
	if s == nil {
		return
	}
	old := smf042First(rec, 1)
	new := smf042First(rec, 2)
	insert(pStmt042Vtoc, id, s.text(0, 8), s.text(8, 8), s.text(16, 8), s.text(24, 8), s.text(32, 6), s.u16(38),
		s.text(40, 4), s.u8(44), s.u8(45), s.u32(48), s.blob(52, 5), s.blob(57, 15), s.u32(72), s.blob(76, 80),
		old.blob(0, 140), new.blob(0, 140))
}
