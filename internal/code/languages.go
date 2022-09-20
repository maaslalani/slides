package code

// cmds: Multiple commands; placeholders can be used
// Placeholders <file>, <name> and <path> can be used.
type cmds [][]string

type reReplacement struct {
	src  string
	repl string
}

// Language represents a programming language with it Extension and Commands to
// execute its programs.
type Language struct {
	// Extension represents the file extension used by this language.
	Extension string
	// Commands  [][]string // placeholders: <name> file name (without
	// extension), <file> file name, <path> path without file name
	Commands           cmds
	OutputsReplacement reReplacement
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
	Java       = "java"
)

// Languages is a map of supported languages with their extensions and commands
// to run to execute the program.
var Languages = map[string]Language{
	Bash: {
		Extension: "sh",
		Commands:  cmds{{"bash", "<file>"}},
		OutputsReplacement: reReplacement{
			src: `<file>: line [0-9]+: `,
		},
	},
	Elixir: {
		Extension: "exs",
		Commands:  cmds{{"elixir", "<file>"}},
	},
	Go: {
		Extension: "go",
		Commands:  cmds{{"go", "run", "<file>"}},
		OutputsReplacement: reReplacement{
			src: `<file>:([0-9]+:){2} `,
		},
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
			// run compiled file
			{"<path>/<name>.run"},
		},
	},
	Java: {
		Extension: "java",
		Commands:  cmds{{"java", "<file>"}},
	},
}
