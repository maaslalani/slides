# Slides

[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://formulae.brew.sh/api/formula/slides.json&query=$.versions.stable&label=homebrew)](https://formulae.brew.sh/formula/slides)

Slides in your terminal.

<p align="center">
  <img src="./assets/demo.gif?raw=true" alt="Slides Presentation" />
</p>

### Installation
Using a package manager:

* Go
```bash
go install github.com/maaslalani/slides
```

* MacOS
```
brew install slides
```

* Arch
```
yay -S slides
```

From source:
```
git clone https://github.com/maaslalani/slides.git
cd slides
go install
```

You can also download a binary from the [releases](https://github.com/maaslalani/slides/releases) page.

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

---

Include ASCII graphs with GraphViz + graph-easy.
https://dot-to-ascii.ggerganov.com/

┌──────────┐     ┌────────────┐     ┌────────┐
│ GraphViz │ ──▶ │ graph-easy │ ──▶ │ slides │
└──────────┘     └────────────┘     └────────┘

```

Checkout the [example slides](./examples).

Then, to present, run:
```
slides presentation.md
```

Live reload is supported by default. It means no need to restart the presentation while continue editing the slides.

`slides` also accepts input through `stdin`:
```
curl http://example.com/slides.md | slides
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

### Configuration
See the [configuration documentation](./docs/configuration)

### Alternatives

**Credits**: This project was heavily inspired by [`lookatme`](https://github.com/d0c-s4vage/lookatme).

* [`lookatme`](https://github.com/d0c-s4vage/lookatme)
* [`sli.dev`](https://sli.dev/)
* [`sent`](https://tools.suckless.org/sent/)

### Development
See the [development documentation](./docs/development)
