# tmux shortcuts

Prefix: C-b

## Most used

C-b g        show this help (cheatsheet)
C-b G        Lazygit popup in current repo
C-b a        floating Pi scratchpad (persists in-memory; C-b d hides, /exit quits)
C-b e        floating scratch shell (persists in-memory; C-b d hides, exit quits)
C-b p        Pi Sessions popup (resume, new window/session, create, delete)
C-b s        Sessions & Projects popup (switch session, open in pane/window/session)
C-b f        Pi Fleet Controller (live agents & subagents, jump, live tail, interrupt)
C-b Space    quick toggle between last two active windows
C-b BSpace   quick toggle between last two active sessions
C-b u        pick/open URL from pane history (macOS browser)
C-b /        native scrollback search, jumps in copy-mode
C-b r        reload tmux config
C-b C-s      save tmux session/layout state
C-b C-r      restore tmux session/layout state
C-b d        detach from tmux, leave session running
C-b t        unbound; tmux clock-mode disabled

## Windows = tmux tabs

Cmd-1..9     jump directly to window (via Ghostty, no prefix needed)
Alt-1..9     jump directly to window (Option/Alt, no prefix needed)
C-b 1..9     jump to window number
C-b c        new window in current directory
C-b n        next window
C-b Space    toggle last active window
C-b < / >    swap window position left / right
C-b ,        rename window
C-b &        close window

## Panes = splits inside a window

C-b |        split horizontally, left/right
C-b -        split vertically, top/bottom
C-b h/j/k/l  move between panes
C-b H/J/K/L  resize pane
C-b T        set current pane title
C-b x        close pane
C-b z        zoom/unzoom current pane

## Sessions & Projects

C-b s        fuzzy switch active sessions or open projects
C-b p        browse/resume Pi agent sessions (Ctrl-A to archive)
C-b P        browse/unarchive Pi sessions archive (Ctrl-U to unarchive)
C-b f        Pi Fleet Controller (real-time agent status, jump, live tail)
C-b $        rename session
C-b d        detach

From shell:

tmux new -s name                 create session
tmux attach -t name              attach session
tmux new-session -A -s main      attach or create main
tmux ls                          list sessions
tmux kill-session -t name        kill session

## Copy mode

C-b [        enter copy mode; freezes this client's live pane view
C-b /        search scrollback and jump to match
h/j/k/l      move cursor one step
C-u/C-d      animated 16-line scroll up/down
PageUp/Down   page up/down, large jump
mouse wheel  scroll history
/            search
v            begin selection
y            copy selection to clipboard
q            quit copy mode

## Popup Controls Standard (Pi Sessions & Projects)

Enter        Open in current pane (gracefully terminates running command first)
Ctrl-w       Open in new Window (matches prefix c)
Ctrl-s       Open in new Session (matches prefix s)
Esc / Ctrl-c Cancel / dismiss popup

## Pi Fleet Controller Controls (prefix f)

Enter        Jump to agent or subagent pane
Ctrl-t       Zoom preview pane (50% -> 75% -> 99% full-screen with zero flicker)
Ctrl-w       Open live tail in a side split pane (42% width)
Ctrl-k       Interrupt running agent or targeted subagent
Esc / Ctrl-c Close popup
