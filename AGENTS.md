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
- `fields_test.go`: `fields` decoding (signed, packed time/interval/integer, hex digits, short section) and IBM hex floating point.
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
| 60 (VVDS update) | `smf060.go` | `smf060` |
| 61 (catalog define) | `smfcatalog.go` | `smf061` |
| 62 (VSAM component or cluster opened) | `smf062.go` | `smf062`, `smf062_vol` |
| 64 (VSAM component or cluster status) | `smf064.go` | `smf064`, `smf064_ext` |
| 65 (catalog delete) | `smfcatalog.go` | `smf065` |
| 66 (catalog alter) | `smfcatalog.go` | `smf066` |
| 70 (RMF processor activity) | `smf070.go`, `rmf.go` | `smf070` (header + RMF product section) and per subtype: 1 `smf070_cpc`, `smf070_cpu`, `smf070_asid`, `smf070_lpar`, `smf070_lcpu`, `smf070_cpuid`; 2 `smf070_icsf` |
| 71 (RMF paging activity) | `smf071.go`, `rmf.go` | `smf071` (header + RMF product section), `smf071_paging` |
| 72 (RMF workload activity, storage, serialization delay) | `smf072.go`, `rmf.go` | `smf072` (header + RMF product section) and per subtype: 3 `smf072_wlm`, `smf072_served`, `smf072_rg`, `smf072_scp`, `smf072_rtd`, `smf072_wrm`, `smf072_dn`; 4 `smf072_stor`, `smf072_swap`; 5 `smf072_ser`, `smf072_cms`, `smf072_local`, `smf072_cmlo`, `smf072_cmlr`, `smf072_latch`, `smf072_enq`, `smf072_qscan` |
| 73 (RMF channel path activity) | `smf073.go`, `rmf.go` | `smf073` (header + RMF product section), `smf073_ctl`, `smf073_chp`, `smf073_ext` |
| 74 (RMF activity of several resources) | `smf074.go`, `rmf.go` | `smf074` (header + RMF product section) and per subtype: 1 `smf074_devctl`, `smf074_dev`; 2 `smf074_xcfctl`, `smf074_xcfsys`, `smf074_xcfmbr`; 3 `smf074_omvs`; 5 `smf074_cachectl`, `smf074_cachedev`, `smf074_cachestat`, `smf074_raid`; 6 `smf074_hfsglobal`, `smf074_hfsbuf` |
| 75 (RMF page data set activity) | `smf075.go`, `rmf.go` | `smf075` (header + RMF product section), `smf075_page` (only `data/smfdata` has type 75) |
| 78 (RMF virtual storage and I/O queuing) | `smf078.go`, `rmf.go` | `smf078` (header + RMF product section) and per subtype: 2 `smf078_common`; 3 `smf078_cfgctl`, `smf078_cfg`, `smf078_ioq`, `smf078_ioqglob`, `smf078_iop`, `smf078_hpav` |
| 80 (RACF processing) | `smf080.go` | `smf080`, `smf080_reloc`, `smf080_reloc2`, and reference tables `smf080_events`, `smf080_qualifiers`, `smf080_types` |
| 82 (ICSF) | `smf082.go` | `smf082` (header, every subtype), `smf082_t20` (subtype 20) |
| 84 (JES monitoring facility) | `smf084.go` | `smf084` (header, JES product and general sections) and, for subtype 21, `smf084_mem`, `smf084_res` (only `data/smfdata` has type 84) |
| 92 (z/OS UNIX file system activity) | `smf092.go` | `smf092` (header, subsystem and identification sections) and `smf092_open` (10), `smf092_close` (11, 16), `smf092_delete` (14), `smf092_count` (17) |

Type 6 has five writers (External Writer, JES2, JES3, PSF, IP PrintWay) that share the header up to `SMF6SBS`; their sections follow in `SMF6PAD1` bit order, each with its own 2-byte length. Only IP PrintWay sections are decoded, since it is the only writer in the sample data. The extended-mode fields (`SMF6BYTD`, `SMF6URI`, `SMF6ACCT`) are implemented from the manual but untested: every sample record is basic mode.

Type 60 stores the raw VVR/NVR in `SMF60CRC` (BLOB, length in its first two bytes) and, in `SMF60RQT`, the 8-byte VVDS request time in microseconds that version 2 records append after it. The manual gives that field no name; `SMF60RQT` is ours. Types 61, 65 and 66 share one layout and one implementation (`smfcatalog.go`: `smfCatalogColumns` and `smfCatalogRow`): the raw catalog record in `SMFnnCRC`, optionally followed by a 16-byte STCKE timestamp (`SMFnnETS`, converted with `stcketime`). The sample records (version 1) do not have it, so it is NULL there and only unit-tested. Type 65 adds `SMF65FNC` (`S` = data set scratched, `U` = only catalog entries changed); type 66 adds `SMF66FNC` (`R` = renamed) and the new name `SMF66NNM` (no sample record is a rename).

