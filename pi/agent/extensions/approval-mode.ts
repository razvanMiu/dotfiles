import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import * as path from "path";

export type ExecutionMode = "full" | "supervised" | "plan";

let currentMode: ExecutionMode = "full";

/** In-memory auto-approved rules for the active session */
const sessionRules = new Set<string>();

/** Pending user instructions keyed by toolCallId to be injected into tool_result */
const pendingInstructions = new Map<string, string>();

/** Helper to check if a tool call matches any session auto-approval rule */
function matchesRule(toolName: string, input: any): boolean {
  if (sessionRules.has(`tool:${toolName}`)) return true;

  if (toolName === "bash") {
    const cmd = ((input?.command as string) || "").trim();
    if (sessionRules.has(`bash:exact:${cmd}`)) return true;
    for (const rule of sessionRules) {
      if (rule.startsWith("bash:prefix:")) {
        const prefix = rule.slice("bash:prefix:".length);
        if (cmd === prefix || cmd.startsWith(prefix + " ")) return true;
      }
    }
    return false;
  }

  if (toolName === "write" || toolName === "edit") {
    const filePath = (input?.path as string) || "";
    if (sessionRules.has(`file:exact:${filePath}`)) return true;
    for (const rule of sessionRules) {
      if (rule.startsWith("file:dir:")) {
        const dir = rule.slice("file:dir:".length);
        if (filePath.startsWith(dir)) return true;
      }
      if (rule.startsWith("file:ext:")) {
        const ext = rule.slice("file:ext:".length);
        if (filePath.endsWith(ext)) return true;
      }
    }
    return false;
  }

  return false;
}

/** Helper to extract binary and subcommand for bash matching */
function extractBaseCommand(cmd: string): { base: string; sub?: string } {
  const clean = cmd
    .replace(/^(sudo|env\s+[A-Za-z0-9_]+=.*?|\w+=.*?)\s+/, "")
    .trim();
  const parts = clean.split(/\s+/);
  const base = parts[0] || "";
  let sub: string | undefined = undefined;
  if (
    parts.length > 1 &&
    ["git", "npm", "pnpm", "yarn", "cargo", "docker", "bun", "make", "python", "uv"].includes(base)
  ) {
    sub = `${base} ${parts[1]}`;
  }
  return { base, sub };
}

