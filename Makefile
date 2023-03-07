GO                 ?= go
BINDIR := $(CURDIR)/bin

all: build

.PHONY: bin
bin:
	GOBIN=$(BINDIR) $(GO) install antrea.io/antrea-ui/...

.PHONY: test
test:
	$(GO) test -v ./...

# code linting
.golangci-bin:
	@echo "===> Installing Golangci-lint <==="
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $@ v1.48.0

.PHONY: golangci
golangci: .golangci-bin
	@echo "===> Running golangci <==="
	@GOOS=linux .golangci-bin/golangci-lint run -c .golangci.yml

.PHONY: golangci-fix
golangci-fix: .golangci-bin
	@echo "===> Running golangci-fix <==="
	@GOOS=linux .golangci-bin/golangci-lint run -c .golangci.yml --fix

.PHONY: clean
clean:
	rm -rf bin
	rm -rf .golangci-bin

.PHONY: build
build:
	docker build -t antrea/antrea-ui-frontend -f build/Dockerfile.frontend .
	docker build -t antrea/antrea-ui-backend -f build/Dockerfile.backend .
