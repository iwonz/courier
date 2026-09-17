import { LitElement, css, html } from "lit";
import { courierMarkSource } from "../identity-assets";
import { cubicBezierPath } from "../geometry";

export { courierMarkSource };

export class CourierBrand extends LitElement {
  static styles = css`
    :host { display: inline-flex; min-width: 0; color: var(--courier-color-text, #0e0f0d); font-family: var(--courier-font-sans, sans-serif); }
    .lockup { display: inline-flex; min-width: 0; align-items: center; gap: 0.58rem; color: inherit; }
    img { flex: 0 0 auto; width: 2.45rem; height: 2.1rem; object-fit: contain; filter: drop-shadow(0 0.32rem 0.55rem rgb(14 15 13 / 0.16)); }
    strong { font-family: var(--courier-font-display, sans-serif); font-size: 1.08rem; font-weight: 820; letter-spacing: -0.045em; }
  `;

  protected render() {
    return html`<span class="lockup"><img part="mark" src=${courierMarkSource} alt="" decoding="async"><strong part="wordmark">Courier</strong></span>`;
  }
}

export class CourierMascot extends LitElement {
  static properties = { alt: { type: String }, eager: { type: Boolean }, mobileSource: { type: String, attribute: "mobile-source" }, source: { type: String } };
  static styles = css`:host { display: block; } picture { display: contents; } img { display: block; width: 100%; height: auto; }`;
  alt = "";
  eager = false;
  mobileSource = "";
  source = "";
  protected render() {
    return html`<picture>${this.mobileSource ? html`<source media="(max-width: 44rem)" srcset=${this.mobileSource}>` : ""}<img part="image" src=${this.source} alt=${this.alt} decoding="async" loading=${this.eager ? "eager" : "lazy"} fetchpriority=${this.eager ? "high" : "auto"}></picture>`;
  }
}

export class CourierStatus extends LitElement {
  static properties = { tone: { type: String, reflect: true } };
  static styles = css`
    :host { display: inline-flex; width: max-content; align-items: center; gap: 0.45rem; color: var(--courier-color-muted, #68645e); font-family: var(--courier-font-mono, monospace); font-size: 0.6875rem; font-weight: 700; }
    i { width: 0.48rem; height: 0.48rem; border: 1px solid currentColor; border-radius: 50%; background: currentColor; }
    :host([tone="signal"]) { color: var(--courier-success, #39765d); }
    :host([tone="warning"]) { color: var(--courier-warning, #a86618); }
    :host([tone="danger"]) { color: var(--courier-danger, #b7443e); }
  `;
  tone: "neutral" | "signal" | "warning" | "danger" = "neutral";
  protected render() { return html`<i aria-hidden="true"></i><slot></slot>`; }
}

export class CourierRoute extends LitElement {
  static properties = { source: { type: String }, destination: { type: String } };
  static styles = css`
    :host { display: grid; color: var(--courier-color-text, #0e0f0d); font-family: var(--courier-font-mono, monospace); }
    .route { display: grid; grid-template-columns: minmax(0, 1fr) minmax(3rem, 0.55fr) minmax(0, 1fr); align-items: center; gap: 0.65rem; }
    .node { overflow: hidden; padding: 0.62rem 0; border-bottom: 1px solid var(--courier-color-border, #cbc5b8); font-size: 0.75rem; text-overflow: ellipsis; white-space: nowrap; }
    .node:last-child { text-align: right; }
    .connector { position: relative; height: 2rem; }
    svg { position: absolute; inset: 0; width: 100%; height: 100%; overflow: visible; color: var(--courier-color-border-strong, #8f8a81); }
    path { fill: none; stroke: currentColor; stroke-width: 1.5; vector-effect: non-scaling-stroke; }
    .terminal { position: absolute; top: 50%; width: 0.65rem; height: 0.65rem; aspect-ratio: 1; border: 2px solid var(--courier-color-accent, #ad431d); border-radius: 50%; background: var(--courier-color-canvas, #f2efe6); transform: translateY(-50%); }
    .terminal.source { left: 0; transform: translate(-50%, -50%); }
    .terminal.destination { right: 0; transform: translate(50%, -50%); }
  `;
  source = "";
  destination = "";
  protected render() {
    const path = cubicBezierPath({ x: 1, y: 15 }, { x: 99, y: 15 });
    return html`<div class="route"><span class="node">${this.source}</span><span class="connector" aria-hidden="true"><svg viewBox="0 0 100 30" preserveAspectRatio="none"><path d=${path}></path></svg><i class="terminal source"></i><i class="terminal destination"></i></span><span class="node">${this.destination}</span></div>`;
  }
}
