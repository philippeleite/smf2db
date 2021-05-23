package main

import (
	"bufio"
	_ "database/sql"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"
	// _ "github.com/mattn/go-sqlite3"
)

func main() {
	file, err := os.Open("smftest.rdw")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(8)
	}
	defer func() {
		err = file.Close()
	}()
	rdr := bufio.NewReader(file)
	record := make([]byte, 64*1024)
	var count int = 0
	var offset int = 0
	matrix := make([]uint64, 256)
	var buffer = make([]byte, 2)
	var reclen, recseg uint16

	start := time.Now()
	for {
		_, err = io.ReadFull(rdr, buffer)
		reclen = binary.BigEndian.Uint16(buffer)
		// spdlog::debug("Record length: %d (0x%04x)", reclen, reclen)
		if err == io.EOF {
			break
		}
		_, err = io.ReadFull(rdr, buffer)
		recseg = binary.BigEndian.Uint16(buffer)
		// spdlog::debug("Record segment: {0:d} (0x{0:04x})", recseg)
		// Check if segment descriptor is valid -> recseg in (0, 256, 512, 768)
		if recseg&0xfcff != 0 {
			fmt.Println("Invalid file format")
			fmt.Println("SMF file must be transferred preserving RDW information")
			os.Exit(8)
		}
		bufrec := make([]byte, reclen-4)
		_, err = io.ReadFull(rdr, bufrec)
		copy(record[offset:], bufrec)
		offset += int(reclen - 4)
		if recseg == 0 || recseg == 512 {
			recnum := int(record[1])
			matrix[recnum] += 1
			// smfparse(record)
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
			fmt.Printf("  %4v   %10v\n", i, matrix[i])
		}
	}
	fmt.Printf("Elapsed time: %v seconds\n", elapsed.Seconds())
}
