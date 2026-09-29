# plugin-file

The `file` multi-role state-provision verb — a host-coupled plugin that both
**checks** a path's state and **renders** an idempotent file-creation.

## What it provides

| Capability | Surface |
|---|---|
| `verb:file` | the `file:` check verb (`do:assert`) — stat a path and assert `exists` / `mode` / `owner` / `group` / `filetype` / `contains` / `link_target` / `sha256` |

The verb is a `kit.CheckVerbProvider` (CHECK) that also implements
`kit.ProvisionActor` (ACT), so charly registers the multi-role adapter:

- **CHECK** — stats the path via the live check engine and asserts the requested
  attributes.
- **ACT** — renders an idempotent runtime file-creation (`mkdir` / `touch` /
  `cat` heredoc + `chmod`).

It is **host-coupled** (compiled-in only): the check runs against the live
`kit.CheckContext` and the act renders a provision script. The `contains` /
`link_target` matchers accept a bare scalar and default it to the `contains`
(substring) operator, reusing `sdk.MatchAll`.

## How to use it

Compose the plugin candy in a box or check bed's `candy:` list:

```yaml
- '@github.com/opencharly/plugin-file/candy/plugin-file:<tag>'
```

Then author the verb in a plan:

```yaml
- check: /etc/hostname exists
  id: file-verb-dispatches
  file:
    file: /etc/hostname
    exists: true
  context: [runtime]
```

## Layout

- `candy/plugin-file/` — the plugin module: `plugin.go` (the
  `kit.CheckVerbProvider` + `kit.ProvisionActor` + `NewCheckVerb()`/`NewMeta()`),
  `schema/file.cue` (the self-contained `#FileInput`),
  `params/cue_types_gen.go`, `cmd/serve/main.go`.
- `charly.yml` — the root project manifest (`discover: candy`).
- `.github/workflows/tag-on-merge.yml` — CalVer tag + `CHANGELOG/` on merge.

## Related

- Owning skill: `/charly-internals:plugin` — the plugin/provider model, including
  the `kit` check-verb and provision-actor contracts. This candy carries no
  `skill:` entity of its own; the gap is tracked in
  [opencharly/opencharly#291](https://github.com/opencharly/opencharly/issues/291).
- [`opencharly/charly`](https://github.com/opencharly/charly) — the charly CLI.
