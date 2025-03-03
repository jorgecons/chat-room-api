## install the tools needed to run all the others commands
.PHONY: install-tools
install-tools:
	go install github.com/daixiang0/gci@latest;
	go install mvdan.cc/gofumpt@latest;
	brew install mockery && brew cleanup mockery;

.PHONY: format
format:
	gofmt -w .; gofumpt -l -w .; go mod tidy

## To run this command install the following tool https://github.com/daixiang0/gci
.PHONY: sort-imports
sort-imports:
	@gci write --skip-generated -s standard -s default  -s "prefix(github.com/)" -s "prefix(`(head -n 1 ./go.mod | sed 's/^module //')`)" .


.PHONY: mocks
mocks:
	mockery --config mockery.yml

.PHONY: run-all
run-all:
	$(MAKE) mocks;
	$(MAKE) format;
	$(MAKE) sort-imports;

