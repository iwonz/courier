import { LitElement, css, html } from "lit";
import { controlStyles } from "../styles";

export class CourierButton extends LitElement {
  static properties = {
    disabled: { type: Boolean, reflect: true },
    type: { type: String, reflect: true },
    variant: { type: String, reflect: true },
  };

  static styles = [controlStyles, css`
    :host { display: inline-flex; }
    button { width: 100%; }
    :host([variant="primary"]) button {
      border-color: var(--courier-color-accent, #d4ff45);
      color: var(--courier-color-accent-ink, #151714);
      background: var(--courier-color-accent, #d4ff45);
      font-weight: 800;
    }
    :host([variant="primary"]) button:hover:not(:disabled) {
      background: color-mix(in srgb, var(--courier-color-accent, #d4ff45) 86%, white);
    }
  `];

  disabled = false;
  type: "button" | "submit" = "button";
  variant: "primary" | "secondary" = "secondary";

  protected render() {
    return html`<button type=${this.type} ?disabled=${this.disabled}><slot></slot></button>`;
  }
}
