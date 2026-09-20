import { mkdir, readFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { chromium } from "@playwright/test";

const packageRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const repositoryRoot = resolve(packageRoot, "..", "..");
const output = resolve(repositoryRoot, "docs", "assets", "courier-relay-pixel-v2-contact-sheet.png");
const assets = [
  ["neutral-v1 reference", "courier-relay-pixel-neutral-v1.webp"],
  ["mark-v2", "courier-relay-pixel-mark-v2.webp"],
  ["route-v2", "courier-relay-pixel-route-v2.webp"],
  ["delivery-v2", "courier-relay-pixel-delivery-v2.webp"],
  ["admin-v2", "courier-relay-pixel-admin-v2.webp"],
];
const images = await Promise.all(assets.map(async ([label, name]) => [label, `data:image/webp;base64,${(await readFile(resolve(packageRoot, "assets", name))).toString("base64")}`]));
const cells = images.map(([label, source]) => `<figure><img src="${source}" alt=""><figcaption>${label}</figcaption></figure>`).join("");
const browser = await chromium.launch({ headless: true });
const page = await browser.newPage({ viewport: { width: 1320, height: 850 }, deviceScaleFactor: 1 });
try {
  await page.setContent(`<style>
    *{box-sizing:border-box}html,body{margin:0;width:1320px;height:850px;font-family:ui-monospace,SFMono-Regular,Menlo,monospace;background:#111827;color:#f8fafc}
    h1{height:70px;margin:0;padding:24px 30px;font-size:20px;letter-spacing:.04em}
    section{display:grid;grid-template-columns:repeat(5,1fr);height:250px;padding:12px 20px;gap:10px}
    section.light{background:#f8fafc;color:#111827}section.dark{background:#111827;color:#f8fafc}
    section.checker{color:#111827;background-color:#fff;background-image:linear-gradient(45deg,#d7dce3 25%,transparent 25%),linear-gradient(-45deg,#d7dce3 25%,transparent 25%),linear-gradient(45deg,transparent 75%,#d7dce3 75%),linear-gradient(-45deg,transparent 75%,#d7dce3 75%);background-size:24px 24px;background-position:0 0,0 12px,12px -12px,-12px 0}
    figure{margin:0;display:grid;grid-template-rows:1fr 24px;place-items:center;min-width:0}img{width:190px;height:190px;object-fit:contain;image-rendering:pixelated}figcaption{font-size:12px;font-weight:700;background:color-mix(in srgb,currentColor 12%,transparent);padding:4px 8px}
  </style><h1>Relay identity QA — canonical neutral-v1 and v2 derivatives</h1><section class="light">${cells}</section><section class="dark">${cells}</section><section class="checker">${cells}</section>`);
  await mkdir(dirname(output), { recursive: true });
  await page.screenshot({ path: output, fullPage: true });
} finally {
  await browser.close();
}
console.log(output);
