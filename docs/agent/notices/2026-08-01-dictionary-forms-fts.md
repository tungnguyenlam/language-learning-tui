# Dictionary Forms FTS Index

Status: active
Scope: `internal/storage/sqlite` dictionary search
Related: `migrations.go` migration 26, `dictionary.go` `Search` and `ImportEntries`

## Why It Matters

`dictionary_fts.forms` remains denormalized and `UNINDEXED` for compatibility with the main dictionary table. Normal inflection lookups use `dictionary_forms_fts`, joined by matching FTS rowids, so typing in the dictionary does not scan all imported entries.

## Required Behavior

Any complete dictionary rebuild must delete and repopulate both FTS tables with matching rowids. Do not change the main-table rowid assignment or update only one of the two indexes.

## Revisit When

The dictionary storage is normalized into ordinary tables or the main FTS table indexes `forms` directly.
