# gasgo

[![GoDoc](https://godoc.org/github.com/crgimenes/gasgo?status.png)](https://pkg.go.dev/github.com/gasgo/atomic?gasgo=doc)
[![Go project version](https://badge.fury.io/go/github.com%2Fcrgimenes%2Fgasgo.svg)](https://badge.fury.io/go/github.com%2Fcrgimenes%2Fgasgo)
[![MIT Licensed](https://img.shields.io/badge/license-MIT-green.svg)](https://tldrlegal.com/license/mit-license)

## WIP

**This is a work in progress.**

Currently writing the specification and detailing how it works.

---

## Equivalent types between Go and SQL databases

| Go value type                         | PostgreSQL column type              | SQLite column type                    | Caveat                                                                        |
| ------------------------------------- | ----------------------------------- | ------------------------------------- | ----------------------------------------------------------------------------  |
| `bool`                                | `boolean`                           | `INTEGER` (0 = false, 1 = true)       | SQLite has no native boolean; uses 0/1 integers.                              |
| `int` (= `int64` on 64-bit) / `int64` | `bigint`                            | `INTEGER`                             | Both store 64-bit signed integers.                                            |
| `int32`                               | `integer`                           | `INTEGER`                             | 32-bit signed in PostgreSQL; SQLite stores as 64-bit but preserves the value. |
| `uint`, `uint64`                      | `bigint` + `CHECK (col >= 0)`       | `INTEGER`                             | PostgreSQL lacks unsigned ints; enforce non-negative with a CHECK constraint. |
| `float64`                             | `double precision`                  | `REAL`                                | 64-bit IEEE-754 on both.                                               |
| `string`                              | `text` / `varchar(n)`               | `TEXT`                                | Variable-length UTF-8.                                                 |
| `[]byte`                              | `bytea`                             | `BLOB`                                | Raw binary data.                                                       |
| `time.Time`                           | `timestamptz`                       | `TEXT` (ISO-8601) or `INTEGER` epoch  | Go drivers convert automatically.                                       |
| “Decimal” (`string` or `big.Rat`)     | `numeric(p,s)`                      | `NUMERIC` (affinity)                  | Recommended for exact monetary values.                                |

---
## Contributing

- Fork the repo on GitHub
- Clone the project to your own machine
- Create a *branch* with your modifications `git checkout -b fantastic-feature`.
- Then _commit_ your changes `git commit -m 'Implementation of new fantastic feature'`
- Make a _push_ to your _branch_ `git push origin fantastic-feature`.
- Submit a **Pull Request** so that we can review your changes
