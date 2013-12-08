package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func fetchAndRenderDoc(user, repo, doc string, defaultTemplate string) (string, error) {
	template := make(chan string)
	go func() {
		buf, err := ioutil.ReadFile("docs/template.html")
		if err != nil {
			template <- defaultTemplate
			return
		}
		template <- string(buf)
	}()

	// https://github.com/github/markup/blob/master/lib/github/markups.rb#L1
	mdExts := map[string]bool{
		".md":        true,
		".mkdn":      true,
		".mdwn":      true,
		".mdown":     true,
		".markdown":  true,
		".litcoffee": true,
	}
	if ok, _ := mdExts[path.Ext(doc)]; !ok {
		doc += ".md"
	}

	buf, err := ioutil.ReadFile("docs/" + doc)
	if err != nil {
		return "# Page not found", err
	}

	bodyStr := string(buf)

	url := "https://api.github.com/markdown/raw"
	if os.Getenv("ACCESS_TOKEN") != "" {
		url += "?access_token=" + os.Getenv("ACCESS_TOKEN")
	}
	resp, err := http.Post(url, "text/x-markdown", strings.NewReader(bodyStr))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Fix relative links
	bodyStr, err = fixRelativeLinks(doc, repo, string(body))
	if err != nil {
		return "", err
	}

	output := strings.Replace(<-template, "{{CONTENT}}", bodyStr, 1)
	if user != "" {
		output = strings.Replace(output, "{{NAME}}", repo, -1)
	}
	if repo != "" {
		output = strings.Replace(output, "{{USER}}", user, -1)
	}
	return output, nil
}
