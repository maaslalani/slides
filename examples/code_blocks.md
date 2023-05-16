# Code blocks

Slides allows you to execute code blocks directly inside your slides!

Just press `ctrl+e` and the result of the code block will be displayed as virtual text in your slides.

Currently supported languages:

<!-- Use comments in your markdown! -->

* `bash`
* `zsh`
* `fish`
* `elixir`
* `go`
* `javascript`
* `python`
* `ruby`
* `perl`
* `rust`
* `java`
* `cpp`
* `swift`
<!-- * `secret` -->

---

### Bash

```bash
ls
```

---

### Zsh

```zsh
ls
```

---

### Fish

```fish
ls
```

---

### Elixir

```elixir
IO.puts "Hello, world!"
```

---

### Go

Use `///` to hide verbose code but still allow the ability to execute it.

If you press `y` to copy (yank) this code block it will return the full snippet.

And, if you press `ctrl+e` it will run the program without error, even though
what is being displayed is not a valid go program because we have commented out
some boilerplate to focus on the important parts.

```go
///package main
///
import "fmt"
///
///func main() {
fmt.Println("Hello, world!")
///}
```

---

### Javascript

```javascript
console.log("Hello, world!")
```

---

### Lua

```lua
print("Hello, World!")
```

---

### Python

```python
print("Hello, world!")
```

---

### Ruby

```ruby
puts "Hello, world!"
```

---

### Perl

```perl
print ("hello, world");
```

---

### Rust

```rust
fn main() {
    println!("Hello, world!");
}
```

---

### Java
```java
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, world!");
    }
}
```

---

### Julia
```julia
println("Hello, world!")
```

---

### C++
```cpp
#include <iostream>

int main() {
    std::cout << "Hello, world!" << std::endl;
    return 0;
}
```

---

### Swift
```swift
print("Hello, world!")
```
