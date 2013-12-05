# viewdocs-preview

**WORK IN PROGRESS**

Preview [Viewdocs](http://viewdocs.io/) documentation before pushing changes
back to your repositories.

## Installation

Right now it is only possible to install the `viewdocs-preview` command from
sources. I'll provide precompiled releases of it as soon as things are stable.

Assuming you have your [`$GOPATH`](http://golang.org/doc/code.html#GOPATH)
configured properly, run:

```
go get github.com/fgrehm/viewdocs-preview
cd $GOPATH/src/github.com/fgrehm/viewdocs-preview
go build
```

Then drop the generated `viewdocs-preview` executable on a directory available
on your `$PATH`.

If you think you'll reach GitHub's [Rate Limit](http://developer.github.com/v3/#rate-limiting)
of 60 requests per hour while working on your docs, please set the `ACCESS_TOKEN`
environmetal variable to your [GitHub access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use).
