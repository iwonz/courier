"use strict";

const assert = require("node:assert/strict");
const path = require("node:path");
const test = require("node:test");

const runner = require("../bin/courier.js");

test("binary path uses the packaged host executable", () => {
  assert.equal(path.basename(runner.binaryPath("linux")), "courier");
  assert.equal(path.basename(runner.binaryPath("win32")), "courier.exe");
});

test("runner forwards arguments and stable status", () => {
  let invocation;
  const status = runner.run(["version"], "linux", (binary, args, options) => {
    invocation = { binary, args, options };
    return { status: 7 };
  });
  assert.equal(status, 7);
  assert.deepEqual(invocation.args, ["version"]);
  assert.deepEqual(invocation.options, { stdio: "inherit" });
  assert.throws(() => runner.run([], "linux", () => ({ error: new Error("spawn") })), /spawn/u);
  assert.equal(runner.run([], "linux", () => ({ status: null })), 1);
});

test("runner mirrors a child signal", () => {
  let signal;
  assert.equal(runner.run([], "linux", () => ({ signal: "SIGTERM" }), (_pid, value) => { signal = value; }), 1);
  assert.equal(signal, "SIGTERM");
});
