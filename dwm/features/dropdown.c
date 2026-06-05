#ifndef DWM_CONTEXT
#define DWM_CONTEXT
#include "../dwm.c"
#else
static time_t dropdownpending[MAXDROPDOWNS];
static int widgetcols[MAXDROPDOWNS];
static int widgetrows[MAXDROPDOWNS];

static Client *
finddropdown(int dropdown)
{
	Client *c;
	Monitor *m;

	if (dropdown < 0 || dropdown >= (int)LENGTH(dropdowns))
		return NULL;
	for (m = mons; m; m = m->next)
		for (c = m->clients; c; c = c->next)
			if (c->dropdown == dropdown)
				return c;
	return NULL;
}

static int
dropdownvisibleon(Client *c, Monitor *m)
{
	return c && c->mon == m && ISVISIBLE(c);
}

static int
dropdownprotected(Client *c)
{
	return c && ISDROPDOWN(c);
}

static unsigned int
dropdowntags(Client *c, Monitor *m)
{
	return ISDROPDOWN(c) ? TAGMASK : m->tagset[m->seltags];
}

static void
cleardropdownpending(Client *c)
{
	if (c && ISDROPDOWN(c) && c->dropdown < MAXDROPDOWNS)
		dropdownpending[c->dropdown] = 0;
}

static void
setruledropdown(Client *c, const Rule *r)
{
	if (r->dropdown >= 0 && r->dropdown < (int)LENGTH(dropdowns))
		c->dropdown = r->dropdown;
}

static void
storedropdownsize(Client *c)
{
	if (!c || !ISDROPDOWN(c) || !c->mon || c->dropdown >= MAXDROPDOWNS)
		return;
	c->mon->dropw[c->dropdown] = c->w;
	c->mon->droph[c->dropdown] = c->h;
}

static void
resizewidgetdropdown(Client *c)
{
	int cols, rows, w, h, maxw, maxh;

	if (!c || !ISDROPDOWN(c) || c->dropdown >= MAXDROPDOWNS)
		return;
	cols = widgetcols[c->dropdown];
	rows = widgetrows[c->dropdown];
	if (cols <= 0 && rows <= 0)
		return;
	updatesizehints(c);
	maxw = MAX(1, c->mon->ww - 2 * c->bw);
	maxh = MAX(1, c->mon->wh - 2 * c->bw);
	w = c->w;
	h = c->h;
	if (cols > 0 && c->incw > 0)
		w = c->basew + cols * c->incw;
	if (rows > 0 && c->inch > 0)
		h = c->baseh + rows * c->inch;
	c->w = MIN(w, maxw);
	c->h = MIN(h, maxh);
}

static void
placedropdown(Client *c)
{
	Monitor *m;
	int dropdown, maxw, maxh, defaultw, defaulth, w, h, x;

	if (!c || !ISDROPDOWN(c) || !(m = c->mon) || c->dropdown >= (int)LENGTH(dropdowns))
		return;

	dropdown = c->dropdown;
	maxw = MAX(1, m->ww - 2 * c->bw);
	maxh = MAX(1, m->wh - 2 * c->bw);
	defaultw = dropdowns[dropdown].wfact == 0.0
		? c->w
		: dropdowns[dropdown].wfact > 1.0
			? (int)dropdowns[dropdown].wfact
			: MAX(1, (int)(maxw * (dropdowns[dropdown].wfact > 0 ? dropdowns[dropdown].wfact : 1.0)));
	defaulth = dropdowns[dropdown].hfact == 0.0
		? c->h
		: dropdowns[dropdown].hfact > 1.0
			? (int)dropdowns[dropdown].hfact
			: MAX(1, (int)(maxh * (dropdowns[dropdown].hfact > 0 ? dropdowns[dropdown].hfact : 0.5)));
	w = MIN((m->dropw[dropdown] && dropdowns[dropdown].wfact != 0.0) ? m->dropw[dropdown] : defaultw, maxw);
	h = MIN((m->droph[dropdown] && dropdowns[dropdown].hfact != 0.0) ? m->droph[dropdown] : defaulth, maxh);
	x = m->wx + (int)((maxw - w) * dropdowns[dropdown].xfact);
	c->isfloating = 1;
	resize(c, x, m->wy, w, h, 0);
}

static void hidevisibledropdowns(int except);

static void
managedropdown(Client *c)
{
	if (!ISDROPDOWN(c))
		return;
	hidevisibledropdowns(c->dropdown);
	cleardropdownpending(c);
	c->mon = selmon;
	c->tags = TAGMASK;
	c->isfloating = 1;
	resizewidgetdropdown(c);
	placedropdown(c);
}

static void
movedropdown(Client *c, Monitor *m)
{
	Monitor *oldmon;

	if (!c || c->mon == m)
		return;
	oldmon = c->mon;
	unfocus(c, 1);
	detach(c);
	detachstack(c);
	c->mon = m;
	attach(c);
	attachstack(c);
	arrange(oldmon);
}

