import { LitElement, css, html } from "lit";
import { translate, type Locale } from "../i18n";

export function progressRatio(value: number, total: number): number {
  if (!Number.isFinite(value) || !Number.isFinite(total) || total <= 0) {
    return 0;
  }
  return Math.min(1, Math.max(0, value / total));
}

export class CourierProgress extends LitElement {
  static properties = {
    value: { type: Number },
    total: { type: Number },
    label: { type: String },
    locale: { type: String },
  };

  static styles = css`
    :host {
      display: grid;
      gap: var(--courier-space-2, 0.5rem);
      color: var(--courier-text, #151714);
      font-family: var(--courier-font-sans, sans-serif);
    }
    .track {
      overflow: hidden;
      height: 0.625rem;
      border: 1px solid var(--courier-line, #c7ccc0);
      border-radius: 999px;
      background: var(--courier-surface, #fff);
    }
    .fill {
      height: 100%;
      background: var(--courier-signal, #d4ff45);
      transform-origin: left;
      transition: transform var(--courier-duration, 160ms) var(--courier-ease, ease);
    }
    output { color: var(--courier-muted, #51574d); font-family: var(--courier-font-mono, monospace); }
  `;

  value = 0;
  total = 0;
  label = "";
  locale: Locale = "en";

  protected render() {
    const ratio = progressRatio(this.value, this.total);
    const label = this.label || translate(this.locale, "progress.label");
    const total = Number.isFinite(this.total) && this.total > 0 ? this.total : 0;
    const value = Number.isFinite(this.value) ? Math.max(0, Math.min(this.value, total)) : 0;
    return html`<div
      class="track"
      role="progressbar"
      aria-label=${label}
      aria-valuemin="0"
      aria-valuemax=${total}
      aria-valuenow=${value}
    ><div class="fill" style=${`transform: scaleX(${ratio})`}></div></div>
    <output>${Math.round(ratio * 100)}%</output>`;
  }
}
