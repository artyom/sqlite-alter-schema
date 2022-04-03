# SQL Migrations in Go

This repository demonstrates how a Go program working with an SQLite database can handle database schema migrations.

See [the initial code version][init-ver] that does not have a concept of schema version
and how it had to change to support both new empty databases
and existing ones with the old schema version,
automatically upgrading them.

[init-ver]: https://github.com/artyom/sqlite-alter-schema/blob/fdd8552e5238ec4b9f9f129e2a9d4d4aa8e229ae/main.go
