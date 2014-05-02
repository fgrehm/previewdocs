OSARCH="windows/386 windows/amd64 linux/386 linux/amd64 darwin/386 darwin/amd64"

default: previewdocs

previewdocs: previewdocs.go
	go build -v -o previewdocs

gox-toolchain:
	go get github.com/mitchellh/gox
	gox -build-toolchain -osarch=$(OSARCH)

release:
	@test -z '$(version)' && echo 'version parameter not provided to `make`!' && exit 1 || return 0
	rm -rf build/*
	gox -osarch=$(OSARCH) -output="build/previewdocs_{{.OS}}_{{.Arch}}" -verbose || { echo "Did you build the toolchain for gox with 'make gox-toolchain'?"; exit 1; }
	git tag $(version)
	git push && git push --tags
	gh release create -d -a build/ $(version)
