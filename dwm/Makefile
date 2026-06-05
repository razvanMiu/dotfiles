# dwm - dynamic window manager
# See LICENSE file for copyright and license details.

include config.mk

SRC = drw.c dwm.c util.c
OBJ = ${SRC:.c=.o}
GO ?= go
USER_BIN ?= ${HOME}/.local/bin

all: dwm

local: all tuis

status-restart: ${USER_BIN}/dwm-status
	pkill -x dwm-status || true
	mkdir -p $${XDG_STATE_HOME:-$${HOME}/.local/state}/dwm
	${USER_BIN}/dwm-status > $${XDG_STATE_HOME:-$${HOME}/.local/state}/dwm/status.log 2>&1 & echo $$! > $${XDG_STATE_HOME:-$${HOME}/.local/state}/dwm/status.pid

tuis: ${USER_BIN}/dwm-widgets ${USER_BIN}/dwm-widgetctl ${USER_BIN}/dwm-status

${USER_BIN}/dwm-widgets: FORCE widgets/panel/main.go widgets/panel/go.mod widgets/panel/go.sum
	mkdir -p ${USER_BIN}
	cd widgets/panel && ${GO} build -o $@ .

${USER_BIN}/dwm-widgetctl: FORCE widgets/panel/cmd/dwm-widgetctl/main.go widgets/panel/go.mod widgets/panel/go.sum
	mkdir -p ${USER_BIN}
	cd widgets/panel && ${GO} build -o $@ ./cmd/dwm-widgetctl

${USER_BIN}/dwm-status: FORCE widgets/status/main.go widgets/status/go.mod
	mkdir -p ${USER_BIN}
	cd widgets/status && ${GO} build -o $@ .

.c.o:
	${CC} -c ${CFLAGS} $<

${OBJ}: config.h config.mk

config.h:
	cp config.def.h $@

dwm: ${OBJ}
	${CC} -o $@ ${OBJ} ${LDFLAGS}

clean:
	rm -f dwm ${OBJ} dwm-${VERSION}.tar.gz

dist: clean
	mkdir -p dwm-${VERSION}
	cp -R LICENSE Makefile README config.def.h config.mk\
		dwm.1 drw.h util.h ${SRC} dwm.png transient.c dwm-${VERSION}
	tar -cf dwm-${VERSION}.tar dwm-${VERSION}
	gzip dwm-${VERSION}.tar
	rm -rf dwm-${VERSION}

install: all
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f dwm ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/dwm
	mkdir -p ${DESTDIR}${MANPREFIX}/man1
	sed "s/VERSION/${VERSION}/g" < dwm.1 > ${DESTDIR}${MANPREFIX}/man1/dwm.1
	chmod 644 ${DESTDIR}${MANPREFIX}/man1/dwm.1

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/dwm\
		${DESTDIR}${MANPREFIX}/man1/dwm.1

FORCE:

.PHONY: all local tuis status-restart clean dist install uninstall FORCE
