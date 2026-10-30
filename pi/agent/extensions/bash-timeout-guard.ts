import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";

/**
 * Default timeout in seconds for shell commands when none is specified by the model.
 * Can be customized via environment variable PI_BASH_DEFAULT_TIMEOUT (e.g. 90).
 */
const DEFAULT_TIMEOUT_SECONDS = Number(process.env.PI_BASH_DEFAULT_TIMEOUT) || 60;

export default function (pi: ExtensionAPI) {
  pi.on("tool_call", async (event) => {
    if (event.toolName === "bash") {
      const input = event.input as { command?: string; timeout?: number };
      if (input && input.timeout === undefined) {
        input.timeout = DEFAULT_TIMEOUT_SECONDS;
      }
    }
  });
}
