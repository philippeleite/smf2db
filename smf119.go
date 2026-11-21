package main

import (
	"database/sql"
	"sync"
)

// Type 119 (TCP/IP statistics). The layouts come from z/OS Communications Server: IP Programmer's Guide and
// Reference (SC27-3659, data/halx001_v3r1.pdf, appendix "Type 119 SMF records"). Every record has a
// triplet count at rec[20] and triplets from rec[24]: the first one locates the TCP/IP identification
// section, stored with the header in smf119; the others locate the subtype's sections, at
// rec[24+8*k] for section k. A subtype whose sections occur once per record gets one row (sections
// that are absent leave their columns NULL); interface and server port sections get one row each,
// with SEQ. Variable-length sections (data set, user and host names) are stored as text.
// Not decoded, as the sample data has none: every subtype except 1, 2, 3, 5, 6, 7, 10, 70 and 72,
// the AT-TLS and IP filter sections of subtype 2, the SOCKS and load module sections of FTP records,
// the IPv6 IP and ICMP statistics of subtype 5 and the IPv6 HOME address sections of subtype 6.
// The field tables were generated from the manual with tools/genfields.py (IP addresses, times,
// dates, TOD clocks and SSL session IDs fixed by hand).

var pStmt119 *sql.Stmt
var once119 sync.Once
var count119 int

var smf119HeaderFields = fields{
	{"SMF119HDTime", 2, 'T', 4},
	{"SMF119HDDate", 6, 'd', 4},
	{"SMF119HDSID", 10, 'c', 4},
	{"SMF119HDSubType", 18, 'u', 2},
}

// smf119Part is one section of a one-row-per-record table: k is its triplet number (1 for the first
// subtype section), f its fields, or, when f is nil, name is the column holding the whole section as text.
type smf119Part struct {
	k    int
	f    fields
	name string
}

type smf119Table struct {
	subtypes []int
	name     string
	parts    []smf119Part
	stmt     *sql.Stmt
}

var smf119Tables = []*smf119Table{
	{[]int{1}, "smf119_conn_init", []smf119Part{{1, smf119InitFields, ""}}, nil},
	{[]int{2}, "smf119_conn_term", []smf119Part{{1, smf119TermFields, ""}, {2, smf119TelnetFields, ""},
		{4, fields{{"SMF119AP_TTAPPLDATA", 0, 'b', 40}}, ""}}, nil},
	{[]int{3}, "smf119_ftpc", []smf119Part{{1, smf119FtpcFields, ""}, {2, nil, "SMF119FT_FCFileName"},
		{4, smf119FtpcSecFields, ""}, {5, nil, "SMF119FT_FCIUserID"}}, nil},
	{[]int{5}, "smf119_stat", []smf119Part{{1, smf119IpFields, ""}, {2, smf119TcpFields, ""}, {3, smf119UdpFields, ""},
		{4, smf119IcmpFields, ""}, {7, smf119StorFields, ""}}, nil},
	{[]int{10}, "smf119_udp_close", []smf119Part{{1, smf119UdpCloseFields, ""}}, nil},
	{[]int{70}, "smf119_ftps", []smf119Part{{1, smf119FtpsFields, ""}, {2, nil, "SMF119FT_FSHostname"},
		{3, nil, "SMF119FT_FSFileName1"}, {4, nil, "SMF119FT_FSFileName2"}, {5, smf119FtpsSecFields, ""}}, nil},
	{[]int{72}, "smf119_ftps_logon", []smf119Part{{1, smf119LogonFields, ""}, {2, smf119LogonSecFields, ""}}, nil},
}

var smf119Sections = []*rmfSection{
	{6, 32, "smf119_intf", smf119IntfFields, nil},
	{7, 32, "smf119_tcp_port", smf119TcpPortFields, nil},
	{7, 40, "smf119_udp_port", smf119UdpPortFields, nil},
}

func smf119Init() {
	pStmt119 = createTable("smf119", "ID INT NOT NULL, "+smf119HeaderFields.columns()+", "+smf119IdFields.columns())
	for _, t := range smf119Tables {
		cols := "ID INT NOT NULL"
		for _, p := range t.parts {
			if p.f == nil {
				cols += ", " + p.name + " TEXT"
			} else {
				cols += ", " + p.f.columns()
			}
		}
		t.stmt = createTable(t.name, cols)
	}
	rmfCreateSections(smf119Sections)
}

