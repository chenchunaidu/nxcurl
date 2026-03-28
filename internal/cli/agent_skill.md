---
name: nxcurl
description: HTTP CLI for terminals and coding agents — run requests, JSONL history, env substitution {{VAR}}, HAR/Postman import, structured --json output.
---

# nxcurl — agent skill

**Tool name (binary):** `nxcurl`  
**Purpose:** Execute real HTTP requests, persist exchanges to local history, substitute secrets from named env files, and optionally emit **machine-readable JSON** for agents.

**Global inputs (most subcommands):**

| Flag | Meaning |
|------|---------|
| `-e`, `--env <name>` | Load `~/.nxcurl/envs/<name>.json` for `{{KEY}}` substitution in URL, headers, and body |
| `--json` | Print structured JSON on stdout (use for agent turns) |

**Data layout:**

| Path | Role |
|------|------|
| `~/.nxcurl/history.jsonl` | Append-only log of requests/responses |
| `~/.nxcurl/envs/<name>.json` | Flat JSON object of string keys → values for `{{KEY}}` |
| `~/.nxcurl/collections/` | Imported or saved request JSON (for `send`) |
| `~/.nxcurl/docs/SKILL.md` | Copy of this document (refreshed by `nxcurl docs`) |

---

## Command: `nxcurl run`

**Tool name:** `nxcurl` (subcommand `run`)

**Inputs:**

- Positional: `<url>` — request URL (may contain `{{VAR}}` placeholders)
- `-X`, `--request <method>` — HTTP method (default `GET` if omitted)
- `-H`, `--header 'Name: value'` — repeatable
- `-d`, `--data <body>` — request body
- `--no-history` — do not append this exchange to history
- Global: `-e` / `--env`, `--json`

**Triggers (when / how to use):**

- Ad-hoc API calls when you need a **one-off** request with explicit method, headers, and body.
- Use `--json` when an agent must parse status, headers, and body (structured fields when the CLI provides them).

**Example:**

```bash
nxcurl run '{{BASE}}/v1/items' -X GET -H 'Authorization: Bearer {{TOKEN}}' -e prod --json
```

---

## Command: `nxcurl send`

**Tool name:** `nxcurl` (subcommand `send`)

**Inputs:**

- Positional: `<request.json>` — JSON with `name`, `method`, `url`, `headers` (object), `body` (string)
- `--no-history` — skip writing this exchange to history
- Global: `-e` / `--env`, `--json`

**Triggers:**

- Run a **saved** request file (often under `~/.nxcurl/collections/` after `import har` / `import postman`, or hand-authored).
- Prefer over `run` when the request definition is fixed and you only swap environment or need JSON output.

**Example:**

```bash
nxcurl send ~/.nxcurl/collections/<collection_dir>/request.json -e prod --json
```

---

## Command: `nxcurl history list`

**Tool name:** `nxcurl` (subcommand `history list`)

**Inputs:**

- `--limit <n>` — max rows (default 50)
- Global: `--json`

**Triggers:**

- List **recent** exchanges (newest first) to obtain **ids** for `history show` / `history replay`.
- Agents should pass `--json` for a parseable list.

**Example:**

```bash
nxcurl history list --limit 20 --json
```

---

## Command: `nxcurl history show`

**Tool name:** `nxcurl` (subcommand `history show`)

**Inputs:**

- Positional: `<id>` — from `history list`
- Global: `--json`

**Triggers:**

- Load **one** stored exchange in full (request + response headers and bodies) for debugging or summarizing a past call.

**Example:**

```bash
nxcurl history show <id> --json
```

---

## Command: `nxcurl history replay`

**Tool name:** `nxcurl` (subcommand `history replay`)

**Inputs:**

- Positional: `<id>` — from `history list`
- Global: `-e` / `--env`, `--json` (same request shape; env file can change resolved `{{VAR}}` values)

**Triggers:**

- **Re-execute** a logged request; appends a **new** row to history (does not remove the original).
- Use when you want the exact same method/URL/headers/body as a prior run, optionally with a different named environment.

**Example:**

```bash
nxcurl history replay <id> -e staging --json
```

---

## Command: `nxcurl import har`

**Tool name:** `nxcurl` (subcommand `import har`)

**Inputs:**

- Positional: `<file.har>` — HTTP Archive 1.x

**Triggers:**

- Turn a browser HAR export into per-request JSON files under `~/.nxcurl/collections/` for `nxcurl send`.

**Example:**

```bash
nxcurl import har ./capture.har
```

---

## Command: `nxcurl import postman`

**Tool name:** `nxcurl` (subcommand `import postman`)

**Inputs:**

- Positional: `<collection.json>` — Postman Collection v2.0 / v2.1

**Triggers:**

- Import a Postman collection into `~/.nxcurl/collections/` as individual request JSON files.

**Example:**

```bash
nxcurl import postman ./API.postman_collection.json
```

---

## Command: `nxcurl env init`

**Tool name:** `nxcurl` (subcommand `env init`)

**Inputs:**

- Positional: `<name>` — creates `~/.nxcurl/envs/<name>.json` as `{}` if it does not exist (errors if it exists)

**Triggers:**

- Create a new named environment before editing key/value secrets for `{{KEY}}` substitution.

**Example:**

```bash
nxcurl env init prod
```

---

## Command: `nxcurl env path`

**Tool name:** `nxcurl` (subcommand `env path`)

**Inputs:**

- Positional: `<name>`

**Triggers:**

- Print the absolute path to `~/.nxcurl/envs/<name>.json` so tools or users can open or validate the file.

**Example:**

```bash
nxcurl env path prod
```

---

## Command: `nxcurl docs`

**Tool name:** `nxcurl` (subcommand `docs`)

**Inputs:** none

**Triggers:**

- Refresh `~/.nxcurl/docs/SKILL.md` and **print this skill** to stdout (agents can run once per session or after upgrades).

**Example:**

```bash
nxcurl docs
```

---

## Agent workflow tips

1. Prefer `nxcurl … --json` whenever the model must parse responses or history.
2. Store secrets in `~/.nxcurl/envs/<name>.json`; use `{{KEY}}` in URL, `-H`, and `-d`; pass `-e <name>`.
3. After `import`, discover request files under `~/.nxcurl/collections/` and run `nxcurl send <path> -e … --json`.
4. Run `nxcurl docs` after upgrading the binary so on-disk `SKILL.md` matches the build you are using.
