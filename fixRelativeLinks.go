package main

import (
	"strings"
	"path"
	"bytes"

	"code.google.com/p/go.net/html"
)

func fixRelativeLinks(doc string, repo string, body string) (string, error) {
	n, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i, a := range n.Attr {
				if a.Key == "href" {
					fs := strings.Index(a.Val, "/")
					fc := strings.Index(a.Val, ":")
					fh := strings.Index(a.Val, "#")
					if fs == 0 || fh == 0 ||
					(fc >= 0 && fc < fs) ||
					(fh >= 0 && fh < fs) {
						continue
					}
					dir := path.Dir(doc)
					n.Attr[i].Val = "/" + repo + "/" + dir + "/" + a.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	b := new(bytes.Buffer)
	if err := html.Render(b, n); err != nil {
		return "", err
	}
	return b.String(), nil
}
