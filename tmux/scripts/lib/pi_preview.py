"""Shared geometry and ANSI palette for the fzf pickers in tmux/scripts.

The pickers import this module from their inline ``python3 -c`` preview and
list blocks. fzf runs preview commands with an arbitrary working directory, so
the caller must put this file's directory on sys.path explicitly instead of
relying on cwd; every picker appends its own ``lib`` directory as the final
argument of the python invocation::

    python3 -c "<code>" "$item" "$tmux_w" "$SCRIPT_DIR/lib"

    import sys
    sys.dont_write_bytecode = True   # keep __pycache__ out of the tracked tree
    sys.path.insert(0, sys.argv[-1])
    from pi_preview import C_CYAN, C_RESET, parse_tmux_width, preview_div_len

Module scope holds constants only: nothing here reads files, inspects argv, or
keeps state between calls, so callers never have to reset anything.
"""

DEFAULT_TMUX_WIDTH = 160

# A popup spans 90% of the tmux client; the preview pane takes 50% of the
# popup's inner width (client width minus the two border columns).
POPUP_RATIO = 0.90
PREVIEW_COLUMN_RATIO = 0.50
PREVIEW_DIVIDER_MIN = 35
PREVIEW_DIVIDER_MARGIN = 10

C_RESET = '\033[0m'
C_BOLD = '\033[1m'
C_DIM = '\033[2m'
C_RED = '\033[1;31m'
C_GREEN = '\033[1;32m'
C_ORANGE = '\033[1;33m'
C_BLUE = '\033[1;34m'
C_MAGENTA = '\033[1;35m'
C_CYAN = '\033[1;36m'


def parse_tmux_width(raw):
    """Coerce a tmux client width argument to int, defaulting to 160."""
    try:
        return int(raw)
    except (TypeError, ValueError):
        return DEFAULT_TMUX_WIDTH


def popup_widths(tmux_w):
    """Return ``(popup_w, inner_w)`` for a popup on a client of width ``tmux_w``."""
    popup_w = int(tmux_w * POPUP_RATIO)
    return popup_w, popup_w - 2


def preview_widths(tmux_w):
    """Return ``(popup_w, inner_w, preview_w)`` for a 50% right-hand preview pane."""
    popup_w, inner_w = popup_widths(tmux_w)
    return popup_w, inner_w, int(inner_w * PREVIEW_COLUMN_RATIO)


def preview_div_len(tmux_w):
    """Return the divider length for the preview body: ``max(35, preview_w - 10)``."""
    _popup_w, _inner_w, preview_w = preview_widths(tmux_w)
    return max(PREVIEW_DIVIDER_MIN, preview_w - PREVIEW_DIVIDER_MARGIN)
