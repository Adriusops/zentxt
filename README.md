<p align="center">
  <img src="assets/Logo_Banner_Icon_white.png" width="700" alt="ZenTxt logo">
</p>

zentxt is a stupidly simple tool for versioning text documents — no commands, no learning curve. Inspired by Steve Jobs' philosophy of simplicity and Apple's minimalist design, zentxt lets anyone track file changes effortlessly.

No branches. No jargon. Just your files, with a memory.

---

## Getting started

Download the latest release for your platform → [GitHub Releases](https://github.com/Adriusops/zentxt/releases)

Double-click to open. That's it.

> **macOS users** — if macOS blocks the app on first launch, go to System Settings → Privacy & Security → "Open Anyway".
> - Or run this command in Terminal: `xattr -cr /path/to/zentxt.app`

> **Windows users** — if SmartScreen warns you, click "More info" → "Run anyway".

---

## For developers

```bash
git clone https://github.com/Adriusops/zentxt
cd zentxt
wails3 dev
```

## Stack

**Backend** — Go with Wails v3. A single native binary, no runtime dependencies.

**Database** — SQLite via modernc.org/sqlite. Local-first by design — your data stays on your machine, in a single file.

**Frontend** — Vanilla JS + Vite. Fast, minimal, no framework overhead.

**Diff engine** — go-diff. Computes character-level diffs between versions on the fly, no storage overhead.

**CI/CD** — GitHub Actions. Every push to `main` compiles the project and runs the test suite automatically. Every tag ships cross-platform binaries for Mac, Windows and Linux.

No cloud. No configuration. No account required.

---

## Roadmap

### Diff management
Currently, ZenTxt stores every version of your file — but choosing what to keep is still manual.
The next step is to let you **review changes between two versions line by line**, accept or discard specific edits, and merge them into a new version — similar to resolving a conflict in Git.

### Rich text format support
ZenTxt currently supports `.txt` and `.md` files.
Planned support for the most common text-based formats:

- `.docx` — extract plain text from the OOXML structure (ZIP + XML), diff the content, ignore formatting metadata
- `.odt` — same approach, LibreOffice's open format
- `.rtf` — strip RTF tags, diff the underlying text
- `.pdf` — on the roadmap but more complex: PDF is a rendering format, not an editing one. Text extraction is doable but imperfect; scanned PDFs would require OCR.

---

## Philosophy

zentxt does one thing well. It remembers your text files so you don't have to.

> "Simplicity is the ultimate sophistication." — Steve Jobs

---

## A personal note

zentxt is a solo side project, built with love and a genuine belief that software can be both useful and simple. It's also how I'm learning Go and DevOps in the open — one commit at a time.

If you believe in what this is trying to be, feel free to star, fork, or just say hi. 👋
