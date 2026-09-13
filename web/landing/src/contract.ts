import generated from "./contract.generated.json";

export interface LandingCommand {
  name: string;
  usage: string;
  system: boolean;
  flags: string[];
}

export interface LandingFlag {
  name: string;
  syntax: string;
  repeatable: boolean;
  default: string;
  appliesTo: string[];
  conflicts: string[];
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
