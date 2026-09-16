import { LitElement, css, html } from "lit";

export interface PointerPosition {
  readonly x: number;
  readonly y: number;
}

function clamp(value: number): number {
  return Math.min(100, Math.max(0, value));
}

export function pointerPosition(bounds: Pick<DOMRect, "left" | "top" | "width" | "height">, clientX: number, clientY: number): PointerPosition {
  if (bounds.width <= 0 || bounds.height <= 0) return { x: 50, y: 50 };
  return {
    x: clamp(((clientX - bounds.left) / bounds.width) * 100),
    y: clamp(((clientY - bounds.top) / bounds.height) * 100),
  };
}

export class CourierScene extends LitElement {
  static properties = {
    eager: { type: Boolean },
    mobileSource: { type: String, attribute: "mobile-source" },
    source: { type: String },
  };

  static styles = css`
    :host { --scene-pointer-x: 72%; --scene-pointer-y: 42%; position: absolute; display: block; inset: 0; overflow: hidden; background: var(--courier-graphite-900, #151714); pointer-events: auto; isolation: isolate; }
    courier-mascot { position: absolute; inset: 0; width: 100%; height: 100%; }
    courier-mascot::part(image) { width: 100%; height: 100%; object-fit: cover; filter: saturate(0.88) contrast(1.02); }
    .base { z-index: 0; transform: none; }
    .refracted { z-index: 1; clip-path: ellipse(clamp(5rem, 12vw, 11rem) clamp(4rem, 10vw, 8.5rem) at var(--scene-pointer-x) var(--scene-pointer-y)); opacity: 0.74; transform: scale(1.018); transform-origin: var(--scene-pointer-x) var(--scene-pointer-y); filter: saturate(1.08) contrast(1.03); will-change: clip-path, transform; }
    .glow { position: absolute; z-index: 2; inset: 0; background: radial-gradient(ellipse clamp(8rem, 23vw, 20rem) clamp(6rem, 17vw, 15rem) at var(--scene-pointer-x) var(--scene-pointer-y), color-mix(in srgb, var(--courier-signal, #d4ff45) 17%, transparent), transparent 66%), radial-gradient(ellipse clamp(5rem, 12vw, 10rem) clamp(8rem, 18vw, 15rem) at calc(var(--scene-pointer-x) + 3%) calc(var(--scene-pointer-y) - 2%), rgb(255 255 255 / 0.09), transparent 72%); mix-blend-mode: screen; pointer-events: none; }
    .veil { position: absolute; z-index: 3; inset: 0; background: linear-gradient(90deg, rgb(10 12 10 / 0.14), transparent 28% 72%, rgb(10 12 10 / 0.24)), linear-gradient(180deg, color-mix(in srgb, var(--courier-color-canvas, #f3f4e9) 16%, transparent), transparent 25% 78%, rgb(10 12 10 / 0.3)); pointer-events: none; }
    @media (pointer: coarse), (hover: none) {
      :host { --scene-pointer-x: 68%; --scene-pointer-y: 40%; }
      .refracted { display: none; }
      .glow { opacity: 0.74; }
    }
    @media (prefers-reduced-motion: reduce) {
      .refracted { display: none; }
      .glow { opacity: 0.68; }
    }
  `;

  eager = false;
  mobileSource = "";
  source = "";
  private frame = 0;
  private pending?: PointerPosition;

  disconnectedCallback(): void {
    if (this.frame) globalThis.cancelAnimationFrame(this.frame);
    this.frame = 0;
    super.disconnectedCallback();
  }

  private tracksPointer(): boolean {
    return globalThis.matchMedia?.("(hover: hover) and (pointer: fine)").matches === true && !globalThis.matchMedia?.("(prefers-reduced-motion: reduce)").matches;
  }

  private move(event: PointerEvent): void {
    if (!this.tracksPointer()) return;
    this.pending = pointerPosition(this.getBoundingClientRect(), event.clientX, event.clientY);
    if (this.frame) return;
    this.frame = globalThis.requestAnimationFrame(() => {
      this.frame = 0;
      const position = this.pending;
      if (!position) return;
      this.style.setProperty("--scene-pointer-x", `${position.x}%`);
      this.style.setProperty("--scene-pointer-y", `${position.y}%`);
    });
  }

  private leave(): void {
    this.pending = undefined;
    if (this.frame) globalThis.cancelAnimationFrame(this.frame);
    this.frame = 0;
    this.style.removeProperty("--scene-pointer-x");
    this.style.removeProperty("--scene-pointer-y");
  }

  protected render() {
    return html`<courier-mascot class="base" ?eager=${this.eager} alt="" .source=${this.source} .mobileSource=${this.mobileSource} @pointermove=${this.move} @pointerleave=${this.leave}></courier-mascot><courier-mascot class="refracted" aria-hidden="true" alt="" .source=${this.source} .mobileSource=${this.mobileSource}></courier-mascot><span class="glow" aria-hidden="true"></span><span class="veil" aria-hidden="true"></span>`;
  }
}
