static time_t dropdownpending[MAXDROPDOWNS];

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
placedropdown(Client *c)
{
	Monitor *m;
	int dropdown, maxw, maxh, defaultw, defaulth;

	if (!c || !ISDROPDOWN(c) || !(m = c->mon) || c->dropdown >= (int)LENGTH(dropdowns))
		return;

	dropdown = c->dropdown;
	maxw = MAX(1, m->ww - 2 * c->bw);
	maxh = MAX(1, m->wh - 2 * c->bw);
	defaultw = MAX(1, (int)(maxw * (dropdowns[dropdown].wfact > 0 ? dropdowns[dropdown].wfact : 1.0)));
	defaulth = MAX(1, (int)(maxh * (dropdowns[dropdown].hfact > 0 ? dropdowns[dropdown].hfact : 0.5)));
	c->isfloating = 1;
	resize(c, m->wx, m->wy,
	       MIN(m->dropw[dropdown] ? m->dropw[dropdown] : defaultw, maxw),
	       MIN(m->droph[dropdown] ? m->droph[dropdown] : defaulth, maxh),
	       0);
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
	placedropdown(c);
	arrange(c->mon);
	focus(c);
	XRaiseWindow(dpy, c->win);
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
		if (dropdownvisibleon(c, selmon))
			hidedropdown(c, 1);
		else
			showdropdown(c);
	} else {
		hidevisibledropdowns(dropdown);
		spawnarg.v = dropdowns[dropdown].cmd;
		dropdownpending[dropdown] = time(NULL);
		spawn(&spawnarg);
	}
}
