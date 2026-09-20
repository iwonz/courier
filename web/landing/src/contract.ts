import generated from "./contract.generated.json";

export interface LandingCommand {
  name: string;
  path: string;
  usage: string;
  system: boolean;
  arguments: LandingArgument[];
  flags: string[] | null;
}

export interface LandingArgument {
  name: string;
  kind: "endpoint" | "uuid" | "command-path";
  required: boolean;
  prefix?: string;
  omitWhenFlag?: string;
}

export interface LandingFlag {
  name: string;
  syntax: string;
  valueKind: "boolean" | "text" | "unsigned" | "enum" | "host-port" | "ip-cidr" | "path" | "regex" | "quantity" | "rate";
  choices: string[] | null;
  placeholder: string;
  repeatable: boolean;
  default: string;
  appliesTo: string[];
  conflicts: string[] | null;
  requires: string[] | null;
}

export interface LandingRoute {
  name: string;
  source: string[];
  destination: string[];
  allowedFlags: string[];
}

export interface LandingContract {
  contractVersion: string;
  targetRelease: string;
  endpoints: { name: string; syntax: string }[];
  commands: LandingCommand[];
  flags: LandingFlag[];
  routes: LandingRoute[];
  examples: string[];
}

export const contractData = generated as LandingContract;
