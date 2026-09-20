import { createHash } from "node:crypto";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { chromium } from "@playwright/test";
import {
  siAlpinelinux,
  siApple,
  siArchlinux,
  siCurl,
  siDebian,
  siFedora,
  siGithub,
  siHomebrew,
  siLinux,
  siManjaro,
  siNpm,
  siPnpm,
  siRedhat,
  siUbuntu,
  siYarn,
} from "simple-icons";

const packageRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const outputRoot = resolve(packageRoot, "assets", "brands");
const sourceRoot = resolve(packageRoot, "scripts", "brand-sources");
await mkdir(outputRoot, { recursive: true });

const simpleIcons = [
  ["alpine-linux", siAlpinelinux],
  ["apple", siApple],
  ["arch-linux", siArchlinux],
  ["curl", siCurl],
  ["debian", siDebian],
  ["fedora", siFedora],
  ["homebrew", siHomebrew],
  ["linux", siLinux],
  ["manjaro", siManjaro],
  ["npm", siNpm],
  ["pnpm", siPnpm],
  ["red-hat", siRedhat],
  ["ubuntu", siUbuntu],
  ["yarn", siYarn],
];

const entries = simpleIcons.map(([name, icon]) => ({
  name,
  label: icon.title,
  source: icon.source,
  revision: "simple-icons@16.31.0",
  color: `#${icon.hex}`,
  svg: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path fill="#${icon.hex}" d="${icon.path}"/></svg>`,
}));
entries.push(
  { name: "github-light", label: "GitHub", source: siGithub.source, revision: "simple-icons@16.31.0", color: "#181717", svg: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path fill="#181717" d="${siGithub.path}"/></svg>` },
  { name: "github-dark", label: "GitHub", source: siGithub.source, revision: "simple-icons@16.31.0", color: "#FFFFFF", svg: `<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path fill="#FFFFFF" d="${siGithub.path}"/></svg>` },
  { name: "windows", label: "Windows", source: "https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks", revision: "windows-11-four-pane", color: "#0078D4", svg: '<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path fill="#0078D4" d="M1 1h10v10H1zm12 0h10v10H13zM1 13h10v10H1zm12 0h10v10H13z"/></svg>' },
  { name: "powershell", label: "PowerShell", source: "https://github.com/PowerShell/PowerShell", revision: "5e35e5ac892815ec4991f74425989e38ce71ccd0", color: "official-multicolor", svg: await readFile(resolve(sourceRoot, "powershell.svg"), "utf8") },
  { name: "scoop", label: "Scoop", source: "https://github.com/ScoopInstaller/scoopinstaller.github.io", revision: "9834decbc5acc60b8e1548d518b40d7774658127", color: "official-multicolor", svg: await readFile(resolve(sourceRoot, "scoop.svg"), "utf8") },
);

const browser = await chromium.launch({ headless: true });
const page = await browser.newPage({ viewport: { width: 128, height: 128 }, deviceScaleFactor: 1 });
const metadata = [];
try {
  for (const entry of entries) {
    const svg = entry.svg.replace(/<\?xml[^>]*>/g, "").replace(/<!--.*?-->/gs, "");
    await page.setContent(`<style>html,body{margin:0;width:128px;height:128px;background:transparent}#icon{display:grid;width:128px;height:128px;place-items:center}#icon>svg{display:block;width:112px!important;height:112px!important;max-width:112px;max-height:112px}</style><div id="icon">${svg}</div>`);
    const path = resolve(outputRoot, `${entry.name}.png`);
    const data = await page.locator("#icon").screenshot({ path, omitBackground: true, animations: "disabled" });
    await writeFile(path, data);
    metadata.push({
      path: `${entry.name}.png`,
      label: entry.label,
      mediaType: "image/png",
      width: 128,
      height: 128,
      transparent: true,
      bytes: data.length,
      sha256: createHash("sha256").update(data).digest("hex"),
      source: entry.source,
      revision: entry.revision,
      color: entry.color,
    });
  }
} finally {
  await browser.close();
}

console.log(JSON.stringify(metadata, null, 2));
