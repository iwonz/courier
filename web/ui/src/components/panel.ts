import { LitElement, css, html } from "lit";

export class CourierPanel extends LitElement {
  static properties = { heading: { type: String } };

  static styles = css`
    :host {
      display: block;
      border: 1px solid var(--courier-line, #c7ccc0);
      border-radius: var(--courier-radius-md, 0.5rem);
      color: var(--courier-text, #151714);
      background: var(--courier-surface, #fff);
      font-family: var(--courier-font-sans, sans-serif);
    }
    section { padding: var(--courier-space-6, 1.5rem); }
    h2 { margin: 0 0 var(--courier-space-4, 1rem); font-size: 1.125rem; }
  `;

  heading = "";

  protected render() {
    return html`<section aria-labelledby="courier-panel-heading">
      <h2 id="courier-panel-heading">${this.heading}</h2>
      <slot></slot>
    </section>`;
  }
}
