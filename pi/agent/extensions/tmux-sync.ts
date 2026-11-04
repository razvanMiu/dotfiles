import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import { exec } from "child_process";

/**
 * Tmux Pi Sync Extension
 * Synchronizes Pi agent state (thinking, waiting/confirming, done)
 * to the enclosing tmux window tab via the @pi_state window option.
 */
export default function (pi: ExtensionAPI) {
  const pane = process.env.TMUX_PANE;
  if (!pane) return; // Not running inside tmux

  const setTmuxState = (state: string | null) => {
    if (state === null) {
      exec(`tmux set-window-option -u -t "${pane}" @pi_state 2>/dev/null`);
    } else {
      exec(
        `tmux set-window-option -t "${pane}" @pi_state "${state}" 2>/dev/null`,
      );
    }
  };

  // Clean up on process termination
  process.on("exit", () => setTmuxState(null));
  process.on("SIGINT", () => setTmuxState(null));
  process.on("SIGTERM", () => setTmuxState(null));

  pi.on("session_shutdown", async () => {
    setTmuxState(null);
  });

  // When agent starts processing a turn
  pi.on("turn_start", async () => {
    setTmuxState("󱚣");
  });

  // When agent finishes/settles
  pi.on("agent_settled", async () => {
    setTmuxState("✓");
  });

  // Clear "done" badge when user starts a new message
  pi.on("input", async () => {
    setTmuxState(null);
  });
}