export default function (pi: ExtensionAPI) {
  // Clear rules when a new session starts
  pi.on("session_start", async () => {
    sessionRules.clear();
    pendingInstructions.clear();
  });

  // 1. Register /mode command
  pi.registerCommand("mode", {
    description:
      "Switch execution mode (full / supervised / plan) or manage session auto-approval rules",
    getArgumentCompletions: (prefix) => {
      const modes = ["full", "supervised", "plan", "rules", "clear-rules"];
      return modes
        .filter((m) => m.startsWith(prefix))
        .map((m) => ({ value: m, label: m }));
    },
    handler: async (args, ctx) => {
      const target = args.trim().toLowerCase();

      if (target === "clear-rules") {
        sessionRules.clear();
        ctx.ui.notify("Session auto-approval rules cleared.", "info");
        return;
      }

      if (target === "rules") {
        if (sessionRules.size === 0) {
          ctx.ui.notify("No auto-approval rules active for this session.", "info");
        } else {
          const list = Array.from(sessionRules).map((r) => `• ${r}`);
          await ctx.ui.select(`Active session rules (${sessionRules.size}):`, [
            ...list,
            "Clear all rules",
            "Back",
          ]).then((choice) => {
            if (choice === "Clear all rules") {
              sessionRules.clear();
              ctx.ui.notify("Session auto-approval rules cleared.", "info");
            }
          });
        }
        return;
      }

      let newMode: ExecutionMode;
      if (target === "full" || target === "supervised" || target === "plan") {
        newMode = target as ExecutionMode;
      } else {
        const choice = await ctx.ui.select("Select execution mode:", [
          "full - auto-execute all tools (default)",
          "supervised - Claude Code-style interactive approvals & steering",
          "plan - strictly read-only exploration",
          `rules - view / clear auto-approval rules (${sessionRules.size} active)`,
        ]);
        if (!choice) return;
        if (choice.startsWith("full")) newMode = "full";
        else if (choice.startsWith("supervised")) newMode = "supervised";
        else if (choice.startsWith("plan")) newMode = "plan";
        else {
          sessionRules.clear();
          ctx.ui.notify("Session auto-approval rules cleared.", "info");
          return;
        }
      }

      currentMode = newMode;
      updateStatus(ctx);
      ctx.ui.notify(`Execution mode: ${newMode}`, "info");
    },
  });

  // 2. Register /plan shortcut alias to toggle plan mode
  pi.registerCommand("plan", {
    description: "Toggle read-only plan mode",
    handler: async (_args, ctx) => {
      currentMode = currentMode === "plan" ? "full" : "plan";
      updateStatus(ctx);
      ctx.ui.notify(`Execution mode: ${currentMode}`, "info");
    },
  });

  // 3. Register shortcut ctrl+alt+m to cycle mode (full -> supervised -> plan)
  pi.registerShortcut("ctrl+alt+m", {
    description: "Cycle execution mode (full / supervised / plan)",
    handler: async (ctx) => {
      if (currentMode === "full") currentMode = "supervised";
      else if (currentMode === "supervised") currentMode = "plan";
      else currentMode = "full";
      updateStatus(ctx);
      ctx.ui.notify(`Execution mode: ${currentMode}`, "info");
    },
  });

  function updateStatus(ctx: any) {
    if (currentMode === "full") {
      ctx.ui.setStatus("mode", undefined);
    } else if (currentMode === "supervised") {
      const rulesBadge = sessionRules.size > 0 ? ` (${sessionRules.size} rules)` : "";
      ctx.ui.setStatus("mode", `🔒 supervised${rulesBadge}`);
    } else if (currentMode === "plan") {
      ctx.ui.setStatus("mode", "🛡️ plan (read-only)");
    }
  }

  // 4. Intercept tool_call with Claude Code approval UX
  pi.on("tool_call", async (event, ctx) => {
    if (currentMode === "full") return undefined;

    const toolName = event.toolName;
    const input = event.input as any;

    // Plan Mode: read-only gate
    if (currentMode === "plan") {
      if (toolName === "write" || toolName === "edit") {
        return {
          block: true,
          reason: `Plan mode is active (read-only). File modification to "${input?.path}" blocked. Switch mode with /mode full or /mode supervised to permit writes.`,
        };
      }
      if (toolName === "bash") {
        const cmd = (input?.command || "").trim();
        const modifying =
          /\b(rm|mv|cp|mkdir|touch|git\s+(commit|push|merge|rebase|reset|checkout\s+-b)|npm\s+(install|i|uninstall)|yarn\s+(add|remove)|pip\s+install|chmod|chown)\b/i;
        if (modifying.test(cmd)) {
          return {
            block: true,
            reason: `Plan mode is active (read-only). Modifying shell command blocked: "${cmd}". Switch mode with /mode full or /mode supervised to permit shell writes.`,
          };
        }
      }
      return undefined;
    }

    // Supervised Mode: Claude Code-style interactive approvals
    if (currentMode === "supervised") {
      if (!ctx.hasUI) {
        return {
          block: true,
          reason: "Supervised mode requires interactive confirmation UI.",
        };
      }

      // Safe read-only inspection tools: always allow without prompt
      if (
        toolName === "read" ||
        toolName === "grep" ||
        toolName === "find" ||
        toolName === "ls" ||
        toolName === "recall"
      ) {
        return undefined;
      }

      // Safe bash commands: always allow
      if (toolName === "bash") {
        const cmd = (input?.command || "").trim();
        const safeBash =
          /^(pwd|ls(\s+-[a-zA-Z0-9]+)*|git\s+(status|diff|log|branch)(\s+.*)?|which\s+.*|whoami|uname(\s+.*)?)$/;
        if (safeBash.test(cmd)) {
          return undefined;
        }
      }

      // Check if already covered by an auto-approval rule
      if (matchesRule(toolName, input)) {
        return undefined;
      }

      // Construct dynamic approval menu tailored to the tool call
      let title = `Approve ${toolName}?`;
      let options: string[] = [];

      if (toolName === "bash") {
        const cmd = ((input?.command as string) || "").trim();
        const preview = cmd.length > 100 ? cmd.slice(0, 97) + "..." : cmd;
        title = `Approve Bash Command?\n$ ${preview}`;

        const { base, sub } = extractBaseCommand(cmd);
        options = [
          "✓ Approve once",
          ...(sub ? [`✓ Approve & always allow '${sub}' this session`] : []),
          ...(base && base !== sub ? [`✓ Approve & always allow '${base} *' this session`] : []),
          "✓ Approve & always allow this exact command this session",
          "💬 Approve with instruction...",
          "✗ Reject",
          "💬 Reject with instruction...",
        ];
      } else if (toolName === "write" || toolName === "edit") {
        const filePath = (input?.path as string) || "";
        const ext = path.extname(filePath);
        const dir = path.dirname(filePath);
        const actionLabel = toolName === "write" ? "Write file" : "Edit file";
        title = `Approve ${actionLabel}?\nPath: ${filePath}`;

        options = [
          "✓ Approve once",
          ...(dir && dir !== "." && dir !== "/"
            ? [`✓ Approve & always allow changes in '${dir}/**' this session`]
            : []),
          ...(ext ? [`✓ Approve & always allow changes to '*${ext}' this session`] : []),
          "✓ Approve & always allow this exact file this session",
          "💬 Approve with instruction...",
          "✗ Reject",
          "💬 Reject with instruction...",
        ];
      } else {
        title = `Approve Tool Execution?\nTool: ${toolName}`;
        options = [
          "✓ Approve once",
          `✓ Approve & always allow tool '${toolName}' this session`,
          "💬 Approve with instruction...",
          "✗ Reject",
          "💬 Reject with instruction...",
        ];
      }

      // Prompt user
      const choice = await ctx.ui.select(title, options);

      // Handle Rejection
      if (!choice || choice === "✗ Reject") {
        return { block: true, reason: "Tool execution rejected by user." };
      }

      if (choice === "💬 Reject with instruction...") {
        const instruction = await ctx.ui.input(
          "Feedback for agent (why rejected / what to do instead):",
          "",
        );
        const reason = instruction?.trim()
          ? `Action rejected by user. Feedback / instructions: ${instruction.trim()}`
          : "Tool execution rejected by user.";
        return { block: true, reason };
      }

      // Handle Approval with Instruction
      if (choice === "💬 Approve with instruction...") {
        const instruction = await ctx.ui.input(
          "Guidance / instruction to execute alongside this step:",
          "",
        );
        if (instruction?.trim()) {
          pendingInstructions.set(event.toolCallId, instruction.trim());
          ctx.ui.notify("Instruction queued for agent.", "info");
        }
        return undefined;
      }

      // Handle Auto-Approval Rules
      if (choice.includes("always allow")) {
        if (toolName === "bash") {
          const cmd = ((input?.command as string) || "").trim();
          const { base, sub } = extractBaseCommand(cmd);

          if (sub && choice.includes(`'${sub}'`)) {
            sessionRules.add(`bash:prefix:${sub}`);
            ctx.ui.notify(`Rule added: always allow '${sub}'`, "info");
          } else if (base && choice.includes(`'${base} *'`)) {
            sessionRules.add(`bash:prefix:${base}`);
            ctx.ui.notify(`Rule added: always allow '${base} *'`, "info");
          } else {
            sessionRules.add(`bash:exact:${cmd}`);
            ctx.ui.notify(`Rule added: always allow exact command`, "info");
          }
        } else if (toolName === "write" || toolName === "edit") {
          const filePath = (input?.path as string) || "";
          const ext = path.extname(filePath);
          const dir = path.dirname(filePath);

          if (dir && dir !== "." && choice.includes(`'${dir}/**'`)) {
            sessionRules.add(`file:dir:${dir}`);
            ctx.ui.notify(`Rule added: always allow '${dir}/**'`, "info");
          } else if (ext && choice.includes(`'*${ext}'`)) {
            sessionRules.add(`file:ext:${ext}`);
            ctx.ui.notify(`Rule added: always allow '*${ext}'`, "info");
          } else {
            sessionRules.add(`file:exact:${filePath}`);
            ctx.ui.notify(`Rule added: always allow '${filePath}'`, "info");
          }
        } else {
          sessionRules.add(`tool:${toolName}`);
          ctx.ui.notify(`Rule added: always allow tool '${toolName}'`, "info");
        }
        updateStatus(ctx);
        return undefined;
      }

      // Default: "✓ Approve once"
      return undefined;
    }

    return undefined;
  });

  // 5. Inject queued user instructions into tool_result so the agent sees them
  pi.on("tool_result", async (event) => {
    const instruction = pendingInstructions.get(event.toolCallId);
    if (!instruction) return undefined;
    pendingInstructions.delete(event.toolCallId);

    const note = `\n\n[User Steering & Guidance]: ${instruction}`;
    const content = [...event.content];
    const last = content[content.length - 1];

    if (last && last.type === "text") {
      content[content.length - 1] = {
        type: "text",
        text: last.text + note,
      };
    } else {
      content.push({ type: "text", text: note });
    }

    return { content };
  });
}
