package main

import (
	"log"
	"os/exec"
	"regexp"
	"strings"
)

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
