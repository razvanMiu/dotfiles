import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import * as path from "path";

/** Tracks sessions where an auto-name prompt has already been shown */
const promptedSessions = new Set<string>();

export default function (pi: ExtensionAPI) {
  // Clear tracking on new session start
  pi.on("session_start", async () => {
    if (promptedSessions.size > 100) promptedSessions.clear();
  });

  async function deriveSessionNameWithAI(ctx: any): Promise<string | undefined> {
    const model = ctx.model;
    if (!model || !ctx.modelRegistry?.complete) {
      return undefined;
    }

    // Extract user requests and context from transcript entries
    let promptContext = "";
    try {
      const entries = ctx.sessionManager?.getEntries?.() || [];
      const userMessages: string[] = [];
      for (const entry of entries) {
        if (entry.type === "message" && entry.message?.role === "user") {
          const content = entry.message.content;
          if (typeof content === "string") {
            userMessages.push(content);
          } else if (Array.isArray(content)) {
            for (const item of content) {
              if (item.type === "text" && item.text) {
                userMessages.push(item.text);
              }
            }
          }
        }
      }
      promptContext = userMessages.slice(0, 3).join("\n---\n").slice(0, 800);
    } catch {
      // Ignore
    }

    if (!promptContext.trim()) return undefined;

    const systemPrompt = `You are a session naming assistant. Your task is to output a single canonical session name for an AI agent conversation.

Format:
<project>_<verb-task>[_<#ticket>]

Rules:
1. <project>: A concise lowercase slug for the project, repository, or tool (e.g. "dotfiles", "eea-frontend", "volto-eea-kitkat", "redmine", etc.). Infer it naturally from the working directory and conversation.
2. <verb-task>: 2 to 5 lowercase hyphenated words describing the primary task or action (e.g. "add-approval-mode", "fix-catalog-indexing", "upgrade-cypress-tests").
3. <#ticket>: Optional. If an issue or ticket ID is mentioned (e.g. #310194 or issue 39124), append it as "_#<number>". Otherwise omit.
4. Do NOT include any date prefix (dates are tracked separately in the session picker).
5. Output ONLY the session name string. No markdown, no quotes, no explanations, no trailing punctuation.`;

    const userPrompt = `Working directory: ${ctx.cwd || process.cwd()}

User conversation:
${promptContext}`;

    try {
      const res = await ctx.modelRegistry.complete(
        model,
        {
          systemPrompt,
          messages: [
            {
              role: "user",
              content: userPrompt,
              timestamp: Date.now(),
            },
          ],
        },
        { maxTokens: 40 },
      );

      const raw = res?.content?.[0]?.type === "text" ? res.content[0].text : "";
      const cleaned = raw.trim().replace(/^["'`]|["'`]$/g, "").trim();

      // Validate format: <project>_<verb-task>[_<#ticket>]
      if (/^[a-z0-9_-]+_[a-z0-9_-]+$/i.test(cleaned)) {
        return cleaned;
      }
      return undefined;
    } catch {
      return undefined;
    }
  }

  async function promptSessionNaming(ctx: any, isManual = false) {
    if (!ctx.hasUI) return;

    // Call the model for the suggestion
    const candidate = await deriveSessionNameWithAI(ctx);
    if (!candidate) {
      if (isManual) {
        ctx.ui.notify("Could not generate session name suggestion.", "warning");
      }
      return;
    }

    const choice = await ctx.ui.select(`Apply session name?\n${candidate}`, [
      `✓ Accept: ${candidate}`,
      "✎ Edit name...",
      "✗ Skip / Keep unnamed",
    ]);

    if (!choice || choice.startsWith("✗ Skip")) {
      return;
    }

    if (choice.startsWith("✓ Accept")) {
      pi.setSessionName(candidate);
      ctx.ui.notify(`Session named: ${candidate}`, "info");
      return;
    }

    if (choice.startsWith("✎ Edit name")) {
      const edited = await ctx.ui.editor("Edit session name:", candidate);
      if (edited && edited.trim()) {
        const finalName = edited.replace(/[\r\n]+/g, " ").trim();
        pi.setSessionName(finalName);
        ctx.ui.notify(`Session named: ${finalName}`, "info");
      }
    }
  }

  // Hook agent_settled to prompt on first turn completion
  pi.on("agent_settled", async (event, ctx) => {
    // Skip if manually aborted
    if (event?.aborted) return;

    // Check if session is already named
    const currentName = pi.getSessionName();
    if (currentName) return;

    // Skip if ephemeral session (no disk file)
    const sessionFile = ctx.sessionManager?.getSessionFile?.();
    if (!sessionFile) return;

    // Only prompt once per session to avoid nagging
    const sessionId = path.basename(sessionFile);
    if (promptedSessions.has(sessionId)) return;
    promptedSessions.add(sessionId);

    // Ask for confirmation with AI-derived suggestion
    await promptSessionNaming(ctx);
  });

  // Register manual command /autoname
  pi.registerCommand("autoname", {
    description:
      "Ask the active AI model to suggest and apply a canonical session name",
    handler: async (_args, ctx) => {
      await promptSessionNaming(ctx, true);
    },
  });
}
