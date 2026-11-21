package main

import (
	"database/sql"
	"sync"
)

// Type 118 (TCP/IP statistics, superseded by type 119). The layouts come from z/OS Communications Server:
// IP Programmer's Guide and Reference (SC27-3659, data/halx001_v3r1.pdf, appendix "Type 118 SMF records"),
// not from the SMF manual. smf118 holds the header of every record (subtype at rec[18], as in the standard
// SMF header). Flat records, whose fields start at rec[20], go to smf118_api (subtypes 1 and 2, TCP API
// initialization and termination), smf118_ftpc (3, FTP client) and smf118_telnet (20 and 21, TN3270E
// server LOGN/LOGF). Subtype 5 (TCPIPSTATISTICS) has subsystem, IP, ICMP, TCP and UDP areas located by
// triplets at the standard positions (rec[24]...); the manual's table shows them 2 bytes off, the records
// do not. Those areas go to one smf118_stat row. Not decoded, as the sample data has none: the TN3270
// client (4) and FTP server (70-75) subtypes, and the FTP client's variable z/OS UNIX file names.
// The field tables were generated from the manual with tools/genfields.py (offsets moved to rec, IPv4
// addresses, times and dates fixed by hand).

var pStmt118, pStmt118Api, pStmt118Ftpc, pStmt118Telnet, pStmt118Stat *sql.Stmt
var once118 sync.Once
var count118 int

var smf118HeaderFields = fields{
	{"SMFHDTME", 2, 'T', 4},
	{"SMFHDDTE", 6, 'd', 4},
	{"SMFHDSID", 10, 'c', 4},
	{"SMFHDSUB", 18, 'u', 2},
}

var smf118SubsysFields = fields{
	{"SMFSubProc", 0, 'c', 8},
	{"SMFSubASID", 8, 'u', 4},
	{"SMFSubTime", 12, 't', 8},
	{"SMFSubFlag", 20, 'u', 4},
}

var smf118IpFields = fields{
	{"imirecv", 0, 'u', 4}, {"imihdrer", 4, 'u', 4}, {"imiadrer", 8, 'u', 4}, {"imiunprt", 16, 'u', 4},
	{"imidisc", 20, 'u', 4}, {"imidelvr", 24, 'u', 4}, {"imoreqst", 28, 'u', 4}, {"imodisc", 32, 'u', 4},
	{"imonorte", 36, 'u', 4}, {"imrsmtos", 40, 'u', 4}, {"imrsmreq", 44, 'u', 4}, {"imrsmok", 48, 'u', 4},
	{"imrsmfld", 52, 'u', 4}, {"imfragok", 56, 'u', 4}, {"imfrgfld", 60, 'u', 4}, {"imrgcre", 64, 'u', 4},
	{"imrtdisc", 68, 'u', 4}, {"imrsmmax", 72, 'u', 4}, {"imrmsact", 76, 'u', 4}, {"imrsmful", 80, 'u', 4},
}

var smf118TcpFields = fields{
	{"tcp_RtoAlgorithm", 0, 'u', 4}, {"tcp_RtoMin", 4, 'u', 4}, {"tcp_RtoMax", 8, 'u', 4},
	{"tcp_MaxConn", 12, 'u', 4}, {"tcp_ActiveOpens", 16, 'u', 4}, {"tcp_PassiveOpens", 20, 'u', 4},
	{"tcp_AttemptFails", 24, 'u', 4}, {"tcp_EstabResets", 28, 'u', 4}, {"tcp_CurrEstab", 32, 'u', 4},
	{"tcp_InSegs", 36, 'u', 4}, {"tcp_OutSegs", 40, 'u', 4}, {"tcp_RetransSegs", 44, 'u', 4},
	{"tcp_InErrs", 48, 'u', 4}, {"tcp_OutRsts", 52, 'u', 4},
}

var smf118UdpFields = fields{
	{"usindgrm", 0, 'u', 4}, {"usnoprts", 4, 'u', 4}, {"usinerrs", 8, 'u', 4}, {"usotdgrm", 12, 'u', 4},
}

func smf118Init() {
	pStmt118 = createTable("smf118", "ID INT NOT NULL, "+smf118HeaderFields.columns())
	pStmt118Api = createTable("smf118_api", "ID INT NOT NULL, "+smf118ApiFields.columns())
	pStmt118Ftpc = createTable("smf118_ftpc", "ID INT NOT NULL, "+smf118FtpcFields.columns())
	pStmt118Telnet = createTable("smf118_telnet", "ID INT NOT NULL, "+smf118TelnetFields.columns())
	pStmt118Stat = createTable("smf118_stat", "ID INT NOT NULL, "+smf118SubsysFields.columns()+", "+
		smf118IpFields.columns()+", "+smf118TcpFields.columns()+", "+smf118UdpFields.columns())
}

