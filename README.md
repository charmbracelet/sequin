# sequin

Parse ANSI sequences into human-readable format.

```console
$ printf '\033[48;2;255;0;0m' | sequin
CSI 48;2;255;0;0m: Set background color to RGB(255,0,0)
```

## State

Mostly generated with AI/manual edits using Claude and
[charmbracelet/x/ansi][ansi].

This is very pre-alpha-do-not-use-yet state.

[ansi]: https://pkg.go.dev/github.com/charmbracelet/x/ansi

## TODO

why this seq breaks it?

```
\x1b]22;wait\x07; \x1b7 \x1b8 \x1b[?6n \x1b[6n \x1b]12;#ff00ff\x07
```
