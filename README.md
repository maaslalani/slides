# Slides

Slides in your terminal.

<p align="center">
  <img src="./assets/slides.gif?raw=true" alt="Slides Presentation" />
</p>

### Installation
[![Homebrew](https://img.shields.io/badge/dynamic/json.svg?url=https://formulae.brew.sh/api/formula/slides.json&query=$.versions.stable&label=homebrew)](https://formulae.brew.sh/formula/slides)
[![Snapcraft](https://snapcraft.io/slides/badge.svg)](https://snapcraft.io/slides)
[![AUR](https://img.shields.io/aur/version/slides)](https://aur.archlinux.org/packages/slides)

* MacOS
```
brew install slides
```

* Arch
```
yay -S slides
```

* Nixpkgs (unstable)
```
nix-env -iA nixpkgs.slides
```

* Any Linux Distro running `snapd`

```
sudo snap install slides
```

* Go
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

You can add a code block with `~~~` and write a command to run *before* displaying
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

Go to the next slide with any of the following keys:
* <kbd>space</kbd>
* <kbd>right</kbd>
* <kbd>down</kbd>
* <kbd>enter</kbd>
* <kbd>n</kbd>
* <kbd>j</kbd>
* <kbd>l</kbd>

Go to the previous slide with any of the following keys:
* <kbd>left</kbd>
* <kbd>up</kbd>
* <kbd>p</kbd>
* <kbd>h</kbd>
* <kbd>k</kbd>

Press <kbd>ctrl+e</kbd> on a slide with a code block to execute it.

### Configuration

### Theme

`slides` allows you to customize your presentation's theme.

If you want to use your own custom [theme.json](./styles/theme.json), add the following to the top of your `presentation.md`:
```yaml
---
theme: ./path/to/theme.json
---
```

Check out the provided [theme.json](./styles/theme.json) to use as a base for your custom theme.

You can also pass the URL to a theme and `slides` will fetch it from the remote before presenting so that you don't need to keep the theme file on disk.

```yaml
---
theme: https://github.com/maaslalani/slides/raw/main/styles/theme.json
---
```

### Author

If you want to customize the name displayed on the bottom left corner of `slides`,
you may pass a string to the `author` field in the metadata section of your `slides.md` file.
This value defaults to your OS user's full name.

```yaml
---
author: Gopher
---
```

### Date

To customize the format of the date string being used,
you may pass a string formatted the way you want using the date `2006-01-02`.
This value defaults to `2006-01-02`.

```yaml
---
date: January 2, 2006
---
```

To hide the date, you may also pass an empty string `""`.

### Paging

To customize the page information displayed on the bottom right corner of
`slides`, you may pass a string containing 0 or more `%d` directives which will
automatically be formatted with the current slide number and total slide
number.
This value defaults to `Slide %d / %d`.

### Alternatives

**Credits**: This project was heavily inspired by [`lookatme`](https://github.com/d0c-s4vage/lookatme).

* [`lookatme`](https://github.com/d0c-s4vage/lookatme)
* [`sli.dev`](https://sli.dev/)
* [`sent`](https://tools.suckless.org/sent/)

### Development
See the [development documentation](./docs/development)
