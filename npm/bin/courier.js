#!/usr/bin/env node
"use strict";

const childProcess = require("node:child_process");
const path = require("node:path");

function binaryPath(platform = process.platform) {
  return path.join(__dirname, "..", "vendor", platform === "win32" ? "courier.exe" : "courier");
}

function run(args = process.argv.slice(2), platform = process.platform, spawn = childProcess.spawnSync, kill = process.kill) {
  const result = spawn(binaryPath(platform), args, { stdio: "inherit" });
  if (result.error) {
    throw result.error;
  }
  if (result.signal) {
    kill(process.pid, result.signal);
    return 1;
  }
  return result.status === null ? 1 : result.status;
}

if (require.main === module) {
  try {
    process.exitCode = run();
  } catch (error) {
    process.stderr.write(`courier: ${error.message}\nRun your package manager without --ignore-scripts so the native binary can be installed.\n`);
    process.exitCode = 1;
  }
}

module.exports = { binaryPath, run };
