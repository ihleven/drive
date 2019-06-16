SHELL := /bin/bash

VERSION         :=      0.7
BINARY_NAME     :=      cloud11



default: backend

backend:
	gin --all

serve:
	cd vue && echo "entering vue subdir" && \
    yarn serve

build:
	cd vue && yarn build
	go build -o cloud11


deploy:
	GOOS=linux GOARCH=amd64 go build -o drive-unix
	ssh ihle@web569.webfaction.com 'rm -f ~ihle/webapps/drive/drive-unix'
	scp drive-unix ihle@web569.webfaction.com:~ihle/webapps/drive/
	ssh ihle@web569.webfaction.com 'rm -rf ~ihle/webapps/drive/_static'
	scp -r _static ihle@web569.webfaction.com:~ihle/webapps/drive/

clean:
	rm -f gin-bin
	rm -rf _static
	bash -c "[ -e $(BINARY_NAME) ] && rm $(BINARY_NAME) || true"
	bash -c "[ -e $(BINARY_NAME)-unix ] && rm $(BINARY_NAME)-unix || true"


# Keeping `./main.go` with just a `cli` and `./lib/*.go` with actual 
# logic, `tests` usually reside under `./lib` (or some other subdirectories).
#
# By using the `./...` notation, all the non-vendor packages are going
# to be tested if they have test files.
test:
	go test ./... -v


# Just like `test`, formatting what matters. As `main.go` is in the root,
# `go fmt` the root package. Then just `cd` to what matters to you (`vendor`
# doesn't matter).
#
# By using the `./...` notation, all the non-vendor packages are going
# to be formatted (including test files).
fmt:
	go fmt ./... -v



# This is pretty much an optional thing that I tend always to include.
#
# Goreleaser is a tool that allows anyone to integrate a binary releasing
# process to their pipelines. 
#
# Here in this target With just a simple `make release` you can have a 
# `tag` created in GitHub with multiple builds if you wish. 
#
# See more at `gorelease` github repo.
release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: install test fmt release