func smf118Parse(rec []byte) {
	once118.Do(smf118Init)

	s := section(rec)
	count118 += 1
	id := count118
	insert(pStmt118, append([]interface{}{id}, smf118HeaderFields.values(s)...)...)

	flat := func(stmt *sql.Stmt, f fields) { insert(stmt, append([]interface{}{id}, f.values(s)...)...) }
	switch s.num(18, 2) {
	case 1, 2:
		flat(pStmt118Api, smf118ApiFields)
	case 3:
		flat(pStmt118Ftpc, smf118FtpcFields)
	case 20, 21:
		flat(pStmt118Telnet, smf118TelnetFields)
	case 5:
		area := func(at int) section {
			off, length, n := tripletAt(rec, at)
			if n == 0 {
				return nil
			}
			return s.sub(off, length)
		}
		row := append([]interface{}{id}, smf118SubsysFields.values(area(24))...)
		row = append(row, smf118IpFields.values(area(32))...)
		row = append(row, smf118TcpFields.values(area(48))...) // area(40) is ICMP, without documented fields
		insert(pStmt118Stat, append(row, smf118UdpFields.values(area(56))...)...)
	}
}

var smf118TelnetFields = fields{
	{"SMFTNTCM", 20, 'c', 4},
	{"SMFTNTLU", 24, 'c', 8},
	{"SMFTNTAP", 32, 'c', 8},
	{"SMFTNTIA", 40, 'u', 4},
	{"SMFTNTRA", 44, 'a', 4},
	{"SMFTNTLA", 48, 'a', 4},
	{"SMFTNTST", 52, 'c', 8},
	{"SMFTNTHN", 60, 'c', 8},
	{"SMFTNTIN", 70, 'u', 4},
	{"SMFTNTOU", 74, 'u', 4},
	{"SMFTNTLF", 78, 'T', 4},
	{"SMFTNTPD", 82, 'd', 4},
	{"SMFTNTRP", 86, 'u', 2},
	{"SMFTNTLP", 88, 'u', 2},
}

var smf118ApiFields = fields{
	{"SMFAPIST", 20, 'c', 4},
	{"SMFAPILA", 24, 'a', 4},
	{"SMFAPIRA", 28, 'a', 4},
	{"SMFAPILP", 32, 'u', 2},
	{"SMFAPIRP", 34, 'u', 2},
	{"SMFAPIIN", 36, 'u', 4},
	{"SMFAPIOU", 40, 'u', 4},
	{"SMFAPIUO", 44, 'u', 2},
	{"SMFAPIUL", 46, 'u', 2},
	{"SMFAPINM", 48, 'c', 8},
	{"SMFAPIJI", 56, 'c', 8},
	{"SMFAPIJS", 64, 'T', 4},
	{"SMFAPIJD", 68, 'd', 4},
}

var smf118FtpcFields = fields{
	{"SMFFTCCM", 20, 'c', 4},
	{"SMFFTCCY", 24, 'c', 4},
	{"SMFFTCSA", 28, 'a', 4},
	{"SMFFTCSL", 32, 'a', 4},
	{"SMFFTCCP", 36, 'u', 2},
	{"SMFFTCCF", 38, 'u', 2},
	{"SMFFTCSU", 44, 'c', 8},
	{"SMFFTCFM", 52, 'c', 1},
	{"SMFFTCMO", 53, 'c', 1},
	{"SMFFTCST", 54, 'c', 1},
	{"SMFFTCDT", 55, 'c', 1},
	{"SMFFTCRS", 56, 'T', 4},
	{"SMFFTCRE", 60, 'T', 4},
	{"SMFFTCBC", 64, 'u', 4},
	{"SMFFTCXD", 68, 'c', 1},
	{"SMFFTCSN", 72, 'c', 44},
	{"SMFFTCEM", 116, 'c', 8},
	{"SMFFTCTC", 184, 'c', 8},
	{"SMFFTCHN", 192, 'c', 8},
	{"SMFFTCF1", 200, 'u', 2},
	{"SMFFTCF2", 202, 'u', 2},
	{"SMFFTCYF", 204, 'f', 8},
	{"SMFFTCIG", 212, 'u', 4},
}
