## Scope
- Lock the v2.1 curriculum structure across the canonical section map
- Add validator support for item status and require honest path validation
- Remove superseded legacy lesson surfaces after their canonical replacements exist
- Keep `ARCHITECTURE.md` thin and move lesson-level truth into `curriculum.v2.json`

## Affected Sections / Lesson IDs
- s00: HC.1-HC.5, trim deprecated HC.6-HC.8 lesson surfaces
- s03: FE.8-FE.10 registration and sequencing support
- s04: TI.1-TI.5, CO.1-CO.3, ST.1-ST.6 canonical relocation and renumbering
- s06: API.1-API.8, DB.1-DB.8 canonical lesson surfaces
- s07: GC.0-GC.5, CT.1-CT.5, TM.1-TM.3, CP.1-CP.6, SY.1-SY.4
- s08: TE.1-TE.10, PR.1-PR.6, PD.1-PD.3
- s09: ARCH.1-ARCH.3, SEC.1-SEC.5
- s10: CFG.1-CFG.3, SL.1-SL.3, GS.1-GS.3, OPS.1-OPS.3, DOCKER.1-DOCKER.6, DEPLOY.1-DEPLOY.3, CG.1-CG.3
- s11: FG.1 flagship entry point

## Why
- The repo needs one validator-enforced v2.1 contract for section layout, lesson IDs, and canonical paths.
- Legacy folder trees create drift and hide broken references unless they are removed once replacements exist.
- Code generation belongs in s10 production context, and the flagship needs an explicit FG.1 entry point after that move.

## Acceptance
- `ARCHITECTURE.md`, `curriculum.v2.json`, and the validator agree on the final 12-section map
- superseded legacy lesson trees are deleted from the repo
- local verification passes: `go build ./...`, `go vet ./...`, `go test ./...`, `go run ./scripts/validate_curriculum.go`
