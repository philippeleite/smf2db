# AGENTS.md — smf2db

Guide for AI agents (and humans) working in this repository.

## What it is

`smf2db` is a Go command-line utility that reads an **IBM z/OS SMF (System Management Facilities) file**, transferred from the mainframe in binary **with the RDW preserved**, and writes the decoded records to a **SQLite** database.

- Version: `0.1.0` (`version` variable in `main.go`)
- Module: `smf2db`, `go 1.16`, a single `main` package (all `.go` files live in the root)
- Only dependency: `github.com/mattn/go-sqlite3 v1.14.7`, which **requires CGO**

## Build, run and test

```sh
# go-sqlite3 needs CGO + a C compiler (gcc/mingw on Windows).
# This environment defaults to CGO_ENABLED=0: the binary builds, but fails at runtime when opening the database.
CGO_ENABLED=1 go build -o smf2db .

# Without gcc, zig works as the C compiler:
CGO_ENABLED=1 CC="zig cc -target x86_64-windows-gnu" go build -o smf2db .

go test ./...      # unit tests; the SMF file test is skipped without CGO
CGO_ENABLED=1 CC="zig cc -target x86_64-windows-gnu" go test ./...   # full run, ~7 s
go vet ./...
```

Tests:
- `zutil_test.go`: date/time/TOD/EBCDIC helpers.
- `main_test.go`: `smfParse` skips a record whose parser panics.
- `smf006_test.go`: `SMF6ACCT` decoding, using the example from the manual.
- `smf042_test.go`: `section` bounds handling and the `smf042Chain` backward-pointer guard.
- `smftest_test.go` (`TestSmfTestFile`): runs `main()` on `data/smftest.rdw` (363 MB, 815,993 records, 2017) and checks the per-type counts printed by smf2db, and the row count of every parsed type's main table (`smfNNN`), against the RECORDS WRITTEN column of the IFASMFDP report in `data/smftest.txt`. It also fails if smf2db writes anything to the log (skipped record, failed `Prepare`/`Exec`), which catches child tables left empty by a broken statement. It skips when those files are missing or SQLite has no CGO. It only checks counts, not field values, and only for types registered in `parseTable`.
- `data/smfdata` (26 MB, 2019) is another local sample, with no reference report.

Usage:

```
smf2db [options] smfFile [dbFile]
  smfFile  input SMF file, with RDW
  dbFile   output SQLite database (default: smf.db)
  -t       record types to process: "4,5,30" or ranges "14:18" (surrounding ()[]{} accepted)
  -v       show version (does not exit)
  -h       help
```

- **`dbFile` is deleted (`os.Remove`) on every run.** Never point it at a database that must be kept.
- At the end it prints the record count per type (every type read, including those without a parser) and the elapsed time, then runs `VACUUM`.
- Git-ignored: `/smf2db`, `/data/`, `/smftest.txt`, `/smf.db`, `/build/`, `/docs/`, `.idea/`, `.claude/`, `CLAUDE.md`. SMF test data is kept outside the repository.

## Architecture

| File | Role |
|---|---|
| `main.go` | CLI, SQLite setup, RDW file reading, segment reassembly, dispatch by type |
| `zutil.go` | Conversions: `smfdate`, `filedate`, `smftime`, `bcd2asc`, `asc2bcd`, `b3toInt`, `divmod` |
| `zutil_test.go` | Unit tests for the helpers |
| `smfNNN.go` | One parser per SMF record type |

### Flow in `main.go`

1. For each block, read the RDW: 2-byte length + 2-byte segment descriptor (big-endian).
2. Valid descriptors: `0x0000` (complete), `0x0100` (first), `0x0200` (last), `0x0300` (middle). Any other value (`recseg&0xfcff != 0`) aborts with "SMF file must be transferred preserving RDW information".
3. Segments are concatenated into a 64 KiB buffer (`record`). When the record ends (`0x0000` or `0x0200`), `record[:offset]` is dispatched.
4. `smfParse` reads the type from `rec[1]` and calls the function registered in `parseTable`.

**Important:** the `rec` slice passed to a parser **does not include the RDW**. Offsets are therefore 4 bytes lower than in IBM's documentation: `rec[0]` = flag, `rec[1]` = type, `rec[2:6]` = time, `rec[6:10]` = date, `rec[10:14]` = SID. Section offsets read from the record itself (offset/length/number triplets) are relative to the RDW, so they are used as `off - 4`.

