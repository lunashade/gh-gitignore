# gh-gitignore

Print gitignore template using [Gitignore API](https://docs.github.com/ja/rest/gitignore).
Written in Go.

## Install & Usage

```bash
$ gh extension install lunashade/gh-gitignore
$ gh gitignore > .gitignore
```

## development

1. make sure remove extension `gh-gitignore`.
2. clone this repo
3. Run `go build && gh extension install .` to symlink install the extension
4. Edit source, build and test it.
