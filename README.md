# nxcurl

**nxcurl** is a small HTTP CLI aimed at the same workflows as Postman, but tuned for **terminal use** and **coding agents** (Cursor, Claude Code, and similar). It runs real requests, keeps a **JSONL history**, supports **environments** with `{{variable}}` substitution, can **import HAR and Postman Collection v2**, and prints either **human-friendly** output or **`--json`** structured output for models.

## Install (from source)

```bash
go build -o nxcurl ./cmd/nxcurl
```

Put the binary on your `PATH`. Data lives under `~/.nxcurl/`.

Embed version when building locally:

```bash
go build -ldflags="-s -w -X github.com/chenchunaidu/nxcurl/internal/version.Version=v0.1.0" -o nxcurl ./cmd/nxcurl
```

## Prebuilt binaries (GitHub Releases)

Releases follow the same automation as [tars](https://github.com/chenchunaidu/tars): each push to **`main`** that does not already have a tag on `HEAD` gets a new **patch** tag (`v0.0.1` → `v0.0.2`), then CI builds archives and publishes a GitHub Release. Pushing a **`v*`** tag yourself runs the same build and release for that tag.

| Platform | Archive (stable name on each release) |
|----------|----------------------------------------|
| Linux x86_64 | `nxcurl_linux_amd64.tar.gz` |
| Linux arm64 | `nxcurl_linux_arm64.tar.gz` |
| macOS x86_64 | `nxcurl_darwin_amd64.tar.gz` |
| macOS Apple Silicon | `nxcurl_darwin_arm64.tar.gz` |
| Windows x86_64 | `nxcurl_windows_amd64.zip` |
| Windows arm64 | `nxcurl_windows_arm64.zip` |

Versioned filenames (e.g. `nxcurl_v0.0.2_linux_amd64.tar.gz`) and **`checksums.txt`** are attached to the same release. **Latest** direct downloads (replace `OWNER/REPO`):

`https://github.com/OWNER/REPO/releases/latest/download/nxcurl_darwin_arm64.tar.gz`

## [tars](https://github.com/chenchunaidu/tars) integration

- Install the binary with **tars** using a formula (see `packaging/tars-formula.example.json` — fill `url` / `sha256` from your GitHub Release assets).
- After install, run **`nxcurl docs`**: it refreshes `~/.nxcurl/docs/SKILL.md` and **prints** the full agent-oriented skill (inputs, triggers, tool name per command) to stdout. Point agents at that output or the on-disk file, or fold it into your workflow with **`tars connect`** / project rules.
- For agent turns, prefer **`nxcurl … --json`** so the model gets parseable status, headers, and body (plus `response_body_json` when the body is JSON).

## Commands

| Command | Purpose |
|--------|---------|
| `nxcurl run <url>` | HTTP request (`-X`, `-H`, `-d`, `-e` / `--env`, `--json`, `--no-history`) |
| `nxcurl send <file.json>` | Run a saved request (e.g. after `import`) |
| `nxcurl history list` | Newest-first log (`--limit`, `--json`) |
| `nxcurl history show <id>` | One stored exchange (`--json`) |
| `nxcurl history replay <id>` | Re-run a past request |
| `nxcurl import har <file.har>` | Writes request JSON under `~/.nxcurl/collections/` |
| `nxcurl import postman <collection.json>` | Same for Postman v2 |
| `nxcurl env init <name>` | Create `~/.nxcurl/envs/<name>.json` |
| `nxcurl env path <name>` | Print env file path |
| `nxcurl docs` | Print agent-oriented `SKILL.md` (and write `~/.nxcurl/docs/SKILL.md`) |

## Environments

1. `nxcurl env init prod`
2. Edit `~/.nxcurl/envs/prod.json` as a flat JSON object: `{ "API_KEY": "…", "BASE": "https://api.example.com" }`
3. Use in requests: `nxcurl run '{{BASE}}/v1/foo' -H 'Authorization: Bearer {{API_KEY}}' -e prod`

## Layout

| Path | Role |
|------|------|
| `~/.nxcurl/history.jsonl` | Append-only request/response log |
| `~/.nxcurl/envs/*.json` | Named environments |
| `~/.nxcurl/collections/*/` | Imported / saved request JSON files |
| `~/.nxcurl/docs/SKILL.md` | Agent-oriented command reference (refreshed by `nxcurl docs`) |

## Status

Early MVP: core run/history/env/import/docs are in place; collection editing UI, auth helpers, and richer Postman/HAR edge cases can be extended incrementally.
