default: previewdocs

previewdocs: previewdocs.go
	go build -v -o previewdocs
