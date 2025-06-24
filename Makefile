#!/usr/bin/make -f
# set variables
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%H')
ifeq (,$(VERSION))
  VERSION := $(shell git describe --tags)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif
LEDGER_ENABLED ?= true
COSMOS_SDK_VERSION := $(shell go list -m github.com/cosmos/cosmos-sdk | sed 's:.* ::')
CMT_VERSION := $(shell go list -m github.com/cometbft/cometbft | sed 's:.* ::')
DOCKER := $(shell which docker)

# process build tags
build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
	ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
   endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=echofi \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=echofid \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/cometbft/cometbft/version.TMCoreSemVer=$(CMT_VERSION)

ifeq ($(LINK_STATICALLY),true)
  ldflags += -linkmode=external -extldflags "-Wl,-z,muldefs -static"
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

###############################################################################
###                                  Build                                  ###
###############################################################################

verify:
	@echo "🔎 - Verifying Dependencies ..."
	@go mod verify > /dev/null 2>&1
	@go mod tidy
	@echo "✅ - Verified dependencies successfully!"
	@echo ""

go-cache: verify
	@echo "📥 - Downloading and caching dependencies..."
	@go mod download
	@echo "✅ - Downloaded and cached dependencies successfully!"
	@echo ""

install: go-cache
	@echo "🔄 - Installing Echofi..."
	@go install $(BUILD_FLAGS) -mod=readonly ./cmd/echofid
	@echo "✅ - Installed Echofi successfully! Run it using 'echofid'!"
	@echo ""
	@echo "====== Install Summary ======"
	@echo "Echofi: $(VERSION)"
	@echo "Cosmos SDK: $(COSMOS_SDK_VERSION)"
	@echo "Comet: $(CMT_VERSION)"
	@echo "============================="

build: go-cache
	@echo "🔄 - Building Echofi..."
	@if [ "$(OS)" = "Windows_NT" ]; then \
		GOOS=windows GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o bin/echofid.exe ./cmd/echofid; \
	else \
		go build -mod=readonly $(BUILD_FLAGS) -o bin/echofid ./cmd/echofid; \
	fi
	@echo "✅ - Built Echofi successfully! Run it using './bin/echofid'!"
	@echo ""
	@echo "====== Install Summary ======"
	@echo "Echofi: $(VERSION)"
	@echo "Cosmos SDK: $(COSMOS_SDK_VERSION)"
	@echo "Comet: $(CMT_VERSION)"
	@echo "============================="

###############################################################################
###                                Test                                 	###
###############################################################################
TEST_FLAGS := -v

test:
	go test $(TEST_FLAGS) ./...


###############################################################################
###                                Linting                                  ###
###############################################################################
golangci_lint_cmd=golangci-lint
golangci_version=v1.60.1

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run --timeout=10m

###############################################################################
###                             e2e interchain test                         ###
###############################################################################

ictest-basic: rm-testcache
	cd interchaintest && go test -race -v -run TestBasicEchofiStart .

rm-testcache:
	go clean -testcache