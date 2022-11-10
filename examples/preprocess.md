# Slides

You can add a code block with three tildes (~) and write a command to run before displaying
the slides, the text inside the code block will be passed as stdin to the command
and the code block will be replaced with the stdout of the command.

```
~~~graph-easy --as=boxart
[ A ] - to -> [ B ]
~~~
```

The above will be pre-processed to look like:

NOTE: You need `graph-easy` installed and in your `$PATH`

```
┌───┐  to   ┌───┐
│ A │ ────> │ B │
└───┘       └───┘
```

For security reasons, you must pass a file that has execution permissions
for the slides to be pre-processed.

```
chmod +x file.md
```

---

~~~sd replaced processed
This content will be passed in as stdin and will be replaced.
~~~

---


Any command will work

~~~echo "You can do whatever, really"
This doesn't matter, since it will be replaced by the stdout
of the command above because the command will disregard stdin.
~~~
---


You can use this to import snippets of code from other files:

~~~xargs cat
examples/import.md
~~~

---


## More pre-process examples:

### PlantUML

```
~~~plantuml -utxt -pipe
@startuml
A --> B: to
@enduml
~~~
```

The above will be pre-processed to look like:

NOTE: You need `plantuml` installed and in your `$PATH`

```
┌─┐          ┌─┐
│A│          │B│
└┬┘          └┬┘
 │    to      │
 │ ─ ─ ─ ─ ─ >│
┌┴┐          ┌┴┐
│A│          │B│
└─┘          └─┘

