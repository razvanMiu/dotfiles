#!/usr/bin/env bun
/**
 * check-session-integrity.ts
 *
 * Validates the tree / DAG integrity of Pi session transcript files (.jsonl).
 * Verifies that every entry links back to a valid parentId, ensuring the linear
 * context chain from leaf to root is 100% walkable by Pi's SessionManager.
 *
 * Usage:
 *   bun scripts/check-session-integrity.ts [--dir <path>] [--repair] [--quiet]
 */

import { existsSync, readdirSync, readFileSync, writeFileSync, copyFileSync } from "node:fs";
import { join } from "node:path";
import { homedir } from "node:os";

interface SessionEntry {
	type: string;
	id?: string;
	parentId?: string | null;
	[key: string]: unknown;
}

const args = process.argv.slice(2);
let repair = false;
let quiet = false;
let sessionDir = process.env.PI_CODING_AGENT_SESSION_DIR || join(homedir(), ".local/state/pi/sessions");

for (let i = 0; i < args.length; i++) {
	const arg = args[i];
	if (arg === "--repair") {
		repair = true;
	} else if (arg === "--quiet") {
		quiet = true;
	} else if (arg === "--dir" && i + 1 < args.length) {
		sessionDir = args[++i];
	} else if (arg === "-h" || arg === "--help") {
		console.log("Usage: bun scripts/check-session-integrity.ts [--dir <path>] [--repair] [--quiet]");
		process.exit(0);
	}
}

if (!existsSync(sessionDir)) {
	if (!quiet) console.error(`Session directory not found: ${sessionDir}`);
	process.exit(1);
}

const collectFiles = (dir: string): string[] => {
	if (!existsSync(dir)) return [];
	return readdirSync(dir)
		.filter((f) => f.endsWith(".jsonl"))
		.map((f) => join(dir, f))
		.sort();
};

const filePaths = [
	...collectFiles(sessionDir),
	...collectFiles(join(sessionDir, "archive")),
];

let totalFiles = 0;
let healthyFiles = 0;
let brokenFiles = 0;
const issues: string[] = [];

for (const filePath of filePaths) {
	const file = filePath.startsWith(sessionDir) ? filePath.slice(sessionDir.length + 1) : filePath;
	totalFiles++;

	const content = readFileSync(filePath, "utf-8");
	const lines = content.split("\n").filter((l) => l.trim().length > 0);
	if (lines.length === 0) continue;

	const entries: SessionEntry[] = [];
	let parseErrors = 0;

	for (const line of lines) {
		try {
			entries.push(JSON.parse(line));
		} catch {
			parseErrors++;
		}
	}

	if (parseErrors > 0) {
		issues.push(`${file}: ${parseErrors} JSON parse errors encountered`);
		brokenFiles++;
		continue;
	}

	const byId = new Map<string, SessionEntry>();
	for (const entry of entries) {
		if (entry.id) {
			byId.set(entry.id, entry);
		}
	}

	// Trace backwards from the leaf entry
	const leaf = entries[entries.length - 1];
	let current: SessionEntry | undefined = leaf;
	let walkedCount = 0;
	const visited = new Set<string>();

	while (current) {
		walkedCount++;
		if (current.id) {
			if (visited.has(current.id)) {
				// Cycle detected
				issues.push(`${file}: cyclic parentId reference detected at ${current.id}`);
				break;
			}
			visited.add(current.id);
		}
		const pid = current.parentId;
		current = pid ? byId.get(pid) : undefined;
	}

	// Check for broken links: every entry after the initial entry (index >= 2) must have a valid parentId
	const disconnected: number[] = [];
	for (let i = 2; i < entries.length; i++) {
		const entry = entries[i];
		// If an entry is session_info or normal entry, its parentId must point to an earlier entry
		if (entry.type !== "session" && (!entry.parentId || !byId.has(entry.parentId))) {
			disconnected.push(i + 1);
		}
	}

	// Walk from leaf must reach the root entry (entries[1])
	const reachedRoot = current === undefined && Boolean(entries[1]?.id && visited.has(entries[1].id));
	const isConnected = disconnected.length === 0 && reachedRoot;

	if (!isConnected) {
		brokenFiles++;
		const msg = `${file}: ${disconnected.length} disconnected entries (lines: ${disconnected.slice(0, 5).join(", ")}${disconnected.length > 5 ? "..." : ""}), walked ${walkedCount}/${entries.length}`;
		issues.push(msg);

		if (repair) {
			const backupPath = `${filePath}.bak-${Date.now()}`;
			copyFileSync(filePath, backupPath);

			let fixedCount = 0;
			// Start at index 2: entries[0] is session metadata, entries[1] is the DAG root (parentId === null)
			for (let i = 2; i < entries.length; i++) {
				const prev = entries[i - 1];
				const curr = entries[i];
				if (!curr.parentId || !byId.has(curr.parentId)) {
					if (prev.id) {
						curr.parentId = prev.id;
						fixedCount++;
					}
				}
			}

			const newContent = entries.map((e) => JSON.stringify(e)).join("\n") + "\n";
			writeFileSync(filePath, newContent, "utf-8");
			if (!quiet) console.log(`  Repaired ${file} (${fixedCount} links repaired, backup saved to ${backupPath})`);
		}
	} else {
		healthyFiles++;
	}
}

if (!quiet) {
	console.log(`Pi Session Integrity Audit:`);
	console.log(`  Directory:     ${sessionDir}`);
	console.log(`  Total files:   ${totalFiles}`);
	console.log(`  Healthy files: ${healthyFiles}`);
	console.log(`  Broken files:  ${brokenFiles}`);

	if (issues.length > 0) {
		console.log(`\nIssues detected:`);
		for (const issue of issues) {
			console.log(`  - ${issue}`);
		}
	} else {
		console.log(`  Status:        All sessions 100% DAG connected.`);
	}
}

process.exit(brokenFiles === 0 ? 0 : 1);
