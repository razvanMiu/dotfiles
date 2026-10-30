#!/usr/bin/env bun

import { existsSync, readFileSync } from "node:fs";

type InstanceConfig = { url: string; apiKeyEnv: string };
type Config = { defaultInstance?: string; instances: Record<string, InstanceConfig> };
type Args = string[];
type RedmineCurrentUserResponse = { user?: { id?: number } };

const CONFIG_PATH = process.env.REDMINE_CONFIG ?? `${process.env.XDG_CONFIG_HOME ?? `${process.env.HOME}/.config`}/redmine/config.json`;

main().catch((error) => {
  console.error(`[redmine] ${error instanceof Error ? error.message : String(error)}`);
  process.exit(1);
});

async function main() {
  const args = process.argv.slice(2);
  if (args.length === 0 || args[0] === "help" || args[0] === "--help" || args[0] === "-h") {
    usage();
    return;
  }

  const command = args.shift();
  if (command === "instances") {
    const config = loadConfig();
    for (const [name, instance] of Object.entries(config.instances)) {
      const marker = config.defaultInstance === name ? "*" : " ";
      console.log(`${marker} ${name}\t${instance.url}\tapiKeyEnv=${instance.apiKeyEnv}`);
    }
    return;
  }

  if (command === "trackers") {
    const instanceName = args.shift();
    const [, instance] = resolveInstance(instanceName);
    const data = (await redmineFetch(instance, "GET", "/trackers.json")) as { trackers?: Array<{ id: number; name: string }> };
    console.log(JSON.stringify(data.trackers ?? [], null, 2));
    return;
  }

  if (command === "statuses") {
    const instanceName = args.shift();
    const [, instance] = resolveInstance(instanceName);
    const data = (await redmineFetch(instance, "GET", "/issue_statuses.json")) as {
      issue_statuses?: Array<{ id: number; name: string; is_closed?: boolean }>;
    };
    console.log(JSON.stringify(data.issue_statuses ?? [], null, 2));
    return;
  }

  if (command === "priorities") {
    const instanceName = args.shift();
    const [, instance] = resolveInstance(instanceName);
    const data = (await redmineFetch(instance, "GET", "/enumerations/issue_priorities.json")) as {
      issue_priorities?: Array<{ id: number; name: string; is_default?: boolean }>;
    };
    console.log(JSON.stringify(data.issue_priorities ?? [], null, 2));
    return;
  }

  if (command === "projects") {
    const instanceName = args.shift();
    const [, instance] = resolveInstance(instanceName);
    const options = parseOptions(args);
    const params = { limit: maybeNumber(options.limit) ?? 50 };
    const data = (await redmineFetch(instance, "GET", appendQuery("/projects.json", params))) as {
      projects?: Array<{ id: number; name: string; identifier: string }>;
    };
    console.log(JSON.stringify(data.projects ?? [], null, 2));
    return;
  }

  if (command === "issues" && args[0] === "list") {
    args.shift();
    await listIssues(args);
    return;
  }

  if (command === "issue") {
    const sub = args.shift();
    if (sub === "get") return getIssue(args);
    if (sub === "create") return createIssue(args);
    if (sub === "update") return updateIssue(args);
    if (sub === "comment") return commentIssue(args);
  }

  if (command === "request") {
    await rawRequest(args);
    return;
  }

  throw new Error(`unknown command: ${command ?? ""}`);
}

function usage() {
  console.log(`Usage: redmine.ts <command> [args]

Commands:
  instances
  trackers <instance>
  statuses <instance>
  priorities <instance>
  projects <instance> [--limit N]
  issues list <instance> [--project ID] [--directly-assigned-to-me] [--assigned-to me|ID] [--status open|closed|*|ID] [--limit N] [--query key=value]
  issue get <instance> <id> [--include journals,attachments,relations,children]
  issue create <instance> --project ID --subject TEXT [--description TEXT] [--tracker-id ID] [--priority-id ID] [--assigned-to-id ID] [--parent-id ID] [--custom-field ID=VALUE] [--execute]
  issue update <instance> <id> [--status-id ID] [--assigned-to-id ID] [--priority-id ID] [--tracker-id ID] [--parent-id ID] [--subject TEXT] [--description TEXT] [--notes TEXT] [--custom-field ID=VALUE] [--execute]
  issue comment <instance> <id> --notes TEXT [--execute]
  request <instance> <METHOD> <path> [--data-json JSON] [--execute]

Mutating issue commands are dry-run by default. Add --execute to send POST/PUT/DELETE.
Always ask for human permission before applying any write or mutation operation.
Set REDMINE_CONFIG to override ${CONFIG_PATH}.
Use with-secrets -- redmine.ts ... so API key env vars are populated safely.`);
}

function loadConfig(): Config {
  if (!existsSync(CONFIG_PATH)) {
    throw new Error(`config not found: ${CONFIG_PATH}; copy ~/.config/redmine/config.example.json to config.json and edit URLs/env names`);
  }
  const config = JSON.parse(readFileSync(CONFIG_PATH, "utf8")) as Config;
  if (!config.instances || typeof config.instances !== "object") throw new Error("config.instances missing");
  return config;
}

