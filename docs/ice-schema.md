# ICE v1 schema

manifest.json: version, authority, files (path, sha256, package, imports), symbols (name, kind, file, start/end lines, optional receiver and syntax call names), implementations (interface, asserted expression, file, line). Companion architecture/symbol/dependency/test/implementation maps are deterministically derived from this index. integration-seams and decisions mirror authoritative docs JSON.

Build uses go/parser; function calls are syntactic references and test maps list calls inside Test* functions. Implementations list explicit `var _ Interface = (*Type)(nil)` assertions. `go test` checks those assertions. It does not infer all assignable types. Configuration contracts, gaps and decisions live in hashed docs with symbols pointing to Go contracts. Generated ICE content is excluded from source hashes to avoid self-reference; validation compares each companion artifact against freshly derived content.
