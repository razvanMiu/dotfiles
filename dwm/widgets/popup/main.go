package main

/*
#cgo pkg-config: x11
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <X11/keysym.h>
#include <stdlib.h>

static int evtype(XEvent *e) { return e->type; }
static int button_x(XEvent *e) { return e->xbutton.x; }
static int button_y(XEvent *e) { return e->xbutton.y; }
static KeySym key_sym(XEvent *e) { return XLookupKeysym(&e->xkey, 0); }
static void set_override_redirect(Display *d, Window w) {
	XSetWindowAttributes attrs;
	attrs.override_redirect = True;
	XChangeWindowAttributes(d, w, CWOverrideRedirect, &attrs);
}
static void set_class_hint(Display *d, Window w, char *class_name) {
	XClassHint hint;
	hint.res_name = class_name;
	hint.res_class = class_name;
	XSetClassHint(d, w, &hint);
}
static void draw_string(Display *d, Drawable w, GC gc, int x, int y, char *s, int n) {
	XDrawString(d, w, gc, x, y, s, n);
}
*/
import "C"

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"
)

const (
	base     = 0x24273a
	surface0 = 0x363a4f
	surface1 = 0x494d64
	text     = 0xcad3f5
	subtext0 = 0xa5adcb
	mauve    = 0xc6a0f6
	green    = 0xa6da95
	red      = 0xed8796
	blue     = 0x8aadf4
)

type app struct {
	dpy           *C.Display
	win           C.Window
	gc            C.GC
	w, h          int
	panel         string
	visible       bool
	calendarMonth time.Time
}

func main() {
	panel := flag.String("panel", "calendar", "initial panel: calendar or audio")
	width := flag.Int("w", 420, "window width in pixels")
	height := flag.Int("h", 360, "window height in pixels")
	x := flag.Int("x", -1, "window x position; default right aligned")
	y := flag.Int("y", 28, "window y position")
	managed := flag.Bool("managed", false, "let dwm manage the window instead of using override_redirect")
	ipc := flag.Bool("ipc", true, "listen for dwm-popupctl commands on a Unix socket")
	hidden := flag.Bool("hidden", false, "start hidden and wait for IPC show/toggle commands")
	flag.Parse()

	dpy := C.XOpenDisplay(nil)
	if dpy == nil {
		fmt.Fprintln(os.Stderr, "dwm-popup: cannot open X display")
		os.Exit(1)
	}
	defer C.XCloseDisplay(dpy)

	screen := C.XDefaultScreen(dpy)
	root := C.XRootWindow(dpy, screen)
	sw := int(C.XDisplayWidth(dpy, screen))
	if *x < 0 {
		*x = sw - *width - 12
	}

	win := C.XCreateSimpleWindow(dpy, root, C.int(*x), C.int(*y), C.uint(*width), C.uint(*height), 0, base, base)
	if !*managed {
		C.set_override_redirect(dpy, win)
	}
	name := C.CString("dropdown-widgets-popup")
	defer C.free(unsafe.Pointer(name))
	C.XStoreName(dpy, win, name)
	C.set_class_hint(dpy, win, name)

	C.XSelectInput(dpy, win, C.ExposureMask|C.ButtonPressMask|C.KeyPressMask|C.StructureNotifyMask)
	if !*hidden {
		C.XMapRaised(dpy, win)
		C.XSetInputFocus(dpy, win, C.RevertToParent, C.CurrentTime)
	}

	gc := C.XCreateGC(dpy, C.Drawable(win), 0, nil)
	defer C.XFreeGC(dpy, gc)

	a := app{
		dpy:           dpy,
		win:           win,
		gc:            gc,
		w:             *width,
		h:             *height,
		panel:         normalizePanel(*panel),
		visible:       !*hidden,
		calendarMonth: firstOfMonth(time.Now()),
	}
	if a.visible {
		a.draw()
	}

	commands := make(chan command)
	if *ipc {
		cleanup, err := startIPC(commands)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dwm-popup: IPC disabled: %v\n", err)
		} else {
			defer cleanup()
			fmt.Fprintf(os.Stderr, "dwm-popup: listening on %s\n", socketPath())
		}
	}

	for {
		for C.XPending(dpy) > 0 {
			var ev C.XEvent
			C.XNextEvent(dpy, &ev)
			switch C.evtype(&ev) {
			case C.Expose:
				a.draw()
			case C.ButtonPress:
				a.click(int(C.button_x(&ev)), int(C.button_y(&ev)))
			case C.KeyPress:
				if !a.key(C.key_sym(&ev)) {
					return
				}
			}
		}
		select {
		case cmd := <-commands:
			if !a.handle(cmd) {
				return
			}
		case <-time.After(16 * time.Millisecond):
		}
	}
}

func normalizePanel(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "audio" {
		return s
	}
	return "calendar"
}

func firstOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

func (a *app) click(x, y int) {
	if y >= 12 && y <= 44 {
		switch {
		case x >= 12 && x < 128:
			a.panel = "calendar"
		case x >= 136 && x < 236:
			a.panel = "audio"
		}
		a.draw()
	}
}

func (a *app) key(ks C.KeySym) bool {
	if ks == C.XK_Escape || ks == C.XK_q {
		a.hide()
		return true
	}
	if ks == C.XK_h {
		a.calendarMonth = a.calendarMonth.AddDate(0, -1, 0)
		a.draw()
	}
	if ks == C.XK_l {
		a.calendarMonth = a.calendarMonth.AddDate(0, 1, 0)
		a.draw()
	}
	return true
}

