import { createHash } from "node:crypto";
import { readFile, readdir } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const packageRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const assetRoot = resolve(packageRoot, "assets");
const manifest = JSON.parse(await readFile(resolve(assetRoot, "provenance.json"), "utf8"));
const brandAssetRoot = resolve(packageRoot, "src", "brand-assets");
const brandManifest = JSON.parse(await readFile(resolve(brandAssetRoot, "provenance.json"), "utf8"));

if (manifest.schemaVersion !== 2 || !Array.isArray(manifest.assets) || manifest.assets.length === 0) {
  throw new Error("asset provenance manifest is invalid");
}
if (brandManifest.schemaVersion !== 1 || !Array.isArray(brandManifest.assets) || brandManifest.assets.length === 0) {
  throw new Error("brand asset provenance manifest is invalid");
}

function uint24(data, offset) {
  return data[offset] | (data[offset + 1] << 8) | (data[offset + 2] << 16);
}

function rasterDimensions(data, mediaType) {
  if (mediaType === "image/png") {
    if (data.toString("hex", 0, 8) !== "89504e470d0a1a0a") throw new Error("invalid PNG header");
    return { width: data.readUInt32BE(16), height: data.readUInt32BE(20) };
  }
  if (data.toString("ascii", 0, 4) !== "RIFF" || data.toString("ascii", 8, 12) !== "WEBP") throw new Error("invalid WebP header");
  const chunk = data.toString("ascii", 12, 16);
  if (chunk === "VP8X") return { width: uint24(data, 24) + 1, height: uint24(data, 27) + 1 };
  if (chunk === "VP8 ") return { width: data.readUInt16LE(26) & 0x3fff, height: data.readUInt16LE(28) & 0x3fff };
  if (chunk === "VP8L") {
    const bits = data.readUInt32LE(21);
    return { width: (bits & 0x3fff) + 1, height: ((bits >>> 14) & 0x3fff) + 1 };
  }
  throw new Error(`unsupported WebP chunk: ${chunk}`);
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
  if (asset.mediaType === "image/png" || asset.mediaType === "image/webp") {
    if (!asset.prompt || !Number.isInteger(asset.width) || !Number.isInteger(asset.height)) {
      throw new Error(`raster provenance is incomplete: ${asset.path}`);
    }
    const dimensions = rasterDimensions(data, asset.mediaType);
    if (dimensions.width !== asset.width || dimensions.height !== asset.height) {
      throw new Error(`asset dimensions mismatch: ${asset.path}`);
    }
  }
}

const actual = new Set(await readdir(assetRoot));
if (actual.size !== expected.size || [...actual].some((name) => !expected.has(name))) {
  throw new Error("asset directory and provenance manifest differ");
}

const landingAssets = manifest.assets.filter((asset) => asset.sequence?.name === "organic-panorama");
const landingBytes = landingAssets.reduce((total, asset) => total + asset.bytes, 0);
if (landingAssets.length !== 2 || landingBytes > 700 * 1024) {
  throw new Error("landing panorama asset count or aggregate budget is invalid");
}
for (const orientation of ["wide", "portrait"]) {
  const sequence = landingAssets.filter((asset) => asset.sequence.orientation === orientation);
  if (sequence.length !== 1 || sequence.some((asset) => asset.sequence.position !== 1 || asset.sequence.total !== 1 || asset.bytes > 420 * 1024)) {
    throw new Error(`landing ${orientation} panorama sequence or per-file budget is invalid`);
  }
}

const compactMark = manifest.assets.find((asset) => asset.path === "courier-mark-v2.webp");
if (!compactMark || compactMark.mediaType !== "image/webp" || compactMark.bytes > 96 * 1024) {
  throw new Error("ImageGen compact mark is missing or exceeds its budget");
}

const productScenes = manifest.assets.filter((asset) => /^(delivery-access|admin-operations)-/.test(asset.path));
if (productScenes.length !== 4 || productScenes.some((asset) => asset.bytes > 200 * 1024)) {
  throw new Error("delivery/admin scene count or per-file budget is invalid");
}

const expectedBrandAssets = new Set(["provenance.json"]);
for (const asset of brandManifest.assets) {
  if (!/^[a-z0-9][a-z0-9.-]*\.svg$/.test(asset.path) || expectedBrandAssets.has(asset.path)) {
    throw new Error(`invalid or duplicate brand asset path: ${asset.path}`);
  }
  if (!/^https:\/\/github\.com\/[A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+$/.test(asset.repository)
    || !/^[0-9a-f]{40}$/.test(asset.revision)
    || !asset.sourcePath
    || asset.license !== "MIT"
    || !asset.copyright) {
    throw new Error(`brand asset provenance is incomplete: ${asset.path}`);
  }
  expectedBrandAssets.add(asset.path);
  const data = await readFile(resolve(brandAssetRoot, asset.path));
  const digest = createHash("sha256").update(data).digest("hex");
  if (data.length !== asset.bytes || digest !== asset.sha256 || data.toString("utf8").includes("<script")) {
    throw new Error(`brand asset integrity mismatch: ${asset.path}`);
  }
}

const actualBrandAssets = new Set(await readdir(brandAssetRoot));
if (actualBrandAssets.size !== expectedBrandAssets.size || [...actualBrandAssets].some((name) => !expectedBrandAssets.has(name))) {
  throw new Error("brand asset directory and provenance manifest differ");
}

const uiPackage = JSON.parse(await readFile(resolve(packageRoot, "package.json"), "utf8"));
const simpleIconsRoot = resolve(packageRoot, "..", "node_modules", "simple-icons");
const simpleIconsPackage = JSON.parse(await readFile(resolve(simpleIconsRoot, "package.json"), "utf8"));
if (uiPackage.dependencies?.["simple-icons"] !== "16.31.0"
  || simpleIconsPackage.version !== "16.31.0"
  || simpleIconsPackage.license !== "CC0-1.0") {
  throw new Error("Simple Icons must remain pinned to the reviewed CC0-1.0 release");
}

const simpleIconData = JSON.parse(await readFile(resolve(simpleIconsRoot, "data", "simple-icons.json"), "utf8"));
const iconLicenses = new Map(simpleIconData.map((icon) => [icon.slug, icon.license?.type ?? "CC0-1.0"]));
if (iconLicenses.get("debian") !== "CC-BY-SA-3.0" || iconLicenses.get("yarn") !== "CC-BY-4.0") {
  throw new Error("reviewed Simple Icons brand-license metadata changed");
}

const searchable = [
  await readFile(resolve(packageRoot, "NOTICE.md"), "utf8"),
  ...await Promise.all(manifest.assets.filter((asset) => asset.mediaType.includes("svg")).map((asset) => readFile(resolve(assetRoot, asset.path), "utf8"))),
].join("\n").toLowerCase();
if (searchable.includes("local.adguard.org") || searchable.includes("<script")) {
  throw new Error("executable reference content was retained");
}

console.log(`verified ${manifest.assets.length} Courier identity assets and ${brandManifest.assets.length} pinned brand assets`);
