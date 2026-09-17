import { LitElement, css, html, type PropertyValues } from "lit";

export class CourierScene extends LitElement {
  static properties = { eager: { type: Boolean }, mobileSource: { type: String, attribute: "mobile-source" }, source: { type: String }, active: { state: true } };
  static styles = css`
    :host { position: absolute; display: block; inset: 0; overflow: hidden; background: var(--courier-color-canvas, #f2efe6); pointer-events: none; isolation: isolate; }
    courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; pointer-events: none; }
    courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; }
    .veil { position: absolute; inset: 0; background: linear-gradient(180deg, rgb(14 15 13 / 0.02), transparent 26% 76%, rgb(14 15 13 / 0.08)); pointer-events: none; }
  `;
  eager = false;
  mobileSource = "";
  source = "";
  private active = false;
  private proximityObserver?: IntersectionObserver;

  connectedCallback(): void { super.connectedCallback(); this.configureLoading(); }
  disconnectedCallback(): void { this.proximityObserver?.disconnect(); this.proximityObserver = undefined; super.disconnectedCallback(); }
  protected updated(changed: PropertyValues<this>): void { if (changed.has("eager") && this.eager) this.activate(); }
  private configureLoading(): void {
    if (this.active) return;
    if (this.eager) { this.activate(); return; }
    if (!globalThis.IntersectionObserver) { this.active = true; return; }
    this.proximityObserver = new IntersectionObserver((entries) => {
      if (entries.some((entry) => entry.isIntersecting)) this.activate();
    }, { rootMargin: "100% 0px", threshold: 0 });
    this.proximityObserver.observe(this);
  }
  private activate(): void { this.active = true; this.proximityObserver?.disconnect(); this.proximityObserver = undefined; }
  protected render() { return html`${this.active ? html`<courier-mascot class="base" ?eager=${this.eager} alt="" .source=${this.source} .mobileSource=${this.mobileSource}></courier-mascot>` : ""}<span class="veil" aria-hidden="true"></span>`; }
}
