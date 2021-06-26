package process

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	md := `
# Slide

~~~sd Replace Process
Replace
~~~

Hello

~~~sd Replace Process
Replace
Multi-line input
~~~

~~~echo -n World
Hello
~~~

---

# Next Slide

GraphViz Test

~~~graph-easy --as=boxart
digraph {
  A -> B
}
~~~
`

	got := Parse(md)
	want := []Block{{
		Command: "sd Replace Process",
		Input:   "Replace",
		Raw:     "~~~sd Replace Process\nReplace\n~~~",
	}, {
		Command: "sd Replace Process",
		Input:   "Replace\nMulti-line input",
		Raw:     "~~~sd Replace Process\nReplace\nMulti-line input\n~~~",
	}, {
		Command: "echo -n World",
		Input:   "Hello",
		Raw:     "~~~echo -n World\nHello\n~~~",
	}, {
		Command: "graph-easy --as=boxart",
		Input:   "digraph {\n  A -> B\n}",
		Raw:     "~~~graph-easy --as=boxart\ndigraph {\n  A -> B\n}\n~~~",
	}}

	if !reflect.DeepEqual(got, want) {
		t.Log(want)
		t.Log(got)
		t.Fatal("Did not parse blocks correctly")
	}
}
