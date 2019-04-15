export PATH := $(PATH):/usr/lib/go-1.10/bin

BIN_DIR = $(DESTDIR)/usr/bin
INSTALL = install -m 0644
INSTALL_PROGRAM = install

all: build

depsinit:
	cd $(CURDIR)/src/ec2-cost-analysis && GOPATH=$(CURDIR) dep init

deps:
	cd $(CURDIR)/src/ec2-cost-analysis && GOPATH=$(CURDIR) dep ensure -update

build:
	GOPATH=$(CURDIR) go build -o bin/ec2-cost-analysis ec2-cost-analysis

format:
	GOPATH=$(CURDIR) gofmt -w src/ec2-cost-analysis/*.go

# https://docs.aws.amazon.com/lambda/latest/dg/lambda-go-how-to-create-deployment-package.html
package:
	GOPATH=$(CURDIR) GOOS=linux go build -o ec2-cost-analysis ec2-cost-analysis/awslambda
	zip ec2-cost-analysis.zip ec2-cost-analysis
