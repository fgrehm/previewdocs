# viewdocs-preview

**WORK IN PROGRESS**

Preview [Viewdocs](http://viewdocs.io/) documentation before pushing changes
back to your repositories.

## Installation

```
mkdir -p $GOPATH/src/github.com/fgrehm
cd $GOPATH/src/github.com/fgrehm
git clone https://github.com/fgrehm/viewdocs-preview.git
cd viewdocs-preview
go get
go build
```

Then drop the generated `viewdocs-preview` executable on a directory available
on your `$PATH`.

If you think you'll reach GitHub's [Rate Limit](http://developer.github.com/v3/#rate-limiting)
of 60 requests per hour while working on your docs, please set the `ACCESS_TOKEN`
environmetal variable to your [GitHub access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use).

### TODO

* Fix relative linking
