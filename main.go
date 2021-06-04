package main

import (
	"bufio"
	"database/sql"
	"encoding/binary"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"time"
)

var db *sql.DB
var err error

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	_ = os.Remove("./smf.db")
	db, err = sql.Open("sqlite3", "./smf.db")
	handleErr(err)
	_, err = db.Exec("PRAGMA cache_size = 400000;")
	handleErr(err)
	_, err = db.Exec("PRAGMA synchronous = OFF;")
	handleErr(err)
	_, err = db.Exec("PRAGMA journal_mode = OFF;")
	handleErr(err)
	_, err = db.Exec("PRAGMA locking_mode = EXCLUSIVE;")
	handleErr(err)
	_, err = db.Exec("PRAGMA count_changes = OFF;")
	handleErr(err)
	_, err = db.Exec("PRAGMA temp_store = MEMORY;")
	handleErr(err)
	_, err = db.Exec("PRAGMA auto_vacuum = NONE ;")
	handleErr(err)
	//smfInit()
}

func main() {
	file, err := os.Open("testdata/smftest.rdw")
	//file, err := os.Open("testdata/safrasmf")
	handleErr(err)
	defer func() {
		err = file.Close()
		err = db.Close()
	}()
	rdr := bufio.NewReader(file)
	var record = make([]byte, 64*1024)
	var count int = 0
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
		offset += (uint(reclen) - 4)
		if recseg == 0 || recseg == 512 {
			recnum := uint(record[1])
			matrix[recnum] += 1
			smfParse(record[:offset])
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
