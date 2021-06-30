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
	Python     = "python"
	Ruby       = "ruby"
	Perl       = "perl"
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
