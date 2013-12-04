# viewdocs-preview

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
