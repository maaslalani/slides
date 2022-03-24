# Welcome to Slides
A terminal based presentation tool

## Everything is markdown
In fact this entire presentation is a markdown file.

Press `n` to go to the next slide.

---

# Display Code

```go
package main

import "fmt"

func main() {
  // You can show code in slides
  // Press ctrl+e to execute this code directly in slides
  fmt.Println("Tada!")
}
```

---

# h1

You can use everything in markdown!
* Like bulleted list
* You know the deal

1. Numbered lists too

## h2

| Tables | Too    |
| ------ | ------ |
| Even   | Tables |


### h3

#### h4
##### h5
###### h6

---

# Graphs

```
digraph {
    rankdir = LR;
    a -> b;
    b -> c;
}
```
```
┌───┐     ┌───┐     ┌───┐
│ a │ ──▶ │ b │ ──▶ │ c │
└───┘     └───┘     └───┘
```
---

All you need to do is separate slides with triple dashes
`---` on a separate line, like so:

```markdown
# Slide 1
Some stuff

--- 

# Slide 2
Some other stuff
```
