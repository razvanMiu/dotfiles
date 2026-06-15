#ifndef DWM_CONTEXT
#define DWM_CONTEXT
#include "../dwm.c"
#else

static Client *lastmaster, *laststack;

static int
tiledpos(Client *target, int *count)
{
	Client *c;
	int i, pos = -1;

	for (i = 0, c = nexttiled(selmon->clients); c; c = nexttiled(c->next), i++)
		if (c == target)
			pos = i;
	*count = i;
	return pos;
}

static Client *
tiledclient(int pos)
{
	Client *c;
	int i;

	for (i = 0, c = nexttiled(selmon->clients); c; c = nexttiled(c->next), i++)
		if (i == pos)
			return c;
	return NULL;
}

static void
focusclient(Client *c)
{
	if (c) {
		focus(c);
		restack(selmon);
	}
}

static int
validslot(Client *c, int stack)
{
	int count, pos;

	if (!c || c->mon != selmon || !ISVISIBLE(c) || c->isfloating)
		return 0;
	pos = tiledpos(c, &count);
	return pos >= 0 && (pos >= selmon->nmaster) == stack;
}

static void
focustile(const Arg *arg)
{
	Client *target = NULL;
	int count, pos, first, last, stack;

	pos = tiledpos(selmon->sel, &count);
	if (pos < 0)
		return;
	stack = pos >= selmon->nmaster;

	if (arg->i == DirUp || arg->i == DirDown) {
		first = stack ? selmon->nmaster : 0;
		last = stack ? count - 1 : MIN(selmon->nmaster, count) - 1;
		if (first == last)
			return;
		pos += arg->i == DirDown ? 1 : -1;
		focusclient(tiledclient(pos < first ? last : pos > last ? first : pos));
		return;
	}

	if (!stack && arg->i == DirRight && count > selmon->nmaster) {
		lastmaster = selmon->sel;
		target = validslot(laststack, 1) ? laststack : tiledclient(selmon->nmaster);
	} else if (stack && arg->i == DirLeft) {
		laststack = selmon->sel;
		target = validslot(lastmaster, 0) ? lastmaster : tiledclient(0);
	}
	focusclient(target);
}

void
focusdir(const Arg *arg)
{
	Arg stackarg;

	if (!selmon->sel || (selmon->sel->isfullscreen && lockfullscreen))
		return;
	if (selmon->lt[selmon->sellt]->arrange == tile && !selmon->sel->isfloating) {
		focustile(arg);
	} else if (arg->i == DirUp || arg->i == DirDown) {
		stackarg.i = arg->i == DirDown ? +1 : -1;
		focusstack(&stackarg);
	}
}

#endif /* DWM_CONTEXT */