// smf119Section returns the first section located by triplet k, or nil when absent.
func smf119Section(rec []byte, k int) section {
	off, length, n := tripletAt(rec, 24+8*k)
	if n == 0 {
		return nil
	}
	return section(rec).sub(off, length)
}

func smf119Parse(rec []byte) {
	once119.Do(smf119Init)

	count119 += 1
	id := count119
	s := section(rec)
	insert(pStmt119, append(append([]interface{}{id}, smf119HeaderFields.values(s)...), smf119IdFields.values(smf119Section(rec, 0))...)...)

	subtype := s.num(18, 2)
	for _, t := range smf119Tables {
		for _, st := range t.subtypes {
			if st != subtype {
				continue
			}
			row := []interface{}{id}
			for _, p := range t.parts {
				sec := smf119Section(rec, p.k)
				if p.f != nil {
					row = append(row, p.f.values(sec)...)
				} else if sec != nil {
					row = append(row, ebcdic(sec))
				} else {
					row = append(row, nil)
				}
			}
			insert(t.stmt, row...)
		}
	}
	rmfInsertSections(smf119Sections, rec, id)
}

var smf119IdFields = fields{
	{"SMF119TI_SYSName", 0, 'c', 8},
	{"SMF119TI_SysplexName", 8, 'c', 8},
	{"SMF119TI_Stack", 16, 'c', 8},
	{"SMF119TI_ReleaseID", 24, 'c', 8},
	{"SMF119TI_Comp", 32, 'c', 8},
	{"SMF119TI_ASName", 40, 'c', 8},
	{"SMF119TI_UserID", 48, 'c', 8},
	{"SMF119TI_ASID2", 58, 'u', 2},
	// skipped SMF119TI_ASID at 58: redefines an earlier field
	{"SMF119TI_Reason", 60, 'u', 1},
	{"SMF119TI_RecordID", 61, 'u', 1},
}

var smf119InitFields = fields{
	{"SMF119AP_TIRName", 0, 'c', 8},
	{"SMF119AP_TIConnID", 8, 'u', 4},
	{"SMF119AP_TISubTask", 16, 'u', 4},
	{"SMF119AP_TIRIP", 20, 'a', 16},
	{"SMF119AP_TILIP", 36, 'a', 16},
	{"SMF119AP_TIRPort", 52, 'u', 2},
	{"SMF119AP_TILPort", 54, 'u', 2},
	{"SMF119AP_TITime", 56, 'T', 4},
	{"SMF119AP_TIDate", 60, 'd', 4},
	{"SMF119AP_TISTCK", 64, 't', 8},
}

