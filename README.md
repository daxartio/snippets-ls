# snippets-ls

The ls is an integration designed for Helix Editor, enabling the utilization of snippets from Visual Studio Code as auto-completion suggestions within the Helix Editor's Language Server Protocol (LSP) environment.

## Installation

### Install via go

```sh
go install github.com/daxartio/snippets-ls@latest
```

Don't forget to append `~/go/bin` to your `$PATH`.

## Usage

Create your own snippets follow [VSCode syntax](https://code.visualstudio.com/docs/editor/userdefinedsnippets#_create-your-own-snippets). Alternatively, you can make use of [pre-existing](https://github.com/microsoft/vscode-go/blob/master/snippets/go.json) [sample](https://github.com/rust-lang/vscode-rust/blob/master/snippets/rust.json) for various programming languages.

Update your configuration file located at `~/.config/helix/languages.toml`:

```toml
[[language]]
name = "go"
formatter = { command = "goimports"}
language-servers = ["gopls", "snippets-ls"]

[language-server.snippets-ls]
command = "snippets-ls"
args = ["-lang", "go"]
```

Subsequently, as you start working on your file, input a snippet prefix to observe the suggestion.
If it does not work, take a look at `~/.cache/helix/helix.log` for additional insights.
