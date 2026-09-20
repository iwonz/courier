import * as React from "react";
import { GithubBrandIcon } from "./brand-icons-react";

export const pixelIconNames = [
  "activity", "archive", "arrow-up", "browser-share", "browser-upload", "check", "chevron-down", "chevron-up", "copy", "database", "download", "file", "folder", "folder-in", "folder-out", "moon", "package", "parcel", "policy", "refresh", "retry", "route", "save", "server", "server-in", "server-out", "shield", "stop", "sun", "system", "terminal", "upload", "webhook-in", "webhook-out", "workflow",
] as const;

export type PixelIconName = typeof pixelIconNames[number];
export type IconName = PixelIconName;

interface PixelGlyph { readonly path: string; readonly viewBox?: string; }

const arrowUp = "M7 1h2v2h2v2h2v2H9v6H7V7H3V5h2V3h2z";
const arrowDown = "M7 15h2v-2h2v-2h2V9H9V3H7v6H3v2h2v2h2z";
const folder = "M1 4h5V2h3v2h6v10H1zm2 2v6h10V6z";
const server = "M2 1h12v6H2zm2 2v2h2V3zm-2 5h12v7H2zm2 2v3h2v-3z";
const webhook = "M6 1h4v2H8v2h2v2H6V5H4V3h2zm4 6h4v2h-2v2h-2v2H8V9h2zM2 8h4v2H4v2h4v2H4v-1H2z";

const glyphs: Record<PixelIconName, PixelGlyph> = {
  activity: { path: "M1 8h3l2-5h2l2 9 2-4h3v2h-2l-2 5H9L7 7l-1 3H1z" },
  archive: { path: "M2 1h12v4h-1v10H3V5H2zm2 2v1h8V3zm1 3v7h6V6zm2 1h2v2H7z" },
  "arrow-up": { path: arrowUp },
  "browser-share": { path: "M1 2h14v12H1zm2 3v7h10V5zm4 1h2v2h2v2H9v1H7v-1H5V8h2z" },
  "browser-upload": { path: `M1 2h14v12H1zm2 3v7h10V5z${arrowUp}` },
  check: { path: "M2 8h2v2h2v2h2v-2h2V8h2V6h2V4h-2v1h-2v2H8v2H6V7H4V6H2z" },
  "chevron-down": { path: "M2 5h2v2h2v2h4V7h2V5h2v4h-2v2h-2v2H6v-2H4V9H2z" },
  "chevron-up": { path: "M6 3h4v2h2v2h2v4h-2V9h-2V7H6v2H4v2H2V7h2V5h2z" },
  copy: { path: "M5 1h10v10h-3V8h1V3H7v1H5zm-4 3h10v11H1zm2 2v7h6V6z" },
  database: { path: "M3 1h10v2h2v10h-2v2H3v-2H1V3h2zm0 2v2h10V3zm0 4v2h10V7zm0 4v2h10v-2z" },
  download: { path: arrowDown },
  file: { path: "M3 1h7l3 3v11H3zm2 2v10h6V6H8V3zm5 0v2h2z" },
  folder: { path: folder },
  "folder-in": { path: `${folder}M7 6h2v3h2v2H9v2H7v-2H5V9h2z` },
  "folder-out": { path: `${folder}M7 12h2V9h2V7H9V5H7v2H5v2h2z` },
  moon: { path: "M5 1h5v2H8v2H6v4h2v2h2v2H5v-2H3V9H1V5h2V3h2z" },
  package: { path: "M2 4h12v10H2zm2 2v6h8V6zM5 1h6v2H5z" },
  parcel: { path: "M1 4l7-3 7 3v9l-7 3-7-3zm3 1l4 2 4-2-4-2zm-1 2v5l4 2V9zm10 0L9 9v5l4-2z" },
  policy: { path: "M3 2h4V1h2v1h4v13H3zm2 2v9h6V4zm1 2h4v2H6zm0 4h3v2H6z" },
  refresh: { path: "M4 2h7V1h2v2h2v4H9V5h2V4H5v2H3v4H1V6h1V4h2zm1 10h6v-2h2V6h2v4h-1v2h-2v2H5v1H3v-2H1V9h6v2H5z" },
  retry: { path: "M4 2h7V1h2v2h2v4H9V5h2V4H5v2H3v4H1V6h1V4h2zm1 10h6v-2h2V6h2v4h-1v2h-2v2H5v1H3v-2H1V9h6v2H5z" },
  route: { path: "M1 2h4v4H1zm10 8h4v4h-4zM5 3h3v2H5zm3 2h2v2H8zm2 2h2v3h-2z" },
  save: { path: "M2 1h10l2 2v12H2zm2 2v4h7V3zm1 6v4h6V9zM8 3h2v3H8z" },
  server: { path: server },
  "server-in": { path: `${server}M7 2h2v2h2v2H9v1H7V6H5V4h2z` },
  "server-out": { path: `${server}M7 14h2v-2h2v-2H9V9H7v1H5v2h2z` },
  shield: { path: "M2 2h12v7h-2v3h-2v2H6v-2H4V9H2zm2 2v5h2v2h4V9h2V4z" },
  stop: { path: "M2 2h12v12H2zm3 3v6h6V5z" },
  sun: { path: "M7 1h2v3H7zM2 3h2v2H2zm10 0h2v2h-2zM5 5h6v6H5zm-4 2h3v2H1zm11 0h3v2h-3zM2 11h2v2H2zm10 0h2v2h-2zM7 12h2v3H7z" },
  system: { path: "M1 2h14v10H9v1h3v2H4v-2h3v-1H1zm2 2v6h10V4z" },
  terminal: { path: "M1 2h14v12H1zm2 3v2h2v2h2V7H5V5zm5 5h5v2H8z" },
  upload: { path: arrowUp },
  "webhook-in": { path: `${webhook}M6 6h2v2h2v2H8v1H6v-1H4V8h2z` },
  "webhook-out": { path: `${webhook}M8 6h2v2h2v2h-2v1H8v-1H6v-2h2z` },
  workflow: { path: "M1 1h5v5H1zm9 0h5v5h-5zM5 3h6v2H5zm2 2h2v6H7zm-2 6h6v4H5z" },
};

