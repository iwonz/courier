"use strict";

const assert = require("node:assert/strict");
const crypto = require("node:crypto");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const test = require("node:test");
const zlib = require("node:zlib");

const installer = require("../install.js");

test("target selection supports the complete release matrix", () => {
  assert.deepEqual(installer.targetFor("darwin", "x64"), { os: "darwin", arch: "amd64", binary: "courier" });
  assert.deepEqual(installer.targetFor("linux", "arm64"), { os: "linux", arch: "arm64", binary: "courier" });
  assert.deepEqual(installer.targetFor("win32", "x64"), { os: "windows", arch: "amd64", binary: "courier.exe" });
  assert.throws(() => installer.targetFor("freebsd", "x64"), /unsupported npm platform/u);
  assert.throws(() => installer.targetFor("linux", "ia32"), /unsupported npm platform/u);
});

test("checksum parsing selects an exact asset", () => {
  const digest = "A".repeat(64);
  assert.equal(installer.parseChecksum(`${digest} *asset.tar.gz\n`, "asset.tar.gz"), digest.toLowerCase());
  assert.throws(() => installer.parseChecksum(`${digest}  other.tar.gz`, "asset.tar.gz"), /has no entry/u);
});

test("tar extraction accepts one regular binary and rejects unsafe input", () => {
  const archive = makeTar([{ name: "README.md", data: Buffer.from("docs") }, { name: "courier", data: Buffer.from("binary") }]);
  assert.equal(installer.extractBinary(archive, "courier").toString(), "binary");
  assert.throws(() => installer.extractBinary(makeTar([{ name: "../courier", data: Buffer.alloc(0) }]), "courier"), /unsafe/u);
  assert.throws(() => installer.extractBinary(makeTar([{ name: "/courier", data: Buffer.alloc(0) }]), "courier"), /unsafe/u);
  assert.throws(() => installer.extractBinary(makeTar([{ name: "dir\\courier", data: Buffer.alloc(0) }]), "courier"), /unsafe/u);
  assert.throws(() => installer.extractBinary(makeTar([{ name: "courier", data: Buffer.alloc(0), type: "5" }]), "courier"), /invalid release archive entry/u);
  assert.throws(() => installer.extractBinary(makeTar([{ name: "courier", data: Buffer.alloc(0) }, { name: "courier", data: Buffer.alloc(0) }]), "courier"), /invalid release archive entry/u);
  assert.throws(() => installer.extractBinary(makeTar([{ name: "README", data: Buffer.alloc(0) }]), "courier"), /has no courier/u);

  const invalidSize = rawHeader("courier", Buffer.alloc(0));
  invalidSize.fill(122, 124, 136);
  assert.throws(() => installer.extractBinary(zlib.gzipSync(invalidSize), "courier"), /invalid release archive size/u);
  const truncated = rawHeader("courier", Buffer.from("x"), 100);
  assert.throws(() => installer.extractBinary(zlib.gzipSync(truncated), "courier"), /truncated/u);
});

test("install verifies before atomically replacing the packaged binary", async (t) => {
  const packageRoot = fs.mkdtempSync(path.join(os.tmpdir(), "courier-npm-test-"));
  t.after(() => fs.rmSync(packageRoot, { recursive: true, force: true }));
  fs.writeFileSync(path.join(packageRoot, "package.json"), JSON.stringify({ version: "1.2.3" }));
  const archive = makeTar([{ name: "courier", data: Buffer.from("native") }]);
  const asset = "courier_1.2.3_linux_amd64.tar.gz";
  const digest = crypto.createHash("sha256").update(archive).digest("hex");
  const fakeDownload = async (url) => url.endsWith("checksums.txt") ? Buffer.from(`${digest}  ${asset}\n`) : archive;
  const destination = await installer.install({ packageRoot, platform: "linux", architecture: "x64", download: fakeDownload, releaseBase: "https://example.invalid" });
  assert.equal(fs.readFileSync(destination, "utf8"), "native");
  assert.equal(fs.readdirSync(path.dirname(destination)).some((name) => name.endsWith(".tmp")), false);

  fs.writeFileSync(destination, "old");
  await assert.rejects(
    installer.install({ packageRoot, platform: "linux", architecture: "x64", download: async (url) => url.endsWith("checksums.txt") ? Buffer.from(`${"0".repeat(64)}  ${asset}\n`) : archive }),
    /checksum mismatch/u,
  );
  assert.equal(fs.readFileSync(destination, "utf8"), "old");
  await assert.rejects(installer.install({ packageRoot, version: "development", download: fakeDownload }), /invalid Courier package version/u);
});

test("archive safety rejects empty path components", () => {
  for (const name of ["", "dir//courier", "dir/../courier"]) {
    assert.throws(() => installer.safeTarName(name), /unsafe/u);
  }
  assert.doesNotThrow(() => installer.safeTarName("dir/courier"));
});

function makeTar(entries) {
  const blocks = [];
  for (const entry of entries) {
    const data = entry.data || Buffer.alloc(0);
    const header = rawHeader(entry.name, data, data.length, entry.type || "0");
    blocks.push(header, data, Buffer.alloc((512 - (data.length % 512)) % 512));
  }
  blocks.push(Buffer.alloc(1024));
  return zlib.gzipSync(Buffer.concat(blocks));
}

function rawHeader(name, data, declaredSize = data.length, type = "0") {
  const header = Buffer.alloc(512);
  header.write(name, 0, 100, "utf8");
  writeOctal(header, 100, 8, 0o755);
  writeOctal(header, 108, 8, 0);
  writeOctal(header, 116, 8, 0);
  writeOctal(header, 124, 12, declaredSize);
  writeOctal(header, 136, 12, 0);
  header.fill(32, 148, 156);
  header.write(type, 156, 1, "ascii");
  header.write("ustar\0", 257, 6, "ascii");
  header.write("00", 263, 2, "ascii");
  const checksum = header.reduce((sum, value) => sum + value, 0);
  writeOctal(header, 148, 8, checksum);
  return header;
}

function writeOctal(buffer, offset, length, value) {
  const text = value.toString(8).padStart(length - 2, "0") + "\0 ";
  buffer.write(text, offset, length, "ascii");
}
