import * as React from "react";
import { cn } from "./lib/utils";

export const pixelBrandIconNames = [
  "alpine-linux", "arch-linux", "curl", "debian", "fedora", "homebrew", "linux", "manjaro", "npm", "npx", "pnpm", "powershell", "red-hat", "scoop", "ubuntu", "wget", "yarn",
] as const;

export type PixelBrandIconName = typeof pixelBrandIconNames[number];
export type BrandIconName = PixelBrandIconName;

interface PixelBrandGlyph { readonly label: string; readonly path: string; }

const download = "M7 1h2v7h2V6h2v4h-2v2H9v2H7v-2H5v-2H3V6h2v2h2zM2 14h12v2H2z";
const glyphs: Record<PixelBrandIconName, PixelBrandGlyph> = {
  "alpine-linux": { label: "Alpine Linux", path: "M1 14L7 2h2l2 4 1-2h1l2 10h-2l-1-6-1 2-3-5-5 9zm4-2h6L8 6z" },
  "arch-linux": { label: "Arch Linux", path: "M1 15L7 1h2l6 14h-3L8 7l-2 5h4l1 3z" },
  curl: { label: "curl", path: "M4 3h8v2H6v2H4v4h2v2h6v2H4v-2H2V5h2zM9 7h5v2H9z" },
  debian: { label: "Debian", path: "M7 1h4v2H7V2H5v2H3v2H1v5h2v2h2v2h5v-2H6v-2H4V7h2V5h5v2H8v2h4V7h2V4h-2V2h-1V1z" },
  fedora: { label: "Fedora", path: "M4 2h8v2H8v2h4v3H8v5H4V9h2V4H4zm4 6h2V7H8z" },
  homebrew: { label: "Homebrew", path: "M3 4h8v2h3v5h-3v3H4V6H3zm2 2v6h4V6zm6 2v1h1V8zM5 1h2v2H5zm4 0h2v2H9z" },
  linux: { label: "Linux", path: "M6 1h4v2h2v3h2v7h-3v2H5v-2H2V6h2V3h2zm0 3v3h4V4zM4 9v4h2V9zm6 0v4h2V9z" },
  manjaro: { label: "Manjaro", path: "M1 1h14v14H1zm3 3v8h2V4zm4 0v3h4V4zm0 5v3h4V9z" },
  npm: { label: "npm", path: "M1 4h14v8h-2V6H3v4h2V7h2v5H1zm8 3h3v3h-1V8H9z" },
  npx: { label: "npx", path: "M1 4h14v8h-2V6H3v4h2V7h2v5H1zm8 3h3v3h-1V8H9z" },
  pnpm: { label: "pnpm", path: "M1 1h4v4H1zm5 0h4v4H6zm5 0h4v4h-4zM1 6h4v4H1zm5 0h4v4H6zm5 0h4v4h-4zM1 11h4v4H1zm5 0h4v4H6zm5 0h4v4h-4z" },
  powershell: { label: "PowerShell", path: "M2 2h12v12H2zm3 3v2h2v2h2v2H7V9H5V7H3V5zm4 6h3v2H9z" },
  "red-hat": { label: "Red Hat", path: "M5 2h6v2h2v2h2v4H1V6h2V4h2zm-2 9h10v2h-2v2H5v-2H3z" },
  scoop: { label: "Scoop", path: "M4 1h8v4h2v10H2V5h2zm2 2v2h4V3zM4 7v6h8V7zm2 1h4v2H6z" },
  ubuntu: { label: "Ubuntu", path: "M6 2h4v2H6zM2 5h3v3H2zm9 0h3v3h-3zM5 5h2v2H5zm4 0h2v2H9zM4 8h2v3h4V8h2v4h-2v2H6v-2H4z" },
  wget: { label: "GNU Wget", path: download },
  yarn: { label: "Yarn", path: "M6 1h4v2h2v2h2v6h-2v2h-2v2H6v-2H4v-2H2V5h2V3h2zm0 3v2H4v4h2v2h4v-2h2V6h-2V4z" },
};

export function resolvePixelBrandIcon(name: string): PixelBrandIconName {
  return pixelBrandIconNames.includes(name as PixelBrandIconName) ? name as PixelBrandIconName : "linux";
}

export interface PixelBrandIconProps extends React.ComponentPropsWithoutRef<"svg"> { readonly name: string; readonly label?: string; }

export function PixelBrandIcon({ name, label, className, ...props }: PixelBrandIconProps): React.JSX.Element {
  const glyph = glyphs[resolvePixelBrandIcon(name)];
  return <svg viewBox="0 0 16 16" fill="currentColor" shapeRendering="crispEdges" role={label ? "img" : undefined} aria-label={label} aria-hidden={label ? undefined : true} data-brand-name={glyph.label} className={cn("size-4", className)} {...props}><path d={glyph.path} /></svg>;
}

export const brandIconNames = pixelBrandIconNames;
export const resolveBrandIcon = resolvePixelBrandIcon;
export const BrandIcon = PixelBrandIcon;
