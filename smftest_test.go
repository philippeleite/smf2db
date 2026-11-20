package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const (
	smfTestFile    = "data/smftest.rdw"
	smfTestReport  = "data/smftest.txt" // IFASMFDP SUMMARY ACTIVITY REPORT of the job that wrote smftest.rdw
	smfTestMissing = "test data not found (data/ is git-ignored and kept locally)"
)

// TestSmfTestFile converts data/smftest.rdw and checks the per-type record counts printed by smf2db,
// and the row counts of every parsed type's main table, against the RECORDS WRITTEN column of the
// IFASMFDP report in data/smftest.txt.
func TestSmfTestFile(t *testing.T) {
	for _, f := range []string{smfTestFile, smfTestReport} {
		if _, err := os.Stat(f); err != nil {
			t.Skip(smfTestMissing)
		}
	}
	requireSQLite(t)

	expected, expectedTotal := readIfasmfdpReport(t, smfTestReport)
	dbFile := filepath.Join(t.TempDir(), "smftest.db")
	actual, actualTotal := runMain(t, smfTestFile, dbFile)

	if actualTotal != expectedTotal {
		t.Errorf("records read: got %d, want %d", actualTotal, expectedTotal)
	}
	for typ, want := range expected {
		if got := actual[typ]; got != want {
			t.Errorf("type %d record count: got %d, want %d", typ, got, want)
		}
	}
	for typ, got := range actual {
		if _, ok := expected[typ]; !ok {
			t.Errorf("type %d: read %d records, not in the IFASMFDP report", typ, got)
		}
	}

	conn, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for typ := range parseTable {
		want := expected[typ]
		if want == 0 {
			continue // no records of this type in the file, so its tables are never created
		}
		var got uint64
		if err := conn.QueryRow(fmt.Sprintf("SELECT count(*) FROM smf%03d", typ)).Scan(&got); err != nil {
			t.Errorf("type %d table: %v", typ, err)
			continue
		}
		if got != want {
			t.Errorf("table smf%03d rows: got %d, want %d", typ, got, want)
		}
	}
}

// requireSQLite skips the test when go-sqlite3 was built without CGO: main would otherwise
// hit log.Fatal on the first statement and kill the test binary.
func requireSQLite(t *testing.T) {
	t.Helper()
	conn, err := sql.Open("sqlite3", ":memory:")
	if err == nil {
		err = conn.Ping()
		conn.Close()
	}
	if err != nil {
		t.Skipf("SQLite unavailable, build with CGO_ENABLED=1: %v", err)
	}
}

// readIfasmfdpReport returns RECORDS WRITTEN per type (types with 0 written are left out) and the TOTAL.
func readIfasmfdpReport(t *testing.T, path string) (map[int]uint64, uint64) {
	t.Helper()
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	counts := map[int]uint64{}
	var total uint64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		fields := strings.Fields(sc.Text())
		if len(fields) > 0 && fields[0] == "0" { // ASA carriage control (double space)
			fields = fields[1:]
		}
		// TYPE READ PERCENT % AVG MIN MAX WRITTEN
		if len(fields) != 8 || fields[3] != "%" {
			continue
		}
		written, err := strconv.ParseUint(strings.ReplaceAll(fields[7], ",", ""), 10, 64)
		if err != nil {
			continue // header line
		}
		if fields[0] == "TOTAL" {
			total = written
		} else if typ, err := strconv.Atoi(fields[0]); err == nil && written > 0 {
			counts[typ] = written
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}
	if total == 0 || len(counts) == 0 {
		t.Fatalf("%s: no SUMMARY ACTIVITY REPORT counts found", path)
	}
	return counts, total
}

// runMain runs main() on smfFile and parses its stdout summary into per-type counts and the total read.
func runMain(t *testing.T, smfFile, dbFile string) (map[int]uint64, uint64) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	out := make(chan string)
	go func() {
		b, _ := io.ReadAll(r)
		out <- string(b)
	}()

	stdout, args := os.Stdout, os.Args
	os.Stdout, os.Args = w, []string{"smf2db", smfFile, dbFile}
	defer func() { os.Stdout, os.Args = stdout, args }()
	main()
	w.Close()
	os.Stdout = stdout
	text := <-out

	counts := map[int]uint64{}
	var total uint64
	for _, line := range strings.Split(text, "\n") {
		if n, err := fmt.Sscanf(line, "SMF Records read: %d", &total); n == 1 && err == nil {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			continue
		}
		typ, err1 := strconv.Atoi(fields[0])
		cnt, err2 := strconv.ParseUint(fields[1], 10, 64)
		if err1 == nil && err2 == nil {
			counts[typ] = cnt
		}
	}
	if total == 0 {
		t.Fatalf("no summary in smf2db output:\n%s", text)
	}
	return counts, total
}
