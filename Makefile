TARGET=sirdsc

SRCDIR=src

PREFIX=/usr

.PHONY: all install clean fmt

all: $(TARGET)

$(TARGET) clean fmt:
	gomake -C src $@

install: all
	install -m 755 -D "$(SRCDIR)/$(TARGET)" "$(PREFIX)/bin/$(TARGET)"