var smf119TermFields = fields{
	{"SMF119AP_TTRName", 0, 'c', 8},
	{"SMF119AP_TTConnID", 8, 'u', 4},
	{"SMF119AP_TTTTLSCS", 12, 'u', 1},
	{"SMF119AP_TTTTLSPS", 13, 'u', 1},
	{"SMF119AP_TTTermCode", 14, 'u', 1},
	{"SMF119AP_TTSMCStatus", 15, 'u', 1},
	{"SMF119AP_TTSubtask", 16, 'u', 4},
	{"SMF119AP_TTSTime", 20, 'T', 4},
	{"SMF119AP_TTSDate", 24, 'd', 4},
	{"SMF119AP_TTETime", 28, 'T', 4},
	{"SMF119AP_TTEDate", 32, 'd', 4},
	{"SMF119AP_TTRIP", 36, 'a', 16},
	{"SMF119AP_TTLIP", 52, 'a', 16},
	{"SMF119AP_TTRPort", 68, 'u', 2},
	{"SMF119AP_TTLPort", 70, 'u', 2},
	{"SMF119AP_TTInBytes", 72, 'u', 8},
	{"SMF119AP_TTOutBytes", 80, 'u', 8},
	{"SMF119AP_TTSWS", 88, 'u', 4},
	{"SMF119AP_TTMSWS", 92, 'u', 4},
	{"SMF119AP_TTCWS", 96, 'u', 4},
	{"SMF119AP_TTSMS", 100, 'u', 4},
	{"SMF119AP_TTRTT", 104, 'u', 4},
	{"SMF119AP_TTRVA", 108, 'u', 4},
	{"SMF119AP_TTStatus", 112, 'u', 1},
	{"SMF119AP_TTTOS", 113, 'u', 1},
	{"SMF119AP_TTXRT", 114, 'u', 2},
	{"SMF119AP_TTProf", 116, 'c', 32},
	{"SMF119AP_TTPol", 148, 'c', 32},
	{"SMF119AP_TTInSeg", 180, 'u', 8},
	{"SMF119AP_TTOutSeg", 188, 'u', 8},
	{"SMF119AP_TTSSTCK", 196, 't', 8},
	{"SMF119AP_TTESTCK", 204, 't', 8},
	{"SMF119AP_TTDupAcksRcvd", 212, 'u', 4},
	{"SMF119AP_TTLclSMCLinkId", 216, 'u', 4},
	{"SMF119AP_TTSMCReason", 224, 'u', 4},
	{"SMF119AP_TTSMCFlags", 228, 'u', 1},
	{"SMF119AP_TTSMCDStatus", 229, 'u', 1},
	{"SMF119AP_TTSMCFlags2", 230, 'u', 1},
	{"SMF119AP_TTLclSMCBufSz", 232, 'u', 4},
	{"SMF119AP_TTRmtSMCBufSz", 236, 'u', 4},
	{"SMF119AP_TTXRT32", 240, 'u', 4},
	{"SMF119AP_TTSMCDReason", 244, 'u', 4},
	{"SMF119AP_TTSMCv2Eid", 248, 'c', 32},
	// skipped SMF119AP_TTSMCDEid at 248: redefines an earlier field
}

var smf119TelnetFields = fields{
	{"SMF119AP_TTTelLUName", 0, 'c', 8},
	{"SMF119AP_TTTelAppl", 8, 'c', 8},
	{"SMF119AP_TTTelLogmode", 16, 'c', 8},
	{"SMF119AP_TTTelStatus", 24, 'u', 4},
	{"SMF119AP_TTTelTermCode", 28, 'u', 1},
}

var smf119FtpcFields = fields{
	{"SMF119FT_FCCmd", 0, 'c', 4},
	{"SMF119FT_FCFType", 4, 'c', 4},
	{"SMF119FT_FCDRIP", 8, 'a', 16},
	{"SMF119FT_FCDLIP", 24, 'a', 16},
	{"SMF119FT_FCDRPort", 40, 'u', 2},
	{"SMF119FT_FCDLPort", 42, 'u', 2},
	{"SMF119FT_FCCRIP", 44, 'a', 16},
	{"SMF119FT_FCCLIP", 60, 'a', 16},
	{"SMF119FT_FCCRPort", 76, 'u', 2},
	{"SMF119FT_FCCLPort", 78, 'u', 2},
	{"SMF119FT_FCRUser", 80, 'c', 8},
	{"SMF119FT_FCLUser", 88, 'c', 8},
	{"SMF119FT_FCType", 96, 'c', 1},
	{"SMF119FT_FCMode", 97, 'c', 1},
	{"SMF119FT_FCStruct", 98, 'c', 1},
	{"SMF119FT_FCDSType", 99, 'c', 1},
	{"SMF119FT_FCSTime", 100, 'T', 4},
	{"SMF119FT_FCSDate", 104, 'd', 4},
	{"SMF119FT_FCETime", 108, 'T', 4},
	{"SMF119FT_FCEDate", 112, 'd', 4},
	{"SMF119FT_FCDur", 116, 'u', 4},
	{"SMF119FT_FCBytes", 120, 'u', 8},
	{"SMF119FT_FCLReply", 128, 'c', 4},
	{"SMF119FT_FCM1", 132, 'c', 8},
	{"SMF119FT_FCHostname", 140, 'c', 8},
	{"SMF119FT_FCBytesFloat", 156, 'f', 8},
	{"SMF119FT_FCCConnID", 164, 'u', 4},
	{"SMF119FT_FCDConnID", 168, 'u', 4},
}

