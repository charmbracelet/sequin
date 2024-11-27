# Sequin

<p>
    <img width="356" alt="Charm Sequin Logo" src="https://github.com/user-attachments/assets/5d41c071-c27f-412e-8a00-4e1258c70ceb"><br>
    <a href="https://github.com/charmbracelet/sequin/releases"><img src="https://img.shields.io/github/release/charmbracelet/sequin.svg" alt="Latest Release"></a>
    <a href="https://github.com/charmbracelet/sequin/actions"><img src="https://github.com/charmbracelet/sequin/workflows/build/badge.svg" alt="Build Status"></a>
</p>

Human-readable ANSI sequences.

<p><img width="766" src="https://github.com/user-attachments/assets/60cce45d-be8a-4d77-a063-a3cf3d952966"></p>

Sequin is a small utility that can help you debug your CLIs and TUIs. It's also
great for describing escape sequences you might not understand, and exploring
what TUIs are doing under the hood.

There are lots more use cases too, like inspecting golden files such as the
ones used by [`teatest`][teatest] to crystalize [Bubble Tea][bubbletea] output.

Are you using Sequin in an interesting way? We’d love to hear about it.

## Installation

Use a package manager:

```bash
# macOS or Linux
brew install charmbracelet/tap/sequin

# Arch Linux (btw)
yay -S sequin-bin

# Nix
nix-shell -p nur.repos.charmbracelet.sequin
```

<details>
<summary>Debian/Ubuntu</summary>

```bash
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
sudo apt update && sudo apt install sequin
```

</details>

<details>
<summary>Fedora/RHEL</summary>

```bash
echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
sudo yum install sequin
```

</details>

Or, download it:

- [Packages][releases] are available in Debian and RPM formats
- [Binaries][releases] are available for Linux, macOS, Windows, FreeBSD, OpenBSD, and NetBSD

[releases]: https://github.com/charmbracelet/sequin/releases

Or, just install it with `go`:

```sh
go install github.com/charmbracelet/sequin@latest
```

<details>
<summary>Shell Completions</summary>

All the packages and archives come with pre-generated completion files for Bash,
ZSH, Fish, and PowerShell.

If you built it from source, you can generate them with:

```bash
sequin completion bash -h
sequin completion zsh -h
sequin completion fish -h
sequin completion powershell -h
```

If you use a package (like Homebrew, Debs, etc), the completions should be set
up automatically, given your shell is configured properly.

</details>

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

```bash
cat ./testdata/MyCuteApp.golden | sequin
```

<p><img src="https://github.com/user-attachments/assets/16367a79-0ee3-40e1-95ae-adc46f411192" width="580"></p>

To generate golden files for your TUIs have a look at [`golden`][golden] and [`teatest`][teatest] from the [`/x`][x] project.

### Fake TTY - executing commands

You can also execute commands directly in sequin:

```bash
sequin -- ls -1 go.*
```

<p><img src="https://github.com/user-attachments/assets/efd9f511-130d-49e8-ba8f-31e1e3d86920" width="450"></p>

## Pro Mode: Syntax Highlighting for Raw Sequences

One of the pain points that we find when reading raw ANSI output is
that it’s hard to visually separate sequences from regular text. For situations
like this you can use the `--raw`/`-r` flag to simply highlight sequences inline:

```bash
git -c status.color=always status -sb | sequin -r && echo
```

<p><img src="https://github.com/user-attachments/assets/c3b19a81-934e-4b87-b86d-2aa2a25b8c5d" width="450"></p>

## How it all works

Sequin relies heavily on our glorious [`ansi`][ansi] package, currently in the
elusive [`/x`][x] project. Whilst traversing the strings, Sequin pretty prints
what the sequences are and what they’re doing.

[ansi]: https://pkg.go.dev/github.com/charmbracelet/x/ansi
[bubbletea]: https://github.com/charmbracelet/bubbletea
[golden]: https://pkg.go.dev/github.com/charmbracelet/x/exp/golden
[teatest]: https://github.com/charmbracelet/x/tree/main/exp/teatest
[x]: https://github.com/charmbracelet/x

## Is it done?

No! Common sequences are implemented, but there is still plenty of work to
do. For instance, APC sequences are not supported yet. If you notice one
of such missing sequences, or want to work on any other area of the project,
feel free to open a PR. 💘

## FAQ

<details>
<summary>
Why isn't sequin using my terminal's background color? (with multiplexer)
</summary>
related: <a href="https://github.com/charmbracelet/sequin/issues/28" target="_blank">https://github.com/charmbracelet/sequin/issues/28</a>
<h2>Check your multiplexer's version</h2>
<h3>Tmux</h3>
<p>This could be a few things depending on your environment. First thing to
consider: are you using a terminal multiplexer (e.g. tmux)? If so, check your
tmux version is 3.4 or higher. (Tmux added support for OSC10/11 queries in 3.4 see
<a href="https://github.com/tmux/tmux/blob/master/CHANGES#L152" target="_blank">https://github.com/tmux/tmux/blob/master/CHANGES#L152)</a>. Not sure which version of tmux you're using? <code>tmux -V</code></p>

<h3>GNU Screen</h3>
<p>If you're a GNU screen user, check that your version is 4.99.0 or higher to
ensure OSC10/11 can be used. Not sure which version of screen you're using? <code>screen --version</code></p>

<h2>Check your config</h2>
<p>Make sure your multiplexer config has your default terminal and <code>COLORTERM</code> set.</p>
<h3><code>tmux.conf</code></h3>
<pre>
set -g default-terminal "xterm-ghostty"
setenv -g COLORTERM "truecolor"
</pre>

<h3><code>~/.screenrc</code></h3>
<pre>
term xterm-256color
setenv COLORTERM truecolor
defbce on
</pre>

<p>Note: if you don't set <code>defbce on</code>, you will get a different background color on cells
cleared with an erase/insert/scroll/clear operation in the regular background
color.</p>
</details>

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
