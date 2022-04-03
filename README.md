# SQL Migrations in Go

This repository demonstrates how a Go program working with an SQLite database can handle database schema migrations.

See [the initial code version][init-ver] that does not have a concept of schema version
and [how it had to change][new-ver] to support both new empty databases
and existing ones with the old schema version,
automatically upgrading them.

[init-ver]: https://github.com/artyom/sqlite-alter-schema/blob/fdd8552e5238ec4b9f9f129e2a9d4d4aa8e229ae/main.go
[new-ver]: https://github.com/artyom/sqlite-alter-schema/commit/54779f851dbbc7617586be48c969e5aee863e1cb#diff-2873f79a86c0d8b3335cd7731b0ecf7dd4301eb19a82ef7a1cba7589b5252261
