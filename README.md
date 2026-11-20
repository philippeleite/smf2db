# smf2db

**smf2db** converts IBM z/OS **SMF (System Management Facilities)** data into a **SQLite** database, so you can analyze mainframe activity with plain SQL on any workstation, with no mainframe tools, SAS/MXG or licensed reporting products.

It reads an SMF dump file transferred from z/OS in binary (with its record descriptor words), decodes the supported record types field by field, and writes one table per record (plus child tables for repeating sections), using the IBM field names as column names.

```text
$ smf2db -t 4,5,30,42 SMF.DUMP smf.db
SMF Records read: 815993
┌──────────────────┐
│ Type │   Total   │
└──────────────────┘
     4         9339
     5         2357
    30        38168
    42       102264
   ...
Elapsed time: 9.29 seconds
```

## Features

- **Fast**: a single pass over the file with prepared statements and SQLite tuned for bulk loading (no journal, exclusive lock). About 800,000 records (363 MB) load in under 10 seconds.
- **IBM field names**: every column is named after the field in the IBM manual (`SMF30JBN`, `S42DSION`, `JFCBDSNM`, …), so the manual is your data dictionary.
- **Repeating sections become child tables**, linked to their parent row by the `ID` column (devices, volumes, DD statements, data sets, …).
- **Optional sections and fields** that a record does not contain are stored as `NULL`.
- **Tolerant**: a malformed record is logged and skipped instead of stopping the whole conversion.
- **No server**: the result is a single SQLite file you can open with the `sqlite3` CLI, DB Browser for SQLite, DBeaver, Python, R, Excel (ODBC), Grafana and so on.

## Supported record types

