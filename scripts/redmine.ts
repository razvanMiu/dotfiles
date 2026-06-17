#!/usr/bin/env bun

type InstanceConfig = { url: string; apiKeyEnv: string };
type Config = { defaultInstance?: string; instances: Record<string, InstanceConfig> };
type Args = string[];

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
  issues list <instance> [--project ID] [--assigned-to me|ID] [--status open|closed|*|ID] [--limit N] [--query key=value]
  issue get <instance> <id> [--include journals,attachments]
  issue create <instance> --project ID --subject TEXT [--description TEXT] [--tracker-id ID] [--priority-id ID] [--assigned-to-id ID] [--execute]
  issue update <instance> <id> [--status-id ID] [--assigned-to-id ID] [--subject TEXT] [--description TEXT] [--notes TEXT] [--execute]
  issue comment <instance> <id> --notes TEXT [--execute]
  request <instance> <METHOD> <path> [--data-json JSON] [--execute]

Mutating issue commands are dry-run by default. Add --execute to send POST/PUT/DELETE.
Set REDMINE_CONFIG to override ${CONFIG_PATH}.
Use with-secrets -- redmine.ts ... so API key env vars are populated safely.`);
}

function loadConfig(): Config {
  const file = Bun.file(CONFIG_PATH);
  if (!file.exists()) {
    throw new Error(`config not found: ${CONFIG_PATH}; copy ~/.config/redmine/config.example.json to config.json and edit URLs/env names`);
  }
  const config = JSON.parse(require("fs").readFileSync(CONFIG_PATH, "utf8")) as Config;
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

async function redmineFetch(instance: InstanceConfig, method: string, path: string, body?: unknown) {
  const url = path.startsWith("http") ? path : `${instance.url}${path.startsWith("/") ? "" : "/"}${path}`;
  const response = await fetch(url, {
    method,
    headers: {
      "Accept": "application/json",
      "Content-Type": "application/json",
      "X-Redmine-API-Key": apiKey(instance),
    },
    body: body === undefined ? undefined : JSON.stringify(body),
  });
  const text = await response.text();
  if (!response.ok) {
    throw new Error(`${method} ${url} failed: HTTP ${response.status} ${response.statusText}${text ? `\n${text}` : ""}`);
  }
  if (!text.trim()) return undefined;
  try { return JSON.parse(text); } catch { return text; }
}

function parseOptions(args: Args): Record<string, string | boolean | string[]> {
  const options: Record<string, string | boolean | string[]> = {};
  for (let i = 0; i < args.length; i++) {
    const arg = args[i];
    if (!arg.startsWith("--")) throw new Error(`unexpected argument: ${arg}`);
    const key = arg.slice(2);
    if (key === "execute" || key === "dry-run") {
      options[key] = true;
      continue;
    }
    const value = args[++i];
    if (value === undefined) throw new Error(`missing value for --${key}`);
    if (key === "query") {
      const list = Array.isArray(options[key]) ? options[key] as string[] : [];
      list.push(value);
      options[key] = list;
    } else {
      options[key] = value;
    }
  }
  return options;
}

function needString(options: Record<string, unknown>, key: string): string {
  const value = options[key];
  if (typeof value !== "string" || value.length === 0) throw new Error(`--${key} is required`);
  return value;
}

function maybeNumber(value: unknown): number | undefined {
  if (value === undefined) return undefined;
  const n = Number(value);
  if (!Number.isInteger(n)) throw new Error(`expected integer, got ${value}`);
  return n;
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
  const params: Record<string, string | number | undefined> = {
    project_id: options.project as string | undefined,
    assigned_to_id: options["assigned-to"] as string | undefined,
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
  setIf(issue, "subject", options.subject);
  setIf(issue, "description", options.description);
  setIf(issue, "notes", options.notes);
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
  console.log(JSON.stringify({
    dryRun: true,
    method,
    url: `${instance.url}${path}`,
    apiKeyEnv: instance.apiKeyEnv,
    body,
    note: "Add --execute to send this request. Run through with-secrets so the API key is populated.",
  }, null, 2));
}
