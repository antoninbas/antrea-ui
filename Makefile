GO                 ?= go
BINDIR := $(CURDIR)/bin
GOMOCK_VERSION := v1.6.0
GOLANGCI_LINT_VERSION := v1.51.2
GOLANGCI_LINT_BINDIR  := .golangci-bin
GOLANGCI_LINT_BIN     := $(GOLANGCI_LINT_BINDIR)/$(GOLANGCI_LINT_VERSION)/golangci-lint

all: build

.PHONY: bin
bin:
	GOBIN=$(BINDIR) $(GO) install antrea.io/antrea-ui/...

.PHONY: test
test:
	$(GO) test -v ./...

# code linting
$(GOLANGCI_LINT_BIN):
	@echo "===> Installing Golangci-lint <==="
	@rm -rf $(GOLANGCI_LINT_BINDIR)/* # delete old versions
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCI_LINT_BINDIR)/$(GOLANGCI_LINT_VERSION) $(GOLANGCI_LINT_VERSION)

.PHONY: golangci
golangci: $(GOLANGCI_LINT_BIN)
	@echo "===> Running golangci <==="
	@GOOS=linux $(GOLANGCI_LINT_BIN) run -c .golangci.yml

.PHONY: golangci-fix
golangci-fix: $(GOLANGCI_LINT_BIN)
	@echo "===> Running golangci-fix <==="
	@GOOS=linux $(GOLANGCI_LINT_BIN) run -c .golangci.yml --fix

.mockgen-bin:
	@echo "===> Installing Mockgen <==="
	GOBIN=$(CURDIR)/$@ $(GO) install github.com/golang/mock/mockgen@$(GOMOCK_VERSION)


.PHONY: generate
generate: .mockgen-bin
	PATH=$(CURDIR)/.mockgen-bin:$$PATH $(GO) generate antrea.io/antrea-ui/...


.PHONY: clean
clean:
	rm -rf bin
	rm -rf $(GOLANGCI_LINT_BINDIR)
	rm -rf .mockgen-bin

.PHONY: build
build:
	docker build -t antrea/antrea-ui-frontend -f build/Dockerfile.frontend .
	docker build -t antrea/antrea-ui-backend -f build/Dockerfile.backend .
