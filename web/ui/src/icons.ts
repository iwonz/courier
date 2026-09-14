import { LitElement, css, html } from "lit";

export const iconNames = [
  "apple", "archive", "copy", "download", "folder", "github", "homebrew", "language-en", "language-ru", "linux", "moon", "npm", "package", "parcel", "pnpm", "receipt", "retry", "route", "scoop", "server", "shield", "sun", "system", "terminal", "upload", "windows", "yarn",
] as const;
export type IconName = typeof iconNames[number];

const paths: Record<IconName, string> = {
  apple: "M15 5c1-1 1-3 1-3-2 0-3 1-4 3m6 7c-1-2-2-3-4-3-1 0-2 1-3 1s-2-1-3-1c-3 0-5 3-5 6 0 4 3 8 5 8 1 0 2-1 3-1s2 1 3 1c2 0 4-3 5-6-2-1-3-2-3-4 0-2 1-3 2-4z",
  archive: "M3 3h18v5H3zM5 8v13h14V8M9 12h6",
  copy: "M8 3h13v13M3 8h13v13H3z",
  download: "M12 3v13m-5-5 5 5 5-5M4 15v6h16v-6",
  folder: "M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3",
  github: "M9 19c-5 1-5-2-7-3m14 6v-3.6c0-1 .1-1.7-.4-2.2 3.2-.4 6.4-1.6 6.4-7.1 0-1.6-.6-3-1.7-4 .2-.5.7-2.3-.2-4.6 0 0-1.4-.5-4.7 1.7a16 16 0 0 0-8.6 0C6.4 1 5 1.5 5 1.5 4.1 3.8 4.6 5.6 4.8 6.1a7 7 0 0 0-1.7 4c0 5.5 3.2 6.7 6.4 7.1-.4.4-.8 1.1-.8 2.2V23",
  homebrew: "M6 4h11l-1 15H8zM17 7h2a2 2 0 0 1 2 2v3a2 2 0 0 1-2 2h-2M5 22h13",
  "language-en": "M3 4h18v14H9l-4 3v-3H3z",
  "language-ru": "M3 4h18v14H9l-4 3v-3H3z",
  linux: "M12 2c-3 0-4 3-4 6-2 2-3 5-3 8l3-1 1 5 3-2 3 2 1-5 3 1c0-3-1-6-3-8 0-3-1-6-4-6zM9 8h.01M15 8h.01M10 11h4",
  moon: "M20 16a8 8 0 0 1-12-10 8 8 0 1 0 12 10z",
  npm: "M2 6h20v12H2zM6 15V9h5v6m0-6h4v6m0-6h3v6",
  package: "m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5",
  parcel: "m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5v5",
  pnpm: "M3 3h5v5H3zM10 3h5v5h-5zM17 3h4v5h-4zM3 10h5v5H3zm7 0h5v5h-5zm7 0h4v5h-4zM10 17h5v4h-5zm7 0h4v4h-4z",
  receipt: "M5 2h14v20l-3-2-4 2-4-2-3 2zM8 7h8M8 11h8m-8 5 2 2 5-4",
  retry: "M3 10a9 9 0 1 1 1 7M3 3v7h7M12 7v5l3 2",
  route: "M2 3h6v6H2zM16 15h6v6h-6zM11 6h8v6m-3-3 3 3 3-3M13 18H5v-6m-3 3 3-3 3 3",
  scoop: "M5 8h14l-2 13H7zM4 8h16M8 8V5a4 4 0 0 1 8 0v3",
  server: "M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3 0h6M6 10v4m12-4v4",
  shield: "m12 2 8 3v7c0 5-8 10-8 10S4 17 4 12V5zM8 11l3 3 5-6",
  sun: "M12 8a4 4 0 1 0 0 8 4 4 0 0 0 0-8zM12 2v3m0 14v3M4.9 4.9 7 7m10 10 2.1 2.1M2 12h3m14 0h3M4.9 19.1 7 17M17 7l2.1-2.1",
  system: "M3 4h18v13H3zM8 21h8M12 17v4",
  terminal: "M4 6l5 6-5 6m7 0h9",
  upload: "M12 16V3m-5 5 5-5 5 5M4 15v6h16v-6",
  windows: "M3 4l8-1v8H3zm10-1 8-1v9h-8zM3 13h8v8l-8-1zm10 0h8v9l-8-1z",
  yarn: "M12 3a9 9 0 1 0 9 9M8 16c4-1 7-4 9-8m-8 1c3 1 5 4 5 8m-5-5c-1-3 0-5 2-6",
};

const badges: Partial<Record<IconName, string>> = { "language-en": "EN", "language-ru": "RU" };

export function resolveIcon(name: string): IconName {
  return iconNames.includes(name as IconName) ? name as IconName : "parcel";
}

export class CourierIcon extends LitElement {
  static properties = {
    name: { type: String },
    label: { type: String },
  };

  static styles = css`
    :host {
      display: inline-flex;
      width: 1.5rem;
      height: 1.5rem;
      color: currentColor;
    }
    svg { width: 100%; height: 100%; }
    text { fill: currentColor; stroke: none; font-family: var(--courier-font-mono, monospace); font-size: 6px; font-weight: 850; letter-spacing: -0.04em; text-anchor: middle; }
  `;

  name = "parcel";
  label = "";

  protected render() {
    const name = resolveIcon(this.name);
    return html`<svg
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="1.7"
      stroke-linecap="square"
      stroke-linejoin="miter"
      role=${this.label ? "img" : "presentation"}
      aria-hidden=${this.label ? "false" : "true"}
      aria-label=${this.label || undefined}
    ><path d=${paths[name]}></path><text x="12" y="13.1">${badges[name] ?? ""}</text></svg>`;
  }
}
