## Summary
Migrate remaining sections (s05-s11) from 01-foundations/ subfolder structure to root-level folders, following the pattern from PR #327.

## Current Structure to Fix
- s05 Composition -> 06-composition/
- s06 Strings -> 07-strings-and-text/
- s07 Modules -> 08-modules-and-packages/
- s08 CLI/IO -> 09-io-and-cli/
- s09 Web/DB -> 10-web-and-database/
- s10 Concurrency -> 11-concurrency/
- s11 Patterns -> 12-concurrency-patterns/
- s12 Quality -> 13-quality-and-performance/
- s13 Architecture -> 14-application-architecture/
- s14 Flagship -> 15-code-generation/

## Expected New Structure
- 05-composition/
- 06-strings-and-text/
- 05-packages-io/01-modules-and-packages/
- 08-cli-and-io/
- 06-backend-db/01-web-and-database/
- 07-concurrency/01-concurrency/
- 07-concurrency/02-concurrency-patterns/
- 08-quality-test/01-quality-and-performance/
- 13-application-architecture/
- 11-flagship/02-code-generation/

## Changes Needed
1. Move folders to root level
2. Update curriculum.v2.json paths
3. Fix lesson file paths (README.md, main.go)
4. Update README.md stage table
5. Update validator paths
6. Run validation and tests