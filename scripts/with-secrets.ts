#!/usr/bin/env bun

import { execFileSync, spawn } from "node:child_process";
import { existsSync, readFileSync } from "node:fs";
import { join } from "node:path";

const BWS_CONTROL_KEYS = ["BWS_ACCESS_TOKEN", "BWS_PROJECT_ID", "BWS_SERVER_URL"] as const;

type Env = Record<string, string>;
type BwsSecret = { id: string; key: string; value: string };
type Options = {
  envFile: string;
  passBwsEnv: boolean;
  command: string[];
};
type SecretRef = {
  envKey: string;
  secretKey: string;
  optional: boolean;
};

const options = parseArgs(process.argv.slice(2));
const fileEnv = loadDotEnv(options.envFile);
const mergedEnv = mergeEnv(fileEnv, process.env);
const refs = collectSecretRefs(mergedEnv);
const childEnv: Env = { ...mergedEnv };

if (refs.length > 0) {
  const { BWS_ACCESS_TOKEN, BWS_PROJECT_ID, BWS_SERVER_URL } = mergedEnv;

  if (!BWS_ACCESS_TOKEN) {
    fail(`BWS_ACCESS_TOKEN missing in ${options.envFile} or current environment`);
  }

  if (!BWS_PROJECT_ID) {
    fail(`BWS_PROJECT_ID missing in ${options.envFile} or current environment`);
  }

  const bws = bwsBinary();
  const bwsEnv: Env = {
    ...mergedEnv,
    BWS_ACCESS_TOKEN,
    BWS_PROJECT_ID,
    ...(BWS_SERVER_URL ? { BWS_SERVER_URL } : {}),
  };

  const secrets = listProjectSecrets(bws, BWS_PROJECT_ID, bwsEnv);
  resolveSecretRefs(refs, secrets, childEnv);
}

if (!options.passBwsEnv) {
  for (const key of BWS_CONTROL_KEYS) {
    delete childEnv[key];
  }
}

runChild(options.command, childEnv);

function parseArgs(args: string[]): Options {
  const defaultEnvFile = join(process.env.XDG_CONFIG_HOME ?? join(process.env.HOME ?? "", ".config"), ".env");
  let envFile = defaultEnvFile;
  let passBwsEnv = false;
  const command: string[] = [];

  for (let i = 0; i < args.length; i += 1) {
    const arg = args[i];

    if (arg === "--") {
      command.push(...args.slice(i + 1));
      break;
    }

    if (arg === "-h" || arg === "--help") {
      usage(0);
    }

    if (arg === "--pass-bws-env") {
      passBwsEnv = true;
      continue;
    }

    if (arg === "--env-file") {
      const value = args[i + 1];
      if (!value) fail("--env-file requires a path");
      envFile = value;
      i += 1;
      continue;
    }

    if (arg.startsWith("--env-file=")) {
      envFile = arg.slice("--env-file=".length);
      continue;
    }

    command.push(...args.slice(i));
    break;
  }

  if (command.length === 0) {
    usage(1);
  }

  return { envFile, passBwsEnv, command };
}

function usage(exitCode: number): never {
  const configHome = process.env.XDG_CONFIG_HOME ?? "$HOME/.config";
  console.error(`Usage: with-secrets [options] -- <command> [args...]

Loads ${configHome}/.env, resolves bws:<secret-key> references from Bitwarden Secrets Manager, then runs the command.

Options:
  --env-file <path>   Use a different env file
  --pass-bws-env      Pass BWS_* control variables to the child command
  -h, --help          Show this help

Secret references:
  REQUIRED=bws:SECRET_KEY     fail if SECRET_KEY is missing or duplicated
  OPTIONAL=bws?:SECRET_KEY    unset OPTIONAL if SECRET_KEY is missing
`);
  process.exit(exitCode);
}

function loadDotEnv(path: string): Env {
  if (!existsSync(path)) {
    fail(`${path} missing`);
  }

  const env: Env = {};
  const lines = readFileSync(path, "utf8").split(/\r?\n/);

  for (let index = 0; index < lines.length; index += 1) {
    const lineNumber = index + 1;
    const raw = lines[index];
    const line = raw.trim();

    if (!line || line.startsWith("#")) continue;

    const equals = line.indexOf("=");
    if (equals <= 0) {
      fail(`${path}:${lineNumber}: expected KEY=value`);
    }

    const key = line.slice(0, equals).trim();
    const rawValue = line.slice(equals + 1).trim();

    if (!/^[A-Za-z_][A-Za-z0-9_]*$/.test(key)) {
      fail(`${path}:${lineNumber}: invalid env key ${JSON.stringify(key)}`);
    }

    env[key] = parseDotEnvValue(rawValue, path, lineNumber);
  }

  return env;
}

