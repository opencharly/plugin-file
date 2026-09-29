# AGENTS.md — plugin-file

Standalone plugin repo for the `file` capability (`verb:file`) — the multi-role
state-provision verb (CHECK + ACT), a host-coupled `kit.CheckVerbProvider` /
`kit.ProvisionActor`. The plugin is a Go module at `candy/plugin-file/` (module
path `github.com/opencharly/plugin-file/candy/plugin-file`); the root `charly.yml`
only declares `discover: candy` so the repo is a project and its candy is
scanned.

Canonical files:

- `candy/plugin-file/charly.yml` — the `plugin-file:` candy entity (`plugin:`
  block, `plan:` check).
- `candy/plugin-file/plugin.go` — the `kit.CheckVerbProvider` +
  `kit.ProvisionActor` (`NewCheckVerb()` + `NewMeta()` + `RunVerb` / act).
- `candy/plugin-file/schema/file.cue` — the self-contained `#FileInput`.
- `candy/plugin-file/params/cue_types_gen.go` — generated params (do not
  hand-edit).
- `.github/workflows/tag-on-merge.yml` — CalVer tag + `CHANGELOG/` on merge.
- `README.md` — user overview only; never agent guidance.

## Load these skills first (R0)

- `/charly-internals:plugin` — the plugin authoring reference: the `plugin:`
  block, the unified Provider model, the `kit` check-verb / provision-actor
  contracts, the per-plugin CUE-schema contract, placement. Load before touching
  the provider or schema.
- `/charly-check:check` — the declarative check-verb surface the `file:` verb is
  authored through.
- `/charly-internals:git-workflow` — before any git/PR action.

## Build / validate / test

- `go build ./...` in `candy/plugin-file/` — compile the plugin module.
- `go test ./...` in `candy/plugin-file/` — the plugin's Go tests.
- `charly box validate` at the repo root — the structural check (the candy +
  `plugin:` block, CUE schema).
- The merge gate is the **org-wide** `charly/pr-validator` (required check
  `validate / validate`, defined in `opencharly/.github`); this repo has **no**
  per-repo candy gate.

## Modify this repo

- Edit the `plugin-file:` candy entity, the Go source, and `schema/file.cue`
  **together** — the schema is the single source for the `params/` struct.
- The verb is host-coupled and compiled-in; keep both roles (CHECK `RunVerb` and
  ACT provision) working, and reuse `sdk.MatchAll` for matcher evaluation.

## Landing

- PR-only. Every change lands through a pull request; the org-required
  `charly/pr-validator` validates the diff and body and arms native auto-merge on
  PASS. Direct pushes to `main` are blocked.
- History lives in `CHANGELOG/` (written by `tag-on-merge` at merge time); the PR
  body IS the changelog.
- The authoritative rulebook is the umbrella `AGENTS.md` in
  `opencharly/opencharly` and `charly/AGENTS.md` in the charly repo. Do not
  restate its rules here.