### Supported types and generated tables

| Type | File | Tables |
|---|---|---|
| 4 (step termination) | `smf004.go` | `smf004`, `smf004_dev` |
| 5 (job termination) | `smf005.go` | `smf005` |
| 6 (output writer) | `smf006.go` | `smf006` (common header, every writer), `smf006_ipw` (IP PrintWay sections only, `SMF6SBS` = 9) |
| 14 (dataset INPUT/RDBACK) | `smf014.go` | `smf014`, `smf014_ucb` |
| 15 (dataset OUTPUT/UPDAT/…) | `smf015.go` | `smf015`, `smf015_ucb` |
| 17 (scratch) | `smf017.go` | `smf017`, `smf017_vol` |
| 18 (rename) | `smf018.go` | `smf018`, `smf018_vol` |
| 21 (tape error statistics by volume) | `smf021.go` | `smf021` |
| 30 (common address space work) | `smf030.go` | `smf030` + `smf030_{ioac,comp,proc,excp,acct,stor,perf,oper,appc,appc_cum,unix,arm,usag,mers,cntr,zedc}` |
| 41 (DIV / VLF) | `smf041.go` | `smf041`, `smf041_{prod,obja,obju,obji,vlfs}` |
| 42 (DFSMS statistics) | `smf042.go`, `smf042misc.go`, `smf042io.go` | `smf042` (header + product section, every subtype) and per subtype: 1 `smf042_bmf`, `smf042_bmf_sc`; 2 `smf042_cu`, `smf042_cu_vol`; 3 `smf042_ea`; 4 `smf042_ext`; 5 `smf042_sc`, `smf042_vol`, `smf042_vol_{vtoc,vtix,vvds}`; 6 `smf042_job`, `smf042_ds`; 9 `smf042_x37`; 21 `smf042_mdel`, `smf042_mdel_alias`; 24 `smf042_madd`, `smf042_madd_alias`; 27 `smf042_vtoc` |

Type 6 has five writers (External Writer, JES2, JES3, PSF, IP PrintWay) that share the header up to `SMF6SBS`; their sections follow in `SMF6PAD1` bit order, each with its own 2-byte length. Only IP PrintWay sections are decoded, since it is the only writer in the sample data. The extended-mode fields (`SMF6BYTD`, `SMF6URI`, `SMF6ACCT`) are implemented from the manual but untested: every sample record is basic mode.

Type 42 notes:
- The subtype triplets sit at fixed offsets after the product triplet. `SMF42NT` (number of triplets) is 0 in most real records, so it is ignored.
- Subtype 5 volumes and subtype 6 data sets are linked lists (`S42VTNXT`, `S42DSNXT`); `smf042Chain` walks them and stops if a pointer goes backwards. Child rows link by `ID` (and `S42VTSER` for the volume component tables).
- Not decoded: subtype 4 concurrent copy / VCC sessions, subtype 5 volume metrics / system I/O / background sections, subtype 6 synchronous I/O sections, and every subtype not listed above (their records still get the `smf042` row).
- Sample data quirks: `SMF42LYY`/`SMF42CYY` hold the binary year since 1900 and `SMF42LDD`/`SMF42CDD` the binary day of year (the manual says `0cyy`/`dddF`); `SMF42LUI`/`SMF42PUI` are RACF tokens (ICHRUTKN), stored as BLOB though the manual says EBCDIC; TOD fields (`SMF42PTS`, `S42EXSTM`, …) are UTC, while SMF header times are local.

`smf015.go` is a copy of `smf014.go` with renamed identifiers (same layout). **A fix in one almost always applies to the other.**

## Parser pattern (follow it when adding a new type)

Every `smfNNN.go` follows the same template. Copy the closest file (`smf017.go` is the simplest; `smf041.go`/`smf030.go` for records with triplets):

```go
var p1StmtNNN, p2StmtNNN *sql.Stmt   // one prepared statement per table
var onceNNN sync.Once
var countNNN int                     // produces the ID column

func smfNNNInit() {   // DROP TABLE IF EXISTS + CREATE TABLE + db.Prepare(INSERT ... VALUES (?,...))
}

func smfNNNParse(rec []byte) {
    onceNNN.Do(smfNNNInit)           // tables are created only if the type appears in the file
    bcd2asc(rec[a:b])                // convert text fields IN PLACE, before string()
    countNNN++
    p1StmtNNN.Exec(countNNN, ...)    // main record
    for ... { p2StmtNNN.Exec(countNNN, ...) }  // repeating sections -> child table with the same ID
}
```

