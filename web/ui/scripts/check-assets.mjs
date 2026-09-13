import { createHash } from "node:crypto";
import { readFile, readdir } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const packageRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const assetRoot = resolve(packageRoot, "assets");
const manifest = JSON.parse(await readFile(resolve(assetRoot, "provenance.json"), "utf8"));

if (manifest.schemaVersion !== 1 || !Array.isArray(manifest.assets) || manifest.assets.length === 0) {
  throw new Error("asset provenance manifest is invalid");
}

const expected = new Set(["provenance.json"]);
for (const asset of manifest.assets) {
  if (!/^[a-z0-9][a-z0-9.-]*$/.test(asset.path) || expected.has(asset.path)) {
    throw new Error(`invalid or duplicate asset path: ${asset.path}`);
  }
  expected.add(asset.path);
  const data = await readFile(resolve(assetRoot, asset.path));
  const digest = createHash("sha256").update(data).digest("hex");
  if (data.length !== asset.bytes || digest !== asset.sha256) {
    throw new Error(`asset integrity mismatch: ${asset.path}`);
  }
}

const actual = new Set(await readdir(assetRoot));
if (actual.size !== expected.size || [...actual].some((name) => !expected.has(name))) {
  throw new Error("asset directory and provenance manifest differ");
}

const searchable = [
  await readFile(resolve(packageRoot, "NOTICE.md"), "utf8"),
  ...await Promise.all(manifest.assets.filter((asset) => asset.mediaType.includes("svg")).map((asset) => readFile(resolve(assetRoot, asset.path), "utf8"))),
].join("\n").toLowerCase();
if (searchable.includes("local.adguard.org") || searchable.includes("<script")) {
  throw new Error("executable reference content was retained");
}

console.log(`verified ${manifest.assets.length} Courier identity assets`);
