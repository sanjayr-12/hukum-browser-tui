# Hukum

A terminal browser written in Go. It uses TinyFish APIs to search the web and render modern, JS-heavy websites as clean Markdown right in your terminal.

## What it does

- Full-page search prompt (like Telescope or FZF in Neovim)
- Omnibox: search queries or open URLs directly
- Follow links on any page by pressing `f`
- Clean markdown rendering with syntax highlighting and text wrap
- History navigation (back, reload, scroll)

## Setup

Put your TinyFish API key in a `.env` file in this directory:

```bash
TINYFISH_API_KEY="your_api_key_here"
```

Or export it in your shell:

```bash
export TINYFISH_API_KEY="your_api_key_here"
```

## Running

Build and run:

```bash
go build -o hukum .
./hukum
```

Open a specific site:

```bash
./hukum https://news.ycombinator.com
```

Or start with a search query:

```bash
./hukum -search "golang bubble tea"
```

## Shortcuts

### Browser view
- `Ctrl+O` or `Ctrl+G` or `/` or `o`: Open search / URL prompt
- `f`: Open link picker for current page
- `b`: Go back
- `r`: Reload current page
- `j` / `k` or Arrow keys: Scroll down / up
- `d` / `u` or PageDown / PageUp: Scroll half page
- `g` / `G`: Top / bottom
- `q` or `Ctrl+C`: Quit

### Search & link picker
- `Enter`: Run search, open URL, or open selected link
- `Ctrl+N` / `Ctrl+P` or `Tab` / `Shift+Tab` or `Down` / `Up`: Navigate items
- `Esc`: Close prompt and go back to reading
