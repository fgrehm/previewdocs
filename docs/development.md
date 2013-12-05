# Development

viewdocs-preview is written in [Go](http://golang.org/) and interacts with the [GitHub
Markdown API](http://developer.github.com/v3/markdown/).

If you want to hack on it, first you'll need to [get your GitHub access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use)
and make sure it is available at all times from the `ACCESS_TOKEN` environmental
variable.

After that, just sing that same old song:

```sh
go get github.com/fgrehm/viewdocs-preview
cd $GOPATH/src/github.com/fgrehm/viewdocs-preview
go run viewdocs.go
```

Then visit `http://localhost:8888/viewdocs-preview` on your browser and enjoy!
