# ZenTxt — Context for Claude Code

## Your mission

You are only allowed to work on the **frontend**. Do not touch any backend files, migrations, or Go logic. The backend is built and understood by the developer — modifying it would create a comprehension debt they want to avoid.

Frontend files live in `static/` and Go template files in `templates/` (to be created).

---

## What is ZenTxt

ZenTxt is a stupidly simple tool for versioning text documents. No commands, no learning curve. Inspired by Steve Jobs' philosophy of simplicity and Apple's minimalist design.

**Core philosophy:**
- One thing, done well
- No jargon (say "Restore" not "Revert", "Save version" not "Commit")
- Feels calm, not technical
- Local-first — no cloud, no account

---

## Project architecture

```
zentxt/
├── cmd/zentxt/main.go          ← Entry point, starts Fiber server
├── internal/
│   ├── api/routes.go           ← All HTTP routes (DO NOT MODIFY)
│   ├── storage/db.go           ← SQLite connection (DO NOT MODIFY)
│   └── versioning/             ← Business logic (DO NOT MODIFY)
│       ├── file.go
│       ├── version.go
│       └── diff.go
├── migrations/                 ← SQL migrations (DO NOT MODIFY)
├── static/                     ← Frontend assets (CSS, JS) ← YOUR ZONE
├── templates/                  ← Go HTML templates ← YOUR ZONE
├── Dockerfile
└── go.mod
```

---

## Tech stack — Frontend only

- **HTMX** — HTTP requests from HTML attributes, no JavaScript framework
- **Tailwind CSS** — utility-first CSS, via CDN for simplicity
- **Go templates** — server-side HTML rendering via Fiber

Load both via CDN in the base template:
```html
<script src="https://unpkg.com/htmx.org@1.9.10"></script>
<link href="https://cdn.tailwindcss.com" rel="stylesheet">
```

---

## API endpoints (read-only reference)

All endpoints are already implemented. The frontend consumes them.

| Method | Endpoint | Description | Request body |
|--------|----------|-------------|--------------|
| `GET` | `/files` | List all tracked files | — |
| `POST` | `/files` | Track a new file | `{"name": "string", "path": "string"}` |
| `GET` | `/files/:id` | Get a file | — |
| `POST` | `/files/:id/versions` | Save a new version | `{"path": "string", "author": "string", "message": "string", "content": "string"}` |
| `GET` | `/files/:id/versions` | List all versions | — |
| `GET` | `/files/:id/versions/:version_id` | Get a specific version | — |
| `GET` | `/files/:id/diff?v1=&v2=` | Compare two versions | — |
| `PATCH` | `/files/:id/restore/:version_id` | Restore a version | — |

---

## Screens to build

### 1. Home — List of tracked files
- Display all tracked files from `GET /files`
- Each file shows: name, path, creation date
- Click on a file → goes to its timeline
- Drag & drop zone to add a new file (calls `POST /files`)

### 2. File timeline
- Display all versions of a file from `GET /files/:id/versions`
- Each version shows: version number, author, message, date
- Button "Restore" on each version → calls `PATCH /files/:id/restore/:version_id`
- Button "Compare" between two versions → goes to diff view

### 3. Diff view
- Side-by-side or inline diff between two versions
- Green = added, red = removed, gray = unchanged
- Data from `GET /files/:id/diff?v1=&v2=`
- Human language only — no technical jargon

---

## Design guidelines

- **Minimal** — lots of whitespace, no visual noise
- **Calm** — muted colors, no aggressive contrasts
- **Human** — soft typography, friendly copy
- Inspiration: Obsidian, Linear, Apple Notes
- Color palette suggestion: white background, `gray-900` text, `indigo-500` accents
- Font suggestion: Inter or system-ui

---

## How to add a new route in Fiber for a template

If you need to add a route that renders an HTML template, add it in `internal/api/routes.go` — but ask the developer first before modifying any backend file.

The developer prefers to add backend routes themselves. Provide the template and the route code separately so they can integrate it manually.

---

## What good looks like

ZenTxt should feel like the love child of Obsidian and Apple Notes. When someone opens it for the first time, they should think "oh, this is nice" — not "where do I start?".

Every screen should be explainable in one sentence.
