PKG_NAME=go-strutil

VERSION				:= $(shell git describe --tags --always --dirty="-dev")
DATE				:= $(shell date -u '+%Y-%m-%d-%H%M UTC')
VERSION_FLAGS		:= -ldflags='-X "main.Version=$(VERSION)" -X "main.BuildTime=$(DATE)"'
PLATFORM        	:=$(shell uname -a)
CMD_RM          	:=$(shell which rm)
CMD_CC          	:=$(shell which gcc)
CMD_STRIP       	:=$(shell which strip)
CMD_DIFF        	:=$(shell which diff)
CMD_RM          	:=$(shell which rm)
CMD_BASH        	:=$(shell which bash)
CMD_CP          	:=$(shell which cp)
CMD_AR          	:=$(shell which ar)
CMD_RANLIB      	:=$(shell which ranlib)
CMD_MV          	:=$(shell which mv)
CMD_AWK				:=$(shell which awk)
CMD_SED				:=$(shell which sed)
CMD_TAIL        	:=$(shell which tail)
CMD_FIND        	:=$(shell which find)
CMD_LDD         	:=$(shell which ldd)
CMD_MKDIR       	:=$(shell which mkdir)
CMD_TEST        	:=$(shell which test)
CMD_SLEEP       	:=$(shell which sleep)
CMD_SYNC        	:=$(shell which sync)
CMD_LN          	:=$(shell which ln)
CMD_ZIP        		:=$(shell which zip)
CMD_MD5SUM      	:=$(shell which md5sum)
CMD_READELF     	:=$(shell which readelf)
CMD_GDB         	:=$(shell which gdb)
CMD_FILE        	:=$(shell which file)
CMD_ECHO        	:=$(shell which echo)
CMD_NM          	:=$(shell which nm)
CMD_GO				:=$(shell which go)
CMD_GOLINT			:=$(shell which golint)
CMD_GOIMPORTS		:=$(shell which goimport)
CMD_MAKE2HELP		:=$(shell which make2help)
CMD_GLIDE			:=$(shell which glide)
CMD_GOVER			:=$(shell which gover)
CMD_GOVERALLS		:=$(shell which goveralls)
CMD_CILINT			:=$(shell which golangci-lint)
CMD_CURL			:=$(shell which curl)

PATH_REPORT=report
PATH_RACE_REPORT=$(PKG_NAME).race.report
PATH_CONVER_PROFILE=$(PKG_NAME).coverprofile
PATH_PROF_CPU=$(PKG_NAME).cpu.prof
PATH_PROF_MEM=$(PKG_NAME).mem.prof
PATH_PROF_BLOCK=$(PKG_NAME).block.prof
PATH_PROF_MUTEX=$(PKG_NAME).mutex.prof

VER_GOLANG=$(shell go version | awk '{print $$3}' | sed -e "s/go//;s/\.//g")
GOLANGV18_OVER=$(shell [ "$(VER_GOLANG)" -ge "180" ] && echo 1 || echo 0)
GOMOD_FOUND=$(shell go --help 2>&1 | fgrep "module maintenance" | awk '{print $$1}')
GOMOD_SUPPORT=$(shell [ "$(GOMOD_FOUND)" = "mod" ] && echo 1 || echo 0)

all: clean setup

## Setup Build Environment
setup::
	@$(CMD_ECHO) -e "\033[1;40;32mSetup Build Environment.\033[01;m\x1b[0m"
ifeq ($(GOMOD_SUPPORT),1)
	@$(CMD_GO) mod tidy
	@$(CMD_GO) mod verify
else
	@$(CMD_GO) get github.com/Masterminds/glide
	@$(CMD_GO) get github.com/Songmu/make2help/cmd/make2help
	@$(CMD_GO) get github.com/davecgh/go-spew/spew
	@$(CMD_GO) get github.com/k0kubun/pp
	@$(CMD_GO) get github.com/mattn/goveralls
	@$(CMD_GO) get golang.org/x/tools/cmd/cover
	@$(CMD_GO) get github.com/modocache/gover
	@$(CMD_GO) get github.com/dustin/go-humanize
	@$(CMD_GO) get golang.org/x/lint/golint
	@$(CMD_GO) get github.com/awalterschulze/gographviz
	@GO111MODULE=off $(CMD_GO) get github.com/golang/dep/cmd/dep
