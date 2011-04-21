include $(GOROOT)/src/Make.inc

TARG=sirdsc
GOFILES=$(wildcard src/*.go)

include $(GOROOT)/src/Make.cmd
