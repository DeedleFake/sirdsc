TARGET=sirdsc

SRCDIR=src

PREFIX=/usr

.PHONY: all install clean

all: $(TARGET)

$(TARGET):
	gomake -C src $@

install: all
	install -m 755 -D "$(SRCDIR)/$(TARGET)" "$(PREFIX)/bin/$(TARGET)"

clean:
	gomake -C src $@