Type 62: after the 10-byte volume entries come the 24-byte SMS classes (always present in the samples, zeros when `SMF62IND` bit 4 is off; stored only when it is on) and the statistics section. The 66-byte encryption block between them is assumed to exist only when bit 5 (encrypted) is on, which no sample record has (see the `shortcut:` comment). Sample statistics sections are 12 bytes, so `SMF62DET`, `SMF62DKL`, `SMF62JobID` and `SMF62SysplexName` are NULL there.

Type 64: 26-byte extent entries follow `SMF64ESL`, then the statistics section with its own length (`SMF64SLN`); bytes after that length are ignored (records with `SMF64RIN` X'01' have no extents but end with 26 zero bytes). The "change since OPEN" fields (`SMF64DIL`…`SMF64DEP`, `SMF64DAU`) are signed; `SMF64CHR` is documented as EBCDIC but is a binary RBA/CI.

Type 70 notes:
- Sections with many fields are declared as `fields` tables (`fields.go`) generated from the manual with `tools/genfields.py`, then reviewed by hand (`SMF70HOF` signed, `SMF70GJT`/`SMF70_IPL_TIME` TOD, `SMF70SER` hex digits). Three offsets in the manual are inconsistent (`SMF70I01`, `SMF70P07`, `SMF70P11`); the generator picks the one that continues the sequence and the table keeps a comment.
- `rmf.go` has the RMF product section, common to types 70–79 (`rmfProductFields(prefix)`), `rmfColumns`/`rmfRow` for the main table of any RMF type (SMF header + product section), and `rmfSection` to map triplets to child tables; new RMF types should use them, as `smf071.go`, `smf072.go` and `smf074.go` do. The reassembly area of broken large records is not decoded; `SMFnnRAN` shows whether a record is one.
- Logical processors are linked to their partition by `SMF70LPN` (each partition is followed by its `SMF70BDN` processors), and to the CPU type name by `SMF70CIX`, which is also stored, 1-based, in `smf070_cpuid`.
- Durations are stored as recorded: `SMF70WAT`/`SMF70PAT` in TOD units (bit 51 = 1 µs, divide by 4096 for µs), `SMF70PDT`/`SMF70EDT`/`SMF70ONT` in microseconds. ICSF counters (`R702*`) are IBM hex floating point, stored as REAL.
- Type 71 has one paging data section per record (`smf071_paging`); only `data/smfdata` contains type 71 records. Its HFP averages (`SMF71LVS`) match the integer ones (`SMF71AVS`), which confirms the HFP conversion on real data. The swap placement section is not decoded (no valid values since z/OS V1R13).
- Not decoded (absent from the samples): subtype 1 logical core and tenant resource group sections; subtype 2 CCA, accelerator and PKCS11 sections.

Type 72 notes:
- `smf072Sections` maps each triplet to a table; every child row has `ID` and `SEQ` (1-based position of the section within its triplet). Triplets the manual maps with one structure share a table and are told apart by a field: `smf072_cms` by `R725CMTY` (4 triplets), `smf072_latch` by `R725LATY` (2), `smf072_enq` by `R725ENSC` (3).
- Response time distribution: the first `smf072_rtd` section of a record (`SEQ` 1) is the bucket map (50…400% and X'FFFFFFFF'), not counts. A period's 14 counters (`R723RTD01`…`R723RTD14`, unnamed in the manual) are at `SEQ = R723CRTX + 1`; their sum matches `R723CRCP` (completed transactions) for 456 of 460 non-heterogeneous periods in the samples. Heterogeneous report class periods (`R723CRS1` bit 0) have invalid distributions, as the manual says.
- IBM names with `#` (`R723MTV#`, `R723SCS#`, `R724PER#`) are double-quoted by `fields.columns()`; queries must quote them too.
- 16-byte sums of squares (`R725*SQ`) use the `w` kind (REAL, approximate).

Type 73: the 48-byte channel measurement data (`SMF73CCM`) has one layout per measurement group (`SMF73CMG` 1, 2, 3); `smf073_chp` has the columns of all three and fills only the channel's group. Fields documented as "float" are IBM hex floating point (4 or 8 bytes): a HiperSockets data unit size decodes to 4096 as HFP and 576 as IEEE. Every record has a section for each possible CHPID; `SMF73FG3` X'02' marks real entries.

Type 74: sections are mapped with `rmfSection` (`rmf.go`: `rmfCreateSections`/`rmfInsertSections`, also used by type 72). The RAID rank/extent pool section has two layouts at the same offsets, chosen by `R7451FLG` (1 rank, 2 extent pool); `smf074_raid` has both column sets and fills one, like type 73's measurement groups (all sample rows are extent pool data). Not decoded (no sample data): subtypes 4 and 7–10, XCF path sections, cache device extension sections, HFS file system sections. `tools/genfields.py` also accepts the `l_float`/`s_float` formats (format names are matched case-insensitively) and strips footnote markers such as `(2)` from names.