var smf119FtpcSecFields = fields{
	{"SMF119FT_FCMechanism", 0, 'c', 1},
	{"SMF119FT_FCCProtect", 1, 'c', 1},
	{"SMF119FT_FCDProtect", 2, 'c', 1},
	{"SMF119FT_FCLoginMech", 3, 'c', 1},
	{"SMF119FT_FCProtoLevel", 4, 'c', 8},
	{"SMF119FT_FCCipherSpec", 12, 'c', 20},
	{"SMF119FT_FCProtBuffSize", 32, 'u', 4},
	{"SMF119FT_FCCipher", 36, 'c', 2},
	{"SMF119FT_FCFips140", 38, 'u', 1},
	{"SMF119FT_FCCipher4", 39, 'c', 4},
	{"SMF119FT_FCSessReuse", 43, 'c', 1},
	{"SMF119FT_FCCSSLSessID", 46, 'x', 32},
	{"SMF119FT_FCDSSLSessID", 80, 'x', 32},
	{"SMF119FT_FCKeyShare", 112, 'c', 4},
}

var smf119IpFields = fields{
	{"SMF119AP_TSIPDuration", 0, 'u', 8},
	{"SMF119AP_TSIPRecData", 8, 'u', 4},
	{"SMF119AP_TSIPDscData", 12, 'u', 4},
	{"SMF119AP_TSIPDscDAddr", 16, 'u', 4},
	{"SMF119AP_TSIPAttFwdData", 20, 'u', 4},
	{"SMF119AP_TSIPDscDUnkPr", 24, 'u', 4},
	{"SMF119AP_TSIPDscDOth", 28, 'u', 4},
	{"SMF119AP_TSIPDlvData", 32, 'u', 4},
	{"SMF119AP_TSIPXData", 36, 'u', 4},
	{"SMF119AP_TSIPXDscOth", 40, 'u', 4},
	{"SMF119AP_TSIPXDscRoute", 44, 'u', 4},
	{"SMF119AP_TSIPTimeouts", 48, 'u', 4},
	{"SMF119AP_TSIPRecDRsbm", 52, 'u', 4},
	{"SMF119AP_TSIPRsmb", 56, 'u', 4},
	{"SMF119AP_TSIPFailRsmb", 60, 'u', 4},
	{"SMF119AP_TSIPRecFgmt", 64, 'u', 4},
	{"SMF119AP_TSIPDscDFgmt", 68, 'u', 4},
	{"SMF119AP_TSIPXFgmt", 72, 'u', 4},
	{"SMF119AP_TSIPRouteDisc", 76, 'u', 4},
	{"SMF119AP_TSIPMaxRsmb", 80, 'u', 4},
	{"SMF119AP_TSIPCurRsmb", 84, 'u', 4},
	{"SMF119AP_TSIPRsmbFlags", 88, 'u', 4},
	{"SMF119AP_TSIPInCalls", 92, 'u', 4},
	{"SMF119AP_TSIPInUerrs", 96, 'u', 4},
	{"SMF119AP_TSIPIDMem", 100, 'u', 4},
	{"SMF119AP_TSIPODSync", 104, 'u', 4},
	{"SMF119AP_TSIPODAsyn", 108, 'u', 4},
	{"SMF119AP_TSIPODMem", 112, 'u', 4},
}

