# Welcome to previewdocs

[Viewdocs](http://viewdocs.io/) is [Read the Docs](https://readthedocs.org/)
meets [Gist.io](http://gist.io/) for simple project documentation. It renders
Markdown from your repository's `docs` directory as simple static pages.

`previewdocs` is a *work in progress* CLI tool that helps you preview
changes to your documentation before pushing the code back to your repository.

### Installation

Download the latest released version for your platform from [GitHub](https://github.com/jingweno/gh/releases/latest),
place it somewhere on your `PATH` and make it executable.

If you think you'll reach GitHub's [Rate Limit](http://developer.github.com/v3/#rate-limiting)
of 60 requests per hour while working on your docs, please create a [GitHub access token](https://help.github.com/articles/creating-an-access-token-for-command-line-use)
and place the token on `$HOME/.previewdocsrc` or set it as the the `ACCESS_TOKEN`
environmetal variable.

### Getting Started

If you are new to Viewdocs, just make a `docs` directory in your GitHub project
repository and put an `index.md` file in there to get started.

Then run `previewdocs` from the project's root and browse to:

	http://localhost:8888/project-name

Any other Markdown files in your `docs` directory are available as a subpath,
including files in directories. You can update pages and hit F5 to see the
changes as you go instead of pushing the code back to the GitHub repository
and waiting for Viewdocs cache to expire.

This page is an example of what documentation will look like by default.
Here is [another example page](example). The source for
these pages are in the [docs directory](https://github.com/fgrehm/previewdocs/tree/master/docs)
of the project.

### Custom layouts

Viewdocs supports custom layouts for your docs. You can make your own `docs/template.html`
based on the [default viewdocs template](https://github.com/progrium/viewdocs/blob/master/docs/template.html)
and your pages will be rendered with that template.

### More information

I highly recommend you [read the source](https://github.com/fgrehm/previewdocs/blob/master/previewdocs.go)
of this app. It's less than 200 lines of Go. If you want to hack on `previewdocs`,
[check this out](development.md).

<br />
Enjoy!<br />
[Fábio Rehm](http://twitter.com/fgrehm)
