import { createHash } from "node:crypto";
import { readFile, readdir } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const packageRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const assetRoot = resolve(packageRoot, "assets");
const manifest = JSON.parse(await readFile(resolve(assetRoot, "provenance.json"), "utf8"));
const brandAssetRoot = resolve(packageRoot, "assets", "brands");
const brandManifest = JSON.parse(await readFile(resolve(brandAssetRoot, "provenance.json"), "utf8"));

if (manifest.schemaVersion !== 4 || !Array.isArray(manifest.assets) || manifest.assets.length === 0) {
  throw new Error("asset provenance manifest is invalid");
}
if (brandManifest.schemaVersion !== 2 || !Array.isArray(brandManifest.assets) || brandManifest.assets.length === 0) {
  throw new Error("brand asset provenance manifest is invalid");
}

function uint24(data, offset) {
  return data[offset] | (data[offset + 1] << 8) | (data[offset + 2] << 16);
}

function rasterDimensions(data, mediaType) {
  if (mediaType === "image/png") {
    if (data.toString("hex", 0, 8) !== "89504e470d0a1a0a") throw new Error("invalid PNG header");
    return { width: data.readUInt32BE(16), height: data.readUInt32BE(20), alpha: [4, 6].includes(data[25]) };
  }
  if (data.toString("ascii", 0, 4) !== "RIFF" || data.toString("ascii", 8, 12) !== "WEBP") throw new Error("invalid WebP header");
  const chunk = data.toString("ascii", 12, 16);
  if (chunk === "VP8X") return { width: uint24(data, 24) + 1, height: uint24(data, 27) + 1, alpha: Boolean(data[20] & 0x10) };
  if (chunk === "VP8 ") return { width: data.readUInt16LE(26) & 0x3fff, height: data.readUInt16LE(28) & 0x3fff, alpha: false };
  if (chunk === "VP8L") {
    const bits = data.readUInt32LE(21);
    return { width: (bits & 0x3fff) + 1, height: ((bits >>> 14) & 0x3fff) + 1, alpha: Boolean((bits >>> 28) & 1) };
  }
  throw new Error(`unsupported WebP chunk: ${chunk}`);
}

const expected = new Set(["provenance.json", "brands"]);
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
    if (dimensions.width !== asset.width || dimensions.height !== asset.height || (asset.transparent && !dimensions.alpha)) {
      throw new Error(`asset dimensions mismatch: ${asset.path}`);
    }
  }
}

const actual = new Set(await readdir(assetRoot));
if (actual.size !== expected.size || [...actual].some((name) => !expected.has(name))) {
  throw new Error("asset directory and provenance manifest differ");
}

const relayBudgets = new Map([
  ["courier-relay-pixel-mark-v2.webp", 24 * 1024],
  ["courier-relay-pixel-neutral-v1.webp", 64 * 1024],
  ["courier-relay-pixel-route-v2.webp", 64 * 1024],
  ["courier-relay-pixel-delivery-v2.webp", 48 * 1024],
  ["courier-relay-pixel-admin-v2.webp", 48 * 1024],
]);
const relayAssets = manifest.assets.filter((asset) => relayBudgets.has(asset.path));
const canonical = manifest.identityFamily?.canonicalAsset;
const expectedInvariants = [
  "one compact near-square pigeon body",
  "one large orange-ringed pigeon eye and short ivory beak",
  "blue-gray head, pale folded wings, compact dark tail and coral feet",
  "one small left-side earpiece and one cobalt courier satchel",
];
const expectedVariations = ["pose", "role equipment", "carried or attached object"];
const expectedForbidden = ["body proportions", "physiology", "base plumage", "eye geometry", "armor", "helmet", "visor", "police or military styling"];
if (relayAssets.length !== relayBudgets.size
  || canonical !== "courier-relay-pixel-neutral-v1.webp"
  || manifest.identityFamily?.revision !== "relay-pixel-v2"
  || manifest.identityFamily?.qaContactSheet !== "docs/assets/courier-relay-pixel-v2-contact-sheet.png"
  || !/identity-preserving .*derivative/.test(manifest.identityFamily?.generationRule ?? "")
  || manifest.identityFamily?.invariants?.join("\n") !== expectedInvariants.join("\n")
  || manifest.identityFamily?.allowedVariations?.join("\n") !== expectedVariations.join("\n")
  || manifest.identityFamily?.forbiddenVariations?.join("\n") !== expectedForbidden.join("\n")
  || relayAssets.some((asset) => asset.mediaType !== "image/webp"
    || asset.width !== asset.height
    || !asset.transparent
    || asset.bytes > relayBudgets.get(asset.path)
    || !asset.prompt
    || !asset.lineage
    || !Array.isArray(asset.consumers)
    || asset.consumers.length === 0
    || (asset.path === canonical ? asset.identityRevision !== "relay-pixel-v1" || asset.identityReference !== "self" : asset.identityRevision !== manifest.identityFamily.revision || asset.identityReference !== canonical))
  || relayAssets.reduce((total, asset) => total + asset.bytes, 0) > 248 * 1024) {
  throw new Error("the consistent transparent Relay pixel family is missing or exceeds its budgets");
}
const contactSheet = await readFile(resolve(packageRoot, "..", "..", manifest.identityFamily.qaContactSheet));
const contactSheetDimensions = rasterDimensions(contactSheet, "image/png");
if (contactSheetDimensions.width !== 1320 || contactSheetDimensions.height !== 850) {
  throw new Error("the Relay v2 QA contact sheet is missing or has unexpected dimensions");
}

