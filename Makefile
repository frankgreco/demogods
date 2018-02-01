ALL_SRC 	= $(shell find . -name "*.go" | grep -v -e vendor)
BINARY		=	$(shell echo $${PWD\#\#*/})
PASS			=	$(shell printf "\033[32mPASS\033[0m")
FAIL			=	$(shell printf "\033[31mFAIL\033[0m")
COLORIZE	=	sed ''/PASS/s//$(PASS)/'' | sed ''/FAIL/s//$(FAIL)/''
PACKAGES 	=	$(shell go list ./... | grep -v /vendor/)

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(ALL_SRC) install fmt binary

.PHONY: install
install:
	dep version || (go get github.com/golang/dep/cmd/dep && dep version)
	dep ensure -v -vendor-only # assumes updated Gopkg.lock

.PHONY: fmt
fmt:
	gofmt -e -s -l -w $(ALL_SRC)

.PHONY: binary
binary:
	go build -o ${BINARY}

.PHONY: docker
docker:
	docker build -t frankgreco/demogods:local .

.PHONY: test
test:
	@bash -c "set -e; set -o pipefail; go test -v -race $(PACKAGES) | $(COLORIZE)"

.PHONY: lint
lint:
	@go vet $(PACKAGES)
	@cat /dev/null > lint.log
	@$(foreach pkg, $(PACKAGES), golint $(pkg) >> lint.log || true;)
	@[ ! -s "lint.log" ] || (echo "Lint Failures" | cat - lint.log && false)
	@gofmt -e -s -l $(ALL_SRC) > fmt.log
	@[ ! -s "fmt.log" ] || (echo "Go Fmt Failures, run 'make fmt'" | cat - fmt.log && false)

.PHONY: install_ci
install_ci: install
	go get github.com/golang/lint/golint
