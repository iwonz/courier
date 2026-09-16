import { LitElement, css, html } from "lit";
import alpineSource from "simple-icons/icons/alpinelinux.svg";
import archSource from "simple-icons/icons/archlinux.svg";
import curlSource from "simple-icons/icons/curl.svg";
import debianSource from "simple-icons/icons/debian.svg";
import fedoraSource from "simple-icons/icons/fedora.svg";
import homebrewSource from "simple-icons/icons/homebrew.svg";
import linuxSource from "simple-icons/icons/linux.svg";
import manjaroSource from "simple-icons/icons/manjaro.svg";
import npmSource from "simple-icons/icons/npm.svg";
import pnpmSource from "simple-icons/icons/pnpm.svg";
import redhatSource from "simple-icons/icons/redhat.svg";
import ubuntuSource from "simple-icons/icons/ubuntu.svg";
import yarnSource from "simple-icons/icons/yarn.svg";
import powershellSource from "../brand-assets/powershell.svg";
import scoopSource from "../brand-assets/scoop.svg";

export const brandIconNames = [
  "alpine-linux", "arch-linux", "curl", "debian", "fedora", "homebrew", "linux", "manjaro", "npm", "npx", "pnpm", "powershell", "red-hat", "scoop", "ubuntu", "wget", "yarn",
] as const;
export type BrandIconName = typeof brandIconNames[number];

interface MaskBrand { readonly kind: "mask"; readonly source: string; readonly title: string }
interface GlyphBrand { readonly kind: "glyph"; readonly title: string }
type Brand = MaskBrand | GlyphBrand;

const brands: Record<BrandIconName, Brand> = {
  "alpine-linux": { kind: "mask", source: alpineSource, title: "Alpine Linux" },
  "arch-linux": { kind: "mask", source: archSource, title: "Arch Linux" },
  curl: { kind: "mask", source: curlSource, title: "curl" },
  debian: { kind: "mask", source: debianSource, title: "Debian" },
  fedora: { kind: "mask", source: fedoraSource, title: "Fedora" },
  homebrew: { kind: "mask", source: homebrewSource, title: "Homebrew" },
  linux: { kind: "mask", source: linuxSource, title: "Linux" },
  manjaro: { kind: "mask", source: manjaroSource, title: "Manjaro" },
  npm: { kind: "mask", source: npmSource, title: "npm" },
  npx: { kind: "mask", source: npmSource, title: "npx" },
  pnpm: { kind: "mask", source: pnpmSource, title: "pnpm" },
  powershell: { kind: "mask", source: powershellSource, title: "PowerShell" },
  "red-hat": { kind: "mask", source: redhatSource, title: "Red Hat" },
  scoop: { kind: "mask", source: scoopSource, title: "Scoop" },
  ubuntu: { kind: "mask", source: ubuntuSource, title: "Ubuntu" },
  wget: { kind: "glyph", title: "GNU Wget" },
  yarn: { kind: "mask", source: yarnSource, title: "Yarn" },
};

export function resolveBrandIcon(name: string): BrandIconName {
  return brandIconNames.includes(name as BrandIconName) ? name as BrandIconName : "linux";
}

export class CourierBrandIcon extends LitElement {
  static properties = {
    name: { type: String },
    label: { type: String },
  };

  static styles = css`
    :host { display: inline-flex; width: 1.5rem; height: 1.5rem; color: currentColor; }
    svg, .mask { display: block; width: 100%; height: 100%; }
    .mask { background: currentColor; mask: var(--brand-mask) center / contain no-repeat; -webkit-mask: var(--brand-mask) center / contain no-repeat; }
  `;

  name = "linux";
  label = "";

  protected render() {
    const brand = brands[resolveBrandIcon(this.name)];
    const label = this.label || undefined;
    if (brand.kind === "glyph") {
      return html`<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linecap="square" stroke-linejoin="miter" role=${label ? "img" : "presentation"} aria-label=${label} aria-hidden=${label ? "false" : "true"}><path d="M12 3v12m-4-4 4 4 4-4M4 17v4h16v-4"></path></svg>`;
    }
    return html`<span class="mask" style=${`--brand-mask: url("${brand.source}")`} role=${label ? "img" : "presentation"} aria-label=${label} aria-hidden=${label ? "false" : "true"}></span>`;
  }
}
