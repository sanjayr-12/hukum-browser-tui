# AGENTS.md

Guidelines and architecture rules for developing Hukum.

## Overview
Hukum is a keyboard-driven terminal web browser written in Go using Bubble Tea, Lip Gloss, Glamour, and the TinyFish HTTP APIs (Search and Fetch). It provides a fast, lightweight way to search the web and read modern JavaScript-rendered websites as clean Markdown in the terminal.

---

## Codebase Structure
The codebase follows a decoupled layout with tests separated from application packages:

```
hukum/
├── .env                  # Local environment variables (API keys, git-ignored)
├── .gitignore            # Ignored files (binaries, env, build artifacts)
├── AGENTS.md             # Development standards and guidelines for agents
├── README.md             # Simple, casual user guide
├── go.mod / go.sum       # Go module and dependency definitions
├── main.go               # Application entrypoint
├── internal/             # Application source packages
│   ├── config/           # Configuration and .env loader
│   │   └── config.go
│   ├── tinyfish/         # TinyFish API HTTP client (Search and Fetch)
│   │   └── client.go
│   └── ui/               # Bubble Tea TUI components and styles
│       ├── links.go      # Markdown AST and link extractor
│       ├── modal.go      # Floating prompt and list picker
│       ├── model.go      # Root browser state, key bindings, Glamour rendering
│       └── styles.go     # Lip Gloss styles and layout helpers
└── test/                 # Dedicated tests directory (never mix with source)
    ├── client_test.go    # TinyFish API integration tests
    ├── links_test.go     # Link extractor and URL parser unit tests
    └── render_test.go    # Glamour terminal rendering tests
```

---

## Architectural Rules

1. **Separation of Tests and Application Code**:
   - Never mix test files (`*_test.go`) into `internal/` packages.
   - All unit, integration, and rendering tests belong in the top-level `test/` directory.

2. **Decoupled Architecture**:
   - `internal/config`: Handles flag parsing and `.env` loading.
   - `internal/tinyfish`: Pure Go client using standard `net/http`. Never depends on `ui` or Bubble Tea.
   - `internal/ui`: Manages presentation and interaction via Bubble Tea's Model-Update-View lifecycle.
   - Root `main.go`: Assembles dependencies and launches the Bubble Tea runtime.

3. **Zero External Runtime Dependencies**:
   - Do not shell out to Node.js, python, or external CLI tools.
   - All network operations must go directly through Go's standard library `net/http` to TinyFish endpoints.

4. **No Emojis**:
   - Do not use emojis anywhere in the codebase: not in UI labels, not in error messages, not in logs, not in comments, and not in markdown documentation.
   - Use clean ASCII or simple symbols (e.g. `>` for selection cursor, `[Web]`, `-`, `|`).

5. **Code Comments**:
   - Do not write obvious or repetitive comments that describe what the code already says.
   - Only add comments when explaining non-obvious business logic, workarounds, or protocol quirks.

6. **UI & Styling Consistency**:
   - Use a full-page terminal layout spanning 100% width and height (never shrink or center dialogs with empty margins).
   - Theme: Authentic Atom One Dark palette:
     - Background: `#282C34` (canvas), `#21252B` (header/status bars), `#2C313A` (search container)
     - Selection Highlight: `#3E4451`
     - Foreground/Text: `#ECEFF4` (bright/headings), `#ABB2BF` (body), `#5C6370` (muted/metadata)
     - Accents: `#61AFEF` (blue/links), `#56B6C2` (cyan), `#98C379` (green), `#C678DD` (purple badge), `#E06C75` (error red)
   - Markdown rendering must wrap to terminal width and re-render on terminal resize events (`tea.WindowSizeMsg`).

7. **Git Restrictions**:
   - Agents are strictly prohibited from running git write or execution commands (such as `git init`, `git add`, `git commit`, `git push`, `git checkout -b`).
   - Git may only be used for read-only inspection (`git status`, `git diff`, `git log`) if explicitly needed.

8. **Documentation Tone**:
   - Documentation must be casual, direct, and straightforward.
   - Avoid buzzwords, corporate marketing jargon, artificial excitement, and emoji decorations.