Type 78: in subtype 3 the configuration control sections and the I/O queuing global section each carry their own triplet (at offset 4, relative to the start of the section) to the configuration data sections (`smf078_cfg`) and IOP initiative queue sections (`smf078_iop`); those rows have `ID`, `SEQ` of the enclosing section and `SUBSEQ`. Not decoded (no sample data): subtype 2 private area and subpool sections. Types 76 and 77 are skipped: the samples have no records of them.

Type 84: not an RMF record; its header has `SMF84SBS`, `SMF84SGN`, `SMF84FL1` and `SMF84VER` before the subtype (at manual offset 24), and the header time is binary although the manual says packed. The JES product section times are packed `hhmmsstF` (`H` kind in `fields.go`). Subtype 21's JES2 resource usage section carries two triplets, relative to its start, for the memory area and resource entries. Not decoded (no sample data): subtypes 1–10 (JES3).

Type 80: the layout is in the RACF manual (`data/icha300_v3r1.pdf`, SA23-2288: "Format of SMF type 80 records" and the tables of event codes and relocate section data), not in the SMF manual. Relocate sections (1-byte type/length) and extended-length ones (2-byte) are located by `SMF80REL`/`SMF80RL2`, offsets from `SMF80FLG` (= `rec[0]`). Data of types listed as EBCDIC is stored as text, everything else as hex digits (data type 6, command-related data, is not decoded). `smf080Events`, `smf080Qualifiers` and `smf080Types` were extracted from the manual's tables by script and are loaded into the reference tables for joins.

Type 82: the layout is in the ICSF System Programmer's Guide (`data/csfb200_icsf_spg_hcr77e0.pdf`, SC14-7507, appendix "ICSF SMF records"). Subtypes below 40 have their data right after the 20-byte header; 40 and above start with an ICSF header (not decoded: no sample data). Subtype 20's TOD clocks use the `M` kind (microseconds since 1900), so `SMF82t20_DQAP - SMF82t20_NQAP` is the coprocessor time in microseconds.

Type 92: the three header triplets (subsystem, identification, data) are read at their fixed offsets; subtype 17 has no identification section. Open times (`SMF92OTO`, `SMF92CTO`) are the raw TOD clock and include leap seconds (26 s ahead of the record time in the samples), while `SMF92CTC` is UTC, so `SMF92CTC` is usually earlier than `SMF92CTO`; values are stored as recorded and the quirk is documented in `smf092.go`. Not decoded (no sample data): subtypes 1–9, 12, 13, 15, 50–59.

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
- **Header times**: the `T` kind decodes SMF binary times (hundredths of a second since midnight), so header fields can go in a `fields` table. The `M` kind stores a TOD clock as integer microseconds since 1900, for fields whose differences matter more than their date.
- **Newer helpers** (used by `smf021.go`, `smf042*.go`; prefer them for new parsers): `createTable(name, columns)` creates the table and prepares an INSERT with one placeholder per column (no hand-counted `?` lists); `insert(stmt, args...)` executes and logs; `section` and `tripletAt` (in `smf042.go`) read fields with bounds checks, returning NULL for a field past the end of a shorter section, and locate sections by triplet; `todtime`/`stcketime` convert TOD/STCKE clock values. For sections with dozens of fields, declare a `fields` table (`fields.go`: name, offset, kind, length) and use `.columns()` / `.values(section)`; generate it from the manual with `tools/genfields.py` (see its docstring) and review the kinds by hand.
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

Record layouts: IBM *z/OS 3.1 MVS System Management Facilities (SMF)*, SA38-0667-60, kept locally at `data/ieag200_v3r1.pdf` (git-ignored). This is the source of truth for every parser offset, except type 80, whose layout is in *z/OS 3.1 Security Server RACF Macros and Interfaces*, SA23-2288-60, at `data/icha300_v3r1.pdf`, and type 82, in *z/OS 3.1 Cryptographic Services ICSF System Programmer's Guide*, SC14-7507-60, at `data/csfb200_icsf_spg_hcr77e0.pdf`.

- Search it as text: `pdftotext -layout data/ieag200_v3r1.pdf <out>.txt`, then grep for the chapter heading, e.g. `Record type 41 (X'29') --`.
- Read field tables with `-raw`, not `-layout`: `-layout` interleaves the Length/Format column with other rows. Find the page with `pdftotext -f N -l N` + grep for a field name, then `pdftotext -f N -l M -raw data/ieag200_v3r1.pdf - | grep -E "^[0-9]+ [0-9A-F]+ SMF21"` prints clean `offset hex name length format` rows.
- Field tables list `Offsets (dec hex) | Name | Length | Format | Description`. The manual's offsets include the 4-byte RDW (`SMFxLEN` at 0, `SMFxSEG` at 2), so subtract 4 to get the index into `rec`.