Then register the function in `parseTable` in `main.go`.

Conventions:
- **ID**: a sequential counter per type; child tables (`_dev`, `_vol`, `_ucb`, …) reuse the parent record's `ID` for joins. No foreign keys or indexes.
- **Column names = IBM field names** (`SMF30JBN`, `JFCBDSNM`, …). Each `Exec` argument carries a comment with the field name. Keep it that way: it is the only documentation of the mapping.
- **Binary integers**: `binary.BigEndian.Uint16/32/64`, 3-byte fields with `b3toInt`, 1-byte fields with `int(rec[n])`.
- **Dates/times** become text `YYYY-MM-DD HH:MM:SS.mmm`: `smfdate` (packed `0cyydddF`) returns the date only, `smftime` (hundredths of a second since midnight) returns the time only with date `0000-00-00`. A zero date/time becomes `0000-00-00 00:00:00.000`. `filedate` handles the 3-byte JFCB date.
- **Newer helpers** (used by `smf021.go`, `smf042*.go`; prefer them for new parsers): `createTable(name, columns)` creates the table and prepares an INSERT with one placeholder per column (no hand-counted `?` lists); `insert(stmt, args...)` executes and logs; `section` (in `smf042.go`) reads fields with bounds checks, returning NULL for a field past the end of a shorter section; `todtime` converts 8-byte TOD clock values.
- **Text**: `ebcdic(b)` converts a field and returns it as a string (newer parsers use it). `bcd2asc` converts **EBCDIC → ASCII** (the name is misleading: it is not BCD) by modifying the buffer. Strings keep their space padding; nothing is trimmed.
- **Optional fields** (sections that may be absent or depend on flags): an `interface{}` variable holding `nil`, stored as NULL; the column has no `NOT NULL`. See `smf014.go` (DASD vs. tape, extended sections) and `smf030.go` (`SMF30CDO`/`SMF30USO` only when the header is long enough).
- **Errors**: SQL failures are only logged (`log.Println`/`log.Printf`) and processing continues; only open/file errors stop the program (`handleErr` → `log.Fatalln`). Parsers share a global `err` variable.
- **Language**: documentation and code comments are written in English.

## Performance

SQLite is opened with `_mutex=no`, `_synchronous=off`, `_journal_mode=off`, `_locking_mode=exclusive`, `_temp_store=memory`, a large cache, plus `PRAGMA count_changes=OFF`. There is no explicit transaction; each `Exec` autocommits (cheap because the journal is off). Processing is sequential and single-threaded: parsers rely on global state and are not safe for concurrency.

## Known issues (verified in the code)

- Parsers do not check bounds themselves. `smfParse` recovers from their panics, logs the record type and length, and skips the record; rows already inserted for that record before the panic stay in the database (partial data). The reader loop in `main.go` is not covered: an RDW length below 4 or a spanned record larger than 64 KiB still panics.
- `-v` prints the version but does not exit; with no arguments it falls through to usage and exits with code 8.
- `smf030.go` has `IsFirst/IsAdditional/More/IsLast` (subtype 30 records split across several records), but the calls are commented out: split records are stored as independent records.
- `main.go` only prints (does not abort on) read errors in the middle of a record; a truncated file can produce garbage in the last record.
- No test checks parsed field values against the manual; the sample SMF files live only in the git-ignored `data/`.

## Reference

Record layouts: IBM *z/OS 3.1 MVS System Management Facilities (SMF)*, SA38-0667-60, kept locally at `data/ieag200_v3r1.pdf` (git-ignored). This is the source of truth for every parser offset.

- Search it as text: `pdftotext -layout data/ieag200_v3r1.pdf <out>.txt`, then grep for the chapter heading, e.g. `Record type 41 (X'29') --`.
- Read field tables with `-raw`, not `-layout`: `-layout` interleaves the Length/Format column with other rows. Find the page with `pdftotext -f N -l N` + grep for a field name, then `pdftotext -f N -l M -raw data/ieag200_v3r1.pdf - | grep -E "^[0-9]+ [0-9A-F]+ SMF21"` prints clean `offset hex name length format` rows.
- Field tables list `Offsets (dec hex) | Name | Length | Format | Description`. The manual's offsets include the 4-byte RDW (`SMFxLEN` at 0, `SMFxSEG` at 2), so subtract 4 to get the index into `rec`.
