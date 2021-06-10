# Slides

Slides in your terminal.

##### Credits
This project is heavily inspired by [`lookatme`](https://github.com/d0c-s4vage/lookatme).
`slides` is a more minimal version of [`lookatme`](https://github.com/d0c-s4vage/lookatme) and written in Go.

### Demo
![slides](../assets/demo.gif?raw=true)

### Installation
```
go get github.com/maaslalani/slides
```

### Usage
Create a simple markdown file that contains your slides:

```
# Welcome to Slides
A terminal based presentation tool

---

## Everything is markdown
In fact, this entire presentation is a markdown file.

---

## Everything happens in your terminal
Create slides and present them without ever leaving your terminal.

```

Checkout the [example slides](./examples).

Then, to present, run:
```
slides presentation.md
```

You are also able to pass in slides through `stdin`, this allows you to `curl` and present remote files:
```
curl https://example.com/slides.md | slides
```

Go to the next slide with any of the following keys:
* <kbd>space</kbd>
* <kbd>right</kbd>
* <kbd>down</kbd>
* <kbd>enter</kbd>
* <kbd>n</kbd>
* <kbd>k</kbd>
* <kbd>l</kbd>

Go to the previous slide with any of the following keys:
* <kbd>left</kbd>
* <kbd>up</kbd>
* <kbd>p</kbd>
* <kbd>h</kbd>
* <kbd>j</kbd>

### Development
```
make
```
