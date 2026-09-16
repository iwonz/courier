import { LitElement, css, html } from "lit";
import type { IconName } from "../icons";

export class CourierIconLink extends LitElement {
  static properties = {
    href: { type: String },
    icon: { type: String },
    label: { type: String },
    target: { type: String },
    rel: { type: String },
  };

  static styles = css`
    :host { display: inline-flex; }
    a {
      box-sizing: border-box;
      display: inline-grid;
      width: var(--courier-control-frame-size, 2.35rem);
      height: var(--courier-control-frame-size, 2.35rem);
      place-items: center;
      border: var(--courier-control-frame-border-width, 1px) solid var(--courier-control-frame-border, var(--courier-color-border, #c8cdbf));
      border-radius: var(--courier-control-frame-radius, var(--courier-radius-md, 0.625rem));
      color: var(--courier-control-frame-color, var(--courier-color-muted, #596054));
      background: var(--courier-control-frame-surface, color-mix(in srgb, var(--courier-color-field, #e7e9dc) 68%, transparent));
      box-shadow: var(--courier-control-frame-shadow, inset 0 1px 2px rgb(16 18 15 / 0.07));
      text-decoration: none;
      transition: color var(--courier-duration, 160ms) var(--courier-ease, ease), border-color var(--courier-duration, 160ms) var(--courier-ease, ease), background var(--courier-duration, 160ms) var(--courier-ease, ease), transform var(--courier-duration, 160ms) var(--courier-ease, ease);
    }
    a:hover { border-color: var(--courier-control-frame-hover-border, var(--courier-color-accent, #d4ff45)); color: var(--courier-control-frame-hover-color, var(--courier-color-text, #151714)); background: var(--courier-control-frame-hover-surface, color-mix(in srgb, var(--courier-color-surface-raised, #fff) 72%, transparent)); }
    a:active { transform: translateY(1px); }
    a:focus-visible { outline: var(--courier-control-frame-focus-width, 3px) solid var(--courier-control-frame-focus, var(--courier-beak, #ff8758)); outline-offset: 2px; }
    courier-icon { width: 1.05rem; height: 1.05rem; }
  `;

  href = "";
  icon: IconName = "github";
  label = "";
  target = "_blank";
  rel = "noopener noreferrer";

  protected render() {
    return html`<a href=${this.href} target=${this.target} rel=${this.rel} aria-label=${this.label} title=${this.label}><courier-icon name=${this.icon}></courier-icon></a>`;
  }
}