var smf119TcpFields = fields{
	{"SMF119AP_TSTCDuration", 0, 'u', 8},
	{"SMF119AP_TSTCAlg", 8, 'u', 4},
	{"SMF119AP_TSTCMinRet", 12, 'u', 4},
	{"SMF119AP_TSTCMxRet", 16, 'u', 4},
	{"SMF119AP_TSTCMxCon", 20, 'u', 4},
	{"SMF119AP_TSTCOpenConn", 24, 'u', 4},
	{"SMF119AP_TSTCPassConn", 28, 'u', 4},
	{"SMF119AP_TSTCOFails", 32, 'u', 4},
	{"SMF119AP_TSTCConReset", 36, 'u', 4},
	{"SMF119AP_TSTCEstab", 40, 'u', 4},
	{"SMF119AP_TSTCInSegs", 44, 'u', 4},
	{"SMF119AP_TSTCOSegs", 48, 'u', 4},
	{"SMF119AP_TSTCRxSegs", 52, 'u', 4},
	{"SMF119AP_TSTCInErrs", 56, 'u', 4},
	{"SMF119AP_TSTCReset", 60, 'u', 4},
	{"SMF119AP_TSTCConCls", 64, 'u', 4},
	{"SMF119AP_TSTCConAttD", 68, 'u', 4},
	{"SMF119AP_TSTCTWRef", 72, 'u', 4},
	{"SMF119AP_TSTCHOKAck", 76, 'u', 4},
	{"SMF119AP_TSTCHOKDat", 80, 'u', 4},
	{"SMF119AP_TSTCIDupAck", 84, 'u', 4},
	{"SMF119AP_TSTCDscChecksum", 88, 'u', 4},
	{"SMF119AP_TSTCDscLen", 92, 'u', 4},
	{"SMF119AP_TSTCDscInsData", 96, 'u', 4},
	{"SMF119AP_TSTCDscOldTime", 100, 'u', 4},
	{"SMF119AP_TSTCICmpDupSeg", 104, 'u', 4},
	{"SMF119AP_TSTCIPartDupSeg", 108, 'u', 4},
	{"SMF119AP_TSTCICmpSegsWin", 112, 'u', 4},
	{"SMF119AP_TSTCIPartSegsWin", 116, 'u', 4},
	{"SMF119AP_TSTCIOOrder", 120, 'u', 4},
	{"SMF119AP_TSTCISegCls", 124, 'u', 4},
	{"SMF119AP_TSTCIWinPr", 128, 'u', 4},
	{"SMF119AP_TSTCIWinUp", 132, 'u', 4},
	{"SMF119AP_TSTCOWinPr", 136, 'u', 4},
	{"SMF119AP_TSTCOWinUp", 140, 'u', 4},
	{"SMF119AP_TSTCODlAck", 144, 'u', 4},
	{"SMF119AP_TSTCOKApr", 148, 'u', 4},
	{"SMF119AP_TSTCRxTim", 152, 'u', 4},
	{"SMF119AP_TSTCRxMTU", 156, 'u', 4},
	{"SMF119AP_TSTCPathM", 160, 'u', 4},
	{"SMF119AP_TSTCDropPr", 164, 'u', 4},
	{"SMF119AP_TSTCDropKA", 168, 'u', 4},
	{"SMF119AP_TSTCDropF2", 172, 'u', 4},
	{"SMF119AP_TSTCDropRx", 176, 'u', 4},
	{"SMF119AP_TSTCEphPortExh", 180, 'u', 4},
	{"SMF119AP_TSTCEphPortAvail", 184, 'u', 2},
	{"SMF119AP_TSTCEphPortInUse", 186, 'u', 2},
	{"SMF119AP_TSTCEphPortMxUse", 188, 'u', 2},
	{"SMF119AP_TSSMCRCurrEstabLnks", 192, 'u', 4},
	{"SMF119AP_TSSMCRLnkActTimeOut", 196, 'u', 4},
	{"SMF119AP_TSSMCRActLnkOpened", 200, 'u', 4},
	{"SMF119AP_TSSMCRPasLnkOpened", 204, 'u', 4},
	{"SMF119AP_TSSMCRLnksClosed", 208, 'u', 4},
	{"SMF119AP_TSSMCRCurrEstab", 212, 'u', 4},
	{"SMF119AP_TSSMCRActiveOpened", 216, 'u', 4},
	{"SMF119AP_TSSMCRPassiveOpened", 220, 'u', 4},
	{"SMF119AP_TSSMCRConnClosed", 224, 'u', 4},
	{"SMF119AP_TSSMCRInSegs", 232, 'u', 8},
	{"SMF119AP_TSSMCROutSegs", 240, 'u', 8},
	{"SMF119AP_TSSMCRInRsts", 248, 'u', 4},
	{"SMF119AP_TSSMCROutRsts", 252, 'u', 4},
}

