package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var DefaultTemplate string

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
