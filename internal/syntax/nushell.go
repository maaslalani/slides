package syntax

import (
	"regexp"
	"strings"
)

// NushellHighlighter provides basic syntax highlighting for Nushell code
// by converting it to a format that works better with existing shell lexers
type NushellHighlighter struct{}

// NewNushellHighlighter creates a new nushell highlighter
func NewNushellHighlighter() *NushellHighlighter {
	return &NushellHighlighter{}
}

// PreprocessForShellHighlighting converts nushell syntax to be more compatible
// with shell syntax highlighters like bash
func (n *NushellHighlighter) PreprocessForShellHighlighting(code string) string {
	// Apply various transformations to make nushell code more shell-like
	// for better syntax highlighting

	processed := code

	// Convert nushell pipes to regular pipes (already compatible)
	// No change needed for basic pipes

	// Convert nushell comments (# already compatible with shell)
	// No change needed

	// Convert nushell variables ($var -> $var, already compatible)
	// No change needed for basic variables

	// Convert nushell record syntax {key: value} to something more shell-like
	// This is a basic approximation for highlighting purposes
	recordRegex := regexp.MustCompile(`\{([^}]+)\}`)
	processed = recordRegex.ReplaceAllStringFunc(processed, func(match string) string {
		// Keep the braces but make it look more like shell arrays
		return strings.ReplaceAll(match, ":", "=")
	})

	// Convert nushell list syntax [item1, item2] to shell-like arrays
	listRegex := regexp.MustCompile(`\[([^\]]+)\]`)
	processed = listRegex.ReplaceAllString(processed, "($1)")

	// Convert nushell where clauses to look more like shell conditionals
	whereRegex := regexp.MustCompile(`\|\s*where\s+`)
	processed = whereRegex.ReplaceAllString(processed, "| grep ")

	// Convert nushell $it references to shell-like variables
	processed = strings.ReplaceAll(processed, "$it", "$_")

	return processed
}

// GetFallbackLanguage returns the best fallback language for nushell highlighting
func (n *NushellHighlighter) GetFallbackLanguage() string {
	return "bash"
}

// ShouldUsePreprocessing determines if the code should be preprocessed
// based on common nushell patterns
func (n *NushellHighlighter) ShouldUsePreprocessing(code string) bool {
	nushellPatterns := []string{
		"{",           // Records
		"where $it",   // Nushell where clauses
		"| get ",      // Nushell get command
		"| select ",   // Nushell select command
		"| each ",     // Nushell each command
		"| from ",     // Nushell from command
		"| to ",       // Nushell to command
		"| str ",      // Nushell str command
		"| math ",     // Nushell math command
		"| length",    // Nushell length command
		"| first",     // Nushell first command
		"| last",      // Nushell last command
		"| sort-by",   // Nushell sort-by command
		"| group-by",  // Nushell group-by command
		"| where ",    // Nushell where command
		"| into ",     // Nushell into command
		"| append",    // Nushell append command
		"| prepend",   // Nushell prepend command
		"open ",       // Nushell open command
		"save ",       // Nushell save command
		"sys",         // Nushell sys command
		"date now",    // Nushell date command
		"http get",    // Nushell http command
		"format date", // Nushell format command
	}

	for _, pattern := range nushellPatterns {
		if strings.Contains(code, pattern) {
			return true
		}
	}

	return false
}

// LanguageMapping defines fallback mappings for nushell identifiers
var LanguageMapping = map[string]string{
	"nu": "bash",
}

// GetMappedLanguage returns the mapped language for highlighting,
// or the original if no mapping exists
func GetMappedLanguage(lang string) string {
	if mapped, exists := LanguageMapping[lang]; exists {
		return mapped
	}
	return lang
}

// ProcessNushellCode applies preprocessing if needed and returns the language to use
func ProcessNushellCode(code, language string) (processedCode, highlightLanguage string) {
	if language == "nu" {
		highlighter := NewNushellHighlighter()

		if highlighter.ShouldUsePreprocessing(code) {
			return highlighter.PreprocessForShellHighlighting(code), highlighter.GetFallbackLanguage()
		}

		return code, highlighter.GetFallbackLanguage()
	}

	return code, language
}

// PreprocessNushellMarkdown preprocesses markdown to replace nushell language identifiers
// with bash for better syntax highlighting in code blocks
func PreprocessNushellMarkdown(markdown string) string {
	// Regular expression to match code blocks with nushell language identifiers
	// Matches both ```nu and ```nushell
	nushellCodeBlockRegex := regexp.MustCompile(`(?m)^` + "```(?:nu)" + `$`)

	// Replace nushell identifiers with bash for better syntax highlighting
	processed := nushellCodeBlockRegex.ReplaceAllString(markdown, "```bash")

	return processed
}