var smf119UdpFields = fields{
	{"SMF119AP_TSUDDuration", 0, 'u', 8},
	{"SMF119AP_TSUDRecData", 8, 'u', 8},
	{"SMF119AP_TSUDRecNoPort", 16, 'u', 4},
	{"SMF119AP_TSUDNoRec", 20, 'u', 4},
	{"SMF119AP_TSUDXmtData", 24, 'u', 8},
	{"SMF119AP_TSUDEphPortExh", 32, 'u', 4},
	{"SMF119AP_TSUDEphPortAvail", 36, 'u', 2},
	{"SMF119AP_TSUDEphPortInUse", 38, 'u', 2},
	{"SMF119AP_TSUDEphPortMxUse", 40, 'u', 2},
}

var smf119IcmpFields = fields{
	{"SMF119AP_TSICDuration", 0, 'u', 8},
	{"SMF119AP_TSICInMsg", 8, 'u', 4},
	{"SMF119AP_TSICInError", 12, 'u', 4},
	{"SMF119AP_TSICInDstUnreach", 16, 'u', 4},
	{"SMF119AP_TSICInTimeExcd", 20, 'u', 4},
	{"SMF119AP_TSICInParmProb", 24, 'u', 4},
	{"SMF119AP_TSICInSrcQuench", 28, 'u', 4},
	{"SMF119AP_TSICInRedirect", 32, 'u', 4},
	{"SMF119AP_TSICInEcho", 36, 'u', 4},
	{"SMF119AP_TSICInEchoRep", 40, 'u', 4},
	{"SMF119AP_TSICInTstamp", 44, 'u', 4},
	{"SMF119AP_TSICInTstampRep", 48, 'u', 4},
	{"SMF119AP_TSICInAddrMask", 52, 'u', 4},
	{"SMF119AP_TSICInAddrMRep", 56, 'u', 4},
	{"SMF119AP_TSICOutMsg", 60, 'u', 4},
	{"SMF119AP_TSICOutError", 64, 'u', 4},
	{"SMF119AP_TSICOutDstUnreach", 68, 'u', 4},
	{"SMF119AP_TSICOutTimeExcd", 72, 'u', 4},
	{"SMF119AP_TSICOutParmProb", 76, 'u', 4},
	{"SMF119AP_TSICOutSrcQuench", 80, 'u', 4},
	{"SMF119AP_TSICOutRedirect", 84, 'u', 4},
	{"SMF119AP_TSICOutEcho", 88, 'u', 4},
	{"SMF119AP_TSICOutEchoRep", 92, 'u', 4},
	{"SMF119AP_TSICOutTstamp", 96, 'u', 4},
	{"SMF119AP_TSICOutTstampRep", 100, 'u', 4},
	{"SMF119AP_TSICOutAddrMask", 104, 'u', 4},
	{"SMF119AP_TSICOutAddrMRep", 108, 'u', 4},
}

var smf119StorFields = fields{
	{"SMF119AP_TSSTECSACurrent", 0, 'u', 8},
	{"SMF119AP_TSSTECSAFree", 8, 'u', 8},
	{"SMF119AP_TSSTPrivateCurrent", 16, 'u', 8},
	{"SMF119AP_TSSTPrivateFree", 24, 'u', 8},
	{"SMF119AP_TSSTSMCRFixedCurrent", 32, 'u', 8},
	{"SMF119AP_TSSTSMCRFixedMax", 40, 'u', 8},
	{"SMF119AP_TSSTSMCRSendCurrent", 48, 'u', 8},
	{"SMF119AP_TSSTSMCRSendMax", 56, 'u', 8},
	{"SMF119AP_TSSTSMCRRecvCurrent", 64, 'u', 8},
	{"SMF119AP_TSSTSMCRRecvMax", 72, 'u', 8},
	{"SMF119AP_TSSTSMCDFixedCurrent", 80, 'u', 8},
	{"SMF119AP_TSSTSMCDFixedMax", 88, 'u', 8},
	{"SMF119AP_TSSTZAGGCURRENT", 96, 'u', 8},
	{"SMF119AP_TSSTZAGGMAX", 104, 'u', 8},
}

