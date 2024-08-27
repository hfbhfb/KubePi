GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOARCH=$(shell go env GOARCH)
GOOS=$(shell go env GOOS )

BASEPATH := $(shell pwd)
BUILDDIR=$(BASEPATH)/dist/usr/local/bin
KUBEPIDIR=$(BASEPATH)/web/kubepi
DASHBOARDDIR=$(BASEPATH)/web/dashboard
TERMINALDIR=$(BASEPATH)/web/terminal
GOTTYDIR=$(BASEPATH)/thirdparty/gotty
MAIN= $(BASEPATH)/cmd/server/main.go
APP_NAME=kubepi-server

tmp: 
	@echo "tmp do nothing"

my_all: 
	@echo "first rsync ; and to newdir "
	#rsync -avh --progress /d/projs/fork8sdir/code-ref/kube-pi /d/tmp/forbuild/
	cd /d/tmp/forbuild/kube-pi;  make clean -f prepare-Makefile;  git pull origin capture; make norunhere324334288

norunhere324334288: my_web_terminal my_web_dashboard my_web_kubepi my_bin_go my_final_docker
	@echo "finish all"

my_final_docker:
	docker build -f Dockerfile-final -t kubeoperator/final-docker-kubepi:master .

my_bin_go:
	docker build -f Dockerfile-go -t kubeoperator/bin-go-kubepi:master .

my_web_kubepi:
	docker build -f Dockerfile-web-kubepi -t kubeoperator/web-kubepi:master .

my_web_dashboard:
	docker build -f Dockerfile-web-dashboard -t kubeoperator/web-dashboard:master .

my_web_terminal:
	docker build -f Dockerfile-web-terminal -t kubeoperator/web-terminal:master .


build_web_kubepi:
	cd $(KUBEPIDIR) && npm install && npm run-script build
build_web_dashboard:
	cd $(DASHBOARDDIR) && npm install && npm run-script build
build_web_terminal:
	cd $(TERMINALDIR) && npm install && npm run-script build

build_web: build_web_kubepi build_web_dashboard build_web_terminal

build_bin:
	GOOS=$(GOOS) GOARCH=$(GOARCH)  $(GOBUILD) -trimpath  -ldflags "-s -w"  -o $(BUILDDIR)/$(APP_NAME) $(MAIN)

build_gotty:
	cd $(GOTTYDIR) && make && mkdir -p  ${BUILDDIR} && mv gotty ${BUILDDIR}

build_all: build_web build_gotty build_bin

build_docker:
	docker build -t kubeoperator/kubepi-server:master .
