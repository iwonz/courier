import * as React from "react";
import alpineLinuxSource from "../assets/brands/alpine-linux.png";
import appleSource from "../assets/brands/apple.png";
import archLinuxSource from "../assets/brands/arch-linux.png";
import curlSource from "../assets/brands/curl.png";
import debianSource from "../assets/brands/debian.png";
import fedoraSource from "../assets/brands/fedora.png";
import githubDarkSource from "../assets/brands/github-dark.png";
import githubLightSource from "../assets/brands/github-light.png";
import homebrewSource from "../assets/brands/homebrew.png";
import linuxSource from "../assets/brands/linux.png";
import manjaroSource from "../assets/brands/manjaro.png";
import npmSource from "../assets/brands/npm.png";
import pnpmSource from "../assets/brands/pnpm.png";
import powershellSource from "../assets/brands/powershell.png";
import redHatSource from "../assets/brands/red-hat.png";
import scoopSource from "../assets/brands/scoop.png";
import ubuntuSource from "../assets/brands/ubuntu.png";
import windowsSource from "../assets/brands/windows.png";
import yarnSource from "../assets/brands/yarn.png";
import { cn } from "./lib/utils";

export const brandIconNames = [
  "alpine-linux", "apple", "arch-linux", "curl", "debian", "fedora", "homebrew", "linux", "manjaro", "npm", "pnpm", "powershell", "red-hat", "scoop", "ubuntu", "windows", "yarn",
] as const;

export type BrandIconName = typeof brandIconNames[number];

const assets: Readonly<Record<BrandIconName, { readonly label: string; readonly source: string }>> = {
  "alpine-linux": { label: "Alpine Linux", source: alpineLinuxSource },
  apple: { label: "Apple", source: appleSource },
  "arch-linux": { label: "Arch Linux", source: archLinuxSource },
  curl: { label: "curl", source: curlSource },
  debian: { label: "Debian", source: debianSource },
  fedora: { label: "Fedora", source: fedoraSource },
  homebrew: { label: "Homebrew", source: homebrewSource },
  linux: { label: "Linux", source: linuxSource },
  manjaro: { label: "Manjaro", source: manjaroSource },
  npm: { label: "npm", source: npmSource },
  pnpm: { label: "pnpm", source: pnpmSource },
  powershell: { label: "PowerShell", source: powershellSource },
  "red-hat": { label: "Red Hat", source: redHatSource },
  scoop: { label: "Scoop", source: scoopSource },
  ubuntu: { label: "Ubuntu", source: ubuntuSource },
  windows: { label: "Windows", source: windowsSource },
  yarn: { label: "Yarn", source: yarnSource },
};

export function resolveBrandIcon(name: string): BrandIconName {
  return brandIconNames.includes(name as BrandIconName) ? name as BrandIconName : "linux";
}

export interface BrandIconProps extends Omit<React.ComponentPropsWithoutRef<"img">, "src" | "width" | "height"> { readonly name: string; readonly label?: string; }

export function BrandIcon({ name, label, className, alt, ...props }: BrandIconProps): React.JSX.Element {
  const asset = assets[resolveBrandIcon(name)];
  return <img src={asset.source} width="128" height="128" alt={alt ?? label ?? ""} role={label ? "img" : undefined} aria-label={label} aria-hidden={label || alt ? undefined : true} data-brand-name={asset.label} className={cn("size-4 object-contain", className)} {...props} />;
}

export function GithubBrandIcon({ className, ...props }: React.HTMLAttributes<HTMLSpanElement>): React.JSX.Element {
  return <span aria-hidden="true" className={cn("inline-grid size-4 place-items-center", className)} {...props}>
    <img src={githubLightSource} width="128" height="128" alt="" className="size-4 object-contain dark:hidden" />
    <img src={githubDarkSource} width="128" height="128" alt="" className="hidden size-4 object-contain dark:block" />
  </span>;
}

export const resolvePixelBrandIcon = resolveBrandIcon;
export const PixelBrandIcon = BrandIcon;
export const pixelBrandIconNames = brandIconNames;
export type PixelBrandIconName = BrandIconName;
