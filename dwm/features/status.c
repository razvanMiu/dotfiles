#ifndef DWM_CONTEXT
#define DWM_CONTEXT
#include "../dwm.c"
#else
int
statussegment(int *offset, char *text, size_t size)
{
	int i = 0;

	if (!stext[*offset])
		return 0;
	while (stext[*offset] && stext[*offset] != STATUSSEP) {
		if (i < size - 1)
			text[i++] = stext[*offset];
		(*offset)++;
	}
	text[i] = '\0';
	if (stext[*offset] == STATUSSEP)
		(*offset)++;
	return 1;
}

int
statuswidth(void)
{
	char text[sizeof(stext)];
	int offset = 0, w = 0;

	while (statussegment(&offset, text, sizeof(text)))
		w += TEXTW(text);
	return w ? w + 2 : 0;
}

void
drawstatus(Monitor *m, int tw)
{
	char text[sizeof(stext)];
	int offset = 0, x, w;

	if (!tw)
		return;
	x = m->ww - tw;
	while (statussegment(&offset, text, sizeof(text))) {
		w = TEXTW(text);
		drw_setscheme(drw, scheme[SchemeNorm]);
		drw_text(drw, x, 0, w, bh, lrpad / 2, text, 0);
		x += w;
	}
}

void
statusclick(XButtonPressedEvent *ev)
{
	char text[sizeof(stext)], idx[16], button[16];
	int i = 0, matched = 0, offset = 0, rx, start, tw, w;

	tw = statuswidth();
	start = selmon->ww - tw;
	rx = ev->x - start;
	if (!tw || rx < 0 || rx > tw)
		return;

	while (statussegment(&offset, text, sizeof(text))) {
		w = TEXTW(text);
		if (rx <= w) {
			matched = 1;
			break;
		}
		rx -= w;
		i++;
	}
	if (!matched)
		return;

	if (i < LENGTH(statussegments) && statussegments[i].func) {
		statussegments[i].func(&statussegments[i].arg);
		return;
	}

	snprintf(idx, sizeof(idx), "%d", i);
	snprintf(button, sizeof(button), "%u", ev->button);
	if (fork() == 0) {
		if (dpy)
			close(ConnectionNumber(dpy));
		setsid();
		setenv("STATUS_INDEX", idx, 1);
		setenv("STATUS_BUTTON", button, 1);
		setenv("STATUS_SEGMENT", i < LENGTH(statussegments) ? statussegments[i].name : "", 1);
		setenv("STATUS_TEXT", stext, 1);
		execvp(statusclickcmd[0], (char *const *)statusclickcmd);
		exit(EXIT_SUCCESS);
	}
}

#endif /* DWM_CONTEXT */
