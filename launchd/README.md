# macOS login agents

Authored launchd configuration lives here. macOS discovers user login agents in
`~/Library/LaunchAgents/`, so install symlinks there rather than duplicate config.
These agents are macOS-only; do not install them on Linux.

## Built-in keyboard: swap section and grave keys

`com.razvan.keyboard-grave-section.plist` configures a swap between keyboard HID
usages `0x35` (grave/tilde) and `0x64` (the extra ISO key, used for section on this
layout). It targets only the Apple internal keyboard, not external keyboards or
trackpads. Once permitted, the swap affects every app, including modified key
combinations.

The active input source stays unchanged. With `USInternational-PC`, grave is a
dead key: press it followed by Space when testing a literal backquote.

The agent uses Apple's built-in `/usr/bin/hidutil` with no extra software or sudo.
On newer macOS, non-modifier key remapping may require **Input Monitoring**
approval for `hidutil`, even when the mapping is stored and the agent exits zero.

1. Open **System Settings → Privacy & Security → Input Monitoring**.
2. Click **+**, then **Command–Shift–G**, enter `/usr/bin/hidutil`, and add it.
3. Enable `hidutil`, then reapply the agent using `launchctl kickstart` below.

This permission must be granted manually; installing the plist does not grant it.
See the [mac-remap-keys setup notes](https://github.com/mvp/mac-remap-keys#installation)
for this newer-macOS prerequisite. Do not install that project's files alongside
this agent; both would own the same mapping.

`RunAtLoad` reapplies the mapping at each GUI login after a reboot. It does not
change the FileVault/pre-login keyboard. HID mappings are lost if the keyboard
service is removed; if that happens, rerun the agent with `launchctl kickstart`
below. No periodic overwrite or permanent background process is used.

The built-in keyboard had no `UserKeyMapping` before installation. This agent
owns that keyboard's entire mapping array; merge future remaps into the plist,
or remove this agent before handing ownership to another remapping tool.

### Install and apply

```sh
mkdir -p "$HOME/Library/LaunchAgents"
ln -s "${XDG_CONFIG_HOME:-$HOME/.config}/launchd/com.razvan.keyboard-grave-section.plist" \
  "$HOME/Library/LaunchAgents/com.razvan.keyboard-grave-section.plist"
launchctl bootstrap "gui/$(id -u)" \
  "$HOME/Library/LaunchAgents/com.razvan.keyboard-grave-section.plist"
```

### Verify / reapply

A populated `UserKeyMapping` and successful agent exit verify configuration,
**not actual typing**. Test both physical keys in TextEdit after granting
permission; use grave followed by Space with `USInternational-PC`. Also test
Ghostty's Control + grave shortcut. Physical typing and reboot behavior require
user confirmation before claiming the swap works.

```sh
launchctl print "gui/$(id -u)/com.razvan.keyboard-grave-section"
hidutil property --matching '{"VendorID":1452,"Product":"Apple Internal Keyboard / Trackpad","PrimaryUsagePage":1,"PrimaryUsage":6}' \
  --get UserKeyMapping

# Reapply if the keyboard service was reset:
launchctl kickstart "gui/$(id -u)/com.razvan.keyboard-grave-section"
```

### Undo

Unload the agent and remove its symlink so it does not return at next login.
Then clear the mapping on the same built-in keyboard only:

```sh
launchctl bootout "gui/$(id -u)/com.razvan.keyboard-grave-section"
rm "$HOME/Library/LaunchAgents/com.razvan.keyboard-grave-section.plist"
hidutil property --matching '{"VendorID":1452,"Product":"Apple Internal Keyboard / Trackpad","PrimaryUsagePage":1,"PrimaryUsage":6}' \
  --set '{"UserKeyMapping":[]}'
```

Keep the authored files in this directory if you want to enable the swap later.

## Docker: weekly storage prune

`com.razvan.docker-prune.plist` runs `~/.config/scripts/docker-prune` weekly
(Sunday 04:00; a missed run fires on next wake/login) so the Colima VM disk
cannot pile up unattended.

The script removes stopped containers, dangling images, and build cache older
than 72h. Aggressive mode (`scripts/docker-prune -a`) also removes unused
images older than 7 days. Volumes are **never** auto-pruned; the script lists
them at the end for manual review.

The agent pins `PATH`, `DOCKER_CONFIG`, and `COLIMA_HOME` as absolute paths
because launchd does not source shell config. Output is appended to
`~/.local/state/docker-prune.log`.

Install (symlink + bootstrap, same pattern as the keyboard agent):

```sh
ln -sfn ~/.config/launchd/com.razvan.docker-prune.plist ~/Library/LaunchAgents/
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.razvan.docker-prune.plist
```

Remove:

```sh
launchctl bootout gui/$(id -u)/com.razvan.docker-prune
rm ~/Library/LaunchAgents/com.razvan.docker-prune.plist
```