endif
	@$(CMD_CURL) -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.16.0
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Build the go-strutil
build::
	@$(CMD_ECHO) -e "\033[1;40;32mBuild the go-strutil.\033[01;m\x1b[0m"
	@$(CMD_GO) build
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Build the go-strutil for development
devbuild::
	@$(CMD_ECHO) -e "\033[1;40;32mBuild the go-strutil.\033[01;m\x1b[0m"
	@$(CMD_GO) build -x -v -gcflags="-N -l" 
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Normal)
lint: setup
	@$(CMD_ECHO) -e "\033[1;40;32mRun a LintChecker (Normal).\033[01;m\x1b[0m"
	@$(CMD_GO) vet $$($(CMD_GLIDE) novendor)
	@for pkg in $$($(CMD_GLIDE) novendor -x); do \
		$(CMD_GOLINT) -set_exit_status $$pkg || exit $$?; \
	done
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run a LintChecker (Strict)
strictlint: setup
	@$(CMD_ECHO) -e "\033[1;40;32mRun a LintChecker (Strict).\033[01;m\x1b[0m"
	@$(CMD_CILINT) run $$($(CMD_GLIDE) novendor)
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test with Data Race Detection
testassert: clean
	@$(CMD_ECHO) -e "\033[1;40;32mRun Go Test.\033[01;m\x1b[0m"
	@$(CMD_GO) test -v -test.parallel 4 -race -run Test_strutils_Assert*
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report of data race detection in $(PATH_REPORT)/doc/$(PATH_RACE_REPORT).pid\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Run Go Test with Data Race Detection
test: clean
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO) -e "\033[1;40;32mRun Go Test.\033[01;m\x1b[0m"
	@GORACE="log_path=$(PATH_REPORT)/doc/$(PATH_RACE_REPORT)" $(CMD_GO) test -tags unittest -v -test.parallel 4 -race -coverprofile=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE)
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report of data race detection in $(PATH_REPORT)/doc/$(PATH_RACE_REPORT).pid\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Send a report of coverage profile to coveralls.io
coveralls::
	@$(CMD_GO) get github.com/mattn/goveralls
	@$(CMD_ECHO) -e "\033[1;40;32mSend a report of coverage profile to coveralls.io.\033[01;m\x1b[0m"
	@$(CMD_GOVERALLS) -coverprofile=$(PATH_REPORT)/raw/$(PATH_CONVER_PROFILE) -service=travis-ci
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Generate a report about coverage
cover: test
	@$(CMD_ECHO) -e "\033[1;40;32mGenerate a report about coverage.\033[01;m\x1b[0m"
	@$(CMD_GO) tool cover -func=$(PATH_CONVER_PROFILE) -o $(PATH_CONVER_PROFILE).txt
	@$(CMD_GO) tool cover -html=$(PATH_CONVER_PROFILE)  -o $(PATH_CONVER_PROFILE).html
	@$(CMD_ECHO) -e "\033[1;40;36mGenerated a report file : $(PATH_CONVER_PROFILE).html\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Profiling
pprof: clean
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO) -e "\033[1;40;32mGenerate profiles.\033[01;m\x1b[0m"
	@$(CMD_ECHO) -e "\033[1;40;33mGenerate a CPU profile.\033[01;m\x1b[0m"
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -cpuprofile=$(PATH_REPORT)/raw/$(PATH_PROF_CPU)
	@$(CMD_ECHO) -e "\033[1;40;33mGenerate a Memory profile.\033[01;m\x1b[0m"
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -memprofile=$(PATH_REPORT)/raw/$(PATH_PROF_MEM)
	@$(CMD_ECHO) -e "\033[1;40;33mGenerate a Block profile.\033[01;m\x1b[0m"
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -blockprofile=$(PATH_REPORT)/raw/$(PATH_PROF_BLOCK)
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_ECHO) -e "\033[1;40;33mGenerate a Mutex profile.\033[01;m\x1b[0m"
	@$(CMD_GO) test -tags unittest -test.parallel 4 -bench . -benchmem -mutexprofile=$(PATH_REPORT)/raw/$(PATH_PROF_MUTEX)
endif
	@$(CMD_MV) -f *.test $(PATH_REPORT)/raw/
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Generate report fo profiling
report: pprof
	@$(CMD_MKDIR) -p $(PATH_REPORT)/raw/ $(PATH_REPORT)/doc/
	@$(CMD_ECHO) -e "\033[1;40;33mGenerate all report in text format.\033[01;m\x1b[0m"
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_CPU) > $(PATH_REPORT)/doc/$(PATH_PROF_CPU).txt
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MEM) > $(PATH_REPORT)/doc/$(PATH_PROF_MEM).txt
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_BLOCK) > $(PATH_REPORT)/doc/$(PATH_PROF_BLOCK).txt
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_GO) tool pprof -text $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MUTEX) > $(PATH_REPORT)/doc/$(PATH_PROF_MUTEX).txt
endif
	@$(CMD_ECHO) -e "\033[1;40;33mGenerate all report in pdf format.\033[01;m\x1b[0m"
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_CPU) > $(PATH_REPORT)/doc/$(PATH_PROF_CPU).pdf
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MEM) > $(PATH_REPORT)/doc/$(PATH_PROF_MEM).pdf
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_BLOCK) > $(PATH_REPORT)/doc/$(PATH_PROF_BLOCK).pdf
ifeq ($(GOLANGV18_OVER),1)
	@$(CMD_GO) tool pprof -pdf $(PATH_REPORT)/raw/$(PKG_NAME).test $(PATH_REPORT)/raw/$(PATH_PROF_MUTEX) > $(PATH_REPORT)/doc/$(PATH_PROF_MUTEX).pdf
endif
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

## Show Help
help::
	@$(CMD_MAKE2HELP) $(MAKEFILE_LIST)

## Clean-up
clean::
	@$(CMD_ECHO) -e "\033[1;40;32mClean-up.\033[01;m\x1b[0m"
	@$(CMD_RM) -rfv *.coverprofile *.swp *.core *.html *.prof *.test *.report ./$(PATH_REPORT)/*
	@$(CMD_ECHO) -e "\033[1;40;36mDone\033[01;m\x1b[0m"

.PHONY: clean cover coveralls help lint pprof report run setup strictlint test
