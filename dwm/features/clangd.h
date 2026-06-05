/*
 * clangd-only preamble for included dwm feature fragments.
 *
 * features/any.c are not standalone translation units; the real build includes
 * them from dwm.c after all dwm types, globals, macros, and config values have
 * been declared. clangd may still parse the opened feature file directly, so
 * .clangd force-includes this lightweight context for editor diagnostics only.
 */
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>
#include <X11/Xlib.h>

#define MAXDROPDOWNS 8
#define LENGTH(X) 8
#define TAGMASK ((1u << 9) - 1)
#define STATUSSEP '\037'
#define MAX(A, B)               ((A) > (B) ? (A) : (B))
#define MIN(A, B)               ((A) < (B) ? (A) : (B))
#define ISDROPDOWN(C)           ((C) && (C)->dropdown >= 0)
#define ISVISIBLE(C)            ((C) && (C)->tags)
#define TEXTW(X)                ((int)sizeof(X) + lrpad)

typedef union {
	int i;
	unsigned int ui;
	float f;
	const void *v;
} Arg;

typedef struct Client Client;
typedef struct Monitor Monitor;
typedef struct Drw Drw;
typedef struct Clr Clr;

typedef struct {
	const char **cmd;
	float wfact;
	float hfact;
	float xfact;
} Dropdown;

typedef struct {
	int dropdown;
} Rule;

typedef struct {
	const char *name;
	void (*func)(const Arg *);
	const Arg arg;
} StatusSegment;

struct Client {
	int dropdown;
	int isfloating;
	int x, y, w, h, bw;
	unsigned int tags;
	Window win;
	Client *next;
	Client *snext;
	Monitor *mon;
};

struct Monitor {
	int wx, wy, ww, wh;
	int seltags;
	unsigned int tagset[2];
	int dropw[MAXDROPDOWNS];
	int droph[MAXDROPDOWNS];
	Client *clients;
	Client *stack;
	Client *sel;
	Monitor *next;
};

enum { SchemeNorm, SchemeSel };

extern char stext[256];
extern int bh, lrpad;
extern Client *nexttiled(Client *c);
extern Clr **scheme;
extern Display *dpy;
extern Drw *drw;
extern Monitor *mons, *selmon;
extern const Dropdown dropdowns[];
extern const StatusSegment statussegments[];
extern const char *statusclickcmd[];

void arrange(Monitor *m);
void attach(Client *c);
void attachstack(Client *c);
void detach(Client *c);
void detachstack(Client *c);
void drw_setscheme(Drw *drw, Clr *scm);
int drw_text(Drw *drw, int x, int y, unsigned int w, unsigned int h, unsigned int lpad, const char *text, int invert);
void focus(Client *c);
void resize(Client *c, int x, int y, int w, int h, int interact);
void setmfact(const Arg *arg);
void spawn(const Arg *arg);
void unfocus(Client *c, int setfocus);
void updatesizehints(Client *c);