function resolveInstance(name: string | undefined): [string, InstanceConfig] {
  const config = loadConfig();
  const instanceName = name || config.defaultInstance;
  if (!instanceName) throw new Error("instance required and no defaultInstance configured");
  const instance = config.instances[instanceName];
  if (!instance) throw new Error(`unknown instance: ${instanceName}`);
  if (!instance.url || !instance.apiKeyEnv) throw new Error(`instance ${instanceName} must define url and apiKeyEnv`);
  return [instanceName, { url: instance.url.replace(/\/+$/, ""), apiKeyEnv: instance.apiKeyEnv }];
}

function apiKey(instance: InstanceConfig): string {
  const key = process.env[instance.apiKeyEnv];
  if (!key) throw new Error(`${instance.apiKeyEnv} missing; run through with-secrets or export it locally`);
  return key;
}

function buildInstanceUrl(instance: InstanceConfig, path: string): string {
  if (/^[a-zA-Z][a-zA-Z\d+\-.]*:/.test(path) || path.startsWith("//")) {
    throw new Error(`path must be relative to instance URL, got: ${path}`);
  }
  const base = new URL(`${instance.url.replace(/\/+$/, "")}/`);
  const resolved = new URL(path.replace(/^\/+/, ""), base);
  if (resolved.origin !== base.origin || !resolved.pathname.startsWith(base.pathname)) {
    throw new Error(`path escapes instance base URL: ${path}`);
  }
  return resolved.toString();
}

async function redmineFetch(instance: InstanceConfig, method: string, path: string, body?: unknown) {
  const url = buildInstanceUrl(instance, path);
  const response = await fetch(url, {
    method,
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
      "X-Redmine-API-Key": apiKey(instance),
    },
    body: body === undefined ? undefined : JSON.stringify(body),
    redirect: "manual",
  });
  const text = await response.text();
  if (!response.ok) {
    throw new Error(`${method} ${url} failed: HTTP ${response.status} ${response.statusText}${text ? `\n${text}` : ""}`);
  }
  if (!text.trim()) return undefined;
  try {
    return JSON.parse(text);
  } catch {
    return text;
  }
}

function parseOptions(args: Args): Record<string, string | boolean | string[]> {
  const options: Record<string, string | boolean | string[]> = {};
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (arg === undefined || !arg.startsWith("--")) throw new Error(`unexpected argument: ${arg ?? ""}`);
    const key = arg.slice(2);
    if (key === "execute" || key === "dry-run" || key === "directly-assigned-to-me") {
      options[key] = true;
      continue;
    }
    const value = args[++i];
    if (value === undefined || value.startsWith("--")) throw new Error(`missing value for --${key}`);
    if (key === "query" || key === "custom-field") {
      const list = Array.isArray(options[key]) ? (options[key] as string[]) : [];
      list.push(value);
      options[key] = list;
    } else {
      options[key] = value;
    }
  }
  return options;
}

function parseCustomFields(rawList: unknown): Array<{ id: number; value: string }> | undefined {
  if (!rawList) return undefined;
  const list = Array.isArray(rawList) ? rawList : [rawList];
  return list.map((item) => {
    if (typeof item !== "string") throw new Error(`invalid custom field: ${String(item)}`);
    const eq = item.indexOf("=");
    if (eq <= 0) throw new Error(`--custom-field must be ID=VALUE, got: ${item}`);
    const id = maybeNumber(item.slice(0, eq));
    if (id === undefined) throw new Error(`--custom-field ID must be an integer, got: ${item}`);
    return { id, value: item.slice(eq + 1) };
  });
}

function needString(options: Record<string, unknown>, key: string): string {
  const value = options[key];
  if (typeof value !== "string" || value.length === 0) throw new Error(`--${key} is required`);
  return value;
}

function maybeNumber(value: unknown): number | undefined {
  if (value === undefined) return undefined;
  if (typeof value !== "string" || !/^-?\d+$/.test(value.trim())) {
    throw new Error(`expected integer, got ${String(value)}`);
  }
  return Number(value);
}

function appendQuery(path: string, params: Record<string, string | number | undefined>) {
  const query = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    if (value !== undefined) query.set(key, String(value));
  }
  const qs = query.toString();
  return qs ? `${path}?${qs}` : path;
}

async function listIssues(args: Args) {
  const instanceName = args.shift();
  const [, instance] = resolveInstance(instanceName);
  const options = parseOptions(args);
  if (options["directly-assigned-to-me"] && options["assigned-to"] !== undefined) {
    throw new Error("use either --directly-assigned-to-me or --assigned-to, not both");
  }
  const assignedTo = options["directly-assigned-to-me"]
    ? await currentUserId(instance)
    : (options["assigned-to"] as string | undefined);
  const params: Record<string, string | number | undefined> = {
    project_id: options.project as string | undefined,
    assigned_to_id: assignedTo,
    status_id: options.status as string | undefined,
    limit: maybeNumber(options.limit) ?? 25,
  };
  for (const item of (options.query as string[] | undefined) ?? []) {
    const [key, ...rest] = item.split("=");
    if (!key || rest.length === 0) throw new Error(`--query must be key=value, got ${item}`);
    params[key] = rest.join("=");
  }
  console.log(JSON.stringify(await redmineFetch(instance, "GET", appendQuery("/issues.json", params)), null, 2));
}

