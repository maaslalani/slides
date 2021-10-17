package code

// base: A filename is appended to the end of the arguments.
// Example: base{"py"} -> `py test.py`
// Placeholders <file>, <name> and <path> can be used.
type base []string

// single: Same as base, except that the filename is not appended.
type single []string

// multi: Multiple commands; placeholders can be used
type multi [][]string

// ----

type Language struct {
	Extension string
	// Commands  [][]string // placeholders: <name> file name (without extension),
	// <file> file name, <path> path without file name
	Cmds interface{} // type of multi, single, base
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
		Cmds:      base{"bash"},
	},
	Elixir: {
		Extension: "exs",
		Cmds:      base{"elixir"},
	},
	Go: {
		Extension: "go",
		Cmds:      base{"go", "run"},
	},
	Javascript: {
		Extension: "js",
		Cmds:      base{"node"},
	},
	Lua: {
		Extension: "lua",
		Cmds:      base{"lua"},
	},
	Ruby: {
		Extension: "rb",
		Cmds:      base{"ruby"},
	},
	Python: {
		Extension: "py",
		Cmds:      base{"python"},
	},
	Perl: {
		Extension: "pl",
		Cmds:      base{"perl"},
	},
	Rust: {
		Extension: "rs",
		Cmds: multi{
			// compile code
			{"rustc", "<file>", "-o", "<path>/<name>.run"},
			// grant execute permissions
			{"chmod", "+x", "<path>/<name>.run"},
			// run compiled file
			{"<path>/<name>.run"},
		},
	},
}
