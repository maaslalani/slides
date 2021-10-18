package code

// cmds: Multiple commands; placeholders can be used
// Placeholders <file>, <name> and <path> can be used.
type cmds [][]string

// ----

type Language struct {
	Extension string
	// Commands  [][]string // placeholders: <name> file name (without extension),
	// <file> file name, <path> path without file name
	Commands cmds
}

// Supported Languages
const (
	Bash       = "bash"
	Elixir     = "elixir"
	Go         = "go"
	Javascript = "javascript"
	Lua        = "lua"
	Perl       = "perl"
	Python     = "python"
	Ruby       = "ruby"
	Rust       = "rust"
)

var Languages = map[string]Language{
	Bash: {
		Extension: "sh",
		Commands:  cmds{{"bash", "<file>"}},
	},
	Elixir: {
		Extension: "exs",
		Commands:  cmds{{"elixir", "<file>"}},
	},
	Go: {
		Extension: "go",
		Commands:  cmds{{"go", "run", "<file>"}},
	},
	Javascript: {
		Extension: "js",
		Commands:  cmds{{"node", "<file>"}},
	},
	Lua: {
		Extension: "lua",
		Commands:  cmds{{"lua", "<file>"}},
	},
	Ruby: {
		Extension: "rb",
		Commands:  cmds{{"ruby", "<file>"}},
	},
	Python: {
		Extension: "py",
		Commands:  cmds{{"python", "<file>"}},
	},
	Perl: {
		Extension: "pl",
		Commands:  cmds{{"perl", "<file>"}},
	},
	Rust: {
		Extension: "rs",
		Commands: cmds{
			// compile code
			{"rustc", "<file>", "-o", "<path>/<name>.run"},
			// grant execute permissions
			{"chmod", "+x", "<path>/<name>.run"},
			// run compiled file
			{"<path>/<name>.run"},
		},
	},
}
