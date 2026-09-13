import { LitElement, css, html } from "lit";
import { controlStyles } from "../styles";

export class CourierButton extends LitElement {
  static properties = {
    disabled: { type: Boolean, reflect: true },
    variant: { type: String, reflect: true },
  };

  static styles = [controlStyles, css`
    :host { display: inline-flex; }
    button { width: 100%; }
    :host([variant="primary"]) button {
      border-color: var(--courier-signal, #d4ff45);
      color: var(--courier-signal-ink, #151714);
      background: var(--courier-signal, #d4ff45);
      font-weight: 800;
    }
  `];

  disabled = false;
  variant: "primary" | "secondary" = "secondary";

  protected render() {
    return html`<button type="button" ?disabled=${this.disabled}><slot></slot></button>`;
  }
}