const font = manifest.assets.find((asset) => asset.path === "pixelify-sans-v1.woff2");
const fontLicense = manifest.assets.find((asset) => asset.path === "pixelify-sans-ofl-1.1.txt");
if (!font || !fontLicense
  || font.mediaType !== "font/woff2"
  || font.bytes > 48 * 1024
  || font.revision !== "39df74aba80df8157546034b878e8be1eb565ced"
  || font.repository !== "https://github.com/eifetx/Pixelify-Sans"
  || font.license !== "OFL-1.1"
  || fontLicense.revision !== font.revision
  || !(await readFile(resolve(assetRoot, fontLicense.path), "utf8")).includes("SIL OPEN FONT LICENSE Version 1.1")) {
  throw new Error("the pinned local Pixelify Sans font or OFL notice is invalid");
}

const expectedBrandAssets = new Set(["provenance.json"]);
for (const asset of brandManifest.assets) {
  if (!/^[a-z0-9][a-z0-9.-]*\.png$/.test(asset.path) || expectedBrandAssets.has(asset.path)) {
    throw new Error(`invalid or duplicate brand asset path: ${asset.path}`);
  }
  if (!/^https:\/\//.test(asset.source)
    || !asset.revision
    || !asset.color
    || !asset.license
    || asset.mediaType !== "image/png"
    || asset.width !== 128
    || asset.height !== 128
    || !asset.transparent) {
    throw new Error(`brand asset provenance is incomplete: ${asset.path}`);
  }
  expectedBrandAssets.add(asset.path);
  const data = await readFile(resolve(brandAssetRoot, asset.path));
  const digest = createHash("sha256").update(data).digest("hex");
  const dimensions = rasterDimensions(data, asset.mediaType);
  if (data.length !== asset.bytes || digest !== asset.sha256 || dimensions.width !== 128 || dimensions.height !== 128 || !dimensions.alpha) {
    throw new Error(`brand asset integrity mismatch: ${asset.path}`);
  }
}

const brandNames = new Set(brandManifest.assets.map((asset) => asset.path));
if (brandNames.has("npx.png")
  || brandNames.has("wget.png")
  || !brandNames.has("github-light.png")
  || !brandNames.has("github-dark.png")
  || brandManifest.assets.find((asset) => asset.path === "github-light.png")?.color !== "#181717"
  || brandManifest.assets.find((asset) => asset.path === "github-dark.png")?.color !== "#FFFFFF") {
  throw new Error("text-only channels or GitHub theme variants violate the brand policy");
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

const sourceRoots = [packageRoot, resolve(packageRoot, "..", "landing"), resolve(packageRoot, "..", "data"), resolve(packageRoot, "..", "admin")].map((root) => resolve(root, "src"));
const sourceEntries = await Promise.all(sourceRoots.map(async (root) => ({ root, files: (await readdir(root, { recursive: true })).filter((name) => /\.(?:ts|tsx|css)$/.test(name)) })));
const browserSource = (await Promise.all(sourceEntries.flatMap(({ root, files }) => files.map((name) => readFile(resolve(root, name), "utf8"))))).join("\n");
if (browserSource.includes("lucide-react")
  || browserSource.includes("🇬🇧")
  || browserSource.includes("🇷🇺")
  || browserSource.includes("courier-relay-tech")
  || browserSource.includes("courier-relay-mark-v2")) {
  throw new Error("legacy icons, emoji flags, or Relay assets remain in browser source");
}

const brandComponent = await readFile(resolve(packageRoot, "src", "brand-icons-react.tsx"), "utf8");
if (!brandComponent.includes("<img") || brandComponent.includes("<svg") || brandComponent.includes("courier-pixel-image") || brandComponent.includes("pixelated")) {
  throw new Error("third-party brands must render only as normal local raster images");
}

console.log(`verified ${relayAssets.length} consistent Relay sprites, the pinned display font, and ${brandManifest.assets.length} local raster brand assets`);
