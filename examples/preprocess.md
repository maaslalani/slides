# Slides

~~~sd replaced processed
This content will be replaced and passed into stdin
of the command above
~~~

---

Any command will work

~~~echo "You can do whatever, really"
This doesn't matter, since it will be replaced by the stdout
of the command above
~~~

---

Pre-process Graphs

~~~graph-easy --as=boxart
[ A ] - to -> [ B ]
~~~

The above will be pre-processed to look like:

```
┌───┐  to   ┌───┐
│ A │ ────> │ B │
└───┘       └───┘
```
