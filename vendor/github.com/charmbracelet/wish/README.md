# Wish

<p>
    <picture>
        <source srcset="https://stuff.charm.sh/wish/wish-header.webp" type="image/webp">
        <img style="width: 247px" src="https://stuff.charm.sh/wish/wish-header.png" alt="A nice rendering of a star, anthropomorphized somewhat by means of a smile, with the words ‘Charm Wish’ next to it">
    </picture><br>
    <a href="https://github.com/charmbracelet/wish/releases"><img src="https://img.shields.io/github/release/charmbracelet/wish.svg" alt="Latest Release"></a>
    <a href="https://pkg.go.dev/github.com/charmbracelet/wish?tab=doc"><img src="https://godoc.org/github.com/golang/gddo?status.svg" alt="GoDoc"></a>
    <a href="https://github.com/charmbracelet/wish/actions"><img src="https://github.com/charmbracelet/wish/workflows/Build/badge.svg" alt="Build Status"></a>
    <a href="https://codecov.io/gh/charmbracelet/wish"><img alt="Codecov branch" src="https://img.shields.io/codecov/c/github/charmbracelet/wish/main.svg"></a>
    <a href="https://goreportcard.com/report/github.com/charmbracelet/wish"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/charmbracelet/wish"></a>
</p>


Make SSH apps, just like that! 💫

SSH is an excellent platform to build remotely accessible applications on. It
offers secure communication without the hassle of HTTPS certificates, it has
user identification with SSH keys and it's accessible from anywhere with a
terminal. Powerful protocols like Git work over SSH and you can even render
TUIs directly over an SSH connection.

Wish is an SSH server with sensible defaults and a collection of middleware that
makes building SSH apps easy. Wish is built on [gliderlabs/ssh][gliderlabs/ssh]
and should be easy to integrate into any existing projects.

## What are SSH Apps?

Usually, when we think about SSH, we think about remote shell access into servers,
most commonly through `openssh-server`.

That's a perfectly valid and probably the most common use of SSH, but it can do so much more than that.
Just like HTTP, SMTP, FTP and others, SSH is a protocol!
It is a cryptographic network protocol for operating network services securely over an unsecured network. [^1]

[^1]: https://en.wikipedia.org/wiki/Secure_Shell

That means, among other things, that we can write custom SSH servers, without touching `openssh-server`,
so we can securely do more things than just providing a shell.

Wish is a library that helps writing these kind of apps using Go.

## Middleware

Wish middlewares are analogous to those in several HTTP frameworks.
They are essentially SSH handlers that you can use to do specific tasks,
and then call the next middleware.

Notice that middlewares are composed from first to last,
which means the last one is executed first.

### Bubble Tea

The [`bubbletea`](bubbletea) middleware makes it easy to serve any
[Bubble Tea][bubbletea] application over SSH. Each SSH session will get their own
`tea.Program` with the SSH pty input and output connected. Client window
dimension and resize messages are also natively handled by the `tea.Program`.

You can see a demo of the Wish middleware in action at: `ssh git.charm.sh`

### Git

The [`git`](git) middleware adds `git` server functionality to any ssh server.
It supports repo creation on initial push and custom public key based auth.

This middleware requires that `git` is installed on the server.

### Logging

The [`logging`](logging)  middleware provides basic connection logging. Connects
are logged with the remote address, invoked command, TERM setting, window
dimensions and if the auth was public key based. Disconnect will log the remote
address and connection duration.

### Access Control

Not all applications will support general SSH connections. To restrict access
to supported methods, you can use the [`activeterm`](activeterm) middleware to
only allow connections with active terminals connected and the
[`accesscontrol`](accesscontrol) middleware that lets you specify allowed
commands.

## Default Server

Wish includes the ability to easily create an always authenticating default SSH
server with automatic server key generation.

## Examples

There are examples for a standalone [Bubble Tea application](examples/bubbletea)
and [Git server](examples/git) in the [examples](examples) folder.

## Apps Built With Wish

* [Soft Serve](https://github.com/charmbracelet/soft-serve)
* [Wishlist](https://github.com/charmbracelet/wishlist)
* [SSHWordle](https://github.com/davidcroda/sshwordle)
* [clidle](https://github.com/ajeetdsouza/clidle)
* [ssh-warm-welcome](https://git.vvvvvvaria.org/decentral1se/ssh-warm-welcome)

[bubbletea]: https://github.com/charmbracelet/bubbletea
[gliderlabs/ssh]: https://github.com/gliderlabs/ssh

## Pro Tip

When building various Wish applications locally you can add the following to
your `~/.ssh/config` to avoid having to clear out `localhost` entries in your
`~/.ssh/known_hosts` file:

```
Host localhost
    UserKnownHostsFile /dev/null
```

## License

[MIT](https://github.com/charmbracelet/wish/raw/main/LICENSE)

***

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge-unrounded.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
