# sequin

Parse ANSI sequences into human-readable format.

```console
$ printf '\033[48;2;255;0;0m\033[m' | sequin
CSI "\x1b[48;2;255;0;0m": Background color: {255 0 0 255}
CSI "\x1b[m": Reset style
```

## Use cases

The most basic use case is to describe what an escape sequence does when you
don't know it.

You can `printf '\x1betc' | sequin` and get an explanation!

More complex use cases might include checking [teatest's][] golden files, e.g.:

```console
$ cat ./testdata/TestApp.golden | sequin
CSI "\x1b[?25l": Disable private mode "cursor visibility"
CSI "\x1b[?2004h": Enable private mode "bracketed paste"
Control code "\r": carriage return
Text: "Hi. This program will exit in 10 seconds. To quit sooner press any key"
Control code "\n": line feed
CSI "\x1b[70D": Cursor left 70
CSI "\x1b[A": Cursor up 1
CSI "\x1b[70D": Cursor left 70
CSI "\x1b[2K": Erase entire line
Text: "Hi. This program will exit in 9 seconds. To quit sooner press any key."
Control code "\n": line feed
CSI "\x1b[70D": Cursor left 70
CSI "\x1b[2K": Erase entire line
Control code "\r": carriage return
CSI "\x1b[?2004l": Disable private mode "bracketed paste"
CSI "\x1b[?25h": Enable private mode "cursor visibility"
CSI "\x1b[?1002l": Disable private mode "mouse cell motion"
CSI "\x1b[?1003l": Disable private mode "mouse all motion"
CSI "\x1b[?1006l": Disable private mode "mouse SGR ext"
```

You can also use it to check the output of any program[^pipe], for instance, `ls` and `glow`:

```console
$ ls -1 --color=always | sequin
CSI "\x1b[38;5;4m": Foreground color: 4
CSI "\x1b[1m": Bold
Text: "folder"
CSI "\x1b[0m": Reset style
Control code "\n": line feed
Text: "file.txt"

$ glow -s dark README.md | sequin
Control code "\n": line feed
CSI "\x1b[;;1m": , , Bold
CSI "\x1b[0m": Reset style
CSI "\x1b[;;1m": , , Bold
CSI "\x1b[0m": Reset style
Text: "  "
CSI "\x1b[;;1m": , , Bold
Text: " "
CSI "\x1b[0m": Reset style
CSI "\x1b[;;1m": , , Bold
Text: "sequin"
```

So you may also use it to debug applications, and of course, to learn more!

## How it works

It relies heavily on our [ansi][] package, and whilst traversing the strings,
pretty prints what the sequences are doing.

Check [ansi][] out to learn more!

[teatest]: https://github.com/charmbracelet/x/tree/main/exp/teatest
[ansi]: https://github.com/charmbracelet/x/tree/main/ansi

## Current state

Common sequences are implemented, but there is still plenty of work to do!
If you notice one of such missing sequences, feel free to open a PR. 😄

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note!

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/lipgloss/raw/master/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source

[^pipe]:
    Beware: some programs might render their output as plain text when it
    detects that the output isn't a terminal (e.g. when redirecting to a file,
    or to another program, like we are doing in the examples).
    Usually there are ways around this, like setting `CLICOLOR_FORCE=1` or flags
    to force colorful output. You might need to check what works in your case!
