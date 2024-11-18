# Sequin

<p>
    <a href="https://github.com/charmbracelet/sequin/releases"><img src="https://img.shields.io/github/release/charmbracelet/sequin.svg" alt="Latest Release"></a>
    <a href="https://github.com/charmbracelet/sequin/actions"><img src="https://github.com/charmbracelet/sequin/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Human-readable ANSI sequences.

<p><img src="https://github.com/user-attachments/assets/e48d3dce-b619-4a10-848f-2b66532a31e7" width="800" alt="Sequin doing its thing"></p>

Sequin is a small utility that can help you debug your CLIs and TUIs. It's also
great for describing escape sequences you might not understand, and exploring
what TUIs are doing under the hood.

There are lots more use cases too, like inspecting golden files such as the
ones used by [`teatest`][teatest] to crystalize [Bubble Tea][bubbletea] output.

Are you using Sequin in an interesting way? We’d love to hear about it.

## Examples

### Describing escape sequences

Just use `printf` to send some sequences to `sequin` for an explanation.

```bash
printf "\x1b[38;5;4mCiao, \x1b[1;7mBaby.\x1b[0m\n" | sequin
```

<p><img src="https://github.com/user-attachments/assets/5df48244-2e20-4742-b403-39c7534d10b8" width="550" alt="Sequin with printf"></p>

### Reading the output of a program

You can also use it to check the output of any program, for instance, `ls` or `git`.

```bash
ls -1 --color=always | sequin
```

<p><img src="https://github.com/user-attachments/assets/75166b2b-7cf5-4d78-97b2-901d00474591" width="400"></p>

```bash
git -c status.color=always status -sb | sequin
```

<p><img src="https://github.com/user-attachments/assets/ebd71c27-198b-4bad-9096-91bc1e944bad" width="450"></p>

So yeah, it’s great for debugging applications, and of course, learning about ANSI!

> [!NOTE]
> Many programs render their output as plain text when output isn't a terminal
> (i.e. when redirecting to a file or piping to a program, like `sequin`). This
> is a good thing, except in this case when we actually want ANSI sequences so
> we can inspect them. Thankfully there are usually ways to force colors, like
> by setting `CLICOLOR_FORCE=1` or with flags to force ANSI output. If you're
> not seeing sequences be sure to to check what works in the case of your
> specific program.

### Examining golden files

Golden file for TUIs contain ANSI, which can be easily inspected with `sequin`:

```console
$ cat ./testdata/MyCuteApp.golden | sequin
```

<p><img src="https://github.com/user-attachments/assets/16367a79-0ee3-40e1-95ae-adc46f411192" width="580"></p>

To generate golden files for your TUIs have a look at [`golden`][golden] and [`teatest`][teatest] from the [`/x`][x] project.

### Visualizing raw sequences

One of the paint points that we find when working with ANSI escape sequences is
that they are just a bunch of "special" characters and text. Sometimes instead
of explaining the sequences line by line, you just want to highlight them to
visually see what's going on.

```bash
git -c status.color=always status -sb | sequin -r && echo
```

<p><img src="https://github.com/user-attachments/assets/2f812aa8-ba40-435f-858d-60f5082c7fad" width="450"></p>

## How it works

It relies heavily on our glorious [`ansi`][ansi] package, currently in the
elusive [`/x`][x] project. Whilst traversing the strings, Sequin pretty prints
what the sequences are and what they’re doing.

[ansi]: https://pkg.go.dev/github.com/charmbracelet/x/ansi
[bubbletea]: https://github.com/charmbracelet/bubbletea
[golden]: https://pkg.go.dev/github.com/charmbracelet/x/exp/golden
[teatest]: https://github.com/charmbracelet/x/tree/main/exp/teatest
[x]: https://github.com/charmbracelet/x

## Is it done?

No! Common sequences are implemented, but there is still plenty of work to
do. For instance, [APC](https://www.apc.fr) sequences are not supported yet. If
you notice one of such missing sequences, or want to work on any other area of
the project, feel free to open a PR. 💘

## Contributing

We love contributions. We recommend checking out [our contribution
guidelines][contributing] for faster responses on our end.

[contributing]: https://github.com/charmbracelet/.github/blob/main/CONTRIBUTING.md

## Feedback

We’d love to hear your thoughts on this project. Feel free to drop us a note.

- [Twitter](https://twitter.com/charmcli)
- [The Fediverse](https://mastodon.social/@charmcli)
- [Discord](https://charm.sh/chat)

## License

[MIT](https://github.com/charmbracelet/sequin/raw/master/LICENSE)

---

Part of [Charm](https://charm.sh).

<a href="https://charm.sh/"><img alt="The Charm logo" src="https://stuff.charm.sh/charm-badge.jpg" width="400"></a>

Charm热爱开源 • Charm loves open source
