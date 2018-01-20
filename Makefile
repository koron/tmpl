build:
	go build -i -v ./cmd/tmpl

tags:
	gotags -f tags -R .

deps-list:
	@go list -f '{{join .Imports "\n"}}' ./... | sort -u | grep -v `go list`

deps-update:
	@go list -f '{{join .Imports "\n"}}' ./... | sort -u | grep -v `go list` | xargs go get -u -d -v

clean:
	rm -f tags
	rm -f tmpl

.PHONY: build tags deps-list deps-update clean