var smf119IntfFields = fields{
	{"SMF119IS_IFDuration", 0, 'u', 8},
	{"SMF119IS_IFLnkHome", 8, 'a', 16},
	{"SMF119IS_IFName", 24, 'c', 16},
	{"SMF119IS_IFDevName", 40, 'c', 16},
	{"SMF119IS_IFDesc", 56, 'c', 18},
	{"SMF119IS_IFFlags", 74, 'u', 1},
	{"SMF119IS_IFActualMtu", 76, 'u', 4},
	{"SMF119IS_IFSPeed", 80, 'u', 4},
	{"SMF119IS_IFHSpeed", 84, 'u', 4},
	{"SMF119IS_IFInBytes", 88, 'u', 8},
	{"SMF119IS_IFInUniC", 96, 'u', 8},
	{"SMF119IS_IFInBroadC", 104, 'u', 8},
	{"SMF119IS_IFInMultiC", 112, 'u', 8},
	{"SMF119IS_IFInDisc", 120, 'u', 4},
	{"SMF119IS_IFInError", 124, 'u', 4},
	{"SMF119IS_IFInUProt", 128, 'u', 4},
	{"SMF119IS_IFOutBytes", 132, 'u', 8},
	{"SMF119IS_IFOutUniC", 140, 'u', 8},
	{"SMF119IS_IFOutBroadC", 148, 'u', 8},
	{"SMF119IS_IFOutMultiC", 156, 'u', 8},
	{"SMF119IS_IFOutDisc", 164, 'u', 4},
	{"SMF119IS_IFOutError", 168, 'u', 4},
	{"SMF119IS_IFOQL", 172, 'u', 4},
	{"SMF119IS_IFPNetID", 224, 'c', 16},
}

var smf119TcpPortFields = fields{
	{"SMF119SP_TCDuration", 0, 'u', 8},
	{"SMF119SP_TCRName", 8, 'c', 8},
	{"SMF119SP_TCBindIP", 16, 'a', 16},
	{"SMF119SP_TCPort", 32, 'u', 2},
	{"SMF119SP_TCConn", 36, 'u', 4},
	{"SMF119SP_TCBinds", 40, 'u', 4},
	{"SMF119SP_TCBusySrv", 44, 'u', 4},
	{"SMF119SP_TCSynAttack", 48, 'u', 4},
	{"SMF119SP_TCHighwater", 52, 'u', 4},
	{"SMF119SP_TCNumConns", 56, 'u', 4},
}

var smf119UdpPortFields = fields{
	{"SMF119SP_UDDuration", 0, 'u', 8},
	{"SMF119SP_UDRName", 8, 'c', 8},
	{"SMF119SP_UDBindIP", 16, 'a', 16},
	{"SMF119SP_UDPort", 32, 'u', 2},
	{"SMF119SP_UDIDgrams", 36, 'u', 8},
	{"SMF119SP_UDODgrams", 44, 'u', 8},
	{"SMF119SP_UDIBytes", 52, 'u', 8},
	{"SMF119SP_UDOBytes", 60, 'u', 8},
}

var smf119UdpCloseFields = fields{
	{"SMF119UD_UCRname", 0, 'c', 8},
	{"SMF119UD_UCConnID", 8, 'u', 4},
	{"SMF119UD_UCSubTask", 12, 'u', 4},
	{"SMF119UD_UCOTime", 16, 'T', 4},
	{"SMF119UD_UCODate", 20, 'd', 4},
	{"SMF119UD_UCCTime", 24, 'T', 4},
	{"SMF119UD_UCCDate", 28, 'd', 4},
	{"SMF119UD_UCRIP", 32, 'a', 16},
	{"SMF119UD_UCLIP", 48, 'a', 16},
	{"SMF119UD_UCRPort", 64, 'u', 2},
	{"SMF119UD_UCLPort", 66, 'u', 2},
	{"SMF119UD_UCType", 68, 'u', 1},
	{"SMF119UD_UCReason", 69, 'u', 1},
	{"SMF119UD_UCInDgrams", 72, 'u', 8},
	{"SMF119UD_UCOutDgrams", 80, 'u', 8},
	{"SMF119UD_UCInBytes", 88, 'u', 8},
	{"SMF119UD_UCOutBytes", 96, 'u', 8},
}

