package main

import (
	"net/http"
	"strings"
)

func parseRequest(r *http.Request) (repo, doc string) {
	path := strings.Split(r.RequestURI, "/")
	repo = path[1]

	if len(path) < 3 || (len(path) == 3 && strings.HasSuffix(r.RequestURI, "/")) {
		doc = "index.md"
	} else {
		doc = strings.Join(path[2:], "/")
		if strings.HasSuffix(doc, "/") {
			doc = doc[:len(doc)-1]
		}
	}
	return
}
