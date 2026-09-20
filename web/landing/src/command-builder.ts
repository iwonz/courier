import type { LandingArgument, LandingCommand, LandingFlag } from "./contract";

export type ShellMode = "posix" | "powershell";
export type ArgumentValues = Readonly<Record<string, string>>;
export type FlagValues = Readonly<Record<string, readonly string[]>>;

export interface BuiltCommand {
  readonly command: string;
  readonly valid: boolean;
  readonly invalidFields: readonly string[];
}

const safeToken = /^[A-Za-z0-9_@%+=:,./-]+$/;
const uuid = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
const unsigned = /^(0|[1-9][0-9]*)$/;

export function quoteToken(value: string, shell: ShellMode): string {
  if (safeToken.test(value)) return value;
  if (shell === "powershell") return `'${value.replaceAll("'", "''")}'`;
  return `'${value.replaceAll("'", `'"'"'`)}'`;
}

export function flagEnabled(values: FlagValues, name: string): boolean {
  return (values[name] ?? []).some((value) => value.trim() !== "");
}

export function argumentValid(argument: LandingArgument, value: string, flags: FlagValues): boolean {
  if (argument.omitWhenFlag && flagEnabled(flags, argument.omitWhenFlag)) return true;
  const clean = value.trim();
  if (!clean) return !argument.required;
  if (argument.kind === "uuid") return uuid.test(clean);
  return true;
}

export function flagValueValid(flag: LandingFlag, value: string): boolean {
  const clean = value.trim();
  if (flag.valueKind === "boolean") return clean === "true";
  if (!clean) return false;
  if (flag.valueKind === "unsigned") return unsigned.test(clean);
  if (flag.valueKind === "enum") return (flag.choices ?? []).includes(clean);
  return true;
}

export function updateFlagValues(current: FlagValues, flag: LandingFlag, values: readonly string[], allFlags: readonly LandingFlag[]): FlagValues {
  const next: Record<string, readonly string[]> = { ...current, [flag.name]: [...values] };
  if (values.some((value) => value.trim() !== "")) {
    for (const candidate of allFlags) {
      if ((flag.conflicts ?? []).includes(candidate.name) || (candidate.conflicts ?? []).includes(flag.name)) delete next[candidate.name];
    }
  } else {
    for (const candidate of allFlags) {
      if ((candidate.requires ?? []).includes(flag.name)) delete next[candidate.name];
    }
  }
  return next;
}

export function buildCommand(command: LandingCommand | undefined, argumentsByName: ArgumentValues, flagsByName: FlagValues, flags: readonly LandingFlag[], shell: ShellMode): BuiltCommand {
  if (!command) return { command: "", valid: false, invalidFields: [] };
  const tokens = ["courier", ...command.path.split(" ").filter(Boolean)];
  const invalidFields: string[] = [];
  for (const argument of command.arguments) {
    const value = (argumentsByName[argument.name] ?? "").trim();
    if (!argumentValid(argument, value, flagsByName)) invalidFields.push(argument.name);
    if (!value || argument.omitWhenFlag && flagEnabled(flagsByName, argument.omitWhenFlag)) continue;
    if (argument.prefix) tokens.push(argument.prefix);
    if (argument.kind === "command-path") tokens.push(...value.split(" ").filter(Boolean));
    else tokens.push(value);
  }
  for (const name of command.flags ?? []) {
    const flag = flags.find((candidate) => candidate.name === name);
    if (!flag) continue;
    const values = (flagsByName[name] ?? []).map((value) => value.trim()).filter(Boolean);
    if (!values.length) continue;
    if ((flag.requires ?? []).some((dependency) => !flagEnabled(flagsByName, dependency))) invalidFields.push(name);
    if (flag.valueKind === "boolean") {
      if (flagValueValid(flag, values[0]!)) tokens.push(`--${flag.name}`);
      else invalidFields.push(name);
      continue;
    }
    for (const value of values) {
      if (!flagValueValid(flag, value)) invalidFields.push(name);
      tokens.push(`--${flag.name}`, value);
    }
  }
  return { command: tokens.map((token) => quoteToken(token, shell)).join(" "), valid: invalidFields.length === 0, invalidFields };
}