func (a *app) handle(cmd command) bool {
	message := "ok"
	switch cmd.action {
	case "show":
		a.show(cmd.panel)
	case "hide":
		a.hide()
	case "toggle":
		if a.visible && a.panel == cmd.panel {
			a.hide()
		} else {
			a.show(cmd.panel)
		}
	case "quit":
		message = "bye"
		if cmd.resp != nil {
			cmd.resp <- message
		}
		return false
	default:
		message = "error: unknown command"
	}
	if cmd.resp != nil {
		cmd.resp <- message
	}
	return true
}

func (a *app) show(panel string) {
	a.panel = normalizePanel(panel)
	if !a.visible {
		C.XMapRaised(a.dpy, a.win)
		a.visible = true
	} else {
		C.XRaiseWindow(a.dpy, a.win)
	}
	C.XSetInputFocus(a.dpy, a.win, C.RevertToParent, C.CurrentTime)
	a.draw()
}

func (a *app) hide() {
	if !a.visible {
		return
	}
	C.XUnmapWindow(a.dpy, a.win)
	C.XFlush(a.dpy)
	a.visible = false
}

func (a *app) draw() {
	if !a.visible {
		return
	}
	a.rect(0, 0, a.w, a.h, base)
	a.rect(0, 0, a.w, 56, surface0)
	a.tab(12, 12, 116, 32, "Calendar", a.panel == "calendar")
	a.tab(136, 12, 100, 32, "Audio", a.panel == "audio")
	a.text(a.w-108, 33, subtext0, "Esc/q closes")
	if a.panel == "audio" {
		a.drawAudio()
	} else {
		a.drawCalendar()
	}
	C.XFlush(a.dpy)
}

func (a *app) tab(x, y, w, h int, label string, active bool) {
	bg := uint64(surface1)
	fg := uint64(subtext0)
	if active {
		bg = mauve
		fg = base
	}
	a.rect(x, y, w, h, bg)
	a.text(x+14, y+21, fg, label)
}

func (a *app) drawCalendar() {
	month := a.calendarMonth
	a.text(24, 88, text, month.Format("January 2006"))
	a.text(250, 88, subtext0, "h/l month")

	gridX, gridY := 24, 118
	cellW := max(40, (a.w-48)/7)
	cellH := 34
	for i, day := range []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"} {
		a.text(gridX+i*cellW+8, gridY, mauve, day)
	}
	weeks := calendarGrid(month)
	today := time.Now()
	for r, week := range weeks {
		for c, d := range week {
			x := gridX + c*cellW
			y := gridY + 18 + r*cellH
			if d.Month() != month.Month() {
				a.rect(x+2, y-14, cellW-4, cellH-4, base)
				a.text(x+10, y+6, surface1, fmt.Sprintf("%2d", d.Day()))
				continue
			}
			if sameDay(d, today) {
				a.rect(x+2, y-14, cellW-4, cellH-4, mauve)
				a.text(x+10, y+6, base, fmt.Sprintf("%2d", d.Day()))
			} else {
				a.rect(x+2, y-14, cellW-4, cellH-4, surface0)
				a.text(x+10, y+6, text, fmt.Sprintf("%2d", d.Day()))
			}
		}
	}

	a.text(24, a.h-64, subtext0, "Agenda prototype")
	a.text(24, a.h-40, text, "• Pixel layout: no terminal cells, no Kitty")
}

func (a *app) drawAudio() {
	a.text(24, 88, text, "Audio")
	a.meter(24, 120, a.w-48, 26, 0.62, "Output  Built-in Audio", green)
	a.meter(24, 176, a.w-48, 26, 0.18, "Input   Microphone", blue)
	a.meter(24, 232, a.w-48, 26, 0.74, "App     Browser", mauve)
	a.text(24, a.h-40, subtext0, "Placeholder only; no pactl writes yet")
}

func (a *app) meter(x, y, w, h int, pct float64, label string, color uint64) {
	a.text(x, y-10, text, label)
	a.rect(x, y, w, h, surface0)
	a.rect(x, y, int(float64(w)*pct), h, color)
	a.text(x+w-48, y+18, base, fmt.Sprintf("%2.0f%%", pct*100))
}

func calendarGrid(month time.Time) [6][7]time.Time {
	start := firstOfMonth(month)
	offset := (int(start.Weekday()) + 6) % 7 // Monday first.
	cursor := start.AddDate(0, 0, -offset)
	var grid [6][7]time.Time
	for r := range grid {
		for c := range grid[r] {
			grid[r][c] = cursor
			cursor = cursor.AddDate(0, 0, 1)
		}
	}
	return grid
}

func sameDay(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func (a *app) rect(x, y, w, h int, color uint64) {
	if w <= 0 || h <= 0 {
		return
	}
	C.XSetForeground(a.dpy, a.gc, C.ulong(color))
	C.XFillRectangle(a.dpy, C.Drawable(a.win), a.gc, C.int(x), C.int(y), C.uint(w), C.uint(h))
}

func (a *app) text(x, y int, color uint64, s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.XSetForeground(a.dpy, a.gc, C.ulong(color))
	C.draw_string(a.dpy, C.Drawable(a.win), a.gc, C.int(x), C.int(y), cs, C.int(len(s)))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
