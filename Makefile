include $(GOROOT)/src/Make.$(GOARCH)

TARG=sirdsc
GOFILES=$(wildcard src/*.go)

include $(GOROOT)/src/Make.cmd
