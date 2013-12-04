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
on your `$PATH` and make sure [your GitHub access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use)
is available from the `ACCESS_TOKEN` environmetal variable at all times.

### TODO

* Add support for configuring user / project names to fill in `{{USER}}` and `{{NAME}}` on templates
* Fix relative linking