function parseDotEnvValue(value: string, path: string, lineNumber: number): string {
  if (value.length < 2) return value;

  const quote = value[0];
  const endQuote = value[value.length - 1];

  if ((quote === '"' || quote === "'") && endQuote === quote) {
    const inner = value.slice(1, -1);
    return quote === '"' ? inner.replace(/\\n/g, "\n").replace(/\\r/g, "\r").replace(/\\t/g, "\t").replace(/\\"/g, '"').replace(/\\\\/g, "\\") : inner;
  }

  if (quote === '"' || quote === "'" || endQuote === '"' || endQuote === "'") {
    fail(`${path}:${lineNumber}: unmatched quote`);
  }

  return value;
}

function mergeEnv(fileEnv: Env, processEnv: NodeJS.ProcessEnv): Env {
  const merged: Env = { ...fileEnv };

  for (const [key, value] of Object.entries(processEnv)) {
    if (typeof value === "string") {
      merged[key] = value;
    }
  }

  return merged;
}

function collectSecretRefs(env: Env): SecretRef[] {
  const refs: SecretRef[] = [];

  for (const [envKey, value] of Object.entries(env)) {
    if (value.startsWith("bws?:")) {
      refs.push({ envKey, secretKey: value.slice("bws?:".length), optional: true });
    } else if (value.startsWith("bws:")) {
      refs.push({ envKey, secretKey: value.slice("bws:".length), optional: false });
    }
  }

  for (const ref of refs) {
    if (!ref.secretKey) {
      fail(`${ref.envKey}: empty Bitwarden secret key reference`);
    }
  }

  return refs;
}

function bwsBinary(): string {
  try {
    return execFileSync("sh", ["-c", "command -v bws"], { stdio: ["ignore", "pipe", "ignore"] }).toString().trim();
  } catch {
    fail("bws not found on PATH; install Bitwarden Secrets Manager CLI as a user binary, e.g. ~/.local/bin/bws");
  }
}

function listProjectSecrets(bws: string, projectId: string, env: Env): BwsSecret[] {
  try {
    const raw = execFileSync(bws, ["secret", "list", projectId, "--output", "json"], {
      env,
      stdio: ["ignore", "pipe", "inherit"],
    });

    const parsed = JSON.parse(raw.toString()) as unknown;
    if (!Array.isArray(parsed)) {
      fail("bws secret list returned non-array JSON");
    }

    return parsed.map((secret) => {
      if (!isBwsSecret(secret)) {
        fail("bws secret list returned an unexpected secret shape");
      }
      return secret;
    });
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    fail(`bws secret list failed: ${message}`);
  }
}

function isBwsSecret(value: unknown): value is BwsSecret {
  return (
    typeof value === "object" &&
    value !== null &&
    typeof (value as BwsSecret).id === "string" &&
    typeof (value as BwsSecret).key === "string" &&
    typeof (value as BwsSecret).value === "string"
  );
}

function resolveSecretRefs(refs: SecretRef[], secrets: BwsSecret[], env: Env): void {
  const byKey = new Map<string, BwsSecret[]>();

  for (const secret of secrets) {
    const existing = byKey.get(secret.key) ?? [];
    existing.push(secret);
    byKey.set(secret.key, existing);
  }

  for (const ref of refs) {
    const matches = byKey.get(ref.secretKey) ?? [];

    if (matches.length > 1) {
      fail(`${ref.envKey}: multiple Bitwarden secrets found with key ${ref.secretKey}`);
    }

    if (matches.length === 0) {
      if (ref.optional) {
        delete env[ref.envKey];
        continue;
      }

      fail(`${ref.envKey}: Bitwarden secret key ${ref.secretKey} not found in project`);
    }

    env[ref.envKey] = matches[0].value;
  }
}

function runChild(command: string[], env: Env): never {
  const [cmd, ...args] = command;
  const child = spawn(cmd, args, { env, stdio: "inherit", shell: false });

  child.on("error", (error) => {
    fail(`failed to start ${cmd}: ${error.message}`);
  });

  child.on("exit", (code, signal) => {
    if (signal) {
      process.kill(process.pid, signal);
    } else {
      process.exit(code ?? 0);
    }
  });

  return undefined as never;
}

function fail(message: string): never {
  console.error(`[with-secrets] ${message}`);
  process.exit(1);
}
