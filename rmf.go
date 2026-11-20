package main

// rmfColumns returns the columns of an RMF record's main table: the SMF header and the RMF product
// section. p is the field name prefix, such as "SMF70".
func rmfColumns(p string) string {
	return "ID INT NOT NULL, " + p + "TME DATETIME NOT NULL, " + p + "DTE DATETIME NOT NULL, " +
		p + "SID TEXT NOT NULL, " + p + "SSI TEXT NOT NULL, " + p + "STY INT NOT NULL, " +
		rmfProductFields(p).columns()
}

// rmfRow decodes the values matching rmfColumns. The product section is located by the triplet at
// offset 28, the same in every RMF record type.
func rmfRow(rec []byte, id int, p string) []interface{} {
	s := section(rec)
	off, length, _ := tripletAt(rec, 24)
	return append([]interface{}{id, s.time(2), s.date(6), s.text(10, 4), s.text(14, 4), s.u16(18)},
		rmfProductFields(p).values(s.sub(off, length))...)
}

// rmfProductFields returns the layout of the RMF product section, which is the same in every RMF
// record type (70 to 79); p is the field name prefix, such as "SMF70". The reassembly area of broken
// records is not decoded; RAN tells whether the record is part of a broken large record.
func rmfProductFields(p string) fields {
	return fields{
		{p + "MFV", 0, 'p', 2},  // RMF version
		{p + "PRD", 2, 'c', 8},  // product name
		{p + "IST", 10, 'h', 4}, // interval start time
		{p + "DAT", 14, 'd', 4}, // interval start date
		{p + "INT", 18, 'i', 4}, // interval duration, in seconds
		{p + "SAM", 24, 'u', 4}, // number of samples
		{p + "FLA", 30, 'u', 2}, // flags
		{p + "CYC", 36, 'p', 4}, // sampling cycle, in milliseconds
		{p + "MVS", 40, 'c', 8}, // z/OS level
		{p + "IML", 48, 'u', 1},
		{p + "PRF", 49, 'u', 1},
		{p + "PTN", 50, 'u', 1}, // PR/SM partition number
		{p + "SRL", 51, 'u', 1}, // record level
		{p + "IET", 52, 'x', 8}, // interval expiration token
		{p + "LGO", 60, 's', 8}, // GMT to local time offset, in TOD clock units
		{p + "RAN", 74, 'u', 2}, // reassembly area indicator
		{p + "OIL", 76, 'u', 2}, // original interval length, in seconds
		{p + "SYN", 78, 'u', 2}, // SYNC value, in seconds
		{p + "GIE", 80, 't', 8}, // projected gathering interval end (GMT)
		{p + "XNM", 88, 'c', 8}, // sysplex name
		{p + "SNM", 96, 'c', 8}, // system name
	}
}
