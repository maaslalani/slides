# Slides

Slides in your terminal.

<p align="center">
  <img src="./assets/slides.gif?raw=true" alt="Slides Presentation" />
</p>

### Installation
[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://formulae.brew.sh/api/formula/slides.json&query=$.versions.stable&label=homebrew)](https://formulae.brew.sh/formula/slides)
[![Snapcraft](https://snapcraft.io/slides/badge.svg)](https://snapcraft.io/slides)
[![AUR](https://img.shields.io/aur/version/slides)](https://aur.archlinux.org/packages/slides)

<details>
<summary>Instructions</summary>

#### MacOS
```
brew install slides
```
#### Arch
```
yay -S slides
```
#### Nixpkgs (unstable)
```
nix-env -iA nixpkgs.slides
```
#### Any Linux Distro running `snapd`
```
sudo snap install slides
```
#### Go
```
go install github.com/maaslalani/slides@latest
```
From source:
```
git clone https://github.com/maaslalani/slides.git
cd slides
go install
```

You can also download a binary from the [releases](https://github.com/maaslalani/slides/releases) page.

</details>


### Usage
Create a simple markdown file that contains your slides:

````markdown
# Welcome to Slides
A terminal based presentation tool

---

## Everything is markdown
In fact, this entire presentation is a markdown file.

---

## Everything happens in your terminal
Create slides and present them without ever leaving your terminal.

---

## Code execution
```go
package main

import "fmt"

func main() {
  fmt.Println("Execute code directly inside the slides")
}
```

You can execute code inside your slides by pressing `<C-e>`,
the output of your command will be displayed at the end of the current slide.

---

## Pre-process slides

You can add a code block with three tildes (`~`) and write a command to run *before* displaying
the slides, the text inside the code block will be passed as `stdin` to the command
and the code block will be replaced with the `stdout` of the command.

~~~graph-easy --as=boxart
[ A ] - to -> [ B ]
~~~

The above will be pre-processed to look like:

┌───┐  to   ┌───┐
│ A │ ────> │ B │
└───┘       └───┘

For security reasons, you must pass a file that has execution permissions
for the slides to be pre-processed. You can use `chmod` to add these permissions.

```bash
chmod +x file.md
```

````

Checkout the [example slides](https://github.com/maaslalani/slides/tree/main/examples).

Then, to present, run:
```
slides presentation.md
```

If given a file name, `slides` will automatically look for changes in the file and update the presentation live.

`slides` also accepts input through `stdin`:
```
curl http://example.com/slides.md | slides
```

Go to the first slide with the following key sequence:
* <kbd>g</kbd> <kbd>g</kbd>

Go to the next slide with any of the following key sequences:
* <kbd>space</kbd>
* <kbd>right</kbd>
* <kbd>down</kbd>
* <kbd>enter</kbd>
* <kbd>n</kbd>
* <kbd>j</kbd>
* <kbd>l</kbd>
* number + any of the above (go forward n slides)

Go to the previous slide with any of the following key sequences:
* <kbd>left</kbd>
* <kbd>up</kbd>
* <kbd>p</kbd>
* <kbd>h</kbd>
* <kbd>k</kbd>
* number + any of the above (go back n slides)

Go to a specific slide with the following key sequence:
* number + <kbd>G</kbd>

Go to the last slide with the following key:
* <kbd>G</kbd>

### Code Execution

If slides finds a code block on the current slides it can execute the code
block and display the result as virtual text on the screen.

Press <kbd>ctrl+e</kbd> on a slide with a code block to execute it and display the result.

### Configuration

`slides` allows you to customize your presentation's look and feel with metadata at the top of your `slides.md`.

> This section is entirely optional, `slides` will use sensible defaults if this section or any field in the section is omitted.

```yaml
---
theme: ./path/to/theme.json
author: Gopher
date: January 2, 2006
paging: Slide %d / %d
---
```

* `theme`: Path to `json` file containing a [glamour
  theme](https://github.com/charmbracelet/glamour/tree/master/styles), can also
  be a link to a remote `json` file which slides will fetch before presenting.
* `author`: A `string` to display on the bottom-left corner of the presentation
  view. Defaults to the OS current user's full name. Can be empty to hide the author.
* `date`: A `string` that is used to format today's date in the native Go
  format `2006-01-02` or in the `YYYY-MM-DD` format. If the date is not a valid
  format, the string will be displayed. Defaults to `2006-01-02`.
* `paging`: A `string` that contains 0 or more `%d` directives. The first `%d`
  will be replaced with the current slide number and the second `%d` will be
  replaced with the total slides count. Defaults to `Slide %d / %d`.

#### Date format

Given the date _January 02, 2006_:

| Value  | Translates to |
|--------|---------------|
| `YYYY` | 2006          |
| `YY`   | 06            |
| `MMMM` | January       |
| `MMM`  | Jan           |
| `MM`   | 01            |
| `mm`   | 1             |
| `DD`   | 02            |
| `dd`   | 2             |

### Alternatives

**Credits**: This project was heavily inspired by [`lookatme`](https://github.com/d0c-s4vage/lookatme).

* [`lookatme`](https://github.com/d0c-s4vage/lookatme)
* [`sli.dev`](https://sli.dev/)
* [`sent`](https://tools.suckless.org/sent/)

### Development
See the [development documentation](./docs/development)
