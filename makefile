# this makefile used in console environment.
# copy this file to project base directory.

# SET BIN NAME BY USER
binName:=pgo2-demo

goBin:=go

######## DO NOT CHANGE THE FLOWING CONTENT ########

# absolute path of makefile
mkPath:=$(abspath $(firstword $(MAKEFILE_LIST)))

# absolute base directory of project
baseDir:=$(strip $(patsubst %/, %, $(dir $(mkPath))))

binDir:=$(baseDir)/bin
pkgDir:=$(baseDir)/pkg

.PHONY: start stop build update pgo2 init

start: build
	$(binDir)/$(binName)

stop:
	-killall $(binName)

build:
	[ -d $(binDir) ] || mkdir $(binDir)
	$(goBin) build -o $(binDir)/$(binName) $(baseDir)/cmd/$(binName)/main.go

update:
	cd $(baseDir) && $(goBin) mod get

install:
	cd $(baseDir) && $(goBin) mod download

pgo2:
	cd $(baseDir) && $(goBin) get -u github.com/pinguo/pgo2

vendor:
	cd $(baseDir) && $(goBin) mod vendor

init:
	[ -d $(baseDir)/cmd ] || mkdir $(baseDir)/cmd
	[ -d $(baseDir)/runtime ] || mkdir $(baseDir)/runtime
	[ -d $(baseDir)/cmd/$(binName) ] || mkdir $(baseDir)/cmd/$(binName)
	[ -d $(baseDir)/web ] || mkdir $(baseDir)/web
	[ -d $(baseDir)/web/template ] || mkdir $(baseDir)/web/template
	[ -d $(baseDir)/web/static ] || mkdir $(baseDir)/web/static
	[ -d $(baseDir)/configs ] || mkdir $(baseDir)/configs
	[ -d $(baseDir)/assets ] || mkdir $(baseDir)/assets
	[ -d $(baseDir)/build ] || mkdir $(baseDir)/build
	[ -d $(pkgDir) ] || mkdir $(pkgDir)
	[ -d $(pkgDir)/command ] || mkdir $(pkgDir)/command
	[ -d $(pkgDir)/controller ] || mkdir $(pkgDir)/controller
	[ -d $(pkgDir)/lib ] || mkdir $(pkgDir)/lib
	[ -d $(pkgDir)/model ] || mkdir $(pkgDir)/model
	[ -d $(pkgDir)/service ] || mkdir $(pkgDir)/service
	[ -d $(pkgDir)/struct ] || mkdir $(pkgDir)/struct
	[ -d $(pkgDir)/test ] || mkdir $(pkgDir)/test
	[ -f $(pkgDir)/go.mod ] || (cd $(baseDir) && $(goBin) mod init $(binName))

help:
	@echo "make start       build and start $(binName)"
	@echo "make stop        stop process $(binName)"
	@echo "make build       build $(binName)"
	@echo "make update      go mod get"
	@echo "make install     go mod download"
	@echo "make pgo2        go mod get -u pgo2"
	@echo "make init        init project"
	@echo "make vendor      go mod vendor"
