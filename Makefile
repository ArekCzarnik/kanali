ALL_SRC := $(shell find . -name "*.go" | grep -v -e vendor \
        -e ".*/\..*" \
        -e ".*/_.*" \
        -e ".*/mocks.*")

FILES = $(shell go list ./... | grep -v /vendor/)
ALL_PACKAGES := $(shell glide novendor)
SRC_PACKAGES := $(shell glide novendor | grep -v -e test)

BINARY=kanali
RACE=-race
GOTEST=go test -v $(RACE)
GOLINT=golint
GOVET=go vet
GOINSTALL=go install $(RACE)
GOFMT=gofmt
FMT_LOG=fmt.log
LINT_LOG=lint.log

PASS=$(shell printf "\033[32mPASS\033[0m")
FAIL=$(shell printf "\033[31mFAIL\033[0m")
COLORIZE=sed ''/PASS/s//$(PASS)/'' | sed ''/FAIL/s//$(FAIL)/''

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(ALL_SRC) unit_test fmt lint

.PHONY: install
install:
	dep version || go get github.com/golang/dep/cmd/dep
	dep ensure -v -vendor-only # assumes updated Gopkg.lock

.PHONY: binary
binary:
	CGO_ENABLED=1 ./scripts/binary.sh $(VERSION)

.PHONY: fmt
fmt:
	$(GOFMT) -e -s -l -w $(ALL_SRC)
	./scripts/updateLicenses.sh

.PHONY: cover
cover:
	./scripts/cover.sh $(shell go list $(ALL_PACKAGES))
	go tool cover -html=cover.out -o cover.html

.PHONY: unit_test
unit_test:
	bash -c "set -e; set -o pipefail; $(GOTEST) $(SRC_PACKAGES) | $(COLORIZE)"

.PHONY: e2e_test
e2e_test:
	bash -c "set -e; set -o pipefail; $(GOTEST) ./test/e2e | $(COLORIZE)"

.PHONY: lint
lint:
	@$(GOVET) $(ALL_PACKAGES)
	@cat /dev/null > $(LINT_LOG)
	@$(foreach pkg, $(ALL_PACKAGES), $(GOLINT) $(pkg) >> $(LINT_LOG) || true;)
	@[ ! -s "$(LINT_LOG)" ] || (echo "Lint Failures" | cat - $(LINT_LOG) && false)
	@$(GOFMT) -e -s -l $(ALL_SRC) > $(FMT_LOG)
	@[ ! -s "$(FMT_LOG)" ] || (echo "Go Fmt Failures, run 'make fmt'" | cat - $(FMT_LOG) && false)

.PHONY: install_ci
install_ci: install
	go get github.com/wadey/gocovmerge
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover
	go get github.com/golang/lint/golint

.PHONY: test_ci
test_ci:
	@./scripts/cover.sh $(shell go list $(SRC_PACKAGES))
	make lint

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi