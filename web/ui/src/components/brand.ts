import { LitElement, css, html } from "lit";

export class CourierBrand extends LitElement {
  static properties = { product: { type: String } };

  static styles = css`
    :host { display: inline-flex; min-width: 0; color: var(--courier-color-text, #151714); font-family: var(--courier-font-sans, sans-serif); }
    .lockup { display: inline-flex; min-width: 0; align-items: center; gap: 0.625rem; color: inherit; }
    svg { flex: 0 0 auto; width: 2.25rem; height: 2.25rem; }
    .words { display: grid; line-height: 1; }
    strong { font-family: var(--courier-font-display, sans-serif); font-size: 1.125rem; font-weight: 850; letter-spacing: -0.04em; }
    small { margin-top: 0.25rem; color: var(--courier-color-muted, #596054); font-family: var(--courier-font-mono, monospace); font-size: 0.625rem; font-weight: 700; letter-spacing: 0.09em; text-transform: uppercase; }
  `;

  product = "";

  protected render() {
    return html`<span class="lockup">
      <svg viewBox="0 0 40 40" aria-hidden="true">
        <rect x="1" y="1" width="38" height="38" rx="4" fill="var(--courier-graphite-900, #151714)" stroke="var(--courier-color-border-strong, #8e9587)"></rect>
        <path d="M10 30c-1.5-3.2-1.4-7 .3-10.2C12.5 15.5 16.4 12 21.5 12c3.8 0 7 1.8 8.8 4.7l-2.5 3.6c1 2.8.8 6.4-.7 9.7H10z" fill="var(--courier-paper-50, #f3f4e9)"></path>
        <path d="m29 15 8 3.1-8 4z" fill="var(--courier-beak, #ff8758)"></path>
        <circle cx="24" cy="16" r="2.2" fill="var(--courier-graphite-900, #151714)"></circle>
        <circle cx="24" cy="16" r=".8" fill="var(--courier-color-accent, #d4ff45)"></circle>
        <path d="M10 25h18.2c.1 1.7-.3 3.5-1.1 5H10c-.6-1.6-.8-3.3-.7-5z" fill="var(--courier-color-accent, #d4ff45)"></path>
        <path d="M15 25v5m8-5v5" stroke="var(--courier-graphite-900, #151714)" stroke-width="1.2"></path>
      </svg>
      <span class="words"><strong>Courier</strong><small>${this.product}</small></span>
    </span>`;
  }
}

export class CourierMascot extends LitElement {
  static properties = {
    alt: { type: String },
    eager: { type: Boolean },
    source: { type: String },
  };

  static styles = css`
    :host { display: block; }
    img { display: block; width: 100%; height: auto; filter: drop-shadow(0 1.5rem 2rem rgb(16 18 15 / 0.18)); }
  `;

  alt = "";
  eager = false;
  source = "";

  protected render() {
    return html`<img part="image" src=${this.source} alt=${this.alt} decoding="async" loading=${this.eager ? "eager" : "lazy"} fetchpriority=${this.eager ? "high" : "auto"}>`;
  }
}

export class CourierStatus extends LitElement {
  static properties = { tone: { type: String, reflect: true } };

  static styles = css`
    :host {
      display: inline-flex;
      width: max-content;
      align-items: center;
      gap: 0.45rem;
      color: var(--courier-color-muted, #596054);
      font-family: var(--courier-font-mono, monospace);
      font-size: 0.6875rem;
      font-weight: 750;
      letter-spacing: 0.075em;
      text-transform: uppercase;
    }
    i { width: 0.5rem; height: 0.5rem; border: 1px solid currentColor; border-radius: 50%; background: currentColor; box-shadow: 0 0 0 3px color-mix(in srgb, currentColor 14%, transparent); }
    :host([tone="signal"]) { color: var(--courier-success, #76a51f); }
    :host([tone="warning"]) { color: var(--courier-warning, #c78300); }
    :host([tone="danger"]) { color: var(--courier-danger, #ff6b5f); }
  `;

  tone: "neutral" | "signal" | "warning" | "danger" = "neutral";

  protected render() {
    return html`<i aria-hidden="true"></i><slot></slot>`;
  }
}

export class CourierRoute extends LitElement {
  static properties = {
    source: { type: String },
    destination: { type: String },
  };

  static styles = css`
    :host { display: grid; color: var(--courier-color-text, #151714); font-family: var(--courier-font-mono, monospace); }
    .route { display: grid; grid-template-columns: minmax(0, 1fr) minmax(3rem, 0.55fr) minmax(0, 1fr); align-items: center; gap: 0.65rem; }
    .node { overflow: hidden; padding: 0.65rem 0.75rem; border: 1px solid var(--courier-color-border, #c8cdbf); border-radius: var(--courier-radius-sm, 0.25rem); background: var(--courier-color-surface, #fafbf3); font-size: 0.75rem; text-overflow: ellipsis; white-space: nowrap; }
    .line { position: relative; height: 1px; color: var(--courier-color-border-strong, #8e9587); background: currentColor; }
    .line::before { content: ""; position: absolute; top: -0.2rem; left: 0; width: 0.45rem; height: 0.45rem; border-radius: 50%; background: var(--courier-color-accent, #d4ff45); }
    .line::after { content: ""; position: absolute; top: -0.22rem; right: 0; width: 0.4rem; height: 0.4rem; border-top: 1px solid currentColor; border-right: 1px solid currentColor; transform: rotate(45deg); }
  `;

  source = "";
  destination = "";

  protected render() {
    return html`<div class="route"><span class="node">${this.source}</span><span class="line" aria-hidden="true"></span><span class="node">${this.destination}</span></div>`;
  }
}
