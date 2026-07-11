package syntax

import (
	"testing"
)

func TestPreprocessNushellMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Replace nu language identifier",
			input: `# Test Slide

` + "```nu" + `
echo "Hello from Nushell!"
` + "```" + `

Some other content.`,
			expected: `# Test Slide

` + "```bash" + `
echo "Hello from Nushell!"
` + "```" + `

Some other content.`,
		},
		{
			name: "Multiple nu code blocks",
			input: `# Slide 1

` + "```nu" + `
echo "First block"
` + "```" + `

---

# Slide 2

` + "```nu" + `
echo "Second block"
` + "```" + ``,
			expected: `# Slide 1

` + "```bash" + `
echo "First block"
` + "```" + `

---

# Slide 2

` + "```bash" + `
echo "Second block"
` + "```" + ``,
		},
		{
			name: "Leave other languages unchanged",
			input: `# Mixed Languages

` + "```python" + `
print("Hello Python")
` + "```" + `

` + "```nu" + `
echo "Hello Nushell"
` + "```" + `

` + "```go" + `
fmt.Println("Hello Go")
` + "```" + ``,
			expected: `# Mixed Languages

` + "```python" + `
print("Hello Python")
` + "```" + `

` + "```bash" + `
echo "Hello Nushell"
` + "```" + `

` + "```go" + `
fmt.Println("Hello Go")
` + "```" + ``,
		},
		{
			name:     "No code blocks",
			input:    "# Just a title\n\nSome regular markdown content.",
			expected: "# Just a title\n\nSome regular markdown content.",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreprocessNushellMarkdown(tt.input)
			if result != tt.expected {
				t.Errorf("PreprocessNushellMarkdown() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestNushellHighlighterPreprocessForShellHighlighting(t *testing.T) {
	highlighter := NewNushellHighlighter()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Convert records to shell-like syntax",
			input:    `{name: "Alice", age: 30}`,
			expected: `{name= "Alice", age= 30}`,
		},
		{
			name:     "Convert lists to shell arrays",
			input:    `[1, 2, 3, 4, 5]`,
			expected: `(1, 2, 3, 4, 5)`,
		},
		{
			name:     "Convert where clauses",
			input:    `ls | where size > 100`,
			expected: `ls | grep size > 100`,
		},
		{
			name:     "Convert $it references",
			input:    `[1, 2, 3] | where $it > 1`,
			expected: `(1, 2, 3) | grep $_ > 1`,
		},
		{
			name: "Complex nushell pipeline",
			input: `{users: [{name: "John", age: 25}, {name: "Jane", age: 30}]}
| get users
| where $it.age > 26`,
			expected: `{users= ({name= "John", age= 25}, {name= "Jane", age= 30})}
| get users
| grep $_.age > 26`,
		},
		{
			name:     "No nushell syntax",
			input:    `echo "Hello World"`,
			expected: `echo "Hello World"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := highlighter.PreprocessForShellHighlighting(tt.input)
			if result != tt.expected {
				t.Errorf("PreprocessForShellHighlighting() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestShouldUsePreprocessing(t *testing.T) {
	highlighter := NewNushellHighlighter()

	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "Contains nushell record syntax",
			code:     `{name: "Alice"}`,
			expected: true,
		},
		{
			name:     "Contains where $it",
			code:     `[1, 2, 3] | where $it > 1`,
			expected: true,
		},
		{
			name:     "Contains nushell get command",
			code:     `data | get name`,
			expected: true,
		},
		{
			name:     "Contains nushell select command",
			code:     `data | select name age`,
			expected: true,
		},
		{
			name:     "Contains nushell math command",
			code:     `[1, 2, 3] | math sum`,
			expected: true,
		},
		{
			name:     "Contains nushell date command",
			code:     `date now`,
			expected: true,
		},
		{
			name:     "Contains nushell http command",
			code:     `http get https://api.example.com`,
			expected: true,
		},
		{
			name:     "Contains sys command",
			code:     `sys`,
			expected: true,
		},
		{
			name:     "Basic shell command",
			code:     `echo "Hello World"`,
			expected: false,
		},
		{
			name:     "Regular bash syntax",
			code:     `ls -la | grep ".txt"`,
			expected: false,
		},
		{
			name:     "Empty code",
			code:     ``,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := highlighter.ShouldUsePreprocessing(tt.code)
			if result != tt.expected {
				t.Errorf("ShouldUsePreprocessing() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetMappedLanguage(t *testing.T) {
	tests := []struct {
		name     string
		language string
		expected string
	}{
		{
			name:     "Map nu (alternative test)",
			language: "nu",
			expected: "bash",
		},
		{
			name:     "Leave other languages unchanged",
			language: "python",
			expected: "python",
		},
		{
			name:     "Leave go unchanged",
			language: "go",
			expected: "go",
		},
		{
			name:     "Leave empty string unchanged",
			language: "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMappedLanguage(tt.language)
			if result != tt.expected {
				t.Errorf("GetMappedLanguage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestProcessNushellCode(t *testing.T) {
	tests := []struct {
		name             string
		code             string
		language         string
		expectedCode     string
		expectedLanguage string
	}{
		{
			name:             "Process nushell code with nu identifier",
			code:             `{name: "Alice"} | get name`,
			language:         "nu",
			expectedCode:     `{name= "Alice"} | get name`,
			expectedLanguage: "bash",
		},
		{
			name:             "Simple nushell code without preprocessing",
			code:             `echo "Hello World"`,
			language:         "nu",
			expectedCode:     `echo "Hello World"`,
			expectedLanguage: "bash",
		},
		{
			name:             "Non-nushell language unchanged",
			code:             `print("Hello Python")`,
			language:         "python",
			expectedCode:     `print("Hello Python")`,
			expectedLanguage: "python",
		},
		{
			name:             "Go code unchanged",
			code:             `fmt.Println("Hello Go")`,
			language:         "go",
			expectedCode:     `fmt.Println("Hello Go")`,
			expectedLanguage: "go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processedCode, highlightLanguage := ProcessNushellCode(tt.code, tt.language)

			if processedCode != tt.expectedCode {
				t.Errorf("ProcessNushellCode() code = %q, want %q", processedCode, tt.expectedCode)
			}

			if highlightLanguage != tt.expectedLanguage {
				t.Errorf("ProcessNushellCode() language = %q, want %q", highlightLanguage, tt.expectedLanguage)
			}
		})
	}
}

func TestGetFallbackLanguage(t *testing.T) {
	highlighter := NewNushellHighlighter()

	expected := "bash"
	result := highlighter.GetFallbackLanguage()

	if result != expected {
		t.Errorf("GetFallbackLanguage() = %q, want %q", result, expected)
	}
}