static void
hidedropdown(Client *c, int refocus);

static void
hidevisibledropdowns(int except)
{
	Client *c;
	Monitor *m;

	for (m = mons; m; m = m->next)
		for (c = m->clients; c; c = c->next)
			if (ISDROPDOWN(c) && c->dropdown != except && ISVISIBLE(c))
				hidedropdown(c, 0);
}

static void
showdropdown(Client *c)
{
	if (!c)
		return;
	hidevisibledropdowns(c->dropdown);
	if (c->mon != selmon)
		storedropdownsize(c);
	movedropdown(c, selmon);
	c->tags = TAGMASK;
	c->isfloating = 1;
	resizewidgetdropdown(c);
	placedropdown(c);
	arrange(c->mon);
	focus(c);
	XRaiseWindow(dpy, c->win);
}

static void
showdropdownid(int dropdown)
{
	Client *c;
	Arg arg;

	arg.i = dropdown;
	if ((c = finddropdown(dropdown)))
		showdropdown(c);
	else
		toggledropdown(&arg);
}

static void
setwidgettab(const char *tab)
{
	char cmd[256];
	const char *home;

	home = getenv("HOME");
	if (home && *home)
		snprintf(cmd, sizeof cmd, "%s/.local/bin/dwm-widgetctl %s", home, tab);
	else
		snprintf(cmd, sizeof cmd, "dwm-widgetctl %s", tab);
	(void)system(cmd);
}

static void
togglewidget(const Arg *arg)
{
	char tab[32];
	int cols, dropdown, rows;

	cols = rows = 0;
	if (!arg->v || sscanf(arg->v, "%d %31s %d %d", &dropdown, tab, &cols, &rows) < 2)
		return;
	if (dropdown >= 0 && dropdown < MAXDROPDOWNS) {
		widgetcols[dropdown] = cols;
		widgetrows[dropdown] = rows;
	}
	setwidgettab(tab);
	showdropdownid(dropdown);
}

static void
hidedropdown(Client *c, int refocus)
{
	if (!c)
		return;
	storedropdownsize(c);
	c->tags = 0;
	if (refocus)
		focus(NULL);
	arrange(c->mon);
}

static void
placedropdowns(Monitor *m)
{
	Client *c;

	if (!m)
		return;
	for (c = m->clients; c; c = c->next)
		if (ISDROPDOWN(c) && ISVISIBLE(c))
			placedropdown(c);
}

static void
placedropdownifvisible(Client *c)
{
	if (ISDROPDOWN(c) && ISVISIBLE(c))
		placedropdown(c);
}

static void
replaceifdropdown(Client *c)
{
	storedropdownsize(c);
	if (ISDROPDOWN(c))
		placedropdown(c);
}

static void
raisedropdown(Monitor *m)
{
	Client *c;

	if (!m)
		return;
	for (c = m->clients; c; c = c->next)
		if (ISDROPDOWN(c) && ISVISIBLE(c))
			XRaiseWindow(dpy, c->win);
}

static void
dropdownsetmfact(const Arg *arg)
{
	Client *c;
	Monitor *m;
	int dropdown, maxh, delta;

	c = selmon ? selmon->sel : NULL;
	if (!ISDROPDOWN(c) || !ISVISIBLE(c)) {
		setmfact(arg);
		return;
	}

	m = c->mon;
	dropdown = c->dropdown;
	maxh = MAX(1, m->wh - 2 * c->bw);
	delta = arg->f ? (int)(m->wh * arg->f) : arg->i;
	if (!delta)
		delta = arg->f > 0 ? 1 : -1;

	m->dropw[dropdown] = c->w;
	m->droph[dropdown] = MAX(1, MIN(c->h + delta, maxh));
	placedropdown(c);
	arrange(m);
	focus(c);
	XRaiseWindow(dpy, c->win);
}

static void
toggledropdown(const Arg *arg)
{
	Client *c;
	int dropdown = arg->i;
	Arg spawnarg;

	if (dropdown < 0 || dropdown >= (int)LENGTH(dropdowns))
		return;
	if (dropdownpending[dropdown] && time(NULL) - dropdownpending[dropdown] < 2)
		return;
	if ((c = finddropdown(dropdown))) {
		if (dropdownvisibleon(c, selmon)) {
			if (selmon->sel == c)
				hidedropdown(c, 1);
			else
				showdropdown(c);
		} else
			showdropdown(c);
	} else {
		hidevisibledropdowns(dropdown);
		spawnarg.v = dropdowns[dropdown].cmd;
		dropdownpending[dropdown] = time(NULL);
		spawn(&spawnarg);
	}
}

#endif /* DWM_CONTEXT */
