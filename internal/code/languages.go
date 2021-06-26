package code

type Language struct {
	Extension string
	Command   []string
}

// Supported Languages
const (
	Bash   = "bash"
	Go     = "go"
	Ruby   = "ruby"
	Python = "python"
	Elixir = "elixir"
)

var Languages = map[string]Language{
	Bash: {
		Extension: "sh",
		Command:   []string{"bash"},
	},
	Go: {
		Extension: "go",
		Command:   []string{"go", "run"},
	},
	Ruby: {
		Extension: "rb",
		Command:   []string{"ruby"},
	},
	Python: {
		Extension: "py",
		Command:   []string{"python"},
	},
	Elixir: {
		Extension: "exs",
		Command:   []string{"elixir"},
	},
}
