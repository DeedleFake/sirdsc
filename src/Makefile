include $(GOROOT)/src/Make.inc

TARG=sirdsc
GOFILES=$(wildcard *.go)

include $(GOROOT)/src/Make.cmd

.PHONY: fmt

fmt:
	gofmt -w $(GOFILES)
