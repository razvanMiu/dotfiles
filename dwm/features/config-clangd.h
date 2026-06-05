/* clangd-only preamble for config.h/config.def.h standalone parsing. */
#include <stddef.h>
#include <X11/X.h>
#include <X11/keysym.h>

typedef union {
	int i;
	unsigned int ui;
	float f;
	const void *v;
} Arg;

typedef struct Client Client;
typedef struct Monitor Monitor;

typedef struct {
	const char *class;
	const char *instance;
	const char *title;
	unsigned int tags;
	int isfloating;
	int dropdown;
	int monitor;
} Rule;

typedef struct {
	const char *symbol;
	void (*arrange)(Monitor *);
} Layout;

typedef struct {
	unsigned int mod;
	KeySym keysym;
	void (*func)(const Arg *);
	const Arg arg;
} Key;

typedef struct {
	unsigned int click;
	unsigned int mask;
	unsigned int button;
	void (*func)(const Arg *);
	const Arg arg;
} Button;

typedef struct {
	const char **cmd;
	float wfact;
	float hfact;
	float xfact;
} Dropdown;

typedef struct {
	const char *name;
	void (*func)(const Arg *);
	const Arg arg;
} StatusSegment;

enum { SchemeNorm, SchemeSel };
enum { ClkTagBar, ClkLtSymbol, ClkStatusText, ClkWinTitle, ClkClientWin, ClkRootWin, ClkLast };

void dropdownsetmfact(const Arg *arg);
void focusmon(const Arg *arg);
void focusstack(const Arg *arg);
void incnmaster(const Arg *arg);
void killclient(const Arg *arg);
void monocle(Monitor *m);
void movemouse(const Arg *arg);
void quit(const Arg *arg);
void resizemouse(const Arg *arg);
void setlayout(const Arg *arg);
void spawn(const Arg *arg);
void tag(const Arg *arg);
void tagmon(const Arg *arg);
void tile(Monitor *m);
void toggleaudiowidget(const Arg *arg);
void togglebar(const Arg *arg);
void togglecalendarwidget(const Arg *arg);
void togglefloating(const Arg *arg);
void toggledropdown(const Arg *arg);
void toggletag(const Arg *arg);
void toggleview(const Arg *arg);
void view(const Arg *arg);
void zoom(const Arg *arg);
