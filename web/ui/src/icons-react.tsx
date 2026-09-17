import * as React from "react";
import {
  Apple,
  Archive,
  Box,
  Check,
  Copy,
  Download,
  Folder,
  FolderInput,
  FolderOutput,
  Globe2,
  HardDriveDownload,
  HardDriveUpload,
  Monitor,
  Moon,
  Package,
  PanelsTopLeft,
  ReceiptText,
  RefreshCw,
  Route,
  Server,
  ServerCog,
  ShieldCheck,
  Sun,
  Terminal,
  Upload,
  Webhook,
} from "lucide-react";

export const iconNames = [
  "apple", "archive", "browser-share", "browser-upload", "check", "copy", "download", "folder", "folder-in", "folder-out", "github", "moon", "package", "parcel", "receipt", "retry", "route", "server", "server-in", "server-out", "shield", "sun", "system", "terminal", "upload", "webhook-in", "webhook-out", "windows",
] as const;
export type IconName = typeof iconNames[number];

const icons: Record<IconName, React.ElementType> = {
  apple: Apple,
  archive: Archive,
  "browser-share": Globe2,
  "browser-upload": Upload,
  check: Check,
  copy: Copy,
  download: Download,
  folder: Folder,
  "folder-in": FolderInput,
  "folder-out": FolderOutput,
  github: GitHubMark,
  moon: Moon,
  package: Package,
  parcel: Box,
  receipt: ReceiptText,
  retry: RefreshCw,
  route: Route,
  server: Server,
  "server-in": HardDriveDownload,
  "server-out": HardDriveUpload,
  shield: ShieldCheck,
  sun: Sun,
  system: Monitor,
  terminal: Terminal,
  upload: Upload,
  "webhook-in": Webhook,
  "webhook-out": Webhook,
  windows: PanelsTopLeft,
};

export function resolveIcon(name: string): IconName {
  return iconNames.includes(name as IconName) ? name as IconName : "parcel";
}

export interface IconProps extends React.ComponentPropsWithoutRef<"svg"> {
  name: string;
  label?: string;
}

export function Icon({ name, label, ...props }: IconProps): React.JSX.Element {
  const Component = icons[resolveIcon(name)];
  return <Component role={label ? "img" : undefined} aria-hidden={label ? undefined : true} aria-label={label} {...props} />;
}

export function GitHubMark(props: React.ComponentPropsWithoutRef<"svg">): React.JSX.Element {
  return <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true" {...props}><path d="M9 19c-5 1-5-2-7-3m14 6v-3.6c0-1 .1-1.7-.4-2.2 3.2-.4 6.4-1.6 6.4-7.1 0-1.6-.6-3-1.7-4 .2-.5.7-2.3-.2-4.6 0 0-1.4-.5-4.7 1.7a16 16 0 0 0-8.6 0C6.4 1 5 1.5 5 1.5 4.1 3.8 4.6 5.6 4.8 6.1a7 7 0 0 0-1.7 4c0 5.5 3.2 6.7 6.4 7.1-.4.4-.8 1.1-.8 2.2V23" /></svg>;
}

export { GitHubMark as Github, Monitor, Moon, Sun, Copy, Check, Route, ServerCog };