| Type | Description | Tables |
|---|---|---|
| 4 | Step termination | `smf004`, `smf004_dev` |
| 5 | Job termination | `smf005` |
| 6 | Output writer (all writers' common header; IP PrintWay sections) | `smf006`, `smf006_ipw` |
| 14 | Data set activity: INPUT or RDBACK close | `smf014`, `smf014_ucb` |
| 15 | Data set activity: OUTPUT, UPDAT, INOUT or OUTIN close | `smf015`, `smf015_ucb` |
| 17 | Scratch (delete) non-VSAM data set | `smf017`, `smf017_vol` |
| 18 | Rename non-VSAM data set | `smf018`, `smf018_vol` |
| 21 | Tape error statistics by volume | `smf021` |
| 30 | Common address space work (job/step accounting) | `smf030` and `smf030_{ioac,comp,proc,excp,acct,stor,perf,oper,appc,appc_cum,unix,arm,usag,mers,cntr,zedc}` |
| 41 | Data-in-virtual objects and VLF statistics | `smf041` and `smf041_{prod,obja,obju,obji,vlfs}` |
| 42 | DFSMS statistics and configuration (subtypes 1–6, 9, 21, 24, 27) | `smf042` and `smf042_*` (see below) |
| 60 | VSAM volume data set (VVDS) updated | `smf060` |
| 61 | Integrated catalog facility define activity | `smf061` |
| 62 | VSAM component or cluster opened | `smf062`, `smf062_vol` |
| 64 | VSAM component or cluster status (CLOSE, EOV, out of space) | `smf064`, `smf064_ext` |
| 65 | Integrated catalog facility delete activity | `smf065` |
| 66 | Integrated catalog facility alter activity | `smf066` |
| 70 | RMF processor activity: CPU, PR/SM partitions, ICSF (subtypes 1 and 2) | `smf070` and `smf070_{cpc,cpu,asid,lpar,lcpu,cpuid,icsf}` |
| 71 | RMF paging activity | `smf071`, `smf071_paging` |
| 72 | RMF workload activity (3), storage data (4) and serialization delay (5) | `smf072` and `smf072_*` (see AGENTS.md) |
| 73 | RMF channel path activity | `smf073`, `smf073_ctl`, `smf073_chp`, `smf073_ext` |
| 74 | RMF device (1), XCF (2), OMVS kernel (3), cache subsystem (5) and HFS (6) activity | `smf074` and `smf074_{devctl,dev,xcfctl,xcfsys,xcfmbr,omvs,cachectl,cachedev,cachestat,raid,hfsglobal,hfsbuf}` |
| 75 | RMF page data set activity | `smf075`, `smf075_page` |
| 78 | RMF virtual storage (2) and I/O queuing (3) activity | `smf078` and `smf078_{common,cfgctl,cfg,ioq,ioqglob,iop,hpav}` |
| 80 | RACF processing (logons, resource access, profile changes) | `smf080`, `smf080_reloc`, `smf080_reloc2`, plus reference tables `smf080_events`, `smf080_qualifiers`, `smf080_types` |
| 84 | JES monitoring facility: JES2 resource usage (subtype 21) | `smf084`, `smf084_mem`, `smf084_res` |

Records of any other type are counted in the summary but not stored.

Type 42 subtypes and their tables:

| Subtype | Content | Tables |
|---|---|---|
| all | Header and product section | `smf042` |
| 1 | Buffer Manager Facility (BMF) totals and storage class summary | `smf042_bmf`, `smf042_bmf_sc` |
| 2 | Control unit cache and SMS volume status | `smf042_cu`, `smf042_cu_vol` |
| 3 | SMS configuration change | `smf042_ea` |
| 4 | Extended format data set (EXT) copy statistics | `smf042_ext` |
| 5 | Storage class and volume response time (VTOC, VTOC index, VVDS) | `smf042_sc`, `smf042_vol`, `smf042_vol_vtoc`, `smf042_vol_vtix`, `smf042_vol_vvds` |
| 6 | Job and data set I/O statistics | `smf042_job`, `smf042_ds` |
| 9 | B37/D37/E37 abend | `smf042_x37` |
| 21 | PDS/PDSE member delete | `smf042_mdel`, `smf042_mdel_alias` |
| 24 | PDS/PDSE member add/replace | `smf042_madd`, `smf042_madd_alias` |
| 27 | VTOC audit log | `smf042_vtoc` |

## Building

smf2db is written in Go and uses [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3), which needs **CGO and a C compiler**.

```sh
git clone https://github.com/philippeleite/smf2db.git
cd smf2db
CGO_ENABLED=1 go build -o smf2db .
```

- **Linux / macOS**: install `gcc` or `clang` (for example `build-essential` or the Xcode command line tools).
- **Windows**: install a MinGW-w64 `gcc`, or use [Zig](https://ziglang.org/) as the C compiler:

  ```sh
  CGO_ENABLED=1 CC="zig cc -target x86_64-windows-gnu" go build -o smf2db.exe .
  ```

A binary built with `CGO_ENABLED=0` compiles but fails at run time with *"go-sqlite3 requires cgo to work"*.

## Getting SMF data from z/OS

1. Dump the records you need with **IFASMFDP** (or IFASMFDL for log streams) into a sequential data set. For example:

   ```jcl
   //SMFDUMP  EXEC PGM=IFASMFDP
   //DUMPIN   DD DISP=SHR,DSN=SYS1.MANX
   //DUMPOUT  DD DSN=MY.SMF.DUMP,DISP=(NEW,CATLG),
   //            SPACE=(CYL,(100,100),RLSE),
   //            DCB=(RECFM=VBS,LRECL=32760,BLKSIZE=0)
   //SYSPRINT DD SYSOUT=*
   //SYSIN    DD *
     INDD(DUMPIN,OPTIONS(DUMP))
     OUTDD(DUMPOUT,TYPE(4:6,14:18,21,30,41,42,60:62,64:66,70:75,78,80,84))
   /*
   ```

2. Transfer the data set to your workstation **in binary, keeping the record descriptor words (RDW)**. With the z/OS FTP server:

   ```text
   ftp> binary
   ftp> quote site rdw
   ftp> get 'MY.SMF.DUMP' smf.dump
   ```

   A plain binary transfer without the RDW loses the record boundaries. smf2db detects that and stops with *"SMF file must be transferred preserving RDW information"*.

## Usage

```text
smf2db [options] smfFile [dbFile]

  smfFile   input SMF file, with RDW
  dbFile    output SQLite database (default: smf.db)

Options:
  -t types  record types to load, comma-separated, with ranges: -t 30  -t 14,15  -t 4:6,30,42
  -v        print the version
  -h        print this help
```

> **Warning:** `dbFile` is deleted and recreated on every run. Use a new file name for each load you want to keep.

At the end smf2db prints how many records of each type it read (including types it does not decode) and the elapsed time.

Examples:

```sh
# Everything smf2db understands, into smf.db
smf2db smf.dump

# Only job/step accounting and DFSMS statistics, into a named database
smf2db -t 30,42 smf.dump accounting.db

# Data set activity (open/close, delete, rename, VSAM open/close), VVDS updates and catalog defines, deletes and alters
smf2db -t 14:18,60:62,64:66 smf.dump datasets.db
```

## Database layout

- Each record becomes one row in its type's main table (`smf030`, `smf042`, …). The `ID` column numbers the records of that type in file order.
- Repeating sections go to child tables with the same `ID`, so you join them with `USING (ID)`.
- Column names are the IBM field names, so look up a column's meaning and units in the manual.
- **Text** fields are converted from EBCDIC to ASCII and keep the mainframe's blank padding. Use `trim()` when you compare or display them (`WHERE trim(SMF30JBN) = 'PAYROLL'`).
- **Dates and times** are text in `YYYY-MM-DD HH:MM:SS.mmm` form. SMF keeps date and time in separate fields, so a date column holds `2017-04-26 00:00:00.000` and a time column holds `0000-00-00 14:00:00.160`. To combine them: `substr(SMF30DTE, 1, 10) || ' ' || substr(SMF30TME, 12)`.
- Fields in hardware clock (TOD) format, such as `SMF42PTS` or `S42EXSTM`, are converted to UTC timestamps. SMF header times are local time.
- **Numbers** are stored as recorded. Many are not in plain units: for example `SMF30CPT` is in hundredths of a second and `SMF21DBW` in 4 KB blocks. Check the manual.
- **Binary structures** that smf2db does not decode, such as a raw VVR or a RACF token, are stored as `BLOB`.

## Example queries

Open the database with `sqlite3 smf.db` (`.mode column` and `.headers on` make the output readable), or with any SQLite client.

### Jobs and steps (type 30)

The 10 steps that used the most CPU (subtype 4 = step end):

```sql
SELECT trim(j.SMF30JBN) AS job,
       trim(j.SMF30STM) AS step,
       trim(j.SMF30PGM) AS program,
       p.SMF30CPT / 100.0 AS cpu_sec,
       p.SMF30CPS / 100.0 AS srb_sec
FROM smf030 j
JOIN smf030_proc p USING (ID)
WHERE j.SMF30TYP = 4
ORDER BY p.SMF30CPT DESC
LIMIT 10;
```

Steps that ended with a non-zero completion code:

```sql
SELECT substr(j.SMF30DTE, 1, 10) || ' ' || substr(j.SMF30TME, 12, 8) AS ended,
       trim(j.SMF30JBN) AS job,
       trim(j.SMF30STM) AS step,
       printf('%04X', c.SMF30SCC) AS completion_code
FROM smf030 j
JOIN smf030_comp c USING (ID)
WHERE j.SMF30TYP = 4 AND c.SMF30SCC <> 0
ORDER BY ended;
```

EXCP counts per DD statement for one job:

```sql
SELECT trim(j.SMF30STM) AS step,
       trim(x.SMF30DDN) AS ddname,
       x.SMF30BLK AS blocks
FROM smf030 j
JOIN smf030_excp x USING (ID)
WHERE j.SMF30TYP = 4 AND trim(j.SMF30JBN) = 'MYJOB'
ORDER BY x.SMF30BLK DESC;
```

### Data sets (types 14, 15, 17, 18)

Most frequently used data sets (one type 14/15 record per close):

```sql
SELECT trim(dsn) AS dataset,
       sum(input) AS input_closes,
       sum(output) AS output_closes
FROM (SELECT JFCBDSNM AS dsn, 1 AS input, 0 AS output FROM smf014
      UNION ALL
      SELECT JFCBDSNM, 0, 1 FROM smf015)
GROUP BY dsn
ORDER BY count(*) DESC
LIMIT 20;
```

Who deleted which data sets, and on which volume:

```sql
SELECT substr(SMF17DTE, 1, 10) || ' ' || substr(SMF17TME, 12, 8) AS deleted_at,
       trim(SMF17JBN) AS job,
       trim(SMF17DSN) AS dataset,
       v.SMF17FVL AS volume
FROM smf017
JOIN smf017_vol v USING (ID)
WHERE trim(SMF17DSN) LIKE 'PROD.%'
ORDER BY deleted_at;
```

Data set renames:

```sql
SELECT trim(SMF18JBN) AS job,
       trim(SMF180DS) AS old_name,
       trim(SMF18NDS) AS new_name
FROM smf018;
```

### DFSMS I/O performance (type 42)

Busiest data sets, with their I/O-weighted average response time:

```sql
SELECT trim(S42DSN) AS dataset,
       trim(S42DSVOL) AS volume,
       sum(S42DSION) AS ios,
       round(sum(S42DSR1U * 1.0 * S42DSION) / sum(S42DSION)) AS avg_response_us
FROM smf042_ds
WHERE S42DSION > 0
GROUP BY S42DSN, S42DSVOL
ORDER BY ios DESC
LIMIT 20;
```

Response time and cache hit ratio by storage class:

```sql
SELECT trim(S42SCRNN) AS storage_class,
       sum(S42SCION) AS ios,
       round(sum(S42SCR1U * 1.0 * S42SCION) / sum(S42SCION)) AS avg_response_us,
       round(100.0 * sum(S42SCHIT) / nullif(sum(S42SCCND), 0), 1) AS cache_hit_pct
FROM smf042_sc
GROUP BY S42SCRNN
ORDER BY ios DESC;
```

Most frequently updated PDS/PDSE members:

```sql
SELECT trim(SMF42PJB) AS job,
       trim(SMF42PDS) AS library,
       trim(SMF42PMN) AS member,
       count(*) AS updates
FROM smf042_madd
GROUP BY 1, 2, 3
ORDER BY updates DESC
LIMIT 20;
```

### Processor activity (type 70)

Busy percentage of each PR/SM partition's logical processors, by CPU type and RMF interval (`SMF70PDT` is in microseconds, `SMF70INT` in seconds):

```sql
SELECT substr(r.SMF70DAT, 1, 10) || ' ' || substr(r.SMF70IST, 12, 5) AS interval_start,
       trim(p.SMF70LPM) AS partition,
       trim(c.SMF70CIN) AS cpu_type,
       count(*) AS logical_cpus,
       round(100.0 * sum(l.SMF70PDT) / 1e6 / (r.SMF70INT * count(*)), 1) AS busy_pct
FROM smf070 r
JOIN smf070_lcpu l USING (ID)
JOIN smf070_lpar p ON p.ID = l.ID AND p.SMF70LPN = l.SMF70LPN
JOIN smf070_cpuid c ON c.ID = l.ID AND c.SMF70CIX = l.SMF70CIX
WHERE trim(p.SMF70LPM) <> 'PHYSICAL'
GROUP BY r.ID, p.SMF70LPM, c.SMF70CIN
ORDER BY interval_start, partition, cpu_type;
```

### Workload activity (type 72)

Transactions, average response time and CPU service units by WLM service class period (`R723CTET` is in 1024-microsecond units):

```sql
SELECT trim(w.R723MCNM) AS service_class,
       p.R723CPER AS period,
       sum(p.R723CRCP) AS transactions,
       round(sum(p.R723CTET) * 1.024 / nullif(sum(p.R723CRCP), 0), 1) AS avg_response_ms,
       round(sum(p.R723CCPU + p.R723CSRB)) AS cpu_service_units
FROM smf072 r
JOIN smf072_wlm w USING (ID)
JOIN smf072_scp p USING (ID)
WHERE r.SMF72STY = 3
GROUP BY w.R723MCNM, p.R723CPER
ORDER BY transactions DESC;
```

Some IBM field names contain `#` (for example `R723MTV#`); quote them in SQL: `SELECT "R723MTV#" FROM smf072_wlm`.

### Security (type 80)

RACF access violations by class, user and resource, with readable event names from the reference tables (relocate data type 17 is the class, 1 the resource name):

```sql
SELECT e.EVENT, q.DESCRIPTION,
       trim(r.SMF80USR) AS userid,
       c.SMF80DTA AS class,
       n.SMF80DTA AS resource,
       count(*) AS events
FROM smf080 r
JOIN smf080_events e USING (SMF80EVT)
JOIN smf080_qualifiers q ON q.SMF80EVT = r.SMF80EVT AND q.SMF80EVQ = r.SMF80EVQ
LEFT JOIN smf080_reloc c ON c.ID = r.ID AND c.SMF80DTP = 17
LEFT JOIN smf080_reloc n ON n.ID = r.ID AND n.SMF80DTP = 1
WHERE r.SMF80DES & 0x8000            -- violation flag
GROUP BY 1, 2, 3, 4, 5
ORDER BY events DESC
LIMIT 20;
```

### Tapes, printing and catalog activity

Tape volumes by data written (`SMF21DBW` counts 4 KB blocks), with error counters:

```sql
SELECT trim(SMF21VOL) AS volser,
       count(*) AS demounts,
       round(sum(SMF21DBW) * 4096 / 1e9, 2) AS gb_written,
       sum(SMF21TR + SMF21TW + SMF21PR + SMF21PW) AS errors
FROM smf021
GROUP BY SMF21VOL
ORDER BY gb_written DESC;
```

Bytes sent to each IP PrintWay printer:

```sql
SELECT trim(SMF6PRTQ) AS print_queue,
       SMF6IP1 || '.' || SMF6IP2 || '.' || SMF6IP3 || '.' || SMF6IP4 AS printer_ip,
       count(*) AS datasets,
       sum(SMF6BYTE) AS bytes
FROM smf006_ipw
GROUP BY 1, 2
ORDER BY bytes DESC;
```

VVDS updates per job and operation (`UP`date, `IN`sert, `DE`lete), with the average request time:

```sql
SELECT trim(SMF60JNM) AS job,
       SMF60SUB AS operation,
       count(*) AS vvr_updates,
       round(avg(SMF60RQT)) AS avg_request_us
FROM smf060
GROUP BY 1, 2
ORDER BY vvr_updates DESC;
```

### Exploring the schema

```sql
-- Tables created by this load
SELECT name FROM sqlite_master WHERE type = 'table' ORDER BY name;

-- Columns of a table
SELECT name, type FROM pragma_table_info('smf042_ds');
```

## Reference

Record layouts come from IBM **z/OS MVS System Management Facilities (SMF)**, SA38-0667 (z/OS 3.1 edition), except type 80 (RACF), which comes from **z/OS Security Server RACF Macros and Interfaces**, SA23-2288. They describe every field, its units and its flag bits. The manuals are not distributed with this project.

Some records differ from the manual. smf2db follows the data, and [AGENTS.md](AGENTS.md) lists each known difference. For example, many type 42 records leave the triplet count at 0, and some fields documented as EBCDIC actually contain binary data.

## Limitations

- Only the record types listed above are decoded. Within some types, not every section is decoded:
  - type 6: only the IP PrintWay sections (JES2, JES3, PSF and External Writer records get only the common header row);
  - type 42: only the subtypes in the table above, and within them not the concurrent copy sessions (4), the volume metrics, system I/O and background sections (5), or the synchronous I/O sections (6).
- Records are processed in a single thread, and the whole database is rebuilt on each run. There is no incremental load.
- Spanned records larger than 64 KB are not supported.

## Development

```sh
go vet ./...
go test ./...                    # unit tests
CGO_ENABLED=1 go test ./...      # also runs the end-to-end test when data/smftest.rdw is present
```

The end-to-end test converts a local SMF file and compares the record counts against the IFASMFDP report of the dump job. Sample data is not included in the repository.

[AGENTS.md](AGENTS.md) describes the code structure, the parser conventions and how to add a new record type.

## License

smf2db is released under the [MIT License](LICENSE).
