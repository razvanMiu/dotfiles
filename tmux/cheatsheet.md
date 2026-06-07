# tmux shortcuts

Prefix: C-b

## Most used

C-b g        show this help
C-b ?        show raw tmux key table
C-b o        SessionX: switch/create sessions and projects
C-b p        project picker: attach/create tmux session for a repo
C-b e        scratch shell popup in current directory
C-b u        pick/open URL from pane history
C-b f        fuzzy-filter scrollback; selected line copied to clipboard
C-b /        native scrollback search, jumps in copy-mode
C-b r        reload tmux config
C-b C-s      save tmux session/layout state
C-b C-r      restore tmux session/layout state
C-b d        detach from tmux, leave session running
C-b t        unbound; tmux clock-mode disabled

## Windows = tmux tabs

C-b c        new window in current directory
C-b 1..9     jump to window number
C-b n        next window
C-b p        previous window
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

## Sessions

C-b o        fuzzy switch/create via SessionX
C-b s        built-in session picker
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
C-b f        fuzzy-filter scrollback; copy selected line to clipboard
h/j/k/l      move cursor one step
C-u/C-d      animated 16-line scroll up/down
PageUp/Down   page up/down, large jump
mouse wheel  scroll history
/            search
v            begin selection
y            copy selection to clipboard
q            quit copy mode

## SessionX popup

C-b o        open
Type name    filter sessions/projects, Enter switches or creates
Ctrl-w       window mode
Ctrl-t       tree mode
Ctrl-x       browse ~/.config
?            toggle preview
Alt-Backspace / configured kill key: delete selected session
