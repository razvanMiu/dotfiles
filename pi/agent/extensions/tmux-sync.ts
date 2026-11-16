import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import { exec, execSync } from "child_process";
import * as fs from "fs";
import * as os from "os";
import * as path from "path";

interface SubagentRun {
  id?: string;
  role: string;
  task: string;
  state: "working" | "done" | "error";
  startedAt: number;
}

interface FleetEntry {
  pid: number;
  sessionId?: string;
  name: string;
  project: string;
  cwd: string;
  sessionFile?: string;
  tmux?: {
    session: string;
    windowIndex: string;
    windowName: string;
    paneId: string;
  };
  state: "working" | "blocked" | "done" | "idle";
  stateDetail?: string;
  currentTool?: string;
  subagents: SubagentRun[];
  startedAt: number;
  updatedAt: number;
}

/**
 * Tmux Pi Sync & Fleet Controller Extension
 * 1. Synchronizes active agent state to enclosing tmux window tab (@pi_state).
 * 2. Maintains zero-overhead atomic fleet registry in ~/.local/state/pi/fleet/<pid>.json.
 * 3. Plays ambient macOS audio chimes on state transitions (blocked -> Ping, done -> Glass).
 */
export default function (pi: ExtensionAPI) {
  // If this process is a spawned subagent child, do not register as a standalone fleet session
  // or overwrite the parent's tmux window status!
  if (
    process.env.PI_SUBAGENT_PARENT_SESSION ||
    process.env.PI_SUBAGENT_RUNNER_PID ||
    process.env.PI_SUBAGENTS_WORKFLOW_RUN_ID
  ) {
    return;
  }

  const sessionStartTime = Date.now();
  let currentState: FleetEntry["state"] = "idle";
  let stateDetail: string | undefined = undefined;
  let currentTool: string | undefined = undefined;
  let activeSubagents: SubagentRun[] = [];

  const tmuxBin = fs.existsSync("/opt/homebrew/bin/tmux")
    ? "/opt/homebrew/bin/tmux"
    : "tmux";

  // Resolve tmux metadata reliably
  const resolveTmuxInfo = (): FleetEntry["tmux"] | undefined => {
    const currentPane = process.env.TMUX_PANE;
    if (currentPane) {
      try {
        const out = execSync(
          `"${tmuxBin}" display-message -p -t "${currentPane}" "#{session_name}:::#{window_index}:::#{window_name}:::#{pane_id}"`,
          { encoding: "utf8", timeout: 1000 },
        ).trim();
        const [session, windowIndex, windowName, paneId] = out.split(":::");
        if (session && windowIndex && paneId) {
          return {
            session,
            windowIndex,
            windowName: windowName || "",
            paneId,
          };
        }
      } catch {}
    }

    // Fallback: match by pane PID against parent process
    try {
      const list = execSync(
        `"${tmuxBin}" list-panes -a -F "#{session_name}:::#{window_index}:::#{window_name}:::#{pane_id}:::#{pane_pid}"`,
        { encoding: "utf8", timeout: 1000 },
      )
        .trim()
        .split("\n");
      for (const line of list) {
        const [s, w, wn, p, p_pid] = line.split(":::");
        if (p_pid === String(process.ppid) || p_pid === String(process.pid)) {
          return {
            session: s,
            windowIndex: w,
            windowName: wn || "",
            paneId: p,
          };
        }
      }
    } catch {}

    return undefined;
  };

  let cachedTmuxInfo = resolveTmuxInfo();

  // Set tmux window tab option (@pi_state)
  const setTmuxState = (state: string | null) => {
    const pane = cachedTmuxInfo?.paneId || process.env.TMUX_PANE;
    if (!pane) return;
    if (state === null) {
      exec(`"${tmuxBin}" set-window-option -u -t "${pane}" @pi_state 2>/dev/null`);
    } else {
      exec(
        `"${tmuxBin}" set-window-option -t "${pane}" @pi_state "${state}" 2>/dev/null`,
      );
    }
  };

  // Play subtle macOS ambient chimes
  const playAudio = (sound: string) => {
    if (process.env.PI_FLEET_SOUNDS === "0") return;
    const soundPath = `/System/Library/Sounds/${sound}`;
    exec(`afplay "${soundPath}" 2>/dev/null &`);
  };

  // Directory for active fleet state
  const fleetDir = path.join(
    process.env.XDG_STATE_HOME || path.join(os.homedir(), ".local", "state"),
    "pi",
    "fleet",
  );

  let lastKnownName = "";
  let lastKnownSessionFile = "";
  let lastKnownSessionId = "";
  let cwd = process.cwd();

  const writeFleetEntry = () => {
    try {
      fs.mkdirSync(fleetDir, { recursive: true, mode: 0o700 });
      const filePath = path.join(fleetDir, `${process.pid}.json`);
      const tmpPath = `${filePath}.tmp.${Date.now()}`;

      const name = pi.getSessionName() || lastKnownName || "unnamed";
      const project = name.includes("_") ? name.split("_")[0] : path.basename(cwd);

      if (!cachedTmuxInfo) {
        cachedTmuxInfo = resolveTmuxInfo();
      }

      const entry: FleetEntry = {
        pid: process.pid,
        sessionId: lastKnownSessionId || undefined,
        name,
        project: project || "default",
        cwd,
        sessionFile: lastKnownSessionFile || undefined,
        tmux: cachedTmuxInfo,
        state: currentState,
        stateDetail,
        currentTool,
        subagents: activeSubagents,
        startedAt: sessionStartTime,
        updatedAt: Date.now(),
      };

      fs.writeFileSync(tmpPath, JSON.stringify(entry, null, 2), {
        mode: 0o600,
      });
      fs.renameSync(tmpPath, filePath);
    } catch {}
  };

  const deleteFleetEntry = () => {
    try {
      const filePath = path.join(fleetDir, `${process.pid}.json`);
      if (fs.existsSync(filePath)) {
        fs.unlinkSync(filePath);
      }
    } catch {}
  };

  // Clean up on exit or signals
  const cleanup = () => {
    deleteFleetEntry();
    setTmuxState(null);
  };
  process.on("exit", cleanup);
  process.on("SIGINT", cleanup);
  process.on("SIGTERM", cleanup);

  // Hook session start to populate session metadata
  pi.on("session_start", async (_event, ctx) => {
    cwd = ctx.cwd || process.cwd();
    const sm = ctx.sessionManager as any;
    if (sm) {
      lastKnownSessionFile = sm.getSessionFile?.() || "";
      lastKnownSessionId = sm.getSessionId?.() || "";
      lastKnownName = pi.getSessionName() || sm.getSessionName?.() || "";
    }
    cachedTmuxInfo = resolveTmuxInfo();
    writeFleetEntry();
  });

  pi.on("session_shutdown", async () => {
    cleanup();
  });

  // Turn started: working
  pi.on("turn_start", async () => {
    currentState = "working";
    stateDetail = undefined;
    setTmuxState("󱚣");
    writeFleetEntry();
  });

  // Tool called: capture active tool & track subagents
  pi.on("tool_call", async (event) => {
    currentTool = event.toolName;

    if (event.toolName === "subagent") {
      const input = (event.input as any) || {};
      const agentRole = input.agent || "worker";
      const taskBrief = typeof input.task === "string" ? input.task.slice(0, 100) : "workflow task";

      // Track single child or workflow tasks
      if (input.workflow === "parallel" && Array.isArray(input.args?.tasks)) {
        for (const t of input.args.tasks) {
          activeSubagents.push({
            id: event.toolCallId,
            role: t.agent || "worker",
            task: typeof t.task === "string" ? t.task.slice(0, 100) : "parallel subtask",
            state: "working",
            startedAt: Date.now(),
          });
        }
      } else if (input.task) {
        activeSubagents.push({
          id: event.toolCallId,
          role: agentRole,
          task: taskBrief,
          state: "working",
          startedAt: Date.now(),
        });
      }
    }

    stateDetail = event.toolName === "bash"
      ? `bash: ${((event.input as any)?.command || "").slice(0, 50)}`
      : event.toolName;

    writeFleetEntry();
  });

  // Tool finished: resolve completed subagents
  pi.on("tool_result", async (event) => {
    if (event.toolName === "subagent") {
      activeSubagents = activeSubagents.filter((s) => s.id !== event.toolCallId);
    }
    currentTool = undefined;
    writeFleetEntry();
  });

  // UI prompt opened: BLOCKED (waiting on approval or user input)
  pi.on("ui_prompt_start", async (event: any) => {
    currentState = "blocked";
    stateDetail = event?.title || "Waiting for user input/approval";
    setTmuxState("⏸");
    playAudio("Ping.aiff");
    writeFleetEntry();
  });

  // UI prompt closed: resume working
  pi.on("ui_prompt_end", async () => {
    if (currentState === "blocked") {
      currentState = "working";
      stateDetail = undefined;
      setTmuxState("󱚣");
      writeFleetEntry();
    }
  });

  // Agent settled: turn done
  pi.on("agent_settled", async () => {
    currentState = "done";
    stateDetail = "Turn completed";
    currentTool = undefined;
    activeSubagents = [];
    setTmuxState("✓");
    playAudio("Glass.aiff");
    writeFleetEntry();
  });

  // User typing: idle
  pi.on("input", async () => {
    currentState = "idle";
    stateDetail = undefined;
    setTmuxState(null);
    writeFleetEntry();
  });

  // Initial entry on extension load
  writeFleetEntry();
}