var smf119FtpsFields = fields{
	{"SMF119FT_FSOper", 0, 'u', 1},
	{"SMF119FT_FSCmd", 4, 'c', 4},
	{"SMF119FT_FSFType", 8, 'c', 4},
	{"SMF119FT_FSDRIP", 12, 'a', 16},
	{"SMF119FT_FSDLIP", 28, 'a', 16},
	{"SMF119FT_FSDRPort", 44, 'u', 2},
	{"SMF119FT_FSDLPort", 46, 'u', 2},
	{"SMF119FT_FSCRIP", 48, 'a', 16},
	{"SMF119FT_FSCLIP", 64, 'a', 16},
	{"SMF119FT_FSCRPort", 80, 'u', 2},
	{"SMF119FT_FSCLPort", 82, 'u', 2},
	{"SMF119FT_FSSUser", 84, 'c', 8},
	{"SMF119FT_FSType", 92, 'c', 1},
	{"SMF119FT_FSMode", 93, 'c', 1},
	{"SMF119FT_FSStruct", 94, 'c', 1},
	{"SMF119FT_FSDsType", 95, 'c', 1},
	{"SMF119FT_FSSTime", 96, 'T', 4},
	{"SMF119FT_FSSDate", 100, 'd', 4},
	{"SMF119FT_FSETime", 104, 'T', 4},
	{"SMF119FT_FSEDate", 108, 'd', 4},
	{"SMF119FT_FSDur", 112, 'u', 4},
	{"SMF119FT_FSBytes", 116, 'u', 8},
	{"SMF119FT_FSLReply", 124, 'c', 4},
	{"SMF119FT_FSM1", 128, 'c', 8},
	{"SMF119FT_FSM2", 144, 'c', 8},
	{"SMF119FT_FSBytesFloat", 152, 'f', 8},
	{"SMF119FT_FSCConnID", 160, 'u', 4},
	{"SMF119FT_FSDConnID", 164, 'u', 4},
	{"SMF119FT_FSSessionID", 168, 'c', 15},
}

var smf119FtpsSecFields = fields{
	{"SMF119FT_FSMechanism", 0, 'c', 1},
	{"SMF119FT_FSCProtect", 1, 'c', 1},
	{"SMF119FT_FSDProtect", 2, 'c', 1},
	{"SMF119FT_FSLoginMech", 3, 'c', 1},
	{"SMF119FT_FSProtoLevel", 4, 'c', 8},
	{"SMF119FT_FSCipherSpec", 12, 'c', 20},
	{"SMF119FT_FSProtoBufSize", 32, 'u', 4},
	{"SMF119FT_FSCipher", 36, 'c', 2},
	{"SMF119FT_FSFips140", 38, 'u', 1},
	{"SMF119FT_FSCipher4", 39, 'c', 4},
	{"SMF119FT_FSSessReuse", 43, 'c', 1},
	{"SMF119FT_FSCSSLSessID", 46, 'x', 32},
	{"SMF119FT_FSDSSLSessID", 80, 'x', 32},
	{"SMF119FT_FSKeyShare", 112, 'c', 4},
}

var smf119LogonFields = fields{
	{"SMF119FT_FFRIP", 0, 'a', 16},
	{"SMF119FT_FFLIP", 16, 'a', 16},
	{"SMF119FT_FFRPort", 32, 'u', 2},
	{"SMF119FT_FFLPort", 34, 'u', 2},
	{"SMF119FT_FFUserID", 36, 'c', 8},
	{"SMF119FT_FFReason", 44, 'u', 1},
	{"SMF119FT_FFCConnID", 48, 'u', 4},
	{"SMF119FT_FFSessionID", 52, 'c', 15},
}

var smf119LogonSecFields = fields{
	{"SMF119FT_FFMechanism", 0, 'c', 1},
	{"SMF119FT_FFCProtect", 1, 'c', 1},
	{"SMF119FT_FFDProtect", 2, 'c', 1},
	{"SMF119FT_FFLoginMech", 3, 'c', 1},
	{"SMF119FT_FFProtoLevel", 4, 'c', 8},
	{"SMF119FT_FFCipherSpec", 12, 'c', 20},
	{"SMF119FT_FFProtBuffSize", 32, 'u', 4},
	{"SMF119FT_FFCipher", 36, 'c', 2},
	{"SMF119FT_FFFips140", 38, 'u', 1},
	{"SMF119FT_FFCipher4", 39, 'c', 4},
	{"SMF119FT_FFSessReuse", 43, 'c', 1},
	{"SMF119FT_FFCSSLSessID", 46, 'x', 32},
}
