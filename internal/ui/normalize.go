package ui

import (
	"regexp"
	"strings"
)

var (
	// Ad & Sponsored Content Regexes
	reAdBlock = regexp.MustCompile(
		`(?im)(?:^[a-z0-9_.-]+(?:\s*[•·*|]|\s+promoted)?\s*\n+)?` +
			`^(?:ad|sponsored|promoted|advertisement)\b[^\n]*\n+` +
			`(?:[^\n]+\n+){1,4}` +
			`^(?:learn more|shop now|sign up|visit site|get offer|install now|download|apply now)\b[^\n]*\n*` +
			`(?:https?://[^\n]+|[a-z0-9-]+\.[a-z]{2,}[^\n]*\n*)?`,
	)

	reCookieBanner = regexp.MustCompile(
		`(?im)^.*(?:we use cookies|cookie policy|by continuing to browse this site|accept all cookies|manage cookie preferences).*$\n*`,
	)

	reSocialShare = regexp.MustCompile(
		`(?im)^\s*(?:share on (?:facebook|twitter|x|linkedin|reddit)|tweet this|email (?:this )?article)\s*$\n*`,
	)

	reTopNavDump = regexp.MustCompile(
		`(?s)^\s*(?:\*\s+[^\n]+\n)+\s*(#\s+)`,
	)

	// Code Block & Formatting Regexes
	reAsterisks = regexp.MustCompile(`\*{3,4}([^*]+)\*{3,4}`)

	// Language on previous line attached to code fence
	reLangCodeWithLine = regexp.MustCompile(
		`(?im)^([a-zA-Z0-9_#+]+)\s*\n+[` + "`" + `\s]{3,}\s*([^\n]+)$`,
	)
	reLangCode = regexp.MustCompile(
		`(?im)^([a-zA-Z0-9_#+]+)\s*\n+[` + "`" + `\s]{3,}\n`,
	)

	// Inline code block inside backticks
	reInlineCodeBlock = regexp.MustCompile(
		`(?m)^[` + "`" + `\s]{3,}\s*([^` + "`" + `\n]+)\s*[` + "`" + `\s]{3,}$`,
	)

	// Code fence with trailing code/backticks
	reTrailingFence = regexp.MustCompile(
		`(?m)\}\s*[` + "`" + `\s]{3,}$`,
	)

	// Any line containing strictly 3 or more backticks and whitespace
	reOnlyBackticks = regexp.MustCompile(
		`(?m)^[` + "`" + `\s]{3,}$`,
	)

	// Table & Feed (e.g. Hacker News) Regexes
	reIsolatedDomain = regexp.MustCompile(
		`\n\(\s*\n([^\n]+)\n\s*\)`,
	)

	reNumberedItem = regexp.MustCompile(
		`(?m)^(\d+)\.\s*\n+([^\n]+)`,
	)

	reHNMetadata = regexp.MustCompile(
		`(\d+\s+points)\s*\n+by\s*\n+([^\n]+)\s*\n+([^\n]+ago)\s*\n+\|\s*\n+hide\s*\n+\|\s*\n+([^\n]+comments?)`,
	)

	reHNNoScore = regexp.MustCompile(
		`by\s*\n+([^\n]+)\s*\n+([^\n]+ago)\s*\n+\|\s*\n+hide\s*\n+\|\s*\n+([^\n]+comments?)`,
	)

	reHeaderSpacing = regexp.MustCompile(
		`(?m)([^\n])\n(#{1,6}\s+)`,
	)
)

// NormalizeMarkdown transforms raw scraped web markdown into clean, well-structured
// Markdown: strips ads and cookie clutter, repairs broken code blocks, triggers Chroma
// syntax highlighting, and formats table-extracted feeds cleanly.
func NormalizeMarkdown(raw string) string {
	if raw == "" {
		return ""
	}

	text := raw

	// 1. Strip standalone ad blocks and promotional spam
	text = reAdBlock.ReplaceAllString(text, "")
	text = reCookieBanner.ReplaceAllString(text, "")
	text = reSocialShare.ReplaceAllString(text, "")
	text = reTopNavDump.ReplaceAllString(text, "$1")

	// 2. Normalize asterisks: ****Text**** -> **Text**
	text = reAsterisks.ReplaceAllString(text, "**$1**")

	// 3. Attach language label on preceding line to the code fence (e.g. Java\n``` -> ```java\n)
	text = reLangCodeWithLine.ReplaceAllString(text, "```$1\n$2")
	text = reLangCode.ReplaceAllString(text, "```$1\n")

	// 4. Fix code on the same line as fence
	text = reInlineCodeBlock.ReplaceAllString(text, "```\n$1\n```")

	// 5. Fix closing braces with attached backticks
	text = reTrailingFence.ReplaceAllString(text, "}\n```")

	// 6. Normalize any lines containing strictly backticks and whitespace to clean triple backticks
	text = reOnlyBackticks.ReplaceAllString(text, "```")

	// 7. Table and feed cleanups (Hacker News, forums)
	text = reIsolatedDomain.ReplaceAllString(text, " ($1)\n")
	text = reNumberedItem.ReplaceAllString(text, "\n\n$1. **$2**\n")
	text = reHNMetadata.ReplaceAllString(text, "   *$1 by $2 $3 | $4*\n")
	text = reHNNoScore.ReplaceAllString(text, "   *by $1 $2 | $3*\n")

	// 8. Ensure markdown headers have a blank line before them
	text = reHeaderSpacing.ReplaceAllString(text, "$1\n\n$2")

	// 9. Normalize excessive blank lines
	for strings.Contains(text, "\n\n\n") {
		text = strings.ReplaceAll(text, "\n\n\n", "\n\n")
	}

	return strings.TrimSpace(text)
}
