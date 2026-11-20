package main

import (
	"bufio"
	"database/sql"
	"encoding/binary"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB
var err error
var version = "0.1.0"

var parseTable = map[int]func([]byte){
	4:  smf004Parse,
	5:  smf005Parse,
	6:  smf006Parse,
	14: smf014Parse,
	15: smf015Parse,
	17: smf017Parse,
	18: smf018Parse,
	21: smf021Parse,
	30: smf030Parse,
	41: smf041Parse,
	42: smf042Parse,
	60: smf060Parse,
	61: smf061Parse,
	62: smf062Parse,
	64: smf064Parse,
	65: smf065Parse,
	66: smf066Parse,
	70: smf070Parse,
}

func smfParse(rec []byte) {
	if len(rec) < 2 {
		log.Printf("skipping SMF record too short to hold a type (%d bytes)", len(rec))
		return
	}
	f, ok := parseTable[int(rec[1])]
	if !ok {
		return
	}
	// Parsers index the record by fixed and self-described offsets without bounds checks;
	// a truncated or malformed record must be skipped instead of aborting the whole run.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("skipping malformed SMF type %d record (%d bytes): %v", rec[1], len(rec), r)
		}
	}()
	f(rec)
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type subtypes [256]bool

var typeFilter = false

func (s *subtypes) String() string {
	buff := strings.Builder{}
	buff.WriteString("[ ")
	for i := 0; i < len(*s); i++ {
		buff.WriteString(fmt.Sprintf("%t ", s[i]))
	}
	buff.WriteString("]")
	return buff.String()
}

func (s *subtypes) Set(value string) error {
	typeFilter = true
	trimmed := strings.Trim(value, "()[]{}")
	conv2int := func(s string) int {
		str, err := strconv.Atoi(s)
		if err != nil {
			fmt.Printf("invalid record type: %v", s)
			os.Exit(1)
		}
		return str
	}
	for _, num := range strings.Split(trimmed, ",") {
		if strings.Contains(num, ":") {
			typeRange := strings.Split(num, ":")
			min := conv2int(typeRange[0])
			max := conv2int(typeRange[1])
			if min > max {
				fmt.Println("Invalid record type range:", num)
				os.Exit(1)
			}
			for i := min; i <= max; i++ {
				s[i] = true
			}
		} else {
			tmp := conv2int(num)
			s[tmp] = true
		}
	}
	return nil
}

var typeTable subtypes

func main() {

	// Parse command line parameters
	var help bool
	var vFlag bool
	flag.BoolVar(&help, "h", false, "Usage help")
	flag.BoolVar(&vFlag, "v", false, "Show version info")
	flag.Var(&typeTable, "t", "Record types separated by comma")
	flag.Parse()
	if help {
		printUsage()
		os.Exit(0)
	}
	if vFlag {
		fmt.Printf("smf2db - SMF to SQLite utility version %s\n", version)
	}

	var smfFile string
	var dbFile string
	switch len(flag.Args()) {
	case 1:
		smfFile = flag.Arg(0)
		dbFile = "smf.db"
	case 2:
		smfFile = flag.Arg(0)
		dbFile = flag.Arg(1)
	default:
		printUsage()
		os.Exit(8)
	}

	// SMF file
	_, err := os.Lstat(smfFile)
	if err != nil {
		fmt.Println("Invalid SMF file:", err)
		os.Exit(1)
	}

	file, err := os.Open(smfFile)
	handleErr(err)
	defer func() {
		err = file.Close()
		err = db.Close()
	}()

	// Database initializations
	_ = os.Remove(dbFile)
	var strBuff strings.Builder
	strBuff.Grow(200)
	strBuff.WriteString("file:")
	strBuff.WriteString(dbFile)
	strBuff.WriteString("?_mutex=no")
	strBuff.WriteString("&_cache_size=400000")
	strBuff.WriteString("&_synchronous=off")
	strBuff.WriteString("&_journal_mode=off")
	strBuff.WriteString("&_locking_mode=exclusive")
	strBuff.WriteString("&_temp_store=memory")
	strBuff.WriteString("&_auto_vacuum=none")
	db, err = sql.Open("sqlite3", strBuff.String())
	handleErr(err)
	_, err = db.Exec("PRAGMA count_changes = OFF;")
	handleErr(err)
	_, err = db.Exec("PRAGMA temp_store = MEMORY;")
	handleErr(err)

	// SMF parsing
	rdr := bufio.NewReaderSize(file, 32768)
	var record = make([]byte, 64*1024)
	var count = 0
	var offset uint = 0
	var matrix = make([]uint64, 256)
	var buffer = make([]byte, 2)
	var reclen, recseg uint16

	start := time.Now()
	for {
		_, err = io.ReadFull(rdr, buffer)
		if err == io.EOF {
			break
		}
		reclen = binary.BigEndian.Uint16(buffer)
		// fmt.Printf("Record length: %d (0x%04x)", reclen, reclen)

		_, err = io.ReadFull(rdr, buffer)
		if err != nil {
			fmt.Print(err)
		}
		recseg = binary.BigEndian.Uint16(buffer)
		// fmt.Printf("Record segment: %d (0x%04x)", recseg, recseg)

		// Segment descriptor:
		//  0x0000 -> record starts and ends in the block
		//  0x0100 -> record begins but does not end in this block
		//  0x0200 -> record ends but does not begin in this block
		//  0x0300 -> record neither starts nor ends in this block

		// Check if segment descriptor is valid -> recseg in (0, 256, 512, 768)
		if recseg&0xfcff != 0 {
			log.Println("Invalid file format")
			log.Fatal("SMF file must be transferred preserving RDW information")
		}
		_, err = io.ReadFull(rdr, record[offset:offset+uint(reclen)-4])
		if err != nil {
			fmt.Print(err)
		}
		offset += uint(reclen) - 4
		if recseg == 0 || recseg == 512 {
			recnum := uint(record[1])
			matrix[recnum] += 1
			if !typeFilter {
				smfParse(record[:offset])
			} else if typeTable[recnum] {
				smfParse(record[:offset])
			}
			count++
			offset = 0
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("SMF Records read: %d\n", count)
	fmt.Println("┌──────────────────┐")
	fmt.Println("│ Type │   Total   │")
	fmt.Println("└──────────────────┘")
	for i := 0; i < len(matrix); i++ {
		if matrix[i] > 0 {
			fmt.Printf("  %4d   %10d\n", i, matrix[i])
		}
	}
	fmt.Printf("Elapsed time: %v seconds\n", elapsed.Seconds())
	_, err = db.Exec("VACUUM;")
	handleErr(err)
}

func printUsage() {
	fmt.Println("smf2db [options] smfFile [dbFile]")
	fmt.Println("  smfFile: input SMF file with RDW information")
	fmt.Println("   dbFile: optional output SQLite database name (default: smf.db)")
	fmt.Println()
	fmt.Println("  Options:")
	flag.PrintDefaults()
}