export function resolvePixelIcon(name: string): PixelIconName {
  return pixelIconNames.includes(name as PixelIconName) ? name as PixelIconName : "parcel";
}

export interface PixelIconProps extends React.ComponentPropsWithoutRef<"svg"> { readonly name: string; readonly label?: string; }

export function PixelIcon({ name, label, ...props }: PixelIconProps): React.JSX.Element {
  const glyph = glyphs[resolvePixelIcon(name)];
  return <svg viewBox={glyph.viewBox ?? "0 0 16 16"} fill="currentColor" shapeRendering="crispEdges" role={label ? "img" : undefined} aria-hidden={label ? undefined : true} aria-label={label} {...props}><path d={glyph.path} /></svg>;
}

export const iconNames = pixelIconNames;
export const resolveIcon = resolvePixelIcon;
export const Icon = PixelIcon;

export function Github(props: React.HTMLAttributes<HTMLSpanElement>): React.JSX.Element { return <GithubBrandIcon {...props} />; }
export function Copy(props: Omit<PixelIconProps, "name">): React.JSX.Element { return <PixelIcon name="copy" {...props} />; }
export function Check(props: Omit<PixelIconProps, "name">): React.JSX.Element { return <PixelIcon name="check" {...props} />; }
export function Route(props: Omit<PixelIconProps, "name">): React.JSX.Element { return <PixelIcon name="route" {...props} />; }
export function ServerCog(props: Omit<PixelIconProps, "name">): React.JSX.Element { return <PixelIcon name="server" {...props} />; }
