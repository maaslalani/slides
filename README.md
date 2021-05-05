# Slides

Slides in your terminal

### Installation
```
go get github.com/maaslalani/slides
```

### Usage
Create a simple markdown file that contains your slides:

```
# Welcome to Slides
A terminal based presentation tool

~~~

## Everything is markdown
In fact, this entire presentation is a markdown file.

~~~

## Everything happens in your terminal
Create slides and present them without ever leaving your terminal.

```

Checkout the [example slides](./examples).

Then, to present, run:
```
slides presentation.md
```

Go to the next slide with any of the following keys:
* <kbd>right</kbd>
* <kbd>up</kbd>
* <kbd>enter</kbd>
* <kbd>n</kbd>
* <kbd>k</kbd>
* <kbd>l</kbd>

Go to the previous slide with any of the following keys:
* <kbd>left</kbd>
* <kbd>down</kbd>
* <kbd>p</kbd>
* <kbd>h</kbd>
* <kbd>j</kbd>

### Demo
![slides](../assets/slides.gif?raw=true)

### Development
```
make
```
