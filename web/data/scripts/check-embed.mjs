import { mkdtemp, readdir, readFile, rm } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join, relative } from "node:path";
import { fileURLToPath } from "node:url";
import { build } from "vite";

const root = fileURLToPath(new URL("../", import.meta.url));
const committed = fileURLToPath(new URL("../../../internal/webdelivery/assets/", import.meta.url));
const temporary = await mkdtemp(join(tmpdir(), "courier-data-assets-"));

async function files(directory, base = directory) {
  const result = [];
  for (const entry of await readdir(directory, { withFileTypes: true })) {
    const path = join(directory, entry.name);
    if (entry.isDirectory()) result.push(...await files(path, base));
    else result.push(relative(base, path));
  }
  return result.sort();
}

try {
  await build({ root, logLevel: "silent", build: { outDir: temporary, emptyOutDir: true } });
  const expected = await files(committed);
  const actual = await files(temporary);
  if (JSON.stringify(actual) !== JSON.stringify(expected)) throw new Error(`embedded data UI file set is stale: ${actual.join(", ")}`);
  for (const name of actual) {
    const [built, stored] = await Promise.all([readFile(join(temporary, name)), readFile(join(committed, name))]);
    if (!built.equals(stored)) throw new Error(`embedded data UI asset is stale: ${name}`);
  }
  console.log(`verified ${actual.length} embedded data UI assets`);
} finally {
  await rm(temporary, { recursive: true, force: true });
}
