# qdrant-cli

`qdrant-cli` is a lightweight command-line tool for inspecting and managing **Qdrant** collections over the **gRPC API**.
It is designed for quick operational tasks such as listing collections, describing collection configuration, scrolling points, and inspecting payload data in human-friendly or machine-friendly formats.

The tool is intentionally simple, dependency-light, and scriptable.

---

## Features

* List existing Qdrant collections
* Create and delete collections
* Describe collection configuration and statistics
* Scroll points from a collection
* Filter points by payload field
* Multiple output formats:

  * **TSV** (default, script-friendly)
  * **Table** (human-readable)
  * **JSON** (raw API output)
* Select which payload fields to display
* Works directly with Qdrant **gRPC** (no REST dependency)

---

## Installation

### Build from source

```bash
go build -o qdrant-cli
```

### Requirements

* Go 1.20+
* Running Qdrant instance with gRPC enabled (default port `6334`)

---

## Global Options

| Flag       | Description                           | Default       |
| ---------- | ------------------------------------- | ------------- |
| `--host`   | Qdrant host                           | `localhost`   |
| `--port`   | Qdrant gRPC port                      | `6334`        |
| `--format` | Output format: `tsv`, `table`, `json` | `tsv`         |
| `--filter` | Payload filter (`key=value`)          | *(none)*      |
| `--fields` | Fields to display (comma-separated)   | `id,filename` |
| `--limit`  | Scroll page size                      | `5`           |

---

## Commands

### List collections

```bash
qdrant-cli --list
```

**Table output (default):**

```
COLLECTION
----------
cbf_images
test_vectors
```

**JSON output:**

```bash
qdrant-cli --list --format json
```

---

### Describe a collection

```bash
# just run make if it is installed on your system
make
# or if you want bare metal
qdrant-cli --info cbf_images
```

**Table output:**

```
FIELD               VALUE
-----               -----
Status              GREEN
Segments            4
Points              3
Indexed vectors     0
Shard number        1
Replication factor  1
Write consistency   1
On-disk payload     true
Vector size         50176
Distance            Cosine
```

**JSON output:**

```bash
qdrant-cli --info cbf_images --format json
```

---

### Create a collection

```bash
qdrant-cli --create cbf_images --dim 50176
```

Creates a collection with cosine distance and the specified vector dimension.

---

### Delete a collection

```bash
qdrant-cli --delete cbf_images
```

⚠️ This permanently deletes the collection.

---

### Scroll points (inspect records)

```bash
qdrant-cli --scroll cbf_images
```

#### Default output (TSV)

One row per point, tab-separated:

```text
7f789aa0-012a-4a37-a297-e3b228ae8486	CeO2_test1_PIL10_007_180.cbf
c17068f2-4622-4d20-9935-0bb5410640e8	CeO2_test1_PIL10_007_179.cbf
```

This format is ideal for scripting:

```bash
qdrant-cli --scroll cbf_images | awk '{print $2}'
```

---

### Select which fields to show

```bash
qdrant-cli --scroll cbf_images --fields id,filename,width,height
```

Output:

```text
7f789aa0-012a-4a37-a297-e3b228ae8486	CeO2_test1_PIL10_007_180.cbf	2463	2527
```

---

### Table format (human-readable)

```bash
qdrant-cli --scroll cbf_images --format table
```

```
ID                                    FIELD     VALUE
--                                    -----     -----
7f789aa0-012a-4a37-a297-e3b228ae8486  filename  CeO2_test1_PIL10_007_180.cbf
                                      width     2463
                                      height    2527
```

---

### JSON format (raw API output)

```bash
qdrant-cli --scroll cbf_images --format json
```

Useful for debugging or programmatic processing.

---

### Filter by payload field

```bash
qdrant-cli --scroll cbf_images --filter engine=cbf2go
```

Filter syntax:

```
payload_key=value
```

Currently supports exact keyword match.

---

## Design Notes

* Uses **Qdrant gRPC client**, not REST
* Output formats are explicitly separated:

  * TSV for automation
  * Table for humans
  * JSON for debugging
* Payload field order is deterministic and user-controlled
* Numeric point IDs are skipped by default in TSV output (UUID-focused workflow)

---

## Example Use Cases

* Inspect vector collections on a production Qdrant instance
* Verify ingestion pipelines
* Pipe UUIDs and filenames into shell scripts
* Quickly audit payload metadata
* Debug collection configuration without opening dashboards

---

## Future Improvements (Ideas)

* Pagination / streaming scroll
* `--sort` and `--where` expressions
* CSV output
* Vector preview / norms
* Collection stats summary mode

---

## License

MIT (or your preferred license)
