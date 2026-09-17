import { gzipSync } from "node:zlib";
import { readFile } from "node:fs/promises";
import { resolve } from "node:path";

const limits = {
  "assets/landing.js": 45 * 1024,
  "assets/landing-index.css": 9 * 1024,
};

for (const [relativePath, limit] of Object.entries(limits)) {
  const data = await readFile(resolve("dist", relativePath));
  const bytes = gzipSync(data, { level: 9 }).length;
  if (bytes > limit) throw new Error(`${relativePath} exceeds its gzip budget: ${bytes} > ${limit}`);
  console.log(`${relativePath}: ${bytes}/${limit} gzip bytes`);
}
