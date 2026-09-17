import * as React from "react";
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
import powershellSource from "./brand-assets/powershell.svg";
import scoopSource from "./brand-assets/scoop.svg";
import { Download } from "lucide-react";
import { cn } from "./lib/utils";

export const brandIconNames = [
  "alpine-linux", "arch-linux", "curl", "debian", "fedora", "homebrew", "linux", "manjaro", "npm", "npx", "pnpm", "powershell", "red-hat", "scoop", "ubuntu", "wget", "yarn",
] as const;
export type BrandIconName = typeof brandIconNames[number];

const brands = {
  "alpine-linux": [alpineSource, "Alpine Linux"],
  "arch-linux": [archSource, "Arch Linux"],
  curl: [curlSource, "curl"],
  debian: [debianSource, "Debian"],
  fedora: [fedoraSource, "Fedora"],
  homebrew: [homebrewSource, "Homebrew"],
  linux: [linuxSource, "Linux"],
  manjaro: [manjaroSource, "Manjaro"],
  npm: [npmSource, "npm"],
  npx: [npmSource, "npx"],
  pnpm: [pnpmSource, "pnpm"],
  powershell: [powershellSource, "PowerShell"],
  "red-hat": [redhatSource, "Red Hat"],
  scoop: [scoopSource, "Scoop"],
  ubuntu: [ubuntuSource, "Ubuntu"],
  yarn: [yarnSource, "Yarn"],
} as const;

export function resolveBrandIcon(name: string): BrandIconName {
  return brandIconNames.includes(name as BrandIconName) ? name as BrandIconName : "linux";
}

export interface BrandIconProps extends React.HTMLAttributes<HTMLSpanElement> {
  name: string;
  label?: string;
}

export function BrandIcon({ name, label, className, ...props }: BrandIconProps): React.JSX.Element {
  const resolved = resolveBrandIcon(name);
  const title = resolved === "wget" ? "GNU Wget" : brands[resolved][1];
  const accessible = label ?? title;
  if (resolved === "wget") {
    return <span role="img" aria-label={accessible} className={cn("inline-grid size-4 place-items-center", className)} {...props}><Download className="size-full" /></span>;
  }
  const source = brands[resolved][0];
  const style = { "--courier-brand-mask": `url("${source}")` } as React.CSSProperties;
  return <span role="img" aria-label={accessible} className={cn("inline-block size-4 bg-current [mask:var(--courier-brand-mask)_center/contain_no-repeat]", className)} style={style} {...props} />;
}
