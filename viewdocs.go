package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"code.google.com/p/vitess/go/cache"
)

const CacheCapacity = 256 * 1024 * 1024 // 256MB
const CacheTTL = 60                     // raw.github.com cache TTL is ~120

var DefaultTemplate string

type CacheValue struct {
	Value     string
	CreatedAt int64
}

func (cv *CacheValue) Size() int {
	return len(cv.Value)
}

func parseRequest(r *http.Request) (doc string, err error) {
	path := strings.Split(r.RequestURI, "/")

	if len(path) == 1 || (len(path) == 2 && strings.HasSuffix(r.RequestURI, "/")) {
		doc = "index"
	} else {
		doc = strings.Join(path[1:], "/")
		if strings.HasSuffix(doc, "/") {
			doc = doc[:len(doc)-1]
		}
	}
	return
}

func fetchAndRenderDoc(doc string) (string, error) {
	template := make(chan string)
	go func() {
		buf, err := ioutil.ReadFile("docs/template.html")
		if err != nil {
			template <- DefaultTemplate
			return
		}
		template <- string(buf)
	}()
	buf, err := ioutil.ReadFile("docs/" + doc + ".md")
	if err != nil {
		return "# Page not found", err
	}

	bodyStr := string(buf)
	resp, err := http.Post("https://api.github.com/markdown/raw?access_token="+os.Getenv("ACCESS_TOKEN"), "text/x-markdown", strings.NewReader(bodyStr))
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	output := strings.Replace(<-template, "{{CONTENT}}", string(body), 1)
	// output = strings.Replace(output, "{{NAME}}", repo, -1)
	// output = strings.Replace(output, "{{USER}}", user, -1)
	return output, nil
}

func main() {
	if os.Getenv("ACCESS_TOKEN") == "" {
		// TODO: Add direct link to Development section of the README
		log.Fatal("ACCESS_TOKEN was not found!")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	lru := cache.NewLRUCache(CacheCapacity)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/favicon.ico" {
			return
		}
		switch r.Method {
		case "GET":
			doc, err := parseRequest(r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			log.Printf("Building docs for '%s'", doc)
			key := doc
			value, ok := lru.Get(key)
			var output string
			if !ok {
				output, err = fetchAndRenderDoc(doc)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(err.Error()))
					return
				}
				lru.Set(key, &CacheValue{output, time.Now().Unix()})
				log.Println("CACHE MISS:", key, lru.StatsJSON())
			} else {
				output = value.(*CacheValue).Value
				if time.Now().Unix()-value.(*CacheValue).CreatedAt > CacheTTL {
					lru.Delete(key)
				}
			}
			w.Write([]byte(output))
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
