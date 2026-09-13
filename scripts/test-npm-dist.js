"use strict";

const assert = require("node:assert/strict");
const childProcess = require("node:child_process");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");

const { install, targetFor } = require("../npm/install.js");

async function main() {
  const metadata = JSON.parse(fs.readFileSync(path.join("dist", "metadata.json"), "utf8"));
  const target = targetFor();
  const asset = `courier_${metadata.version}_${target.os}_${target.arch}.tar.gz`;
  const packageRoot = fs.mkdtempSync(path.join(os.tmpdir(), "courier-npm-dist-"));
  try {
    fs.writeFileSync(path.join(packageRoot, "package.json"), JSON.stringify({ version: metadata.version }));
    const destination = await install({
      packageRoot,
      version: metadata.version,
      download: async (url) => fs.readFileSync(url.endsWith("checksums.txt") ? path.join("dist", "checksums.txt") : path.join("dist", asset)),
    });
    const output = childProcess.execFileSync(destination, ["version"], { encoding: "utf8" });
    assert.match(output, new RegExp(`^courier ${metadata.version.replace(/[.*+?^${}()|[\]\\]/gu, "\\$&")}`));
    process.stdout.write(`Verified npm bootstrap against the ${metadata.version} snapshot\n`);
  } finally {
    fs.rmSync(packageRoot, { recursive: true, force: true });
  }
}

main().catch((error) => {
  process.stderr.write(`${error.stack || error.message}\n`);
  process.exitCode = 1;
});
