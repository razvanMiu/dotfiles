import type { ExtensionAPI } from "@earendil-works/pi-coding-agent";
import { exec } from "child_process";
import * as path from "path";

export default function (pi: ExtensionAPI) {
  pi.on("agent_settled", async (event, ctx) => {
    // Skip notification if the run was manually aborted with Escape
    if (event?.aborted) return;

    const project = path.basename(ctx?.cwd || process.cwd());
    const title = "Pi Coding Agent";
    const subtitle = `Project: ${project}`;
    const message = "Task completed and ready for review.";

    // Trigger native macOS notification banner with sound
    const script = `display notification "${message}" with title "${title}" subtitle "${subtitle}" sound name "Glass"`;
    exec(`osascript -e '${script}'`);
  });
}
