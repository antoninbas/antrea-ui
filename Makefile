GO                 ?= go
BINDIR := $(CURDIR)/bin
GOMOCK_VERSION := v1.6.0

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


.mockgen-bin:
	@echo "===> Installing Mockgen <==="
	GOBIN=$(CURDIR)/$@ $(GO) install github.com/golang/mock/mockgen@$(GOMOCK_VERSION)


.PHONY: generate
generate: .mockgen-bin
	PATH=$(CURDIR)/.mockgen-bin:$$PATH $(GO) generate antrea.io/antrea-ui/...


.PHONY: clean
clean:
	rm -rf bin
	rm -rf .golangci-bin
	rm -rf .mockgen-bin

.PHONY: build
build:
	docker build -t antrea/antrea-ui-frontend -f build/Dockerfile.frontend .
	docker build -t antrea/antrea-ui-backend -f build/Dockerfile.backend .
