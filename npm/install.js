"use strict";

const crypto = require("node:crypto");
const fs = require("node:fs");
const https = require("node:https");
const path = require("node:path");
const zlib = require("node:zlib");

const MAX_DOWNLOAD = 512 * 1024 * 1024;
const REPOSITORY = "iwonz/courier";

function targetFor(platform = process.platform, architecture = process.arch) {
  const operatingSystems = { darwin: "darwin", linux: "linux", win32: "windows" };
  const architectures = { x64: "amd64", arm64: "arm64" };
  const osName = operatingSystems[platform];
  const archName = architectures[architecture];
  if (!osName || !archName) {
    throw new Error(`unsupported npm platform: ${platform}/${architecture}`);
  }
  return { os: osName, arch: archName, binary: osName === "windows" ? "courier.exe" : "courier" };
}

function parseChecksum(manifest, assetName) {
  for (const line of manifest.split(/\r?\n/u)) {
    const match = /^([0-9a-f]{64})\s+\*?(.+)$/iu.exec(line.trim());
    if (match && match[2] === assetName) {
      return match[1].toLowerCase();
    }
  }
  throw new Error(`checksums.txt has no entry for ${assetName}`);
}

function download(url, redirects = 5) {
  return new Promise((resolve, reject) => {
    const parsed = new URL(url);
    if (parsed.protocol !== "https:") {
      reject(new Error(`refusing non-HTTPS download: ${parsed.protocol}`));
      return;
    }
    const request = https.get(parsed, { headers: { "User-Agent": "@iwonz/courier" } }, (response) => {
      if (response.statusCode >= 300 && response.statusCode < 400 && response.headers.location) {
        response.resume();
        if (redirects === 0) {
          reject(new Error("too many release download redirects"));
          return;
        }
        resolve(download(new URL(response.headers.location, parsed).toString(), redirects - 1));
        return;
      }
      if (response.statusCode !== 200) {
        response.resume();
        reject(new Error(`release download failed with HTTP ${response.statusCode}`));
        return;
      }
      const declaredLength = Number(response.headers["content-length"] || 0);
      if (declaredLength > MAX_DOWNLOAD) {
        response.resume();
        reject(new Error("release download exceeds the size limit"));
        return;
      }
      const chunks = [];
      let length = 0;
      response.on("data", (chunk) => {
        length += chunk.length;
        if (length > MAX_DOWNLOAD) {
          request.destroy(new Error("release download exceeds the size limit"));
          return;
        }
        chunks.push(chunk);
      });
      response.on("end", () => resolve(Buffer.concat(chunks, length)));
    });
    request.on("error", reject);
  });
}

function tarString(block, offset, length) {
  return block.subarray(offset, offset + length).toString("utf8").replace(/\0.*$/u, "");
}

function safeTarName(name) {
  if (!name || name.includes("\\") || name.startsWith("/") || name.split("/").some((part) => part === ".." || part === "")) {
    throw new Error(`unsafe release archive entry: ${JSON.stringify(name)}`);
  }
}

function extractBinary(archive, binaryName) {
  const tar = zlib.gunzipSync(archive, { maxOutputLength: MAX_DOWNLOAD });
  let offset = 0;
  let result;
  while (offset + 512 <= tar.length) {
    const header = tar.subarray(offset, offset + 512);
    if (header.every((byte) => byte === 0)) {
      break;
    }
    const name = [tarString(header, 345, 155), tarString(header, 0, 100)].filter(Boolean).join("/");
    safeTarName(name);
    const sizeText = tarString(header, 124, 12).trim();
    const size = Number.parseInt(sizeText || "0", 8);
    if (!Number.isSafeInteger(size) || size < 0) {
      throw new Error(`invalid release archive size for ${name}`);
    }
    const dataStart = offset + 512;
    const dataEnd = dataStart + size;
    if (dataEnd > tar.length) {
      throw new Error(`truncated release archive entry: ${name}`);
    }
    if (name === binaryName) {
      const type = header[156];
      if (result || (type !== 0 && type !== 48)) {
        throw new Error(`invalid release archive entry: ${binaryName}`);
      }
      result = Buffer.from(tar.subarray(dataStart, dataEnd));
    }
    offset = dataStart + Math.ceil(size / 512) * 512;
  }
  if (!result) {
    throw new Error(`release archive has no ${binaryName}`);
  }
  return result;
}

async function install(options = {}) {
  const packageRoot = options.packageRoot || __dirname;
  const metadata = JSON.parse(fs.readFileSync(path.join(packageRoot, "package.json"), "utf8"));
  const version = options.version || metadata.version;
  if (!/^\d+\.\d+\.\d+(?:-[0-9A-Za-z.-]+)?$/u.test(version)) {
    throw new Error(`invalid Courier package version: ${version}`);
  }
  const target = targetFor(options.platform, options.architecture);
  const asset = `courier_${version}_${target.os}_${target.arch}.tar.gz`;
  const releaseBase = options.releaseBase || process.env.COURIER_RELEASE_BASE_URL || `https://github.com/${REPOSITORY}/releases/download`;
  const fetch = options.download || download;
  const [archive, manifestBuffer] = await Promise.all([
    fetch(`${releaseBase}/v${version}/${asset}`),
    fetch(`${releaseBase}/v${version}/checksums.txt`),
  ]);
  const expected = parseChecksum(manifestBuffer.toString("utf8"), asset);
  const actual = crypto.createHash("sha256").update(archive).digest("hex");
  if (actual !== expected) {
    throw new Error(`checksum mismatch for ${asset}`);
  }
  const binary = extractBinary(archive, target.binary);
  const vendor = path.join(packageRoot, "vendor");
  fs.mkdirSync(vendor, { recursive: true, mode: 0o700 });
  const staged = path.join(vendor, `.${target.binary}.${process.pid}.tmp`);
  const destination = path.join(vendor, target.binary);
  try {
    fs.writeFileSync(staged, binary, { mode: 0o700, flag: "wx" });
    fs.renameSync(staged, destination);
  } finally {
    fs.rmSync(staged, { force: true });
  }
  return destination;
}

if (require.main === module) {
  install().then(
    (destination) => process.stdout.write(`Installed Courier binary at ${destination}\n`),
    (error) => {
      process.stderr.write(`@iwonz/courier install failed: ${error.message}\n`);
      process.exitCode = 1;
    },
  );
}

module.exports = { download, extractBinary, install, parseChecksum, safeTarName, targetFor };
