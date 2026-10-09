import { expect, it } from "vitest";
import { contractData } from "./contract";

it("contains only the generated public contract projection", () => {
  expect(contractData.contractVersion).toMatch(/^\d+\.\d+\.\d+$/);
  expect(contractData.targetRelease).toMatch(/^\d+\.\d+\.\d+$/);
  expect(contractData.commands.map((command) => command.name)).toContain("from");
  expect(contractData.commands.map((command) => command.name)).toContain("update");
  expect(contractData.routes.map((route) => route.name)).toContain("path-to-path");
  expect(contractData.flags.map((flag) => flag.name)).toContain("archive");
  expect(contractData.flags.map((flag) => flag.name)).toContain("force-source-creation");
  expect(contractData.flags.every((flag) => !flag.required && flag.description.en && flag.description.ru)).toBe(true);
  expect(contractData.commands.flatMap((command) => command.arguments).every((argument) => argument.description.en && argument.description.ru)).toBe(true);
  expect(contractData.scopeLabels["path-to-web"]?.ru).toContain("браузер");
  expect(contractData.endpoints.map((endpoint) => endpoint.name)).toContain("ssh");
  expect(contractData.examples.length).toBeGreaterThan(0);
  expect(JSON.stringify(contractData)).not.toContain('"status":"planned"');
});
