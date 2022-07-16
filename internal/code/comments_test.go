package code

import "testing"

func TestHidesComments(t *testing.T) {
	content := `
///package main
///
///import "fmt"
///
///func main() {
  fmt.Println("Hello, world!")
///}`

	expected := `
  fmt.Println("Hello, world!")`

	if HideComments(content) != expected {
		t.Errorf("Expected %s, got %s", expected, HideComments(content))
	}
}

func TestNoComments(t *testing.T) {
	content := `
package main

import "fmt"

func main() {
  fmt.Println("Hello, world!")
}`
	expected := content

	if HideComments(content) != expected {
		t.Errorf("Expected %s, got %s", expected, HideComments(content))
	}
	if RemoveComments(content) != expected {
		t.Errorf("Expected %s, got %s", expected, HideComments(content))
	}
}

func TestRemoveComments(t *testing.T) {
	content := `
///package main
///
///import "fmt"
///
///func main() {
  fmt.Println("Hello, world!")
///}`

	expected := `
package main

import "fmt"

func main() {
  fmt.Println("Hello, world!")
}`

	if RemoveComments(content) != expected {
		t.Errorf("Expected %s, got %s", expected, RemoveComments(content))
	}
}
