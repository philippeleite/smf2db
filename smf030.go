package main

import (
	"database/sql"
	"encoding/binary"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var p1Stmt030 *sql.Stmt
var p2Stmt030 *sql.Stmt

func smf030Init() {
	sqlStmt := `
	DROP TABLE IF EXISTS smf030;
	CREATE TABLE smf030 (
		SMF30TME   DATETIME   NOT NULL,
		SMF30DTE   DATETIME   NOT NULL,
		SMF30SID   TEXT       NOT NULL,
		SMF30WID   TEXT       NOT NULL
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	sqlStmt2 := `
	DROP TABLE IF EXISTS smf030_vol;
	CREATE TABLE smf030_vol (
		SMF30TME   DATETIME   NOT NULL,
		SMF30DTE   DATETIME   NOT NULL,
		SMF30FVL   TEXT       NOT NULL
	);
	`
	_, err = db.Exec(sqlStmt2)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt2)
	}

	sqlStmt3 := `
		INSERT INTO smf030
		VALUES (
			?,
			?,
			?,
			?
		);
`
	p1Stmt030, err = db.Prepare(sqlStmt3)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt3)
	}

	sqlStmt4 := `
		INSERT INTO smf030_vol
		VALUES (
			?,
			?,
			?
		);
`

	p2Stmt030, err = db.Prepare(sqlStmt4)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt4)
	}
}

func smf030Parse(rec []byte) {

	SMF30STP := int(binary.BigEndian.Uint16(rec[18:20]))   // SMF30STP
	SMF30SOF := int(binary.BigEndian.Uint32(rec[20:24]))   // SMF30SOF
	SMF30SLN := int(binary.BigEndian.Uint16(rec[24:26]))   // SMF30SLN
	SMF30SON := int(binary.BigEndian.Uint16(rec[26:28]))   // SMF30SON
	SMF30IOF := int(binary.BigEndian.Uint32(rec[28:32]))   // SMF30IOF
	SMF30ILN := int(binary.BigEndian.Uint16(rec[32:34]))   // SMF30ILN
	SMF30ION := int(binary.BigEndian.Uint16(rec[34:36]))   // SMF30ION
	SMF30UOF := int(binary.BigEndian.Uint32(rec[36:40]))   // SMF30UOF
	SMF30ULN := int(binary.BigEndian.Uint16(rec[40:42]))   // SMF30ULN
	SMF30UON := int(binary.BigEndian.Uint16(rec[42:44]))   // SMF30UON
	SMF30TOF := int(binary.BigEndian.Uint32(rec[44:48]))   // SMF30TOF
	SMF30TLN := int(binary.BigEndian.Uint16(rec[48:50]))   // SMF30TLN
	SMF30TON := int(binary.BigEndian.Uint16(rec[50:52]))   // SMF30TON
	SMF30COF := int(binary.BigEndian.Uint32(rec[52:56]))   // SMF30COF
	SMF30CLN := int(binary.BigEndian.Uint16(rec[56:58]))   // SMF30CLN
	SMF30CON := int(binary.BigEndian.Uint16(rec[58:60]))   // SMF30CON
	SMF30AOF := int(binary.BigEndian.Uint32(rec[60:64]))   // SMF30AOF
	SMF30ALN := int(binary.BigEndian.Uint16(rec[64:66]))   // SMF30ALN
	SMF30AON := int(binary.BigEndian.Uint16(rec[66:68]))   // SMF30AON
	SMF30ROF := int(binary.BigEndian.Uint32(rec[68:72]))   // SMF30ROF
	SMF30RLN := int(binary.BigEndian.Uint16(rec[72:74]))   // SMF30RLN
	SMF30RON := int(binary.BigEndian.Uint16(rec[74:76]))   // SMF30RON
	SMF30POF := int(binary.BigEndian.Uint32(rec[76:80]))   // SMF30POF
	SMF30PLN := int(binary.BigEndian.Uint16(rec[80:82]))   // SMF30PLN
	SMF30PON := int(binary.BigEndian.Uint16(rec[82:84]))   // SMF30PON
	SMF30OOF := int(binary.BigEndian.Uint32(rec[84:88]))   // SMF30OOF
	SMF30OLN := int(binary.BigEndian.Uint16(rec[88:90]))   // SMF30OLN
	SMF30OON := int(binary.BigEndian.Uint16(rec[90:92]))   // SMF30OON
	SMF30EOF := int(binary.BigEndian.Uint32(rec[92:96]))   // SMF30EOF
	SMF30ELN := int(binary.BigEndian.Uint16(rec[96:98]))   // SMF30ELN
	SMF30EON := int(binary.BigEndian.Uint16(rec[98:100]))  // SMF30EON
	SMF30EOR := int(binary.BigEndian.Uint16(rec[100:102])) // SMF30EOR
	SMF30EOS := int(binary.BigEndian.Uint32(rec[104:108])) // SMF30EOS
	SMF30DRO := int(binary.BigEndian.Uint32(rec[108:112])) // SMF30DRO
	SMF30DRL := int(binary.BigEndian.Uint16(rec[112:114])) // SMF30DRL
	SMF30DRN := int(binary.BigEndian.Uint16(rec[114:116])) // SMF30DRN
	SMF30ARO := int(binary.BigEndian.Uint32(rec[116:120])) // SMF30ARO
	SMF30ARL := int(binary.BigEndian.Uint16(rec[120:122])) // SMF30ARL
	SMF30ARN := int(binary.BigEndian.Uint16(rec[122:124])) // SMF30ARN
	SMF30OPO := int(binary.BigEndian.Uint32(rec[124:128])) // SMF30OPO
	SMF30OPL := int(binary.BigEndian.Uint16(rec[128:130])) // SMF30OPL
	SMF30OPN := int(binary.BigEndian.Uint16(rec[130:132])) // SMF30OPN
	SMF30OPM := int(binary.BigEndian.Uint32(rec[132:136])) // SMF30OPM
	SMF30UDO := int(binary.BigEndian.Uint32(rec[136:140])) // SMF30UDO
	SMF30UDL := int(binary.BigEndian.Uint16(rec[140:142])) // SMF30UDL
	SMF30UDN := int(binary.BigEndian.Uint16(rec[142:144])) // SMF30UDN
	SMF30UDS := int(binary.BigEndian.Uint32(rec[144:148])) // SMF30UDS
	SMF30RMO := int(binary.BigEndian.Uint32(rec[148:152])) // SMF30RMO
	SMF30RML := int(binary.BigEndian.Uint16(rec[152:154])) // SMF30RML
	SMF30RMN := int(binary.BigEndian.Uint16(rec[154:156])) // SMF30RMN
	SMF30RMS := int(binary.BigEndian.Uint32(rec[156:160])) // SMF30RMS
	SMF30MOF := int(binary.BigEndian.Uint32(rec[160:164])) // SMF30MOF
	SMF30MLN := int(binary.BigEndian.Uint16(rec[164:166])) // SMF30MLN
	SMF30MNO := int(binary.BigEndian.Uint16(rec[166:168])) // SMF30MNO
	SMF30MOS := int(binary.BigEndian.Uint32(rec[168:172])) // SMF30MOS
	SMF30CDO := int(binary.BigEndian.Uint32(rec[172:176])) // SMF30CDO
	SMF30CDL := int(binary.BigEndian.Uint16(rec[176:178])) // SMF30CDL
	SMF30CDN := int(binary.BigEndian.Uint16(rec[178:180])) // SMF30CDN
	SMF30USO := int(binary.BigEndian.Uint32(rec[180:184])) // SMF30USO
	SMF30USL := int(binary.BigEndian.Uint16(rec[184:186])) // SMF30USL
	SMF30USN := int(binary.BigEndian.Uint16(rec[186:188])) // SMF30USN
	numVols := int(rec[131])

	bcd2asc(rec[10:14]) // SMF30SID
	bcd2asc(rec[14:18]) // SMF30WID

	_, err = p1Stmt030.Exec(
		smftime(rec[2:6]),  // SMF30TME
		smfdate(rec[6:10]), // SMF30DTE
		string(rec[10:14]), // SMF30SID
		string(rec[14:18]), // SMF30WID
	)
	if err != nil {
		log.Println(err)
	}

	devOff := 132
	for i := 0; i < numVols; i++ {
		bcd2asc(rec[devOff+2 : devOff+8])
		_, err = p2Stmt030.Exec(
			smftime(rec[2:6]),              // SMF30TME
			smfdate(rec[6:10]),             // SMF30DTE
			string(rec[devOff+2:devOff+8]), // SMF30FVL
		)
		if err != nil {
			log.Println(err)
		}
		devOff += 8
	}
}
