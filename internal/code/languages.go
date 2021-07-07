package code

type Language struct {
	Extension string
	Command   []string
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
)

var Languages = map[string]Language{
	Bash: {
		Extension: "sh",
		Command:   []string{"bash"},
	},
	Elixir: {
		Extension: "exs",
		Command:   []string{"elixir"},
	},
	Go: {
		Extension: "go",
		Command:   []string{"go", "run"},
	},
	Javascript: {
		Extension: "js",
		Command:   []string{"node"},
	},
	Lua: {
		Extension: "lua",
		Command:   []string{"lua"},
	},
	Ruby: {
		Extension: "rb",
		Command:   []string{"ruby"},
	},
	Python: {
		Extension: "py",
		Command:   []string{"python"},
	},
	Perl: {
		Extension: "pl",
		Command:   []string{"perl"},
	},
}
