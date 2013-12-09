package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"
)

var DefaultTemplate string

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

func fixRelativeLinks(doc, repo, body string) (string, error) {
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

func fetchAndRenderDoc(user, repo, doc string) (string, error) {
	template := make(chan string)
	go func() {
		buf, err := ioutil.ReadFile("docs/template.html")
		if err != nil {
			template <- DefaultTemplate
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

	output := strings.Replace(<-template, "{{CONTENT}}", string(body), 1)
	output = strings.Replace(output, "{{NAME}}", repo, -1)
	output = strings.Replace(output, "{{USER}}", user, -1)

	// Fix relative links
	output, err = fixRelativeLinks(doc, repo, output)
	if err != nil {
		return "", err
	}

	return output, nil
}

func grabUserAndRepo() (user, repo string) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Printf("Error fetching github user and repository\nERROR: %s", err)
	} else {
		output := strings.Trim(string(out), "\n")
		reg := regexp.MustCompile(`([^:/]+)/([\w.-]+)\.git$`)
		matches := reg.FindStringSubmatch(output)

		if len(matches) > 0 {
			user = matches[1]
			repo = matches[2]
		} else {
			log.Fatalf("Unable to parse your GitHub user and repository from '%s'. Please open an issue on https://github.com/fgrehm/previewdocs", output)
		}
	}

	return
}

func main() {
	if os.Getenv("ACCESS_TOKEN") == "" {
		log.Println("WARNING: ACCESS_TOKEN was not found, you'll be subject to GitHub's Rate Limiting of 60 requests per hour. " +
			"Please read http://developer.github.com/v3/#rate-limiting for more information")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	resp, err := http.Get("https://raw.github.com/progrium/viewdocs/master/docs/template.html")
	if err != nil || resp.StatusCode == 404 {
		log.Fatal("Unable to fetch default template")
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	DefaultTemplate = string(body)

	user, gitHubRepo := grabUserAndRepo()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/" {
			http.Redirect(w, r, "/" + gitHubRepo, 301)
			return
		}
		if r.RequestURI == "/favicon.ico" {
			return
		}
		switch r.Method {
		case "GET":
			requestedRepo, doc := parseRequest(r)
			if requestedRepo != gitHubRepo {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("Invalid repository '" + requestedRepo +  "'"))
				return
			}
			log.Printf("Building docs for '%s'", doc)
			output, err := fetchAndRenderDoc(user, gitHubRepo, doc)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(output))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
