import { LitElement, css, html } from "lit";

export class CourierPanel extends LitElement {
  static properties = { heading: { type: String } };

  static styles = css`
    :host {
      display: block;
      overflow: hidden;
      border: 1px solid var(--courier-color-border, #c8cdbf);
      border-radius: var(--courier-radius-md, 0.5rem);
      color: var(--courier-color-text, #151714);
      background: var(--courier-color-surface-raised, #fff);
      font-family: var(--courier-font-sans, sans-serif);
      box-shadow: 0 1px 0 rgb(16 18 15 / 0.04);
    }
    section { padding: var(--courier-space-6, 1.5rem); }
    h2 {
      margin: 0 0 var(--courier-space-4, 1rem);
      font-family: var(--courier-font-mono, monospace);
      font-size: 0.6875rem;
      letter-spacing: 0.09em;
      text-transform: uppercase;
    }
  `;

  heading = "";

  protected render() {
    if (!this.heading) {
      return html`<section><slot></slot></section>`;
    }
    return html`<section aria-labelledby="courier-panel-heading">
      <h2 id="courier-panel-heading">${this.heading}</h2>
      <slot></slot>
    </section>`;
  }
}