async function currentUserId(instance: InstanceConfig): Promise<number> {
  const data = (await redmineFetch(instance, "GET", "/users/current.json")) as RedmineCurrentUserResponse;
  const id = data.user?.id;
  if (typeof id !== "number" || !Number.isInteger(id)) throw new Error("/users/current.json did not return user.id");
  return id;
}

async function getIssue(args: Args) {
  const instanceName = args.shift();
  const id = args.shift();
  if (!id) throw new Error("issue id required");
  const [, instance] = resolveInstance(instanceName);
  const options = parseOptions(args);
  const path = appendQuery(`/issues/${encodeURIComponent(id)}.json`, { include: options.include as string | undefined });
  console.log(JSON.stringify(await redmineFetch(instance, "GET", path), null, 2));
}

async function createIssue(args: Args) {
  const instanceName = args.shift();
  const [, instance] = resolveInstance(instanceName);
  const options = parseOptions(args);
  const issue: Record<string, unknown> = {
    project_id: needString(options, "project"),
    subject: needString(options, "subject"),
  };
  setIf(issue, "description", options.description);
  setIf(issue, "tracker_id", maybeNumber(options["tracker-id"]));
  setIf(issue, "priority_id", maybeNumber(options["priority-id"]));
  setIf(issue, "assigned_to_id", maybeNumber(options["assigned-to-id"]));
  setIf(issue, "parent_issue_id", maybeNumber(options["parent-id"]));
  const customFields = parseCustomFields(options["custom-field"]);
  if (customFields && customFields.length > 0) issue.custom_fields = customFields;
  await maybeMutate(instance, "POST", "/issues.json", { issue }, Boolean(options.execute));
}

async function updateIssue(args: Args) {
  const instanceName = args.shift();
  const id = args.shift();
  if (!id) throw new Error("issue id required");
  const [, instance] = resolveInstance(instanceName);
  const options = parseOptions(args);
  const issue: Record<string, unknown> = {};
  setIf(issue, "status_id", maybeNumber(options["status-id"]));
  setIf(issue, "assigned_to_id", maybeNumber(options["assigned-to-id"]));
  setIf(issue, "priority_id", maybeNumber(options["priority-id"]));
  setIf(issue, "tracker_id", maybeNumber(options["tracker-id"]));
  setIf(issue, "parent_issue_id", maybeNumber(options["parent-id"]));
  setIf(issue, "subject", options.subject);
  setIf(issue, "description", options.description);
  setIf(issue, "notes", options.notes);
  const customFields = parseCustomFields(options["custom-field"]);
  if (customFields && customFields.length > 0) issue.custom_fields = customFields;
  if (Object.keys(issue).length === 0) throw new Error("no update fields provided");
  await maybeMutate(instance, "PUT", `/issues/${encodeURIComponent(id)}.json`, { issue }, Boolean(options.execute));
}

async function commentIssue(args: Args) {
  const instanceName = args.shift();
  const id = args.shift();
  if (!id) throw new Error("issue id required");
  const [, instance] = resolveInstance(instanceName);
  const options = parseOptions(args);
  const issue = { notes: needString(options, "notes") };
  await maybeMutate(instance, "PUT", `/issues/${encodeURIComponent(id)}.json`, { issue }, Boolean(options.execute));
}

async function rawRequest(args: Args) {
  const instanceName = args.shift();
  const method = args.shift()?.toUpperCase();
  const path = args.shift();
  if (!method || !path) throw new Error("request requires <instance> <METHOD> <path>");
  const [, instance] = resolveInstance(instanceName);
  const options = parseOptions(args);
  const data = typeof options["data-json"] === "string" ? JSON.parse(options["data-json"] as string) : undefined;
  const mutates = !["GET", "HEAD", "OPTIONS"].includes(method);
  if (mutates && !options.execute) return dryRun(method, instance, path, data);
  console.log(JSON.stringify(await redmineFetch(instance, method, path, data), null, 2));
}

function setIf(target: Record<string, unknown>, key: string, value: unknown) {
  if (value !== undefined) target[key] = value;
}

async function maybeMutate(instance: InstanceConfig, method: string, path: string, body: unknown, execute: boolean) {
  if (!execute) return dryRun(method, instance, path, body);
  const result = await redmineFetch(instance, method, path, body);
  if (result !== undefined) console.log(JSON.stringify(result, null, 2));
  else console.log("ok");
}

function dryRun(method: string, instance: InstanceConfig, path: string, body: unknown) {
  console.log(
    JSON.stringify(
      {
        dryRun: true,
        method,
        url: buildInstanceUrl(instance, path),
        apiKeyEnv: instance.apiKeyEnv,
        body,
        note: "Add --execute to send this request. Run through with-secrets so the API key is populated.",
      },
      null,
      2
    )
  );
}
