import { LitElement, css, html } from "lit";

export const iconNames = ["archive", "copy", "download", "folder", "parcel", "receipt", "retry", "route", "server", "shield", "upload"] as const;
export type IconName = typeof iconNames[number];

const paths: Record<IconName, string> = {
  archive: "M3 3h18v5H3zM5 8v13h14V8M9 12h6",
  copy: "M8 3h13v13M3 8h13v13H3z",
  download: "M12 3v13m-5-5 5 5 5-5M4 15v6h16v-6",
  folder: "M3 6h7l2 3h9v12H3zM3 6V3h7l2 3h9v3",
  parcel: "m3 7 9-5 9 5v11l-9 4-9-4zM3 7l9 5 9-5M12 12v10M8 4l9 5v5",
  receipt: "M5 2h14v20l-3-2-4 2-4-2-3 2zM8 7h8M8 11h8m-8 5 2 2 5-4",
  retry: "M3 10a9 9 0 1 1 1 7M3 3v7h7M12 7v5l3 2",
  route: "M2 3h6v6H2zM16 15h6v6h-6zM11 6h8v6m-3-3 3 3 3-3M13 18H5v-6m-3 3 3-3 3 3",
  server: "M3 2h18v8H3zM3 14h18v8H3zM7 6h1m3 0h6M7 18h1m3 0h6M6 10v4m12-4v4",
  shield: "m12 2 8 3v7c0 5-8 10-8 10S4 17 4 12V5zM8 11l3 3 5-6",
  upload: "M12 16V3m-5 5 5-5 5 5M4 15v6h16v-6",
};

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
    ><path d=${paths[name]}></path></svg>`;
  }
}
