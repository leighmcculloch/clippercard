release:
	git tag -a $(VERSION) -m "$(VERSION)"
	git push origin $(VERSION)
	goreleaser --rm-dist --parallelism 1 || (git tag -d $(VERSION) && git push --delete origin $(VERSION))

setup:
	go get github.com/goreleaser/goreleaser
