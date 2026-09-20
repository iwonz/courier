import { describe, expect, it } from "vitest";
import { argumentValid, buildCommand, flagEnabled, flagValueValid, quoteToken, updateFlagValues } from "./command-builder";
import { contractData, type LandingArgument, type LandingFlag } from "./contract";

const flag = (name: string): LandingFlag => contractData.flags.find((candidate) => candidate.name === name)!;
const command = (name: string) => contractData.commands.find((candidate) => candidate.name === name)!;

describe("command builder", () => {
  it("quotes safe and unsafe POSIX and PowerShell tokens", () => {
    expect(quoteToken("./safe-path", "posix")).toBe("./safe-path");
    expect(quoteToken("it's here", "posix")).toBe("'it'\"'\"'s here'");
    expect(quoteToken("it's here", "powershell")).toBe("'it''s here'");
  });

  it("validates argument and parameter kinds", () => {
    const optional: LandingArgument = { name: "command", kind: "command-path", required: false };
    const required: LandingArgument = { name: "source", kind: "endpoint", required: true };
    const stop: LandingArgument = { name: "uuid", kind: "uuid", required: true, omitWhenFlag: "all" };
    expect(argumentValid(optional, "", {})).toBe(true);
    expect(argumentValid(required, "", {})).toBe(false);
    expect(argumentValid(required, "./source", {})).toBe(true);
    expect(argumentValid(stop, "bad", {})).toBe(false);
    expect(argumentValid(stop, "bad", { all: ["true"] })).toBe(true);
    expect(argumentValid(stop, "550e8400-e29b-41d4-a716-446655440000", {})).toBe(true);
    expect(flagValueValid(flag("archive"), "true")).toBe(true);
    expect(flagValueValid(flag("archive"), "false")).toBe(false);
    expect(flagValueValid(flag("listen"), "")).toBe(false);
    expect(flagValueValid(flag("listen"), "127.0.0.1:8080")).toBe(true);
    expect(flagValueValid(flag("auth-attempts"), "0")).toBe(true);
    expect(flagValueValid(flag("auth-attempts"), "01")).toBe(false);
    expect(flagValueValid(flag("auth"), "basic")).toBe(true);
    expect(flagValueValid(flag("auth"), "token")).toBe(false);
    expect(flagValueValid({ ...flag("auth"), choices: null }, "basic")).toBe(false);
    expect(flagEnabled({}, "archive")).toBe(false);
    expect(flagEnabled({ archive: ["", "true"] }, "archive")).toBe(true);
  });

  it("clears reciprocal conflicts and dependent values", () => {
    expect(updateFlagValues({ extract: ["true"] }, flag("archive"), ["true"], contractData.flags)).toEqual({ archive: ["true"] });
    expect(updateFlagValues({ extract: ["true"], "max-extracted-size": ["1GiB"] }, flag("extract"), [], contractData.flags)).toEqual({ extract: [] });
    expect(updateFlagValues({ archive: ["true"] }, flag("extract"), [""], contractData.flags)).toEqual({ archive: ["true"], extract: [""] });
  });

  it("builds ordered commands without materializing defaults", () => {
    const result = buildCommand(command("from"), { source: "folder with space", destination: "host:/srv/it's" }, {
      extract: ["true"], exclude: ["*.tmp", "old files/*"], "max-extracted-size": ["2GiB"], auth: [""],
    }, contractData.flags, "posix");
    expect(result.valid).toBe(true);
    expect(result.command).toBe("courier from 'folder with space' to 'host:/srv/it'\"'\"'s' --extract --exclude '*.tmp' --exclude 'old files/*' --max-extracted-size 2GiB");
    expect(result.command).not.toContain("127.0.0.1:8080");
    expect(buildCommand(command("from"), { source: "a", destination: "b" }, { "auth-attempts": ["01"] }, contractData.flags, "posix")).toEqual(expect.objectContaining({ valid: false, invalidFields: ["auth-attempts"] }));
  });

  it("supports --all, help paths, empty selection, and missing contract flags", () => {
    expect(buildCommand(command("servers-stop"), { uuid: "bad" }, { all: ["true"] }, contractData.flags, "powershell")).toEqual({ command: "courier servers stop --all", valid: true, invalidFields: [] });
    expect(buildCommand(command("servers-stop"), { uuid: "bad" }, {}, contractData.flags, "posix").valid).toBe(false);
    expect(buildCommand(command("help"), { command: "servers stop" }, {}, contractData.flags, "posix").command).toBe("courier help servers stop");
    expect(buildCommand(undefined, {}, {}, contractData.flags, "posix")).toEqual({ command: "", valid: false, invalidFields: [] });
    const synthetic = { ...command("servers"), flags: ["missing"] };
    expect(buildCommand(synthetic, {}, { missing: ["value"] }, contractData.flags, "posix")).toEqual({ command: "courier servers", valid: true, invalidFields: [] });
    expect(buildCommand({ ...command("servers-stop"), arguments: [] }, {}, { all: ["false"] }, contractData.flags, "posix")).toEqual(expect.objectContaining({ valid: false, invalidFields: ["all"] }));
    expect(buildCommand(command("from"), { source: "a", destination: "b" }, { "max-extracted-size": ["2GiB"] }, contractData.flags, "posix")).toEqual(expect.objectContaining({ valid: false, invalidFields: ["max-extracted-size"] }));
  });
});
